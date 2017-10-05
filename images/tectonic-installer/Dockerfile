FROM golang:1.9.1-stretch

ENV TERRAFORM_VERSION="0.10.7"

RUN apt-get update \
    && apt-get install --no-install-recommends -y -q \
    unzip

# Install Terraform
RUN curl https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip | funzip > /usr/local/bin/terraform && chmod +x /usr/local/bin/terraform

ENV PROJECT_DIR /go/src/github.com/coreos/tectonic-installer

ADD . $PROJECT_DIR/
