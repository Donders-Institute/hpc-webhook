#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DATABASE" <<-EOSQL
INSERT INTO qaas (id, hash, username)
VALUES 
    (1, '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08', 'jonsno'),
    (2, '1286d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a10', 'foobar'),
    (3, '2086d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a12', 'somguy'),
    (4, '2486d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f42424', 'dccnuser');
EOSQL