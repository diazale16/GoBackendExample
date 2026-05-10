# Go + PostgreSQL + S3 Demo Application

Educational project demonstrating integration with PostgreSQL and S3-compatible storage, designed to facilitate migration from Supabase to a self-hosted VPS.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Router (Gin)                            │
│                     HTTP Request Handling                       │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Controllers                               │
│            HTTP Input/Output + Validation                       │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Services                                 │
│                  Business Logic Layer                           │
└─────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repositories                               │
│                    Data Access Layer                            │
└─────────────────────────────────────────────────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    ▼                         ▼
        ┌───────────────────┐    ┌───────────────────┐
        │    PostgreSQL     │    │    S3 Storage    │
        │   (Metadata)      │    │   (Files)        │
        └───────────────────┘    └───────────────────┘
```

## Features

- **GORM with AutoMigrate**: Schema management directly from code
- **PostgreSQL**: Primary database with connection pooling
- **S3-Compatible Storage**: Works with Supabase Storage, MinIO, AWS S3
- **RESTful API**: Clean endpoints for CRUD operations
- **File Upload**: Multipart handling with metadata storage
- **Soft Deletes**: GORM's built-in DeletedAt support

## Project Structure

```
cmd/
  main.go                 # Application entry point

internal/
  config/
    config.go            # Environment configuration
  database/
    database.go          # PostgreSQL connection + AutoMigrate
  storage/
    storage.go           # S3 client implementation

models/
  user.go                # User model + DTOs
  document.go            # Document model + DTOs

repositories/
  user.go                # User data access
  document.go            # Document data access

services/
  user.go                # User business logic
  document.go            # Document business logic

controllers/
  user.go                # User HTTP handlers
  document.go            # Document HTTP handlers

routes/
  routes.go              # Gin router setup

.env.example             # Environment variables template
```

## Installation

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- S3-compatible storage (Supabase, MinIO, AWS)

### Steps

1. Clone the repository:
```bash
git clone <repository-url>
cd supabase-migration-demo
```

2. Install dependencies:
```bash
go mod tidy
```

3. Copy and configure environment:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run cmd/main.go
```

## Configuration

Create a `.env` file with the following variables:

```env
# PostgreSQL Connection
# Format: postgresql://user:password@host:port/database?options
DATABASE_URL=postgresql://postgres:password@localhost:5432/mydb?sslmode=disable

# S3-Compatible Storage
# For Supabase: https://your-project.supabase.co/storage/v1/object
# For MinIO: http://localhost:9000
# For AWS S3: leave empty or use https://s3.amazonaws.com
S3_ENDPOINT=https://your-project.supabase.co/storage/v1/object
S3_REGION=us-east-1
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_BUCKET=documents

# Server
SERVER_PORT=8080
```

### Supabase Storage Setup

1. Create a storage bucket named `documents` in your Supabase project
2. Set bucket as public (for simplicity) or configure signed URLs
3. Use the storage URL format: `https://project-ref.supabase.co/storage/v1/object`

### MinIO Setup

```bash
# Run MinIO container
docker run -d -p 9000:9000 -p 9001:9001 --name minio \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data

# Create bucket
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/documents
mc anonymous set public local/documents
```

## API Endpoints

### Health Check

```bash
GET /health
```

### Users

```bash
# Create user
POST /users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}

# Get all users
GET /users

# Get user by ID
GET /users/:id

# Delete user
DELETE /users/:id
```

### Documents

```bash
# Upload document
POST /documents/upload
Content-Type: multipart/form-data

file: <binary file>
user_id: <uuid>

# Get all documents
GET /documents

# Get document by ID
GET /documents/:id

# Delete document
DELETE /documents/:id
```

## Examples

### Create a User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Upload a Document

```bash
curl -X POST http://localhost:8080/documents/upload \
  -F "file=@/path/to/document.pdf" \
  -F "user_id=550e8400-e29b-41d4-a716-446655440000"
```

Response:
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "file_name": "document.pdf",
  "storage_key": "uploads/2024/01/15/660e8400_document.pdf",
  "mime_type": "application/pdf",
  "url": "https://s3.amazonaws.com/documents/uploads/2024/01/15/660e8400_document.pdf",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

## Understanding AutoMigrate

### What AutoMigrate Does

```go
db.AutoMigrate(&User{}, &Document{})
```

AutoMigrate automatically handles:

- **Table Creation**: Creates tables if they don't exist
- **Column Addition**: Adds new columns that don't exist
- **Index Management**: Creates indexes defined in models
- **Constraint Addition**: Adds constraints like foreign keys
- **Type Changes**: Attempts to modify column types (with limitations)

### What AutoMigrate Does NOT Do

- **Column Removal**: Never removes columns (safety)
- **Column Renaming**: Can't detect renames (will create new)
- **Data Migration**: Doesn't migrate existing data
- **Index Removal**: Never drops unused indexes
- **Constraint Changes**: Can't modify existing constraints

### AutoMigrate Limitations

1. **Not Production-Ready**: AutoMigrate can lose data in some cases
2. **No Schema Downsgrades**: Can't rollback changes
3. **Migration State Unknown**: Doesn't track what's been applied
4. **Concurrent Issues**: Not safe to run multiple instances

### Best Practices for Development

```go
// Enable SQL logging to see what AutoMigrate does
config := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
}

// Check current schema before migrating
db.AutoMigrate(&User{}, &Document{})

// Review generated SQL in logs
// CREATE TABLE IF NOT EXISTS... for tables
// CREATE INDEX... for indexes
```

### Production Alternatives

For production, consider using dedicated migration tools:

- **golang-migrate**: Versioned SQL migrations
- ** goose**: Migration tool with up/down support
- **sql-migrate**: Migration management

Example with golang-migrate:
```bash
# Create migration
migrate create -ext sql -dir migrations -seq create_users

# Write migration file
# 000001_create_users.up.sql
# 000001_create_users.down.sql

# Run migrations
migrate -path migrations -database $DATABASE_URL up
```

## Database Schema

### User Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

### Document Table

```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    storage_key VARCHAR(512) NOT NULL,
    mime_type VARCHAR(100),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_deleted_at ON documents(deleted_at);
```

### Relationships

- User has many Documents (1:N)
- Document belongs to User
- CASCADE delete: Deleting a user removes all their documents

## Storage Architecture

### Metadata vs Binary Storage

```
┌────────────────────────┐
│   PostgreSQL (Metadata) │
├────────────────────────┤
│ - id                   │
│ - user_id              │
│ - file_name            │
│ - storage_key          │
│ - mime_type            │
│ - created_at           │
└────────────────────────┘
          │
          │ storage_key references
          ▼
┌────────────────────────┐
│   S3 Storage (Binary)   │
├────────────────────────┤
│ - uploads/             │
│   └── 2024/01/15/      │
│       └── uuid_file.txt│
└────────────────────────┘
```

### File Organization

Files are stored with date-based prefixing:
```
uploads/YYYY/MM/DD/{uuid}_{original_filename}
```

Benefits:
- Natural partitioning for better performance
- Unique filenames prevent collisions
- Date-based cleanup possible

### Deletion Strategy

When deleting a document:
1. Delete metadata from PostgreSQL (transactional)
2. Delete file from S3 (best effort)
3. Log any S3 errors but don't fail the operation

This ensures data consistency in the database while cleaning up storage.

## S3 Configuration Examples

### Supabase Storage

```env
S3_ENDPOINT=https://xyztpq.supabase.co/storage/v1/object
S3_REGION=us-east-1
S3_ACCESS_KEY=your-anon-key-or-service-key
S3_SECRET_KEY=your-service-key
S3_BUCKET=documents
```

### MinIO (Local)

```env
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=documents
```

### AWS S3

```env
S3_ENDPOINT=
S3_REGION=us-east-1
S3_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
S3_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
S3_BUCKET=my-documents-bucket
```

## Error Handling

The application handles errors gracefully:

- **400 Bad Request**: Invalid input, validation errors
- **404 Not Found**: Resource doesn't exist
- **500 Internal Server Error**: Unexpected errors

All errors return JSON:
```json
{
  "error": "descriptive error message"
}
```

## Logging

SQL logging is enabled by default, showing:
- All queries executed
- Query parameters
- Execution time
- Results

This helps understand GORM behavior and debug issues.

## Migration Guide: From Supabase to VPS

### Step 1: Export Data

```bash
# Export users
pg_dump -h $SUPABASE_HOST -U $SUPABASE_USER -d postgres -t users > users.sql

# Export documents metadata
pg_dump -h $SUPABASE_HOST -U $SUPABASE_USER -d postgres -t documents > documents.sql
```

### Step 2: Setup PostgreSQL on VPS

```bash
# Install PostgreSQL
sudo apt install postgresql

# Create database and user
sudo -u postgres psql
CREATE DATABASE myapp;
CREATE USER myappuser WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE myapp TO myappuser;
```

### Step 3: Migrate Data

```bash
# Import data
psql -h localhost -U myappuser -d myapp < users.sql
psql -h localhost -U myappuser -d myapp < documents.sql
```

### Step 4: Setup MinIO on VPS

```bash
# Download MinIO
wget https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio

# Run MinIO
MINIO_ROOT_USER=minioadmin MINIO_ROOT_PASSWORD=minioadmin ./minio server /data &
```

### Step 5: Migrate Files

```bash
# Use AWS CLI to sync from Supabase to MinIO
aws s3 sync s3://supabase-bucket s3://my-bucket \
  --endpoint-url http://localhost:9000 \
  --source-region us-east-1 \
  --region us-east-1
```

### Step 6: Update Configuration

```env
DATABASE_URL=postgresql://myappuser:password@localhost:5432/myapp
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=documents
```

## Testing

Run the application and test endpoints:

```bash
# Start server
go run cmd/main.go

# In another terminal

# Create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com"}'

# Get users
curl http://localhost:8080/users

# Upload file
curl -X POST http://localhost:8080/documents/upload \
  -F "file=@README.md" \
  -F "user_id=<user-uuid-from-previous-response>"

# Get documents
curl http://localhost:8080/documents

# Delete document
curl -X DELETE http://localhost:8080/documents/<doc-uuid>

# Delete user (cascades to documents)
curl -X DELETE http://localhost:8080/users/<user-uuid>
```

## License

MIT