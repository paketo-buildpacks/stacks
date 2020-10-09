#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat >&2 <<EOF
Usage: $0 -s <base/full> [<options>]

  -b, --build-dest <dest>       Destination to tag and publish base image to
  -r, --run-dest <dest>         Destination to tag and publish run image to
  -v, --version <image version> Version to include in image tags
  -p, --publish                 Publish images after creating
EOF
  exit 1
}

update_bionic_image() {
  docker pull ubuntu:bionic
}

build_base_image() {
  tag=$1
  stack_name=$2
  image_name=$3

  base_dir=$(cd "$(dirname "$0")" && cd .. && pwd)
  stack_dir="${base_dir}/bionic"
  arch="x86_64"

  echo "Building base ${image_name} image"

  docker build -t "${tag}" \
    --build-arg "sources=$(cat "${base_dir}/arch/${arch}/sources.list")" \
    --build-arg "packages=$(cat "${base_dir}/packages/${stack_name}/${image_name}")" \
    --no-cache "${stack_dir}/dockerfile/${image_name}"

  echo "${tag}"
}

build_cnb_image() {
  tag=$1
  base_image=$2
  mixins=$3
  image_name=$4

  base_dir=$(cd "$(dirname "$0")" && cd .. && pwd)
  stack_dir="${base_dir}/bionic"

  echo "Building cnb ${image_name} image"

  docker build -t "${tag}" \
    --build-arg "base_image=${base_image}" \
    --build-arg "mixins=${mixins}" \
    "${stack_dir}/cnb/${image_name}"
}

get_mixins() {
  build_image=$1
  run_image=$2
  image_name=$3

  build_packages="$(docker run --rm "${build_image}" dpkg-query -f '${Package}\n' -W)"
  run_packages="$(docker run --rm "${run_image}" dpkg-query -f '${Package}\n' -W)"

  shared_mixins="$(comm -12 <(echo "${build_packages}") <(echo "${run_packages}") | jq -cnR '[inputs | select(length>0)]')"

  if [[ "${image_name}" == "build" ]]; then
    image_only_mixins="$(comm -23 <(echo "${build_packages}") <(echo "${run_packages}") | sed 's/^/build:/g' | jq -cnR '[inputs | select(length>0)]')"
  elif [[ "${image_name}" == "run" ]]; then
    image_only_mixins="$(comm -13 <(echo "${build_packages}") <(echo "${run_packages}") | sed 's/^/run:/g' | jq -cnR '[inputs | select(length>0)]')"
  fi

  echo "${shared_mixins}" "${image_only_mixins}" | jq -c -s add
}

publish() {
  while [ $# -gt 0 ]; do
    docker push "$1"
  done
}

main() {
  version=dev
  build_dest=build
  run_dest=run
  stack_name=
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
      -s|--stack)
	stack_name=$2
	[[ -n ${stack_name} ]] || usage
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

  if [[ "${stack_name}" != "full" && "${stack_name}" != "base" ]]; then
    usage
  fi

  base_build_tag="${build_dest}:${version}-${stack_name}"
  base_run_tag="${run_dest}:${version}-${stack_name}"
  cnb_build_tag="${base_build_tag}-cnb"
  cnb_run_tag="${base_run_tag}-cnb"

  update_bionic_image

  build_base_image "${base_build_tag}" "${stack_name}" "build"
  build_base_image "${base_run_tag}" "${stack_name}" "run"

  build_mixins="$(get_mixins "${base_build_tag}" "${base_run_tag}" "build")"
  run_mixins="$(get_mixins "${base_build_tag}" "${base_run_tag}" "run")"

  build_cnb_image "${cnb_build_tag}" "${base_build_tag}" "${build_mixins}" "build"
  build_cnb_image "${cnb_run_tag}" "${base_run_tag}" "${run_mixins}" "run"

  if [[ -n "${publish}" ]]; then
    publish "${base_build_tag}" "${base_run_tag}" "${cnb_build_tag}" "${cnb_run_tag}"
  fi
}

main "$@"
