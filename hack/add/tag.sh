#!/bin/bash

host="${HOST:-localhost:8080}"

#
# Types
#

curl -X POST ${host}/controls/tag-type -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"Testing",
    "colour": "#807ded",
    "rank": 0
}'

curl -X POST ${host}/controls/tag-type -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"Operating System",
    "colour": "#807ded",
    "rank": 10
}'

curl -X POST ${host}/controls/tag-type -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"Database",
    "colour": "#8aed7d",
    "rank": 20
}'

curl -X POST ${host}/controls/tag-type -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"Language",
    "colour": "#ede97d",
    "rank": 30
}'

#
# Tags
#

curl -X POST ${host}/controls/tag -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"RHEL",
    "tagType": {"id":1}
}'

curl -X POST ${host}/controls/tag -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"PostgreSQL",
    "tagType": {"id":2}
}'

curl -X POST ${host}/controls/tag -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"C++",
    "tagType": {"id":3}
}'

curl -X POST ${host}/controls/tag -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"CRAZY",
    "tagType": {
      "createUser": "tackle",
      "username": "tackle",
      "name":"CRAZY",
      "colour": "#0000",
      "rank": 40
    }
}'

curl -X POST ${host}/controls/tag -d \
'{
    "createUser": "tackle",
    "username": "tackle",
    "name":"CRAZY-TRAIN",
    "tagType": {
      "id": 4,
      "createUser": "tackle",
      "username": "tackle",
      "name":"CRAZY-TRAIN",
      "colour": "#66666",
      "rank": 40
    }
}'
