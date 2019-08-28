#!/usr/bin/env bats

setup() {
    docker build "$BATS_TEST_DIRNAME/testapp" --tag tiny-testapp
}

@test "the testapp runs" {
    docker run --rm tiny-testapp
}

teardown() {
    docker rmi --force tiny-testapp
}