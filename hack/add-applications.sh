#!/bin/bash

host="localhost:8080"

#######################################################
# APPLICATIONS
#######################################################

curl -X POST ${host}/application-inventory/application -d \
'{
    "createUser": "tackle",
    "name":"jeff",
    "tags":["1","2"]
}'
