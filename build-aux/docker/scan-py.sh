#!/bin/bash
set -e
set -o pipefail

. /scripts/imports.sh

scan_python_package() {
  echo >&2 "Getting python dependencies for $1"
  SOURCE=$1
  pushd $(dirname "${SOURCE}") > /dev/null

  python3 -m venv ./venv
  . "./venv/bin/activate"
  pip3 install -U pip==20.2.4 pip-tools==5.3.1
  pip3 --disable-pip-version-check install -r requirements.txt

  {
    pip3 --disable-pip-version-check freeze --exclude-editable | cut -d= -f1 | xargs pip show
    echo ''
  } | sed 's/^---$//' >"$2"

  popd > /dev/null
}

cd /app

# Get dependencies for each requirements.txt
DEPENDENCIES="py_deps.txt"
find . -name requirements.txt -print | while read -r file; do
  scan_python_package "${file}" "${DEPENDENCIES}"
done

# Get dependencies
JSON_DEPS="/temp/py_dependencies.json"
find . -name "${DEPENDENCIES}" -exec cat '{}' \; | /scripts/py-mkopensource --output-type=json \
  --application-type="${APPLICATION_TYPE}" > "${JSON_DEPS}"

# Generate license information
jq -r '.licenseInfo | to_entries | .[] | "* [" + .key + "](" + .value + ")"' "${JSON_DEPS}" >"${PY_LICENSES}"

# Generate dependency information
jq -r '.dependencies[] | .name + "|" + .version + "|" + (.licenses | flatten | join(", "))' "${JSON_DEPS}" > /tmp/deps.txt
generate_opensource /tmp/deps.txt Python "${PY_DEPENDENCIES}"
