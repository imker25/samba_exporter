#!/bin/bash
# #########################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# #########################################################################################

# #########################################################################################
# Script to start the samba_statusd in productive environemnts
#
# - Will ensure that the pipes for comunication between samba_statusd and samba_exporter
#   are correctly setup and then start samba_statusd with any given paramter
# #########################################################################################
echo "Startup samba_statusd with ARGS: $?"
pipe_permissions="660"
pipe_owner="root:samba-exporter"
request_pipe_file="/run/samba_exporter.request.pipe"
response_pipe_file="/run/samba_exporter.response.pipe"
samba_statusd="/usr/local/bin/samba_statusd"

# Check that samba_statusd is installed as expected
if [ ! -f "$samba_statusd" ]; then
    echo "Error: $samba_statusd not found"
    exit 1
fi

# Setup request pipe
if [ -p "$request_pipe_file" ]; then
    rm "$request_pipe_file"
fi
mkfifo "$request_pipe_file"
chown "$pipe_owner" "$request_pipe_file"
chmod "$pipe_permissions" "$request_pipe_file"

# Setup response pipe
if [ -p "$response_pipe_file" ]; then
    rm "$response_pipe_file"
fi
mkfifo "$response_pipe_file"
chown "$pipe_owner" "$response_pipe_file"
chmod "$pipe_permissions" "$response_pipe_file"

# Run samba_statusd with the given arguments as daemon
echo "Starting as daemon: $samba_statusd $?"
$samba_statusd $* &

# Ensure deamon is up
sleep 0.05

# Exit the startup script
exit 0