#!/bin/bash

set -e

if [ -z "${path}" ]; then
  echo '[!] Input $path missing!'
  exit 1
fi


if [ ! -d "${path}" ]; then
  if [[ "${is_create_path}" == "true" ]]; then
    mkdir -p "${path}"
  else
    echo "${path} doesn't exist, but is_create_path=false"
    exit 1
  fi
fi

cd "${path}"
expanded_target_path="$(pwd)"

echo "Change dir to ${expanded_target_path} ..."
envman add --key BITRISE_SOURCE_DIR --value "${expanded_target_path}"
