#!/bin/bash

host="localhost:8080"

#######################################################
# APPLICATIONS
#######################################################

curl -X POST ${host}/application-inventory/application -d \
'{
    "createUser": "tackle",
    "name":"jeff",
    "description": "Forklift",
    "tags":[
      "1",
      "2"
    ]
}'
