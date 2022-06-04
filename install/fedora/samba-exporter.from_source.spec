# Spec file to build a samba_exporter RPM source or binary package 

%global goipath    tobi.backfrak.de/cmd/samba_exporter

Name: samba-exporter
Version: x.x.x
Source0: https://github.com/imker25/samba_exporter/archive/refs/tags/X.X.X-pre.tar.gz
%global tag x.x.x-pre

%gometa
%global debug_package %{nil}

Release: 1
Summary: Prometheus exporter to get metrics of a samba server
License: ASL 2.0
URL: https://github.com/imker25/samba_exporter
Distribution: Fedora
Group: utils
Requires: samba, systemd, filesystem, binutils, man-db 

BuildRequires:  systemd-rpm-macros
BuildRequires:  go-rpm-macros
BuildRequires:  golang(github.com/go-kit/log)
BuildRequires:  golang(github.com/go-kit/log/level)	
BuildRequires:  golang(github.com/prometheus/client_golang/prometheus)
BuildRequires:  golang(github.com/prometheus/client_golang/prometheus/collectors)
BuildRequires:  golang(github.com/prometheus/client_golang/prometheus/promhttp)
BuildRequires:  golang(github.com/prometheus/client_model/go)
BuildRequires:  golang(golang.org/x/sys/unix)
BuildRequires:  golang(gopkg.in/alecthomas/kingpin.v2)
BuildRequires:  rubygem-ronn-ng

%define _rpmfilename %%{NAME}-%%{VERSION}-%%{RELEASE}.%%{ARCH}.rpm
%define _unpackaged_files_terminate_build 0

%gopkg

%prep
%goprep
mkdir -p "%{gobuilddir}/src"
cp -rpv "%{gobuilddir}/src/tobi.backfrak.de/cmd/samba_exporter/samba_exporter-%{tag}/"* "%{gobuilddir}/src/"
echo "%{tag}" > "%{gobuilddir}/src/VersionMaster.txt"
	
%build
GOPATH="$GOPATH:%{gobuilddir}/src/"
export BUILDTAGS="netgo osusergo static_build"
LDFLAGS="-X main.version=%{tag}" \
%gobuild -o %{gobuilddir}/bin/samba_exporter src/tobi.backfrak.de/cmd/samba_exporter
LDFLAGS="-X main.version=%{tag}" \
%gobuild -o %{gobuilddir}/bin/samba_statusd src/tobi.backfrak.de/cmd/samba_statusd
"%{gobuilddir}/src/build/CreateManPage.sh"


%install
"%{gobuilddir}/src/build/InstallProgram.sh" "%{gobuilddir}/src" %{gobuilddir}/bin %{buildroot}/
install -m 664  "%{gobuilddir}/src/LICENSE" "%{buildroot}/usr/share/doc/samba-exporter/LICENSE"
install -d -m 775 "%{buildroot}/usr/share/man/man1/"
install -m 664  "%{gobuilddir}/src/src/man/samba_exporter.1.gz" "%{buildroot}/usr/share/man/man1/samba_exporter.1.gz" 
install -m 664  "%{gobuilddir}/src/src/man/samba_statusd.1.gz" "%{buildroot}/usr/share/man/man1/samba_statusd.1.gz" 
install -m 664  "%{gobuilddir}/src/src/man/start_samba_statusd.1.gz" "%{buildroot}/usr/share/man/man1/start_samba_statusd.1.gz" 


%check
echo "Testing .........................................................................................."

%pre
if [ $1 == 2 ];then
    # Stop services before install in case of package upgrade
    systemctl stop samba_exporter.service
    systemctl stop samba_statusd.service
fi


%post
# Add samba-exporter user if needed
if ! getent group samba-exporter > /dev/null ; then
    groupadd -r samba-exporter
fi
if ! getent passwd samba-exporter > /dev/null ; then
    adduser --system --no-create-home --home-dir /nonexistent --gid samba-exporter --shell /bin/false --comment "samba-exporter daemon" samba-exporter || true
fi
# Ensure the daemons are known
systemctl daemon-reload
if [ $1 == 1 ];then
    # Ensure the daemons start automaticaly in case of package installation
    systemctl enable samba_statusd.service
    systemctl enable samba_exporter.service
fi
# Ensure the daemons run the latest version
systemctl start samba_statusd.service
systemctl start samba_exporter.service
# Ensure man-db is updated
mandb > /dev/null


%preun
if [ $1 == 0 ];then
    request_pipe_file="/run/samba_exporter.request.pipe"
    response_pipe_file="/run/samba_exporter.response.pipe"
    # Stop the services before removing the package
    systemctl stop samba_statusd.service
    systemctl stop samba_exporter.service
    if [ -p "$request_pipe_file" ]; then
        rm "$request_pipe_file"
    fi
    if [ -p "$response_pipe_file" ]; then
        rm "$response_pipe_file"
    fi
fi

%postun
if [ $1 == 0 ];then
    # When the package got removed the service files got deleted. So systemd can now remove the services from its internal db
    systemctl daemon-reload
    if [ -d "/usr/share/doc/samba-exporter" ]; then 
        rm -rf "/usr/share/doc/samba-exporter"
    fi
fi


%description
 This is a prometheus exporter to get metrics of a samba server.
 It uses smbstatus to collect the data and converts the result into
 prometheus style data.
 The prometheus style data can be requested manually on port 9922
 using a http client. Or a prometheus database sever can be configured
 to collect the data by scraping port 9922 on the samba server.


%files
%config(noreplace) "/etc/default/samba_exporter"
%config(noreplace) "/etc/default/samba_statusd"
"/lib/systemd/system/samba_exporter.service"
"/lib/systemd/system/samba_statusd.service"
"/usr/bin/samba_exporter"
"/usr/bin/samba_statusd"
"/usr/bin/start_samba_statusd"
%dir "/usr/share/"
%docdir "/usr/share/doc/"
%docdir "/usr/share/doc/samba-exporter/"
%doc "/usr/share/doc/samba-exporter/README.md"
%license "/usr/share/doc/samba-exporter/LICENSE"
%docdir "/usr/share/doc/samba-exporter/docs/"
%docdir "/usr/share/doc/samba-exporter/docs/DeveloperDocs/"
%doc "/usr/share/doc/samba-exporter/docs/DeveloperDocs/ActionsAndReleases.md"
%doc "/usr/share/doc/samba-exporter/docs/DeveloperDocs/Compile.md"
%doc "/usr/share/doc/samba-exporter/docs/DeveloperDocs/Hints.md"
%doc "/usr/share/doc/samba-exporter/docs/Index.md"
%docdir "/usr/share/doc/samba-exporter/docs/Installation/"
%doc "/usr/share/doc/samba-exporter/docs/Installation/InstallationGuide.md"
%doc "/usr/share/doc/samba-exporter/docs/Installation/SupportedVersions.md"
%docdir "/usr/share/doc/samba-exporter/docs/UserDocs/"
%doc "/usr/share/doc/samba-exporter/docs/UserDocs/Concept.md"
%doc "/usr/share/doc/samba-exporter/docs/UserDocs/ServiceIntegration.md"
%doc "/usr/share/doc/samba-exporter/docs/UserDocs/UserGuide.md"
%docdir "/usr/share/doc/samba-exporter/docs/assets/"
%doc "/usr/share/doc/samba-exporter/docs/assets/Samba-Dashboard.png"
%doc "/usr/share/doc/samba-exporter/docs/assets/samba-exporter.icon.png"
%docdir "/usr/share/doc/samba-exporter/grafana/"
%doc "/usr/share/doc/samba-exporter/grafana/SambaService.json"
%doc "/usr/share/man/man1/samba_exporter.1.gz"
%doc "/usr/share/man/man1/samba_statusd.1.gz"
%doc "/usr/share/man/man1/start_samba_statusd.1.gz"

