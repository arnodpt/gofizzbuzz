#!/bin/bash

source ./.env.sh

curl -Ss -X GET -H "Content-Type: application/json" "http://localhost:${SERVER_PORT}/popular"