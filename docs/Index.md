# Samba exporter

A prometheus exporter for statistic data of the samba file server.

It uses smbstatus to collect the data and converts the result into prometheus style data.
The prometheus style data can be requested manually on port 9922 using a http client. Or a prometheus sever can be configured to collect the data by scraping port 9922 on the samba server.
