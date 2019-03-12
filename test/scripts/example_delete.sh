#!/bin/bash
curl -X DELETE \
  http://localhost:5111/configuration/550e8400-e29b-41d4-a716-446655440001 \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "hash": "550e8400-e29b-41d4-a716-446655440001", 
  "groupname": "dccngroup",
  "username": "dccnuser",
  "description": ""
}'