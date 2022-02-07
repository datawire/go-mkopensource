#!/bin/bash
set -e
set -o pipefail

export DOCKER_BUILDKIT=1

archive_dependencies() {
  tar -vf "$1" -c $2
}

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
  -t "go-deps-builder" --target license_output \
  --output "${BUILD_TMP}" .
popd >/dev/null

######################################################################
# Python dependencies
######################################################################
if [ -n "${PYTHON_PACKAGES}" ]; then
  echo "Scanning Python dependency licenses"
  validate_required_variable PYTHON_IMAGE

  archive_dependencies "${BUILD_SCRIPTS}/docker/python_dependencies.tar" "${PYTHON_PACKAGES}"

  pushd "${BUILD_HOME}" >/dev/null
  docker build \
    -f "${BUILD_HOME}/${SCRIPTS_HOME}/build-aux/docker/py_builder.dockerfile" \
    --build-arg PYTHON_IMAGE="${PYTHON_IMAGE}" \
    --build-arg SCRIPTS_HOME="${SCRIPTS_HOME}" \
    -t "py-deps-builder" \
    --target python_dependency_scanner .
  popd >/dev/null

  docker run --rm --env APPLICATION \
    --volume "$(realpath ${BUILD_TMP})":/temp \
    py-deps-builder /scripts/scan-py.sh
fi

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
    --build-arg NODE_IMAGE="${NODE_IMAGE}" \
    --build-arg APPLICATION_TYPE="${APPLICATION_TYPE}" \
    --build-arg SCRIPTS_HOME="${SCRIPTS_HOME}" \
    -t "js-deps-builder" \
    --target npm_dependency_scanner .
  popd >/dev/null

  docker run --rm \
    --env APPLICATION \
    --env USER_ID=${UID} \
    --volume "$(realpath ${BUILD_TMP})":/temp \
    js-deps-builder /scripts/scan-js.sh
fi

# Generate DEPENDENCY_LICENSES.md
(
  echo -e "${APPLICATION} incorporates Free and Open Source software under the following licenses:\n"
  (
    if [ -f "${BUILD_TMP}/go_licenses.txt" ]; then cat "${BUILD_TMP}/go_licenses.txt"; fi
    if [ -f "${BUILD_TMP}/py_licenses.txt" ]; then cat "${BUILD_TMP}/py_licenses.txt"; fi
    if [ -f "${BUILD_TMP}/js_licenses.txt" ]; then cat "${BUILD_TMP}/js_licenses.txt"; fi
  ) | sort | uniq | sed -e 's/\[\([^]]*\)]()/\1/'
) >"${BUILD_HOME}/DEPENDENCY_LICENSES.md"

# Generate DEPENDENCIES.md
(
  if [ -f "${BUILD_TMP}/go_dependencies.txt" ]; then
    cat "${BUILD_TMP}/go_dependencies.txt"
    echo -e "\n"
  fi

  if [ -f "${BUILD_TMP}/py_dependencies.txt" ]; then
    cat "${BUILD_TMP}/py_dependencies.txt"
    echo -e "\n"
  fi

  if [ -f "${BUILD_TMP}/js_dependencies.txt" ]; then
    cat "${BUILD_TMP}/js_dependencies.txt"
    echo -e "\n"
  fi
) >"${BUILD_HOME}/DEPENDENCIES.md"
