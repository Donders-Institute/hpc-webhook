#!/bin/bash
echo Change Windows to Unix line endings
sed -i 's/\r//g' ../init/01-initialize-database.sh
sed -i 's/\r//g' ../init/02-fill-database.sh
sed -i 's/\r//g' ../../scripts/wait-for-it.sh
docker-compose -f ../../docker-compose-test.yml build
docker-compose -f ../../docker-compose-test.yml up