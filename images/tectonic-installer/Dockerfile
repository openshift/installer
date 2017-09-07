FROM golang:1.8.3-stretch

ENV TERRAFORM_VERSION="0.10.4"

RUN apt-get update \
    && apt-get install --no-install-recommends -y -q \
    unzip

# Install Terraform
RUN curl https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip | funzip > /usr/local/bin/terraform && chmod +x /usr/local/bin/terraform

ENV PROJECT_DIR /go/src/github.com/coreos/tectonic-installer

ADD Makefile $PROJECT_DIR/
ADD config.tf $PROJECT_DIR/
ADD installer/bin/linux/installer $PROJECT_DIR/installer/bin/linux/
ADD platforms $PROJECT_DIR/platforms
ADD modules $PROJECT_DIR/modules

