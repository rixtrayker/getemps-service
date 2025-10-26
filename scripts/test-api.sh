#!/bin/bash

# API Testing Script for GetEmpStatus Service
# ===========================================

set -e

# Configuration
BASE_URL="${BASE_URL:-http://localhost:8080}"
API_ENDPOINT="$BASE_URL/api/GetEmpStatus"
HEALTH_ENDPOINT="$BASE_URL/health"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test function
test_endpoint() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_status="$5"
    local headers="$6"

    log_info "Testing: $name"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            ${headers:+-H "$headers"} \
            -d "$data" \
            "$url")
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" \
            ${headers:+-H "$headers"} \
            "$url")
    fi
    
    body=$(echo "$response" | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    status=$(echo "$response" | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    if [ "$status" = "$expected_status" ]; then
        log_success "$name - Status: $status"
        if command -v jq >/dev/null 2>&1; then
            echo "$body" | jq . 2>/dev/null || echo "$body"
        else
            echo "$body"
        fi
    else
        log_error "$name - Expected: $expected_status, Got: $status"
        echo "$body"
        return 1
    fi
    echo ""
}

# Main test suite
main() {
    log_info "Starting API tests for GetEmpStatus Service"
    log_info "Base URL: $BASE_URL"
    echo ""

    # Test 1: Health Check
    test_endpoint \
        "Health Check" \
        "GET" \
        "$HEALTH_ENDPOINT" \
        "" \
        "200"

    # Test 2: Missing Authentication Token (401)
    test_endpoint \
        "Missing Authentication Token" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1001"}' \
        "401"

    # Test 3: Invalid Token Format (401)
    test_endpoint \
        "Invalid Token Format" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1001"}' \
        "401" \
        "Authorization: InvalidFormat"

    # Test 4: Valid Employee (RED status) - WITH AUTH
    test_endpoint \
        "Valid Employee - RED Status" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1001"}' \
        "200" \
        "Authorization: Bearer test-token-12345"

    # Test 5: Valid Employee (GREEN status) - WITH AUTH
    test_endpoint \
        "Valid Employee - GREEN Status" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1004"}' \
        "200" \
        "Authorization: Bearer test-token-12345"

    # Test 6: Valid Employee (ORANGE status) - WITH AUTH
    test_endpoint \
        "Valid Employee - ORANGE Status" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1002"}' \
        "200" \
        "Authorization: Bearer test-token-12345"

    # Test 7: Invalid National Number - WITH AUTH
    test_endpoint \
        "Invalid National Number" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT9999"}' \
        "404" \
        "Authorization: Bearer test-token-12345"

    # Test 8: Inactive User - WITH AUTH
    test_endpoint \
        "Inactive User" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1003"}' \
        "406" \
        "Authorization: Bearer test-token-12345"

    # Test 9: Insufficient Data - WITH AUTH
    test_endpoint \
        "Insufficient Data" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1011"}' \
        "422" \
        "Authorization: Bearer test-token-12345"

    # Test 10: Invalid Request Format - WITH AUTH
    test_endpoint \
        "Invalid Request Format" \
        "POST" \
        "$API_ENDPOINT" \
        '{"InvalidField": "NAT1001"}' \
        "400" \
        "Authorization: Bearer test-token-12345"

    # Test 11: Empty Request Body - WITH AUTH
    test_endpoint \
        "Empty Request Body" \
        "POST" \
        "$API_ENDPOINT" \
        '{}' \
        "400" \
        "Authorization: Bearer test-token-12345"

    # Test 12: Malformed JSON - WITH AUTH
    test_endpoint \
        "Malformed JSON" \
        "POST" \
        "$API_ENDPOINT" \
        '{"NationalNumber": "NAT1001"' \
        "400" \
        "Authorization: Bearer test-token-12345"

    log_success "All API tests completed!"
}

# Check if service is running
check_service() {
    log_info "Checking if service is running..."
    if curl -s "$HEALTH_ENDPOINT" > /dev/null; then
        log_success "Service is running"
    else
        log_error "Service is not running at $BASE_URL"
        log_info "Start the service with: make up"
        exit 1
    fi
}

# Performance test function
performance_test() {
    log_info "Running basic performance test..."
    
    if command -v hey >/dev/null 2>&1; then
        hey -n 100 -c 5 -m POST \
            -H "Content-Type: application/json" \
            -d '{"NationalNumber": "NAT1001"}' \
            "$API_ENDPOINT"
    else
        log_warning "hey tool not installed. Skipping performance test."
        log_info "Install with: go install github.com/rakyll/hey@latest"
    fi
}

# Usage function
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -p, --performance   Run performance tests"
    echo "  -c, --check         Only check if service is running"
    echo "  -u, --url URL       Set base URL (default: http://localhost:8080)"
    echo ""
    echo "Examples:"
    echo "  $0                          # Run all tests"
    echo "  $0 -p                       # Run tests with performance"
    echo "  $0 -u http://localhost:9090 # Test different URL"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -p|--performance)
            PERFORMANCE=true
            shift
            ;;
        -c|--check)
            check_service
            exit 0
            ;;
        -u|--url)
            BASE_URL="$2"
            API_ENDPOINT="$BASE_URL/api/GetEmpStatus"
            HEALTH_ENDPOINT="$BASE_URL/health"
            shift 2
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Run tests
check_service
main

if [ "$PERFORMANCE" = true ]; then
    performance_test
fi