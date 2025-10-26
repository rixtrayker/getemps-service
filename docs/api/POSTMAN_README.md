# GetEmpStatus API - Postman Collections

This directory contains comprehensive Postman collections for testing the GetEmpStatus API across different scenarios and environments.

## 📁 Collection Overview

### Main Collections

| Collection | Purpose | Test Cases | Environment |
|-----------|---------|------------|-------------|
| **GetEmpStatus_Complete_Collection.json** | Comprehensive API testing | 25+ scenarios | All environments |
| **GetEmpStatus_Performance_Collection.json** | Performance & load testing | 15+ performance tests | Development/Staging |
| **postman_collection.json** | Basic API testing | 9 basic scenarios | Development |

### Environment Files

| Environment | File | Usage |
|-------------|------|-------|
| **Development** | `environments/Development.postman_environment.json` | Local development testing |
| **Docker** | `environments/Docker.postman_environment.json` | Docker container testing |
| **Staging** | `environments/Staging.postman_environment.json` | Staging environment testing |
| **Production** | `environments/Production.postman_environment.json` | Production testing (with auth) |

## 🚀 Quick Start

### 1. Import Collections

1. Open Postman
2. Click "Import" button
3. Drag and drop or select the collection files:
   - `GetEmpStatus_Complete_Collection.json` (recommended)
   - `GetEmpStatus_Performance_Collection.json` (for performance testing)

### 2. Import Environment

1. In Postman, go to "Environments"
2. Click "Import"
3. Select the appropriate environment file:
   - `Development.postman_environment.json` for local testing
   - `Docker.postman_environment.json` for Docker testing

### 3. Set Active Environment

1. In the top-right corner of Postman, select your imported environment
2. Verify the `baseUrl` variable matches your setup:
   - Development: `http://localhost:8080`
   - Docker: `http://localhost:8080`

### 4. Run Tests

- **Individual Test**: Click on any request and hit "Send"
- **Collection Runner**: Use Postman's Collection Runner for automated testing
- **Command Line**: Use Newman for CI/CD integration

## 📋 Test Scenarios

### Complete Collection Test Coverage

#### ✅ Successful Scenarios
- **RED Status** (NAT1001): Average salary < 2000
- **ORANGE Status** (NAT1002): Average salary = 2000  
- **GREEN Status** (NAT1004): Average salary > 2000
- **Tax Deduction** (NAT1005): High salary with tax calculation
- **December Bonus** (NAT1006): Holiday bonus scenario
- **Summer Deduction** (NAT1007): Summer month adjustments

#### ❌ Error Scenarios
- **404**: Invalid national number (NAT9999)
- **406**: Inactive user (NAT1003)
- **422**: Insufficient data (NAT1011)

#### 🔍 Validation Tests
- **400**: Invalid request format
- **400**: Empty request body
- **400**: Malformed JSON
- **400**: Empty national number

#### 🔐 Authentication Tests
- Valid token requests
- Missing token handling
- Invalid token handling

#### 🧮 Business Logic Validation
- Salary calculation verification
- Data type validation
- Status determination logic
- Response structure validation

#### ⚡ Performance Tests
- Response time benchmarks
- Cache performance testing
- Different data scenario performance

### Performance Collection Test Coverage

#### 🚀 Performance Baseline
- Single request benchmarks
- Cache performance validation
- Response time categorization (excellent < 100ms, good < 500ms, acceptable < 1000ms)

#### 📊 Different Data Scenarios
- Simple calculations
- Complex calculations (tax scenarios)
- Seasonal adjustments

#### 🔄 Concurrent Request Simulation
- Multiple user simulation
- Performance statistics calculation
- Load distribution testing

#### 🎯 Error Response Performance
- Fast error responses
- Validation error performance

#### 📈 Progressive Load Testing
- Incremental load testing
- Throughput calculations
- Performance degradation monitoring

## 🛠️ Advanced Usage

### Running with Newman (CLI)

Install Newman:
```bash
npm install -g newman
```

Run complete collection:
```bash
newman run GetEmpStatus_Complete_Collection.json \
  -e environments/Development.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export results.json
```

Run performance tests:
```bash
newman run GetEmpStatus_Performance_Collection.json \
  -e environments/Development.postman_environment.json \
  --iteration-count 10 \
  --delay-request 100
```

### Environment Variables

Each environment includes these variables:

| Variable | Development | Docker | Staging | Production |
|----------|-------------|--------|---------|------------|
| `baseUrl` | http://localhost:8080 | http://localhost:8080 | https://api-staging.example.com | https://api.example.com |
| `authToken` | demo-token-12345678 | docker-test-token-87654321 | Dynamic | Secret |
| `timeoutMs` | 5000 | 10000 | 15000 | 30000 |

### Custom Variables

You can customize these variables for your testing:

```json
{
  "baseUrl": "http://your-custom-url:8080",
  "authToken": "your-test-token",
  "timeoutMs": "5000"
}
```

## 🔧 Test Configuration

### Automated Test Scripts

Each request includes automated test scripts that verify:

- **HTTP Status Codes**: Correct response codes for each scenario
- **Response Structure**: Required fields and data types
- **Business Logic**: Salary calculations and status determination
- **Performance**: Response time thresholds
- **Security**: Authentication and authorization

### Example Test Script

```javascript
pm.test('Response status is 200', function () {
    pm.response.to.have.status(200);
});

pm.test('Response has all required fields', function () {
    var response = pm.response.json();
    pm.expect(response).to.have.property('id');
    pm.expect(response).to.have.property('nationalNumber');
    pm.expect(response).to.have.property('status');
});

pm.test('Status determination is correct', function () {
    var response = pm.response.json();
    var avg = response.salaryDetails.averageSalary;
    var status = response.status;
    
    if (avg > 2000) {
        pm.expect(status).to.eql('GREEN');
    } else if (avg === 2000) {
        pm.expect(status).to.eql('ORANGE');
    } else {
        pm.expect(status).to.eql('RED');
    }
});
```

## 📊 Performance Benchmarks

### Expected Performance Metrics

| Scenario | Expected Response Time | Threshold |
|----------|----------------------|-----------|
| Cache Hit | < 50ms | Excellent |
| Simple Calculation | < 100ms | Good |
| Complex Calculation | < 500ms | Acceptable |
| Error Responses | < 200ms | Fast |
| Database Query | < 1000ms | Acceptable |

### Performance Test Results

The performance collection automatically calculates:
- Average response time
- Min/Max response times
- Throughput (requests/second)
- Performance degradation under load

## 🔍 Troubleshooting

### Common Issues

**Collection Import Fails**
- Ensure you're using Postman v8.0+
- Check file format is valid JSON
- Try importing one file at a time

**Environment Variables Not Working**
- Verify environment is selected in top-right dropdown
- Check variable names match exactly (case-sensitive)
- Ensure environment file imported correctly

**Tests Failing**
- Check service is running: `make up` or `docker-compose up -d`
- Verify baseUrl in environment matches your setup
- Check authentication token if auth is enabled

**Performance Tests Inconsistent**
- Run tests multiple times for average
- Ensure system is not under other load
- Check database connection is stable

### Service Setup

Before running tests, ensure the service is running:

```bash
# Using Make (recommended)
make up
make test-api

# Using Docker Compose
docker-compose up -d
curl http://localhost:8080/health

# Local development
go run cmd/api/main.go
```

## 🔗 Integration with CI/CD

### GitHub Actions Example

```yaml
name: API Tests
on: [push, pull_request]

jobs:
  api-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Start services
        run: docker-compose up -d
        
      - name: Wait for services
        run: sleep 30
        
      - name: Install Newman
        run: npm install -g newman
        
      - name: Run API tests
        run: |
          newman run docs/api/GetEmpStatus_Complete_Collection.json \
            -e docs/api/environments/Docker.postman_environment.json \
            --reporters cli,junit \
            --reporter-junit-export test-results.xml
            
      - name: Publish test results
        uses: dorny/test-reporter@v1
        if: always()
        with:
          name: API Tests
          path: test-results.xml
          reporter: java-junit
```

### Make Integration

The project includes Make targets for Postman testing:

```bash
# Run API tests using the shell script
make test-api

# Run performance tests (requires hey tool)
make test-load

# Check service health
make health-check
```

## 📝 Best Practices

### Test Organization
- Use folders to organize related tests
- Include descriptive test names
- Add pre-request and test scripts for automation
- Use environment variables for configuration

### Performance Testing
- Run performance tests on consistent environments
- Test with realistic data loads
- Monitor both response time and throughput
- Include error scenario performance

### Security Testing
- Test authentication scenarios
- Validate authorization levels
- Check for sensitive data exposure
- Test input validation thoroughly

### Maintenance
- Update collections when API changes
- Review test thresholds regularly
- Keep environment files current
- Document any custom test scenarios

## 📞 Support

For issues with the Postman collections:

1. Check this documentation
2. Verify service setup with `make health-check`
3. Review Postman console for detailed error messages
4. Check the main project README for service configuration

The collections are designed to be comprehensive and self-documenting with automated test scripts and clear naming conventions.