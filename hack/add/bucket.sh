#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/buckets -d \
'{
    "createUser": "tackle",
    "name":"test"
}' | jq -M .
