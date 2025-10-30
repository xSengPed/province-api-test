# Cloudflare Deployment Guide

This guide will help you deploy the Thai Location API with Cloudflare as your public host, providing global CDN, DDoS protection, and SSL certificates.

## Overview

You have several options for deploying with Cloudflare:

1. **Cloudflare Tunnel** (Recommended) - Secure, no open ports needed
2. **Traditional hosting** with Cloudflare proxy
3. **Cloudflare Workers** (for edge deployment)

## Option 1: Cloudflare Tunnel (Recommended)

Cloudflare Tunnel creates a secure connection between your server and Cloudflare without exposing your server's IP.

### Prerequisites
- Domain managed by Cloudflare
- Server with Docker
- Cloudflare account

### Step 1: Install Cloudflared on Your Server

```bash
# For Ubuntu/Debian
curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared.deb

# For CentOS/RHEL
curl -L --output cloudflared.rpm https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-x86_64.rpm
sudo rpm -i cloudflared.rpm

# For macOS
brew install cloudflared
```

### Step 2: Authenticate Cloudflared

```bash
cloudflared tunnel login
```

This will open a browser to authenticate with Cloudflare.

### Step 3: Create a Tunnel

```bash
cloudflared tunnel create thai-location-api
```

This creates a tunnel and generates a UUID. Save this UUID!

### Step 4: Create DNS Record

```bash
# Replace YOUR_TUNNEL_UUID with the actual UUID from step 3
cloudflared tunnel route dns thai-location-api api.yourdomain.com
```

### Step 5: Configure the Tunnel

Create a configuration file at `~/.cloudflared/config.yml`:

```yaml
tunnel: thai-location-api
credentials-file: /home/your-user/.cloudflared/YOUR_TUNNEL_UUID.json

ingress:
  - hostname: api.yourdomain.com
    service: http://localhost:3000
  - service: http_status:404
```

### Step 6: Deploy Your API

Use the enhanced Docker Compose configuration:

```bash
# Use the cloudflare-docker-compose.yml
docker-compose -f cloudflare-docker-compose.yml up -d
```

### Step 7: Start the Tunnel

```bash
cloudflared tunnel run thai-location-api
```

Or run as a service:

```bash
sudo cloudflared service install
sudo systemctl start cloudflared
sudo systemctl enable cloudflared
```

## Option 2: Traditional Hosting with Cloudflare Proxy

### Step 1: Deploy API on Your Server

```bash
# Deploy using standard docker-compose
./deploy.sh
```

### Step 2: Configure Reverse Proxy (Nginx)

Create nginx configuration:

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Step 3: Configure Cloudflare DNS

1. Go to Cloudflare Dashboard
2. Select your domain
3. Go to DNS settings
4. Add an A record:
   - Name: `api`
   - Content: Your server's IP address
   - Proxy status: Proxied (orange cloud)

## Option 3: Enhanced Docker Compose with Cloudflare

I'll create an enhanced docker-compose file for better Cloudflare integration.

## Security Considerations

### Environment Variables
Set these environment variables for production:

```bash
export CLOUDFLARE_API_TOKEN="your_api_token"
export DOMAIN="yourdomain.com"
export SUBDOMAIN="api"
```

### Firewall Rules
With Cloudflare Tunnel, you don't need to open ports. For traditional hosting:

```bash
# Allow only Cloudflare IPs
ufw allow from 173.245.48.0/20
ufw allow from 103.21.244.0/22
ufw allow from 103.22.200.0/22
# ... (all Cloudflare IP ranges)
```

## Performance Optimization

### Cloudflare Settings
1. **Speed > Optimization**
   - Auto Minify: Enable CSS, HTML, JS
   - Brotli: Enable
   - Early Hints: Enable

2. **Caching > Configuration**
   - Browser Cache TTL: 4 hours
   - Development Mode: Off (for production)

3. **Page Rules** (for `/api/*`)
   - Cache Level: Cache Everything
   - Edge Cache TTL: 2 hours
   - Browser Cache TTL: 30 minutes

### API Response Headers
The API already includes appropriate caching headers, but you can enhance them for Cloudflare.

## Monitoring and Analytics

### Cloudflare Analytics
Enable analytics in Cloudflare Dashboard to monitor:
- Request volume
- Response times
- Geographic distribution
- Security threats blocked

### Health Monitoring
Set up health checks:
- Cloudflare Health Checks
- External monitoring services
- Custom alerting

## SSL/TLS Configuration

Cloudflare provides free SSL certificates. Configure:

1. **SSL/TLS > Overview**
   - Encryption mode: Full (strict)

2. **SSL/TLS > Edge Certificates**
   - Always Use HTTPS: On
   - HTTP Strict Transport Security (HSTS): Enable
   - Minimum TLS Version: 1.2

## Example API URLs

After deployment, your API will be available at:

```
https://api.yourdomain.com/health
https://api.yourdomain.com/api/v1/provinces
https://api.yourdomain.com/api/v1/districts
https://api.yourdomain.com/api/v1/subdistricts
```

## Cost Considerations

- **Cloudflare Free Plan**: Suitable for most use cases
- **Cloudflare Pro Plan** ($20/month): Better for high-traffic APIs
- **Server Costs**: Your VPS/cloud server costs

## Troubleshooting

### Common Issues

1. **522 Connection Timed Out**
   - Check if your server is running
   - Verify port 3000 is accessible locally

2. **525 SSL Handshake Failed**
   - Set Cloudflare SSL mode to "Flexible" initially
   - Then upgrade to "Full" or "Full (strict)"

3. **API Returns 404**
   - Check tunnel configuration
   - Verify DNS records
   - Check API server logs

### Debug Commands

```bash
# Check if API is running locally
curl http://localhost:3000/health

# Check tunnel status
cloudflared tunnel info thai-location-api

# View tunnel logs
cloudflared tunnel run thai-location-api --loglevel debug
```

## Next Steps

1. Deploy your API using one of the methods above
2. Configure Cloudflare settings for optimal performance
3. Set up monitoring and alerting
4. Test your API endpoints
5. Update documentation with your domain

Choose the deployment method that best fits your infrastructure and security requirements!