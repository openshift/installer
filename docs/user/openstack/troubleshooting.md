# OpenShift 4 installer on OpenStack troubleshooting

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

This document discusses some troubleshooting options for OpenStack based
deployments. For general tips on troubleshooting the installer, see the [Installer Troubleshooting](../troubleshooting.md) guide.

## View instances logs

OpenStack CLI tools should be installed, then:

`openstack console log show <instance>`

## Machine has ERROR state

This could be because the machine's instance was accidentally destroyed and the cluster API provider cannot recreate it.

You can check the status of machines with the help of the command

```sh
oc get machines -n openshift-machine-api
```

If the broken machine is a master then follow the instructions in the [disaster recovery documentation](https://docs.openshift.com/container-platform/4.1/disaster_recovery/scenario-1-infra-recovery.html).

For workers, you should delete the machine manually with

```sh
oc delete machine -n openshift-machine-api <machine_name>
```

The operation can take up to 5 minutes, during which time the machine will be gracefully removed and all its resources returned to the pool.

A new worker machine for the cluster will soon be created automatically by the [machine-api-operator](https://github.com/openshift/machine-api-operator).

> **Note**
> In future versions of OpenShift all broken machines will be automatically deleted and recovered by the machine-api-operator.

## SSH access to the instances

Get the IP address of the node on the private network:

```sh
openstack server list | grep master
| 0dcd756b-ad80-42f1-987a-1451b1ae95ba | cluster-wbzrr-master-1     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.21                | rhcos           | m1.s2.xlarge |
| 3b455e43-729b-4e64-b3bd-1d4da9996f27 | cluster-wbzrr-master-2     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.18                | rhcos           | m1.s2.xlarge |
| 775898c3-ecc2-41a4-b98b-a4cd5ae56fd0 | cluster-wbzrr-master-0     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.12                | rhcos           | m1.s2.xlarge |
```

And connect to it using the master currently holding the API VIP (and hence the API FIP) as a jumpbox:

```sh
ssh -J core@${FIP} core@<host>
```

## Bootstrap Flavor Issues

The OpenShift installer allows specifying a custom flavor for the temporary bootstrap machine via the `bootstrapFlavor` field in the install config. When this field is misconfigured, the installation will fail before the cluster is created. This section covers how to diagnose and resolve common bootstrap flavor problems.

For background on the `bootstrapFlavor` field and cost optimization guidance, see the [OpenStack customization documentation](customization.md#bootstrap-flavor-optimization).

### Error: Bootstrap Flavor Not Found

If the specified flavor does not exist in your OpenStack environment, the installer will fail during validation with an error similar to:

```
FATAL Failed to validate install-config: [platform.openstack.bootstrapFlavor: Not found: "my-bootstrap-flavor"]
```

**Resolution:**

1. List available flavors in your OpenStack project:

   ```sh
   openstack flavor list
   ```

2. To include private (project-scoped) flavors that may not appear by default:

   ```sh
   openstack flavor list --private
   openstack flavor list --all
   ```

3. Verify that the flavor name in your `install-config.yaml` exactly matches a flavor in the output. Flavor names are case-sensitive.

4. If the required flavor does not exist, ask your OpenStack administrator to create it or choose an existing flavor that meets the [minimum bootstrap resource requirements](customization.md#minimum-recommended-resources-for-bootstrap).

### Error: Bootstrap Flavor Fails Resource Validation

The bootstrap flavor must meet at least the control plane flavor minimum requirements (4 vCPUs, 16 GB RAM, 100 GB disk). If the specified flavor is too small, validation will fail with an error similar to:

```
FATAL Failed to validate install-config: [platform.openstack.bootstrapFlavor: Invalid value: "my-small-flavor": flavor does not meet minimum resource requirements]
```

**Resolution:**

1. Check the resource specifications of available flavors:

   ```sh
   openstack flavor list --long
   ```

   Or inspect a specific flavor:

   ```sh
   openstack flavor show <flavor-name>
   ```

2. Ensure the chosen bootstrap flavor has at least:
   - **4 vCPUs**
   - **16 GB RAM**
   - **100 GB disk**

3. Update the `bootstrapFlavor` field in your `install-config.yaml` to a flavor that meets these minimums, or remove it entirely to inherit the control plane flavor.

> **Note**
> Baremetal flavors skip resource validation and can be used as bootstrap flavors regardless of their reported resource values.

### Error: Insufficient Quota for Bootstrap Instance

Even if the flavor exists and meets resource requirements, OpenStack quota limits may prevent the bootstrap instance from being created. The installation will appear to hang or fail during the infrastructure provisioning phase.

**Symptoms:**

- Installation stalls after printing `Waiting for bootstrap to complete...`
- The bootstrap server is not visible in `openstack server list`
- OpenStack logs or events show quota-exceeded errors

**Resolution:**

1. Check your current quota and usage:

   ```sh
   openstack quota show
   openstack limits show --absolute
   ```

2. Check for quota-related errors in OpenStack compute events:

   ```sh
   openstack server list --all-projects 2>/dev/null | grep bootstrap
   ```

3. Verify there is sufficient quota for the additional vCPUs, RAM, and instances required by the bootstrap machine. Contact your OpenStack administrator to increase quota limits if needed.

4. Consider choosing a smaller (but still qualifying) bootstrap flavor to reduce resource consumption. See [Bootstrap Flavor Optimization](customization.md#bootstrap-flavor-optimization) for guidance on selecting an appropriate flavor.

### Diagnosing Bootstrap Machine Failures

If the bootstrap machine is created but the installation still fails, use the following steps to investigate flavor-related issues.

**Check bootstrap server status:**

```sh
openstack server list | grep bootstrap
```

A healthy bootstrap server should show `ACTIVE` status. An `ERROR` status indicates a provisioning failure, which may be caused by an incompatible flavor.

**View bootstrap server details:**

```sh
openstack server show <bootstrap-server-name>
```

Look at the `fault` and `OS-EXT-STS:task_state` fields for error details. The flavor used is shown in the `flavor` field — verify it matches your intended configuration.

**View bootstrap console logs:**

```sh
openstack console log show <bootstrap-server-name>
```

This can reveal early boot failures caused by insufficient disk space or memory associated with the chosen flavor.

**Check bootstrap machine events:**

```sh
openstack server event list <bootstrap-server-name>
```

This shows a timeline of actions taken on the bootstrap instance, including any scheduling or resource errors.

**SSH to the bootstrap node (if reachable):**

If the bootstrap node is accessible, check for service failures related to insufficient resources:

```sh
ssh -J core@${FIP} core@<bootstrap-ip>
sudo journalctl -b -p err
sudo systemctl --failed
```

### Summary: Common Bootstrap Flavor Issues

| Issue | Error / Symptom | Resolution |
|---|---|---|
| Flavor name misspelled or wrong case | `Not found: "flavor-name"` at validation | Run `openstack flavor list --all` and correct the name |
| Flavor does not exist | `Not found: "flavor-name"` at validation | Create the flavor or choose an existing one |
| Flavor too small | `Invalid value` at validation | Choose a flavor with ≥4 vCPU, 16 GB RAM, 100 GB disk |
| Quota exceeded | Installation hangs; no bootstrap server | Increase quota or choose a smaller qualifying flavor |
| Bootstrap server in ERROR state | `openstack server list` shows ERROR | Check `openstack server show` fault field and console log |

## Cluster destruction if its metadata has been lost

When deploying a cluster, the installer generates metadata in the asset directory that is then used to destroy the cluster. If the metadata were accidentally deleted, the destruction of the cluster terminates with an error

```txt
FATAL Failed while preparing to destroy cluster: open clustername/metadata.json: no such file or directory
```

To avoid this error and successfully destroy the cluster, you need to restore the `metadata.json` file in a temporary asset directory. To do this, you only need to know ID of the cluster you want to destroy.

First, you need to create a temporary directory where the `metadata.json` file will be located. The name and location can be anything, but to avoid possible conflicts, we recommend using `mktemp` command.

```sh
export TMP_DIR=$(mktemp -d -t shiftstack-XXXXXXXXXX)
```

The next step is to restore the `metadata.json` file.

```sh
export CLUSTER_ID=clustername-eiu38 # id of the cluster you want to destroy
echo "{\"infraID\":\"$INFRA_ID\",\"openstack\":{\"identifier\":{\"openshiftClusterID\":\"$INFRA_ID\"}}}" > $TMP_DIR/metadata.json
```

Now you have a working directory and you can destroy the cluster by executing the following command:

```sh
openshift-install destroy cluster --dir $TMP_DIR
```
