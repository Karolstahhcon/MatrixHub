#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"

ROOT_DIR="$(realpath "${DIR}/..")"

PROTOVENDOR_DIR="${ROOT_DIR}/api/proto/vendor"
PROTO_DIR="${ROOT_DIR}/api/proto"
SWAGGER_DIR="${ROOT_DIR}/api/openapiv2"
TS_DIR="${ROOT_DIR}/api/ts"
GO_DIR="${ROOT_DIR}/api/go"
GOBIN="${ROOT_DIR}/bin"
PATH="${GOBIN}:${PATH}"

# Install protoc
install_protoc() {
  local os_name=$(uname -s | tr '[:upper:]' '[:lower:]')
  local arch_name=$(uname -m)
  local protoc_version="33.5"
  local protoc_binary="${ROOT_DIR}/bin/protoc"
  if [[ -x "${protoc_binary}" ]]; then
    # Extract version from "libprotoc <version>"
    local installed_version
    installed_version="$("${protoc_binary}" --version 2>/dev/null | awk '{print $2}')"
    if [[ "${installed_version}" == "${protoc_version}" ]]; then
      echo "protoc ${protoc_version} is already installed."
      return
    else
      echo "Found protoc ${installed_version}, but ${protoc_version} is required. Reinstalling..."
    fi
  fi

  if [[ "$arch_name" == "x86_64" || "$arch_name" == "amd64" ]]; then
    arch_name="x86_64"
  elif [[ "$arch_name" == "aarch64" || "$arch_name" == "arm64" ]]; then
    arch_name="aarch_64"
  else
    echo "Unsupported architecture: $arch_name"
    exit 1
  fi

  if [[ "$os_name" == "linux" ]]; then
    os_name="linux"
  elif [[ "$os_name" == "darwin" ]]; then
    os_name="osx"
  else
    echo "Unsupported OS: $os_name"
    exit 1
  fi

  local protoc_zip="protoc-${protoc_version}-${os_name}-${arch_name}.zip"
  local download_url="https://github.com/protocolbuffers/protobuf/releases/download/v${protoc_version}/${protoc_zip}"
  echo "${download_url}"

  wget -O "${protoc_zip}" "${download_url}"
