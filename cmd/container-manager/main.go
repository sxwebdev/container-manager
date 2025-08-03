package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"container-manager/internal/logger"

	"gopkg.in/yaml.v3"
)

var (
	appName    = "container-manager"
	version    = "local"
	commitHash = "unknown"
	buildDate  = "unknown"
)

var appLogger logger.Logger

type ServiceConfig struct {
	Name        string `yaml:"name"`
	ProjectPath string `yaml:"project_path"`
	ComposeFile string `yaml:"compose_file"`
	Enabled     bool   `yaml:"enabled"`
}

type Config struct {
	BindAddress string          `yaml:"bind_address"`
	Port        string          `yaml:"port"`
	Services    []ServiceConfig `yaml:"services"`
}

type UpdateRequest struct {
	Service string `json:"service"`
	Action  string `json:"action"`
	Target  string `json:"target,omitempty"`
}

type UpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output,omitempty"`
	Error   string `json:"error,omitempty"`
	Service string `json:"service"`
}

type ListServicesResponse struct {
	Services []ServiceInfo `json:"services"`
}

type ServiceInfo struct {
	Name        string `json:"name"`
	ProjectPath string `json:"project_path"`
	ComposeFile string `json:"compose_file"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
}

var config Config

func loadConfig() {
	// Default configuration values
	config = Config{
		BindAddress: "127.0.0.1",
		Port:        "8090",
		Services:    []ServiceConfig{},
	}

	// Load from YAML file
	configPaths := []string{
		"./config.yaml",
		"./container-manager.yaml",
	}

	for _, configPath := range configPaths {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}

		if data, err := os.ReadFile(configPath); err == nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				appLogger.Errorf("Error parsing config file %s: %v", configPath, err)
				continue
			}
			appLogger.Infof("Loaded config from: %s", configPath)
			break
		}
	}

	// Environment variables override config
	if addr := os.Getenv("CONTAINER_MANAGER_BIND_ADDRESS"); addr != "" {
		config.BindAddress = addr
	}
	if port := os.Getenv("CONTAINER_MANAGER_PORT"); port != "" {
		config.Port = port
	}
}

func findService(serviceName string) (*ServiceConfig, error) {
	for _, service := range config.Services {
		if service.Name == serviceName && service.Enabled {
			return &service, nil
		}
	}
	return nil, fmt.Errorf("service '%s' not found or disabled", serviceName)
}

func getServiceStatus(service *ServiceConfig) string {
	cmd := exec.Command("docker", "compose", "-f", service.ComposeFile, "ps", "--format", "json")
	cmd.Dir = service.ProjectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "error"
	}

	if strings.TrimSpace(string(output)) == "" || string(output) == "[]\n" {
		return "stopped"
	}

	return "running"
}

func executeDockerCommand(serviceName, action, target string) UpdateResponse {
	service, err := findService(serviceName)
	if err != nil {
		return UpdateResponse{
			Success: false,
			Error:   err.Error(),
			Service: serviceName,
		}
	}

	appLogger.Infow("Executing action",
		"action", action,
		"service", serviceName,
		"target", target)

	var args []string

	switch action {
	case "pull":
		args = []string{"compose", "-f", service.ComposeFile, "pull"}
		if target != "" {
			args = append(args, target)
		}
	case "up":
		args = []string{"compose", "-f", service.ComposeFile, "up", "-d"}
		if target != "" {
			args = append(args, target)
		}
	case "down":
		args = []string{"compose", "-f", service.ComposeFile, "down"}
		if target != "" {
			args = append(args, target)
		}
	case "restart":
		args = []string{"compose", "-f", service.ComposeFile, "restart"}
		if target != "" {
			args = append(args, target)
		}
	case "stop":
		args = []string{"compose", "-f", service.ComposeFile, "stop"}
		if target != "" {
			args = append(args, target)
		}
	case "logs":
		args = []string{"compose", "-f", service.ComposeFile, "logs", "--tail=50"}
		if target != "" {
			args = append(args, target)
		}
	case "ps":
		args = []string{"compose", "-f", service.ComposeFile, "ps"}
	case "update":
		// Комбинированная команда: pull + up
		pullResp := executeDockerCommand(serviceName, "pull", target)
		if !pullResp.Success {
			return pullResp
		}
		return executeDockerCommand(serviceName, "up", target)
	default:
		return UpdateResponse{
			Success: false,
			Error:   fmt.Sprintf("Unknown action: %s", action),
			Service: serviceName,
		}
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = service.ProjectPath

	appLogger.Infow("Running docker command",
		"command", strings.Join(args, " "),
		"directory", service.ProjectPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		appLogger.Errorw("Docker command failed",
			"error", err,
			"output", string(output),
			"service", serviceName)
		return UpdateResponse{
			Success: false,
			Message: "Docker command failed",
			Output:  string(output),
			Error:   err.Error(),
			Service: serviceName,
		}
	}

	appLogger.Infow("Command completed successfully", "service", serviceName)
	return UpdateResponse{
		Success: true,
		Message: fmt.Sprintf("Action '%s' completed successfully for service '%s'", action, serviceName),
		Output:  string(output),
		Service: serviceName,
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	serviceStatuses := make([]ServiceInfo, 0, len(config.Services))
	for _, service := range config.Services {
		if service.Enabled {
			serviceStatuses = append(serviceStatuses, ServiceInfo{
				Name:        service.Name,
				ProjectPath: service.ProjectPath,
				ComposeFile: service.ComposeFile,
				Enabled:     service.Enabled,
				Status:      getServiceStatus(&service),
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"services":  serviceStatuses,
	})
}

func listServicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	services := make([]ServiceInfo, 0, len(config.Services))
	for _, service := range config.Services {
		services = append(services, ServiceInfo{
			Name:        service.Name,
			ProjectPath: service.ProjectPath,
			ComposeFile: service.ComposeFile,
			Enabled:     service.Enabled,
			Status:      getServiceStatus(&service),
		})
	}

	json.NewEncoder(w).Encode(ListServicesResponse{Services: services})
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdateRequest

	switch r.Method {
	case "GET":
		serviceName := r.URL.Query().Get("service")
		if serviceName == "" {
			serviceName = "default"
		}
		req = UpdateRequest{
			Service: serviceName,
			Action:  r.URL.Query().Get("action"),
			Target:  r.URL.Query().Get("target"),
		}
		if req.Action == "" {
			req.Action = "update"
		}
	case "POST":
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			body, _ := io.ReadAll(r.Body)
			req = UpdateRequest{
				Service: "default",
				Action:  strings.TrimSpace(string(body)),
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(UpdateResponse{
			Success: false,
			Error:   "Only GET and POST methods allowed",
		})
		return
	}

	result := executeDockerCommand(req.Service, req.Action, req.Target)

	if !result.Success {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	appLogger = logger.New(logger.Config{
		Level:  slog.LevelInfo,
		Format: "json", // можно поменять на "text" для более читаемого формата
		ExtraAttrs: map[string]any{
			"service": appName,
			"version": version,
			"build":   buildDate,
			"commit":  commitHash,
		},
	})

	loadConfig()

	bindAddr := fmt.Sprintf("%s:%s", config.BindAddress, config.Port)
	appLogger.Infof("Starting container-manager on %s", bindAddr)
	appLogger.Infow("Configured services", "count", len(config.Services))

	for _, service := range config.Services {
		if !service.Enabled {
			continue
		}

		if _, err := os.Stat(service.ProjectPath); os.IsNotExist(err) {
			appLogger.Warnw("Project directory does not exist for service",
				"service", service.Name,
				"path", service.ProjectPath)
			continue
		}

		composeFullPath := filepath.Join(service.ProjectPath, service.ComposeFile)
		if _, err := os.Stat(composeFullPath); os.IsNotExist(err) {
			appLogger.Warnw("Compose file does not exist for service",
				"service", service.Name,
				"compose_path", composeFullPath)
			continue
		}

		appLogger.Infow("Service configured",
			"service", service.Name,
			"project_path", service.ProjectPath,
			"compose_file", service.ComposeFile)
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/services", listServicesHandler)
	http.HandleFunc("/service", serviceHandler)

	server := &http.Server{
		Addr:         bindAddr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	appLogger.Infof("Server ready at http://%s", bindAddr)
	if err := server.ListenAndServe(); err != nil {
		appLogger.Fatal("Server failed to start", "error", err)
	}
}
