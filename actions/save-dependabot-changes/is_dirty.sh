#!/usr/bin/env bash
set -e

if ! git diff --name-only --exit-code DEPENDENCIES.md; then
    echo '::notice:: DEPENDENCIES.md and/or DEPENDENCY_LICENSES.md changed and they will be committed.'
    echo "DIRTY=true" >> $GITHUB_OUTPUT
    exit 0
fi

if ! git diff --name-only --exit-code DEPENDENCY_LICENSES.md; then
    echo '::notice:: DEPENDENCIES.md and/or DEPENDENCY_LICENSES.md changed and they will be committed.'
    echo "DIRTY=true" >> $GITHUB_OUTPUT
    exit 0
fi

echo "::debug:: There are no changes to save"
echo "DIRTY=false" >> $GITHUB_OUTPUT
