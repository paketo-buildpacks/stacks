#!/usr/bin/env bash

set -euo pipefail

usage() { echo "Usage: $0 [-v <image version>]" 1>&2; exit 1; }

update_bionic_image(){
  docker pull ubuntu:bionic
}

build_build_base_image(){
  local tag=$1

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
  local tag=$1

  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  arch="x86_64"
  tiny_dir=${base_dir}/tiny

  docker build -t "${tag}" "$tiny_dir/dockerfile/run"
}

build_build_cnb_image(){
  local tag=$1
  local base_image=$2
  local description=$3
  local date=$4
  local fully_qualified_base_image=$5
  local package_metadata=$6


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
    --build-arg "package_metadata=${package_metadata}" \
    "$bionic_dir/cnb/build"

}

build_run_cnb_image(){
  local tag=$1
  local base_image=$2
  local description=$3
  local date=$4
  local fully_qualified_base_image=$5
  local package_metadata=$6

  scripts_dir=$(cd "$(dirname $0)" && pwd)
  base_dir=${scripts_dir}/..
  tiny_dir=${base_dir}/tiny

  docker build -t "${tag}" \
    --build-arg "base_image=${base_image}" \
    --build-arg "released=${date}" \
    --build-arg "description=${description}" \
    --build-arg "fully_qualified_base_image=${fully_qualified_base_image}" \
    --build-arg "package_metadata=${package_metadata}" \
    "$tiny_dir/cnb/run"
}

get_mixins() {
  local build_image=$1
  local run_image=$2
  local image_name=$3

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

get_run_package_metadata(){
  local image_name=$1

  container_id="$(docker create "${image_name}" sleep)"
  rm -rf /tmp/tiny-pkgs
  docker cp "${container_id}":/var/lib/dpkg/status.d /tmp/tiny-pkgs
  docker rm "${container_id}" >/dev/null

  pkg_array=()

  for pkg in /tmp/tiny-pkgs/*; do
    if [[ "${pkg}" == "/tmp/tiny-pkgs/status.d" ]]; then
      continue
    fi

    name="$(grep "^Package:" "${pkg}" | cut -d ' ' -f2)"
    version="$(grep "^Version:" "${pkg}" | cut -d ' ' -f2)"
    arch="$(grep "^Architecture:" "${pkg}" | cut -d ' ' -f2)"
    source_package="$(grep "^Source-Package:" "${pkg}" | cut -d ' ' -f2)"
    source_version="$(grep "^Source-Version:" "${pkg}" | cut -d ' ' -f2)"
    source_upstream_version="$(grep "^Source-Upstream-Version:" "${pkg}" | cut -d ' ' -f2)"
    summary="$(grep "^Description:" "${pkg}" | cut -d ' ' -f2-)"

    pkg_array+=("{\"name\":\"${name}\",\"version\":\"${version}\",\"arch\":\"${arch}\",\"summary\":\"${summary}\",\"sourcePackage\":{\"name\":\"${source_package}\",\"version\":\"${source_version}\",\"upstreamVersion\":\"${source_upstream_version}\"}}")
  done

  printf -v joined '%s,' "${pkg_array[@]}"
  echo ["${joined%,}"] | jq -c
}

format_metadata_line(){
    local line=$1
    IFS=';' read -r -a package_array <<< "$line"
    summary="$(echo ${package_array[3]} | sed 's/\"/\\"/g')"
    echo "{\"name\":\"${package_array[0]}\",\"version\":\"${package_array[1]}\",\"arch\":\"${package_array[2]}\",\"summary\":\"$summary\",\"sourcePackage\":{\"name\":\"${package_array[4]}\",\"version\":\"${package_array[5]}\",\"upstreamVersion\":\"${package_array[6]}\"}},"
}

parse_package_list(){
    local package_list=$1
    array_contents="$(
    while IFS='\n' read -r line; do
        format_metadata_line "$line"
    done <<< "$package_list"
    )"
    array_contents="$(echo $array_contents | sed 's/.$//')"
    echo ["$array_contents"]
}

get_build_package_metadata() {
  local image_name=$1

  package_list="$(docker run --rm "${image_name}" dpkg-query -W -f='${binary:Package};${Version};${Architecture};${binary:Summary};${source:Package};${source:Version};${source:Upstream-Version}\n')"

  json_package_list="$(parse_package_list "${package_list}")"

  echo "${json_package_list}" | jq -c
}

publish() {
  while [ $# -gt 0 ]; do
    docker push "$1"
    shift
  done
}

main() {
  local version=dev
  local build_dest=build
  local run_dest=run
  local stack_name=
  local build_base_image=
  local publish=

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

  build_package_metadata="$(get_build_package_metadata "${build_base_tag}")"
  run_package_metadata="$(get_run_package_metadata "${run_base_tag}")"

  build_build_cnb_image "${build_cnb_tag}" "${build_base_tag}" "${build_cnb_description}" "${date}" "" "${build_package_metadata}"
  build_run_cnb_image "${run_cnb_tag}" "${run_base_tag}" "${run_cnb_description}" "${date}" "${fully_qualified_run_base_image}" "${run_package_metadata}"

  if [[ -n "${publish}" ]]; then
    publish "${build_cnb_tag}" "${run_cnb_tag}"
  fi

}

main "$@"
