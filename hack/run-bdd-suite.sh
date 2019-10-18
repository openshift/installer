#! /usr/bin/sh
# Example: ./hack/run-bdd-suite.sh tests/bdd-smoke/suites/config/libvirt

if [ "$IS_CONTAINER" != "" ]; then
  mkdir ~/.ssh
  echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCc4amSo5Z59mH4bUgZ4M8A1sURg+qUKkYjZ9+0fft5/OmSbe34Ge9gUrJ8UKNpmbC+W26v9qKX21NCYct1QOgIlVqyIjZUll0NfcJXToAd1o2GrItMlpQEpXfQGNOjZtWu1uyO187qbGMwnq8e8EH9R1IJZwTI6sOVnc8sSyzydQ==" > ~/.ssh/id_rsa.pub
  
  TAGS=libvirt ./hack/build.sh
  GOPATH_ORI=$GOPATH
  export GOPATH=$GOPATH:/go/src/github.com/openshift/installer/tests/bdd-smoke/vendor
  go build -o "bin/ginkgo" github.com/onsi/ginkgo/ginkgo
  export GOPATH=$GOPATH_ORI
  ./bin/ginkgo "${@}"
else
  podman run --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/installer:z" \
    --workdir /go/src/github.com/openshift/installer \
    docker.io/openshift/origin-release:golang-1.12 \
    ./hack/run-bdd-suite.sh "${@}"
fi;