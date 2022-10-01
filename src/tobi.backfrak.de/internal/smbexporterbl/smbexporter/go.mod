module tobi.backfrak.de/internal/smbexporterbl/smbexporter


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

require github.com/prometheus/client_golang v1.11.0