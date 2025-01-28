# FIPS support is offered via the baremetal-installer image

FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.23-openshift-4.19 AS macbuilder
ENV GO_COMPLIANCE_EXCLUDE=".*"
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.23-openshift-4.19 AS macarmbuilder
ENV GO_COMPLIANCE_EXCLUDE=".*"
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.23-openshift-4.19 AS linuxbuilder
ENV GO_COMPLIANCE_EXCLUDE=".*"
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.23-openshift-4.19 AS linuxarmbuilder
ENV GO_COMPLIANCE_EXCLUDE=".*"
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make -C terraform

FROM registry.ci.openshift.org/ocp/4.16:base-rhel9
WORKDIR /go/src/github.com/openshift/installer
COPY --from=macbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=macarmbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=linuxbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=linuxarmbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
