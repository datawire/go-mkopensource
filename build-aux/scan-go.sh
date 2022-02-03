#!/bin/bash
set -e
set -o pipefail

BUILD_SCRIPTS=$(dirname $(realpath "$0"))
. "${BUILD_SCRIPTS}/docker/imports.sh"

validate_required_variable BUILD_HOME
validate_required_variable BUILD_TMP

check_go_tar() {
  validate_required_variable GO_TAR

  if [ ! -f "${GO_TAR}" ]; then
    echo "Go tar '${GO_TAR}' does not exist" >&2
    exit 1
  fi
}

scan_go_package() {
  echo >&2 "Getting GO dependencies"

   ${BUILD_SCRIPTS}/go-mkopensource --output-format=txt --package=mod --output-type=markdown --gotar="${GO_TAR}" >"${GO_DEPENDENCIES}"

  DEPENDENCY_INFO="${BUILD_TMP}/go_dependencies.json"
   ${BUILD_SCRIPTS}/go-mkopensource --output-format=txt --package=mod --output-type=json --gotar="${GO_TAR}" > "${DEPENDENCY_INFO}"
  jq -r '.licenseInfo | to_entries | .[] | "* [" + .key + "](" + .value + ")"' "${DEPENDENCY_INFO}" | sort >"${GO_LICENSES}"
}


cd "${BUILD_HOME}"

check_go_tar
scan_go_package

# Restore old state
git checkout -- go.mod go.sum
