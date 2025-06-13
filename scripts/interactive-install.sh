#!/bin/bash
## Velo Interactive Installation Script
## This script provides an interactive installation experience for Velo

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_question() {
    echo -e "${BLUE}?${NC} $1"
}

# Check if running as root
if [ $EUID != 0 ]; then
    print_error "Please run this script as root or with sudo"
    exit 1
fi

# Welcome banner
echo -e "${PURPLE}"
cat << "EOF"
 __      __   _          ___         _        _ _           
 \ \    / /__| |___     |_ _|_ _  __| |_ __ _| | |___ _ _ 
  \ \/\/ / -_) / _ \     | || ' \(_-<  _/ _` | | / -_) '_|
   \_/\_/\___|_\___/    |___|_||_/__/\__\__,_|_|_\___|_|  
                                                          
EOF
echo -e "${NC}"

print_info "Welcome to the Velo Interactive Installer!"
print_info "This installer will guide you through setting up Velo on your system."
echo

# Detect OS
OS_TYPE=$(grep -w "ID" /etc/os-release | cut -d "=" -f 2 | tr -d '"')
print_info "Detected OS: $OS_TYPE"

# Get user preferences
echo
print_question "Installation Configuration:"

# Docker network pool
read -p "$(echo -e "${BLUE}?${NC} Docker network pool base (default: 10.0.0.0/8): ")" DOCKER_POOL_BASE
DOCKER_POOL_BASE=${DOCKER_POOL_BASE:-"10.0.0.0/8"}

read -p "$(echo -e "${BLUE}?${NC} Docker network pool size (default: 24): ")" DOCKER_POOL_SIZE
DOCKER_POOL_SIZE=${DOCKER_POOL_SIZE:-24}

# Installation options
echo
print_question "What would you like to install?"
echo "1) Manager node (orchestrates deployments)"
echo "2) Worker node (runs containers)"
echo "3) Both (single-node setup)"
read -p "$(echo -e "${BLUE}?${NC} Choose option [1-3] (default: 3): ")" INSTALL_TYPE
INSTALL_TYPE=${INSTALL_TYPE:-3}

# Service configuration
if [[ "$INSTALL_TYPE" == "1" || "$INSTALL_TYPE" == "3" ]]; then
    echo
    read -p "$(echo -e "${BLUE}?${NC} Start Velo service automatically? [Y/n]: ")" AUTO_START
    AUTO_START=${AUTO_START:-Y}
    
    read -p "$(echo -e "${BLUE}?${NC} API port (default: 37355): ")" API_PORT
    API_PORT=${API_PORT:-37355}
fi

# Confirmation
echo
print_info "Installation Summary:"
echo "  OS Type: $OS_TYPE"
echo "  Docker Pool: $DOCKER_POOL_BASE (size $DOCKER_POOL_SIZE)"
case $INSTALL_TYPE in
    1) echo "  Install Type: Manager node only" ;;
    2) echo "  Install Type: Worker node only" ;;
    3) echo "  Install Type: Both (single-node)" ;;
esac

if [[ "$INSTALL_TYPE" == "1" || "$INSTALL_TYPE" == "3" ]]; then
    echo "  API Port: $API_PORT"
    echo "  Auto-start: $AUTO_START"
fi

echo
read -p "$(echo -e "${BLUE}?${NC} Proceed with installation? [Y/n]: ")" CONFIRM
CONFIRM=${CONFIRM:-Y}

if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    print_warning "Installation cancelled by user"
    exit 0
fi

echo
print_info "Starting installation..."

# Export environment variables for main installer
export DOCKER_ADDRESS_POOL_BASE="$DOCKER_POOL_BASE"
export DOCKER_ADDRESS_POOL_SIZE="$DOCKER_POOL_SIZE"
export VELO_API_PORT="$API_PORT"
export VELO_INSTALL_TYPE="$INSTALL_TYPE"
export VELO_AUTO_START="$AUTO_START"

# Download and run main installer
print_info "Downloading main installation script..."
INSTALLER_URL="https://raw.githubusercontent.com/jasonlovesdoggo/velo/main/install.sh"

# Try to download from repository, fallback to local
if curl -fsSL "$INSTALLER_URL" -o /tmp/velo-install.sh 2>/dev/null; then
    print_info "Running main installation script..."
    bash /tmp/velo-install.sh
    rm -f /tmp/velo-install.sh
else
    print_warning "Could not download installer from repository."
    print_info "Please ensure you have internet connectivity or use the local install.sh script."
    exit 1
fi

print_success "Interactive installation completed!"

# Post-installation instructions
echo
print_info "Next steps:"
case $INSTALL_TYPE in
    1)
        echo "  • Your Velo manager is running on port $API_PORT"
        echo "  • Use 'veloctl --server localhost:$API_PORT deploy' to deploy services"
        echo "  • To add worker nodes, get the join token with 'docker swarm join-token worker'"
        ;;
    2)
        echo "  • Your worker node is ready to join a Velo cluster"
        echo "  • Use 'docker swarm join' with the token from your manager node"
        ;;
    3)
        echo "  • Your single-node Velo cluster is ready!"
        echo "  • Use 'veloctl deploy --service myapp --image nginx:latest' to deploy your first service"
        echo "  • Access the API at localhost:$API_PORT"
        ;;
esac

echo "  • Check service status: systemctl status velo"
echo "  • View logs: journalctl -u velo -f"
echo
print_success "Welcome to Velo! 🚀"