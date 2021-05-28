#!/bin/bash
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
branch_dir="$script_dir/../.."

# Run some installation tests for the samba_exporter package
echo "# ###################################################################"
echo "# Prepare for Installation test"
echo "# ###################################################################"
echo "SAMBA_EXPORTER_PACKAGE_NAME=$SAMBA_EXPORTER_PACKAGE_NAME"
echo "script_dir=$script_dir"
echo "branch_dir=$branch_dir"

if [ -f "./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb" ]; then
    echo "Install package found on : ${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb"
else
    echo "Error: Installer \"./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb\" not found"
    exit 1
fi

echo "# ###################################################################"
echo "sudo apt-get update && sudo apt-get install -y  samba smbclient wget curl coreutils mawk"
# Install dependencies for testing 
sudo apt-get update && sudo apt-get install -y  samba smbclient wget curl coreutils mawk < /bin/true

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

echo "# ###################################################################"
echo "# Installation test"
echo "# ###################################################################"
echo "sudo dpkg --install  \"./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb"
sudo dpkg --install  \"./${SAMBA_EXPORTER_PACKAGE_NAME}_amd64.deb\"
echo "# ###################################################################"
assert "echo \"$?\"" "0"
assert_raises "samba_exporter --help" 0
assert_raises "samba_statusd --help" 0


exporterPID=$(pidof samba_exporter)
echo "$samba_exporter running with PID $exporterPID"
statusdPID=$(pidof samba_statusd)
echo "$samba_statusd running with PID $statusdPID"


echo "# ###################################################################"
echo "# Purge test"
echo "# ###################################################################"
echo "sudo dpkg --purge samba-exporter"
sudo dpkg --purge samba-exporter
assert "echo \"$?\"" "0"
echo "# ###################################################################"

assert_end samba-exporter_InstallationTests
exit 0