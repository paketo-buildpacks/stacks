#!/usr/bin/env bash

set -euo pipefail

usage() { echo "Usage: $0 [-v <image version>]" 1>&2; exit 1; }

update_bionic_image(){
  docker pull ubuntu:bionic
}


build_build_base_image(){
  tag=$1

  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  arch="x86_64"
  bionic_dir=${base_dir}/bionic

  docker build -t "${tag}" \
    --build-arg "sources=$(cat "${base_dir}/arch/${arch}/sources.list")" \
    --build-arg "packages=$(cat "${base_dir}/packages/base/build")" \
    --no-cache "${bionic_dir}/dockerfile/build"
}

build_run_base_image(){
  tag=$1

  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  arch="x86_64"
  tiny_dir=${base_dir}/tiny

  docker build -t "${tag}" "$tiny_dir/dockerfile/run"
}

build_build_cnb_image(){
  tag=$1
  base_image=$2
  description=$3
  date=$4
  fully_qualified_base_image=$5


  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  bionic_dir=${base_dir}/bionic

  # Build CNB images
  docker build -t "${tag}" \
    --build-arg "base_image=${base_image}" \
    --build-arg "stack_id=io.paketo.stacks.tiny" \
    --build-arg "released=${date}" \
    --build-arg "description=${description}" \
    --build-arg "fully_qualified_base_image=${fully_qualified_base_image}" \
    --build-arg "mixins=[]" \
    "$bionic_dir/cnb/build"

}

build_run_cnb_image(){
  tag=$1
  base_image=$2
  description=$3
  date=$4
  fully_qualified_base_image=$5

  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  tiny_dir=${base_dir}/tiny

  docker build -t "${tag}" \
    --build-arg "base_image=${base_image}" \
    --build-arg "released=${date}" \
    --build-arg "description=${description}" \
    --build-arg "fully_qualified_base_image=${fully_qualified_base_image}" \
    "$tiny_dir/cnb/run"
}

publish() {
  while [ $# -gt 0 ]; do
    docker push "$1"
    shift
  done
}

main() {
  version=dev
  build_dest=build
  run_dest=run
  stack_name=
  build_base_image=
  publish=

  while [ $# -gt 0 ]; do
    case $1 in
      -v|--version)
        version=$2
        [[ -n ${version} ]] || usage
        shift
        ;;
      -b|--build-dest)
        build_dest=$2
        [[ -n ${build_dest} ]] || usage
        shift
        ;;
      -r|--run-dest)
        run_dest=$2
        [[ -n ${run_dest} ]] || usage
        shift
        ;;
      -i|--base-build-image)
	build_base_image=$2
        shift
	;;
      -p|--publish)
        publish="true"
        ;;
      --) shift; break;;
      -*) usage;;
      *) break;;
    esac
    shift
  done

  build_base_tag="${build_base_image}"
  run_base_tag="${run_dest}:${version}-tiny"
  build_cnb_tag="${build_dest}:${version}-tiny-cnb"
  run_cnb_tag="${run_base_tag}-cnb"

  update_bionic_image

  if [[ -z "$build_base_image" ]]; then
   build_base_tag="${build_dest}:${version}-tiny"
   build_build_base_image "${build_base_tag}"
  fi

  build_run_base_image "${run_base_tag}"

  fully_qualified_run_base_image=""

  if [[ -n "${publish}" ]]; then
    publish "${run_base_tag}"
    fully_qualified_run_base_image="{\"base-image\":\"$(docker inspect --format='{{index .RepoDigests 0}}' "${run_base_tag}")\"}"
  fi

  date=$(date '+%Y-%m-%d')
  build_cnb_description="ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
  run_cnb_description="distroless-like bionic + glibc + openssl + CA certs"

  build_build_cnb_image "${build_cnb_tag}" "${build_base_tag}" "${build_cnb_description}" "${date}" ""
  build_run_cnb_image "${run_cnb_tag}" "${run_base_tag}" "${run_cnb_description}" "${date}" "${fully_qualified_run_base_image}"

  if [[ -n "${publish}" ]]; then
    publish "${build_cnb_tag}" "${run_cnb_tag}"
  fi

}

main "$@"
