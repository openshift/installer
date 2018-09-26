DO NOT MERGE

[![Build Status](https://travis-ci.org/openshift/installer.svg?branch=master)](https://travis-ci.org/openshift/installer)

# Openshift Installer

The CoreOS and OpenShift teams are now working together to integrate Tectonic and OpenShift into a converged platform.
See the CoreOS blog for any additional details:
https://coreos.com/blog/coreos-tech-to-combine-with-red-hat-openshift

## Hacking

These instructions can be used for AWS:

1. Set you access-key and secret in `~/.aws/credentials`.
    You can create credentials in [the IAM console][aws-iam-console], as documented [here][aws-cli-config] and [here][aws-cli-config-files].

2. Build the project
    ```sh
    bazel build tarball
    ```

    *Note*: the project can optionally be built without installing Bazel, provided Podman is installed:
    ```sh
    podman run --rm -v $PWD:$PWD:Z -w $PWD quay.io/coreos/tectonic-builder:bazel-v0.3 bazel --output_base=.cache build tarball
    ```

3. Extract the tarball
    ```sh
    tar -zxf bazel-bin/tectonic-dev.tar.gz
    ```

4. Create an alias for tectonic
    ```sh
    alias tectonic="${PWD}/tectonic-dev/installer/tectonic"
    ```

5. Edit Tectonic configuration file including the $CLUSTER_NAME
    ```sh
    $EDITOR examples/aws.yaml
    ```

6. Prepare a local configuration.
    The structure behind the YAML input is described [here][godoc-InstallConfig].
    ```sh
    tectonic init --config=examples/aws.yaml
    ```

7. Install Tectonic cluster
    ```sh
    tectonic install --dir=$CLUSTER_NAME
    ```

8. Visit `https://{$CLUSTER_NAME}-api.${BASE_DOMAIN}:6443/console/`.
    You may need to ignore a certificate warning if you did not configure a CA known to your browser.
    Log in with the admin credentials you configured in `aws.yaml`.

9. Teardown Tectonic cluster
    ```sh
    tectonic destroy --dir=$CLUSTER_NAME
    ```

## Managing Dependencies
### Go

We follow a hard flattening approach; i.e. direct and inherited dependencies are installed in the base `vendor/`.

Dependencies are managed with [glide](https://glide.sh/) but committed directly to the repository. If you don't have glide, install the latest release from [https://glide.sh/](https://glide.sh/). We require version 0.12 at a minimum.

The vendor directory is pruned using [glide-vc](https://github.com/sgotti/glide-vc). Follow the [installation instructions](https://github.com/sgotti/glide-vc#install) in the project's README.

To add a new dependency:
- Edit the `glide.yaml` file to add your dependency.
- Ensure you add a `version` field for the sha or tag you want to pin to.
- Revendor the dependencies:

```sh
rm glide.lock
glide install --strip-vendor
glide-vc --use-lock-file --no-tests --only-code
bazel run //:gazelle
```

If it worked correctly it should:
- Clone your new dep to the `/vendor` dir and check out the ref you specified.
- Update `glide.lock` to include your new package, add any transitive dependencies and update its hash.
- Regenerate BUILD.bazel files.

For the sake of your fellow reviewers, commit vendored code separately from any other changes.

## Tests

See [tests/README.md](tests/README.md).

[aws-cli-config]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-quick-configuration
[aws-cli-config-files]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html
[aws-iam-console]: https://console.aws.amazon.com/iam/home#/users
[godoc-InstallConfig]: https://godoc.org/github.com/openshift/installer/pkg/types#InstallConfig
