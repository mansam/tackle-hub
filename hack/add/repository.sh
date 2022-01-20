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
