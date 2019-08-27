module github.com/openshift/installer/pkg/terraform/exec/plugins

go 1.12

require (
	cloud.google.com/go v0.40.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/Azure/azure-sdk-for-go v26.7.0+incompatible
	github.com/Azure/go-autorest v11.9.0+incompatible
	github.com/Unknwon/com v0.0.0-20181010210213-41959bdd855f // indirect
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/dmacvicar/terraform-provider-libvirt v0.6.0
	github.com/hashicorp/go-azure-helpers v0.4.1
	github.com/hashicorp/terraform v0.12.0
	github.com/libvirt/libvirt-go-xml v5.1.0+incompatible // indirect
	github.com/mitchellh/packer v1.3.5 // indirect
	github.com/openshift-metal3/terraform-provider-ironic v0.1.7
	github.com/satori/uuid v1.2.0 // indirect
	github.com/terraform-providers/terraform-provider-azurerm v1.27.1
	github.com/terraform-providers/terraform-provider-google v1.20.0 // indirect
	github.com/terraform-providers/terraform-provider-google/v2 v2.8.0
	github.com/terraform-providers/terraform-provider-ignition v1.0.1
	github.com/terraform-providers/terraform-provider-local v1.2.1
	github.com/terraform-providers/terraform-provider-openstack v1.18.1-0.20190515162737-b1406b8e4894
	github.com/terraform-providers/terraform-provider-random/v2 v2.1.1
	github.com/vrutkovs/terraform-provider-aws/v3 v3.0.0
	google.golang.org/appengine v1.6.1 // indirect
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655 // indirect
)

replace (
	github.com/mitchellh/packer => github.com/hashicorp/packer v1.3.5
	github.com/terraform-providers/terraform-provider-google/v2 => github.com/vrutkovs/terraform-provider-google/v2 v2.8.0
	github.com/terraform-providers/terraform-provider-ignition => github.com/vrutkovs/terraform-provider-ignition v1.0.2-0.20190819094334-ac54201ee306
	github.com/terraform-providers/terraform-provider-random/v2 => github.com/vrutkovs/terraform-provider-random/v2 v2.1.1
)
