# This Dockerfile is a used by CI to test a libvirt cluster launched in a gce instance
# It builds an image containing google-cloud-sdk, ns_wrapper and scripts to launch a VM for a libvirt install.
FROM openshift/origin-release:golang-1.10 AS build
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN TAGS=libvirt_destroy hack/build.sh

FROM centos:7
COPY --from=build /go/src/github.com/openshift/installer/bin/openshift-install /bin/openshift-install
COPY --from=build /go/src/github.com/openshift/installer/images/nested-libvirt/google-cloud-sdk.repo /etc/yum.repos.d/google-cloud-sdk.repo
COPY --from=build /go/src/github.com/openshift/installer/images/nested-libvirt/mock-nss.sh /bin/mock-nss.sh

RUN yum install -y \
    epel-release \
    gettext \
    google-cloud-sdk \
    openssh-clients && \
    yum --enablerepo=epel-testing install -y nss_wrapper && \
    yum -y update && \
    yum clean all
