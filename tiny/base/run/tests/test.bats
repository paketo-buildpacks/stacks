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
@test "the os-release file some contents" {
    local TMP_OSRELEASE=$(mktemp -d)
    docker cp "${container_id}:/etc/os-release" "${TMP_OSRELEASE}"

    [[ "$(grep -v 'VERSION=' "${TMP_OSRELEASE}/os-release" )" == "$(cat "./files/os-release")" && $(grep 'VERSION=' "${TMP_OSRELEASE}/os-release") ]]
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

