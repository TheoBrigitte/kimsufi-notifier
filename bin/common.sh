#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

declare -A OVH_API_ENDPOINTS
OVH_API_ENDPOINTS["ovh-eu"]="https://eu.api.ovh.com/v1"
OVH_API_ENDPOINTS["ovh-ca"]="https://ca.api.ovh.com/v1"
OVH_API_ENDPOINTS["ovh-us"]="https://api.us.ovhcloud.com/1.0"

JQ_VERSION="1.7.1"
JQ_BIN="${SCRIPT_DIR}/jq"

install_tool() {
    local binary="$(basename $1)" ; shift
    local version=$1 ; shift
    local url=$1 ; shift
    local destination="${SCRIPT_DIR}/${binary}"

    if [[ ! -f "${destination}" ]]; then
        echo "> Installing ${binary}"
        (
          set -x
          curl -fsL -o "${destination}" "${url}"
          chmod +x "${destination}"
        )
    fi
}

install_tools() {
  os=""
  case "$OSTYPE" in
    darwin*)  os="macos" ;;
    linux*)   os="linux" ;;
    *)        echo "unsupported OS: $OSTYPE"; exit ;;
  esac

  arch=""
  case "$(uname -m)" in
    x86_64) arch="amd64" ;;
    arm*)   arch="arm64" ;;
    aarch*) arch="arm64" ;;
    *)     echo "unsupported architecture: $(uname -m)"; exit ;;
  esac

  install_tool "${JQ_BIN}" "${JQ_VERSION}" "https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-${os}-${arch}"
}
