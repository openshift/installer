FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.14 AS macbuilder
ARG TAGS=""
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN GOOS=darwin GOARCH=amd64 DEFAULT_ARCH="$(go env GOHOSTARCH)" make -C terraform all
RUN find /go/src/github.com/openshift/installer/terraform/bin -executable ! -name "terraform" -type f -exec rm -f {} \;

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.14 AS macarmbuilder
ARG TAGS=""
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 DEFAULT_ARCH="$(go env GOHOSTARCH)" make -C terraform all
RUN find /go/src/github.com/openshift/installer/terraform/bin -executable ! -name "terraform" -type f -exec rm -f {} \;

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.14 AS linuxbuilder
ARG TAGS=""
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 DEFAULT_ARCH="$(go env GOHOSTARCH)" make -C terraform all
RUN find /go/src/github.com/openshift/installer/terraform/bin -executable ! -name "terraform" -type f -exec rm -f {} \;

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.14 AS linuxarmbuilder
ARG TAGS=""
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 DEFAULT_ARCH="$(go env GOHOSTARCH)" make -C terraform all
RUN find /go/src/github.com/openshift/installer/terraform/bin -executable ! -name "terraform" -type f -exec rm -f {} \;

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.14 AS builder
WORKDIR /go/src/github.com/openshift/installer
#COPY . .
COPY --from=macbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=macarmbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=linuxbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=linuxarmbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
