# AsyncAPI - Asynchronous Report Generation Service

A Go-based microservices architecture for handling asynchronous report generation with JWT authentication, PostgreSQL database, and SQS message queuing.

## 🏗️ Architecture Overview

```mermaid
graph TB
    %% Client Layer
    Client[👤 Client Application]
    
    %% API Gateway / Load Balancer
    LB[🌐 Load Balancer]
    
    %% API Server
    API[🚀 API Server<br/>- JWT Auth<br/>- Request Validation<br/>- Business Logic]
    
    %% Message Queue
    SQS[📨 Amazon SQS<br/>reports_sqs_queue]
    
    %% Worker Service
    Worker[⚙️ SQS Worker<br/>- Report Processing<br/>- File Generation<br/>- Status Updates]
    
    %% Database
    DB[(🗄️ PostgreSQL<br/>- users<br/>- refresh_tokens<br/>- reports)]
    
    %% File Storage
    S3[☁️ Amazon S3<br/>s3://api-reports<br/>- Generated Files<br/>- Presigned URLs]
    
    %% External API
    LoZ[🔗 LoZ API<br/>External Data Source]
    
    %% Request Flow
    Client -->|"POST /auth/signup<br/>POST /auth/signin<br/>POST /auth/refresh"| LB
    Client -->|"POST /reports<br/>GET /reports/{id}"| LB
    
    LB --> API
    
    API -->|"Store User Data<br/>Validate Tokens<br/>Update Report Status"| DB
    API -->|"Queue Report Job<br/>{report_id, user_id, type}"| SQS
    
    SQS -->|"Consume Messages<br/>Process Reports"| Worker
    
    Worker -->|"Fetch Data"| LoZ
    Worker -->|"Update Job Status<br/>Store Metadata"| DB
    Worker -->|"Upload Generated Files<br/>Create Download URLs"| S3
    
    %% Styling
    classDef client fill:#e1f5fe
    classDef api fill:#f3e5f5
    classDef storage fill:#e8f5e8
    classDef queue fill:#fff3e0
    classDef external fill:#fce4ec
    
    class Client client
    class API,LB api
    class DB,S3 storage
    class SQS,Worker queue
    class LoZ external
```

### 🔄 Request Flow

1. **Authentication Flow**:
   - Client registers/signs in via API Server
   - JWT tokens (access + refresh) issued and stored
   - Subsequent requests authenticated via JWT middleware

2. **Report Generation Flow**:
   - Client submits report request to API Server
   - API Server validates request and queues job in SQS
   - API Server immediately returns report ID to client
   - SQS Worker processes job asynchronously:
     - Fetches required data from LoZ API
     - Generates report files
     - Uploads to S3 and creates presigned download URLs
     - Updates job status in database

3. **Report Retrieval Flow**:
   - Client polls API Server for report status
   - API Server returns current status and download URL (if ready)

## 🚀 Features

### Authentication & Authorization
- **JWT-based Authentication**: Secure access and refresh token implementation
- **User Registration**: Email-based signup with bcrypt password hashing  
- **Token Management**: Automatic token refresh and validation middleware

### Asynchronous Report Processing
- **Report Generation**: Async processing of report requests via SQS
- **Job Tracking**: Real-time status updates (pending, processing, completed, failed)
- **File Storage**: Generated reports stored in S3 with presigned download URLs
- **Error Handling**: Comprehensive error tracking and user notification

### Database Design
- **PostgreSQL**: Robust relational database with proper foreign key constraints
- **Migration Support**: Database versioning with golang-migrate
- **Connection Pooling**: Optimized database connections using sqlx

## 📁 Project Structure

```
asyncapi/
├── apiserver/          # HTTP API server implementation
│   ├── handler.go      # HTTP request handlers
│   ├── helpers.go      # Utility functions and middleware
│   ├── jwt.go          # JWT token management
│   ├── middleware.go   # HTTP middleware (logging, auth)
│   └── server.go       # Server setup and routing
├── cmd/
│   └── apiserver/      # Application entry points
├── config/             # Configuration management
├── db/
│   └── migrations/     # Database migration files
├── fixtures/           # Test utilities and fixtures  
├── store/              # Data access layer
│   ├── db.go          # Database connection
│   ├── store.go       # Store initialization
│   └── users.go       # User repository
├── docker-compose.yml  # Local development environment
├── Makefile           # Development commands
└── go.mod             # Go module dependencies
```

## 🛠️ Technology Stack

- **Language**: Go 1.23+
- **Database**: PostgreSQL 15+
- **Message Queue**: Amazon SQS
- **File Storage**: Amazon S3
- **Authentication**: JWT with HS256 signing
- **Password Hashing**: bcrypt
- **Testing**: Testify framework
- **Migrations**: golang-migrate
- **Container**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.23 or later
- PostgreSQL 15+
- Docker & Docker Compose (for local development)
- AWS Account (for SQS and S3 services)

## 🚀 Getting Started

### 1. Clone the Repository
```bash
git clone <repository-url>
cd asyncapi
```

### 2. Environment Setup
Create a `.envrc` file in the project root:
```bash
export DB_NAME=asyncapi
export DB_HOST=localhost
export DB_PORT=5432
export DB_PORT_TEST=5433
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export ENV=dev
export PROJECT_ROOT=/path/to/your/project
export APISERVER_PORT=8080
export APISERVER_HOST=localhost
export JWT_SECRET=your-super-secret-jwt-key
export DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/asyncapi?sslmode=disable
```

### 3. Start Dependencies
```bash
# Start PostgreSQL databases (main and test)
docker-compose up -d

# Run database migrations
make db_migrate
```

### 4. Install Dependencies
```bash
go mod tidy
```

### 5. Run the Application
```bash
# Start the API server
go run cmd/apiserver/main.go
```

The API server will start on `http://localhost:8080`

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
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

## 🌟 Key Features Implementation

### JWT Authentication
- Access tokens (15 minutes expiry)
- Refresh tokens (30 days expiry)  
- Secure token validation middleware
- Automatic token refresh flow

### Async Processing
- SQS-based message queuing for scalable report processing
- Job status tracking with timestamps
- Error handling and retry mechanisms
- S3 integration for report file storage

### Database Design
- UUID primary keys for security
- Proper foreign key relationships
- Timestamp tracking for audit trails
- Optimized queries with proper indexing

---

## 🚧 Status

**This project is currently under construction.** 

While the core authentication system and database schema are implemented and tested, the following components are still in development:

- SQS Worker implementation for async report processing
- S3 integration for file storage and presigned URLs  
- Report generation logic and job queue processing
- Additional API endpoints for report management
- Production deployment configurations
- Comprehensive API documentation

Check back soon for updates as I continue building out the complete async report generation system!