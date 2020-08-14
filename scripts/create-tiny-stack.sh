#!/usr/bin/env bash

set -euo pipefail

usage() { echo "Usage: $0 [-v <image version>]" 1>&2; exit 1; }
version=dev

# add stack arg
while getopts "v" o; do
  case "${o}" in
    v)
      version=${OPTARG}
      [[ -n ${version} ]] || usage
      ;;
    *)
      usage
      ;;
  esac
done

docker pull ubuntu:bionic

scripts_dir=$(cd "$(dirname $0)" && pwd)
base_dir=${scripts_dir}/..
bionic_dir=${base_dir}/bionic
tiny_dir=${base_dir}/tiny
arch="x86_64"

# Define tags
base_build=build:${version}-base
tiny_run=run:${version}-tiny

cnb_tiny_build=build:${version}-tiny-cnb
cnb_tiny_run=run:${version}-tiny-cnb

# Build base images
docker build -t "${base_build}" \
  --build-arg "sources=$(cat "${base_dir}/arch/${arch}/sources.list")" \
  --build-arg "packages=$(cat "${base_dir}/packages/base/build")" \
  --no-cache "${bionic_dir}/dockerfile/build"
docker build -t "${tiny_run}" "$tiny_dir/dockerfile/run"

# Build CNB images
docker build --build-arg "base_image=${base_build}" --build-arg "stack_id=io.paketo.stacks.tiny" -t "${cnb_tiny_build}"  "$bionic_dir/cnb/build"
docker build --build-arg "base_image=${tiny_run}" -t "${cnb_tiny_run}" "$tiny_dir/cnb/run"

echo "To publish these images:"
for image in "${base_build}" "${tiny_run}" "${cnb_tiny_build}" "${cnb_tiny_run}"; do
  echo "  docker push ${image}"
done
