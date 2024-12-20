#!/bin/bash

# export variables defined in .env
set -a && source .env && set +a
docker-compose -p hpc-webhook up -d
