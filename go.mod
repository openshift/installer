module github.com/openshift/installer

go 1.14

require (
	cloud.google.com/go v0.65.0
	github.com/Azure/azure-sdk-for-go v43.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.0
	github.com/Azure/go-autorest/autorest/adal v0.8.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.1
	github.com/Azure/go-autorest/autorest/to v0.3.1-0.20191028180845-3492b2aff503
	github.com/Azure/go-ntlmssp v0.0.0-20191115210519-2b2be6cc8ed4 // indirect
	github.com/ChrisTrenkamp/goxpath v0.0.0-20190607011252-c5096ec8773d // indirect
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667 // indirect
	github.com/antchfx/xpath v1.1.2 // indirect
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.32.3
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/c4milo/gotoolkit v0.0.0-20190525173301-67483a18c17a // indirect
	github.com/clarketm/json v1.14.1
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/ignition/v2 v2.3.0
	github.com/dmacvicar/terraform-provider-libvirt v0.6.2
	github.com/equinix/terraform-provider-equinix-metal v1.7.3-0.20201114035505-e4eeff216bf2
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-playground/validator/v10 v10.2.0
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gophercloud/gophercloud v0.12.1-0.20200821143728-362eb785d617
	github.com/gophercloud/utils v0.0.0-20200918191848-da0e919a012a
	github.com/h2non/filetype v1.0.12
	github.com/hashicorp/go-azure-helpers v0.10.0
	github.com/hashicorp/go-plugin v1.2.2
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/terraform v0.12.21
	github.com/hashicorp/terraform-plugin-sdk v1.15.0
	github.com/hashicorp/vault v1.3.0 // indirect
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c // indirect
	github.com/keybase/go-crypto v0.0.0-20190828182435-a05457805304 // indirect
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/libvirt/libvirt-go-xml v5.10.0+incompatible // indirect
	github.com/masterzen/simplexml v0.0.0-20190410153822-31eea3082786 // indirect
	github.com/masterzen/winrm v0.0.0-20190308153735-1d17eaf15943 // indirect
	github.com/metal3-io/baremetal-operator v0.0.0
	github.com/metal3-io/cluster-api-provider-baremetal v0.0.0
	github.com/mitchellh/cli v1.1.1
	github.com/openshift-metal3/terraform-provider-ironic v0.2.3
	github.com/openshift/api v3.9.1-0.20191111211345-a27ff30ebf09+incompatible
	github.com/openshift/client-go v0.0.0-20201020074620-f8fd44879f7c
	github.com/openshift/cloud-credential-operator v0.0.0-20200316201045-d10080b52c9e
	github.com/openshift/cluster-api v0.0.0-20191129101638-b09907ac6668
	github.com/openshift/cluster-api-provider-equinix-metal v0.1.0
	github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20201027164920-70f2f92e64ab
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20191219173431-2336783d4603
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20200504092944-27473ea1ae43
	github.com/openshift/library-go v0.0.0-20201022113156-a4ff9e1d2900
	github.com/openshift/machine-api-operator v0.2.1-0.20201002104344-6abfb5440597
	github.com/openshift/machine-config-operator v0.0.0
	github.com/ovirt/go-ovirt v0.0.0-20200613023950-320a86f9df27
	github.com/ovirt/terraform-provider-ovirt v0.4.3-0.20200914080915-c4444fb5c201
	github.com/packer-community/winrmcp v0.0.0-20180921211025-c76d91c1e7db // indirect
	github.com/packethost/packngo v0.5.1
	github.com/pborman/uuid v1.2.0
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.10.0
	github.com/satori/uuid v1.2.0 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/stoewer/go-strcase v1.1.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/terraform-provider-openstack/terraform-provider-openstack v1.32.0
	github.com/terraform-providers/terraform-provider-aws v0.0.0
	github.com/terraform-providers/terraform-provider-azurerm v0.0.0
	github.com/terraform-providers/terraform-provider-dns v0.0.0-20191209223915-3fb1c1918eb1
	github.com/terraform-providers/terraform-provider-google v1.20.1-0.20200623174414-27107f2ee160
	github.com/terraform-providers/terraform-provider-ignition/v2 v2.1.0
	github.com/terraform-providers/terraform-provider-local v1.4.0
	github.com/terraform-providers/terraform-provider-random v1.3.2-0.20190925210718-83518d96ae4f
	github.com/terraform-providers/terraform-provider-vsphere v1.16.2
	github.com/ulikunitz/xz v0.5.7
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	github.com/vmware/govmomi v0.22.2
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20201018230417-eeed37f84f13
	google.golang.org/api v0.33.0
	google.golang.org/genproto v0.0.0-20200904004341-0bd0a958aa1d
	google.golang.org/grpc v1.31.1
	gopkg.in/AlecAivazis/survey.v1 v1.8.9-0.20200217094205-6773bdf39b7f
	gopkg.in/ini.v1 v1.51.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.3.0
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/cluster-api-provider-aws v0.0.0
	sigs.k8s.io/cluster-api-provider-azure v0.0.0
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0
	sigs.k8s.io/controller-tools v0.3.1-0.20200617211605-651903477185
	sigs.k8s.io/yaml v1.2.0
)

replace (
	bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d
	github.com/Azure/go-autorest => github.com/tombuildsstuff/go-autorest v14.0.1-0.20200416184303-d4e299a3c04a+incompatible
	github.com/Azure/go-autorest/autorest => github.com/tombuildsstuff/go-autorest/autorest v0.10.1-0.20200416184303-d4e299a3c04a
	github.com/Azure/go-autorest/autorest/azure/auth => github.com/tombuildsstuff/go-autorest/autorest/azure/auth v0.4.3-0.20200416184303-d4e299a3c04a
	github.com/equinix/terraform-provider-equinix-metal => github.com/packethost/terraform-provider-packet v1.7.3-0.20201202165003-a5613b748108
	github.com/go-log/log => github.com/go-log/log v0.1.1-0.20181211034820-a514cf01a3eb // Pinned by MCO

	github.com/hashicorp/terraform => github.com/openshift/terraform v0.12.20-openshift-4 // Pin to fork with deduplicated rpc types v0.12.20-openshift-4
	github.com/hashicorp/terraform-plugin-sdk => github.com/openshift/hashicorp-terraform-plugin-sdk v1.14.0-openshift // Pin to fork with public rpc types
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20200715132148-0f91f62a41fe // Use OpenShift fork
	github.com/metal3-io/cluster-api-provider-baremetal => github.com/openshift/cluster-api-provider-baremetal v0.0.0-20190821174549-a2a477909c1d // Pin OpenShift fork
	github.com/miekg/dns => github.com/miekg/dns v1.0.8 // Pin to terraform-provider-dns
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200601094953-95abe2d2f422 // Pin API
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200929181438-91d71ef2122c // Pin client-go
	github.com/openshift/cluster-api-provider-equinix-metal => github.com/detiber/openshift-provider-packet v0.0.0-20201117162756-512a178614c0
	github.com/openshift/machine-api-operator => github.com/detiber/machine-api-operator v0.0.0-20201113194109-75933fe1fd83

	github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20201009041932-4fe8559913b8 // Pin MCO so it doesn't get downgraded
	github.com/terraform-providers/terraform-provider-aws => github.com/openshift/terraform-provider-aws v1.60.1-0.20200630224953-76d1fb4e5699 // Pin to openshift fork with tag v2.67.0-openshift
	github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.40.1-0.20200707062554-97ea089cc12a // release-2.17.0 branch
	github.com/terraform-providers/terraform-provider-ignition/v2 => github.com/community-terraform-providers/terraform-provider-ignition/v2 v2.1.0
	github.com/terraform-providers/terraform-provider-vsphere => github.com/openshift/terraform-provider-vsphere v1.18.1-openshift-2
	github.com/vmware/govmomi => github.com/vmware/govmomi v0.22.2-0.20200420222347-5fceac570f29
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20201022175424-d30c7a274820 // Pin OpenShift fork
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20200120114645-8a9592f1f87b // Pin OpenShift fork
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20200526112135-319a35b2e38e // Pin OpenShift fork
)
