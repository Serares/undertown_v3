#!/bin/bash

# Environment variables for PostgreSQL connection
PG_USER=${PSQL_USER}
PG_PASS=${PSQL_PASSWORD}
PG_HOST=${PSQL_HOST}
PG_PORT=${PSQL_PORT:-5432} # Default port is 5432
PG_DBNAME=${PSQL_DB}
SSL_MODE=${SSL_MODE:-disable} # Default sslmode is disable

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
goose -dir repositories/migrations/schema "user=$PG_USER dbname=$PG_DBNAME host=$PG_HOST port=$PG_PORT password=$PG_PASS sslmode=disable" up

# If you want Goose to run down migrations, use:
# goose -dir prepositories/migrations/schema "$PG_CONNECTION_STRING" down

echo "Migration completed."
