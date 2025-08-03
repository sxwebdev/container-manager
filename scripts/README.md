# Container Manager Deployment Scripts

This directory contains scripts for automated installation and upgrades of Container Manager system.

## Scripts Overview

### install.sh

Automated installation script for Linux systems with systemd support.

**Features:**

- Platform detection (Linux amd64/arm64)
- Automatic binary download from GitHub releases
- System user creation with security hardening
- Directory structure setup with proper permissions
- Systemd service configuration with security features
- Configuration file creation with sensible defaults
- Service startup and health verification
- Docker integration for container management

**Usage:**

```bash
# Basic installation
curl -L https://raw.githubusercontent.com/sxwebdev/container-manager/master/scripts/install.sh | sudo bash

# Custom installation directories
sudo ./install.sh --install-dir /opt/container-manager/bin --config-dir /etc/container-manager

# Install specific version
sudo ./install.sh --version v1.0.0

# Custom service configuration
sudo ./install.sh --service-name container-manager-prod --user cmuser
```

**Options:**

- `-d, --install-dir DIR`: Binary installation directory (default: `/usr/local/bin`)
- `-c, --config-dir DIR`: Configuration directory (default: `/etc/container-manager`)
- `-D, --data-dir DIR`: Data directory (default: `/var/lib/container-manager`)
- `-s, --service-name NAME`: Systemd service name (default: `container-manager`)
- `-u, --user USER`: Service user (default: `container-manager`)
- `-v, --version VERSION`: Install specific version (default: latest)
- `-h, --help`: Show help message

**Post-Installation:**

- Service: `systemctl status container-manager`
- Configuration: `/etc/container-manager/config.yaml`
- Logs: `journalctl -u container-manager -f`
- Web interface: <http://localhost:8090>

### upgrade.sh

Automated upgrade script for existing Container Manager installations.

**Features:**

- Automatic latest version detection
- Service backup before upgrade
- Graceful service shutdown and restart
- Configuration preservation
- Health checks after upgrade
- Automatic rollback on failure
- Minimal downtime (10-30 seconds)

**Usage:**

```bash
# Upgrade to latest version
curl -L https://raw.githubusercontent.com/sxwebdev/container-manager/master/scripts/upgrade.sh | sudo bash

# Upgrade specific service
sudo ./upgrade.sh container-manager

# Upgrade to specific version
sudo ./upgrade.sh --version v1.0.0

# Dry run (check what would be updated)
sudo ./upgrade.sh --dry-run

# Force upgrade (skip confirmations)
sudo ./upgrade.sh --force
```

**Options:**

- `SERVICE_NAME`: Target service name (positional argument, default: `container-manager`)
- `--version VERSION`: Upgrade to specific version
- `--dry-run`: Show what would be upgraded without making changes
- `--force`: Skip confirmation prompts
- `--help`: Show help message

## Requirements

### System Requirements

- Linux operating system (Ubuntu, CentOS, RHEL, Debian, etc.)
- systemd init system
- Root or sudo privileges
- Internet connection for downloads

### Dependencies

Both scripts automatically check for and require these tools:

- `curl`: For downloading files
- `jq`: For JSON processing
- `tar`: For archive extraction
- `systemctl`: For service management
- `docker`: For container management (runtime dependency)

### Platform Support

- **Architecture**: amd64 (x86_64), arm64 (aarch64)
- **Operating System**: Linux distributions with systemd and Docker
- **Not Supported**: macOS (no systemd), Windows

## Best Practices

### Production Deployment

1. **Docker Access**: Ensure the container-manager user has access to Docker socket
1. **Project Permissions**: Configure proper permissions for Docker Compose project directories
1. **Network Security**: Use firewall rules to restrict access to management interface
1. **Regular Backups**: Backup configuration and Docker volumes regularly
1. **Monitor Resources**: Set up monitoring for both Container Manager and managed containers

### Maintenance

1. **Regular Updates**: Use upgrade script monthly or when security updates are available
1. **Health Monitoring**: Include Container Manager service in your monitoring stack
1. **Configuration Review**: Periodically review and update service configurations
1. **Docker Maintenance**: Regular Docker cleanup and maintenance of managed containers

### Security

1. **User Permissions**: Never run Container Manager as root in production
1. **Docker Security**: Follow Docker security best practices for managed containers
1. **Network Security**: Use firewall rules to restrict access
1. **Regular Audits**: Review service logs and container access patterns
1. **Update Management**: Apply security updates promptly

## Configuration

### Basic Configuration

The default configuration file (`/etc/container-manager/config.yaml`) contains:

```yaml
bind_address: "127.0.0.1"
port: "8090"

services:
  - name: "example-service"
    project_path: "/opt/docker-projects/example"
    compose_file: "docker-compose.yml"
    enabled: true
```

### Service Configuration

Each service in the configuration should specify:

- `name`: Unique identifier for the service
- `project_path`: Path to the directory containing Docker Compose files
- `compose_file`: Name of the Docker Compose file (default: docker-compose.yml)
- `enabled`: Whether the service should be available for management

### Docker Integration

Container Manager requires:

- Docker daemon running and accessible
- User permissions to execute Docker commands
- Access to Docker Compose project directories

## License

These scripts are part of the Container Manager project and are licensed under the MIT License.
