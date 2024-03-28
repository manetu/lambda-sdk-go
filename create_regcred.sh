#!/bin/bash

if [ $# -ne 2 ]
then
    echo "usage: `basename $0` docker-host-port docker-user "
    echo '       with environemt variable $DOCKER_PAT set to docker credential'
    exit 1
fi

HOSTPORT=$1
REGISTRY_USER=$2
REGISTRY_PASSWORD=$DOCKER_PAT

cat <<EOF | yq -r -o=json - | curl --data-binary @- --silent --show-error --fail --header 'Content-Type: application/json' -u ":$MANETU_PAT" $MANETU_GRAPHQL_URL | jq
query: |
  mutation
  {
    create_basicauth_regcred ( 
        hostport: "$HOSTPORT"
        username: "$REGISTRY_USER"
        password: "$REGISTRY_PASSWORD"
    ){
        hostport
        last_updated
        mrn
        version
    }
  }
EOF

