#!/bin/bash
# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ###########################################################################################
# Script to run integration tests for the samba_exporter binaries
# Used inside the ./../docker/Dockerfile.ubuntu container by the GitHub CI/CD workflow
#
# Usage: ./RunIntegrationTests.sh [container]
#        container    optional, tell the sript it runs in the github workflow integration test container
#  ###########################################################################################

# ###########################################################################################
# Environment
# ###########################################################################################
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
branch_dir="$script_dir/../../.."
request_pipe_file="/dev/shm/samba_exporter.request.pipe"
response_pipe_file="/dev/shm/samba_exporter.response.pipe"

if [ "$1" == "container" ]; then
    source "/functions.sh"
    samba_exporter="/samba_exporter/samba_exporter"
    samba_statusd="/samba_statusd/samba_statusd"
    samba_statusd_log="/samba_statusd/samba_statusd.log"
    samba_exporter_log="/samba_exporter/samba_exporter.log"
else
    source "$branch_dir/test/import/functions.sh"
    samba_exporter="$branch_dir/bin/samba_exporter"
    samba_statusd="$branch_dir/bin/samba_statusd"
    samba_statusd_log="$branch_dir/logs/samba_statusd.log"
    samba_exporter_log="$branch_dir/logs/samba_exporter.log"
    if [ -d "$branch_dir/logs/" ]; then
        if [ -f "$samba_statusd_log" ]; then 
            rm "$samba_statusd_log"
        fi
        if [ -f "$samba_exporter_log" ]; then 
            rm "$samba_exporter_log"
        fi
    else 
        mkdir -p "$branch_dir/logs/"
    fi
fi


# ###########################################################################################
# Test code
# ###########################################################################################
echo "# ###################################################################"
echo "$(date) - OS Info"
echo "# ###################################################################"
cat /etc/os-release
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

if [ -p "$request_pipe_file" ]; then
    echo "Delete $request_pipe_file"
    rm "$request_pipe_file"
fi
if [ -p "$response_pipe_file" ]; then
    echo "Delete $response_pipe_file"
    rm "$response_pipe_file"
fi

echo "# ###################################################################"
echo "Start as daemon: $samba_statusd -test-mode -verbose"
$samba_statusd -test-mode -verbose &
sleep 0.2
statusdPID=$(pidof samba_statusd)

# Wait a bit to ensure the process is running
sleep 0.1
echo "# ###################################################################"
echo "$samba_statusd running with PID $statusdPID"
assert_raises "processWithNameIsRunning samba_statusd" 1

echo "# ###################################################################"
echo "Test IPC"
# Show the output of -test-pipe mode for debug propose
echo "$samba_exporter -test-mode -verbose -test-pipe"
$samba_exporter -test-mode -verbose -test-pipe

echo "# ###################################################################"
# Test the response code of -test-pipe mode
assert_raises "$samba_exporter -test-mode -test-pipe" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe" 0

# Test the output of -test-pipe mode
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"PID: 1117; UserID: 1080; GroupID: 117; Machine: 192.168.1.242 (ipv4:192.168.1.242:42296); ProtocolVersion: SMB3_11; Encryption: -; Signing: partial(AES-128-CMAC);\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"Service: IPC$; PID: 1119; Machine: 192.168.1.242; ConnectedAt: 2021-05-16T11:55:36\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"PID: 1120; UserID: 1080; DenyMode: DENY_NONE; Access: 0x80; AccessMode: RDONLY; Oplock: NONE; SharePath: /usr/share/data; Name: .: Time 2021-05-16T12:07:02Z;\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_individual_user_count: 1\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_pid_count: 3\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_smbd_virtual_memory_usage_percent\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_smbd_sum_virtual_memory_usage_percent\"" 0

# Test the -not-expose-* options
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_lock_created_since_seconds\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe -not-expose-user-data | grep \"samba_lock_created_since_seconds\"" 1
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_client_connected_since_seconds\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe -not-expose-client-data | grep \"samba_client_connected_since_seconds\"" 1
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"samba_protocol_version_count\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe -not-expose-encryption-data | grep \"samba_protocol_version_count\"" 1
assert_raises "$samba_exporter -test-mode -verbose -test-pipe | grep \"smbd_io_counter_write_count\"" 0
assert_raises "$samba_exporter -test-mode -verbose -test-pipe -not-expose-pid-data | grep \"smbd_io_counter_write_count\"" 1

normalOutPutLineCount=$( $samba_exporter -test-mode -test-pipe | wc -l )
noEncryptionOutPutLineCount=$( $samba_exporter -test-mode -test-pipe -not-expose-encryption-data | wc -l )
noClientOutPutLineCount=$( $samba_exporter -test-mode -test-pipe -not-expose-client-data | wc -l )
noUserOutPutLineCount=$( $samba_exporter -test-mode -test-pipe -not-expose-user-data | wc -l )
noPidOutPutLineCount=$( $samba_exporter -test-mode -test-pipe -not-expose-pid-data | wc -l )
assert_raises " [ $normalOutPutLineCount == $noEncryptionOutPutLineCount ] " 1
assert_raises " [ $normalOutPutLineCount == $noClientOutPutLineCount ] " 1
assert_raises " [ $normalOutPutLineCount == $noUserOutPutLineCount ] " 1
assert_raises " [ $normalOutPutLineCount == $noPidOutPutLineCount ] " 1
echo "normalOutPutLineCount: $normalOutPutLineCount"
echo "noClientOutPutLineCount: $noClientOutPutLineCount"
echo "noEncryptionOutPutLineCount: $noEncryptionOutPutLineCount"
echo "noUserOutPutLineCount: $noUserOutPutLineCount"
echo "noUserOutPutLineCount: $noPidOutPutLineCount"

echo "# ###################################################################"
echo "Normal test-mode output"
$samba_exporter -test-mode -test-pipe 

echo "# ###################################################################"
echo "Start as daemon: $samba_exporter -test-mode -verbose"
$samba_exporter -test-mode -verbose &
sleep 0.1
exporterPID=$(pidof samba_exporter)

# Wait a bit to ensure the process is running
sleep 0.1
echo "# ###################################################################"
echo "$samba_exporter running with PID $exporterPID"
assert_raises "processWithNameIsRunning samba_exporter" 1

echo "# ###################################################################"
echo "Test samba_exporter webinterface"

# Get the outputs of the promethues web requests for debug propose
echo "# ###################################################################"
echo "Get the enpoint:"
echo "Call: curl http://127.0.0.1:9922"
curl http://127.0.0.1:9922
echo " "
echo "# ###################################################################"
echo "Get metrics"
echo "Call: curl http://127.0.0.1:9922/metrics"
curl http://127.0.0.1:9922/metrics
echo " "
echo "# ###################################################################"

# Test the output of promethues web requests 
assert_raises "curl http://127.0.0.1:9922 | grep \"<p><a href='/metrics'>Metrics</a></p>\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<head><title>Samba Exporter</title></head>\"" 0 
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"promhttp_metric_handler_requests_total\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"process_virtual_memory_max_bytes\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"promhttp_metric_handler_requests_in_flight 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_individual_user_count\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"# HELP samba_individual_user_count The number of users connected to this samba server\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_individual_user_count 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"# TYPE samba_individual_user_count gauge\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"/usr/share/data\"" 0


# End daemons
echo "# ###################################################################"
echo "End $samba_statusd with PID $statusdPID"
kill $statusdPID
echo "End $samba_exporter with PID $exporterPID"
kill $exporterPID

assert_raises "processWithNameIsRunning samba_exporter" 0
assert_raises "processWithNameIsRunning samba_statusd" 0


echo "# ###################################################################"
echo "Test in daemons non verbose mode"
echo "# ###################################################################"
echo "$samba_statusd -test-mode > $samba_statusd_log 2>&1 &"
$samba_statusd -test-mode > $samba_statusd_log 2>&1 &
sleep 0.1
statusdPID=$(pidof samba_statusd)
sleep 0.1
echo "$samba_exporter -test-mode > $samba_exporter_log 2>&1 &"
$samba_exporter -test-mode > $samba_exporter_log 2>&1 &
sleep 0.1
exporterPID=$(pidof samba_exporter)
sleep 0.1
echo "$samba_statusd running with PID $statusdPID"
echo "$samba_exporter running with PID $exporterPID"
echo "# ###################################################################"

assert_raises "processWithNameIsRunning samba_exporter" 1
assert_raises "processWithNameIsRunning samba_statusd" 1

echo "Test Web Interface"
echo "# ###################################################################"
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"/usr/share/data\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<p><a href='/metrics'>Metrics</a></p>\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<head><title>Samba Exporter</title></head>\"" 0 
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_exporter_information\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_information\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_process_per_client_count\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"SMB3_11\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"4.11.6-Ubuntu\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_protocol_version_count\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_signing_method_count\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"partial(AES-128-CMAC)\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_encryption_method_count\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_request_time\"" 0

# End daemons
echo "End $samba_statusd with PID $statusdPID"
kill $statusdPID
echo "End $samba_exporter with PID $exporterPID"
kill $exporterPID
echo "# ###################################################################"

assert_raises "processWithNameIsRunning samba_exporter" 0
assert_raises "processWithNameIsRunning samba_statusd" 0

echo "Print Log files"
echo "cat $samba_exporter_log"
cat $samba_exporter_log
echo "# ###################################################################"
echo "cat $samba_statusd_log"
cat $samba_statusd_log
echo "# ###################################################################"

echo "Check log files"
samba_exporter_log_lines=$(wc -l $samba_exporter_log | awk '{print $1}' )
samba_statusd_log_lines=$(wc -l $samba_statusd_log | awk '{print $1}' )
echo "$samba_exporter_log has $samba_exporter_log_lines lines"
echo "$samba_statusd_log has $samba_statusd_log_lines lines"

assert "echo $samba_exporter_log_lines" "1"
assert "echo $samba_statusd_log_lines" "1"

echo "End Tests"
echo "# ###################################################################"
rm -v /dev/shm/samba_exporter.re*
# Finish test run
assert_end samba-exporter_IntegrationTests

exit 0