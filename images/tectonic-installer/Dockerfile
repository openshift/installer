FROM alpine:3.7

# Docker build does not follow symlinks (see
# https://github.com/moby/moby/issues/1676) Make sure to first copy the tarball
# from bazel-bin/tectonic.tar.gz to the docker build context folder of your
# choice (e.g. the root of the repository).
ADD tectonic-dev.tar.gz /
