#! /bin/bash

CRC=0

RED='\033[0;31m'
NC='\033[0m' # No Color

function fail()
{
    echo -en $RED
    tput bold
    echo -e "$1"
    tput sgr0
    echo -en $NC
}


# Enforce check for copyright statements in Go code
GOSRCFILES=($(find . -name "*.go" | grep -v code-generator))
THISYEAR=$(date +"%Y")
for GOFILE in "${GOSRCFILES[@]}"; do
  if ! grep -q "Licensed under the Apache License, Version 2.0" $GOFILE; then
    fail "Missing copyright/licence statement in ${GOFILE}."
    CRC=$(($CRC + 1))
  fi
done 

# Check for multi-line in new files
if git fetch origin ; then
    GOSRCFILES=($(git diff --name-only origin/main | grep '\.go$'))
    for GOFILE in "${GOSRCFILES[@]}"; do
    if grep -q "Copyright .* Portieris Authors." $GOFILE; then
        YEAR_LINE=$(grep "Copyright .* Portieris Authors." $GOFILE)
        YEARS=($(echo $YEAR_LINE | grep -oE '[0-9]{4}'))
        if [[ ${#YEARS[@]} == 1 ]]; then
            if [[ ${YEARS[0]} != ${THISYEAR} ]]; then
                fail "Single out-of-date copyright in ${GOFILE}."
                CRC=$(($CRC + 1))
            fi
        elif [[ ${#YEARS[@]} == 2 ]]; then
            if [[ ${YEARS[1]} != ${THISYEAR} ]]; then
                fail "Double year copyright with out-of-date second year in ${GOFILE}."
                CRC=$(($CRC + 1))
            fi
        else
            echo "#YEARS was ${#YEARS[@]}"
        fi
    fi
    done
    if [ $CRC -gt 0 ]; then fail "Please run make copyright to add copyright statements and check in the updated file(s).\n"; fi
else
    fail "Failed to get branch information for origin/main, cannot perform copyright check. Make sure that you have a remote called origin in your git project."
    CRC=$(($CRC + 1))
fi
exit $CRC
