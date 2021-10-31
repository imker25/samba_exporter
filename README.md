# samba_exporter

A prometheus exporter for statistic data of the samba file server.

It uses smbstatus to collect the data and converts the result into prometheus style data.
The prometheus style data can be requested manually on port 9922 using a http client. Or a prometheus database sever can be configured to collect the data by scraping port 9922 on the samba server.

## Concept

Since the tool [smbstatus](https://www.samba.org/samba/docs/current/man-html/smbstatus.1.html) provided by the [samba](https://www.samba.org/) project can only run with elevated rights, and a [prometheus exporter](https://prometheus.io/docs/instrumenting/writing_exporters/) always exposes a public network endpoint, the samba_exporter package contains two services.

- **samba_exporter**: The prometheus exporter service that exposes the public network endpoint for the prometheus server to collect the data running as non privileged user
- **samba_statusd**: The service uses `smbstatus` to collect data and return it when requested running as privileged user

Both services can communicate using a named pipe owned by a common group.

## Installation

### Supported Versions

**Ubuntu:**

| Version | Code Name | Supported |
|---------|-----------|-----------|
| Ubnutu 20.04   | Focal Fossa | yes |
| Ubnutu 21.04   | Hirsute Hippo | no |
| Ubnutu 21.10   | Impish Indri | yes |

### Launchpad

The **samba exporter** package is published on [launchpad](https://launchpad.net/~imker/+archive/ubuntu/samba-exporter-ppa). To install from there do the following commands on any supported Ubuntu version:

```sh
sudo add-apt-repository ppa:imker/samba-exporter-ppa
sudo apt-get update
sudo apt-get install samba-exporter
```

### GitHub

Install the [latest Release](https://github.com/imker25/samba_exporter/releases/latest) (only avalible for Ubuntu 20.04) by downloading the debian package and installing it. For example:

```sh
wget https://github.com/imker25/samba_exporter/releases/download/0.1.192-pre/samba-exporter_0.1.192-f6b01a7+ubuntu-20.04_amd64.deb
sudo dpkg --install ./samba-exporter_0.1.192-f6b01a7+ubuntu-20.04_amd64.deb
```

**Hint:** Link and file name needs to be adapted to the latest release.

It's also possible to download and install pre-releases from the GitHub this way.

For manual installation see [Build and manual install](#build-and-manual-install).

## Usage

It's is assumed both services are installed as shown in the [installation](#Installation) section.

By default the prometheus exporter endpoint only listen on localhost. To change this behavior update `/etc/default/samba_exporter` according to your needs and restart the `samba_exporter` service. See [samba_exporter service](#samba_exporter-service) for details.

### samba_statusd service

To change the behavior of the samba_statusd service update the `/etc/default/samba_statusd` according to your needs. You can add any option shown in the help output of `samba_statusd` to the `ARGS` variable.

Get help:

```sh
samba_statusd -help
samba_statusd: Wrapper for smbstatus. Collects data used by the samba_exporter service.
Program Version: 0.1.164-2c3eda2

Usage: ./bin/samba_statusd [options]
Options:
  -help
        Print this help message
  -print-version
        With this flag the program will only print it's version and exit
  -test-mode
        Run the program in test mode. In this mode the program will always return the same test data. To work with samba_exporter both programs needs to run in test mode or not.
  -verbose
        With this flag the program will print verbose output
```

You may not want to start the service with arguments that will exit before listening starts like `-help` or `-print-version`.

To stop, start or restart the service use `systemctl`, e. g.: `sudo systemctl stop samba_statusd`. To read the log output use `journalctl`, e. g. `sudo journalctl -u samba_statusd`.

**Remark:** Due to the services dependencies `samba_exporter` service stops whenever `samba_statusd` stops. And `samba_statusd` always starts when `samba_exporter` is started if not already running.

### samba_exporter service

To change the behavior of the samba_exporter service update the `/etc/default/samba_exporter` according to your needs. You can add any option shown in the help output of `samba_exporter` to the `ARGS` variable.

Get help:

```sh
samba_exporter -help     
samba_exporter: prometheus exporter for the samba file server. Collects data using the samba_statusd service.
Program Version: 0.1.164-2c3eda2

Usage: ./bin/samba_exporter [options]
Options:
  -help
        Print this help message
  -print-version
        With this flag the program will only print it's version and exit
  -request-timeout int
        The timeout for a request to samba_statusd in seconds (default 5)        
  -test-mode
        Run the program in test mode. In this mode the program will always return the same test data. To work with samba_statusd both programs needs to run in test mode or not.
  -test-pipe
        Requests status from samba_statusd and exits. May be combined with -test-mode.
  -verbose
        With this flag the program will print verbose output
  -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9922")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
```

You may not want to start the service with arguments that will exit before listening starts like `-test-pipe`, `-help` or `-print-version`.

To stop, start or restart the service use `systemctl`, e. g.: `sudo systemctl stop samba_exporter`. To read the log output use `journalctl`, e. g. `sudo journalctl -u samba_exporter`.

**Remark:** Due to the services dependencies `samba_exporter` service stops whenever `samba_statusd` stops. And `samba_statusd` always starts when `samba_exporter` is started if not already running.

## Prometheus

To add this exporter to your [prometheus database](https://prometheus.io/) you have to add the endpoint as scrape job to the `/etc/prometheus/prometheus.yml` on your prometheus server. Therefor add the lines shown below:

```yaml
  - job_name: 'Samba exporter node on server.local'
    metrics_path: metrics
    static_configs:
      - targets: ['server.local:9922']
```

Replace `server.local` with the network name of your samba server.

## Grafana

For [grafana](https://grafana.com) an example dashboard is installed with the debian package and can be found at `/usr/share/doc/samba_exporter-Vx.y/grafana/SambaService.json` (Replace x.y with the current version).

When [importing](https://grafana.com/docs/grafana/latest/dashboards/export-import/#import-dashboard) this dashboard you need to change `server.local` to the network name of your samba server.

## Developer Documentation

### Build and manual install

To build the project you need [Go](https://golang.org/) Version 1.16.x and [Java](https://java.com/) Version 11 on your development machine. 
On your target machine, the samba server you want to monitor, you need [samba](https://www.samba.org/) and [systemd](https://www.freedesktop.org/wiki/Software/systemd/) installed.

To build the software change to the repositories directory and run:

```sh
./gradlew build
```

For manual install on the `target` machine do the following copies:

```sh
scp ./bin/samba_exporter <target>:/usr/bin/samba_exporter
scp ./bin/samba_statusd <target>:/usr/bin/samba_statusd 
scp ./install/usr/bin/start_samba_statusd <target>:/usr/bin/start_samba_statusd
scp ./install/lib/systemd/system/samba_statusd.service <target>:/lib/systemd/system/samba_statusd.service
scp ./install/lib/systemd/system/samba_exporter.service <target>:/lib/systemd/system/samba_exporter.service
scp install/etc/default/samba_exporter <target>:/etc/default/samba_exporter
scp install/etc/default/samba_statusd <target>:/etc/default/samba_statusd
```

Now login to your target machine and run the commands below to enable the services and create the needed user and group:

```sh
systemctl daemon-reload
systemctl enable samba_statusd.service
systemctl enable samba_exporter.service
addgroup --system samba-exporter
adduser --system --no-create-home --disabled-login samba-exporter
adduser samba-exporter samba-exporter
```

Finally you are abel to start the services:

```sh
systemctl start samba_statusd.service
systemctl start samba_exporter.service
```

Test the `samba_exporter` by requesting metrics with `curl`:

```sh
curl http://127.0.0.1:9922/metrics 
```

The output of this test should look something like this:

```txt
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
...
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP samba_client_count Number of clients using the samba server
# TYPE samba_client_count gauge
samba_client_count 0
# HELP samba_individual_user_count The number of users connected to this samba server
# TYPE samba_individual_user_count gauge
samba_individual_user_count 0
# HELP samba_locked_file_count Number of files locked by the samba server
# TYPE samba_locked_file_count gauge
samba_locked_file_count 0
# HELP samba_pid_count Number of processes running by the samba server
# TYPE samba_pid_count gauge
samba_pid_count 0
# HELP samba_satutsd_up 1 if the samba_statusd seems to be running
# TYPE samba_satutsd_up gauge
samba_satutsd_up 1
# HELP samba_server_up 1 if the samba server seems to be running
# TYPE samba_server_up gauge
samba_server_up 1
# HELP samba_share_count Number of shares used by clients of the samba server
# TYPE samba_share_count gauge
samba_share_count 0
```

### Man page creation

Man pages are written in [ronn](https://github.com/rtomayko/ronn). The `*.ronn`source files are converted by the script `build/CreateManPage.sh` into man pages.

### CI/CD Pipeline

For continuous integration and deployment this project uses [GitHub Actions](https://github.com/imker25/samba_exporter/actions). The main pipeline is defined in `.github/workflows/ci-jobs.yml`. This pipeline will do:

- On push to any branch on github
  - Build the project and the man pages
  - Run unit tests defined in `*_test.go`
  - Run integration tests from `test/integrationTest/scripts/RunIntegrationTests.sh`
  - Run installation tests from `test/installationTest/RunInstallationTest.sh`
  - Build a debian binary package (`*.deb`)
- On push to main and release/* branch additionally
  - Upload the binary package (`*.deb`) as [GitHub Release](https://github.com/imker25/samba_exporter/releases)
  - In case it's the main branch the release will be a pre release
  - On release/* branches it will be a full public release
  
### Release process

The release process of this project is fully automated. To create a new release of the software use the script `build/PrepareRelease.sh`. Before running the script ensure you got the latest changes from github origin. This script then will:

- Create a release branch from the current state at the main branch
- Update the `VersionMaster.txt` with a new increment version on main branch
- Update the `changelog` with a stub entry for the new version on main branch
- Commit the changes on the main branch
- Push all changes on main and the new release branch to github

Once this changes are pushed to github the CI/CD pipeline will start to run for both, main and the new release/* branch.

After a full public release is done from the the CI/CD run on the release/* branch `.github/workflows/release-jobs.yml` will be triggered. This job runs `build/PublishLaunchpadInDocker.sh` to transfer the just created github release to the [samba-exporter launchpad ppa](https://launchpad.net/~imker/+archive/ubuntu/samba-exporter-ppa) where it will be published as well. During this process a slightly modified version of the sources will be pushed into the corresponding [launchpad git repository](https://code.launchpad.net/~imker/samba-exporter/+git/samba-exporter).

### Developer Hints

In case you want develop this software with [VS Code](https://code.visualstudio.com/) you need to add the repositories root folder to the **GOPATH** within the `VS Code Settings` to get golang extension and golang tools work, e. g.:

```json
{
      "go.gopath": "${env:GOPATH}:${workspaceFolder}",
}
```