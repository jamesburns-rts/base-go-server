#!/bin/bash

cd `dirname $0`

PROJECT_NAME=$(grep project_name ../version.properties | cut -d'=' -f2 | tr -d '\r\n')
PROJECT_DOCKER_IMAGE_LOCAL=$(grep project_docker_image_local ../version.properties | cut -d'=' -f2 | tr -d '\r\n')
PROJECT_VERSION=$(grep project_version ../version.properties | cut -d'=' -f2 | tr -d '\r\n')

echo "Building, testing, and publishing $PROJECT_NAME $PROJECT_VERSION with '$@'"

VERSION=${1:-${PROJECT_VERSION}}
BRANCH=${2:-$(git rev-parse --abbrev-ref HEAD)}

echo "Building packages for branch $BRANCH and version $VERSION"

docker build -f ../build/package/Dockerfile \
    --build-arg BRANCH=$BRANCH \
    --build-arg PROJECT_NAME=$PROJECT_NAME \
    --build-arg VERSION=$VERSION \
    -t $PROJECT_DOCKER_IMAGE_LOCAL:latest \
    -t $PROJECT_DOCKER_IMAGE_LOCAL:$PROJECT_VERSION \
    ..

