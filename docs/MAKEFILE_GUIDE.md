# Makefile Guide for GetEmpStatus Service

This project uses multiple Makefiles to streamline development, testing, and deployment workflows.

## Overview of Makefiles

### 1. Main Makefile (`Makefile`)
**Purpose:** Core development tasks and project management
**Location:** Root directory

Key features:
- Environment setup and dependency management
- Building and running the application
- Basic testing and code quality checks
- Docker Compose orchestration
- Development tools integration

**Most Used Commands:**
```bash
make help           # Show available commands
make setup          # Setup development environment
make dev            # Start with live reload
make build          # Build the application
make test           # Run tests
make up             # Start all services
make down           # Stop all services
```

### 2. Docker Makefile (`docker/Makefile.docker`)
**Purpose:** Docker-specific operations and container management

Key features:
- Multi-stage Docker builds
- Image tagging and registry operations
- Container lifecycle management
- Security scanning and optimization
- Multi-architecture builds

**Most Used Commands:**
```bash
make docker-build     # Build production image
make docker-run       # Run in container
make docker-push      # Push to registry
make docker-clean     # Clean images and containers
```

### 3. Testing Makefile (`scripts/Makefile.test`)
**Purpose:** Comprehensive testing workflows

Key features:
- Unit, integration, and API testing
- Coverage reporting and validation
- Performance and load testing
- Security testing
- Test database management

**Most Used Commands:**
```bash
make test-unit        # Unit tests
make test-coverage    # Coverage report
make test-integration # Integration tests
make test-load-light  # Load testing
```

### 4. Master Makefile (`Makefile.all`)
**Purpose:** Unified interface for all Makefiles

Key features:
- Includes all other Makefiles
- Provides unified help system
- Defines convenient aliases
- Enhanced workflows combining multiple operations

**Most Used Commands:**
```bash
make help-all         # Show all commands
make quick-start      # Complete setup and start
make full-test        # Run all tests
make full-reset       # Complete environment reset
```

## Quick Command Reference

### Essential Commands for Daily Development

```bash
# Getting Started
make setup            # One-time setup
make quick-start      # Complete setup and start

# Development Cycle
make dev              # Start with live reload
make test             # Run tests
make fmt              # Format code
make build            # Build binary

# Service Management
make up               # Start all services
make down             # Stop all services
make restart          # Restart services
make logs             # View logs

# Testing
make test-unit        # Quick unit tests
make test-api         # Test API endpoints
make test-coverage    # Coverage report

# Docker Operations
make docker-build     # Build image
make docker-run       # Run container

# Database
make db-reset         # Reset database
make db-shell         # Database shell
```

### Advanced Commands

```bash
# Code Quality
make lint             # Run linter
make vet              # Run go vet
make security-scan    # Security analysis

# Testing (Advanced)
make test-benchmark   # Performance benchmarks
make test-load-heavy  # Heavy load testing
make test-integration # Integration tests
make test-stress      # Stress testing

# Docker (Advanced)
make docker-multi-arch    # Multi-architecture build
make docker-scan-security # Security scan
make docker-dive         # Image analysis

# CI/CD
make ci               # CI pipeline
make release-check    # Release validation
make prod-deploy      # Production deployment
```

## Workflow Examples

### New Developer Setup
```bash
# Clone repository
git clone <repository-url>
cd getemps-service

# One-command setup
make quick-start

# Alternative manual setup
make setup
make up
make test-api
```

### Daily Development Workflow
```bash
# Start development server
make dev

# In another terminal - run tests
make test

# Before committing
make pre-commit
```

### Testing Workflow
```bash
# Quick tests during development
make test-unit

# Full testing before merge
make test-all

# Performance testing
make test-load-light

# Security testing
make test-security
```

### Docker Development
```bash
# Build and test image
make docker-build
make docker-run

# Push to registry
make docker-push

# Clean up
make docker-clean
```

### Database Management
```bash
# Reset database (development)
make db-reset

# Manual migration
make db-migrate

# Access database
make db-shell
```

## Environment Variables

Some Makefile commands can be customized with environment variables:

```bash
# Docker registry
export DOCKER_REGISTRY=your-registry.com
make docker-push

# Custom image tag
export IMAGE_TAG=v1.2.3
make docker-build

# Test configuration
export TEST_TIMEOUT=60s
make test-all

# Base URL for API tests
export BASE_URL=http://localhost:9090
make test-api
```

## Integration with Development Tools

### With Air (Live Reload)
```bash
# Install air
make install-tools

# Start with live reload
make dev
```

### With IDE Integration
Many IDEs can integrate with Make:
- VSCode: Use "Tasks: Run Task" (Ctrl+Shift+P)
- GoLand: Run configurations can call make commands
- Vim/Neovim: `:make <target>`

### With Git Hooks
Add to `.git/hooks/pre-commit`:
```bash
#!/bin/bash
make pre-commit
```

## Troubleshooting

### Common Issues

**"make: command not found"**
- Install make: `brew install make` (macOS) or `apt-get install make` (Ubuntu)

**"No rule to make target"**
- Check available commands: `make help`
- Ensure you're in the correct directory

**Docker commands fail**
- Ensure Docker is running: `docker version`
- Check Docker Compose: `docker-compose version`

**Tests fail with database errors**
- Reset test database: `make db-reset`
- Check database is running: `make db-up`

**Permission denied on scripts**
- Make scripts executable: `chmod +x scripts/*.sh`

### Getting Help

```bash
# Show main commands
make help

# Show all commands from all Makefiles
make help-all

# Show Docker-specific commands
make docker-help

# Show testing commands
make test-help
```

## Customization

### Adding New Commands

To add commands to the main Makefile:

```makefile
.PHONY: my-command
my-command: ## Description of what this does
	@echo "Running my custom command..."
	@# Your commands here
```

### Creating Project-Specific Makefiles

For project-specific tasks, create `Makefile.custom`:

```makefile
include Makefile

.PHONY: deploy-staging
deploy-staging: ci docker-build ## Deploy to staging
	@echo "Deploying to staging..."
	@# Deployment commands
```

Then use: `make -f Makefile.custom deploy-staging`

## Best Practices

1. **Always use .PHONY** for targets that don't create files
2. **Add help text** with `## Description` for all public targets
3. **Use colors** for better output readability
4. **Check prerequisites** before running commands
5. **Provide aliases** for commonly used commands
6. **Include cleanup** commands for development artifacts
7. **Use environment variables** for configuration
8. **Test Makefile commands** in CI/CD pipelines

## Contributing

When adding new Makefile targets:

1. Choose the appropriate Makefile based on functionality
2. Follow existing naming conventions
3. Add helpful descriptions
4. Test on different environments
5. Update this documentation
6. Consider adding aliases in `Makefile.all`

For questions or suggestions about the Makefile structure, please open an issue.