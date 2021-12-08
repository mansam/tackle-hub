#!/bin/bash

host="localhost:8080"

curl -X POST ${host}/addons/test/tasks -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"Test",
    "data": {
      "application": 1
    }
}'

