#!/bin/bash

REGION=europe-west4
REPOSITORY=cvk-demos2
PROJECT=$(gcloud config get-value project)
ARTIFACT=fract-http2

DOCKERBASE=$REGION-docker.pkg.dev
IMAGE=$DOCKERBASE/$PROJECT/$REPOSITORY/$ARTIFACT

gcloud run deploy $ARTIFACT \
      	--image $IMAGE \
	--region=${REGION} \
      	--allow-unauthenticated

