#!/usr/bin/env bash

set -x
set -eo pipefail

DB_USER=${POSTGRES_USER:=postgres}
DB_PASSWORD=${POSTGRES_PASSWORD:=password}
DB_NAME=${POSTGRES_DB:=newsletter}
DB_PORT=${POSTGRES_PORT:=5432}
DB_HOST=${POSTGRES_HOST:=localhost}

if ! [[ -x "$(command -v goose)" ]]; then
    echo >&2 "ERROR: goose is not installed."
    echo >&2 "Use:"
    echo >&2 "go install github.com/pressly/goose/v3/cmd/goose@latest"
    echo >&2 "to install it"
    exit 1
fi

if [[ -z "${SKIP_DOCKER}" ]]; then
    docker run \
        -e POSTGRES_USER="${DB_USER}" \
        -e POSTGRES_PASSWORD="${DB_PASSWORD}" \
        -e POSTGRES_DB="${DB_NAME}" \
        -p "${DB_PORT}":5432 \
        -d postgres \
        postgres -N 1000
fi

export PGPASSWORD="${DB_PASSWORD}"
until psql -h "${DB_HOST}" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c '\q'; do
    echo >&2 "Postgres is still unavailable - sleeping"
    sleep 1
done

echo >&2 "Postgres is up and running on port ${DB_PORT} - running migrations now"

if [[ -z "${PROD}" ]]; then
    GOOSE_DRIVER=postgres
    GOOSE_DBSTRING=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}
    export GOOSE_DRIVER
    export GOOSE_DBSTRING

    mkdir -p migrations/
    cd migrations/
    #goose create create_subscription_table sql
    goose up
    echo >&2 "Postgres has been migrated, ready to go!"
    cd ..
fi
