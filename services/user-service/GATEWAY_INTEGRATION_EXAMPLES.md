# API Gateway - User Service Integration Examples

This document provides example configurations for integrating the User Service microservice with various API gateway solutions.

## 1. Express.js with http-proxy-middleware

```typescript
// gateway.ts
import express from 'express';
import { createProxyMiddleware } from 'http-proxy-middleware';

const app = express();

// Health check passthrough
app.use(
  '/api/health',
  createProxyMiddleware({
    target: 'http://localhost:3000',
    changeOrigin: true,
    pathRewrite: { '^/api/health': '/health' },
  }),
);

// Users service routes
app.use(
  '/api/users',
  createProxyMiddleware({
    target: 'http://localhost:3000',
    changeOrigin: true,
    pathRewrite: { '^/api/users': '/users' },
  }),
);

app.listen(8000, () => {
  console.log('API Gateway running on port 8000');
  console.log('User Service available at /api/users');
});
```

## 2. NestJS Gateway

```typescript
// gateway.module.ts
import { Module } from '@nestjs/common';
import { HttpModule } from '@nestjs/axios';
import { GatewayController } from './gateway.controller';
import { GatewayService } from './gateway.service';

@Module({
  imports: [HttpModule],
  controllers: [GatewayController],
  providers: [GatewayService],
  exports: [HttpModule],
})
export class GatewayModule {}
```

```typescript
// gateway.controller.ts
import { Controller, All, Req, Res } from '@nestjs/common';
import { Request, Response } from 'express';
import { GatewayService } from './gateway.service';

@Controller('api')
export class GatewayController {
  constructor(private gatewayService: GatewayService) {}

  @All('users/*')
  @All('health')
  async route(@Req() req: Request, @Res() res: Response) {
    const path = req.path.replace('/api', '');
    const result = await this.gatewayService.forwardRequest(
      path,
      req.method,
      req.body,
      req.headers,
    );

    res.status(result.status).send(result.data);
  }
}
```

```typescript
// gateway.service.ts
import { Injectable } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class GatewayService {
  private userServiceUrl = 'http://localhost:3000';

  constructor(private httpService: HttpService) {}

  async forwardRequest(path: string, method: string, body: any, headers: any) {
    try {
      const url = `${this.userServiceUrl}${path}`;

      const response = await firstValueFrom(
        this.httpService.request({
          method: method as any,
          url,
          data: body,
          headers: this.sanitizeHeaders(headers),
        }),
      );

      return {
        status: response.status,
        data: response.data,
      };
    } catch (error: any) {
      return {
        status: error.response?.status || 500,
        data: error.response?.data || { message: 'Gateway error' },
      };
    }
  }

  private sanitizeHeaders(headers: any) {
    const sanitized = { ...headers };
    delete sanitized.host;
    delete sanitized['content-length'];
    return sanitized;
  }
}
```

## 3. Kong Configuration

```yaml
# Kong Gateway - docker-compose.yml
version: '3.8'
services:
  kong:
    image: kong:alpine
    environment:
      KONG_DATABASE: postgres
      KONG_PG_HOST: kong-db
      KONG_PROXY_ACCESS_LOG: /dev/stdout
      KONG_ADMIN_ACCESS_LOG: /dev/stdout
      KONG_PROXY_ERROR_LOG: /dev/stderr
      KONG_ADMIN_ERROR_LOG: /dev/stderr
      KONG_ADMIN_LISTEN: 0.0.0.0:8001
    ports:
      - '8000:8000'
      - '8001:8001'
    depends_on:
      - kong-db
    networks:
      - kong-net

  kong-db:
    image: postgres:15
    environment:
      POSTGRES_DB: kong
      POSTGRES_USER: kong
      POSTGRES_PASSWORD: kong
    networks:
      - kong-net

  user-service:
    build: ./services/user-service
    ports:
      - '3000:3000'
    networks:
      - kong-net

networks:
  kong-net:
    driver: bridge
```

```bash
# Configure Kong - Add upstream
curl -X POST http://localhost:8001/upstreams \
  -d "name=user-service"

# Add target to upstream
curl -X POST http://localhost:8001/upstreams/user-service/targets \
  -d "target=user-service:3000"

# Add service
curl -X POST http://localhost:8001/services \
  -d "name=user-service" \
  -d "host=user-service" \
  -d "port=3000" \
  -d "protocol=http"

# Add route
curl -X POST http://localhost:8001/services/user-service/routes \
  -d "paths[]=/api/users" \
  -d "strip_path=true"

# Add health check route
curl -X POST http://localhost:8001/services/user-service/routes \
  -d "paths[]=/api/health" \
  -d "strip_path=true"
```

## 4. Nginx Configuration

```nginx
# nginx.conf
upstream user_service {
    server localhost:3000;
    keepalive 32;
}

server {
    listen 8000;
    server_name localhost;

    # Users service
    location /api/users {
        proxy_pass http://user_service/users;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # CORS Headers
        add_header 'Access-Control-Allow-Origin' '*' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, PATCH, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Content-Type, Authorization' always;

        if ($request_method = 'OPTIONS') {
            return 204;
        }
    }

    # Health check
    location /api/health {
        proxy_pass http://user_service/health;
        proxy_http_version 1.1;
        proxy_set_header Connection '';
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/s;
    limit_req zone=api_limit burst=200 nodelay;
}
```

## 5. HAProxy Configuration

```
# haproxy.cfg
global
    maxconn 4096
    log localhost local0
    log localhost local1 notice

defaults
    log global
    mode http
    option httplog
    timeout connect 5000
    timeout client 50000
    timeout server 50000

frontend api_gateway
    bind *:8000
    default_backend user_service

backend user_service
    balance roundrobin
    option httpchk GET /health
    server user-service-1 localhost:3000 check inter 2000 rise 2 fall 3

    # Enable keep-alive
    http-reuse safe

    # Add request headers
    http-request set-header X-Forwarded-For %[src]
    http-request set-header X-Forwarded-Proto http

    # Remove path prefix /api/users -> /users
    reqrep ^([^\ :]*\ )/api/users(.*) \1/users\2
```

## 6. AWS API Gateway (Serverless)

```typescript
// serverless.yml
service: api-gateway

provider:
  name: aws
  runtime: nodejs18.x
  region: us-east-1
  environment:
    USER_SERVICE_URL: ${env:USER_SERVICE_URL, 'http://localhost:3000'}

functions:
  proxy:
    handler: gateway/handler.main
    events:
      - http:
          path: /api/users/{proxy+}
          method: ANY
          cors: true
      - http:
          path: /api/health
          method: ANY
          cors: true

plugins:
  - serverless-offline
```

```typescript
// gateway/handler.ts
import axios from 'axios';

export const main = async (event: any) => {
  const userServiceUrl = process.env.USER_SERVICE_URL;
  const path = event.path.replace('/api', '');

  try {
    const response = await axios({
      method: event.httpMethod.toLowerCase(),
      url: `${userServiceUrl}${path}`,
      data: event.body ? JSON.parse(event.body) : undefined,
      headers: event.headers,
    });

    return {
      statusCode: response.status,
      body: JSON.stringify(response.data),
      headers: {
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': '*',
      },
    };
  } catch (error: any) {
    return {
      statusCode: error.response?.status || 500,
      body: JSON.stringify({ message: error.message }),
      headers: { 'Access-Control-Allow-Origin': '*' },
    };
  }
};
```

## 7. Docker Compose - Full Stack

```yaml
# docker-compose.yml
version: '3.8'

services:
  # API Gateway
  api-gateway:
    build:
      context: .
      dockerfile: gateway/Dockerfile
    ports:
      - '8000:8000'
    environment:
      USER_SERVICE_URL: http://user-service:3000
    depends_on:
      - user-service
    networks:
      - api-network

  # User Service
  user-service:
    build:
      context: ./services/user-service
      dockerfile: Dockerfile
    ports:
      - '3000:3000'
    environment:
      NODE_ENV: production
      PORT: 3000
    networks:
      - api-network
    healthcheck:
      test: ['CMD', 'curl', '-f', 'http://localhost:3000/health']
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  api-network:
    driver: bridge
```

```dockerfile
# services/user-service/Dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY dist ./dist

EXPOSE 3000

CMD ["node", "dist/main.js"]
```

## 8. Load Balancing with Multiple Instances

```yaml
# docker-compose-scale.yml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - '8000:8000'
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - user-service-1
      - user-service-2
      - user-service-3

  user-service-1:
    build: ./services/user-service
    environment:
      PORT: 3000

  user-service-2:
    build: ./services/user-service
    environment:
      PORT: 3001

  user-service-3:
    build: ./services/user-service
    environment:
      PORT: 3002
```

## Quick Integration Checklist

- [ ] Choose gateway solution (Express, NestJS, Kong, Nginx, etc.)
- [ ] Configure service discovery (static URL or service registry)
- [ ] Add route mapping `/api/users → localhost:3000/users`
- [ ] Add health check endpoint monitoring
- [ ] Configure CORS headers if needed
- [ ] Add request/response logging
- [ ] Implement rate limiting (if needed)
- [ ] Set up error handling middleware
- [ ] Configure timeouts and retries
- [ ] Test all CRUD operations through gateway
- [ ] Monitor service health
- [ ] Document endpoint mapping

## Testing Through Gateway

```bash
# Create user through gateway
curl -X POST http://localhost:8000/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "firstName": "Test",
    "lastName": "User",
    "password": "password123"
  }'

# Health check
curl http://localhost:8000/api/health

# Get all users
curl http://localhost:8000/api/users
```

---

Choose the configuration that best fits your architecture and requirements!
