#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}🧪 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Function to run tests with clean output
run_tests() {
    local test_path="$1"
    local test_name="$2"
    
    print_status "Running $test_name tests..."
    
    # Run tests and capture output
    output=$(docker-compose run --rm api sh -c "go test -cover $test_path" 2>&1)
    exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        print_success "$test_name tests passed!"
        # Extract and display coverage
        coverage=$(echo "$output" | grep "coverage:" | tail -1)
        if [ ! -z "$coverage" ]; then
            echo -e "${GREEN}📊 $coverage${NC}"
        fi
    else
        print_error "$test_name tests failed!"
        echo "$output"
        return 1
    fi
}

# Main execution
print_status "Starting test suite..."

# Run database tests
run_tests "./db/sqlc/" "Database"
if [ $? -ne 0 ]; then
    print_error "Database tests failed!"
    exit 1
fi

# Run controller tests
run_tests "./internal/apps/accounts/controllers/" "Controller"
if [ $? -ne 0 ]; then
    print_error "Controller tests failed!"
    exit 1
fi

print_success "All tests completed successfully!" 