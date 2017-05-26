# Scaling Tectonic AWS clusters

To scale worker nodes, adjust `tectonic_worker_count` in `terraform.vars` and run:

```
$ terraform apply $ terraform plan \
  -var-file=build/${CLUSTER}/terraform.tfvars \
  -target module.workers \
  platforms/aws
```
