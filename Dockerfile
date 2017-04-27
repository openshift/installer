FROM golang:alpine

ENV TERRAFORM_VERSION=0.9.4

RUN apk add --update git bash make

RUN go get github.com/coreos/bcrypt-tool

WORKDIR $GOPATH/src/github.com/hashicorp/terraform
RUN git clone https://github.com/hashicorp/terraform.git ./ && \
    git checkout v${TERRAFORM_VERSION} && \
    go run scripts/generate-plugins.go && \
    XC_ARCH=amd64 XC_OS=linux ./scripts/build.sh

VOLUME /terraform
WORKDIR /terraform

ENTRYPOINT ["/bin/bash", "-c"]
CMD "bash"