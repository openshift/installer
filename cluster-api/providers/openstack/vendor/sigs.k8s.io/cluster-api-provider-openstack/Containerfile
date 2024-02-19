# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Containerfile for podman
#
# Ideally we would just use Dockerfile to avoid maintaining 2 files. However,
# the Dockerfile is currently using some features [1] that podman doesn't
# support.  We should aim to re-merge these files eventually. In the meantime
# we need to keep them consistent.
#
# [1] * RUN --mount=type=...
#     * podman can't pull dockerfile:1.1-experimental: unsupported docker v2s2
#         media type: "application/vnd.oci.image.layer.v1.tar+gzip"

# Build the manager binary
FROM golang:1.19.6 as builder
WORKDIR /workspace

# Run this with docker build --build_arg goproxy=$(go env GOPROXY) to override the goproxy
ARG goproxy=https://proxy.golang.org
ENV GOPROXY=$goproxy

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the sources
COPY ./ ./

# Build
ARG package=.
ARG ARCH
ARG ldflags

# Do not force rebuild of up-to-date packages (do not use -a) and use the compiler cache folder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -ldflags "${ldflags} -extldflags '-static'" \
    -o manager ${package}

# Production image
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
# Use uid of nonroot user (65532) because kubernetes expects numeric user when applying pod security policies
USER 65532
ENTRYPOINT ["/manager"]
