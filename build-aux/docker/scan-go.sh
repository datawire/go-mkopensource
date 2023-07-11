#!/bin/bash
set -e
set -o pipefail

. /scripts/imports.sh

BUILD_TMP=/temp
mkdir -p "${BUILD_TMP}"

cd /app

GO_VERSION=$(go version | sed -E 's/.*go([1-9\.]*).*/\1/')

ADDITIONAL_GENERATE_ARGS=""
if [[ -n "${UNPARSABLE_PACKAGE}" ]]; then
    ADDITIONAL_GENERATE_ARGS="--unparsable-packages=${UNPARSABLE_PACKAGE} "
fi
if [[ -n "${PROPRIETARY_PACKAGES}" ]]; then
    ADDITIONAL_GENERATE_ARGS="${ADDITIONAL_GENERATE_ARGS} --proprietary-software=${PROPRIETARY_PACKAGES} "
fi

/scripts/go-mkopensource --output-format=txt --package=mod --output-type=markdown --gotar="$(ls /data/go*.src.tar.gz)"  \
    ${ADDITIONAL_GENERATE_ARGS} >"${GO_DEPENDENCIES}"

DEPENDENCY_INFO="${BUILD_TMP}/go_dependencies.json"
/scripts/go-mkopensource --output-format=txt --package=mod --output-type=json --application-type=${APPLICATION_TYPE} \
    ${ADDITIONAL_GENERATE_ARGS} --gotar="$(ls /data/go*.src.tar.gz)" >"${DEPENDENCY_INFO}"

jq -r '.licenseInfo | to_entries | .[] | "* [" + .key + "](" + .value + ")"' "${DEPENDENCY_INFO}" >"${GO_LICENSES}"
