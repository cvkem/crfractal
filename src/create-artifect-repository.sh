#!/bin/bash


PROJECT=$(gcloud config get-value project)

echo "Creating the repository in folder "${PROJECT}

# read exactly one character (-n 1) and treat \ as a normal character
# as no varname is supplied the result is read in $REPLY
read -p "Are you sure  [y/n]? " -n 1 -r
echo    # (optional) move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo cancelled
    exit 1
fi

gcloud services enable artifactregistry.googleapis.com

gcloud artifacts repositories create \
  --location europe-west4 \
  --repository-format docker \
  cvk-demos
