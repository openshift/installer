FROM debian:stretch

ENV TERRAFORM_VERSION="0.11.1"

RUN apt-get update && \
    apt-get install --no-install-recommends -y -q \
    curl make unzip jq awscli wget xvfb xauth ssh openvpn build-essential zlib1g-dev gnupg uuid-runtime && \
    apt-get clean

# Install bazel
RUN echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" > /etc/apt/sources.list.d/bazel.list && \
    curl https://bazel.build/bazel-release.pub.gpg | apt-key add - && \
    apt-get update && \
    apt-get install --no-install-recommends -y -q openjdk-8-jdk bazel git python python-yaml && \
    apt-get clean

# Install Terraform
RUN curl -L https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip | funzip > /usr/local/bin/terraform && chmod +x /usr/local/bin/terraform

# Add user jenkins (uid: 1000), this is needed among others for `ssh` to not
# complain about missing user.
RUN useradd -u 1000 -ms /bin/bash jenkins
RUN mkdir -p /home/jenkins/.ssh && chown jenkins:jenkins /home/jenkins/.ssh
ENV HOME /home/jenkins
