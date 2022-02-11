#!/bin/bash
set -e
set -o pipefail

BUILD_SCRIPTS=$(dirname $(realpath "$0"))
. "${BUILD_SCRIPTS}/docker/imports.sh"

validate_required_variable BUILD_HOME
validate_required_variable BUILD_TMP
validate_required_variable GO_VERSION

download_go_tar() {
  if [ ! -f "${GO_TAR}" ]; then
    curl -o "${GO_TAR}" --fail -L "https://dl.google.com/go/go${GO_VERSION}.src.tar.gz"
  fi
}

scan_go_package() {
   ${BUILD_SCRIPTS}/go-mkopensource --output-format=txt --package=mod --output-type=markdown --gotar="${GO_TAR}" >"${GO_DEPENDENCIES}"

  DEPENDENCY_INFO="${BUILD_TMP}/go_dependencies.json"
   ${BUILD_SCRIPTS}/go-mkopensource --output-format=txt --package=mod --output-type=json --gotar="${GO_TAR}" > "${DEPENDENCY_INFO}"
  jq -r '.licenseInfo | to_entries | .[] | "* [" + .key + "](" + .value + ")"' "${DEPENDENCY_INFO}" | sort >"${GO_LICENSES}"
}


cd "${BUILD_HOME}"

export GO_TAR="${BUILD_TMP}/go${GO_VERSION}.src.tar.gz"
download_go_tar
scan_go_package

# Restore old state
git checkout -- go.mod go.sum
