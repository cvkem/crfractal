#!/bin/bash

ACCOUNT=docker-artifact-repo-push@cvk-sandbox.iam.gserviceaccount.com
LOCATION=europe-west4

echo
echo "This login is valid for an hour"

#gcloud auth print-access-token \
#  --impersonate-service-account  ${ACCOUNT} | docker login \
#  -u oauth2accesstoken \
#  --password-stdin https://${LOCATION}-docker.pkg.dev


gcloud auth print-access-token | docker login \
  -u oauth2accesstoken \
  --password-stdin https://${LOCATION}-docker.pkg.dev
