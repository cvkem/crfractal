#!/bin/bash

REGION=europe-west4
REPOSITORY=cvk-demos
PROJECT=$(gcloud config get-value project)
ARTIFACT=fract-http

DOCKERBASE=$REGION-docker.pkg.dev
IMAGE=$DOCKERBASE/$PROJECT/$REPOSITORY/$ARTIFACT

echo "Building container for project: "$PROJECT
echo "Handling image: "$IMAGE

echo docker build . --tag $IMAGE -f Dockerfile
docker build . --tag $IMAGE -f Dockerfile

# one time configure
#gcloud auth configure-docker $DOCKERBASE

echo docker push $IMAGE
docker push $IMAGE


