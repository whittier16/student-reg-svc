#!/bin/sh

set -e

if [ $# -ne 1 ]
then
    echo "Usage: $0 [SERVICE NAME]"
    exit
fi

TAG="latest"
SERVICE_NAME=${1}

echo "building $SERVICE_NAME"

make linux

docker build -t ${SERVICE_NAME}:latest . -f ./build/package/Dockerfile
