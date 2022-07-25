#!/bin/bash

source ./.env.sh

str1=fizz
str2=buzz
int1=3
int2=5
limit=10

curl -Ss -X GET -H "Content-Type: application/json" "http://localhost:${SERVER_PORT}/fizzbuzz?str1=${str1}&str2=${str2}&int1=${int1}&int2=${int2}&limit=${limit}"