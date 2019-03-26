#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DATABASE" <<-EOSQL
    DROP TABLE IF EXISTS hpc_webhook;
    CREATE TABLE hpc_webhook(
        id          SERIAL PRIMARY KEY,
        hash        CHAR (36) UNIQUE NOT NULL,
        groupname   VARCHAR (32) NOT NULL,
        username    VARCHAR (32) NOT NULL,
        description VARCHAR (255),
        created     TIMESTAMP NOT NULL);
EOSQL