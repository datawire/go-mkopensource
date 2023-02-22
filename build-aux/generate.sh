#!/bin/bash
set -e
set -o pipefail

export DOCKER_BUILDKIT=1

archive_dependencies() {
  tar -vf "$1" -c $2
}

UNPARSABLE_PACKAGE_VALUE=""

if [ $# -gt 0 ]; then
  if [ "$1" = "--unparsable-packages" ]; then
    UNPARSABLE_PACKAGE_VALUE="$2"

  fi
fi


BUILD_SCRIPTS=$(dirname $(realpath "$0"))
. "${BUILD_SCRIPTS}/docker/imports.sh"

validate_required_variable APPLICATION
validate_required_variable APPLICATION_TYPE
validate_required_variable BUILD_HOME
validate_required_variable BUILD_TMP

######################################################################
# Go dependencies
######################################################################
echo "Scanning Go dependency licenses"
validate_required_variable GO_IMAGE
validate_required_variable SCRIPTS_HOME
validate_required_variable GIT_TOKEN

pushd "${BUILD_HOME}" >/dev/null
docker build \
  -f "${BUILD_HOME}/${SCRIPTS_HOME}/build-aux/docker/go_builder.dockerfile" \
  --build-arg APPLICATION_TYPE="${APPLICATION_TYPE}" \
  --build-arg GIT_TOKEN="${GIT_TOKEN}" \
  --build-arg GO_IMAGE="${GO_IMAGE}" \
  --build-arg SCRIPTS_HOME="${SCRIPTS_HOME}" \
  --build-arg UNPARSABLE_PACKAGE="${UNPARSABLE_PACKAGE_VALUE}" \
  -t "go-deps-builder" --target license_output \
  --output "${BUILD_TMP}" .
popd >/dev/null

######################################################################
# Node.Js dependencies
######################################################################
if [ -n "${NPM_PACKAGES}" ]; then
  echo "Scanning Node.Js dependency licenses"
  validate_required_variable NODE_IMAGE

  archive_dependencies "${BUILD_SCRIPTS}/docker/npm_dependencies.tar" "${NPM_PACKAGES}"

  pushd "${BUILD_HOME}" >/dev/null
  docker build \
    -f "${BUILD_HOME}/${SCRIPTS_HOME}/build-aux/docker/js_builder.dockerfile" \
    --build-arg APPLICATION="${APPLICATION}" \
    --build-arg NODE_IMAGE="${NODE_IMAGE}" \
    --build-arg APPLICATION_TYPE="${APPLICATION_TYPE}" \
    --build-arg SCRIPTS_HOME="${SCRIPTS_HOME}" \
    --build-arg EXCLUDED_PKG="${EXCLUDED_PKG}" \
    --build-arg USER_ID="${UID}" \
    -t "js-deps-builder" \
    --target license_output \
    --output "${BUILD_TMP}" .
  popd >/dev/null
fi

# Generate DEPENDENCY_LICENSES.md
(
  echo -e "${APPLICATION} incorporates Free and Open Source software under the following licenses:\n"
  (
    if [ -f "${BUILD_TMP}/go_licenses.txt" ]; then cat "${BUILD_TMP}/go_licenses.txt"; fi
    if [ -f "${BUILD_TMP}/js_licenses.txt" ]; then cat "${BUILD_TMP}/js_licenses.txt"; fi
  ) | sort | uniq | sed -e 's/\[\([^]]*\)]()/\1/'
) >"${BUILD_HOME}/DEPENDENCY_LICENSES.md"

# Generate DEPENDENCIES.md
(
  if [ -f "${BUILD_TMP}/go_dependencies.txt" ]; then
    cat "${BUILD_TMP}/go_dependencies.txt"
    echo -e "\n"
  fi

  if [ -f "${BUILD_TMP}/js_dependencies.txt" ]; then
    cat "${BUILD_TMP}/js_dependencies.txt"
    echo -e "\n"
  fi
) >"${BUILD_HOME}/DEPENDENCIES.md"

# copy go.mod and go.sum
cp  "${BUILD_TMP}/go.mod" "${BUILD_HOME}"
cp  "${BUILD_TMP}/go.sum" "${BUILD_HOME}"
