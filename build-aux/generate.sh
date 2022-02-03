#!/bin/env bash
set -ex
set -o pipefail

archive_dependencies() {
  tar -vf "$1" -c $2
}

#TODO: Remove DOCKER_BUILDKIT
DOCKER_BUILDKIT=0

BUILD_SCRIPTS=$(dirname $(realpath "$0"))
. "${BUILD_SCRIPTS}/docker/imports.sh"

validate_required_variable APPLICATION
validate_required_variable BUILD_HOME
validate_required_variable BUILD_TMP

pushd "${BUILD_HOME}" >/dev/null
#TODO: Enable if statement
if [ -f go.mod ]; then
  validate_required_variable GO_VERSION

  pushd "${BUILD_SCRIPTS}/../cmd/go-mkopensource" >/dev/null
  go build -o "${BUILD_SCRIPTS}/" .
  popd >/dev/null

  "${BUILD_SCRIPTS}/scan-go.sh"
fi

if [ -n "${PYTHON_PACKAGES}" ]; then
  validate_required_variable PYTHON_VERSION

  archive_dependencies "${BUILD_SCRIPTS}/docker/python_dependencies.tar" "${PYTHON_PACKAGES}"

  pushd "${BUILD_SCRIPTS}/../cmd/py-mkopensource" >/dev/null
  go build -o "${BUILD_SCRIPTS}/docker" .
  popd >/dev/null

  pushd "${BUILD_SCRIPTS}/docker" >/dev/null
  docker build -f py_builder.dockerfile --build-arg PYTHON_VERSION="${PYTHON_VERSION}" -t "py-deps-builder" --target python_dependency_scanner .
  popd >/dev/null

  docker run --rm --env APPLICATION \
    --volume "${BUILD_TMP}":/temp \
    py-deps-builder /scripts/scan-py.sh ;\
fi

# Generate LICENSES.md
(
  echo -e "${APPLICATION} incorporates Free and Open Source software under the following licenses:\n"
  (
    if [ -f "${BUILD_TMP}/go_licenses.txt" ]; then cat "${BUILD_TMP}/go_licenses.txt"; fi
    if [ -f "${BUILD_TMP}/py_licenses.txt" ]; then cat "${BUILD_TMP}/py_licenses.txt"; fi
    if [ -f "${BUILD_TMP}/js_licenses.txt" ]; then cat "${BUILD_TMP}/js_licenses.txt"; fi
  ) | sort | uniq
) >"${BUILD_HOME}/LICENSES.md"

# Generate OPENSOURCE.md
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
) >"${BUILD_HOME}/OPENSOURCE.md"

exit 0

validate_required_variable "PYTHON_VERSION"
validate_required_variable "NODE_VERSION"

validate_required_variable "JS_MKOPENSOURCE"
validate_required_variable "PY_MKOPENSOURCE"
validate_required_variable "GO_MKOPENSOURCE"

validate_required_variable "APPLICATION"

SCRIPT_DIR="$(dirname $(realpath $0))"
cd "${HOME}"

#Build images to scan dependencies
pushd "${SCRIPT_DIR}/docker" >/dev/null
cp -a "${JS_MKOPENSOURCE}" "${PY_MKOPENSOURCE}" "${GO_MKOPENSOURCE}" .
docker build -f js_builder.dockerfile --build-arg NODE_VERSION="${NODE_VERSION}" -t "js-deps-builder" --target npm_dependency_scanner .
popd >/dev/null
