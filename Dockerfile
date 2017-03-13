FROM golang:alpine

ENV TERRAFORM_VERSION=0.8.8

RUN apk add --update git bash

ENV TF_DEV=true

WORKDIR $GOPATH/src/github.com/hashicorp/terraform
RUN git clone https://github.com/hashicorp/terraform.git ./ && \
    git checkout v${TERRAFORM_VERSION} && \
    /bin/bash scripts/build.sh

COPY ./plugins $GOPATH/src/github.com/coreos-inc/tectonic-platform-sdk/plugins

RUN go build -o $GOTPAH/bin/terraform-provider-localfile \
  github.com/coreos-inc/tectonic-platform-sdk/plugins/bins/provider-localfile

RUN go build -o $GOPATH/bin/terraform-provider-template \
  github.com/coreos-inc/tectonic-platform-sdk/plugins/bins/provider-template

VOLUME /terraform
WORKDIR /terraform

ENTRYPOINT ["terraform"]
