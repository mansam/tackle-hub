#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/identities -d \
'{
    "createUser": "tackle",
    "kind": "git",
    "name":"jeff",
    "description": "Forklift",
    "user": "userA",
    "password": "passwordA",
    "key": "keyA",
    "settings": "settingsA",
    "repository": 1
}' | jq -M .

curl -X POST ${host}/repositories/1/identities -d \
'{
    "createUser": "tackle",
    "kind": "git",
    "name":"jeff",
    "description": "Forklift",
    "user": "userA",
    "password": "passwordA",
    "key": "keyA",
    "settings": "settingsA"
}' | jq -M .
