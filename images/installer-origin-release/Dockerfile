# This Dockerfile is used to create the image quay.io/coreos/origin-release:bazel used in CI
# to build smoke tests.
# EX: docker build --rm -t this-image:latest .

FROM docker.io/openshift/origin-release:golang-1.10

RUN wget -P /etc/yum.repos.d https://copr.fedorainfracloud.org/coprs/vbatts/bazel/repo/epel-7/vbatts-bazel-epel-7.repo

RUN yum install -y bazel gcc-c++ libtool && \
    yum clean all

WORKDIR /go/src/github.com/openshift/origin
