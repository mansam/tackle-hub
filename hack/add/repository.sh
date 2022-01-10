#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/repositories -d \
'{
    "createUser": "tackle",
    "application": {"id":1},
    "name": "created-directly",
    "type": "git",
    "url": "git://testing"
}' | jq -M .

curl -X POST ${host}/application-inventory/application/1/repositories -d \
'{
    "createUser": "tackle",
    "name": "created-for-application",
    "type": "git",
    "url": "git://testing"
}' | jq -M .
