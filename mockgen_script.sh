#!/usr/bin/env bash

set -euo pipefail

if ! type "mockgen" > /dev/null; then
  go install github.com/golang/mock/mockgen@v1.6.0
  go mod vendor
fi

rm -rf tests/mock/

for file in `find . -name '*.go' | grep -v proto | grep -v vendor | grep -v tests`; do
    if `grep -q 'interface {' ${file}`; then
        dest=${file//internal\//}
        mockgen -source=${file} -destination=tests/mock/${dest}
    fi
done