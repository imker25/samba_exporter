module tobi.backfrak.de/cmd/samba_exporter/smbexporter


require tobi.backfrak.de/internal/commonbl v0.0.0
replace tobi.backfrak.de/internal/commonbl v0.0.0 => ../../../internal/commonbl

require tobi.backfrak.de/cmd/samba_exporter/pipecomunication v0.0.0
replace tobi.backfrak.de/cmd/samba_exporter/pipecomunication v0.0.0 => ../pipecomunication

require tobi.backfrak.de/cmd/samba_exporter/statisticsGenerator v0.0.0
replace tobi.backfrak.de/cmd/samba_exporter/statisticsGenerator v0.0.0 => ../statisticsGenerator

require tobi.backfrak.de/cmd/samba_exporter/smbstatusreader v0.0.0
replace tobi.backfrak.de/cmd/samba_exporter/smbstatusreader v0.0.0 => ../smbstatusreader

require github.com/prometheus/client_golang v1.11.0
