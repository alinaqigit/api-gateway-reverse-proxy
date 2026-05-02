# User Service - API Testing Guide

Quick reference for testing the User Service endpoints using curl.

## Start the Service

```bash
# Development mode (with auto-restart)
npm run start:dev

# Production mode
npm run build
npm run start:prod
```

The service will be available at `http://localhost:3000`

## Health Check

```bash
curl http://localhost:3000/health
```

**Expected Response:**

```json
{
  "status": "healthy",
  "service": "user-service",
  "timestamp": "2026-05-02T10:30:00.000Z"
}
```

## Create User

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "password": "securePassword123"
  }'
```

**Expected Response (201 Created):**

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

## Get All Users

```bash
curl http://localhost:3000/users
```

**Expected Response:**

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

## Get Single User

Replace `{userID}` with actual user ID:

```bash
curl http://localhost:3000/users/550e8400-e29b-41d4-a716-446655440000
```

**Expected Response:**

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

## Update User

```bash
curl -X PATCH http://localhost:3000/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Jonathan",
    "lastName": "Smith"
  }'
```

**Expected Response:**

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

## Delete User

```bash
curl -X DELETE http://localhost:3000/users/550e8400-e29b-41d4-a716-446655440000
```

**Expected Response:**

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

## Using Postman/Insomnia

### Create Request Collection

1. **Create User**
   - `POST` http://localhost:3000/users
   - Body (JSON):
     ```json
     {
       "email": "user@example.com",
       "firstName": "First",
       "lastName": "Last",
       "password": "password"
     }
     ```

2. **Get All Users**
   - `GET` http://localhost:3000/users

3. **Get User by ID**
   - `GET` http://localhost:3000/users/{{userId}}

4. **Update User**
   - `PATCH` http://localhost:3000/users/{{userId}}
   - Body (JSON):
     ```json
     {
       "firstName": "Updated",
       "lastName": "Name"
     }
     ```

5. **Delete User**
   - `DELETE` http://localhost:3000/users/{{userId}}

## Error Scenarios

### Duplicate Email (409 Conflict)

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "firstName": "Jane",
    "lastName": "Doe",
    "password": "password"
  }'
```

**Response:**

```json
{
  "statusCode": 409,
  "message": "User with this email already exists",
  "error": "Conflict"
}
```

### User Not Found (404)

```bash
curl http://localhost:3000/users/invalid-id
```

**Response:**

```json
{
  "statusCode": 404,
  "message": "User with ID invalid-id not found",
  "error": "Not Found"
}
```

## Run Tests

```bash
# E2E Tests
npm run test:e2e

# Unit Tests
npm run test

# Watch Mode
npm run test:watch

# Coverage
npm run test:cov
```

## Frontend Integration Tips

When integrating with your API Gateway, configure routing like:

```
/api/users/* → http://localhost:3000/users/*
/api/health → http://localhost:3000/health
```

Example fetch call:

```javascript
// Create user
fetch('http://api-gateway:3000/api/users', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    firstName: 'John',
    lastName: 'Doe',
    password: 'password123',
  }),
})
  .then((res) => res.json())
  .then((data) => console.log(data));

// Get all users
fetch('http://api-gateway:3000/api/users')
  .then((res) => res.json())
  .then((data) => console.log(data));

// Update user
fetch('http://api-gateway:3000/api/users/userId', {
  method: 'PATCH',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    firstName: 'Updated',
    lastName: 'Name',
  }),
})
  .then((res) => res.json())
  .then((data) => console.log(data));

// Delete user
fetch('http://api-gateway:3000/api/users/userId', {
  method: 'DELETE',
})
  .then((res) => res.json())
  .then((data) => console.log(data));
```
