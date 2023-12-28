#!/usr/bin/env bash

dirty=""
if ! git diff --quiet; then
  dirty="DIRTY-"
fi

hash=$(git rev-parse --short HEAD)

echo "$dirty$1-$hash"