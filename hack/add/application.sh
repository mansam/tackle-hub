#!/bin/bash

host="localhost:8080"

curl -X POST ${host}/application-inventory/application -d \
'{
    "createUser": "tackle",
    "name":"jeff",
    "description": "Forklift",
    "businessService": {
      "id": 1
    },
    "tags":[
      "1",
      "2"
    ]
}'
