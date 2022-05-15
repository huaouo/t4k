#!/usr/bin/env bash

pushd "${0%/*}" >>/dev/null || exit

for d in *; do
  if [ -d "$d/proto" ]; then
    pushd "$d" >/dev/null || exit
    echo "recompiling $d/proto ..."
    protoc --go_out=. --go-grpc_out=. proto/*
    popd >/dev/null || exit
  fi
done

popd >/dev/null || exit
