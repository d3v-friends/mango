#!/bin/bash

ORIGIN=$1
if [ -z "$ORIGIN" ]; then
	ORIGIN="origin"
fi

TAG="v1.0.0"
LATEST="v1.0.0-latest"

# delete tag
git tag -d "$TAG"
git tag -d "$LATEST"
git push -d "$ORIGIN" "$TAG"
git push -d "$ORIGIN" "$LATEST"

# create tag
git tag "$TAG"
git tag "$LATEST"
git push --tags "$ORIGIN" "$TAG"
git push --tags "$ORIGIN" "$LATEST"