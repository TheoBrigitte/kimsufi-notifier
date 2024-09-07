#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE}") && pwd -P)

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
  install_tool "${JQ_BIN}" "${JQ_VERSION}" "https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-linux64"
}
