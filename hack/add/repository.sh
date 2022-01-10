#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/repositories -d \
'{
    "createUser": "tackle",
    "application": 1,
    "name": "created-directly",
    "kind": "git",
    "url": "git://github.com/testing"
}' | jq -M .

curl -X POST ${host}/application-inventory/application/1/repositories -d \
'{
    "createUser": "tackle",
    "name": "created-for-application",
    "kind": "git",
    "url": "git://github.com/testing"
}' | jq -M .
