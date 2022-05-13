#!/usr/bin/env bash

cd "${0%/*}" || exit

for d in *; do
  if [ -d "$d/proto" ]; then
    pushd "$d" > /dev/null || exit
    echo "recompiling $d/proto ..."
    protoc --go_out=. --go-grpc_out=. proto/*
    popd > /dev/null || exit
  fi
done
