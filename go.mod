module github.com/openshift/installer

go 1.12

require (
	cloud.google.com/go/bigtable v1.0.0 // indirect
	cloud.google.com/go/pubsub v1.0.1 // indirect
	github.com/Azure/azure-sdk-for-go v33.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.0
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/awalterschulze/gographviz v0.0.0-20170410065617-c84395e536e1
	github.com/aws/aws-sdk-go v1.25.3
	github.com/containers/image v2.0.0+incompatible
	github.com/coreos/ignition v0.26.0
	github.com/dustinkirkland/golang-petname v0.0.0-20190613200456-11339a705ed2 // indirect
	github.com/emicklei/go-restful v2.10.0+incompatible // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.3 // indirect
	github.com/golang/mock v1.3.1
	github.com/gophercloud/gophercloud v0.4.1-0.20190930034851-863d5406e68f
	github.com/gophercloud/utils v0.0.0-20190527093828-25f1b77b8c03
	github.com/hashicorp/terraform v0.12.12 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/libvirt/libvirt-go v5.0.0+incompatible
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/metal3-io/baremetal-operator v0.0.0-00010101000000-000000000000
	github.com/metal3-io/cluster-api-provider-baremetal v0.0.0-20191010235856-134c3b78ec63
	github.com/openshift/api v3.9.1-0.20190806225813-d2972510af76+incompatible
	github.com/openshift/client-go v0.0.0-20190806162413-e9678e3b850d
	github.com/openshift/cloud-credential-operator v0.0.0-20190905120421-44ed18ef8496
	github.com/openshift/cluster-api v0.0.0-20190619113136-046d74a3bd91
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20190702211226-53df0c29f8e2 // indirect
	github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20190826205919-0cd5daa07e0d
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20190613141010-ecea5317a4ab
	github.com/openshift/installer/pkg/terraform/exec v0.0.0-00010101000000-000000000000
	github.com/openshift/installer/pkg/terraform/exec/plugins v0.0.0-00010101000000-000000000000
	github.com/openshift/library-go v0.0.0-20190129125304-4b9f6ceb6598
	github.com/openshift/machine-config-operator v3.11.0+incompatible // indirect
	github.com/pborman/uuid v1.2.0
	github.com/petar/GoLLRB v0.0.0-20130427215148-53be0d36a84c // indirect
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.10.0
	github.com/shurcooL/vfsgen v0.0.0-20181020040650-a97a25d856ca // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/terraform-providers/terraform-provider-aws v1.60.0 // indirect
	github.com/terraform-providers/terraform-provider-azurerm v1.35.0 // indirect
	github.com/terraform-providers/terraform-provider-local v1.4.0 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20160330182126-9a301d65acbb
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190911201528-7ad0cfa0b7b5
	google.golang.org/api v0.9.0
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	gopkg.in/ini.v1 v1.42.0
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d // indirect
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5
	sigs.k8s.io/cluster-api-provider-aws v0.0.0-00010101000000-000000000000
	sigs.k8s.io/cluster-api-provider-azure v0.0.0-00010101000000-000000000000
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.2.0-beta.2 // indirect
)

replace (
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20191001171423-cd2cdd14084a
	github.com/mitchellh/packer => github.com/hashicorp/packer v1.3.5
	github.com/openshift/installer/pkg/terraform/exec => ./pkg/terraform/exec
	github.com/openshift/installer/pkg/terraform/exec/plugins => ./pkg/terraform/exec/plugins
	github.com/terraform-providers/terraform-provider-google/v2 => github.com/vrutkovs/terraform-provider-google/v2 v2.8.0
	github.com/terraform-providers/terraform-provider-ignition => github.com/abhinavdahiya/terraform-provider-ignition v1.0.2-0.20190513232748-18ce0b36dae1
	github.com/terraform-providers/terraform-provider-random/v2 => github.com/vrutkovs/terraform-provider-random/v2 v2.1.1
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20190619152724-cf06d47b6cee
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20190718103506-6a50a8c59d8a
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20190925224209-945cf044115f
)
