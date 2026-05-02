# User Microservice

A NestJS-based user microservice designed to be part of an API gateway architecture.

## Features

- ✅ **Health Check Endpoint** - Monitor service status
- ✅ **Create Users** - POST /users - Create new users with email validation
- ✅ **List All Users** - GET /users - Retrieve all users
- ✅ **Get User by ID** - GET /users/:id - Fetch specific user details
- ✅ **Update User** - PATCH /users/:id - Update user information
- ✅ **Delete User** - DELETE /users/:id - Remove users

## Prerequisites

- Node.js 18+
- npm

## Installation

```bash
npm install
```

## Running the Service

### Development Mode

```bash
npm run start:dev
```

The service will run on `http://localhost:3000` and automatically reload on file changes.

### Production Mode

```bash
npm run build
npm run start:prod
```

## API Endpoints

### Health Check

```http
GET /health
```

**Response:**

```json
{
  "status": "healthy",
  "service": "user-service",
  "timestamp": "2026-05-02T10:30:00.000Z"
}
```

### Create User

```http
POST /users
Content-Type: application/json

{
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "password": "securePassword123"
}
```

**Response (201 Created):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "createdAt": "2026-05-02T10:30:00.000Z",
  "updatedAt": "2026-05-02T10:30:00.000Z"
}
```

### List All Users

```http
GET /users
```

**Response (200 OK):**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "createdAt": "2026-05-02T10:30:00.000Z",
    "updatedAt": "2026-05-02T10:30:00.000Z"
  }
]
```

### Get User by ID

```http
GET /users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "createdAt": "2026-05-02T10:30:00.000Z",
  "updatedAt": "2026-05-02T10:30:00.000Z"
}
```

### Update User

```http
PATCH /users/550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json

{
  "firstName": "Jonathan",
  "lastName": "Smith"
}
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "Jonathan",
  "lastName": "Smith",
  "createdAt": "2026-05-02T10:30:00.000Z",
  "updatedAt": "2026-05-02T10:30:00.000Z"
}
```

### Delete User

```http
DELETE /users/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "firstName": "Jonathan",
  "lastName": "Smith",
  "createdAt": "2026-05-02T10:30:00.000Z",
  "updatedAt": "2026-05-02T10:30:00.000Z"
}
```

## Testing

### Run Unit Tests

```bash
npm run test
```

### Run E2E Tests

```bash
npm run test:e2e
```

### Run Tests with Coverage

```bash
npm run test:cov
```

## Error Responses

### 400 Bad Request

Invalid request format or missing required fields.

### 404 Not Found

```json
{
  "statusCode": 404,
  "message": "User with ID xxx not found",
  "error": "Not Found"
}
```

### 409 Conflict

```json
{
  "statusCode": 409,
  "message": "User with this email already exists",
  "error": "Conflict"
}
```

## Project Structure

```
src/
├── app.controller.ts          # Main app controller (health check)
├── app.service.ts             # Main app service
├── app.module.ts              # Root module
├── main.ts                     # Application entry point
└── users/
    ├── dto/
    │   ├── create-user.dto.ts       # Create user DTO
    │   ├── update-user.dto.ts       # Update user DTO
    │   └── user-response.dto.ts     # User response DTO
    ├── entities/
    │   └── user.entity.ts           # User entity definition
    ├── users.controller.ts      # User controller
    ├── users.service.ts         # User service (business logic)
    └── users.module.ts          # Users module
```

## Configuration

Set the PORT environment variable to run on a different port:

```bash
PORT=3001 npm run start:dev
```

Default port: **3000**

## API Gateway Integration

This microservice is designed to be fronted by an API gateway. The gateway should:

1. Route `/users*` requests to this service
2. Apply rate limiting and authentication policies
3. Handle request/response transformation if needed
4. Implement circuit breaking for resilience

Example gateway route configuration:

```
/api/users/* -> http://localhost:3000/users/*
/health -> http://localhost:3000/health
```

## Notes

- Currently uses in-memory storage (loses data on restart)
- For production use, integrate with a database (PostgreSQL, MongoDB, etc.)
- Passwords are stored in plain text (implement encryption for production)
- Add authentication/authorization middleware for security

## Future Enhancements

- [ ] Database integration (TypeORM with PostgreSQL)
- [ ] Password hashing (bcrypt)
- [ ] JWT authentication
- [ ] Role-based access control
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Request validation (class-validator)
- [ ] Logging and monitoring
- [ ] Error handling middleware
- [ ] Rate limiting
