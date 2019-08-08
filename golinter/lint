#!/bin/sh

# Strict mode.
set -euo pipefail

# Run 'goimports'.
echo "Formatting with 'goimports'..." >&2
goimports -w -l $(find . -name '*.go' | grep -v '.pb.go') \
  | tee /dev/fd/1 \
  | xargs -0 test -z

# Run 'revive'.
echo "Linting with 'revive'..." >&2
if [ -n "$PLUGIN_REVIVE_CONFIG" ]; then
  REVIVE_CONFIG="$PLUGIN_REVIVE_CONFIG"
fi
if [ -n "$REVIVE_CONFIG" ]; then
  revive -config "$REVIVE_CONFIG" ./...
else
  revive ./...
fi

# Run 'go vet'.
echo "Checking code with 'go vet'..." >&2
go vet ./...
