#!/bin/sh

set -eo pipefail

# Process envvars.
if [ -z "$PLUGIN_PATTERN" ]; then
  echo '$PLUGIN_PATTERN must be set.' >&2
  exit 1
fi

if [ "$PLUGIN_EXCLUDE_DRONE_CONFIG" != true ]; then
  PLUGIN_PATTERN="($PLUGIN_PATTERN|^.drone.yml$)"
fi

# Check number of changed files.
if [ $(git rev-list --count HEAD) -eq 1 ]; then
  exit 0 # no files changed
fi

# Check if pattern matches any files in the prev-commit-diff.
if ! git diff --name-only HEAD~1 HEAD | \
  egrep "$PLUGIN_PATTERN" > /dev/null; then
  exit 78
fi
