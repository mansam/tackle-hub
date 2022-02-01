#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/application-inventory/application/1/repository -d \
'{
    "createUser": "tackle",
    "name": "created-for-application",
    "kind": "git",
    "url": "git://github.com/testing"
}' | jq -M .
