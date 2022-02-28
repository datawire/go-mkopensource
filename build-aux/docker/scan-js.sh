#!/bin/bash
set -e
set -o pipefail

. /scripts/imports.sh

validate_required_variable USER_ID

scan_npm_package() {
  echo >&2 "Getting NPM dependencies for $1"
  SOURCE=$1
  pushd $(dirname "${SOURCE}") >/dev/null

  PKG_NAME=$(jq -r '.name + "@" + .version' package.json)
  if [ -z "${PKG_NAME}" ]; then
    echo >&2 "ERROR: Could not get package name"
    return 1
  fi

  echo >&2 "Analyzing package ${PKG_NAME}"
  npm install -f >&2

  PACKAGE_DEPS="/temp/${PKG_NAME}-licenses.json"
  license-checker --excludePackages "${PKG_NAME}" --customPath "/scripts/customLicenseFormat.json" \
    --json >"${PACKAGE_DEPS}"

  /scripts/js-mkopensource --application-type=${APPLICATION_TYPE} < <(cat "${PACKAGE_DEPS}") >"$2"

  popd >/dev/null
}

cd /app

# Get dependencies for each package.json
DEPENDENCIES="js_deps.json"
find . -name package.json -print | while read -r file; do
  scan_npm_package "${file}" "${DEPENDENCIES}"
done

# Generate license information
(
  find "$(pwd)" -name "${DEPENDENCIES}" -print | while read -r file; do
    echo >&2 "Getting licenses for ${file}"
    jq -r '.licenseInfo | to_entries | .[] | "* [" + .key + "](" + .value + ")"' "${file}"
  done
) >"${JS_LICENSES}"

# Generate dependency information
(
  find "$(pwd)" -name "${DEPENDENCIES}" -print | while read -r file; do
    echo >&2 "Getting dependencies for ${file}"
    jq -r '.dependencies[] | .name + "|" + .version + "|" + (.licenses | flatten | join(", "))' "${file}"
  done
) | sort | uniq >/tmp/deps.txt

generate_opensource /tmp/deps.txt Node.Js "${JS_DEPENDENCIES}"
chown "${USER_ID}" -R /temp
