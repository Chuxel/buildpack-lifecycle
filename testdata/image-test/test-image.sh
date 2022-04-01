#!/bin/bash
set -e
BUILD_PATH="${1:-"."}"
IMAGE_NAME="${2:-"ghcr.io/chuxel/lifecycle/lifecycle"}"
LIFECYCLE_IMAGE_TAG="${3:-latest}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}/../../"

echo "(*) Building lifecycle binaries..."
make clean
make build-linux-amd64 
echo "(*) Packaging lifecycle binaries..."
make package-linux-amd64 
cp -f out/lifecycle-v*+linux.x86-64.tgz "${SCRIPT_DIR}/lifecycle.tgz"

echo "(*) Creating builder.."
cd "${SCRIPT_DIR}"
pack builder create "ghcr.io/chuxel/lifecycle/builder-devcontainer-empty" --pull-policy if-not-present -c builder-devcontainer.toml
rm "${SCRIPT_DIR}/lifecycle.tgz"

pack build lifecycle-image-test-output \
    --clear-cache -v \
    -e BP_DCNB_BUILD_MODE=devcontainer \
    --pull-policy if-not-present \
    --buildpack paketo-buildpacks/go \
    --builder ghcr.io/chuxel/lifecycle/builder-devcontainer-empty \
    --trust-builder \
    --lifecycle-image "${IMAGE_NAME}:${LIFECYCLE_IMAGE_TAG}" \
    -p "${BUILD_PATH}"

echo "(*) Listing layers in output.."
docker run -it --rm --entrypoint /cnb/lifecycle/launcher lifecycle-image-test-output ls /layers