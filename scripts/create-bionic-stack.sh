#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 -s <base/full> [-b <build-dest>] [-r <run-dest>] [-v <image version>]" 1>&2
  exit 1
}

build_images() {
  stack_name=$1
  version=$2
  build_dest=$3
  run_dest=$4

  arch="x86_64"

  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  stack_dir=${base_dir}/bionic

  build="${build_dest}:${version}-${stack_name}"
  run="${run_dest}:${version}-${stack_name}"

  cnb_build="${build_dest}:${version}-${stack_name}-cnb"
  cnb_run="${run_dest}:${version}-${stack_name}-cnb"

  docker pull ubuntu:bionic

  docker build -t "${build}" \
    --build-arg "sources=$(cat "${base_dir}/arch/${arch}/sources.list")" \
    --build-arg "packages=$(cat "${base_dir}/packages/${stack_name}/build")" \
    --no-cache "${stack_dir}/dockerfile/build"

  docker build -t "${run}" \
    --build-arg "sources=$(cat "${base_dir}/arch/${arch}/sources.list")" \
    --build-arg "packages=$(cat "${base_dir}/packages/${stack_name}/run")" \
    --no-cache "${stack_dir}/dockerfile/run"

  docker build --build-arg "base_image=${build}" -t "${cnb_build}"  "${stack_dir}/cnb/build"
  docker build --build-arg "base_image=${run}" -t "${cnb_run}" "${stack_dir}/cnb/run"

  echo "To publish these images:"
  for image in "${build}" "${run}" "${cnb_build}" "${cnb_run}"; do
    echo "  docker push ${image}"
  done
}

main() {
  version=dev
  build_dest=build
  run_dest=run
  stack_name=

  while getopts "v:b:r:s:" o; do
    case "${o}" in
      v)
        version=${OPTARG}
        [[ -n ${version} ]] || usage
        ;;
      b)
        build_dest=${OPTARG}
        [[ -n ${build_dest} ]] || usage
        ;;
      r)
        run_dest=${OPTARG}
        [[ -n ${run_dest} ]] || usage
        ;;
      s)
	stack_name=${OPTARG}
	[[ -n ${stack_name} ]] || usage
	;;
      *)
        usage
        ;;
    esac
  done

  if [[ "${stack_name}" != "full" && "${stack_name}" != "base" ]]; then
    usage
  fi

  build_images "${stack_name}" "${version}" "${build_dest}" "${run_dest}"
}

main "$@"
