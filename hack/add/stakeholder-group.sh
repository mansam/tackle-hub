#!/bin/bash

host="localhost:8080"

curl -X POST ${host}/controls/stakeholder-group -d \
'{
    "createUser": "tackle",
    "name": "Big Dogs",
    "username": "tackle",
    "description": "Group of big dogs."
}'