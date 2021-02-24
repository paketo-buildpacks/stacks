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
  description=$3
  mixins=$4
  image_name=$5
  date=$6
  fully_qualified_base_image=$7
  package_metadata=$8

  base_dir=$(cd "$(dirname "$0")" && cd .. && pwd)
  stack_dir="${base_dir}/bionic"

  echo "Building cnb ${image_name} image"

  if [[ "${package_metadata}" != "" ]]; then
    docker build -t "${tag}" \
      --build-arg "base_image=${base_image}" \
      --build-arg "description=${description}" \
      --build-arg "mixins=${mixins}" \
      --build-arg "released=${date}" \
      --build-arg "fully_qualified_base_image=${fully_qualified_base_image}" \
      --build-arg "package_metadata=${package_metadata}" \
      "${stack_dir}/cnb/${image_name}"
  else
    grep -v "io.paketo.stack.packages" "${stack_dir}/cnb/${image_name}/Dockerfile" | \
      docker build -t "${tag}" \
      --build-arg "base_image=${base_image}" \
      --build-arg "description=${description}" \
      --build-arg "mixins=${mixins}" \
      --build-arg "released=${date}" \
      --build-arg "fully_qualified_base_image=${fully_qualified_base_image}" \
      -
  fi
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

format_metadata_line(){
    line=$1
    IFS=';' read -r -a package_array <<< "$line"
    summary="$(echo ${package_array[3]} | sed 's/\"/\\"/g')"
    echo "{\"name\":\"${package_array[0]}\",\"version\":\"${package_array[1]}\",\"arch\":\"${package_array[2]}\",\"summary\":\"$summary\",\"source\":{\"name\":\"${package_array[4]}\",\"version\":\"${package_array[5]}\",\"upstream-version\":\"${package_array[6]}\"}},"
}

parse_package_list(){
    package_list=$1
    array_contents="$(
    while IFS='\n' read -r line; do
        format_metadata_line "$line"
    done <<< "$package_list"
    )"
    array_contents="$(echo $array_contents | sed 's/.$//')"
    echo ["$array_contents"]
}

get_package_metadata() {
  image_name=$1

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

  build_description=
  run_description=

  if [[ "${stack_name}" == "base" ]]; then
    build_description="ubuntu:bionic + openssl + CA certs + compilers + shell utilities"
    run_description="ubuntu:bionic + openssl + CA certs"
  elif [[ "${stack_name}" == "full" ]]; then
    build_description="ubuntu:bionic + many common C libraries and utilities"
    run_description="ubuntu:bionic + many common C libraries and utilities"
  else
    usage
  fi

  build_base_tag="${build_dest}:${version}-${stack_name}"
  run_base_tag="${run_dest}:${version}-${stack_name}"
  cnb_build_tag="${build_base_tag}-cnb"
  cnb_run_tag="${run_base_tag}-cnb"

  update_bionic_image

  build_base_image "${build_base_tag}" "${stack_name}" "build"
  build_base_image "${run_base_tag}" "${stack_name}" "run"

  fully_qualified_build_base_image=""
  fully_qualified_run_base_image=""

  if [[ -n "${publish}" ]]; then
    publish "${build_base_tag}" "${run_base_tag}"
    fully_qualified_build_base_image="{\"base-image\":\"$(docker inspect --format='{{index .RepoDigests 0}}' "${build_base_tag}")\"}"
    fully_qualified_run_base_image="{\"base-image\":\"$(docker inspect --format='{{index .RepoDigests 0}}' "${run_base_tag}")\"}"
  fi

  build_mixins="$(get_mixins "${build_base_tag}" "${run_base_tag}" "build")"
  run_mixins="$(get_mixins "${build_base_tag}" "${run_base_tag}" "run")"

  build_package_metadata=
  run_package_metadata=
  if [[ "${stack_name}" == "base" ]]; then
    build_package_metadata="$(get_package_metadata "${build_base_tag}")"
    run_package_metadata="$(get_package_metadata "${run_base_tag}")"
  fi

  date="$(date '+%Y-%m-%d')"

  build_cnb_image "${cnb_build_tag}" "${build_base_tag}" "${build_description}" "${build_mixins}" "build" "${date}" "${fully_qualified_build_base_image}" "${build_package_metadata}"
  build_cnb_image "${cnb_run_tag}" "${run_base_tag}" "${run_description}" "${run_mixins}" "run" "${date}" "${fully_qualified_run_base_image}" "${run_package_metadata}"

  if [[ -n "${publish}" ]]; then
    publish "${cnb_build_tag}" "${cnb_run_tag}"
  fi
}

main "$@"
