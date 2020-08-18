#!/usr/bin/env bats

setup() {
    #docker create requires an entry point. As we're copying files and not running, we provide a superficial entrypoint of foo
    container_id="$(docker create tiny foo)"
}

@test "the ca-certificates copyright file exists" {
    check_file_exists /usr/share/doc/ca-certificates/copyright
}

@test "the ca-certificates file exists" {
    check_file_exists /etc/ssl/certs/ca-certificates.crt
}

@test "the ca-certificates directory has been removed" {
    run docker cp ${container_id}:/usr/share/ca-certificates -
    [ "$status" -eq 1 ]
}

@test "the ca-certificates dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/ca-certificates"
}
@test "the base-files dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/base-files"
}
@test "the libssl1.1 dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/libssl1.1"
}
@test "the libc6 dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/libc6"
}
@test "the tzdata dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/tzdata"
}
@test "the netbase dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/netbase"
}
@test "the openssl dpkg status.d exists" {
    check_file_exists "/var/lib/dpkg/status.d/openssl"
}

@test "the /root directory exists" {
    check_file_exists "/root"
}
@test "the /home/nonroot directory exists" {
    check_file_exists "/home/nonroot"
}
@test "the /tmp directory exists" {
    check_file_exists "/tmp"
}

@test "the /etc/services file exists" {
    check_file_exists "/etc/services"
}
@test "the /etc/nsswitch.conf file exists" {
    check_file_exists "/etc/nsswitch.conf"
}

@test "the /etc/passwd file exists" {
    check_file_exists "/etc/passwd"
}
@test "the passwd file some contents" {
    local TMP_PASSWD=$(mktemp -d)
    docker cp "${container_id}:/etc/passwd" "${TMP_PASSWD}"

    [[ "$(cat "${TMP_PASSWD}/passwd")" == "$(cat "./files/passwd")" ]]
}

@test "the /etc/os-release file exists" {
    check_file_exists "/etc/os-release"
}

function remove_tiny_values {
    grep -v 'PRETTY_NAME=' $1 \
      | grep -v 'HOME_URL=' \
      | grep -v 'SUPPORT_URL=' \
      | grep -v 'BUG_REPORT_URL=' \
      | sort
}
@test "the os-release file some contents" {
    local TMP_OSRELEASE=$(mktemp -d)
    local TMP_ORIGINAL_OSRELEASE=$(mktemp -d)

    pushd $TMP_ORIGINAL_OSRELEASE
    apt-get update
    apt download base-files
    BASE_FILES_DEB=$(ls)
    ar p $BASE_FILES_DEB data.tar.xz | unxz | tar x -C .
    popd

    docker cp "${container_id}:/etc/os-release" "${TMP_OSRELEASE}"

    [[     $(grep 'PRETTY_NAME=' "${TMP_OSRELEASE}/os-release")  \
        && $(grep 'HOME_URL=' "${TMP_OSRELEASE}/os-release") \
        && $(grep 'SUPPORT_URL=' "${TMP_OSRELEASE}/os-release") \
        && $(grep 'BUG_REPORT_URL=' "${TMP_OSRELEASE}/os-release") \
        && "$( remove_tiny_values "${TMP_ORIGINAL_OSRELEASE}/etc/os-release" )" == "$( remove_tiny_values "${TMP_OSRELEASE}/os-release" )" ]]

}
@test "the /etc/group file exists" {
    check_file_exists "/etc/group"
}
@test "the groups file some contents" {
    local TMP_GROUP=$(mktemp -d)
    docker cp "${container_id}:/etc/group" "${TMP_GROUP}"

    [[ "$(cat "${TMP_GROUP}/group")" == "$(cat "./files/group")" ]]
}

teardown() {
    docker rm -v ${container_id}
}

check_file_exists() {
    local path=$1

    docker cp ${container_id}:${path} - > /dev/null
}

