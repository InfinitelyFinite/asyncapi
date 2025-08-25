# AsyncAPI - Asynchronous Report Generation Service

A Go-based microservices architecture for handling asynchronous report generation with JWT authentication, PostgreSQL database, and SQS message queuing.

## 🏗️ Architecture Overview

```mermaid
graph TB
    client[client]
    apiserver[apiserver]
    postgres[(postgres)]
    reports_sqs_queue[reports_sqs_queue]
    sqsworker[sqsworker]
    s3[(s3)]
    lozapi[LoZ API]
    
    client -->|"POST /auth/signup<br/>POST /auth/signin<br/>POST /auth/refresh"| apiserver
    client -->|"POST /reports<br/>GET /reports/{report_id}"| apiserver
    
    apiserver --> postgres
    apiserver -->|"SQS { report_id }"| reports_sqs_queue
    
    reports_sqs_queue --> sqsworker
    sqsworker --> lozapi
    sqsworker --> s3
    
    postgres -.->|"users<br/>refresh_tokens<br/>report_jobs"| postgres
    s3 -.->|"s3://api-reports"| s3
```

## 🔄 Request Flow

1. **Authentication**: Client registers/signs in via API Server, receives JWT tokens
2. **Report Creation**: Client submits report request, API server queues job in SQS
3. **Async Processing**: SQS Worker processes job, fetches data from LoZ API, uploads to S3
4. **Report Retrieval**: Client polls for status and gets download URL when ready

## 🚀 Features

- **JWT Authentication**: Secure access and refresh token implementation
- **Async Report Processing**: SQS-based job queuing and processing
- **File Storage**: S3 integration with presigned download URLs
- **Database**: PostgreSQL with migrations and proper schema design

## 📁 Project Structure

```
asyncapi/
├── apiserver/              # HTTP API server implementation
│   ├── handler.go          # HTTP request handlers (auth, reports)
│   ├── helpers.go          # Utility functions and error handling
│   ├── jwt.go              # JWT token generation and parsing
│   ├── middleware.go       # Authentication and logging middleware
│   └── server.go           # Server setup, routing, and lifecycle
├── cmd/
│   ├── apiserver/          # API server entry point
│   ├── awstest/           # AWS services testing utility
│   └── worker/            # SQS worker entry point
├── config/                # Environment-based configuration
├── db/
│   └── migrations/        # Database schema migrations
├── fixtures/              # Test utilities and database setup
├── reports/               # Report generation and processing
│   ├── builder.go         # Core report generation logic
│   ├── loz_client.go      # External API client
│   ├── sqs.go            # SQS message structures
│   └── worker.go         # SQS message consumer
├── store/                 # Data access layer
│   ├── db.go             # Database connection setup
│   ├── store.go          # Store aggregation
│   ├── users.go          # User repository
│   ├── refresh_tokens.go # Token management
│   └── reports.go        # Report data access
├── terraform/            # Infrastructure as Code
├── docker-compose.yml    # Local development environment
└── Makefile             # Development commands
```

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **Database**: PostgreSQL 15+
- **Message Queue**: Amazon SQS
- **File Storage**: Amazon S3
- **Authentication**: JWT with HS256
- **Testing**: Testify framework
- **Container**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Docker & Docker Compose
- AWS Account (for SQS and S3 services)

## 🚀 Getting Started

### 1. Clone and Setup
```bash
git clone <repository-url>
cd asyncapi
```

### 2. Environment Setup
Create `.envrc` file with required environment variables:
```bash
export DB_NAME=asyncapi
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export JWT_SECRET=your-super-secret-jwt-key
export DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/asyncapi?sslmode=disable
# ... other AWS and LocalStack configs
```

### 3. Start Services
```bash
# Start dependencies
docker-compose up -d

# Run migrations
make db_migrate

# Start API server
go run cmd/apiserver/main.go

# Start worker (in separate terminal)
go run cmd/worker/main.go
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./store/...
```

## 📚 API Endpoints

### Authentication
- `POST /auth/signup` - User registration
- `POST /auth/signin` - User login
- `POST /auth/refresh` - Refresh access token

### Reports
- `POST /reports` - Submit new report generation request
- `GET /reports/{report_id}` - Get report status and download URL

### Health Check
- `GET /ping` - Service health check

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(320) NOT NULL UNIQUE,
    hashed_password VARCHAR(96) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    hashed_token VARCHAR(500) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, hashed_token)
);
```

### Reports Table
```sql
CREATE TABLE reports (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    report_type VARCHAR NOT NULL,
    output_file_path VARCHAR,
    download_url VARCHAR,
    download_url_expires_at TIMESTAMPTZ,
    error_message VARCHAR,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, id)
);
```

## 🔧 Development Commands

```bash
# Database operations
make db_login          # Connect to database
make db_migrate        # Run migrations
make db_create_migration name=migration_name  # Create new migration

# Testing
go test ./...          # Run all tests
go test -v ./...       # Run tests with verbose output
```