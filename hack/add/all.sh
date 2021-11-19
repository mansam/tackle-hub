#!/bin/bash

host="localhost:8080"

#######################################################
# ALL
#######################################################

dir=`dirname $0`
cd ${dir}

./tag.sh
./business-service.sh
./application.sh
./review.sh

