#!/bin/bash

NAME="$@"

curl -X 'GET'  -u ":$MANETU_PAT" \
      --data-urlencode "name=$NAME" \
      "$MANETU_URL/api/lambda/hello/greet" \
      -H 'accept: text/plain'
echo ""

