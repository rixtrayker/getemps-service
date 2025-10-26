# GetEmpStatus Web Service - Implementation Plan

## Project Overview

A RESTful web service that retrieves employee information, processes salary data, and calculates employee status based on defined business logic.

**Estimated Time:** 10-16 hours (core + bonus features)

---

## Table of Contents

1. [Technology Stack](#technology-stack)
2. [Project Structure](#project-structure)
3. [Database Design](#database-design)
4. [API Specification](#api-specification)
5. [Business Logic](#business-logic)
6. [Core Components](#core-components)
7. [Bonus Features](#bonus-features)
8. [Docker Setup](#docker-setup)
9. [Implementation Steps](#implementation-steps)
10. [Testing Strategy](#testing-strategy)

---

## Technology Stack

### Core Technologies
- **Language:** Go 1.21+
- **Database:** PostgreSQL 15+ (or MySQL 8+)
- **Container:** Docker & Docker Compose

### Go Packages

#### HTTP Framework & Router
```go
github.com/gin-gonic/gin v1.9.1  // HTTP web framework
```

#### Database
```go
github.com/lib/pq v1.10.9                    // PostgreSQL driver
// OR
github.com/go-sql-driver/mysql v1.7.1       // MySQL driver

github.com/jmoiron/sqlx v1.3.5              // Extensions to database/sql
```

#### Validation
```go
github.com/go-ozzo/ozzo-validation/v4 v4.3.0  // Validation framework
```

#### Configuration
```go
github.com/spf13/viper v1.17.0              // Configuration management
github.com/joho/godotenv v1.5.1             // .env file support
```

#### Caching (Bonus)
```go
github.com/patrickmn/go-cache v2.1.0+incompatible  // In-memory cache
// OR
github.com/go-redis/redis/v8 v8.11.5              // Redis client
```

#### Logging (Bonus)
```go
github.com/sirupsen/logrus v1.9.3           // Structured logger
// OR
go.uber.org/zap v1.26.0                     // High-performance logger
```

#### Retry Mechanism (Bonus)
```go
github.com/avast/retry-go/v4 v4.5.0         // Retry library
```

#### Authentication (Bonus)
```go
github.com/golang-jwt/jwt/v5 v5.0.0         // JWT tokens
```

#### Testing
```go
github.com/stretchr/testify v1.8.4          // Testing toolkit
github.com/DATA-DOG/go-sqlmock v1.5.0       // SQL mock
```

---

## Project Structure

```
getemps-service/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── models/
│   │   ├── user.go                # User entity
│   │   ├── salary.go              # Salary entity
│   │   ├── employee_info.go       # EmpInfo response model
│   │   └── request.go             # Request models
│   ├── validator/
│   │   └── validator.go           # Input validation logic
│   ├── repository/
│   │   ├── interface.go           # Repository interfaces
│   │   ├── user_repository.go     # User data access
│   │   └── salary_repository.go   # Salary data access
│   ├── service/
│   │   ├── process_status.go      # Main business logic
│   │   └── salary_calculator.go   # Salary calculations
│   ├── handler/
│   │   └── employee_handler.go    # HTTP handlers
│   ├── middleware/
│   │   ├── auth.go                # Token validation (bonus)
│   │   ├── logger.go              # Request logging
│   │   └── recovery.go            # Panic recovery
│   ├── cache/                     # Cache layer (bonus)
│   │   └── cache.go
│   ├── logger/                    # Custom logger (bonus)
│   │   └── logger.go
│   └── database/
│       ├── connection.go          # DB connection
│       └── retry.go               # Retry mechanism (bonus)
├── migrations/
│   ├── 001_create_tables.sql      # Database schema
│   └── 002_insert_data.sql        # Sample data
├── docs/
│   └── api/
│       └── postman_collection.json
├── docker/
│   └── Dockerfile
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## Database Design

### Database Schema Files

**Location:** `migrations/001_create_tables.sql`
**Location:** `migrations/002_insert_data.sql`

### Schema Details

#### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    national_number VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_national_number ON users(national_number);
CREATE INDEX idx_is_active ON users(is_active);
```

#### Salaries Table
```sql
CREATE TABLE salaries (
    id SERIAL PRIMARY KEY,
    year INT NOT NULL,
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    salary DECIMAL(10, 2) NOT NULL CHECK (salary >= 0),
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_year_month UNIQUE (user_id, year, month)
);

CREATE INDEX idx_user_id ON salaries(user_id);
CREATE INDEX idx_year_month ON salaries(year, month);
```

#### Logs Table (Bonus Feature)
```sql
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    level VARCHAR(20) NOT NULL,  -- INFO, ERROR, WARN, DEBUG
    message TEXT NOT NULL,
    context JSONB,               -- Additional context data
    endpoint VARCHAR(255),
    request_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_level ON logs(level);
CREATE INDEX idx_created_at ON logs(created_at);
CREATE INDEX idx_endpoint ON logs(endpoint);
```

---

## API Specification

### Endpoint

**URL:** `POST /api/GetEmpStatus`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <token>  (Bonus feature)
```

### Request Body

```json
{
  "NationalNumber": "NAT1001"
}
```

### Response - Success (200 OK)

```json
{
  "id": 1,
  "username": "jdoe",
  "nationalNumber": "NAT1001",
  "email": "jdoe@example.com",
  "phone": "0791111111",
  "isActive": true,
  "salaryDetails": {
    "averageSalary": 1432.50,
    "highestSalary": 1760.00,
    "sumOfSalaries": 7162.50
  },
  "status": "RED",
  "lastUpdated": "2025-10-26T14:30:00Z"
}
```

### Response - Error Cases

#### Invalid National Number (404)
```json
{
  "error": "Invalid National Number"
}
```

#### User Not Active (406)
```json
{
  "error": "User is not Active"
}
```

#### Insufficient Data (422)
```json
{
  "error": "INSUFFICIENT_DATA"
}
```

#### Unauthorized (401) - Bonus
```json
{
  "error": "Unauthorized - Invalid or missing token"
}
```

---

## Business Logic

### Validation Rules

1. **National Number Validation:**
   - Required field
   - Non-empty string
   - Must exist in database

2. **User Active Check:**
   - User must have `is_active = true`

3. **Minimum Data Requirement:**
   - User must have at least 3 salary records

### Salary Adjustments

Applied in order before calculations:

#### 1. Seasonal Adjustments
```go
// December: +10% holiday bonus
if month == 12 {
    adjustedSalary = salary * 1.10
}

// Summer months (June, July, August): -5% deduction
if month == 6 || month == 7 || month == 8 {
    adjustedSalary = salary * 0.95
}
```

#### 2. Tax Deduction
```go
// If total salary > 10,000: apply 7% tax
totalSalary := sum(adjustedSalaries)
if totalSalary > 10000 {
    // Apply 7% deduction to EACH salary
    for i := range adjustedSalaries {
        adjustedSalaries[i] = adjustedSalaries[i] * 0.93
    }
}
```

### Status Determination

After all adjustments and calculations:

```go
averageSalary := sum(adjustedSalaries) / count(adjustedSalaries)

if averageSalary > 2000 {
    status = "GREEN"
} else if averageSalary == 2000 {
    status = "ORANGE"
} else {
    status = "RED"
}
```

### Calculation Example

**Scenario:** User NAT1001 with salaries:
- Jan 2025: 1200
- Feb 2025: 1300
- Mar 2025: 1400
- May 2025: 1500
- Jun 2025: 1600 (summer month)

**Step 1 - Seasonal Adjustments:**
- Jan: 1200 (no adjustment)
- Feb: 1300 (no adjustment)
- Mar: 1400 (no adjustment)
- May: 1500 (no adjustment)
- Jun: 1600 * 0.95 = 1520 (summer deduction)

**Step 2 - Calculate Total:**
- Total = 1200 + 1300 + 1400 + 1500 + 1520 = 6920

**Step 3 - Tax Check:**
- Total (6920) < 10000, so no tax deduction

**Step 4 - Final Calculations:**
- Average: 6920 / 5 = 1384
- Highest: 1520
- Sum: 6920
- Status: RED (average < 2000)

---

## Core Components

### 1. Validator (`internal/validator/validator.go`)

**Responsibilities:**
- Validate request payload structure
- Validate national number format
- Return user-friendly error messages

**Key Functions:**
```go
type EmployeeRequest struct {
    NationalNumber string `json:"NationalNumber"`
}

func (r EmployeeRequest) Validate() error {
    return validation.ValidateStruct(&r,
        validation.Field(&r.NationalNumber,
            validation.Required.Error("National number is required"),
            validation.Length(3, 50).Error("Invalid national number format"),
        ),
    )
}
```

### 2. DataAccess / Repository (`internal/repository/`)

**Responsibilities:**
- All database interactions
- Execute queries and stored procedures
- Handle database transactions
- Implement retry mechanism (bonus)

**Key Interfaces:**
```go
type UserRepository interface {
    GetByNationalNumber(ctx context.Context, nationalNumber string) (*models.User, error)
}

type SalaryRepository interface {
    GetByUserID(ctx context.Context, userID int64) ([]models.Salary, error)
    CountByUserID(ctx context.Context, userID int64) (int, error)
}
```

### 3. EmpInfo Model (`internal/models/employee_info.go`)

**Responsibilities:**
- Represent employee entity with all attributes
- Structure response data

**Structure:**
```go
type EmployeeInfo struct {
    ID             int64          `json:"id"`
    Username       string         `json:"username"`
    NationalNumber string         `json:"nationalNumber"`
    Email          string         `json:"email"`
    Phone          string         `json:"phone"`
    IsActive       bool           `json:"isActive"`
    SalaryDetails  SalaryDetails  `json:"salaryDetails"`
    Status         string         `json:"status"`
    LastUpdated    time.Time      `json:"lastUpdated"`
}

type SalaryDetails struct {
    AverageSalary float64 `json:"averageSalary"`
    HighestSalary float64 `json:"highestSalary"`
    SumOfSalaries float64 `json:"sumOfSalaries"`
}
```

### 4. ProcessStatus Service (`internal/service/process_status.go`)

**Responsibilities:**
- Orchestrate main business logic flow
- Coordinate between repositories
- Apply business rules
- Handle error cases
- Return formatted response

**Key Functions:**
```go
type ProcessStatusService struct {
    userRepo    repository.UserRepository
    salaryRepo  repository.SalaryRepository
    cache       cache.Cache // bonus
    logger      *logrus.Logger // bonus
}

func (s *ProcessStatusService) GetEmployeeStatus(ctx context.Context, nationalNumber string) (*models.EmployeeInfo, error)
```

**Logic Flow:**
1. Check cache (bonus)
2. Validate and fetch user
3. Check if user is active
4. Fetch and validate salary records count
5. Calculate adjusted salaries
6. Apply tax if needed
7. Calculate statistics
8. Determine status
9. Cache result (bonus)
10. Return response

---

## Bonus Features

### 1. Cache Handler (`internal/cache/cache.go`)

**Purpose:** Reduce database calls for frequently accessed employee data

**Implementation:**
```go
type Cache interface {
    Get(key string) (*models.EmployeeInfo, bool)
    Set(key string, value *models.EmployeeInfo, duration time.Duration)
    Delete(key string)
}

// Use in-memory cache (go-cache) or Redis
// Cache key: "emp_status:{nationalNumber}"
// TTL: 5 minutes (configurable)
```

**Strategy:**
- Cache complete employee status response
- Invalidate on data updates
- Use national number as cache key

### 2. Custom Logger (`internal/logger/logger.go`)

**Purpose:** Centralized logging to database and console

**Features:**
- Log levels: DEBUG, INFO, WARN, ERROR
- Store logs in database
- Include request context
- Track API calls and errors

**Implementation:**
```go
type Logger struct {
    db     *sqlx.DB
    logger *logrus.Logger
}

func (l *Logger) Info(ctx context.Context, message string, fields map[string]interface{})
func (l *Logger) Error(ctx context.Context, message string, err error, fields map[string]interface{})
```

**Database Logging:**
```go
INSERT INTO logs (level, message, context, endpoint, request_id, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())
```

### 3. DB Retry Mechanism (`internal/database/retry.go`)

**Purpose:** Retry failed database queries up to 3 times

**Implementation:**
```go
import "github.com/avast/retry-go/v4"

func ExecuteWithRetry(operation func() error) error {
    return retry.Do(
        operation,
        retry.Attempts(3),
        retry.Delay(100 * time.Millisecond),
        retry.DelayType(retry.BackOffDelay),
        retry.OnRetry(func(n uint, err error) {
            log.Printf("Retry attempt %d due to error: %v", n, err)
        }),
    )
}
```

**Usage:**
```go
var user *models.User
err := ExecuteWithRetry(func() error {
    var queryErr error
    user, queryErr = repo.GetByNationalNumber(ctx, nationalNumber)
    return queryErr
})
```

### 4. API Token Implementation (`internal/middleware/auth.go`)

**Purpose:** Validate API token before processing requests

**Implementation:**
```go
func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Unauthorized - Missing token"})
            c.Abort()
            return
        }
        
        // Extract Bearer token
        token := strings.TrimPrefix(authHeader, "Bearer ")
        
        // Validate token (simple or JWT)
        if !isValidToken(token, secretKey) {
            c.JSON(401, gin.H{"error": "Unauthorized - Invalid token"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

**Token Options:**
1. Simple API key validation
2. JWT token with claims
3. OAuth2 integration

---

## Docker Setup

### Dockerfile (`docker/Dockerfile`)

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

### Docker Compose (`docker-compose.yml`)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: getemps_db
    environment:
      POSTGRES_DB: getemps_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
    networks:
      - getemps_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: getemps_redis
    ports:
      - "6379:6379"
    networks:
      - getemps_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  api:
    build:
      context: .
      dockerfile: docker/Dockerfile
    container_name: getemps_api
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres123
      - DB_NAME=getemps_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - API_SECRET_KEY=your-secret-key-here
      - GIN_MODE=release
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - getemps_network
    restart: unless-stopped

volumes:
  postgres_data:

networks:
  getemps_network:
    driver: bridge
```

### Environment Configuration (`.env.example`)

```bash
# Application
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=getemps_db
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5

# Redis (optional - for caching)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Cache
CACHE_ENABLED=true
CACHE_TTL=300  # seconds

# Logging
LOG_LEVEL=info
LOG_TO_DB=true

# Security
API_SECRET_KEY=your-super-secret-key-change-in-production
TOKEN_REQUIRED=true

# Retry Configuration
DB_RETRY_ATTEMPTS=3
DB_RETRY_DELAY=100  # milliseconds
```

---

## Implementation Steps

### Phase 1: Project Setup (1 hour)

1. **Initialize Go Module**
   ```bash
   mkdir getemps-service && cd getemps-service
   go mod init github.com/yourusername/getemps-service
   ```

2. **Install Dependencies**
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/lib/pq
   go get github.com/jmoiron/sqlx
   go get github.com/go-ozzo/ozzo-validation/v4
   go get github.com/spf13/viper
   go get github.com/joho/godotenv
   go get github.com/sirupsen/logrus
   go get github.com/patrickmn/go-cache
   go get github.com/avast/retry-go/v4
   go get github.com/stretchr/testify
   ```

3. **Create Project Structure**
   ```bash
   mkdir -p cmd/api internal/{config,models,validator,repository,service,handler,middleware,cache,logger,database} migrations docs/api docker
   ```

4. **Create Configuration Files**
   - `.env.example`
   - `.gitignore`
   - `docker-compose.yml`

### Phase 2: Database Setup (1 hour)

1. **Create Migration Files**
   - `migrations/001_create_tables.sql`
   - `migrations/002_insert_data.sql`

2. **Setup Database Connection**
   - Implement `internal/database/connection.go`
   - Connection pooling
   - Health checks

3. **Test Database Connection**

### Phase 3: Core Models (1 hour)

1. **Implement Models**
   - `internal/models/user.go`
   - `internal/models/salary.go`
   - `internal/models/employee_info.go`
   - `internal/models/request.go`

2. **Define Structures**
   - Request/Response DTOs
   - Database entities

### Phase 4: Validation Layer (0.5 hour)

1. **Implement Validator**
   - `internal/validator/validator.go`
   - Request validation rules
   - Custom validation functions

### Phase 5: Repository Layer (2 hours)

1. **Define Interfaces**
   - `internal/repository/interface.go`

2. **Implement Repositories**
   - `internal/repository/user_repository.go`
   - `internal/repository/salary_repository.go`

3. **Write Repository Tests**
   - Mock database interactions
   - Test error cases

### Phase 6: Business Logic (3 hours)

1. **Implement Salary Calculator**
   - `internal/service/salary_calculator.go`
   - Seasonal adjustments
   - Tax calculations
   - Status determination

2. **Implement Process Status Service**
   - `internal/service/process_status.go`
   - Orchestration logic
   - Error handling
   - Business rules

3. **Write Unit Tests**
   - Test salary calculations
   - Test business rules
   - Test edge cases

### Phase 7: HTTP Layer (1.5 hours)

1. **Implement Handlers**
   - `internal/handler/employee_handler.go`
   - Request parsing
   - Response formatting
   - Error responses

2. **Setup Router**
   - Route definitions
   - Middleware registration

3. **Implement Basic Middleware**
   - `internal/middleware/logger.go`
   - `internal/middleware/recovery.go`

### Phase 8: Main Application (0.5 hour)

1. **Implement Entry Point**
   - `cmd/api/main.go`
   - Configuration loading
   - Dependency injection
   - Server startup

2. **Graceful Shutdown**

### Phase 9: Bonus Features (3-4 hours)

1. **Cache Implementation** (1 hour)
   - `internal/cache/cache.go`
   - In-memory or Redis
   - Integration with service layer

2. **Custom Logger** (1 hour)
   - `internal/logger/logger.go`
   - Database logging
   - Log rotation

3. **Retry Mechanism** (0.5 hour)
   - `internal/database/retry.go`
   - Integration with repositories

4. **API Token Authentication** (1 hour)
   - `internal/middleware/auth.go`
   - Token validation
   - JWT support

5. **Testing Bonus Features** (0.5 hour)

### Phase 10: Docker & Documentation (1.5 hours)

1. **Create Dockerfile**
   - Multi-stage build
   - Optimization

2. **Create Docker Compose**
   - Service definitions
   - Network configuration
   - Volume management

3. **Write README**
   - Setup instructions
   - Architecture description
   - API documentation
   - Bonus features explanation

4. **Create Postman Collection**
   - API endpoints
   - Example requests
   - Test scenarios

### Phase 11: Testing & Refinement (1 hour)

1. **Integration Testing**
   - Test end-to-end flow
   - Test with Docker

2. **Performance Testing**
   - Load testing
   - Cache effectiveness

3. **Bug Fixes & Optimization**

---

## Testing Strategy

### Unit Tests

**Files to Test:**
- `internal/validator/validator_test.go`
- `internal/service/salary_calculator_test.go`
- `internal/service/process_status_test.go`

**Test Cases:**
```go
// Salary Calculator Tests
- TestSeasonalAdjustments_December
- TestSeasonalAdjustments_Summer
- TestTaxDeduction_AboveThreshold
- TestTaxDeduction_BelowThreshold
- TestStatusDetermination
- TestCompleteCalculationFlow

// Process Status Tests
- TestGetEmployeeStatus_Success
- TestGetEmployeeStatus_InvalidNationalNumber
- TestGetEmployeeStatus_InactiveUser
- TestGetEmployeeStatus_InsufficientData
```

### Integration Tests

**Test Scenarios:**
1. Complete flow with database
2. Cache hit/miss scenarios
3. Retry mechanism on DB failure
4. Token authentication flow

### Manual Testing

**Postman Collection Test Cases:**

1. **Valid Request - GREEN Status**
   - National Number: NAT1004
   - Expected: Status GREEN, Average > 2000

2. **Valid Request - RED Status**
   - National Number: NAT1001
   - Expected: Status RED, Average < 2000

3. **Invalid National Number**
   - National Number: NAT9999
   - Expected: 404 error

4. **Inactive User**
   - National Number: NAT1003
   - Expected: 406 error

5. **Insufficient Data**
   - National Number: NAT1011
   - Expected: 422 error

6. **Missing Token** (Bonus)
   - No Authorization header
   - Expected: 401 error

7. **Invalid Token** (Bonus)
   - Invalid Authorization header
   - Expected: 401 error

---

## API Documentation

### Health Check Endpoint

**URL:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "database": "connected",
  "cache": "connected",
  "timestamp": "2025-10-26T14:30:00Z"
}
```

### Metrics Endpoint (Optional)

**URL:** `GET /metrics`

**Response:**
```json
{
  "total_requests": 1234,
  "cache_hits": 890,
  "cache_misses": 344,
  "average_response_time_ms": 45.2
}
```

---

## Configuration Management

### Configuration Structure (`internal/config/config.go`)

```go
type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Cache    CacheConfig
    Security SecurityConfig
    Logging  LoggingConfig
}

type AppConfig struct {
    Port string
    Env  string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
    MaxOpenConns int
    MaxIdleConns int
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type CacheConfig struct {
    Enabled bool
    TTL     int // seconds
}

type SecurityConfig struct {
    APISecretKey   string
    TokenRequired  bool
}

type LoggingConfig struct {
    Level  string
    ToDB   bool
}
```

---

## Error Handling Strategy

### Custom Error Types

```go
type AppError struct {
    Code    int    `json:"-"`
    Message string `json:"error"`
    Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Error Constructors
func NewNotFoundError(message string) *AppError
func NewValidationError(message string) *AppError
func NewUnauthorizedError(message string) *AppError
func NewInternalError(message string) *AppError
```

### Error Response Handler

```go
func handleError(c *gin.Context, err error) {
    if appErr, ok := err.(*AppError); ok {
        c.JSON(appErr.Code, appErr)
        return
    }
    
    c.JSON(500, gin.H{
        "error": "Internal server error",
    })
}
```

---

## Performance Optimization

### Database Query Optimization

1. **Indexes:**
   - Index on `national_number` for fast lookups
   - Composite index on `user_id` + `year` + `month` for salary queries

2. **Query Strategy:**
   - Use prepared statements
   - Fetch only required columns
   - Use joins instead of N+1 queries

3. **Connection Pooling:**
   - Configure max open connections
   - Configure max idle connections
   - Set connection max lifetime

### Caching Strategy

1. **What to Cache:**
   - Complete employee status response
   - User information (if frequently accessed)

2. **Cache Invalidation:**
   - Time-based expiration (5 minutes)
   - Manual invalidation on data updates

3. **Cache Key Format:**
   ```
   emp_status:{nationalNumber}
   ```

### Response Time Goals

- **Without Cache:** < 100ms
- **With Cache Hit:** < 10ms
- **P95 Response Time:** < 150ms

---

## Security Considerations

### Input Validation
- Validate all user inputs
- Sanitize national number input
- Prevent SQL injection (use parameterized queries)

### Authentication
- Implement token-based authentication
- Use HTTPS in production
- Rotate tokens regularly

### Data Protection
- Hash sensitive data
- Don't log sensitive information
- Use environment variables for secrets

### Rate Limiting (Optional)
```go
// Add rate limiting middleware
import "github.com/gin-contrib/rate"

router.Use(rate.NewRateLimiter(
    100, // max requests
    time.Minute, // per minute
))
```

---

## Deployment Checklist

### Pre-deployment
- [ ] All tests passing
- [ ] Environment variables configured
- [ ] Database migrations ready
- [ ] Docker images built successfully
- [ ] Documentation complete

### Production Configuration
- [ ] Set `GIN_MODE=release`
- [ ] Use strong API secret keys
- [ ] Configure proper database credentials
- [ ] Enable HTTPS
- [ ] Set up monitoring and alerting
- [ ] Configure log rotation
- [ ] Set up backup strategy

### Post-deployment
- [ ] Health check endpoint responding
- [ ] Database connected successfully
- [ ] Cache working properly
- [ ] Logs being written correctly
- [ ] API responding to requests
- [ ] Monitor error rates
- [ ] Monitor response times

---

## Monitoring & Observability

### Logging Strategy

**What to Log:**
- All API requests (endpoint, status, duration)
- Database operations
- Cache hits/misses
- Errors and exceptions
- Business logic decisions

**Log Levels:**
- DEBUG: Detailed diagnostic information
- INFO: General informational messages
- WARN: Warning messages
- ERROR: Error events

### Metrics to Track

1. **Request Metrics:**
   - Total requests
   - Requests per endpoint
   - Response times
   - Error rates

2. **Database Metrics:**
   - Query execution time
   - Connection pool usage
   - Retry attempts

3. **Cache Metrics:**
   - Hit rate
   - Miss rate
   - Cache size

4. **Business Metrics:**
   - Status distribution (GREEN/ORANGE/RED)
   - Average salaries processed
   - Error type distribution

---

## Troubleshooting Guide

### Common Issues

**Issue: Database Connection Failed**
```
Solution:
1. Check database credentials in .env
2. Verify database is running: docker ps
3. Check network connectivity
4. Review database logs
```

**Issue: High Response Times**
```
Solution:
1. Check database query performance
2. Verify cache is working
3. Review connection pool settings
4. Check for N+1 query problems
```

**Issue: Cache Not Working**
```
Solution:
1. Verify Redis is running (if using Redis)
2. Check cache configuration in .env
3. Review cache TTL settings
4. Check cache key generation
```

**Issue: Token Authentication Failing**
```
Solution:
1. Verify token format (Bearer <token>)
2. Check API_SECRET_KEY configuration
3. Verify token hasn't expired
4. Review middleware implementation
```

---

## Future Enhancements

### Potential Features

1. **Advanced Caching:**
   - Distributed caching with Redis cluster
   - Cache warming strategies
   - Predictive caching

2. **Analytics Dashboard:**
   - Real-time statistics
   - Historical trends
   - Salary distribution visualizations

3. **Batch Processing:**
   - Bulk employee status retrieval
   - Scheduled report generation

4. **Advanced Security:**
   - OAuth2 integration
   - Role-based access control
   - API key rotation

5. **Performance:**
   - GraphQL endpoint
   - gRPC support
   - WebSocket for real-time updates

6. **Data Export:**
   - Export to PDF/Excel
   - Email reports
   - Scheduled exports

---

## References & Resources

### Documentation Links
- Gin Framework: https://gin-gonic.com/docs/
- SQLX: http://jmoiron.github.io/sqlx/
- Ozzo Validation: https://github.com/go-ozzo/ozzo-validation
- Viper Config: https://github.com/spf13/viper
- Docker Compose: https://docs.docker.com/compose/

### Best Practices
- Go Project Layout: https://github.com/golang-standards/project-layout
- Effective Go: https://go.dev/doc/effective_go
- Go Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments

---

## File Locations Summary

| File Type | Location | Description |
|-----------|----------|-------------|
| Database Schema | `migrations/001_create_tables.sql` | Tables, indexes, constraints |
| Sample Data | `migrations/002_insert_data.sql` | Initial data inserts |
| Main Entry | `cmd/api/main.go` | Application entry point |
| Configuration | `.env` | Environment variables |
| Docker | `docker-compose.yml` | Container orchestration |
| Dockerfile | `docker/Dockerfile` | Application container |
| API Docs | `docs/api/postman_collection.json` | Postman collection |

---

## Quick Start Commands

```bash
# Setup
git clone <your-repo>
cd getemps-service
cp .env.example .env
# Edit .env with your configuration

# Run with Docker Compose
docker-compose up -d

# Check service health
curl http://localhost:8080/health

# Test API
curl -X POST http://localhost:8080/api/GetEmpStatus \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token-here" \
  -d '{"NationalNumber": "NAT1001"}'

# View logs
docker-compose logs -f api

# Stop services
docker-compose down

# Development (without Docker)
go mod download
go run cmd/api/main.go
```

---

## Contact & Support

For questions or issues:
- Review the README.md
- Check the troubleshooting guide
- Review API documentation
- Check logs for error messages

---

**Project Status:** Ready for Implementation
**Estimated Completion:** 10-16 hours
**Last Updated:** October 26, 2025
