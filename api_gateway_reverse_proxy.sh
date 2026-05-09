#!/bin/bash
# A script to initialize the project environment and run the application according to the developer requirements.

# Color codes
RED=$'\033[0;31m'
GREEN=$'\033[0;32m'
YELLOW=$'\033[1;33m'
BLUE=$'\033[0;34m'
CYAN=$'\033[0;36m'
BOLD=$'\033[1m'
NC=$'\033[0m' # No Color

# Function to print colored messages
print_error() {
    echo "${RED}✗ ERROR:${NC} $1"
}

print_success() {
    echo "${GREEN}✓ SUCCESS:${NC} $1"
}

print_info() {
    echo "${CYAN}ℹ INFO:${NC} $1"
}

print_warning() {
    echo "${YELLOW}⚠ WARNING:${NC} $1"
}

print_header() {
    echo ""
    echo "${BOLD}${BLUE}╔════════════════════════════════════════╗${NC}"
    echo "${BOLD}${BLUE}║${NC}   $1"
    echo "${BOLD}${BLUE}╚════════════════════════════════════════╝${NC}"
    echo ""
}

print_banner() {
    echo "${BOLD}${CYAN}"
    echo "  ╔═══════════════════════════════════════════════╗"
    echo "  ║   🚀 API Gateway & Microservices Manager    ║"
    echo "  ╚═══════════════════════════════════════════════╝"
    echo "${NC}"
}

# Function to display usage information
usage() {
    print_banner
    echo "Usage: $0 [choice] [env-file]"
    echo ""
    echo "Arguments:"
    echo "  [choice]              Menu choice (1/2/3) - optional, will prompt if not provided"
    echo "  [env-file]            Path to environment file (optional, auto-searches for .env.local, .env.dev, .env)"
    echo ""
    echo "Options:"
    echo "  -h, --help            Show this help message and exit"
    echo ""
    echo "${BOLD}Menu choices:${NC}"
    echo "  1                     Run the API Gateway"
    echo "  2                     Run the User Service"
    echo "  3                     Initialize the project environment (choose API Gateway or User Service)"
    echo ""
    echo "${BOLD}Required Environment Variables:${NC}"
    echo ""
    echo "For all options:"
    echo "  ENV                   Either 'dev' or 'prod'"
    echo ""
    echo "For option 1 (API Gateway):"
    echo "  API_GATEWAY_HOSTNAME     Hostname to listen on (e.g., localhost)"
    echo "  API_GATEWAY_TARGET_URL   Target URL to proxy to (e.g., http://localhost:8080)"
    echo "  API_GATEWAY_LISTEN_PORT  Port to listen on (e.g., 3000)"
    echo "  API_GATEWAY_DATABASE_URL PostgreSQL connection URL (e.g., postgresql://user:pass@localhost:5432/dbname)"
    echo "  GOOSE_DRIVER             Database driver for migrations (e.g., postgres)"
    echo "  GOOSE_DBSTRING           Database connection string for Goose (e.g., postgresql://user:pass@localhost/dbname)"
    echo "  GOOSE_MIGRATION_DIR      Path to migration files directory (e.g., migrations)"
    echo ""
    echo "For option 2 (User Service):"
    echo "  USER_SERVICE_HOSTNAME    Hostname to listen on (e.g., localhost)"
    echo "  USER_SERVICE_PORT        Port to run the service on (e.g., 8080)"
    echo "  USER_SERVICE_DB_URL      PostgreSQL connection URL (e.g., postgresql://user:pass@localhost:5432/dbname)"
    echo ""
    echo "${BOLD}Example:${NC}"
    echo "  $0 1                  Run API Gateway (reads from .env.local or similar)"
    echo "  $0 2 .env.prod        Run User Service with prod environment file"
    echo "  $0                    Show interactive menu"
    exit 1
}

# Handle help flag
if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    usage
fi

# Parse arguments - choice can be first or second arg, env file will be the other or second
CHOICE=""
ENV_FILE=""

# First argument could be choice (1-3) or env file
if [ -n "$1" ]; then
    if [[ "$1" =~ ^[1-3]$ ]]; then
        # First arg is choice
        CHOICE="$1"
        ENV_FILE="$2"
    else
        # First arg is env file
        ENV_FILE="$1"
        if [ -n "$2" ] && [[ "$2" =~ ^[1-3]$ ]]; then
            CHOICE="$2"
        fi
    fi
fi

# Validate provided env file exists
if [ -n "$ENV_FILE" ] && [ ! -f "$ENV_FILE" ]; then
    echo "Error: Environment file '$ENV_FILE' not found."
    exit 1
fi

# Auto-search for env files if none provided
if [ -z "$ENV_FILE" ]; then
    for file in .env.local .env.dev .env; do
        if [ -f "$file" ]; then
            ENV_FILE="$file"
            break
        fi
    done
fi

# Source env file if found
if [ -n "$ENV_FILE" ]; then
    source "$ENV_FILE"
    print_success "Environment file loaded: $ENV_FILE"
else
    print_warning "No environment file found (.env.local, .env.dev, or .env). Using system environment variables."
fi

# Check if node and go are setup correctly
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed or not in the PATH."
    exit 1
fi
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in the PATH."
    exit 1
fi

print_success "Node.js and Go are installed"

# Ask user what to do (or use pre-set choice)
if [ -z "$CHOICE" ]; then
    print_banner
    echo "${BOLD}What do you want to do?${NC}"
    echo ""
    echo "  1️⃣   Run the API Gateway"
    echo "  2️⃣   Run the User Service"
    echo "  3️⃣   Initialize the project environment"
    echo ""
    read -p "Enter your choice (1/2/3): " choice
else
    print_header "Executing Choice $CHOICE"
    choice="$CHOICE"
fi

case $choice in
    1)
        print_header "🌐 API GATEWAY"
        
        # Validate ENV is set and valid
        if [ -z "$ENV" ]; then
            print_error "ENV environment variable must be set."
            exit 1
        fi

        if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
            print_error "ENV must be either 'dev' or 'prod', but got '$ENV'."
            exit 1
        fi
        
        # Validate the API Gateway environment variables 
        if [ -z "$API_GATEWAY_HOSTNAME" ] || [ -z "$API_GATEWAY_TARGET_URL" ] || [ -z "$API_GATEWAY_LISTEN_PORT" ] || [ -z "$API_GATEWAY_DATABASE_URL" ] || [ -z "$GOOSE_DRIVER" ] || [ -z "$GOOSE_DBSTRING" ] || [ -z "$GOOSE_MIGRATION_DIR" ]; then
            print_error "Missing required environment variables:"
            echo "  - API_GATEWAY_HOSTNAME"
            echo "  - API_GATEWAY_TARGET_URL"
            echo "  - API_GATEWAY_LISTEN_PORT"
            echo "  - API_GATEWAY_DATABASE_URL"
            echo "  - GOOSE_DRIVER"
            echo "  - GOOSE_DBSTRING"
            echo "  - GOOSE_MIGRATION_DIR"
            exit 1
        fi

        # If the url and port present but not valid then show error message and exit
        if ! [[ "$API_GATEWAY_TARGET_URL" =~ ^http://[a-zA-Z0-9.-]+:[0-9]+$ ]]; then
            print_error "Invalid TARGET_URL format. It should be in the format http://hostname:port."
            exit 1
        fi
        if ! [[ "$API_GATEWAY_LISTEN_PORT" =~ ^[0-9]+$ ]] || [ "$API_GATEWAY_LISTEN_PORT" -le 0 ] || [ "$API_GATEWAY_LISTEN_PORT" -gt 65535 ]; then
            print_error "Invalid LISTEN_PORT. It should be a number between 1 and 65535."
            exit 1
        fi
        
        # Validate database URL format
        if ! [[ "$API_GATEWAY_DATABASE_URL" =~ ^postgresql:// ]]; then
            print_error "Invalid DATABASE_URL format. It should start with 'postgresql://'."
            exit 1
        fi
        
        # Validate Goose driver
        if ! [[ "$GOOSE_DRIVER" =~ ^(postgres|mysql|sqlite3|mssql)$ ]]; then
            print_error "Invalid GOOSE_DRIVER. Must be one of: postgres, mysql, sqlite3, mssql"
            exit 1
        fi
        
        # Validate migration directory exists
        if [ ! -d "./api-gateway/$GOOSE_MIGRATION_DIR" ]; then
            print_warning "Migration directory './api-gateway/$GOOSE_MIGRATION_DIR' does not exist. It will be created when needed."
        fi

        print_info "Configuration:"
        echo "  🔗 Listen Address  : ${BOLD}$API_GATEWAY_HOSTNAME:$API_GATEWAY_LISTEN_PORT${NC}"
        echo "  📦 Target URL      : ${BOLD}$API_GATEWAY_TARGET_URL${NC}"
        echo "  💾 Database URL    : ${BOLD}$API_GATEWAY_DATABASE_URL${NC}"
        echo "  🔄 Goose Driver    : ${BOLD}$GOOSE_DRIVER${NC}"
        echo "  📂 Migration Dir   : ${BOLD}$GOOSE_MIGRATION_DIR${NC}"
        echo "  ⚙️  Environment    : ${BOLD}$ENV${NC}"
        echo ""
        
        export ENV API_GATEWAY_HOSTNAME API_GATEWAY_TARGET_URL API_GATEWAY_LISTEN_PORT API_GATEWAY_DATABASE_URL GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

        # Check if env is dev or prod and start the api gateway accordingly
        if [ "$ENV" == "dev" ]; then
            print_info "Starting API Gateway in ${BOLD}dev${NC} mode with hot reload..."
            cd ./api-gateway && air
        elif [ "$ENV" == "prod" ]; then
            print_info "Building and starting API Gateway in ${BOLD}prod${NC} mode..."
            cd ./api-gateway && go build -o api-gateway && ./api-gateway
        else
            print_error "Invalid environment. Please set ENV to 'dev' or 'prod'."
            exit 1
        fi
        ;;
    2)
        print_header "👥 USER SERVICE"
        
        # Validate ENV is set
        if [ -z "$ENV" ]; then
            print_error "ENV environment variable must be set."
            exit 1
        fi
        
        if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
            print_error "ENV must be either 'dev' or 'prod', but got '$ENV'."
            exit 1
        fi
        
        # Validate USER_SERVICE_HOSTNAME, PORT, and DB URL are set
        if [ -z "$USER_SERVICE_HOSTNAME" ]; then
            print_error "USER_SERVICE_HOSTNAME environment variable must be set."
            exit 1
        fi
        
        if [ -z "$USER_SERVICE_PORT" ]; then
            print_error "USER_SERVICE_PORT environment variable must be set."
            exit 1
        fi

        if [ -z "$USER_SERVICE_DB_URL" ]; then
            print_error "USER_SERVICE_DB_URL environment variable must be set."
            exit 1
        fi
        
        if ! [[ "$USER_SERVICE_PORT" =~ ^[0-9]+$ ]] || [ "$USER_SERVICE_PORT" -le 0 ] || [ "$USER_SERVICE_PORT" -gt 65535 ]; then
            print_error "Invalid USER_SERVICE_PORT. It should be a number between 1 and 65535."
            exit 1
        fi

        if ! [[ "$USER_SERVICE_DB_URL" =~ ^postgresql:// ]]; then
            print_error "Invalid USER_SERVICE_DB_URL format. It should start with 'postgresql://'."
            exit 1
        fi
        
        print_info "Configuration:"
        echo "  🔗 Listen Address  : ${BOLD}$USER_SERVICE_HOSTNAME:$USER_SERVICE_PORT${NC}"
        echo "  💾 Database URL    : ${BOLD}$USER_SERVICE_DB_URL${NC}"
        echo "  ⚙️  Environment    : ${BOLD}$ENV${NC}"
        echo ""
        
        export USER_SERVICE_HOSTNAME USER_SERVICE_PORT USER_SERVICE_DB_URL ENV
        
        # Also export PORT for backward compatibility
        export PORT="$USER_SERVICE_PORT"
        
        # Run user service based on environment
        if [ "$ENV" == "dev" ]; then
            print_info "Starting User Service in ${BOLD}dev${NC} mode with watch mode..."
            cd ./services/user-service && npm run start:dev
        elif [ "$ENV" == "prod" ]; then
            print_info "Building and starting User Service in ${BOLD}prod${NC} mode..."
            cd ./services/user-service && npm run build && npm run start:prod
        else
            print_error "Invalid environment. Please set ENV to 'dev' or 'prod'."
            exit 1
        fi
        ;;
    3)
        print_header "⚙️  INITIALIZING PROJECT"
        echo "${BOLD}Which project do you want to initialize?${NC}"
        echo ""
        echo "  1️⃣   API Gateway (Go)"
        echo "  2️⃣   User Service (Node.js)"
        echo ""
        read -p "Enter your choice (1/2): " init_choice
        
        case $init_choice in
            1)
                print_info "Initializing API Gateway..."
                export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR
                cd ./api-gateway && sudo docker compose up -d && goose up && sqlc generate && go mod tidy && cd ../
                if [ $? -eq 0 ]; then
                    print_success "API Gateway initialized!"
                else
                    print_error "Failed to initialize API Gateway."
                    exit 1
                fi
                ;;
            2)
                print_info "Initializing User Microservice..."
                cd ./services/user-service && npm run init && cd ../../
                if [ $? -eq 0 ]; then
                    print_success "User Service initialized!"
                else
                    print_error "Failed to initialize User Service."
                    exit 1
                fi
                ;;
            *)
                print_error "Invalid choice. Please enter 1 or 2."
                exit 1
                ;;
        esac
        
        print_success "Project environment initialized!"
        ;;
    *)
        print_error "Invalid choice. Please enter 1, 2, or 3."
        usage
        ;;
esac