#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DATABASE" <<-EOSQL
    DROP TABLE IF EXISTS qaas;
    CREATE TABLE qaas(
        id         SERIAL PRIMARY KEY,
        hash       CHAR (64)        NOT NULL,
        username   VARCHAR (32)     NOT NULL,
    );
EOSQL