FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.15 AS macbuilder
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.15 AS macarmbuilder
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.15 AS linuxbuilder
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.20-openshift-4.15 AS linuxarmbuilder
WORKDIR /go/src/github.com/openshift/installer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make -C terraform

FROM registry.ci.openshift.org/ocp/builder:base
WORKDIR /go/src/github.com/openshift/installer
COPY --from=macbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=macarmbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=linuxbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
COPY --from=linuxarmbuilder /go/src/github.com/openshift/installer/terraform/bin/ terraform/bin/
