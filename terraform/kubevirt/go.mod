module github.com/openshift/installer/terraform/kubevirt

go 1.16

require github.com/kubevirt/terraform-provider-kubevirt v0.0.0-00010101000000-000000000000

replace (
	github.com/hashicorp/terraform => github.com/openshift/terraform v0.12.20-openshift-4
	github.com/hashicorp/terraform-plugin-sdk => github.com/openshift/hashicorp-terraform-plugin-sdk v1.14.0-openshift
	github.com/kubevirt/terraform-provider-kubevirt => github.com/nirarg/terraform-provider-kubevirt v0.0.0-20201222125919-101cee051ed3
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)
