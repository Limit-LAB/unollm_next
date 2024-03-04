#!/usr/bin/env bash

set -exo pipefail

echo ${FEATURES}

DIR=$(dirname $0)

COMMIT_HASH=$(bash "${DIR}"/commit-hash.sh)

cd ${DIR}/../ && \
    DOCKER_BUILDKIT=1 docker build -t ghcr.io/limit-lab/unollm_next:"${COMMIT_HASH}" \
    --build-arg GITHUB_USERNAME="${GITHUB_USERNAME}" \
    --build-arg GITHUB_TOKEN="${GITHUB_TOKEN}" \
    -f ./Dockerfile ./

docker tag ghcr.io/limit-lab/unollm_next:"${COMMIT_HASH}" ghcr.io/limit-lab/unollm_next:latest

if [ ! -z ${IMAGE_TAG} ]; then
    docker tag ghcr.io/limit-lab/unollm_next:"${COMMIT_HASH}" ghcr.io/limit-lab/unollm_next:${IMAGE_TAG}
fi
