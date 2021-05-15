#!/bin/bash

# ###########################################################################################
# Script to run integration tests
#
# Usage: ./RunIntegrationTests.sh [container]
#        container    optional, tell the sript it runs in the github workflow integration test container
#  ###########################################################################################

# ###########################################################################################
# Environment
# ###########################################################################################
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
branch_dir="$script_dir/../.."

if [ "$1" == "container" ]; then
    samba_exporter="/samba_exporter/samba_exporter"
    samba_statusd="/samba_statusd/samba_statusd"
else
    samba_exporter="$branch_dir/bin/samba_exporter"
    samba_statusd="$branch_dir/bin/samba_statusd"
fi

# ###########################################################################################
# Test code
# ###########################################################################################
echo "# ###################################################################"
echo "$(date) - Basic tests"
echo "# ###################################################################"
if [ -f "$samba_exporter" ]; then
    echo "Run: $samba_exporter -print-version"
    $samba_exporter -print-version
    if [ "$?" != "0" ]; then 
        echo "Error while print version of $samba_exporter"
        exit 1
    fi
else
    echo "Error $samba_exporter not found"
    exit 1
fi

if [ -f "$samba_statusd" ]; then
    echo "Run: $samba_statusd -print-version"
    $samba_statusd -print-version
    if [ "$?" != "0" ]; then 
        echo "Error while print version of $samba_statusd"
        exit 1
    fi
else
    echo "Error $samba_statusd not found"
    exit 1
fi

echo "# ###################################################################"
echo "$(date) - Prepare for integration testing"
echo "# ###################################################################"
if [ -f "$script_dir/assert.sh" ]; then
    echo "Remove old $script_dir/assert.sh"
    rm -rf "$script_dir/assert.sh"
fi
wget -O "$script_dir/assert.sh" https://raw.githubusercontent.com/lehmannro/assert.sh/v1.1/assert.sh
 
if [ -f "$script_dir/assert.sh" ]; then
    chmod 700 "$script_dir/assert.sh"
    source "$script_dir/assert.sh"
else
    echo "Error while getting https://github.com/lehmannro/assert.sh"
    exit -1
fi

echo "# ###################################################################"
echo "$(date) - Run integration tests"
echo "# ###################################################################"

# Test the version output
assert_raises "$samba_statusd -version | grep Version: &> /dev/null" 0
assert_raises "$samba_exporter -version | grep Version: &> /dev/null" 0

# Test the help output
assert_raises "$samba_statusd -help | grep \"Usage: \" &> /dev/null" 0
assert_raises "$samba_exporter -help | grep \"Usage: \" &> /dev/null" 0


assert_end samba-exporter_IntegrationTests
exit 0