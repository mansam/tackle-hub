#!/bin/bash

host="${HOST:-localhost:8080}"

curl -X POST ${host}/application-inventory/application -d \
'{
    "createUser": "tackle",
    "name":"jeff",
    "description": "Forklift",
    "businessService": "1",
    "tags":[
      "1",
      "2"
    ]
}'
