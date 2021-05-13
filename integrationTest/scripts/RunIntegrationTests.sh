#!/bin/bash

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
branch_dir="$script_dir/../.."

if [ "$1" == "container" ]; then
    samba_exporter="/samba_exporter/samba_exporter"
    samba_statusd="/samba_statusd/samba_statusd"
else
    samba_exporter="$branch_dir/bin/samba_exporter"
    samba_statusd="$branch_dir/bin/samba_statusd"
fi


if [ -f "$samba_exporter" ]; then
    echo "Run $samba_exporter"
    $samba_exporter
else
    echo "Error $samba_exporter not found"
    exit 1
fi

if [ -f "$samba_statusd" ]; then
    echo "Run $samba_statusd"
    $samba_statusd
else
    echo "Error $samba_statusd not found"
    exit 1
fi

exit 0