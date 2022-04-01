#!/bin/bash
IMAGE_NAME="${1:-"ghcr.io/chuxel/lifecycle/lifecycle"}"
LIFECYCLE_IMAGE_TAG="${2:-latest}"

cd "$(dirname "${BASH_SOURCE[0]}")/.."

echo "Building..."
make build
make package

echo "Creating images..."
DOCKER_CLI_EXPERIMENTAL=enabled

#LINUX_AMD64_SHA=$(go run ./tools/image/main.go -lifecyclePath ./out/lifecycle-v*+linux.x86-64.tgz -tag ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}-linux-x86-64 | awk '{print $NF}')
#echo "LINUX_AMD64_SHA: $LINUX_AMD64_SHA"
#LINUX_ARM64_SHA=$(go run ./tools/image/main.go -lifecyclePath ./out/lifecycle-v*+linux.arm64.tgz -tag ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}-linux-arm64 -arch arm64 | awk '{print $NF}')
#echo "LINUX_ARM64_SHA: $LINUX_ARM64_SHA"
#WINDOWS_AMD64_SHA=$(go run ./tools/image/main.go -lifecyclePath ./out/lifecycle-v*+windows.x86-64.tgz -tag ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}-windows -os windows | awk '{print $NF}')
#echo "WINDOWS_AMD64_SHA: $WINDOWS_AMD64_SHA"
#docker manifest create ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG} \
#    ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}-linux-x86-64@${LINUX_AMD64_SHA} \
#    ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}-linux-arm64@${LINUX_ARM64_SHA} \
#    ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}-windows@${WINDOWS_AMD64_SHA}


os=$(uname -s | tr '[:upper:]' '[:lower:]')
architecture=$(uname -m | tr '_' '-')
go run ./tools/image/main.go -daemon -lifecyclePath ./out/lifecycle-v*+${os}.${architecture}.tgz -tag ${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}
echo "${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG} created for ${os} ${architecture}."
