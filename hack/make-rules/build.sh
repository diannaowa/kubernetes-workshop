#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT=$(dirname "${BASH_SOURCE[0]}")/../..
echo $PROJECT_ROOT
if [ -z "$@" ];then
  for i in $(ls $PROJECT_ROOT/cmd);do
    go build -o ${PROJECT_ROOT}/bin/cmd/$i ${PROJECT_ROOT}/cmd/${i}/main.go
  done
  exit 0
fi
go build -o ${PROJECT_ROOT}/bin/$@ ${PROJECT_ROOT}/$@/main.go
