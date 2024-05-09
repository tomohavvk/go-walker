#database
export DB_NAME="walker_service"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="root"
export DB_PASSWORD="password"

#server
export SERVER_HOST="192.168.0.153"
export SERVER_PORT="9001"
export EXCHANGE_API_TYPE="futures"

go run cmd/main.go