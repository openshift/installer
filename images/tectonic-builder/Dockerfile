###
# Update the builder image manually by running:
# (See README.md/##Upstream-and-CoreOS-Terraform)
#
# docker build -t quay.io/coreos/tectonic-builder:<next-semver> -f images/builder/Dockerfile .
# docker push quay.io/coreos/tectonic-builder:<next-semver>
#
# docker build \
# --build-arg TERRAFORM_URL=<upstream terraform download url> \
# -t quay.io/coreos/tectonic-builder:<next-semver>-upstream-terraform \
# -f images/builder/Dockerfile .
# docker push quay.io/coreos/tectonic-builder:<next-semver>-upstream-terraform
###

FROM golang:1.9.2-stretch

### For golang testing stuff
RUN go get -u github.com/golang/lint/golint
RUN go get github.com/jstemmer/go-junit-report

### Tools used by 'make structure-check'
RUN go get github.com/segmentio/terraform-docs
RUN go get github.com/coreos/tectonic-installer/contrib/terraform-examples && \
    rm -rf /go/src/github.com/coreos/tectonic-installer/
RUN go get github.com/bronze1man/yaml2json

### License parser
RUN go get github.com/coreos/license-bill-of-materials

### 'grafiti' for cluster cleanup
ENV GRAFITI_VERSION "v0.1.1"
RUN git clone -q https://github.com/coreos/grafiti.git ${GOPATH}/src/github.com/coreos/grafiti \
    && cd ${GOPATH}/src/github.com/coreos/grafiti \
    && git checkout -q tags/${GRAFITI_VERSION} \
    && make install

# /go needs to be writable by jenkins user like it is in the upstream golang image
RUN chmod 777 -R /go

### Install Shellcheck and Terraform
ENV SHELLCHECK_VERSION="v0.4.6"
ENV TERRAFORM_VERSION="0.11.1"
ENV HOME /opt/home
RUN mkdir -p ${HOME} && \
    chmod 777 -R ${HOME} && \
    apt-get update && \
    apt-get install --no-install-recommends -y -q \
    build-essential sudo curl wget git autoconf automake unzip libtool jq awscli gnupg1 \
    openvpn xauth

# Install Terraform
# TERRAFORM_URL enables us to build the upstream-terraform Tectonic builder
# image (See README.md/##Upstream-and-CoreOS-Terraform)
ARG TERRAFORM_URL=https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip
RUN curl -L ${TERRAFORM_URL} | funzip > /usr/local/bin/terraform && chmod +x /usr/local/bin/terraform

# Install Shellcheck
RUN cd /tmp && \
    wget --quiet https://storage.googleapis.com/shellcheck/shellcheck-${SHELLCHECK_VERSION}.linux.x86_64.tar.xz && \
    tar xJf shellcheck-${SHELLCHECK_VERSION}.linux.x86_64.tar.xz && \
    mv /tmp/shellcheck-${SHELLCHECK_VERSION}/shellcheck /usr/local/bin/shellcheck
