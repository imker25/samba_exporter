#!/bin/bash

samba_exporter="/samba_exporter"
samba_statusd="/samba_statusd"

ls -l "/"
ls -l "$samba_exporter"
ls -l "$samba_statusd"

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