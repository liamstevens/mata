# Docker Compose Example

This example demonstrates Mata in a complete stack with Nginx and multiple backend services.

## Architecture

```
Client → Nginx (Port 8000) → Mata (Port 8080) → [App Service, Analytics Service]
```

## Services

- **Nginx**: Frontend proxy with SSL termination, caching, and load balancing
- **Mata**: Traffic duplicator that sends requests to both app and analytics
- **App**: Main application service (Python HTTP server)
- **Analytics**: Analytics service receiving duplicated traffic

## Usage

1. **Start the stack:**
   ```bash
   make docker-up
   ```

2. **Test the setup:**
   ```bash
   # Send request through the stack
   curl http://localhost:8000
   
   # Both app and analytics services receive the request
   # Response comes from the app service
   ```

3. **View logs:**
   ```bash
   make docker-logs
   ```

4. **Stop the stack:**
   ```bash
   make docker-down
   ```

## What Happens

1. Client makes HTTP request to localhost:8000
2. Nginx receives request and forwards to Mata on port 8080
3. Mata duplicates the request to both:
   - App service (returns response to client)
   - Analytics service (receives copy for monitoring)
4. Client receives response from app service
5. Analytics service can log, process, or analyze the duplicated traffic

## Files

- `docker-compose.yml` - Service definitions
- `Dockerfile` - Mata container build
- `examples/nginx.conf` - Nginx configuration
- `examples/app/index.html` - Main app content
- `examples/analytics/index.html` - Analytics service content