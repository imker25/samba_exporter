module tobi.backfrak.de/internal/smbexporterbl/smbexporter

go 1.21

require tobi.backfrak.de/internal/commonbl v0.0.0

replace tobi.backfrak.de/internal/commonbl v0.0.0 => ../../commonbl

require tobi.backfrak.de/internal/smbexporterbl/pipecomunication v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/pipecomunication v0.0.0 => ../pipecomunication

require tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/statisticsGenerator v0.0.0 => ../statisticsGenerator

require tobi.backfrak.de/internal/smbexporterbl/smbstatusreader v0.0.0

replace tobi.backfrak.de/internal/smbexporterbl/smbstatusreader v0.0.0 => ../smbstatusreader

require tobi.backfrak.de/internal/smbstatusout v0.0.0

replace tobi.backfrak.de/internal/smbstatusout v0.0.0 => ../../smbstatusout

require github.com/prometheus/client_golang v1.19.0

require tobi.backfrak.de/internal/testhelper v0.0.0

replace tobi.backfrak.de/internal/testhelper v0.0.0 => ../../../internal/testhelper

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
