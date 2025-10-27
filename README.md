# GetEmpStatus Web Service

A RESTful web service that retrieves employee information, processes salary data, and calculates employee status based on defined business logic.

## Status Badges

[![CI Pipeline](https://github.com/rixtrayker/getemps-service/workflows/CI%20Pipeline/badge.svg)](https://github.com/rixtrayker/getemps-service/actions)
[![CD Pipeline](https://github.com/rixtrayker/getemps-service/workflows/CD%20Pipeline/badge.svg)](https://github.com/rixtrayker/getemps-service/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/rixtrayker/getemps-service)](https://goreportcard.com/report/github.com/rixtrayker/getemps-service)
[![codecov](https://codecov.io/gh/rixtrayker/getemps-service/branch/main/graph/badge.svg)](https://codecov.io/gh/rixtrayker/getemps-service)
[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

> **Badge Status**: CI/CD pipelines run tests, builds, and deployments. Go Report Card analyzes code quality. CodeCov tracks test coverage. All badges link to their respective detailed reports.

## Features

### Core Features
- Employee status retrieval by national number
- Salary adjustments (seasonal bonuses/deductions)
- Tax calculations
- Status determination (GREEN/ORANGE/RED)
- Comprehensive input validation
- Error handling for edge cases

### Bonus Features
- ✅ In-memory caching (5-minute TTL)
- ✅ Custom structured logging
- ✅ Database retry mechanism
- ✅ Token-based authentication (optional)
- ✅ Health check endpoints
- ✅ Graceful shutdown
- ✅ Docker containerization

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for development)
- Make (for using Makefiles)

### Super Quick Start (Recommended)

```bash
# One-command setup and start
make quick-start
```

This will:
- Set up the development environment
- Start all services with Docker Compose
- Show you what's available

### Manual Setup

1. **Clone and setup**
   ```bash
   git clone https://github.com/rixtrayker/getemps-service.git
   cd getemps-service
   make setup  # or: cp .env.example .env && go mod download
   ```

2. **Start services**
   ```bash
   make up  # or: docker-compose up -d
   ```

3. **Test the API**
   ```bash
   make test-api  # or: scripts/test-api.sh
   ```

### Using Makefiles

This project includes comprehensive Makefiles for all operations:

```bash
# Show all available commands
make help        # Main development commands
make help-all    # All commands from all Makefiles

# Development workflow
make dev         # Start with live reload
make test        # Run tests
make build       # Build binary
make clean       # Clean artifacts

# Docker operations
make up          # Start all services
make down        # Stop all services
make logs        # View logs
make ps          # Show running containers

# Testing
make test-unit          # Unit tests only
make test-integration   # Integration tests
make test-coverage      # Test with coverage
make test-load-light    # Light load testing
```

### Traditional Docker Compose (Alternative)

```bash
# Start services
docker-compose up -d

# Check service health
curl http://localhost:8080/health

# Test API
curl -X POST http://localhost:8080/api/GetEmpStatus \
  -H "Content-Type: application/json" \
  -d '{"NationalNumber": "NAT1001"}'
```

### Development Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Setup PostgreSQL**
   ```bash
   docker run --name postgres -e POSTGRES_PASSWORD=postgres123 -e POSTGRES_DB=getemps_db -p 5432:5432 -d postgres:15-alpine
   ```

3. **Run migrations**
   ```bash
   psql -h localhost -U postgres -d getemps_db -f migrations/001_create_tables.sql
   psql -h localhost -U postgres -d getemps_db -f migrations/002_insert_data.sql
   ```

4. **Run application**
   ```bash
   go run cmd/api/main.go
   ```

## API Testing with Postman

### Postman Collections

The project includes comprehensive Postman collections for API testing:

| Collection | Purpose | Test Cases |
|-----------|---------|------------|
| **GetEmpStatus_Complete_Collection.json** | Full API testing suite | 25+ test scenarios |
| **GetEmpStatus_Performance_Collection.json** | Performance testing | Load & response time tests |
| **postman_collection.json** | Basic API testing | 9 essential tests |

### Environment Files

| Environment | File | Usage |
|-------------|------|-------|
| Development | `environments/Development.postman_environment.json` | Local development |
| Docker | `environments/Docker.postman_environment.json` | Container testing |
| Staging | `environments/Staging.postman_environment.json` | Staging environment |
| Production | `environments/Production.postman_environment.json` | Production testing |

### Using Postman Collections

**Import Collections:**
1. Open Postman
2. Import `docs/api/GetEmpStatus_Complete_Collection.json`
3. Import environment file (e.g., `docs/api/environments/Development.postman_environment.json`)
4. Select environment and run tests

**Command Line Testing:**
```bash
# Install Newman (Postman CLI)
make install-tools  # Includes Newman installation

# Run complete test suite
make test-postman

# Run performance tests
make test-postman-performance

# Run all API tests
make test-api-all
```

**Manual Newman Usage:**
```bash
# Install Newman
npm install -g newman

# Run complete collection
newman run docs/api/GetEmpStatus_Complete_Collection.json \
  -e docs/api/environments/Development.postman_environment.json

# Run with iterations for load testing
newman run docs/api/GetEmpStatus_Performance_Collection.json \
  -e docs/api/environments/Development.postman_environment.json \
  --iteration-count 10 --delay-request 100
```

### Test Coverage

The Postman collections cover:
- ✅ All success scenarios (RED/ORANGE/GREEN status)
- ❌ All error scenarios (404, 406, 422, 400)
- 🔍 Input validation and edge cases
- 🔐 Authentication testing
- 🧮 Business logic validation
- ⚡ Performance benchmarks
- 📊 Response structure validation

See `docs/api/POSTMAN_README.md` for detailed documentation.

## API Documentation

### Endpoint: Get Employee Status

**URL:** `POST /api/GetEmpStatus`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <token> (optional)
```

**Request Body:**
```json
{
  "NationalNumber": "NAT1001"
}
```

**Success Response (200):**
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

**Error Responses:**
- `404` - Invalid National Number
- `406` - User is not Active  
- `422` - INSUFFICIENT_DATA (less than 3 salary records)
- `401` - Unauthorized (if token authentication enabled)

### Health Check

**URL:** `GET /health`

**Response:**
```json
{
  "status": "healthy",
  "timestamp": {
    "now": "2025-10-26T14:30:00Z"
  }
}
```

## Business Logic

### Salary Calculation Process

1. **Seasonal Adjustments** (applied first):
   - December: +10% holiday bonus
   - Summer months (June, July, August): -5% deduction

2. **Tax Deduction** (applied second):
   - If total adjusted salary > 10,000: apply 7% tax to each salary

3. **Status Determination** (based on average):
   - Average > 2000: GREEN
   - Average = 2000: ORANGE  
   - Average < 2000: RED

### Example Calculation

For user NAT1001 with salaries:
- Jan 2025: 1200
- Feb 2025: 1300  
- Mar 2025: 1400
- May 2025: 1500
- Jun 2025: 1600 (summer month)

**Step 1 - Seasonal Adjustments:**
- Jan: 1200 (no change)
- Feb: 1300 (no change)
- Mar: 1400 (no change)
- May: 1500 (no change)
- Jun: 1600 × 0.95 = 1520 (summer deduction)

**Step 2 - Tax Check:**
- Total: 1200 + 1300 + 1400 + 1500 + 1520 = 6920
- Since 6920 < 10000, no tax applied

**Step 3 - Final Calculations:**
- Average: 6920 ÷ 5 = 1384
- Highest: 1520
- Sum: 6920
- Status: RED (average < 2000)

## Architecture

### Project Structure
```
getemps-service/
├── cmd/api/main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── models/                  # Data models
│   ├── validator/               # Input validation
│   ├── repository/              # Database access layer
│   ├── service/                 # Business logic
│   ├── handler/                 # HTTP handlers
│   ├── middleware/              # HTTP middleware
│   ├── cache/                   # Caching layer
│   ├── logger/                  # Custom logging
│   └── database/                # Database utilities
├── migrations/                  # Database migrations
├── docker/                      # Docker configuration
└── docker-compose.yml           # Container orchestration
```

### Technology Stack
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP)
- **Database:** PostgreSQL 15+
- **Cache:** In-memory (go-cache)
- **Validation:** ozzo-validation
- **Logging:** Logrus
- **Config:** Viper
- **Container:** Docker

## Configuration

All configuration is managed through environment variables. See `.env.example` for available options.

### Key Configuration Options

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

# Cache
CACHE_ENABLED=true
CACHE_TTL=300

# Security  
API_SECRET_KEY=your-secret-key
TOKEN_REQUIRED=false

# Logging
LOG_LEVEL=info
LOG_TO_DB=true
```

## Sample Data

The service comes with pre-loaded test data for the following scenarios:

| National Number | Scenario | Expected Status |
|----------------|----------|----------------|
| NAT1001 | Basic calculation with summer month | RED |
| NAT1002 | Average exactly 2000 | ORANGE |
| NAT1004 | High salary average | GREEN |
| NAT1005 | Tax deduction scenario | GREEN |
| NAT1006 | December bonus | GREEN |
| NAT1007 | Multiple summer months | GREEN |
| NAT1003 | Inactive user | 406 Error |
| NAT1011 | Insufficient data | 422 Error |
| NAT9999 | Non-existent user | 404 Error |

## Performance Features

### Caching
- In-memory cache with 5-minute TTL
- Cache key format: `emp_status:{nationalNumber}`
- Automatic cache invalidation

### Database Optimization
- Indexes on frequently queried columns
- Connection pooling
- Retry mechanism for failed queries

### Response Times
- Without cache: < 100ms
- With cache hit: < 10ms
- P95 response time: < 150ms

## Security

### Authentication
- Optional token-based authentication
- Configurable via `TOKEN_REQUIRED` environment variable
- Simple Bearer token format

### Input Validation
- National number format validation
- JSON schema validation
- SQL injection prevention (parameterized queries)

### Error Handling
- Structured error responses
- No sensitive data in error messages
- Comprehensive logging for debugging

## Monitoring & Observability

### Logging
- Structured JSON logging
- Request/response logging
- Database operation logging
- Optional database log storage

### Health Checks
- Application health endpoint
- Database connectivity check
- Service dependency monitoring

## Development

### Makefile Overview

The project includes multiple specialized Makefiles for different aspects:

| Makefile | Purpose | Key Commands |
|----------|---------|--------------|
| `Makefile` | Main development tasks | `make help`, `make dev`, `make test`, `make build` |
| `docker/Makefile.docker` | Docker operations | `make docker-build`, `make docker-run`, `make docker-push` |
| `scripts/Makefile.test` | Testing workflows | `make test-unit`, `make test-coverage`, `make test-load` |
| `Makefile.all` | Unified interface | `make help-all`, `make quick-start` |

### Development Workflow

```bash
# Initial setup
make setup                    # Setup development environment
make install-tools           # Install additional dev tools

# Development cycle
make dev                     # Start with live reload (requires air)
make run                     # Run without live reload
make test                    # Run tests
make fmt                     # Format code
make lint                    # Run linter

# Building and testing
make build                   # Build binary
make test-coverage          # Test with coverage report
make test-integration       # Run integration tests
make test-api              # Test API endpoints

# Docker development
make docker-build          # Build Docker image
make docker-run           # Run in container
make up                   # Start all services
make down                 # Stop all services
```

### Database Operations

```bash
# Database management
make db-up                 # Start database only
make db-down              # Stop database
make db-reset             # Reset database (WARNING: destroys data)
make db-migrate           # Run migrations
make db-shell             # Connect to database shell

# Traditional commands (alternative)
docker-compose down -v    # Reset database
docker-compose up -d      # Start services
psql -h localhost -U postgres -d getemps_db -f migrations/001_create_tables.sql
```

### Testing Commands

```bash
# Unit testing
make test-unit            # Fast unit tests
make test-service         # Test service layer
make test-repository      # Test repository layer

# API testing
make test-api             # Test API endpoints (shell script)
make test-postman         # Run Postman collection tests
make test-postman-performance # Run performance tests
make test-api-all         # Run all API tests

# Integration testing
make test-integration     # Full integration tests
make test-api-integration # API integration tests

# Coverage and reporting
make test-coverage        # Generate coverage report
make test-coverage-html   # HTML coverage report
make test-coverage-check  # Verify coverage threshold

# Performance testing
make test-benchmark       # Run benchmarks
make test-load-light      # Light load test (100 requests)
make test-load-medium     # Medium load test (1000 requests)
make test-stress          # Stress test with multiple scenarios

# Security testing
make test-security        # Security vulnerability scan
make test-auth           # Authentication tests
```

### Docker Commands

```bash
# Image management
make docker-build-prod    # Build production image
make docker-build-dev     # Build development image
make docker-push         # Push to registry
make docker-clean        # Clean images and containers

# Container operations
make docker-run          # Run container
make docker-shell        # Shell into container
make docker-logs         # View container logs

# Advanced operations
make docker-scan-security # Security scan
make docker-benchmark    # Benchmark startup time
make docker-dive         # Analyze layers (requires dive tool)
```

## Troubleshooting

### Common Issues

**Database Connection Failed**
- Check PostgreSQL is running: `docker ps`
- Verify credentials in `.env`
- Check network connectivity

**High Response Times**
- Check cache is enabled
- Verify database indexes
- Review connection pool settings

**Authentication Errors**
- Verify `TOKEN_REQUIRED` setting
- Check Bearer token format
- Validate API secret key

### Logs Location
- Application logs: stdout (JSON format)
- Database logs: `logs` table (if enabled)
- Docker logs: `docker-compose logs -f api`

## Contributing

1. Fork the repository
2. Create feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit pull request

## License

This project is licensed under the MIT License.