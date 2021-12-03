module github.com/openshift/installer/terraform/aws

go 1.16

require github.com/terraform-providers/terraform-provider-aws v1.60.0

replace github.com/terraform-providers/terraform-provider-aws => github.com/openshift/terraform-provider-aws v1.60.1-0.20210622193531-7d13cfbb1a8c // Pin to openshift fork with tag v2.67.0-openshift-1
