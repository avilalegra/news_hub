#!/usr/bin/env bash

PKG_LIST=$(go list ./... | grep -v /vendor/ | tr '\n' ' ')
GOFLAGS="-count=1" go test -covermode=count -coverprofile coverage $PKG_LIST  -p 1 ./...

if [[ "$1" == "--html" ]]; then
  go tool cover -html=coverage
else
  go tool cover -func=coverage
fi