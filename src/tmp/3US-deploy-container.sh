#!/bin/bash

REGION=europe-west4
DEPLOY_REGION=us-west2
REPOSITORY=cvk-demos2
PROJECT=$(gcloud config get-value project)
ARTIFACT=fract-http2

DOCKERBASE=$REGION-docker.pkg.dev
IMAGE=$DOCKERBASE/$PROJECT/$REPOSITORY/$ARTIFACT

gcloud run deploy $ARTIFACT \
      	--image $IMAGE \
	--concurrency 2 \
	--region=${DEPLOY_REGION} \
      	--allow-unauthenticated

