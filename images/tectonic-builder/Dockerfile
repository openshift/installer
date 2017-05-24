###
# Builds a container used for compiling the tectonic-installer binaries and javascript applications.
#
# Node install inspired by: https://github.com/dockerfile/nodejs/blob/master/Dockerfile
#
# Update the builder image manually by running:
#
# docker build -f builder/Dockerfile .
# docker tag <resulting-image-id> quay.io/coreos/tectonic-builder:<next-semver>
# docker push quay.io/coreos/tectonic-builder:<next-semver>
###

FROM golang:1.8-wheezy

### For golang testing stuff
RUN go get -u github.com/golang/lint/golint
RUN go get github.com/jstemmer/go-junit-report

### License parser
RUN go get github.com/estroz/license-bill-of-materials

### Install Terraform, NodeJS and Yarn
ENV TERRAFORM_VERSION="v0.9.6-fcdf494"
ENV NODE_VERSION="v6.10.0"
ENV YARN_VERSION="v0.21.3"
# yarn needs a home writable by any user running the container
ENV HOME /opt/home
RUN mkdir -p ${HOME}
RUN chmod 777 -R ${HOME}

# Needed to get jq in wheezy
RUN echo "deb http://ftp.debian.org/debian wheezy-backports main" > /etc/apt/sources.list.d/backports.list

RUN apt-get update \
    && apt-get install --no-install-recommends -y -q \
    build-essential sudo curl wget git autoconf automake unzip libtool && \
    apt-get install --no-install-recommends -y -q -t wheezy-backports jq

# install terraform
# RUN curl https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip | funzip > /usr/local/bin/terraform && chmod +x /usr/local/bin/terraform
RUN curl -L -o /tmp/terraform.zip https://github.com/coreos/terraform/releases/download/${TERRAFORM_VERSION}/terraform_linux_amd64.zip && \
    unzip -o /tmp/terraform.zip -d /usr/local/bin/ && \
    chmod +x /usr/local/bin/terraform

RUN cd /tmp && \
    wget --quiet -O /tmp/node.tar.gz http://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-linux-x64.tar.gz && \
    tar xf node.tar.gz && \
    rm -f /tmp/node.tar.gz && \
    cd node-* && \
    cp -r lib/node_modules /usr/local/lib/node_modules && \
    cp bin/node /usr/local/bin && \
    ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm
# so any container user can install global node modules if needed
RUN chmod 777 /usr/local/lib/node_modules

RUN cd /tmp && \
    wget --quiet -O /tmp/yarn.tar.gz https://github.com/yarnpkg/yarn/releases/download/${YARN_VERSION}/yarn-${YARN_VERSION}.tar.gz && \
    tar xf yarn.tar.gz && \
    rm -f /tmp/yarn.tar.gz && \
    mv /tmp/dist /usr/local/yarn && \
    ln -s /usr/local/yarn/bin/yarn /usr/local/bin/yarn

# cleanup
RUN rm -rf /tmp/node-v*
RUN rm -f /tmp/yarn.tar.gz
