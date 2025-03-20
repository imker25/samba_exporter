# User Guide

For user documentation please read the man pages:

- [samba_exporter](https://imker25.github.io/samba_exporter/manpages/samba_exporter.1.html)
- [samba_statusd](https://imker25.github.io/samba_exporter/manpages/samba_statusd.1.html)
- [start_samba_statusd](https://imker25.github.io/samba_exporter/manpages/start_samba_statusd.1.html)

In case you installed the package already, you can read the man pages using man. For example:

```bash
man samba_exporter
```

## Exported values

A list of all exported values can be found at the [samba_exporter man page](https://imker25.github.io/samba_exporter/manpages/samba_exporter.1.html).

## Problems and solutions

In case you are running into problems when requesting data from the samba-exporter service. Please first of all stop both required services:

```sh
sudo systemctl stop samba_exporter                        
sudo systemctl stop samba_statusd
``` 

Delete the pipes used for communication between the services:

```sh
sudo rm -f /run/samba_exporter.*     
```

Start the services again:

```sh
sudo systemctl start samba_statusd                          
sudo systemctl start samba_exporter                         
```
