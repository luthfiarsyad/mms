# MMS API Postman Collection

This document explains how to use the Postman collection for the Money Management System (MMS) API.

## Importing the Collection

1. Open Postman
2. Click on "Import" in the top left
3. Select the `MMS_API.postman_collection.json` file
4. The collection will be imported with all endpoints organized by functionality

## Collection Structure

The collection is organized into the following folders:

### 1. Health Check
- **GET /health**: Check if the API server is running

### 2. Authentication
- **POST /api/v1/auth/register**: Register a new user account
- **POST /api/v1/auth/login**: Login with email and password to get access token

### 3. Users
- **POST /api/v1/users**: Create a new user
- **GET /api/v1/users**: Get list of all users
- **GET /api/v1/users/{userId}**: Get user details by ID
- **PUT /api/v1/users/{userId}**: Update user information by ID
- **DELETE /api/v1/users/{userId}**: Delete user by ID

### 4. Transactions
- **POST /api/v1/transactions**: Create a new transaction
- **GET /api/v1/transactions**: Get list of all transactions
- **GET /api/v1/transactions/{transactionId}**: Get transaction details by ID

## Environment Variables

The collection uses the following environment variables:

- `baseUrl`: Base URL of the API (default: `http://localhost:8080`)
- `accessToken`: JWT token obtained from login (to be set manually after login)
- `userId`: Sample user ID for testing (default: `1`)
- `transactionId`: Sample transaction ID for testing (default: `1`)

## Usage Instructions

### 1. Setup Environment
1. In Postman, click on the "Environments" tab
2. Create a new environment called "MMS API"
3. Add the variables mentioned above
4. Set `baseUrl` to your API server address

### 2. Authentication Flow
1. First, use the **Register** endpoint to create a new user
2. Then use the **Login** endpoint with the same credentials
3. Copy the `access_token` from the login response
4. Update the `accessToken` environment variable with this token

### 3. Using Protected Endpoints
All user and transaction endpoints require authentication:
1. Make sure you have a valid `accessToken` set in your environment
2. The collection automatically includes the Authorization header with Bearer token

## Sample Request Bodies

### Register Request
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Login Request
```json
{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Sample Transaction Request (when implemented)
```json
{
  "amount": 100.50,
  "description": "Sample transaction",
  "type": "expense"
}
```

## Implementation Status

**Currently Working:**
- Health Check endpoint
- User Registration
- User Login

**Not Yet Implemented (returns 501):**
- All User CRUD operations (except auth)
- All Transaction operations

## Running the API Server

To test this collection, make sure your MMS API server is running:

```bash
cd /path/to/your/mms/project
go run cmd/main.go
```

The server will start on `http://localhost:8080` by default.

## Expected Responses

### Successful Registration (201)
```json
{
  "id": 1,
  "email": "john.doe@example.com",
  "created_at": "2023-01-01T12:00:00Z"
}
```

### Successful Login (200)
```json
{
  "access_token": "your_paseto_token_here",
  "token_type": "bearer",
  "expires_in": 86400
}
```

### Health Check (200)
```json
{
  "status": "ok"
}
```

### Not Implemented (501)
```json
{
  "error": "not implemented"
}
```

## Troubleshooting

1. **Connection refused**: Make sure the API server is running on the correct port
2. **Invalid credentials**: Check that you're using the correct email and password
3. **Missing token**: Make sure to set the `accessToken` environment variable after login
4. **Token expired**: Login again to get a new token

## Security Notes

- The API uses PASETO V2 tokens for authentication
- Tokens expire after 24 hours
- In production, use HTTPS instead of HTTP
- The symmetric key should be loaded from environment variables in production