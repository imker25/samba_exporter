#!/bin/bash
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
branch_dir="$script_dir/../.."
tmp_dir="$branch_dir/tmp"

# Run some installation tests for the samba_exporter package
echo "# ###################################################################"
echo "# Prepare for Installation test"
echo "# ###################################################################"
echo "SAMBA_EXPORTER_PACKAGE_NAME=$SAMBA_EXPORTER_PACKAGE_NAME"
echo "script_dir=$script_dir"
echo "branch_dir=$branch_dir"
source "$branch_dir/test/import/functions.sh"

if [ -f "./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb" ]; then
    echo "Install package found on : ${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb"
else
    echo "Error: Installer \"./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb\" not found"
    exit 1
fi

echo "# ###################################################################"
echo "sudo apt-get update && sudo apt-get install -y  samba smbclient wget curl coreutils mawk"
# Install dependencies for testing 
sudo apt-get update | cat >> /dev/null
if [ "$?" != "0" ]; then 
    echo "Error while apt-get update"
    exit 1
fi
sudo apt-get install -y  samba smbclient wget curl coreutils mawk < /bin/true | cat >> /dev/null
if [ "$?" != "0" ]; then 
    echo "Error while samba package installation"
    exit 1
fi

echo "# ###################################################################"
if [ -f "$script_dir/assert.sh" ]; then
    echo "Remove old $script_dir/assert.sh"
    rm -rf "$script_dir/assert.sh"
fi
echo "wget -O "$script_dir/assert.sh" https://raw.githubusercontent.com/lehmannro/assert.sh/v1.1/assert.sh"
wget -O "$script_dir/assert.sh" https://raw.githubusercontent.com/lehmannro/assert.sh/v1.1/assert.sh
 
if [ -f "$script_dir/assert.sh" ]; then
    chmod 700 "$script_dir/assert.sh"
    source "$script_dir/assert.sh"
else
    echo "Error while getting https://github.com/lehmannro/assert.sh"
    exit -1
fi

if [ ! -d "$tmp_dir" ]; then
    mkdir -p "$tmp_dir"
fi

echo "# ###################################################################"
echo "# Install package test"
echo "# ###################################################################"
echo "sudo dpkg --install  \"./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb"
sudo dpkg --install  "./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb"
echo "# ###################################################################"
assert "echo \"$?\"" "0"
sleep 0.4
assert_raises "fileExists \"/etc/default/samba_exporter\"" 1
assert_raises "fileExists \"/etc/default/samba_statusd\"" 1

assert_raises "samba_exporter --help" 0
assert_raises "samba_statusd --help" 0


processWithNameIsRunning samba_statusd
processWithNameIsRunning samba_exporter
assert_raises "processWithNameIsRunning samba_statusd" 1
assert_raises "processWithNameIsRunning samba_exporter" 1

echo "Test Jornal for the servives"
sudo journalctl -u samba_exporter.service > $tmp_dir/samba_exporter.service.1.log
sudo journalctl -u samba_statusd.service > $tmp_dir/samba_statusd.service.1.log
samba_exporter_log_lines=$(wc -l $tmp_dir/samba_exporter.service.1.log| awk '{print $1}' )
samba_statusd_log_lines=$(wc -l $tmp_dir/samba_exporter.service.1.log | awk '{print $1}' )
echo "$tmp_dir/samba_exporter.service.1.log has $samba_exporter_log_lines lines"
echo "$tmp_dir/samba_exporter.service.1.log has $samba_statusd_log_lines lines"

assert "echo $samba_exporter_log_lines" "4"
assert "echo $samba_statusd_log_lines" "4"


echo "# ###################################################################"
echo "Test Service start stop"
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<p><a href='/metrics'>Metrics</a></p>\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<head><title>Samba Exporter</title></head>\"" 0 

echo "# ###################################################################"
echo "sudo systemctl stop samba_satutsd "
sudo systemctl stop samba_statusd
assert_raises "processWithNameIsRunning samba_statusd" 0
assert_raises "processWithNameIsRunning samba_exporter" 0

assert_raises "curl http://127.0.0.1:9922/metrics" 7
echo "sudo systemctl start samba_exporter"
sudo systemctl start samba_exporter
sleep 0.1
assert_raises "processWithNameIsRunning samba_statusd" 1
assert_raises "processWithNameIsRunning samba_exporter" 1

assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<p><a href='/metrics'>Metrics</a></p>\"" 0
assert_raises "curl http://127.0.0.1:9922 | grep \"<head><title>Samba Exporter</title></head>\"" 0 

echo "sudo systemctl stop samba_exporter"
sudo systemctl stop samba_exporter
assert_raises "processWithNameIsRunning samba_statusd" 1
assert_raises "processWithNameIsRunning samba_exporter" 0

echo "sudo systemctl start samba_exporter"
sudo systemctl start samba_exporter
sleep 0.1
assert_raises "processWithNameIsRunning samba_statusd" 1
assert_raises "processWithNameIsRunning samba_exporter" 1

assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0

echo "# ###################################################################"
echo "Test Service restart"
exporterPIDBefore=$(pidof samba_exporter)
echo "samba_exporter running with PID $exporterPIDBefore"
statusdPIDBefore=$(pidof samba_statusd)
echo "samba_statusd running with PID $statusdPIDBefore"

echo "sudo systemctl restart samba_statusd"
sudo systemctl restart samba_statusd
sleep 0.1
echo "sudo systemctl restart samba_exporter"
sudo systemctl restart samba_exporter
sleep 0.1
assert_raises "processWithNameIsRunning samba_statusd" 1
assert_raises "processWithNameIsRunning samba_exporter" 1

exporterPIDAfter=$(pidof samba_exporter)
echo "samba_exporter running with PID $exporterPIDAfter"
statusdPIDAfter=$(pidof samba_statusd)
echo "samba_statusd running with PID $statusdPIDAfter"

if [ "$exporterPIDBefore" == "$exporterPIDAfter" ]; then
    asster "echo \"samba_exporter was not restarted\"" ""
fi

if [ "$statusdPIDBefore" == "$statusdPIDAfter" ]; then
    asster "echo \"samba_exporter was not restarted\"" ""
fi

assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_up 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_satutsd_up 1\"" 0

echo "# ###################################################################"
echo "sudo journalctl -u samba_statusd.service "
sudo journalctl -u samba_statusd.service 
echo "# ###################################################################"
echo "sudo journalctl -u samba_statusd.service "
sudo journalctl -u samba_exporter.service 
echo "# ###################################################################"

echo "# ###################################################################"
echo "# Purge package test"
echo "# ###################################################################"
echo "sudo dpkg --purge samba-exporter"
sudo dpkg --purge samba-exporter
assert "echo \"$?\"" "0"
echo "# ###################################################################"

assert_raises "processWithNameIsRunning samba_statusd" 0
assert_raises "processWithNameIsRunning samba_exporter" 0
assert_raises "fileExists \"/etc/default/samba_exporter\"" 0
assert_raises "fileExists \"/etc/default/samba_statusd\"" 0
assert_raises "fileExists \"/usr/local/bin/start_samba_statusd.sh\"" 0
assert_raises "fileExists \"/usr/local/bin/samba_statusd\"" 0
assert_raises "fileExists \"/usr/local/bin/samba_exporter\"" 0
assert_raises "fileExists \"/etc/systemd/system/samba_exporter.service\"" 0
assert_raises "fileExists \"/etc/systemd/system/samba_statusd.service\"" 0

echo "Tests done"
echo "# ###################################################################"
assert_end samba-exporter_InstallationTests
exit 0