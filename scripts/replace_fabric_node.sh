#!/bin/bash

DOCKER_ID=$1

docker_cp_to_build () {
  output=".build$1"

  echo "cp $1 from docker $DOCKER_ID to $output"
  set -x
  rm -rf "$output"
  mkdir -p "$output"
  docker cp "$DOCKER_ID:$1/." "$output/"
  set +x
}

generate_env () {
  echo Generate ENV to $1
  docker exec $DOCKER_ID env \
      | grep -e ^CORE -e ^FABRIC -e ^ORDER \
      | sed 's/\/host\/var/\/var/' \
      | sed "s#/var/hyperledger#$PWD/.build/var/hyperledger#" \
      | sed "s#/etc/hyperledger#$PWD/.build/etc/hyperledger#" \
      > $1
}

docker_cp_to_build /var/hyperledger/
docker_cp_to_build /etc/hyperledger
generate_env docker.env

docker stop $docker_id
