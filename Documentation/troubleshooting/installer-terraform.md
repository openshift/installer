# Troubleshooting Tectonic Installer

## Re-applying over a partially failed run

Re-applying over a failed run is not supported until the installer migrates to Terraform 0.9. If you see the following output:

```
Error applying plan:

3 error(s) occurred:

* ignition_config.master: invalid file "8f8b62c46063e64e3709a2be1317bd83562abf43c3866d7adae023477b57fb47", unknown file id
* ignition_config.worker: invalid file "8f8b62c46063e64e3709a2be1317bd83562abf43c3866d7adae023477b57fb47", unknown file id
* ignition_config.etcd: invalid systemd unit "297994451bbc66a0b227eea50017c4d05fde37b99130d1efc2f6b14b3bd2a283", unknown systemd unit id
```

Destroy your cluster and start over.

## Failed tectonic.service unit

If you see a failed `tectonic.service` unit, you will need to manually copy a file over to complete the install:

```
$ scp ./kubeconfig core@<your api public ip>:/home/core
```

Then SSH to the box:

```
$ ssh core@<your api public ip>
```

And execute this script to finish the install:

```
$ sudo bash /opt/tectonic/tectonic.sh kubeconfig /opt/tectonic/tectonic
```

## Missing Terraform modules

If you see errors about missing modules:

```
Error configuring: 2 error(s) occurred:

* module.bootkube: provider localfile couldn't be found
* module.tectonic: provider localfile couldn't be found
```

These errors indicate that you are not using the customized Terraform binary bundled with the project. You can check this:

```
$ which terraform
/Users/coreos/tectonic-installer/bin/terraform/terraform
```

## Invalid or unknown key: tags

```
$ terraform plan -var-file=terraform.tfvars platforms/aws
2 error(s) occurred:

* module.masters.aws_autoscaling_group.masters: : invalid or unknown key: tags
* module.workers.aws_autoscaling_group.workers: : invalid or unknown key: tags
```

This error indicates that the `.terraformrc` is not being used. Be sure you have exported your config environment variable:

```
export TERRAFORM_CONFIG=$(pwd)/.terraformrc
```

Afterwards, you should be able to execute the desired action.
