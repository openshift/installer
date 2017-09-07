FROM ruby:2.4-slim

ENV TERRAFORM_VERSION="v0.10.4-coreos.1"

COPY ./tests/rspec/Gemfile* /tmp/app/

RUN apt-get update && \
    apt-get install --no-install-recommends -y -q \
    curl make unzip jq awscli ssh openvpn && \
    apt-get clean

RUN cd /tmp/app && bundler install && rm -r /tmp/app

# Install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /bin/kubectl

# Install Terraform
RUN curl -L https://github.com/coreos/terraform/releases/download/${TERRAFORM_VERSION}/linux_amd64.zip | funzip > /usr/local/bin/terraform && chmod +x /usr/local/bin/terraform

# Add user jenkins (uid: 1000), this is needed among others for `ssh` to not
# complain about missing user.
RUN useradd -u 1000 -ms /bin/bash jenkins 
