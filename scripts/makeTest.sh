#!/bin/bash

GOTAGS=$2

failures=0
echo 'mode: atomic' > cover.out
for PACKAGE in $1; do
    go test --tags $GOTAGS -covermode=atomic -coverprofile=cover.tmp $PACKAGE && tail -n +2 cover.tmp >> cover.out
    packageResult=$?
    failures=$((failures + $packageResult))

    if [[ $packageResult > 0 ]]; then
        echo "At least one test failed in package $PACKAGE"
    fi

    rm -f cover.tmp
done

if [[ $failures > 0 ]]; then
    echo "$failures package has failed tests"
    exit 1
else
    echo "all tests pass"
fi
