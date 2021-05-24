#!/bin/bash

request_pipe_file="/run/samba_exporter.request.pipe"
response_pipe_file="/run/samba_exporter.response.pipe"
samba_exporter="/usr/bin/samba_exporter"
samba_statusd="/usr/bin/start_samba_statusd.sh"
samba_statusd_log="/var/log/samba_statusd.log"
samba_exporter_log="/var/log/samba_exporter.log"

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
echo "$(date) - Prepare for System testing"
echo "# ###################################################################"
systemctl daemon-reload
echo "systemctl stop samba_statusd.service"
systemctl stop samba_statusd.service
echo "systemctl stop samba_exporter.service"
systemctl stop samba_statusd.service

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
echo "id samba-exporter"
id samba-exporter
echo "# ###################################################################"
echo "$(date) - Run System tests"
echo "# ###################################################################"

echo "# ###################################################################"
echo "Run in daemons verbose mode"
echo "# ###################################################################"
echo "$samba_statusd -verbose &"
$samba_statusd -verbose  &
sleep 0.1
statusdPID=$(pidof samba_statusd)
echo "$samba_statusd running with PID $statusdPID"
echo " "
echo "$su -s /bin/bash  samba-exporter -c \"samba_exporter -verbose &\""
su -s /bin/bash  samba-exporter -c "$samba_exporter -verbose  &"
sleep 0.1
exporterPID=$(pidof samba_exporter)
echo "$samba_exporter running with PID $exporterPID"

echo "# ###################################################################"
echo "Get the enpoint:"
echo "Call: curl http://127.0.0.1:9922"
curl http://127.0.0.1:9922
echo " "
echo "# ###################################################################"
echo "Get metrics"
echo "Call: curl http://127.0.0.1:9922/metrics"
curl http://127.0.0.1:9922/metrics 
echo "# ###################################################################"

echo "Test Web Interface"
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<p><a href='/metrics'>Metrics</a></p>\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<head><title>Samba Exporter</title></head>\"" 0 

# End daemons
echo "# ###################################################################"
echo "End $samba_statusd with PID $statusdPID"
kill $statusdPID
echo "End $samba_exporter with PID $exporterPID"
kill $exporterPID
echo "# ###################################################################"

echo "# ###################################################################"
echo "Test Services"
echo "# ###################################################################"
echo "systemctl start samba_statusd.service"
systemctl start samba_statusd.service
echo "systemctl start samba_exporter.service"
systemctl start samba_statusd.service
echo "# ###################################################################"
exporterPID=$(pidof samba_exporter)
echo "$samba_exporter running with PID $exporterPID"
exporterPID=$(pidof samba_exporter)
echo "$samba_exporter running with PID $exporterPID"
echo "# ###################################################################"

echo "# ###################################################################"
echo "Get the enpoint:"
echo "Call: curl http://127.0.0.1:9922"
curl http://127.0.0.1:9922
echo " "
echo "# ###################################################################"
echo "Get metrics"
echo "Call: curl http://127.0.0.1:9922/metrics"
curl http://127.0.0.1:9922/metrics 
echo "# ###################################################################"

echo "Test Web Interface"
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<p><a href='/metrics'>Metrics</a></p>\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<head><title>Samba Exporter</title></head>\"" 0 

echo "# ###################################################################"
echo "$(date) End Tests"
echo "# ###################################################################"
# Finish test run
assert_end samba-exporter_IntegrationTests
exit 0