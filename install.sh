#!/bin/bash
## Velo Installation Script
## This script installs Velo, a lightweight, self-hostable deployment and operations platform built on Docker Swarm

set -e # Exit immediately if a command exits with a non-zero status
set -o pipefail # Cause a pipeline to return the status of the last command that exited with a non-zero status

# Environment variables that can be set:
# DOCKER_ADDRESS_POOL_BASE - Custom Docker address pool base (default: 10.0.0.0/8)
# DOCKER_ADDRESS_POOL_SIZE - Custom Docker address pool size (default: 24)
# DOCKER_POOL_FORCE_OVERRIDE - Force override Docker address pool configuration (default: false)
# REGISTRY_URL - Custom registry URL for Docker images (default: ghcr.io)

# Constants
DATE=$(date +"%Y%m%d-%H%M%S")
OS_TYPE=$(grep -w "ID" /etc/os-release | cut -d "=" -f 2 | tr -d '"')
DOCKER_VERSION="27.0"
CURRENT_USER=$USER
VELO_DATA_DIR="/opt/velo"
VELO_CONFIG_DIR="$VELO_DATA_DIR/config"
VELO_LOG_DIR="$VELO_DATA_DIR/logs"
INSTALLATION_LOG_WITH_DATE="$VELO_DATA_DIR/installation-${DATE}.log"

# Check if running as root
if [ $EUID != 0 ]; then
    echo "Please run this script as root or with sudo"
    exit 1
fi

echo -e "Welcome to Velo Installer!"
echo -e "This script will install everything for you. Sit back and relax."
echo -e "Source code: https://github.com/jasonlovesdoggo/velo\n"

# Create necessary directories
mkdir -p $VELO_DATA_DIR
mkdir -p $VELO_CONFIG_DIR
mkdir -p $VELO_LOG_DIR

# Set up logging
mkdir -p "$(dirname "$INSTALLATION_LOG_WITH_DATE")"
exec > >(tee -a "$INSTALLATION_LOG_WITH_DATE") 2>&1

# Docker address pool configuration defaults
DOCKER_ADDRESS_POOL_BASE_DEFAULT="10.0.0.0/8"
DOCKER_ADDRESS_POOL_SIZE_DEFAULT=24

# Check if environment variables were explicitly provided
DOCKER_POOL_BASE_PROVIDED=false
DOCKER_POOL_SIZE_PROVIDED=false
DOCKER_POOL_FORCE_OVERRIDE=${DOCKER_POOL_FORCE_OVERRIDE:-false}

if [ -n "${DOCKER_ADDRESS_POOL_BASE+x}" ]; then
    DOCKER_POOL_BASE_PROVIDED=true
fi

if [ -n "${DOCKER_ADDRESS_POOL_SIZE+x}" ]; then
    DOCKER_POOL_SIZE_PROVIDED=true
fi

# Registry URL configuration
if [ -n "${REGISTRY_URL+x}" ]; then
    echo "Using registry URL from environment variable: $REGISTRY_URL"
else
    REGISTRY_URL="ghcr.io"
    echo "Using default registry URL: $REGISTRY_URL"
fi

# Function to restart Docker service
restart_docker_service() {
    # Check if systemctl is available
    if command -v systemctl >/dev/null 2>&1; then
        systemctl restart docker
        if [ $? -eq 0 ]; then
            echo " - Docker daemon restarted successfully"
        else
            echo " - Failed to restart Docker daemon"
            return 1
        fi
    # Check if service command is available
    elif command -v service >/dev/null 2>&1; then
        service docker restart
        if [ $? -eq 0 ]; then
            echo " - Docker daemon restarted successfully"
        else
            echo " - Failed to restart Docker daemon"
            return 1
        fi
    # If neither systemctl nor service is available
    else
        echo " - Error: No service management system found"
        return 1
    fi
}

# Docker address pool configuration
DOCKER_ADDRESS_POOL_BASE=${DOCKER_ADDRESS_POOL_BASE:-"$DOCKER_ADDRESS_POOL_BASE_DEFAULT"}
DOCKER_ADDRESS_POOL_SIZE=${DOCKER_ADDRESS_POOL_SIZE:-$DOCKER_ADDRESS_POOL_SIZE_DEFAULT}

# Check if daemon.json exists and extract existing address pool configuration
EXISTING_POOL_CONFIGURED=false
if [ -f /etc/docker/daemon.json ]; then
    if command -v jq >/dev/null 2>&1 && jq -e '.["default-address-pools"]' /etc/docker/daemon.json >/dev/null 2>&1; then
        EXISTING_POOL_BASE=$(jq -r '.["default-address-pools"][0].base' /etc/docker/daemon.json 2>/dev/null || true)
        EXISTING_POOL_SIZE=$(jq -r '.["default-address-pools"][0].size' /etc/docker/daemon.json 2>/dev/null || true)

        if [ -n "$EXISTING_POOL_BASE" ] && [ -n "$EXISTING_POOL_SIZE" ] && [ "$EXISTING_POOL_BASE" != "null" ] && [ "$EXISTING_POOL_SIZE" != "null" ]; then
            echo "Found existing Docker network pool: $EXISTING_POOL_BASE/$EXISTING_POOL_SIZE"
            EXISTING_POOL_CONFIGURED=true

            # Check if environment variables were explicitly provided
            if [ "$DOCKER_POOL_BASE_PROVIDED" = false ] && [ "$DOCKER_POOL_SIZE_PROVIDED" = false ]; then
                DOCKER_ADDRESS_POOL_BASE="$EXISTING_POOL_BASE"
                DOCKER_ADDRESS_POOL_SIZE="$EXISTING_POOL_SIZE"
            else
                # Check if force override is enabled
                if [ "$DOCKER_POOL_FORCE_OVERRIDE" = true ]; then
                    echo "Force override enabled - network pool will be updated with $DOCKER_ADDRESS_POOL_BASE/$DOCKER_ADDRESS_POOL_SIZE."
                else
                    echo "Custom pool provided but force override not enabled - using existing configuration."
                    echo "To force override, set DOCKER_POOL_FORCE_OVERRIDE=true"
                    DOCKER_ADDRESS_POOL_BASE="$EXISTING_POOL_BASE"
                    DOCKER_ADDRESS_POOL_SIZE="$EXISTING_POOL_SIZE"
                    DOCKER_POOL_BASE_PROVIDED=false
                    DOCKER_POOL_SIZE_PROVIDED=false
                fi
            fi
        fi
    fi
fi

echo -e "---------------------------------------------"
echo "| Operating System  | $OS_TYPE"
echo "| Docker            | $DOCKER_VERSION"
echo "| Docker Pool       | $DOCKER_ADDRESS_POOL_BASE (size $DOCKER_ADDRESS_POOL_SIZE)"
echo "| Registry URL      | $REGISTRY_URL"
echo -e "---------------------------------------------\n"

echo -e "1. Installing required packages (curl, wget, git, jq, openssl, rsync). "

# Install required packages based on OS type
case "$OS_TYPE" in
arch)
    pacman -Sy --noconfirm --needed curl wget git jq openssl rsync >/dev/null || true
    ;;
alpine)
    sed -i '/^#.*\/community/s/^#//' /etc/apk/repositories
    apk update >/dev/null
    apk add curl wget git jq openssl rsync >/dev/null
    ;;
ubuntu | debian | raspbian)
    apt-get update -y >/dev/null
    apt-get install -y curl wget git jq openssl rsync >/dev/null
    ;;
centos | fedora | rhel | ol | rocky | almalinux | amzn)
    if [ "$OS_TYPE" = "amzn" ]; then
        dnf install -y wget git jq openssl rsync >/dev/null
    else
        if ! command -v dnf >/dev/null; then
            yum install -y dnf >/dev/null
        fi
        if ! command -v curl >/dev/null; then
            dnf install -y curl >/dev/null
        fi
        dnf install -y wget git jq openssl rsync >/dev/null
    fi
    ;;
sles | opensuse-leap | opensuse-tumbleweed)
    zypper refresh >/dev/null
    zypper install -y curl wget git jq openssl rsync >/dev/null
    ;;
*)
    echo "This script only supports Debian, Redhat, Arch Linux, or SLES based operating systems for now."
    exit 1
    ;;
esac

echo -e "2. Check Docker Installation. "

# Function to install Docker
install_docker() {
    curl -s https://releases.rancher.com/install-docker/${DOCKER_VERSION}.sh | sh 2>&1 || true
    if ! [ -x "$(command -v docker)" ]; then
        curl -s https://get.docker.com | sh -s -- --version ${DOCKER_VERSION} 2>&1
        if ! [ -x "$(command -v docker)" ]; then
            echo " - Docker installation failed."
            echo "   Maybe your OS is not supported?"
            echo " - Please visit https://docs.docker.com/engine/install/ and install Docker manually to continue."
            exit 1
        fi
    fi
}

# Install Docker if not already installed
if ! [ -x "$(command -v docker)" ]; then
    echo " - Docker is not installed. Installing Docker. It may take a while."
    case "$OS_TYPE" in
    "almalinux")
        dnf config-manager --add-repo=https://download.docker.com/linux/centos/docker-ce.repo >/dev/null 2>&1
        dnf install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin >/dev/null 2>&1
        if ! [ -x "$(command -v docker)" ]; then
            echo " - Docker could not be installed automatically. Please visit https://docs.docker.com/engine/install/ and install Docker manually to continue."
            exit 1
        fi
        systemctl start docker >/dev/null 2>&1
        systemctl enable docker >/dev/null 2>&1
        ;;
    "alpine")
        apk add docker docker-cli-compose >/dev/null 2>&1
        rc-update add docker default >/dev/null 2>&1
        service docker start >/dev/null 2>&1
        if ! [ -x "$(command -v docker)" ]; then
            echo " - Failed to install Docker with apk. Try to install it manually."
            echo "   Please visit https://wiki.alpinelinux.org/wiki/Docker for more information."
            exit 1
        fi
        ;;
    "arch")
        pacman -Sy docker docker-compose --noconfirm >/dev/null 2>&1
        systemctl enable docker.service >/dev/null 2>&1
        if ! [ -x "$(command -v docker)" ]; then
            echo " - Failed to install Docker with pacman. Try to install it manually."
            echo "   Please visit https://wiki.archlinux.org/title/docker for more information."
            exit 1
        fi
        ;;
    "amzn")
        dnf install docker -y >/dev/null 2>&1
        DOCKER_CONFIG=${DOCKER_CONFIG:-/usr/local/lib/docker}
        mkdir -p "$DOCKER_CONFIG"/cli-plugins >/dev/null 2>&1
        curl -sL "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o $DOCKER_CONFIG/cli-plugins/docker-compose >/dev/null 2>&1
        chmod +x "$DOCKER_CONFIG"/cli-plugins/docker-compose >/dev/null 2>&1
        systemctl start docker >/dev/null 2>&1
        systemctl enable docker >/dev/null 2>&1
        if ! [ -x "$(command -v docker)" ]; then
            echo " - Failed to install Docker with dnf. Try to install it manually."
            exit 1
        fi
        ;;
    "centos" | "fedora" | "rhel")
        if [ -x "$(command -v dnf5)" ]; then
            # dnf5 is available
            dnf config-manager addrepo --from-repofile=https://download.docker.com/linux/"$OS_TYPE"/docker-ce.repo --overwrite >/dev/null 2>&1
        else
            # dnf5 is not available, use dnf
            # shellcheck disable=SC2086
            dnf config-manager --add-repo=https://download.docker.com/linux/$OS_TYPE/docker-ce.repo >/dev/null 2>&1
        fi
        dnf install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin >/dev/null 2>&1
        if ! [ -x "$(command -v docker)" ]; then
            echo " - Docker could not be installed automatically. Please visit https://docs.docker.com/engine/install/ and install Docker manually to continue."
            exit 1
        fi
        systemctl start docker >/dev/null 2>&1
        systemctl enable docker >/dev/null 2>&1
        ;;
    "ubuntu" | "debian" | "raspbian")
        if [ "$OS_TYPE" = "ubuntu" ] && [ "$OS_VERSION" = "24.10" ]; then
            echo " - Installing Docker for Ubuntu 24.10..."
            apt-get update >/dev/null
            apt-get install -y ca-certificates curl >/dev/null
            install -m 0755 -d /etc/apt/keyrings
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
            chmod a+r /etc/apt/keyrings/docker.asc

            # Add the repository to Apt sources
            echo \
                "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
                  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" |
                tee /etc/apt/sources.list.d/docker.list >/dev/null
            apt-get update >/dev/null
            apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin >/dev/null

            if ! [ -x "$(command -v docker)" ]; then
                echo " - Docker installation failed."
                echo "   Please visit https://docs.docker.com/engine/install/ubuntu/ and install Docker manually to continue."
                exit 1
            fi
            echo " - Docker installed successfully for Ubuntu 24.10."
        else
            install_docker
        fi
        ;;
    *)
        install_docker
        ;;
    esac
    echo " - Docker installed successfully."
else
    echo " - Docker is already installed."
fi

echo -e "3. Configure Docker. "

echo " - Network pool configuration: ${DOCKER_ADDRESS_POOL_BASE}/${DOCKER_ADDRESS_POOL_SIZE}"
echo " - To override existing configuration: DOCKER_POOL_FORCE_OVERRIDE=true"

mkdir -p /etc/docker

# Backup original daemon.json if it exists
if [ -f /etc/docker/daemon.json ]; then
    cp /etc/docker/daemon.json /etc/docker/daemon.json.original-"$DATE"
fi

# Create Docker configuration with or without address pools based on whether they were explicitly provided
if [ "$DOCKER_POOL_FORCE_OVERRIDE" = true ] || [ "$EXISTING_POOL_CONFIGURED" = false ]; then
    # First check if the configuration would actually change anything
    if [ -f /etc/docker/daemon.json ]; then
        CURRENT_POOL_BASE=$(jq -r '.["default-address-pools"][0].base' /etc/docker/daemon.json 2>/dev/null)
        CURRENT_POOL_SIZE=$(jq -r '.["default-address-pools"][0].size' /etc/docker/daemon.json 2>/dev/null)

        if [ "$CURRENT_POOL_BASE" = "$DOCKER_ADDRESS_POOL_BASE" ] && [ "$CURRENT_POOL_SIZE" = "$DOCKER_ADDRESS_POOL_SIZE" ]; then
            echo " - Network pool configuration unchanged, skipping update"
            NEED_MERGE=false
        else
            # If force override is enabled or no existing configuration exists,
            # create a new configuration with the specified address pools
            echo " - Creating new Docker configuration with network pool: ${DOCKER_ADDRESS_POOL_BASE}/${DOCKER_ADDRESS_POOL_SIZE}"
            cat >/etc/docker/daemon.json <<EOL
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "default-address-pools": [
    {"base":"${DOCKER_ADDRESS_POOL_BASE}","size":${DOCKER_ADDRESS_POOL_SIZE}}
  ]
}
EOL
            NEED_MERGE=true
        fi
    else
        # No existing configuration, create new one
        echo " - Creating new Docker configuration with network pool: ${DOCKER_ADDRESS_POOL_BASE}/${DOCKER_ADDRESS_POOL_SIZE}"
        cat >/etc/docker/daemon.json <<EOL
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "default-address-pools": [
    {"base":"${DOCKER_ADDRESS_POOL_BASE}","size":${DOCKER_ADDRESS_POOL_SIZE}}
  ]
}
EOL
        NEED_MERGE=true
    fi
else
    # Check if we need to update log settings
    if [ -f /etc/docker/daemon.json ] && jq -e '.["log-driver"] == "json-file" and .["log-opts"]["max-size"] == "10m" and .["log-opts"]["max-file"] == "3"' /etc/docker/daemon.json >/dev/null 2>&1; then
        echo " - Log configuration is up to date"
        NEED_MERGE=false
    else
        # Create a configuration without address pools to preserve existing ones
        cat >/etc/docker/daemon.json.velo <<EOL
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOL
        NEED_MERGE=true
    fi
fi

# If Docker configuration needs to be updated, restart Docker
if [ "$NEED_MERGE" = true ]; then
    echo " - Configuration updated - restarting Docker daemon..."
    restart_docker_service
else
    echo " - Configuration is up to date"
fi

echo -e "4. Initialize Docker Swarm. "

# Check if Docker Swarm is already initialized
if docker info | grep -q "Swarm: active"; then
    echo " - Docker Swarm is already initialized."
else
    echo " - Initializing Docker Swarm..."
    docker swarm init --advertise-addr "$(hostname -i | awk '{print $1}')" || {
        echo " - Failed to initialize Docker Swarm. Please check your network configuration."
        exit 1
    }
    echo " - Docker Swarm initialized successfully."
fi

echo -e "5. Download and install Velo. "

# Create Velo directories
mkdir -p $VELO_DATA_DIR/bin

# Download latest Velo binary
echo " - Downloading Velo binary..."

# Try to get latest release, fallback to building from source
LATEST_RELEASE_URL=$(curl -s https://api.github.com/repos/jasonlovesdoggo/velo/releases/latest | grep "browser_download_url.*Linux_$(uname -m).tar.gz" | cut -d '"' -f 4 2>/dev/null || true)

if [ -n "$LATEST_RELEASE_URL" ]; then
    echo " - Downloading from: $LATEST_RELEASE_URL"
    curl -L -s "$LATEST_RELEASE_URL" -o /tmp/velo.tar.gz
    # Extract the binary
    echo " - Extracting Velo binary..."
    tar -xzf /tmp/velo.tar.gz -C /tmp
    cp /tmp/velo $VELO_DATA_DIR/bin/
    chmod +x $VELO_DATA_DIR/bin/velo
    rm -f /tmp/velo.tar.gz /tmp/velo
else
    echo " - No pre-built release found. Building from source..."
    
    # Install Go if not present
    if ! command -v go >/dev/null 2>&1; then
        echo " - Installing Go..."
        case "$OS_TYPE" in
        ubuntu | debian | raspbian)
            apt-get update -y >/dev/null
            apt-get install -y golang-go >/dev/null
            ;;
        centos | fedora | rhel | ol | rocky | almalinux)
            dnf install -y golang >/dev/null
            ;;
        arch)
            pacman -Sy --noconfirm go >/dev/null
            ;;
        alpine)
            apk add go >/dev/null
            ;;
        *)
            echo " - Go installation not supported for $OS_TYPE. Please install Go manually."
            exit 1
            ;;
        esac
    fi
    
    # Clone and build Velo
    echo " - Cloning Velo repository..."
    git clone https://github.com/jasonlovesdoggo/velo.git /tmp/velo-src >/dev/null 2>&1
    cd /tmp/velo-src
    
    echo " - Building Velo..."
    go mod download >/dev/null 2>&1
    go build -o velo ./cmd/velo >/dev/null 2>&1
    go build -o veloctl ./cmd/cli >/dev/null 2>&1
    go build -o velo-runner ./cmd/runner >/dev/null 2>&1
    
    # Copy binaries
    cp velo $VELO_DATA_DIR/bin/
    cp veloctl $VELO_DATA_DIR/bin/
    cp velo-runner $VELO_DATA_DIR/bin/
    chmod +x $VELO_DATA_DIR/bin/*
    
    # Clean up
    cd /
    rm -rf /tmp/velo-src
fi

# Create a symbolic link to make Velo available system-wide
ln -sf $VELO_DATA_DIR/bin/velo /usr/local/bin/velo

echo -e "6. Create systemd service for Velo. "

# Create systemd service file
cat > /etc/systemd/system/velo.service << EOL
[Unit]
Description=Velo Deployment Platform
After=docker.service
Requires=docker.service

[Service]
Type=simple
ExecStart=/usr/local/bin/velo --manager
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOL

# Reload systemd, enable and start Velo service
systemctl daemon-reload
systemctl enable velo.service
systemctl start velo.service

echo -e "7. Verify installation. "

# Wait for Velo to start
echo " - Waiting for Velo to start..."
sleep 5

# Check if Velo service is running
if systemctl is-active --quiet velo; then
    echo " - Velo service is running."
else
    echo " - Velo service failed to start. Please check the logs with 'journalctl -u velo'."
fi

# Get IP addresses for access
IPV4_PUBLIC_IP=$(curl -4s https://ifconfig.io || true)
IPV6_PUBLIC_IP=$(curl -6s https://ifconfig.io || true)

echo -e "\033[0;35m
 __      __   _          _____           _        _ _           _ 
 \ \    / /__| |___     |_   _|_ _  ___ | |_ __ _| | |___ _  __| |
  \ \/\/ / -_) / _ \      | |/ _\` |(_-< |  _/ _\` | | / -_) |/ _\` |
   \_/\_/\___|_\___/      |_|\__,_|/__/  \__\__,_|_|_\___|_|\__,_|
                                                                  
\033[0m"

echo -e "\nYour Velo instance is ready to use!\n"

if [ -n "$IPV4_PUBLIC_IP" ]; then
    echo -e "You can access Velo through your Public IPV4: http://$IPV4_PUBLIC_IP:37355"
fi

if [ -n "$IPV6_PUBLIC_IP" ]; then
    echo -e "You can access Velo through your Public IPv6: http://[$IPV6_PUBLIC_IP]:37355"
fi

set +e
DEFAULT_PRIVATE_IP=$(ip route get 1 | sed -n 's/^.*src \([0-9.]*\) .*$/\1/p')
PRIVATE_IPS=$(hostname -I 2>/dev/null || ip -o addr show scope global | awk '{print $4}' | cut -d/ -f1)
set -e

if [ -n "$PRIVATE_IPS" ]; then
    echo -e "\nIf your Public IP is not accessible, you can use the following Private IPs:\n"
    for IP in $PRIVATE_IPS; do
        if [ "$IP" != "$DEFAULT_PRIVATE_IP" ]; then
            echo -e "http://$IP:37355"
        fi
    done
fi

echo -e "\nInstallation completed successfully!"
echo -e "Use 'systemctl status velo' to check the service status."
echo -e "Use 'journalctl -u velo' to view the service logs."