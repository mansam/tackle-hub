#!/bin/bash

host="localhost:8080"

curl -X POST ${host}/tasks/1/report -d \
'{
    "createUser": "tackle",
    "updateUser": "tackle",
    "status": "Running",
    "total": 10,
    "completed": 0,
    "detail": "addon started."
}'
