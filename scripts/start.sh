#!/bin/bash
echo Change Windows to Unix line endings
sed -i 's/\r//g' ../init/01-initialize-database.sh
sed -i 's/\r//g' ../cmd/client/wait-for-it.sh
sed -i 's/\r//g' ../cmd/server/wait-for-it.sh
docker-compose -f ../deployments/docker-compose.yml build
docker-compose -f ../deployments/docker-compose.yml up