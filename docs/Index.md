# samba_exporter

A [prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for statistic data of the [samba file server](https://www.samba.org/).

It uses [smbstatus](https://www.samba.org/samba/docs/current/man-html/smbstatus.1.html)  to collect the data and converts the result into prometheus style data.
The prometheus style data can be requested manually on port 9922 using a http client. Or a prometheus database sever can be configured to collect the data by scraping port 9922 on the samba server.

## Next steps

- [Installation](./Installation/InstallationGuide.md)
