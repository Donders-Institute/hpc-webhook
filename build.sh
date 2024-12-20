#!/bin/bash

# export variables defined in .env
set -a && source .env && set +a
docker-compose -f docker-compose.yml build --force-rm 
