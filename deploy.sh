#!/bin/bash

# Thai Location API - Deploy Script for Portainer
# This script deploys the application using Docker Compose

set -e

# Configuration
STACK_NAME="thai-location-api"
COMPOSE_FILE="docker-compose.yml"

echo "🚀 Deploying Thai Location API to Portainer..."

# Check if docker-compose.yml exists
if [ ! -f "$COMPOSE_FILE" ]; then
    echo "❌ Error: $COMPOSE_FILE not found!"
    exit 1
fi

# Function to deploy with docker-compose
deploy_with_compose() {
    echo "📋 Using Docker Compose for deployment..."
    
    # Stop existing services
    echo "🛑 Stopping existing services..."
    docker-compose down --remove-orphans || true
    
    # Pull/build latest images
    echo "📦 Building/pulling latest images..."
    docker-compose build --no-cache
    
    # Start services
    echo "🔄 Starting services..."
    docker-compose up -d
    
    # Wait for services to be ready
    echo "⏳ Waiting for services to be ready..."
    sleep 10
    
    # Health check
    echo "🧪 Performing health check..."
    if curl -f http://localhost:3000/health; then
        echo "✅ Deployment successful!"
        echo "🌐 API is available at: http://localhost:3000"
        echo "📋 Health endpoint: http://localhost:3000/health"
        echo "📚 API documentation:"
        echo "   - Geographies: GET /api/v1/geographies"
        echo "   - Provinces: GET /api/v1/provinces"
        echo "   - Districts: GET /api/v1/districts"
        echo "   - Sub-districts: GET /api/v1/subdistricts"
    else
        echo "❌ Health check failed!"
        echo "📊 Checking container logs..."
        docker-compose logs
        exit 1
    fi
}

# Function to create Portainer stack file
create_portainer_stack() {
    echo "📝 Creating Portainer stack configuration..."
    
    cat > portainer-stack.yml << EOF
version: '3.8'

services:
  thai-location-api:
    image: thai-location-api:latest
    container_name: thai-location-api
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - PORT=3000
    networks:
      - thai-location-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M

networks:
  thai-location-network:
    driver: bridge
EOF

    echo "✅ Portainer stack file created: portainer-stack.yml"
    echo ""
    echo "📋 To deploy in Portainer:"
    echo "1. Copy the contents of portainer-stack.yml"
    echo "2. Go to Portainer -> Stacks -> Add Stack"
    echo "3. Name: thai-location-api"
    echo "4. Paste the stack configuration"
    echo "5. Deploy the stack"
    echo ""
}

# Main deployment function
main() {
    echo "🔧 Choose deployment method:"
    echo "1. Docker Compose (local)"
    echo "2. Generate Portainer stack file"
    echo "3. Both"
    read -p "Enter your choice (1-3): " choice
    
    case $choice in
        1)
            deploy_with_compose
            ;;
        2)
            create_portainer_stack
            ;;
        3)
            deploy_with_compose
            create_portainer_stack
            ;;
        *)
            echo "❌ Invalid choice!"
            exit 1
            ;;
    esac
}

# Show current status
show_status() {
    echo ""
    echo "📊 Current Status:"
    if docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep -q thai-location-api; then
        docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep thai-location-api
    else
        echo "No containers running"
    fi
}

# Check if argument is provided
if [ "$1" = "status" ]; then
    show_status
    exit 0
elif [ "$1" = "stop" ]; then
    echo "🛑 Stopping Thai Location API..."
    docker-compose down
    exit 0
elif [ "$1" = "logs" ]; then
    echo "📋 Showing logs..."
    docker-compose logs -f
    exit 0
elif [ "$1" = "restart" ]; then
    echo "🔄 Restarting Thai Location API..."
    docker-compose restart
    exit 0
fi

# Run main deployment
main

# Show final status
show_status

echo ""
echo "🎉 Deployment completed!"