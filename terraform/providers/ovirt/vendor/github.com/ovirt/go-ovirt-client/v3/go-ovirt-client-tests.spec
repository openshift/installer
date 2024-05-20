Name:		go-ovirt-client-tests
Version:	1.0.0
Release:	0.0.master%{?release_suffix}%{?dist}
License:	ASL 2.0
Summary:	go ovirt client with current patch
Group:		Virtualization/Management
URL:		https://github.com/oVirt/go-ovirt-client
BuildArch:	x86_64
Source:         go-ovirt-client-tests.tar.gz

%description
go ovirt client tests binary compiled with current patch

%install
mkdir -p %{buildroot}/usr/bin/
cp %{_sourcedir}/go-ovirt-client-tests-exe %{buildroot}/usr/bin/

%files
/usr/bin/go-ovirt-client-tests-exe

%changelog
* Sun Mar 20 2022 Eli Mesika <emesika@redhat.com> 1.0.0
- Created
