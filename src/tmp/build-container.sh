#!/bin/bash

REGION=europe-west4
REPOSITORY=cvk-demos
PROJECT=$(gcloud config get-value project)
ARTIFACT=fract-http

DOCKERBASE=$REGION-docker.pkg.dev
IMAGE=$DOCKERBASE/$PROJECT/$REPOSITORY/$ARTIFACT

echo
echo "Building container for project: "$PROJECT
echo "Handling image: "$IMAGE
echo docker build . --tag $IMAGE -f Dockerfile
docker build . --tag $IMAGE -f Dockerfile

echo
# one time configure of docker
#gcloud auth configure-docker $DOCKERBASE

echo
echo docker push $IMAGE
docker push $IMAGE


