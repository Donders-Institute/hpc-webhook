#!/bin/bash
echo Change Windows to Unix line endings
sed -i 's/\r//g' ../init/01-initialize-database.sh
sed -i 's/\r//g' wait-for-it.sh
docker-compose -f ../docker-compose.yml build
docker-compose -f ../docker-compose.yml up