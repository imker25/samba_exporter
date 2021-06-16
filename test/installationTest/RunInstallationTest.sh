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
echo "sudo apt-get update && sudo apt-get install -y  samba smbclient wget curl coreutils mawk nano"
# Install dependencies for testing 
sudo apt-get update | cat >> /dev/null
if [ "$?" != "0" ]; then 
    echo "Error while apt-get update"
    exit 1
fi
sudo apt-get install -y  samba smbclient wget curl coreutils mawk nano < /bin/true | cat >> /dev/null
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
samba_statusd_log_lines=$(wc -l $tmp_dir/samba_statusd.service.1.log | awk '{print $1}' )
echo "$tmp_dir/samba_exporter.service.1.log has $samba_exporter_log_lines lines"
echo "$tmp_dir/samba_exporter.service.1.log has $samba_statusd_log_lines lines"

assert "echo $samba_exporter_log_lines" "4"
assert "echo $samba_statusd_log_lines" "4"
assert_raises "cat $tmp_dir/samba_exporter.service.1.log | grep \"get metrics on http://127.0.0.1:9922/metrics\"" 0

exporterPID=$(pidof samba_exporter)
uidOfexporterPID=$(awk '/^Uid:/{print $2}' /proc/$exporterPID/status)
userOfexporterPID=$(getent passwd "$uidOfexporterPID" | awk -F: '{print $1}')
assert "echo $userOfexporterPID" "samba-exporter"

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
sleep 0.2
echo "sudo systemctl restart samba_exporter"
sudo systemctl restart samba_exporter
sleep 0.2
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
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_client_count 0\"" 0

echo "Restart samba server with updated settings, so a share is provided"
echo "# ###################################################################"
echo "sudo mkdir -p /srv/test"
sudo mkdir -p /srv/test
echo "sudo chmod 777 /srv/test"
sudo chmod 777 /srv/test
echo "sudo cp \"$script_dir/test.smb.conf\" \"/etc/samba/smb.conf\""
sudo cp "$script_dir/test.smb.conf" "/etc/samba/smb.conf"
echo "sudo systemctl restart smbd.service"
sudo systemctl restart smbd.service
echo "# ###################################################################"
echo "sudo systemctl status smbd.service"
sudo systemctl status smbd.service > "$tmp_dir/samba.service.status.1.log"
cat "$tmp_dir/samba.service.status.1.log"
echo "echo \"My awsome test file\" > /srv/test/test.file"
echo "My awsome test file" > /srv/test/test.file
echo "smbclient -L //127.0.0.1"
smbclient -L //127.0.0.1

echo "# ###################################################################"
echo "Mount samba share"
echo "sudo mkdir /mnt/test"
sudo mkdir /mnt/test
echo "sudo mount -t cifs -o username=guest,password=\"\" //127.0.0.1/test /mnt/test/" 
sudo mount -t cifs -o username=guest,password="" //127.0.0.1/test /mnt/test/
echo "sudo cat /mnt/test/test.file" 
sudo cat /mnt/test/test.file
assert "sudo cat /mnt/test/test.file" "My awsome test file"

echo "# ###################################################################"
echo "sudo smbstatus -L -n"
sudo smbstatus -L -n

echo "# ###################################################################"
echo "sudo smbstatus -S -n"
sudo smbstatus -S -n

echo "# ###################################################################"
echo "sudo smbstatus -p -n"
sudo smbstatus -p -n

echo "# ###################################################################"
echo "curl http://127.0.0.1:9922/metrics"
curl http://127.0.0.1:9922/metrics 

assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_client_count 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_share_count 2\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_individual_user_count 1\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_server_information\"" 0
assert_raises "curl http://127.0.0.1:9922/metrics | grep \"samba_pid_count 1\"" 0

echo "# ###################################################################"
echo "sudo journalctl -u samba_statusd.service "
sudo journalctl -u samba_statusd.service 
echo "# ###################################################################"
echo "sudo journalctl -u samba_exporter.service "
sudo journalctl -u samba_exporter.service 


echo "# ###################################################################"
echo "Check logs before purge"
sudo journalctl -u samba_exporter.service > $tmp_dir/samba_exporter.service.2.log
sudo journalctl -u samba_statusd.service > $tmp_dir/samba_statusd.service.2.log
samba_exporter_log_lines=$(wc -l $tmp_dir/samba_exporter.service.2.log| awk '{print $1}' )
samba_statusd_log_lines=$(wc -l $tmp_dir/samba_statusd.service.2.log | awk '{print $1}' )
echo "$tmp_dir/samba_exporter.service.2.log has $samba_exporter_log_lines lines"
echo "$tmp_dir/samba_exporter.service.2.log has $samba_statusd_log_lines lines"
assert "echo $samba_exporter_log_lines" "28"
assert "echo $samba_statusd_log_lines" "16"

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
assert_raises "fileExists \"/run/samba_exporter.request.pipe\"" 0
assert_raises "fileExists \"/run/samba_exporter.response.pipe\"" 0

if getent passwd samba-exporter > /dev/null; then
    assert "echo \"User samba-exporter exists after purgig the package\"" ""
fi


echo "Tests done"
echo "# ###################################################################"
assert_end samba-exporter_InstallationTests
exit 0