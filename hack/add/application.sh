#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/application-inventory/application -d \
'{
    "createUser": "tackle",
    "name":"Dog",
    "description": "Dog application.",
    "businessService": "1",
    "tags":[
      "1",
      "2"
    ]
}' | jq -M .

curl -X POST ${host}/application-inventory/application -d \
'{
    "createUser": "tackle",
    "name":"Cat",
    "description": "Cat application.",
    "repository": {
      "name": "Cat",
      "kind": "git",
      "url": "git://github.com/pet/cat",
      "branch": "/cat"
    },
    "businessService": "1",
    "tags":[
      "1",
      "2"
    ]
}' | jq -M .

