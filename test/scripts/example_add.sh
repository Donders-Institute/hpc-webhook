#!/bin/bash
curl -X PUT \
  http://localhost:443/configuration \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "hash": "550e8400-e29b-41d4-a716-446655440001", 
  "groupname": "dccngroup",
  "username": "dccnuser",
  "script": "qsub.sh",
  "description": "",
  "created": "2019-03-11 10:10:25"
}'