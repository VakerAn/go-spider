#!/usr/bin/env bash

# Assuming it is project directory
CURRENT_DIRECTORY=$(pwd)

CURRENT_DIRECTORY=${CURRENT_DIRECTORY}

CURRENT_DIRECTORY="$CURRENT_DIRECTORY"
DOCKER_COMPOSE_FILE='docker-compose.yml'

APPLICATION_NAME="gomovie.local"
APPLICATION_PORT=8899
DATABASE_CONTAINER_NAME='database-test'
STORE_NETWORK='store_network'
REDIS_CONTAINER_NAME='redis-test'
export COMPOSE_IGNORE_ORPHANS=true

echo 'Creating network...'
docker network create --driver bridge ${STORE_NETWORK} 2>&1 >> /dev/null
if [ "$?" -ne 0 ]; then
    echo 'Network already exists'
fi
echo 'Network is ok'

OS_TYPE=$(uname)
OS_DEBUG_HOST=host.docker.internal
if [ "$OS_TYPE" = "Linux" ]; then
  OS_DEBUG_HOST=172.17.0.1
fi
go install github.com/mitchellh/gox@latest
gox -osarch="linux/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
DEBUG_HOST=${OS_DEBUG_HOST} \
    DATABASE_CONTAINER_NAME=${DATABASE_CONTAINER_NAME} \
    REDIS_CONTAINER_NAME=${REDIS_CONTAINER_NAME} \
    APPLICATION_NAME=${APPLICATION_NAME} \
    APPLICATION_PORT=${APPLICATION_PORT} \
    NETWORK=${STORE_NETWORK} \
    DOCKER_COMPOSE_FILE=${DOCKER_COMPOSE_FILE} \
    docker-compose -f ${CURRENT_DIRECTORY}/docker-compose.yml up -d