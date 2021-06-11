# samba_exporter

A prometheus exporter for statistic data of the samba file server

**Still in development, but basic metrics are working. Tested only on Ubuntu 20.04**

## Concept

Since the tool [smbstatus](https://www.samba.org/samba/docs/current/man-html/smbstatus.1.html) provided by the [samba](https://www.samba.org/) project can only run with elevated rights, and a [prometheus exporter](https://prometheus.io/docs/instrumenting/writing_exporters/) always exposes a public network endpoint, the samba_exporter package contains two services.

- **samba_exporter**: The prometheus exporter service that exposes the public network endpoint for the prometheus server to collect the data running as non privileged user
- **samba_statusd**: The service uses `smbstatus` to collect data and return it when requested running as privileged user

Both services can communicate using a named pipe owned by a common group.

## Installation

Install the latest [Release](https://github.com/imker25/samba_exporter/releases) by downloading the debian package and installing it. For example (Link and Version needs to be adapted to the latest release):

```sh
wget https://github.com/imker25/samba_exporter/releases/download/0.1.192-pre/samba-exporter_0.1.192-f6b01a7+ubuntu-20.04_amd64.deb
sudo dpkg --install ./samba-exporter_0.1.192-f6b01a7+ubuntu-20.04_amd64.deb
```

By default the prometheus exporter endpoint only listen on localhost. To change this behavior update `/etc/default/samba_exporter` on the target machine according to your needs and restart the `samba_exporter` service. See [samba_statusd service](#samba_statusd-service) for details.

For manual installation see [Build and manual install](#build-and-manual-install).

## Usage

It's is assumed both services are installed as shown in the [installation](#Installation) section.

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
  -test-mode
        Run the program in test mode. In this mode the program will always return the same test data. To work with samba_statusd both programs needs to run in test mode or not.
  -test-pipe
        Requests status from samba_statusd and exits. May be combinde with -test-mode.
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

## Build and manual install

To build the project you need [Go](https://golang.org/) Version 1.16.x and [Java](https://java.com/) Version 11 on your development machine. 
On your target machine, the samba server you want to monitor, you need [samba](https://www.samba.org/) and [systemd](https://www.freedesktop.org/wiki/Software/systemd/) installed.

To build the software change to the repositories directory and run:

```sh
./gradlew build
```

For manual install on the `target` machine do the following copies:

```sh
scp ./bin/samba_exporter <target>:/usr/local/bin/samba_exporter
scp ./bin/samba_statusd <target>:/usr/local/bin/samba_statusd 
scp ./install/usr/local/bin/start_samba_statusd.sh <target>:/usr/local/bin/start_samba_statusd.sh
scp ./install/etc/systemd/system/samba_statusd.service <target>:/etc/systemd/system/samba_statusd.service
scp ./install/etc/systemd/system/samba_exporter.service <target>:/etc/systemd/system/samba_exporter.service
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


