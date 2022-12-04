#!/usr/bin/env bash
set -x
set -eo pipefail

if ! [ -x "$(command -v migrate)" ]; then
  echo >&2 "Error: migrate is not installed."
  exit 1
fi

if ! [ -x "$(command -v migrate)" ]; then
  echo >&2 "Error: migrate is not installed."
  echo >&2 "Use:"
  echo >&2 "    brew install golang-migrate"
  echo >&2 "to install it."
  exit 1
fi

DB_USER=${POSTGRES_USER:=postgres}
DB_PASSWORD="${POSTGRES_PASSWORD:=password}"
DB_NAME="${POSTGRES_DB:=user}"
DB_PORT="${POSTGRES_PORT:=5432}"

# Allow to skip Docker if a dockerized Postgres database is already running
if [[ -z "${SKIP_DOCKER}" ]]
then
  docker build -t postgres-local -f ./scripts/Dockerfile ./scripts
  docker run \
      -e POSTGRES_USER=${DB_USER} \
      -e POSTGRES_PASSWORD=${DB_PASSWORD} \
      -e POSTGRES_DB=${DB_NAME} \
      -e POSTGRES_MULTIPLE_DATABASES=request \
      -p "${DB_PORT}":5432 \
      -d postgres-local \
      postgres -N 1000
fi

until PGPASSWORD="${DB_PASSWORD}" psql -h "localhost" -U "${DB_USER}" -p "${DB_PORT}" -d "postgres" -c '\q'; do
  >&2 echo "Postgres is still unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up and running on port ${DB_PORT} - running migrations now!"

export DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}
migrate -path ./user-service/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/user?sslmode=disable" up
migrate -path ./request-service/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/request?sslmode=disable" up

>&2 echo "Postgres has been migrated, ready to go!"