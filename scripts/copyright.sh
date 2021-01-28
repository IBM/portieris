#! /bin/bash

set +x 

GOSRCFILES=($(find . -name "*.go"))
THISYEAR=$(date +"%Y")

for GOFILE in "${GOSRCFILES[@]}"; do
  # If no copyright at all
  if ! grep -q "Licensed under the Apache License, Version 2.0" $GOFILE; then
    echo "// Copyright ${THISYEAR}  Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

$(cat ${GOFILE})" > ${GOFILE}
   fi    
done

# Check for multi-line in new files
GOSRCFILES=($(git diff --name-only origin/master | grep '\.go$'))
for GOFILE in "${GOSRCFILES[@]}"; do
  if grep -q "Copyright .* Portieris Authors." $GOFILE; then
      YEAR_LINE=$(grep "Copyright .* Portieris Authors." $GOFILE)
      YEARS=($(echo $YEAR_LINE | grep -oE '\d{4}'))
      if [[ ${#YEARS[@]} == 1 ]]; then
         if [[ ${YEARS[0]} != ${THISYEAR} ]]; then
            sed -i '' -e "s|Copyright ${YEARS[0]} Portieris Authors.|Copyright ${THISYEAR} Portieris Authors.|" $GOFILE
         fi
      elif [[ ${#YEARS[@]} == 2 ]]; then
        if [[ ${YEARS[1]} != ${THISYEAR} ]]; then
            sed -i '' -e "s|Copyright ${YEARS[0]}-${YEARS[1]} Portieris Authors.|Copyright ${YEARS[0]}, ${THISYEAR} Portieris Authors.|" $GOFILE
        fi
      fi
  fi
done
