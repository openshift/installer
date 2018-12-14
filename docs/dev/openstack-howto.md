# OpenStack HOWTO

It's possible to launch an all-in-one - for test and development - cluster on OpenStack.

## The ocp-doit scripts

We've created a set of scripts that make it easier to setup an all-in-one OpenStack. The OpenStack cluster should work either on a VM or a BM node. We recommend the latter.

The scripts can be found in the following [repo](https://github.com/imain/ocp-doit)


### Configure ocp-doit

Copy the `config_example.sh` file and configure it to match your environment.

### Install OpenStack 

Run the scripts in the following order to install the dependencies, deploy, and configure OpenStack:

```console
$ ./01_install_requirements.sh
$ ./02_run_all_in_one.sh
$ ./03_configure_undercloud.sh
```

There should be a `$HOME/.config/openstack/clouds,yaml` file once the above scripts have completed. This file contains the openstack credentials. It should be possible to query the openstack cloud at this point:

```console
$ export OS_CLOUD=openshift
$ openstack server list
```

### Clone and build the project

```console
$ ./04_ocp_repo_sync.sh
$ ./05_build_ocp_installer.sh
```

### Check the OCP env configs

The `openshift-install` specific configs are stored in the `ocp_install_env.sh` script. Modify this script accordingly to match the newly installed OpenStack cloud.

There shouldn't be need to change any of the options, except for the base domain.

## Cleanup

The following script cleans the OCP deployment. It removes the VMs created by `openshift-install` and all the other resources created for this cluster.

```console
$ ./ocp_cleanup.sh
```

To remove everything, including OpenStack, use the following script:

```console
$ ./tripleo_cleanup.sh
```

### SSH access

For this deployment, we're trying to not depend on Floating IPS. It is possible to create and assign a floating ip to the `$cluster-api` vm but it's not required for a successful deployment. There's a script that would help to ssh into the VMs using `netns`:

```console
$ ./s.sh master-0
```

### Inspect the cluster with kubectl

You'll need a `kubectl` binary on your path and [the kubeconfig from your `cluster` call](../../README.md#kubeconfig).

The following will allow you to query the cluster from the node you ran `openshift-install` from. Note that this will likely require you to have a floating IP set on your api VM.

```sh
export KUBECONFIG="${DIR}/auth/kubeconfig"
kubectl get --all-namespaces pods
```

Alternatively, you can run `kubectl` from the bootstrap or master nodes.
Use `scp` or similar to transfer your local `${DIR}/auth/kubeconfig`, then [SSH in](#ssh-access) and run:

```sh
export KUBECONFIG=/where/you/put/your/kubeconfig
kubectl get --all-namespaces pods
```

## FAQ

### Github Issue Tracker
You might find other reports of your problem in the [Issues tab for this repository][issues_openstack] where we ask you to provide any additional information.
If your issue is not reported, please do.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
