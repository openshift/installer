Name:		tr-ocp-installer
Version:	1.0.1
Release:	0.0.master%{?release_suffix}%{?dist}
License:	ASL 2.0
Summary:	OCP installer with terraform provider patch
Group:		Virtualization/Management
URL:		https://github.com/openshift/installer
BuildArch:	x86_64
Source:         tr-ocp-installer.tar.gz

%description
OCP installer binary compiled with terraform provider patch and terraform provider binary

%install
mkdir -p %{buildroot}/usr/bin/
cp %{_sourcedir}/openshift-install %{_sourcedir}/terraform-provider-ovirt_linux_amd64 %{buildroot}/usr/bin/

%files
/usr/bin/openshift-install
/usr/bin/terraform-provider-ovirt_linux_amd64

%changelog
* Thu Feb 2 2022 Eli Mesika <emesika@redhat.com> 1.0.1
- Added terraform provider executable
* Thu Jan 13 2022 Eli Mesika <emesika@redhat.com> 1.0.0
- Created
