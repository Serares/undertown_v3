#!/bin/bash

# Environment variables for PostgreSQL connection
PG_USER=${PSQL_USER}
PG_PASS=${PSQL_PASSWORD}
PG_HOST=${PSQL_HOST}
PG_PORT=${PSQL_PORT:-5432} # Default port is 5432
PG_DBNAME=${PSQL_DB}
SSL_MODE=${SSL_MODE:-disable} # Default sslmode is disable
GOOSE_COMMAND=$1

# // check if goose_command is empty and print a message before exiting the script
if [ -z "$GOOSE_COMMAND" ]; then
    echo "Please provide a Goose command as an argument."
    exit 1
fi

# Construct the PostgreSQL connection string
PG_CONNECTION_STRING="postgres://${PG_USER}:${PG_PASS}@${PG_HOST}:${PG_PORT}/${PG_DBNAME}?sslmode=${SSL_MODE}"

# Check if Goose is installed
if ! command -v goose &>/dev/null; then
    echo "Goose could not be found, please install it before proceeding."
    exit 1
fi

# Run Goose migrations
echo "Running Goose migrations..."
echo "Connection string: $PG_CONNECTION_STRING"
cd "repositories/migrations/schema"
goose postgres "user=$PG_USER dbname=$PG_DBNAME host=$PG_HOST port=$PG_PORT password="$PG_PASS"" $GOOSE_COMMAND

echo "Migration completed."
