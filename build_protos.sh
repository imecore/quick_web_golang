#!/bin/bash

set -ex

cd "$(dirname "$0")"

BUILDER="docker run --rm -v$(pwd):$(pwd) -w$(pwd) thethingsindustries/protoc:3.1.27"

PROTOS_OUTDIR="protos"

$BUILDER -Iprotos -I/usr/include \
  --go_out=plugins=grpc:"$PROTOS_OUTDIR" \
  --grpc-gateway_out=logtostderr=true:"$PROTOS_OUTDIR" \
  protos/*.proto