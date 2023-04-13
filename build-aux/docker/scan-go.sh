#!/bin/bash
set -e
set -o pipefail

. /scripts/imports.sh

BUILD_TMP=/temp
mkdir -p "${BUILD_TMP}"

cd /app

GO_VERSION=$(go version | sed -E 's/.*go([1-9\.]*).*/\1/')

/scripts/go-mkopensource --output-format=txt --package=mod --output-type=markdown --gotar="$(ls /data/go*.src.tar.gz)" \
  --unparsable-packages="${UNPARSABLE_PACKAGE}" >"${GO_DEPENDENCIES}"

DEPENDENCY_INFO="${BUILD_TMP}/go_dependencies.json"
/scripts/go-mkopensource --output-format=txt --package=mod --output-type=json --application-type=${APPLICATION_TYPE} \
  --unparsable-packages="${UNPARSABLE_PACKAGE}" --gotar="$(ls /data/go*.src.tar.gz)" >"${DEPENDENCY_INFO}"

jq -r '.licenseInfo | to_entries | .[] | "* [" + .key + "](" + .value + ")"' "${DEPENDENCY_INFO}" >"${GO_LICENSES}"
