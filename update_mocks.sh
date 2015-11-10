#! /bin/bash -e

mkdir -p \
  handlers/mock

mockgen github.com/nanobox-io/golang-discovery Discover > handlers/mock/mock.go
