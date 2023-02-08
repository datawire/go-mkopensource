#!/usr/bin/env bash
set -e

git config --local user.name "d6eautomaton"
git config --local user.email "<>"

DESTINATION_BRANCH="${GITHUB_HEAD_REF:-$GITHUB_REF_NAME}"

git checkout "${DESTINATION_BRANCH}"

echo '::notice:: Committing dependabot changes to DEPENDENCIES.md and/or DEPENDENCY_LICENSES.md'
git commit -m  "Updated dependency information after dependabot change." DEPENDENCIES.md DEPENDENCY_LICENSES.md go.mod go.sum

git push
