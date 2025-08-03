# Container Manager

A lightweight web-based tool for managing Docker Compose services with a simple HTTP API and web interface.

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Required-blue.svg)](https://docker.com)

## Features

- 🐳 **Docker Compose Management**: Start, stop, restart, and rebuild Docker Compose services
- 🌐 **Web Interface**: Simple HTTP API for service management
- ⚙️ **Configuration-driven**: YAML-based configuration for easy service setup
- 📊 **Service Status**: Real-time status monitoring of configured services
- 🔧 **Multiple Actions**: Support for up, down, restart, build, and pull operations
- 📝 **Structured Logging**: Built-in structured logging with JSON and text formats
- 🚀 **Lightweight**: Single binary deployment with minimal dependencies

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Linux/macOS/Windows system
- Go 1.24+ (for building from source)

### Installation

#### Option 1: Using Installation Script (Linux only)

```bash
# Install latest version
curl -L https://raw.githubusercontent.com/sxwebdev/container-manager/master/scripts/install.sh | sudo bash

# Custom installation
sudo ./scripts/install.sh --install-dir /opt/container-manager/bin --config-dir /etc/container-manager
```

#### Option 2: Download Binary

Download the latest release from [GitHub Releases](https://github.com/sxwebdev/container-manager/releases)

#### Option 3: Build from Source

```bash
git clone https://github.com/sxwebdev/container-manager.git
cd container-manager
go build -o container-manager ./cmd/container-manager
```

### Configuration

Create a `config.yaml` file:

```yaml
bind_address: "127.0.0.1"
port: "8090"

services:
  - name: "web-app"
    project_path: "/path/to/your/docker-project"
    compose_file: "docker-compose.yml"
    enabled: true

  - name: "database"
    project_path: "/path/to/database/project"
    compose_file: "docker-compose.prod.yml"
    enabled: true

  - name: "monitoring"
    project_path: "/opt/monitoring"
    compose_file: "docker-compose.yml"
    enabled: false
```

### Running

```bash
# Run with default config file (./config.yaml)
./container-manager

# Run with custom config file
./container-manager --config /path/to/config.yaml

# Run as systemd service (after installation)
sudo systemctl start container-manager
sudo systemctl enable container-manager
```

## API Usage

### List Services

```bash
curl http://localhost:8090/services
```

Response:

```json
{
  "services": [
    {
      "name": "web-app",
      "project_path": "/path/to/your/docker-project",
      "compose_file": "docker-compose.yml",
      "enabled": true
    }
  ]
}
```

### Manage Service

```bash
# Start service
curl -X POST "http://localhost:8090/service?action=up&service=web-app"

# Stop service
curl -X POST "http://localhost:8090/service?action=down&service=web-app"

# Restart service
curl -X POST "http://localhost:8090/service?action=restart&service=web-app"

# Rebuild service
curl -X POST "http://localhost:8090/service?action=build&service=web-app"

# Pull latest images
curl -X POST "http://localhost:8090/service?action=pull&service=web-app"

# Target specific container
curl -X POST "http://localhost:8090/service?action=up&service=web-app&target=api"
```

Response:

```json
{
  "success": true,
  "message": "Docker command executed successfully",
  "output": "Container web-app_api_1 started",
  "service": "web-app"
}
```

### Health Check

```bash
curl http://localhost:8090/health
```

Response:

```json
{
  "status": "ok",
  "timestamp": "2025-08-04T02:00:00Z"
}
```

## Configuration Reference

### Global Settings

| Parameter      | Type   | Default       | Description                     |
| -------------- | ------ | ------------- | ------------------------------- |
| `bind_address` | string | `"127.0.0.1"` | IP address to bind the server   |
| `port`         | string | `"8090"`      | Port number for the HTTP server |

### Service Configuration

| Parameter      | Type    | Required | Description                                                     |
| -------------- | ------- | -------- | --------------------------------------------------------------- |
| `name`         | string  | Yes      | Unique identifier for the service                               |
| `project_path` | string  | Yes      | Path to the Docker Compose project directory                    |
| `compose_file` | string  | No       | Docker Compose file name (default: "docker-compose.yml")        |
| `enabled`      | boolean | No       | Whether the service is available for management (default: true) |

## Supported Actions

| Action    | Description     | Docker Compose Command            |
| --------- | --------------- | --------------------------------- |
| `up`      | Start service   | `docker compose up -d [target]`   |
| `down`    | Stop service    | `docker compose down [target]`    |
| `restart` | Restart service | `docker compose restart [target]` |
| `build`   | Build service   | `docker compose build [target]`   |
| `pull`    | Pull images     | `docker compose pull [target]`    |

## Deployment

### Production Deployment

1. **Create dedicated user**:

   ```bash
   sudo useradd -r -s /bin/false container-manager
   ```

2. **Set up directories**:

   ```bash
   sudo mkdir -p /etc/container-manager /var/lib/container-manager
   sudo chown container-manager:container-manager /var/lib/container-manager
   ```

3. **Configure systemd service**: Use the provided installation script or create manually

4. **Security considerations**:
   - Run as non-root user
   - Ensure Docker socket access for the service user
   - Use firewall rules to restrict network access
   - Regular security updates

### Docker Group Permissions

The service user needs access to Docker:

```bash
sudo usermod -aG docker container-manager
```

### Systemd Service

Example systemd service file:

```ini
[Unit]
Description=Container Manager Service
Documentation=https://github.com/sxwebdev/container-manager
After=network.target docker.service
Wants=network.target
Requires=docker.service

[Service]
Type=simple
User=container-manager
Group=container-manager
ExecStart=/usr/local/bin/container-manager
WorkingDirectory=/etc/container-manager
Environment=PATH=/usr/local/bin:/usr/bin:/bin
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Development

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### Logs

View service logs:

```bash
# Systemd logs
journalctl -u container-manager -f

# Check service status
systemctl status container-manager

# Application logs (if running directly)
./container-manager
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📝 [Documentation](https://github.com/sxwebdev/container-manager/wiki)
- 🐛 [Issue Tracker](https://github.com/sxwebdev/container-manager/issues)
- 💬 [Discussions](https://github.com/sxwebdev/container-manager/discussions)
