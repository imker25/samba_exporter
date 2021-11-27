# samba_exporter

A [prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for statistic data of the [samba file server](https://www.samba.org/).

It uses [smbstatus](https://www.samba.org/samba/docs/current/man-html/smbstatus.1.html)  to collect the data and converts the result into prometheus style data.
The prometheus style data can be requested manually on port 9922 using a http client. Or a prometheus database sever can be configured to collect the data by scraping port 9922 on the samba server.

There are packages for several [Ubuntu](https://ubuntu.com/download) and [Debian](https://www.debian.org/) Versions for you ready to install.

## Documentation

For detailed documentation please take a look at the [projects page](https://imker25.github.io/samba_exporter) or read in the [docs](./docs/Index.md) folder.
