#!/bin/bash

DOCKER=docker
DOCKER_PATH=$(command -v $DOCKER)

if [ "$CI_COMMIT_REF_NAME" = "" ]; then
    echo "Starting PostgreSQL docker..."

    if [ -z "$DOCKER_PATH" ]; then
        echo "Docker not found. Checking for podman."

        DOCKER=podman
        PODMAN_PATH=$(command -v podman)
        if [ -z "$PODMAN_PATH" ]; then
            echo "Podman not found. No suitable container engine."
            exit 127
        fi
    fi

    $DOCKER run -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=local_db -d --name=postgres postgres:16-alpine

    printf "\nWaiting for PostgreSQL to be fully available."
    until $DOCKER exec postgres pg_isready >/dev/null 2>/dev/null; do
        printf "."
        sleep 5
    done

    echo "\n PostgreSQL ready!"
else
    echo "PostgreSQL docker provided as a service by GitLab CI"
fi

echo "Running PostgreSQL migrations..."

sh ./run_migration.sh
