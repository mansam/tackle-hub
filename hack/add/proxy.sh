#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/proxies -d \
'{
    "createUser": "tackle",
    "kind": "http",
    "host":"myhost",
    "port": 80,
    "user": "userA",
    "password": "passwordA"
}' | jq -M .

curl -X POST ${host}/proxies -d \
'{
    "createUser": "tackle",
    "kind": "https",
    "host":"myhost",
    "port": 443,
    "user": "userB",
    "password": "passwordB"
}' | jq -M .
