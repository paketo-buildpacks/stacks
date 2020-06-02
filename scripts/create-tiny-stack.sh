#!/usr/bin/env bash

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
base_dir=${scripts_dir}/../base
tiny_dir=${scripts_dir}/../tiny

# Define tags
base_build=cloudfoundry/build:${version}-base
base_run=cloudfoundry/run:${version}-base

cnb_base_build=cloudfoundry/build:${version}-base-cnb
cnb_base_run=cloudfoundry/run:${version}-base-cnb

# Build base images
docker build -t "${base_build}" "$base_dir/dockerfile/build"
docker build -t "${base_run}" "$tiny_dir/dockerfile/run"

# Build CNB images
docker build --build-arg "base_image=${base_build}" --build-arg "stack_id=io.paketo.stacks.tiny" -t "${cnb_base_build}"  "$base_dir/cnb/build"
docker build --build-arg "base_image=${base_run}" -t "${cnb_base_run}" "$tiny_dir/cnb/run"

echo "To publish these images:"
for image in "${base_build}" "${base_run}" "${cnb_base_build}" "${cnb_base_run}"; do
  echo "  docker push ${image}"
done
