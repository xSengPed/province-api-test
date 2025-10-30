#!/bin/bash

# Thai Location API - Build Script
# This script builds the Docker image for the Thai Location API

set -e

# Configuration
IMAGE_NAME="thai-location-api"
IMAGE_TAG="latest"
CONTAINER_NAME="thai-location-api"

echo "🏗️  Building Thai Location API Docker Image..."

# Stop and remove existing container if it exists
if docker ps -a --format 'table {{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "🛑 Stopping existing container..."
    docker stop $CONTAINER_NAME
    echo "🗑️  Removing existing container..."
    docker rm $CONTAINER_NAME
fi

# Remove existing image if it exists
if docker images --format 'table {{.Repository}}:{{.Tag}}' | grep -q "^${IMAGE_NAME}:${IMAGE_TAG}$"; then
    echo "🗑️  Removing existing image..."
    docker rmi ${IMAGE_NAME}:${IMAGE_TAG}
fi

# Build the Docker image
echo "🔨 Building Docker image..."
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

# Test the image
echo "🧪 Testing the image..."
docker run --rm -d --name ${CONTAINER_NAME}-test -p 3001:3000 ${IMAGE_NAME}:${IMAGE_TAG}

# Wait for the container to start
sleep 5

# Health check
if curl -f http://localhost:3001/health; then
    echo "✅ Health check passed!"
else
    echo "❌ Health check failed!"
    docker stop ${CONTAINER_NAME}-test
    exit 1
fi

# Stop test container
docker stop ${CONTAINER_NAME}-test

echo "✅ Build completed successfully!"
echo "📦 Image: ${IMAGE_NAME}:${IMAGE_TAG}"
echo "🚀 Ready for deployment!"

# Optional: Tag for registry
read -p "Do you want to tag this image for a registry? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter registry URL (e.g., registry.domain.com): " REGISTRY_URL
    if [ ! -z "$REGISTRY_URL" ]; then
        FULL_TAG="${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"
        docker tag ${IMAGE_NAME}:${IMAGE_TAG} $FULL_TAG
        echo "🏷️  Tagged as: $FULL_TAG"
        
        read -p "Do you want to push to registry? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker push $FULL_TAG
            echo "📤 Pushed to registry: $FULL_TAG"
        fi
    fi
fi