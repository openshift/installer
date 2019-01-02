# This Dockerfile is a used by CI to publish openshift/origin-v4.0:installer-libvirt
# It builds an image containing openshift-install with libvirt support, 
# which is required for bring-your-own RHCOS

FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN TAGS=libvirt hack/build.sh


FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /go/src/github.com/openshift/installer/bin/openshift-install /bin/openshift-install
RUN yum install -y libvirt-libs && yum clean all
RUN mkdir /output && chown 1000:1000 /output
USER 1000:1000
ENV PATH /bin
ENV HOME /output
WORKDIR /output
ENTRYPOINT ["/bin/openshift-install"]
