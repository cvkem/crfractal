#!/bin/bash

REGION=europe-west4
REPOSITORY=cvk-demos
PROJECT=$(gcloud config get-value project)
ARTIFACT=fract-http

DOCKERBASE=$REGION-docker.pkg.dev
IMAGE=$DOCKERBASE/$PROJECT/$REPOSITORY/$ARTIFACT

gcloud run deploy $ARTIFACT \
      	--image $IMAGE \
	--concurrency 2 \
      	--allow-unauthenticated

