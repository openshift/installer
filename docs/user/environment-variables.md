# Environment Variables

The installer accepts a number of environment variable that allow the interactive prompts to be bypassed. Setting any of the following environment variables to their corresponding value, will cause the installer to use that value instead of prompting.

## General

* `OPENSHIFT_INSTALL_BASE_DOMAIN`:
    The base domain of the cluster. All DNS records will be sub-domains of this base.

    For AWS, this must be a previously-existing public Route 53 zone.  You can check for any already in your account with:

    ```sh
    aws route53 list-hosted-zones --query 'HostedZones[? !(Config.PrivateZone)].Name' --output text
    ```

* `OPENSHIFT_INSTALL_CLUSTER_NAME`:
     The name of the cluster.
     This will be used when generating sub-domains.
* `OPENSHIFT_INSTALL_EMAIL_ADDRESS`:
     The email address of the cluster administrator.
     This will be used to log in to the console.
* `OPENSHIFT_INSTALL_PASSWORD`:
     The password of the cluster administrator.
     This will be used to log in to the console.
* `OPENSHIFT_INSTALL_PLATFORM`:
     The platform onto which the cluster will be installed.
     Valid values are `aws` and `libvirt`.
* `OPENSHIFT_INSTALL_PULL_SECRET`:
     The container registry pull secret for this cluster (e.g. `{"auths": {...}}`).
     You can generate these secrets with the `podman login` command.
* `OPENSHIFT_INSTALL_PULL_SECRET_PATH`:
     As an alternative to `OPENSHIFT_INSTALL_SSH_PUB_KEY`, you can configure this variable with a path containing your pull secret.
* `OPENSHIFT_INSTALL_SSH_PUB_KEY`:
     The SSH public key used to access all nodes within the cluster (e.g. `ssh-rsa AAAA...`).
     This is optional.
* `OPENSHIFT_INSTALL_SSH_PUB_KEY_PATH`:
     As an alternative to `OPENSHIFT_INSTALL_SSH_PUB_KEY`, you can configure this variable with a path containing your SSH public key (e.g. `~/.ssh/id_rsa.pub`).

## Platform-Specific

* `AWS_PROFILE`:
     The AWS profile that corresponds to value in `${HOME}/.aws/credentials`.  If not provided, the default is "default".
* `OPENSHIFT_INSTALL_AWS_REGION`:
    The AWS region to be used for installation.
* `OPENSHIFT_INSTALL_LIBVIRT_URI`:
    The libvirt connection URI to be used.
    This must be accessible from the running cluster.
* `OPENSHIFT_INSTALL_LIBVIRT_IMAGE`:
    The URI for the OS image.
    For example it might be url like `http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz` or
    a local file like `file:///tmp/openshift-install-853528428`.
