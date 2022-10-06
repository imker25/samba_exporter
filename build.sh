#!/bin/bash

# Script to run the mage build workflow for the project.
# - running the script without paramates just recompiles the project
# - running the script with -l list the targets
# - running the script with -h shows the mage help


SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd "$SCRIPT_DIR/build/workflow"
go run ./mage.go -d "$SCRIPT_DIR/build/workflow/magefiles" -w "$SCRIPT_DIR" "$@"
build_exit_code=$?
if [ "$build_exit_code" != 0 ]; then
    echo "Error: 'go run' exit with code '$build_exit_code'"
fi
popd

exit $build_exit_code