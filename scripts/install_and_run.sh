#!/usr/bin/env sh -ue

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source $BASE_DIR/utils.sh

function build() {
  [ -z $GOPATH ] && fail "The GOPATH environment is not set"
  go install github.com/svett/nuvi/cmd/nuvi
}

function run() {
  $GOPATH/bin/nuvi "$@"
}

build
run "$@"
