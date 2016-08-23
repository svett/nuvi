#!/usr/bin/env sh -ue

function main() {
  fetch_submodules
  check_if_redis_exists
  fetch_golang_dependencies
  run_tests
}

function fetch_submodules() {
  echo "Fetching submodules"
  git submodule update --init --recursive
}

function check_if_redis_exists() {
 if ! [ -x "$(command -v redis-server)" ]; then
   fail "Redis server is not installed"
 fi
}

function fetch_golang_dependencies() {
  [ -z $GOPATH ] && fail "The GOPATH environment is not set"
  echo "Getting Golang dependencies"
  go get github.com/onsi/ginkgo/ginkgo
  go get github.com/onsi/gomega
}

function run_tests() {
  echo "Running tests"
  ginkgo .
  ginkgo integration/
}

fail() {
  echo "$*"
  exit 1
}

main
