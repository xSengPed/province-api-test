# Portainer Deployment Guide

This guide will help you deploy the Thai Location API on Portainer.

## Prerequisites

- Docker installed on your server
- Portainer installed and accessible
- Access to your server's command line (for initial image build)

## Step 1: Build the Docker Image

### Option A: Build on Server
```bash
# 1. Upload project files to your server
scp -r . user@your-server:/path/to/thai-location-api/

# 2. SSH to your server
ssh user@your-server

# 3. Navigate to project directory
cd /path/to/thai-location-api/

# 4. Build the image
./build.sh
```

### Option B: Build Locally and Push to Registry
```bash
# 1. Build locally
./build.sh

# 2. Tag for your registry
docker tag thai-location-api:latest your-registry.com/thai-location-api:latest

# 3. Push to registry
docker push your-registry.com/thai-location-api:latest
```

## Step 2: Deploy via Portainer

### Method 1: Using Portainer Stacks (Recommended)

1. **Access Portainer Web Interface**
   - Open your browser and go to your Portainer URL
   - Login with your credentials

2. **Create a New Stack**
   - Go to "Stacks" in the left sidebar
   - Click "Add stack"
   - Name: `thai-location-api`

3. **Add Stack Configuration**
   Copy and paste the following configuration:

   ```yaml
   version: '3.8'

   services:
     thai-location-api:
       image: thai-location-api:latest  # or your-registry.com/thai-location-api:latest
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
   ```

4. **Deploy the Stack**
   - Click "Deploy the stack"
   - Wait for deployment to complete

### Method 2: Using Portainer Container Interface

1. **Go to Containers**
   - Click "Containers" in the left sidebar
   - Click "Add container"

2. **Container Configuration**
   - **Name**: `thai-location-api`
   - **Image**: `thai-location-api:latest`
   - **Port mapping**: Host `3000` → Container `3000`
   - **Restart policy**: Unless stopped
   - **Environment variables**:
     - `PORT=3000`

3. **Advanced Settings**
   - **Resources**: Set CPU limit to 0.5 and memory to 256MB
   - **Health check**: 
     - Command: `wget --no-verbose --tries=1 --spider http://localhost:3000/health`
     - Interval: 30s

4. **Deploy Container**
   - Click "Deploy the container"

## Step 3: Configure Reverse Proxy (Optional)

### With Traefik
Add these labels to your container or stack:

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.thai-location-api.rule=Host(`api.yourdomain.com`)"
  - "traefik.http.routers.thai-location-api.entrypoints=websecure"
  - "traefik.http.routers.thai-location-api.tls.certresolver=letsencrypt"
  - "traefik.http.services.thai-location-api.loadbalancer.server.port=3000"
```

### With Nginx Proxy Manager
1. Create a new proxy host
2. **Domain**: `api.yourdomain.com`
3. **Forward Hostname/IP**: Your server IP
4. **Forward Port**: `3000`
5. Enable SSL if needed

## Step 4: Verify Deployment

1. **Check Container Status**
   - In Portainer, go to Containers
   - Verify the container is running (green status)

2. **Test Health Endpoint**
   ```bash
   curl http://your-server:3000/health
   ```
   
   Expected response:
   ```json
   {
     "status": "OK",
     "message": "Thai Location API is running"
   }
   ```

3. **Test API Endpoints**
   ```bash
   # Get provinces
   curl http://your-server:3000/api/v1/provinces?limit=5
   
   # Get geographies
   curl http://your-server:3000/api/v1/geographies
   ```

## Step 5: Monitor and Maintain

### View Logs
1. In Portainer, go to Containers
2. Click on `thai-location-api`
3. Go to "Logs" tab

### Update Application
1. Build new image with updated code
2. In Portainer, go to Containers
3. Select `thai-location-api` container
4. Click "Recreate" or update the stack

### Backup Configuration
Save your stack configuration or container settings for disaster recovery.

## Troubleshooting

### Container Won't Start
1. Check logs in Portainer
2. Verify image exists: `docker images | grep thai-location-api`
3. Check data files are included in image

### Health Check Fails
1. Verify port 3000 is not blocked
2. Check if data files are accessible
3. Review application logs

### Performance Issues
1. Increase memory limit in container settings
2. Monitor CPU usage
3. Check for memory leaks in logs

### Port Conflicts
1. Change host port mapping (e.g., 3001:3000)
2. Update health check URL accordingly
3. Update proxy configuration

## Production Recommendations

1. **Use a Registry**: Store images in a private registry
2. **Enable SSL**: Use reverse proxy with SSL certificates
3. **Monitor Resources**: Set up alerts for CPU/memory usage
4. **Backup Data**: Regular backups of configuration
5. **Log Management**: Centralized logging solution
6. **Health Monitoring**: External health check monitoring

## Example Production Stack

```yaml
version: '3.8'

services:
  thai-location-api:
    image: registry.yourdomain.com/thai-location-api:v1.0.0
    container_name: thai-location-api-prod
    restart: unless-stopped
    environment:
      - PORT=3000
    networks:
      - proxy
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.thai-api.rule=Host(`api.yourdomain.com`)"
      - "traefik.http.routers.thai-api.entrypoints=websecure"
      - "traefik.http.routers.thai-api.tls.certresolver=letsencrypt"
      - "traefik.http.services.thai-api.loadbalancer.server.port=3000"
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M

networks:
  proxy:
    external: true
```

## Support

For issues and questions:
1. Check container logs in Portainer
2. Review this deployment guide
3. Check the main README.md for API documentation