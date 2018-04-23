#! /bin/bash

GOSRCFILES=($(find . -name "*.go" | grep -v vendor))

SUCCESS=0
for GOFILE in "${GOSRCFILES[@]}"; do
  # If no copyright at all
  if ! grep -q "Licensed under the Apache License, Version 2.0" $GOFILE; then
    echo "$GOFILE is missing a copyright/license"
    SUCCESS=1
  fi 
done

if [[ $SUCCESS -ne 0 ]]; then
  exit 1
fi   

