module github.com/openshift/installer

go 1.12

require (
	cloud.google.com/go/bigtable v1.0.0 // indirect
	cloud.google.com/go/pubsub v1.0.1 // indirect
	github.com/Azure/azure-sdk-for-go v32.5.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.2.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect

	github.com/apparentlymart/go-cidr v1.0.1
	github.com/awalterschulze/gographviz v0.0.0-20170410065617-c84395e536e1
	github.com/aws/aws-sdk-go v1.25.3
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/ignition v0.33.0
	github.com/coreos/ignition/v2 v2.0.1
	github.com/dmacvicar/terraform-provider-libvirt v0.6.0
	github.com/dustinkirkland/golang-petname v0.0.0-20190613200456-11339a705ed2 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/golang/mock v1.3.1
	github.com/gophercloud/gophercloud v0.4.1-0.20190930034851-863d5406e68f
	github.com/gophercloud/utils v0.0.0-20190527093828-25f1b77b8c03
	github.com/hashicorp/go-plugin v1.0.1-0.20190610192547-a1bc61569a26
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/terraform v0.12.12
	github.com/kr/fs v0.1.0 // indirect
	github.com/libvirt/libvirt-go v5.0.0+incompatible
	github.com/metal3-io/baremetal-operator v0.0.0-00010101000000-000000000000
	github.com/metal3-io/cluster-api-provider-baremetal v0.0.0-20191010235856-134c3b78ec63
	github.com/mitchellh/cli v1.0.0
	github.com/openshift-metal3/terraform-provider-ironic v0.1.7
	github.com/openshift/api v3.9.1-0.20191018132714-d0b31d707c46+incompatible
	github.com/openshift/client-go v0.0.0-20191001081553-3b0e988f8cb0
	github.com/openshift/cloud-credential-operator v0.0.0-20190905120421-44ed18ef8496
	github.com/openshift/cluster-api v0.0.0-20191004085540-83f32d3e7070
	github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20190826205919-0cd5daa07e0d
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20190613141010-ecea5317a4ab
	github.com/openshift/library-go v0.0.0-20191003152030-97c62d8a2901
	github.com/openshift/machine-config-operator v0.0.0-00010101000000-000000000000
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.10.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/terraform-providers/terraform-provider-google v1.20.0 // indirect
	github.com/terraform-providers/terraform-provider-google/v2 v2.8.0
	github.com/terraform-providers/terraform-provider-ignition v1.0.1
	github.com/terraform-providers/terraform-provider-local v1.2.1
	github.com/terraform-providers/terraform-provider-openstack v1.18.1-0.20190515162737-b1406b8e4894
	github.com/terraform-providers/terraform-provider-random/v2 v2.1.1
	github.com/vincent-petithory/dataurl v0.0.0-20160330182126-9a301d65acbb
	github.com/vrutkovs/terraform-provider-aws/v3 v3.0.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191024172528-b4ff53e7a1cb
	golang.org/x/tools v0.0.0-20191024191802-c825e86a855b // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/api v0.9.0
	gopkg.in/AlecAivazis/survey.v1 v1.8.7
	gopkg.in/ini.v1 v1.42.0
	k8s.io/api v0.0.0-20191016225839-816a9b7df678
	k8s.io/apimachinery v0.0.0-20191017185446-6e68a40eebf9
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d // indirect
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4
	sigs.k8s.io/cluster-api-provider-aws v0.0.0-00010101000000-000000000000
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.3.0 // indirect
)

replace (
	github.com/Azure/go-autorest v10.15.4+incompatible => github.com/Azure/go-autorest/autorest v0.9.2
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20191001171423-cd2cdd14084a
	github.com/mitchellh/packer => github.com/hashicorp/packer v1.3.5
	github.com/openshift/machine-config-operator => github.com/vrutkovs/machine-config-operator v0.0.0-20191021113908-b6af01302153
	github.com/terraform-providers/terraform-provider-google/v2 => github.com/vrutkovs/terraform-provider-google/v2 v2.8.0
	github.com/terraform-providers/terraform-provider-ignition => github.com/vrutkovs/terraform-provider-ignition v1.0.2-0.20190819094334-ac54201ee306
	github.com/terraform-providers/terraform-provider-random/v2 => github.com/vrutkovs/terraform-provider-random/v2 v2.1.1
	k8s.io/api => k8s.io/api v0.0.0-20190904195148-bacad065d7c3
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190831152136-93cd198ca677
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190904195533-1592ba1f99b8
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20190619152724-cf06d47b6cee
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20190925224209-945cf044115f
)
