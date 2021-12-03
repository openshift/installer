module github.com/openshift/installer/terraform/vsphere

go 1.16

require github.com/hashicorp/terraform-provider-vsphere v1.24.3

replace github.com/hashicorp/terraform-provider-vsphere => github.com/openshift/terraform-provider-vsphere v1.24.3-openshift
