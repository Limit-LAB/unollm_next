#!/usr/bin/env bash

set -euxo pipefail

HASH=$(git rev-list -1 HEAD --abbrev-commit)
if [[ $(git status --porcelain) ]]; then
  HASH=${HASH}-dirty
fi 

echo $HASH
