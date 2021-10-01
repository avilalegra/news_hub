#!/usr/bin/env bash

PKG_LIST=$(go list ./... | grep -v /vendor/ | tr '\n' ' ')
go test -covermode=count -coverprofile coverage $PKG_LIST 

if [[ "$1" == "--html" ]]; then
  go tool cover -html=coverage
else
  go tool cover -func=coverage
fi