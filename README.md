# OpenShift Installer

## Supported Platforms

* [AWS](docs/user/aws/README.md)
* [Libvirt with KVM](docs/dev/libvirt-howto.md) (development only)
* [OpenStack (experimental)](docs/user/openstack/README.md)

## Quick Start

First, install all [build dependencies](docs/dev/dependencies.md).

Clone this repository to `src/github.com/openshift/installer` in your [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable). Then build the `openshift-install` binary with:

```sh
hack/build.sh
```

This will create `bin/openshift-install`. This binary can then be invoked to create an OpenShift cluster, like so:

```sh
bin/openshift-install create cluster
```

The installer will show a series of prompts for user-specific information and use reasonable defaults for everything else.
In non-interactive contexts, prompts can be bypassed by [providing an `install-config.yaml`](docs/user/overview.md#multiple-invocations).

If you have trouble, refer to [the troubleshooting guide](docs/user/troubleshooting.md).

### Connect to the cluster

Details for connecting to your new cluster are printed by the `openshift-install` binary upon completion, and are also available in the `.openshift_install.log` file.

Example output:

```sh
INFO Waiting 10m0s for the openshift-console route to be created...
INFO Install complete!
INFO Run 'export KUBECONFIG=/path/to/auth/kubeconfig' to manage the cluster with 'oc', the OpenShift CLI.
INFO The cluster is ready when 'oc login -u kubeadmin -p 5char-5char-5char-5char' succeeds (wait a few minutes).
INFO Access the OpenShift web-console here: https://console-openshift-console.apps.${CLUSTER_NAME}.${BASE_DOMAIN}:6443
INFO Login to the console with user: kubeadmin, password: 5char-5char-5char-5char
```

### Cleanup

Destroy the cluster and release associated resources with:

```sh
openshift-install destroy cluster
```

Note that you almost certainly also want to clean up the installer state files too, including `auth/`, `terraform.tfstate`, etc.
The best thing to do is always pass the `--dir` argument to `install` and `destroy`.
And if you want to reinstall from scratch, `rm -rf` the asset directory beforehand.


### Troubleshouting installation

```Cleanup``` step *should* tear down everything created during the previous installation, however sometimes the installer may report error messages like the ones below:

``` 
ERROR Error: Error applying plan:
ERROR
ERROR 3 errors occurred:
ERROR     * module.bootstrap.aws_iam_instance_profile.bootstrap: 1 error occurred:
ERROR     * aws_iam_instance_profile.bootstrap: Error creating IAM instance profile test-bootstrap-profile: 

FATAL failed to fetch Cluster: failed to generate asset "Cluster": failed to create cluster: failed to apply using Terraform 
``` 
This happens when the installer tries to create an IAM instance profile with a name that matches an existing instance profile. This usually means that IAM instance profiles from a previous installation by the user didn't get deleted during an earlier cluster teardown.

To fix this problem it is necessary to delete conflicting IAM profiles and retry. 
To delete problematic IAM profiles follow [list-instance-profiles](https://docs.aws.amazon.com/cli/latest/reference/iam/list-instance-profiles.html) and [delete-instance-profile](https://docs.aws.amazon.com/cli/latest/reference/iam/delete-instance-profile.html) which can be summarized to 
```
$ aws iam list-instance-profiles | grep USER
$ aws iam delete-instance-profile --instance-profile-name PROFILE_NAME
``` 
