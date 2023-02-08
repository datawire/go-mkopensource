#!/usr/bin/env bash
set -e

if ! git diff --name-only --exit-code DEPENDENCIES.md; then
    echo "DIRTY=true" >> $GITHUB_OUTPUT
    exit 0
fi

if ! git diff --name-only --exit-code DEPENDENCY_LICENSES.md; then
    echo "DIRTY=true" >> $GITHUB_OUTPUT
    exit 0
fi

echo "There are no changes to save"
echo "DIRTY=false" >> $GITHUB_OUTPUT
