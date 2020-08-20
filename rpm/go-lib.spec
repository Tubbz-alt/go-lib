# Run tests in check section
# disable for bootstrapping
%bcond_with check
%global import_path pkg.deepin.io/lib

%global goipath  pkg.deepin.io/lib
%global forgeurl https://github.com/linuxdeepin/go-lib
%global sname go-dlib
%global release_name server-industry

%global with_debug 1

%if 0%{?with_debug}
%global debug_package   %{nil}
%endif
%gometa

Name:           go-lib
Version:        5.4.5
Release:        2
Summary:        Go bindings for Deepin Desktop Environment development
License:        GPLv3
URL:            http://shuttle.corp.deepin.com/cache/tasks/18990/unstable-amd64/
Source0:        %{sname}_%{version}-%{release_name}.orig.tar.xz
BuildRequires:  compiler(go-compiler)

%description
DLib is a set of Go bindings/libraries for DDE development.
Containing dbus (forking from guelfey), glib, gdkpixbuf, pulse and more.

%package devel
Summary:        %{summary}
BuildArch:      noarch
%if %{with check}
# Required for tests
BuildRequires:  deepin-gir-generator
BuildRequires:  dbus-x11
BuildRequires:  iso-codes
BuildRequires:  mobile-broadband-provider-info
BuildRequires:  golang(github.com/linuxdeepin/go-x11-client)
BuildRequires:  golang(github.com/smartystreets/goconvey/convey)
BuildRequires:  golang(gopkg.in/check.v1)
BuildRequires:  pkgconfig(gio-2.0)
BuildRequires:  pkgconfig(gdk-3.0)
BuildRequires:  pkgconfig(gdk-x11-3.0)
BuildRequires:  pkgconfig(gdk-pixbuf-xlib-2.0)
BuildRequires:  pkgconfig(libpulse)
%endif

%description devel
%{summary}.

Provides: golang(pkg.deepin.io/lib)

This package contains library source intended for
building other packages which use import path with
%{goipath} prefix.

%prep
%setup -q -n  %{sname}-%{version}-%{release_name}
%forgeautosetup -n  %{sname}-%{version}-%{release_name}

%install
install -d -p %{buildroot}/%{gopath}/src/%{import_path}/
for file in $(find . -iname "*.go" -o -iname "*.c" -o -iname "*.h") ; do
    install -d -p %{buildroot}/%{gopath}/src/%{import_path}/$(dirname $file)
    cp -pav $file %{buildroot}/%{gopath}/src/%{import_path}/$file
    echo "%%{gopath}/src/%%{import_path}/$file" >> devel.file-list
done

cp -pav README.md %{buildroot}/%{gopath}/src/%{goipath}/README.md
cp -pav CHANGELOG.md %{buildroot}/%{gopath}/src/%{goipath}/CHANGELOG.md
echo "%%{gopath}/src/%%{goipath}/README.md" >> devel.file-list
echo "%%{gopath}/src/%%{goipath}/CHANGELOG.md" >> devel.file-list

%if %{with check}
%check
%gochecks
%endif

%files devel -f devel.file-list
%doc README.md
%license LICENSE

%changelog
* Thu Jun 11 2020 uoser <uoser@uniontech.com> - 5.4.5
- Update to 5.4.5
