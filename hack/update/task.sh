#!/bin/bash

host="localhost:8080"

curl -X PUT ${host}/tasks/1/report -d \
'{
    "createUser": "tackle",
    "updateUser": "tackle",
    "status": "Running",
    "total": 10,
    "completed": 9,
    "detail": "reading /files/application/dog.java."
}'
