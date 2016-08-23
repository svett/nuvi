#!/usr/bin/env sh -ue

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source $BASE_DIR/utils.sh

function main() {
  pushd $BASE_DIR/..
  fetch_submodules
  check_if_redis_exists
  fetch_golang_dependencies
  run_tests $PWD
  popd
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
  run_ginkgo -r -skip vendor $1
}

function run_ginkgo() {
 $GOPATH/bin/ginkgo $1
}

main
