#!/usr/bin/env sh -ue

function build() {
  go install github.com/svett/nuvi/cmd/nuvi
}

function run() {
  [ -z $GOPATH ] && fail "The GOPATH environment is not set"
  $GOPATH/bin/nuvi "$@"
}

build
run "$@"
