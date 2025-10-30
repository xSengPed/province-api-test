#!/bin/bash

# Thai Location API - Cloudflare Deployment Script
# This script helps deploy the API with Cloudflare integration

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_NAME="thai-location-api"
DOMAIN=""
SUBDOMAIN="api"
DEPLOYMENT_TYPE=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

echo_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

echo_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo "🌥️  Cloudflare Deployment for Thai Location API"
    echo "=================================================="
}

check_prerequisites() {
    echo_info "Checking prerequisites..."
    
    # Check Docker
    if ! command -v docker >/dev/null 2>&1; then
        echo_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose >/dev/null 2>&1; then
        echo_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    echo_success "Prerequisites check passed!"
}

get_deployment_type() {
    echo ""
    echo_info "Choose your Cloudflare deployment method:"
    echo "1. Cloudflare Tunnel (Recommended - No open ports needed)"
    echo "2. Traditional hosting with Cloudflare proxy"
    echo "3. Local development with Cloudflare testing"
    echo ""
    
    read -p "Enter your choice (1-3): " choice
    
    case $choice in
        1)
            DEPLOYMENT_TYPE="tunnel"
            ;;
        2)
            DEPLOYMENT_TYPE="traditional"
            ;;
        3)
            DEPLOYMENT_TYPE="development"
            ;;
        *)
            echo_error "Invalid choice!"
            exit 1
            ;;
    esac
    
    echo_success "Selected deployment type: $DEPLOYMENT_TYPE"
}

get_domain_info() {
    echo ""
    echo_info "Domain Configuration"
    echo "Please enter your domain information:"
    
    read -p "Domain name (e.g., yourdomain.com): " DOMAIN
    read -p "Subdomain (default: api): " subdomain_input
    
    if [ ! -z "$subdomain_input" ]; then
        SUBDOMAIN="$subdomain_input"
    fi
    
    if [ -z "$DOMAIN" ]; then
        echo_error "Domain name is required!"
        exit 1
    fi
    
    echo_success "API will be available at: https://${SUBDOMAIN}.${DOMAIN}"
}

deploy_tunnel() {
    echo_info "Setting up Cloudflare Tunnel deployment..."
    
    # Check if cloudflared is installed
    if ! command -v cloudflared >/dev/null 2>&1; then
        echo_warning "Cloudflared is not installed. Installing..."
        
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
            curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
            sudo dpkg -i cloudflared.deb
            rm cloudflared.deb
        elif [[ "$OSTYPE" == "darwin"* ]]; then
            if command -v brew >/dev/null 2>&1; then
                brew install cloudflared
            else
                echo_error "Please install Homebrew first or manually install cloudflared"
                exit 1
            fi
        else
            echo_error "Unsupported OS. Please install cloudflared manually."
            exit 1
        fi
    fi
    
    echo_info "Building and starting API..."
    docker-compose -f cloudflare-docker-compose.yml build
    docker-compose -f cloudflare-docker-compose.yml up -d thai-location-api
    
    echo_info "Waiting for API to be ready..."
    sleep 10
    
    # Test API
    if curl -f http://localhost:3000/health >/dev/null 2>&1; then
        echo_success "API is running successfully!"
    else
        echo_error "API failed to start. Check logs with: docker-compose logs"
        exit 1
    fi
    
    echo ""
    echo_info "Next steps for Cloudflare Tunnel:"
    echo "1. Run: cloudflared tunnel login"
    echo "2. Run: cloudflared tunnel create ${API_NAME}"
    echo "3. Run: cloudflared tunnel route dns ${API_NAME} ${SUBDOMAIN}.${DOMAIN}"
    echo "4. Create config file at ~/.cloudflared/config.yml with:"
    echo ""
    echo "tunnel: ${API_NAME}"
    echo "credentials-file: /path/to/your/tunnel/credentials.json"
    echo ""
    echo "ingress:"
    echo "  - hostname: ${SUBDOMAIN}.${DOMAIN}"
    echo "    service: http://localhost:3000"
    echo "  - service: http_status:404"
    echo ""
    echo "5. Run: cloudflared tunnel run ${API_NAME}"
    echo ""
}

deploy_traditional() {
    echo_info "Setting up traditional hosting with Cloudflare proxy..."
    
    # Update nginx config with domain
    sed "s/api\.yourdomain\.com/${SUBDOMAIN}.${DOMAIN}/g" nginx.conf > nginx_temp.conf
    mv nginx_temp.conf nginx.conf
    
    echo_info "Building and starting services..."
    docker-compose -f cloudflare-docker-compose.yml --profile traditional-hosting up -d
    
    echo_info "Waiting for services to be ready..."
    sleep 15
    
    # Test API through nginx
    if curl -f http://localhost/health >/dev/null 2>&1; then
        echo_success "Services are running successfully!"
    else
        echo_error "Services failed to start. Check logs with: docker-compose logs"
        exit 1
    fi
    
    echo ""
    echo_info "Next steps for traditional hosting:"
    echo "1. Point your domain ${SUBDOMAIN}.${DOMAIN} to this server's IP"
    echo "2. In Cloudflare DNS settings:"
    echo "   - Add A record: ${SUBDOMAIN} -> Your server IP"
    echo "   - Enable proxy (orange cloud)"
    echo "3. Configure Cloudflare SSL/TLS settings:"
    echo "   - SSL/TLS mode: Full or Full (strict)"
    echo "   - Enable 'Always Use HTTPS'"
    echo "4. Optional: Configure Cloudflare page rules for caching"
    echo ""
}

deploy_development() {
    echo_info "Setting up development environment..."
    
    docker-compose up -d
    
    echo_info "Waiting for API to be ready..."
    sleep 10
    
    if curl -f http://localhost:3000/health >/dev/null 2>&1; then
        echo_success "Development API is running!"
        echo_info "API available at: http://localhost:3000"
        echo_info "Test with: curl http://localhost:3000/api/v1/provinces?limit=5"
    else
        echo_error "API failed to start. Check logs with: docker-compose logs"
        exit 1
    fi
}

create_cloudflare_config() {
    echo_info "Creating Cloudflare configuration files..."
    
    # Create .env file for Cloudflare settings
    cat > .env.cloudflare << EOF
# Cloudflare Configuration
DOMAIN=${DOMAIN}
SUBDOMAIN=${SUBDOMAIN}
FULL_DOMAIN=${SUBDOMAIN}.${DOMAIN}

# Cloudflare Tunnel (if using tunnel deployment)
# CLOUDFLARE_TUNNEL_TOKEN=your_tunnel_token_here

# API Configuration
PORT=3000
NODE_ENV=production
CLOUDFLARE_ENABLED=true
TRUST_PROXY=true
EOF

    echo_success "Created .env.cloudflare with your domain configuration"
}

show_status() {
    echo ""
    echo_info "Current Status:"
    
    if docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep -q "${API_NAME}"; then
        docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep "${API_NAME}"
        echo ""
        
        # Test API
        if curl -f http://localhost:3000/health >/dev/null 2>&1; then
            echo_success "API Health Check: ✅ Healthy"
        else
            echo_warning "API Health Check: ❌ Unhealthy"
        fi
    else
        echo_warning "No containers running"
    fi
}

show_help() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  deploy    Deploy the API with Cloudflare integration"
    echo "  status    Show current deployment status"
    echo "  stop      Stop all services"
    echo "  logs      Show service logs"
    echo "  restart   Restart services"
    echo "  cleanup   Remove all containers and images"
    echo "  help      Show this help message"
    echo ""
}

# Main execution
main() {
    print_header
    
    case "${1:-deploy}" in
        "deploy")
            check_prerequisites
            get_deployment_type
            get_domain_info
            create_cloudflare_config
            
            case $DEPLOYMENT_TYPE in
                "tunnel")
                    deploy_tunnel
                    ;;
                "traditional")
                    deploy_traditional
                    ;;
                "development")
                    deploy_development
                    ;;
            esac
            
            show_status
            echo_success "Deployment completed!"
            ;;
            
        "status")
            show_status
            ;;
            
        "stop")
            echo_info "Stopping services..."
            docker-compose down
            docker-compose -f cloudflare-docker-compose.yml down
            echo_success "Services stopped"
            ;;
            
        "logs")
            echo_info "Showing logs..."
            docker-compose logs -f
            ;;
            
        "restart")
            echo_info "Restarting services..."
            docker-compose restart
            docker-compose -f cloudflare-docker-compose.yml restart
            echo_success "Services restarted"
            ;;
            
        "cleanup")
            echo_warning "This will remove all containers and images. Are you sure? (y/N)"
            read -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                docker-compose down --rmi all --volumes --remove-orphans
                docker-compose -f cloudflare-docker-compose.yml down --rmi all --volumes --remove-orphans
                echo_success "Cleanup completed"
            fi
            ;;
            
        "help"|"-h"|"--help")
            show_help
            ;;
            
        *)
            echo_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"