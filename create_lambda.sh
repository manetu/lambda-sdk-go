#!/bin/bash

SITE_FILE=$1
BASE_64_CONFIG=$(cat $SITE_FILE | base64 )

cat <<EOF | yq -r -o=json - | curl --data-binary @-  --silent --show-error --fail --header 'Content-Type: application/json' -u ":$MANETU_PAT" $MANETU_GRAPHQL_URL | jq
query: |
  mutation
  {
    create_lambda ( 
        config: "$BASE_64_CONFIG"
    ){
        last_updated
        mrn
        paused
        version
    }
  }
EOF

