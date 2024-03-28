#!/bin/bash

cat <<EOF | yq -r -o=json - | curl --data-binary @- --silent --show-error --fail --header 'Content-Type: application/json' -u ":$MANETU_PAT" $MANETU_GRAPHQL_URL | jq
query: |
  {
    get_profile {
       name
       userid
    }
  }
EOF
