#!/bin/bash
curl -X PPUT \
  http://localhost:5111/configuration \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "hash": "550e8400-e29b-41d4-a716-446655440001", 
  "username": "dccnuser"
}'