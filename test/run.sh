#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

DISTRO="${1}"

scripts="check.sh;list.sh;order.sh --dry-run -i dedicated_datacenter=rbx -i region=europe"

echo "> Building image for ${DISTRO}"
docker build -t "${DISTRO}-test" -f "${SCRIPT_DIR}/Dockerfile.${DISTRO}" .
docker run --rm -v "${SCRIPT_DIR}/..:/usr/local" "${DISTRO}-test" sh -c "${scripts}"
