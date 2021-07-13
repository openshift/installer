module github.com/openshift/installer

go 1.14

require (
	cloud.google.com/go v0.65.0
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/Azure/azure-sdk-for-go v51.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/adal v0.9.10
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.1
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-ntlmssp v0.0.0-20191115210519-2b2be6cc8ed4 // indirect
	github.com/ChrisTrenkamp/goxpath v0.0.0-20190607011252-c5096ec8773d // indirect
	github.com/IBM-Cloud/bluemix-go v0.0.0-20210611051827-cdc80c935c05
	github.com/IBM-Cloud/terraform-provider-ibm v1.26.2
	github.com/IBM/go-sdk-core v1.1.0 // indirect
	github.com/IBM/go-sdk-core/v4 v4.10.0 // indirect
	github.com/IBM/go-sdk-core/v5 v5.4.3
	github.com/IBM/networking-go-sdk v0.14.0
	github.com/IBM/platform-services-go-sdk v0.18.16
	github.com/IBM/vpc-go-sdk v1.0.1
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667 // indirect
	github.com/antchfx/xpath v1.1.2 // indirect
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.35.20
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/c4milo/gotoolkit v0.0.0-20190525173301-67483a18c17a // indirect
	github.com/clarketm/json v1.14.1
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/ignition/v2 v2.9.0
	github.com/coreos/stream-metadata-go v0.0.0-20210225230131-70edb9eb47b3
	github.com/dmacvicar/terraform-provider-libvirt v0.6.4-0.20201216193629-2b60d7626ff8
	github.com/fatih/color v1.10.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-playground/validator/v10 v10.2.0
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/google/martian v2.1.1-0.20190517191504-25dcb96d9e51+incompatible // indirect
	github.com/google/uuid v1.2.0
	github.com/gophercloud/gophercloud v0.17.0
	github.com/gophercloud/utils v0.0.0-20210323225332-7b186010c04f
	github.com/h2non/filetype v1.0.12
	github.com/hashicorp/go-azure-helpers v0.13.0
	github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/terraform v0.13.4
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/hashicorp/terraform-provider-kubernetes v1.13.3
	github.com/hashicorp/terraform-provider-vsphere v1.24.3
	github.com/hashicorp/vault v1.3.0 // indirect
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20191119172530-79f836b90111
	github.com/kubevirt/terraform-provider-kubevirt v0.0.0-00010101000000-000000000000
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/masterzen/simplexml v0.0.0-20190410153822-31eea3082786 // indirect
	github.com/masterzen/winrm v0.0.0-20190308153735-1d17eaf15943 // indirect
	github.com/metal3-io/baremetal-operator v0.0.0-20210527161605-4e331bfd4b1d
	github.com/metal3-io/cluster-api-provider-baremetal v0.0.0
	github.com/mitchellh/cli v1.1.1
	github.com/openshift-metal3/terraform-provider-ironic v0.2.6
	github.com/openshift/api v0.0.0-20210622200201-7740bfb97397
	github.com/openshift/client-go v0.0.0-20210409155308-a8e62c60e930
	github.com/openshift/cloud-credential-operator v0.0.0-20200316201045-d10080b52c9e
	github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20201203141909-4dc702fd57a5
	github.com/openshift/cluster-api-provider-ibmcloud v0.0.0-20210702173623-676faba9895d
	github.com/openshift/cluster-api-provider-kubevirt v0.0.0-20201214114543-e5aed9c73f1f
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20191219173431-2336783d4603
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20210622084102-c4f9f269bcb7
	github.com/openshift/library-go v0.0.0-20210408164723-7a65fdb398e2
	github.com/openshift/machine-api-operator v0.2.1-0.20210505133115-b7ef098180db
	github.com/openshift/machine-config-operator v0.0.0
	github.com/ovirt/go-ovirt v0.0.0-20210308100159-ac0bcbc88d7c
	github.com/ovirt/terraform-provider-ovirt v0.99.1-0.20210628083913-78ce17b16709
	github.com/packer-community/winrmcp v0.0.0-20180921211025-c76d91c1e7db // indirect
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/satori/uuid v1.2.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	github.com/terraform-provider-openstack/terraform-provider-openstack v1.37.0
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20200807230610-d5346d47e3af
	github.com/terraform-providers/terraform-provider-azurerm v1.44.1-0.20200911233120-57b2bfc9d42c
	github.com/terraform-providers/terraform-provider-google v1.20.1-0.20200623174414-27107f2ee160
	github.com/terraform-providers/terraform-provider-ignition/v2 v2.1.0
	github.com/terraform-providers/terraform-provider-local v1.4.0
	github.com/terraform-providers/terraform-provider-random v1.3.2-0.20190925210718-83518d96ae4f
	github.com/ulikunitz/xz v0.5.8
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	github.com/vmware/govmomi v0.24.0
	github.com/zclconf/go-cty v1.6.1
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073
	google.golang.org/api v0.33.0
	google.golang.org/genproto v0.0.0-20201110150050-8816d57aaa9a
	google.golang.org/grpc v1.32.0
	gopkg.in/ini.v1 v1.61.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.0-rc.0
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.8.0
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009
	kubevirt.io/client-go v0.29.0
	kubevirt.io/containerized-data-importer v1.10.9
	sigs.k8s.io/cluster-api-provider-aws v0.0.0
	sigs.k8s.io/cluster-api-provider-azure v0.0.0
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0
	sigs.k8s.io/controller-tools v0.4.1
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.57.0
	github.com/IBM-Cloud/terraform-provider-ibm => github.com/openshift/terraform-provider-ibm v1.26.2-openshift
	github.com/go-log/log => github.com/go-log/log v0.1.1-0.20181211034820-a514cf01a3eb // Pinned by MCO
	github.com/hashicorp/terraform => github.com/openshift/terraform v0.12.20-openshift-4 // Pin to fork with deduplicated rpc types v0.12.20-openshift-4
	github.com/hashicorp/terraform-plugin-sdk => github.com/openshift/hashicorp-terraform-plugin-sdk v1.14.0-openshift // Pin to fork with public rpc types
	github.com/hashicorp/terraform-provider-vsphere => github.com/openshift/terraform-provider-vsphere v1.24.3-openshift
	github.com/kubevirt/terraform-provider-kubevirt => github.com/nirarg/terraform-provider-kubevirt v0.0.0-20201222125919-101cee051ed3
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20210527161605-4e331bfd4b1d // Use OpenShift fork
	github.com/metal3-io/cluster-api-provider-baremetal => github.com/openshift/cluster-api-provider-baremetal v0.0.0-20190821174549-a2a477909c1d // Pin OpenShift fork
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200929181438-91d71ef2122c // Pin client-go
	github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20201009041932-4fe8559913b8 // Pin MCO so it doesn't get downgraded
	github.com/terraform-providers/terraform-provider-aws => github.com/openshift/terraform-provider-aws v1.60.1-0.20210622193531-7d13cfbb1a8c // Pin to openshift fork with tag v2.67.0-openshift-1
	github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.44.1-0.20210224232508-7509319df0f4 // Pin to 2.48.0-openshift
	github.com/terraform-providers/terraform-provider-ignition/v2 => github.com/community-terraform-providers/terraform-provider-ignition/v2 v2.1.0
	k8s.io/client-go => k8s.io/client-go v0.21.0-rc.0
	k8s.io/kubectl => k8s.io/kubectl v0.21.0-rc.0
	kubevirt.io/client-go => kubevirt.io/client-go v0.29.0
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20210121023454-5ffc5f422a80
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20210626224711-5d94c794092f
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20210302164104-8498241fa4bd
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.0-alpha.1
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.3.1-0.20200617211605-651903477185
)

// Prevent the following modules from upgrading version as result of terraform-provider-kubernetes module
// The following modules need to be locked to compile correctly with terraform-provider-azure, terraform-provider-google, and terraform-provider-ibm
replace (
	github.com/IBM/vpc-go-sdk => github.com/IBM/vpc-go-sdk v0.7.0
	github.com/apparentlymart/go-cidr => github.com/apparentlymart/go-cidr v1.0.1
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.32.3
	github.com/go-openapi/errors => github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.4
	github.com/go-openapi/validate => github.com/go-openapi/validate v0.19.8
	github.com/hashicorp/go-plugin => github.com/hashicorp/go-plugin v1.2.2
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.7
	google.golang.org/api => google.golang.org/api v0.25.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
