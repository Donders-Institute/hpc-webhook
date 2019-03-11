#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DATABASE" <<-EOSQL
INSERT INTO qaas (id, hash, groupname, username, script, description, created)
VALUES 
    (1, '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08', 'dccngroup', 'jonsno', 'test1.sh', 'Test script', '2019-03-11 10:21:00'),
    (2, '1286d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a10', 'dccngroup', 'foobar', 'test2.sh', 'Test script 2', '2019-03-11 11:21:00'),
    (3, '2086d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a12', 'dccngroup', 'somguy', 'test3.sh', '', '2019-03-11 12:21:00'),
    (4, '2486d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f42424', 'dccngroup', 'dccnuser', 'test4.sh', 'Tryout script', '2019-03-11 13:42:00');
EOSQL