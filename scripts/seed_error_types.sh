#!/bin/bash

# Script to seed and delete error types from the Delícias da Lu API
# Usage: ./seed_error_types.sh --action create --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app --username admin --password password

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
API_URL="https://delicias-da-lu-514609008596.southamerica-east1.run.app"
ACTION=""
USERNAME="josealdo"
PASSWORD="senhaforte"
WAIT_TIME=2

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --api-url)
            API_URL="$2"
            shift 2
            ;;
        --action)
            ACTION="$2"
            shift 2
            ;;
        --username)
            USERNAME="$2"
            shift 2
            ;;
        --password)
            PASSWORD="$2"
            shift 2
            ;;
        --wait)
            WAIT_TIME="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Validate inputs
if [ -z "$ACTION" ] || [ -z "$USERNAME" ] || [ -z "$PASSWORD" ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo "Usage: $0 --action [create|delete|delete-all|create-and-delete] --username USERNAME --password PASSWORD [--api-url URL] [--wait SECONDS]"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Delícias da Lu - Error Types Management${NC}"
echo -e "${BLUE}========================================${NC}"

# Authenticate
echo -e "${BLUE}ℹ️  Authenticating as $USERNAME...${NC}"
AUTH_RESPONSE=$(curl -s -X POST "$API_URL/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo "$AUTH_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Authentication failed${NC}"
    echo "Response: $AUTH_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Authentication successful${NC}"
echo -e "${BLUE}Token: ${TOKEN:0:20}...${NC}"

# Helper function for create operation
create_error_types() {
    echo -e "${BLUE}ℹ️  Creating error types...${NC}"
    
    PAYLOAD='[
        {"identifier":"not-found","title":"Not Found","detail":"The requested resource could not be found.","status":404},
        {"identifier":"invalid-type","title":"Invalid Type","detail":"The provided type is not supported.","status":400},
        {"identifier":"invalid-credentials","title":"Invalid Credentials","detail":"The username or password is incorrect.","status":401},
        {"identifier":"invalid-token","title":"Invalid Token","detail":"The provided token is invalid or expired.","status":401},
        {"identifier":"unauthorized","title":"Unauthorized","detail":"You are not authorized to access this resource.","status":401},
        {"identifier":"type-not-found","title":"Error Type Not Found","detail":"No error type found for the provided identifier.","status":404},
        {"identifier":"instance-not-found","title":"Error Instance Not Found","detail":"No error instance found for the provided identifier.","status":404},
        {"identifier":"invalidDocumentLenght","title":"Invalid Document Lenght","detail":"The provided document does not have a known length (CPF or CNPJ).","status":400},
        {"identifier":"invalidFilter","title":"Invalid Filter","detail":"The filter query parameter is invalid. Valid values are '\''type'\'' and '\''instance'\''.","status":400},
        {"identifier":"unexpectedUnhandledError","title":"Unexpected Error","detail":"An unexpected error occurred. Please contact support.","status":500}
    ]'
    
    RESPONSE=$(curl -s -X POST "$API_URL/v1/admin/error-types" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$PAYLOAD")
    
    if echo "$RESPONSE" | grep -q '"message"'; then
        echo -e "${GREEN}✓ Error types created successfully${NC}"
        echo "$RESPONSE" | grep -o '"message":"[^"]*' || true
    else
        echo -e "${RED}✗ Failed to create error types${NC}"
        echo "Response: $RESPONSE"
        return 1
    fi
}

# Helper function for delete operation
delete_error_type() {
    local IDENTIFIER=$1
    local RESPONSE=$(curl -s -X DELETE "$API_URL/v1/admin/error-types/$IDENTIFIER" \
        -H "Authorization: Bearer $TOKEN")
    
    if echo "$RESPONSE" | grep -q '"message"'; then
        echo -e "${GREEN}✓ Deleted: $IDENTIFIER${NC}"
        return 0
    else
        echo -e "${RED}✗ Failed to delete: $IDENTIFIER${NC}"
        return 1
    fi
}

# Helper function for delete all
delete_all_error_types() {
    echo -e "${BLUE}ℹ️  Deleting all error types...${NC}"
    
    local IDENTIFIERS=(
        "not-found"
        "invalid-type"
        "invalid-credentials"
        "invalid-token"
        "unauthorized"
        "type-not-found"
        "instance-not-found"
        "invalidDocumentLenght"
        "invalidFilter"
        "unexpectedUnhandledError"
    )
    
    local FAILED=0
    for identifier in "${IDENTIFIERS[@]}"; do
        delete_error_type "$identifier" || ((FAILED++))
    done
    
    if [ $FAILED -eq 0 ]; then
        echo -e "${GREEN}✓ All error types deleted successfully${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Failed to delete $FAILED error types${NC}"
        return 1
    fi
}

# Perform action
case $ACTION in
    create)
        create_error_types
        ;;
    delete)
        echo -e "${RED}Error: --action delete requires error identifier${NC}"
        exit 1
        ;;
    delete-all)
        delete_all_error_types
        ;;
    create-and-delete)
        create_error_types
        echo -e "${BLUE}ℹ️  Waiting $WAIT_TIME seconds before deletion...${NC}"
        sleep "$WAIT_TIME"
        delete_all_error_types
        ;;
    *)
        echo -e "${RED}Unknown action: $ACTION${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Operation completed${NC}"
echo -e "${BLUE}========================================${NC}"
