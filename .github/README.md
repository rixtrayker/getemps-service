# GitHub Actions CI/CD Pipelines

This repository includes comprehensive GitHub Actions workflows for continuous integration and deployment.

## Workflows Overview

### 1. CI Pipeline (`ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**
- **Test**: Runs unit tests with PostgreSQL service, generates coverage reports
- **Lint**: Performs code quality checks using golangci-lint
- **Build**: Builds binaries for Linux, Windows, and macOS

**Features:**
- Go module caching for faster builds
- PostgreSQL service for integration tests
- Code coverage reporting with Codecov
- Multi-platform binary builds
- Artifact upload for build outputs

### 2. CD Pipeline (`cd.yml`)

**Triggers:**
- Push to `main` branch
- Git tags starting with `v*`
- Successful completion of CI pipeline

**Jobs:**
- **Docker**: Builds and pushes Docker images to GitHub Container Registry
- **Deploy Staging**: Deploys to staging environment (main branch)
- **Deploy Production**: Deploys to production environment (version tags)

**Features:**
- Multi-architecture Docker builds (amd64, arm64)
- Automatic image tagging based on branch/tag
- Environment-specific deployments
- Docker layer caching for faster builds

### 3. Security Pipeline (`security.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Daily scheduled runs at 2 AM UTC

**Jobs:**
- **Dependency Check**: Scans for vulnerable dependencies using govulncheck and Nancy
- **Code Security**: Static code analysis using Gosec
- **Docker Security**: Container image vulnerability scanning using Trivy
- **Secrets Scan**: Detects exposed secrets using TruffleHog
- **License Check**: Validates dependency licenses using go-licenses

**Features:**
- SARIF report uploads to GitHub Security tab
- Automated vulnerability detection
- License compliance checking
- Secrets detection in code and history

## Setup Instructions

### 1. Repository Secrets

Configure the following secrets in your repository settings:

- `GITHUB_TOKEN`: Automatically provided by GitHub (no setup needed)

### 2. Environment Configuration

Create the following environments in your repository settings:

- **staging**: For staging deployments
- **production**: For production deployments (with protection rules)

### 3. Container Registry

The workflows are configured to use GitHub Container Registry (ghcr.io). Ensure your repository has the necessary permissions:

1. Go to repository Settings → Actions → General
2. Set "Workflow permissions" to "Read and write permissions"
3. Check "Allow GitHub Actions to create and approve pull requests"

### 4. Branch Protection

Configure branch protection rules for `main`:

1. Require status checks to pass before merging
2. Require branches to be up to date before merging
3. Include the following status checks:
   - `Test`
   - `Lint`
   - `Build`
   - `Dependency Vulnerability Scan`
   - `Code Security Analysis`

## Customization

### Database Configuration

The CI pipeline includes a PostgreSQL service for testing. Update the environment variables in `ci.yml` if you need different database settings:

```yaml
env:
  DB_HOST: localhost
  DB_PORT: 5432
  DB_USER: testuser
  DB_PASSWORD: testpassword
  DB_NAME: testdb
  DB_SSLMODE: disable
```

### Deployment Commands

Update the deployment steps in `cd.yml` with your specific deployment commands:

```yaml
- name: Deploy to staging
  run: |
    # Add your staging deployment commands here
    # Examples:
    # kubectl apply -f k8s/staging/
    # docker-compose -f docker-compose.staging.yml up -d
    # aws ecs update-service --cluster staging --service app
```

### Linting Configuration

Customize the linting rules by modifying `.golangci.yml`:

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    # Add or remove linters as needed
```

## Monitoring and Notifications

### Notifications

Configure notifications in repository settings:

1. Go to Settings → Notifications
2. Set up email or Slack notifications for workflow failures
3. Configure security alerts for vulnerability discoveries

## Troubleshooting

### Common Issues

1. **Test failures**: Check database connectivity and environment variables
2. **Lint failures**: Run `golangci-lint run` locally to identify issues
3. **Build failures**: Ensure all dependencies are properly declared in `go.mod`
4. **Security scan failures**: Review and address reported vulnerabilities

### Local Testing

Run the same checks locally:

```bash
# Run tests
go test -v -race ./...

# Run linting
golangci-lint run

# Run security scan
gosec ./...

# Build application
go build -o bin/main ./cmd/api
```

## Performance Optimization

The workflows include several optimizations:

- **Caching**: Go modules and Docker layers are cached
- **Parallel execution**: Jobs run in parallel where possible
- **Conditional execution**: Some jobs only run when necessary
- **Artifact management**: Build outputs are stored with appropriate retention

## Security Best Practices

- Secrets are never logged or exposed
- Minimal permissions are used for each job
- Security scans run on every change
- Dependencies are regularly updated and scanned
- Container images are scanned for vulnerabilities