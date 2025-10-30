# Thai Location API

A high-performance REST API for Thai geographic data (Provinces, Districts, Sub-districts/Tambons) built with Go and Fiber framework.

## Features

- 🚀 Fast and lightweight API built with Go Fiber
- 📍 Complete Thai geographic data (Provinces, Districts, Sub-districts)
- 🔍 Search functionality with Thai and English names
- 📄 Pagination support
- 🐳 Docker containerized for easy deployment
- 🔄 Ready for Portainer deployment
- 🏥 Health check endpoints
- 🌐 CORS enabled
- 📊 Structured JSON responses

## API Endpoints

### Health Check
- `GET /health` - API health status

### Geographies
- `GET /api/v1/geographies` - Get all geographies (regions)

### Provinces
- `GET /api/v1/provinces` - Get all provinces
  - Query params: `geography_id`, `search`, `page`, `limit`
- `GET /api/v1/provinces/{id}` - Get province by ID
- `GET /api/v1/provinces/{id}/districts` - Get districts by province ID

### Districts
- `GET /api/v1/districts` - Get all districts
  - Query params: `province_id`, `search`, `page`, `limit`
- `GET /api/v1/districts/{id}` - Get district by ID
- `GET /api/v1/districts/{id}/subdistricts` - Get sub-districts by district ID

### Sub-districts (Tambons)
- `GET /api/v1/subdistricts` - Get all sub-districts
  - Query params: `district_id`, `zip_code`, `search`, `page`, `limit`
- `GET /api/v1/subdistricts/{id}` - Get sub-district by ID

## API Documentation (Swagger)

Interactive API documentation is available via Swagger UI. After starting the server (or when deployed), open:

```
http://localhost:3000/docs
```

This serves a browsable OpenAPI specification (`/docs/openapi.json`) and provides example requests/responses for the available endpoints.

## Query Parameters

- `search` - Search by Thai or English name
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)
- `geography_id` - Filter provinces by geography
- `province_id` - Filter districts by province
- `district_id` - Filter sub-districts by district
- `zip_code` - Filter sub-districts by postal code

## Example Requests

```bash
# Get all provinces
curl http://localhost:3000/api/v1/provinces

# Search provinces by name
curl "http://localhost:3000/api/v1/provinces?search=กรุงเทพ"

# Get districts in Bangkok (province_id=1)
curl http://localhost:3000/api/v1/provinces/1/districts

# Get sub-districts with pagination
curl "http://localhost:3000/api/v1/subdistricts?page=1&limit=10"

# Search by zip code
curl "http://localhost:3000/api/v1/subdistricts?zip_code=10200"
```

## Response Format

### Success Response
```json
{
  "status": "success",
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 77,
    "total_pages": 4
  }
}
```

### Error Response
```json
{
  "status": "error",
  "error": "Error message"
}
```

## Quick Start

### Using Docker Compose

1. **Build and run:**
   ```bash
   ./build.sh
   ./deploy.sh
   ```

2. **Access the API:**
   ```
   http://localhost:3000
   ```

### Manual Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Run the application:**
   ```bash
   go run .
   ```

## Deployment

### Docker Compose
```bash
# Build and deploy
./deploy.sh

# Check status
./deploy.sh status

# View logs
./deploy.sh logs

# Stop services
./deploy.sh stop

# Restart services
./deploy.sh restart
```

### Portainer

1. **Build the image:**
   ```bash
   ./build.sh
   ```

2. **Generate Portainer stack:**
   ```bash
   ./deploy.sh
   # Choose option 2 to generate portainer-stack.yml
   ```

3. **Deploy in Portainer:**
   - Go to Stacks → Add Stack
   - Name: `thai-location-api`
   - Copy contents from `portainer-stack.yml`
   - Deploy

### Environment Variables

- `PORT` - Server port (default: 3000)

## Data Structure

The API uses the following data structure:

- **Geography**: Regions of Thailand (North, Central, Northeast, West, East, South)
- **Province**: 77 provinces in Thailand
- **District**: Districts (Amphoe) within provinces
- **Sub-district**: Sub-districts (Tambon) within districts

## Development

### Project Structure
```
├── main.go           # Application entry point
├── models.go         # Data structures
├── service.go        # Data loading service
├── handlers.go       # API handlers
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Docker Compose configuration
├── build.sh          # Build script
├── deploy.sh         # Deployment script
└── data/raw/         # JSON data files
    ├── geographies.json
    ├── provinces.json
    ├── districts.json
    └── sub_districts.json
```

### Building
```bash
# Local build
go build -o thai-location-api

# Docker build
docker build -t thai-location-api .
```

## Performance

- Data is loaded into memory on startup for fast access
- Concurrent-safe with read-write mutex
- Efficient indexing for O(1) lookups
- Minimal memory footprint (~50MB container)

## License

MIT License

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if needed
5. Submit a pull request

## Support

For issues and questions, please create an issue in the repository.