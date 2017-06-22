# Scaling Tectonic AWS clusters

To scale worker nodes, adjust `tectonic_worker_count` in `terraform.tfvars`.

Use the `plan` command to check your syntax: 

```
$ terraform plan \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/aws
```

Once you are ready to make the changes live, use `apply`:

```
$ terraform apply \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/aws
```

The new nodes should automatically show up in the Tectonic Console shortly after they boot.
