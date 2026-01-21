#!/bin/bash
set -e

# ===== Config =====
DOCKER_USERNAME="nanil0034"
DOCKER_PASSWORD="dckr_pat___**"

IMAGE_NAME="nanil0034/kindergarten-registry-student"
TAG="${TAG:-96}"

# ENV values
VAULT_ADDR="http://192.168.121.132:8200"
VAULT_TOKEN="**"
SERVICE_NAME="student"

FULL_IMAGE="${IMAGE_NAME}:${TAG}"

echo "Logging into Docker Hub..."
echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

echo "Building Docker image: ${FULL_IMAGE}"
# Using '.' as the build context if you're inside the employeeservice directory
docker build \
  --build-arg VAULT_ADDR="${VAULT_ADDR}" \
  --build-arg VAULT_TOKEN="${VAULT_TOKEN}" \
  --build-arg SERVICE_NAME="${SERVICE_NAME}" \
  -t "${FULL_IMAGE}" \
  .

echo "Pushing Docker image..."
docker push "${FULL_IMAGE}"

echo "Done  Image pushed: ${FULL_IMAGE}"
