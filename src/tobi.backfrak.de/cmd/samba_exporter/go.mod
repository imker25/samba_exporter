module tobi.backfrak.de/cmd/samba_exporter

go 1.21

require tobi.backfrak.de/internal/commonbl v0.0.0

replace tobi.backfrak.de/internal/commonbl v0.0.0 => ../../internal/commonbl

require tobi.backfrak.de/internal/testhelper v0.0.0

replace tobi.backfrak.de/internal/testhelper v0.0.0 => ../../internal/testhelper

require tobi.backfrak.de/internal/smbexporterbl/pipecomunication v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/pipecomunication v0.0.0 => ../../internal/smbexporterbl/pipecomunication

require tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator v0.0.0 => ../../internal/smbexporterbl/statisticsGenerator

require tobi.backfrak.de/internal/smbexporterbl/smbstatusreader v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/smbstatusreader v0.0.0 => ../../internal/smbexporterbl/smbstatusreader

require tobi.backfrak.de/internal/smbexporterbl/smbexporter v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/smbexporter v0.0.0 => ../../internal/smbexporterbl/smbexporter

replace tobi.backfrak.de/internal/smbstatusout v0.0.0 => ../../internal/smbstatusout

require github.com/prometheus/client_golang v1.19.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
