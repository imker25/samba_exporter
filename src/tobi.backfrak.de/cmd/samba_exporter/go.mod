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

require github.com/prometheus/client_golang v1.14.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
