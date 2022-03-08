#!/bin/bash

set -e

WorkPath=$(pwd)

# export tools bin to PATH
PATH="${WorkPath}"/tools:"${PATH}"

if [ ! -f "$WorkPath"/tools/goimports-reviser-darwin ]; then
  {
    mkdir -p tools
    wget https://resource.gocloudcoder.com/goimports-reviser-darwin -o "${WorkPath}"/tools
  }
fi

for i in $(git diff --cached --name-only --diff-filter=ACM -- '*.go' | grep -v ".pb.go" | grep -v ".pb.gw.go" | grep -v "_test.go") ; do
    "$WorkPath"/tools/goimports-reviser-darwin -rm-unused -set-alias -format -file-path "$i"
done
