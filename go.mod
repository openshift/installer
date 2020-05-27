module github.com/openshift/installer

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.3
	github.com/Azure/go-autorest/autorest/adal v0.8.1-0.20191028180845-3492b2aff503
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.1
	github.com/Azure/go-autorest/autorest/to v0.3.1-0.20191028180845-3492b2aff503
	github.com/Azure/go-autorest/autorest/validation v0.2.1-0.20191028180845-3492b2aff503 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20191115210519-2b2be6cc8ed4 // indirect
	github.com/ChrisTrenkamp/goxpath v0.0.0-20190607011252-c5096ec8773d // indirect
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.60.285 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.0.4+incompatible // indirect
	github.com/aliyun/aliyun-tablestore-go-sdk v4.1.3+incompatible // indirect
	github.com/antchfx/xpath v1.1.2 // indirect
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.29.24
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/c4milo/gotoolkit v0.0.0-20190525173301-67483a18c17a // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0 // indirect
	github.com/coreos/ignition v0.35.0
	github.com/dmacvicar/terraform-provider-libvirt v0.6.2
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/go-playground/validator/v10 v10.2.0
	github.com/golang/mock v1.3.1
	github.com/google/go-cmp v0.3.2-0.20191028172631-481baca67f93 // indirect
	github.com/google/martian v2.1.1-0.20190517191504-25dcb96d9e51+incompatible // indirect
	github.com/google/uuid v1.1.1
	github.com/gophercloud/gophercloud v0.8.0
	github.com/gophercloud/utils v0.0.0-20191212191830-4533a07bd492
	github.com/gopherjs/gopherjs v0.0.0-20191106031601-ce3c9ade29de // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/h2non/filetype v1.0.12
	github.com/hashicorp/go-azure-helpers v0.10.0
	github.com/hashicorp/go-plugin v1.0.1
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/memberlist v0.1.5 // indirect
	github.com/hashicorp/serf v0.8.5 // indirect
	github.com/hashicorp/terraform v0.12.21
	github.com/hashicorp/terraform-plugin-sdk v1.8.0
	github.com/hashicorp/vault v1.3.0 // indirect
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c // indirect
	github.com/joyent/triton-go v1.7.0 // indirect
	github.com/keybase/go-crypto v0.0.0-20190828182435-a05457805304 // indirect
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/libvirt/libvirt-go-xml v5.10.0+incompatible // indirect
	github.com/lusis/go-artifactory v0.0.0-20180304164534-a47f63f234b2 // indirect
	github.com/masterzen/simplexml v0.0.0-20190410153822-31eea3082786 // indirect
	github.com/masterzen/winrm v0.0.0-20190308153735-1d17eaf15943 // indirect
	github.com/metal3-io/baremetal-operator v0.0.0
	github.com/metal3-io/cluster-api-provider-baremetal v0.0.0
	github.com/mitchellh/cli v1.0.0
	github.com/opencontainers/image-spec v1.0.2-0.20190823105129-775207bd45b6 // indirect
	github.com/openshift-metal3/terraform-provider-ironic v0.2.1
	github.com/openshift/api v3.9.1-0.20191111211345-a27ff30ebf09+incompatible
	github.com/openshift/client-go v0.0.0-20200320150128-a906f3d8e723
	github.com/openshift/cloud-credential-operator v0.0.0-20200316201045-d10080b52c9e
	github.com/openshift/cluster-api v0.0.0-20191129101638-b09907ac6668
	github.com/openshift/cluster-api-provider-gcp v0.0.1-0.20200120152131-1b09fd9e7156
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20191219173431-2336783d4603
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20200504092944-27473ea1ae43
	github.com/openshift/library-go v0.0.0-20200324092245-db2a8546af81
	github.com/openshift/machine-api-operator v0.2.1-0.20200429102619-d36974451290
	github.com/openshift/machine-config-operator v4.2.0-alpha.0.0.20190917115525-033375cbe820+incompatible
	github.com/ovirt/go-ovirt v0.0.0-20200428093010-9bcc4fd4e6c0
	github.com/ovirt/terraform-provider-ovirt v0.4.3-0.20200406133650-74a154c1d861
	github.com/packer-community/winrmcp v0.0.0-20180921211025-c76d91c1e7db // indirect
	github.com/pborman/uuid v1.2.0
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.2.1 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/satori/uuid v1.2.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/stoewer/go-strcase v1.1.0 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/terraform-providers/terraform-provider-aws v1.60.1-0.20200319153356-a183662640d9
	github.com/terraform-providers/terraform-provider-azuread v0.7.0 // indirect
	github.com/terraform-providers/terraform-provider-azurerm v0.0.0
	github.com/terraform-providers/terraform-provider-google v1.20.1-0.20200204003432-77547e3e7d52
	github.com/terraform-providers/terraform-provider-ignition v1.2.1
	github.com/terraform-providers/terraform-provider-local v1.4.0
	github.com/terraform-providers/terraform-provider-openstack v1.25.0
	github.com/terraform-providers/terraform-provider-random v1.3.2-0.20190925210718-83518d96ae4f
	github.com/terraform-providers/terraform-provider-vsphere v1.16.2
	github.com/ulikunitz/xz v0.5.6
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	github.com/vmware/govmomi v0.22.2
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/api v0.14.0
	gopkg.in/AlecAivazis/survey.v1 v1.8.9-0.20200217094205-6773bdf39b7f
	gopkg.in/ini.v1 v1.51.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8
	sigs.k8s.io/cluster-api-provider-aws v0.0.0
	sigs.k8s.io/cluster-api-provider-azure v0.0.0
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0
	sigs.k8s.io/controller-tools v0.3.0
)

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0 // Pin non-versioned import to v22.0.0
	github.com/go-log/log => github.com/go-log/log v0.1.1-0.20181211034820-a514cf01a3eb // Pinned by MCO
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.6.2 // Pin to version required by terraform fork
	github.com/hashicorp/terraform => github.com/openshift/hashicorp-terraform v0.12.20-openshift-2 // Pin to fork with deduplicated rpc types
	github.com/hashicorp/terraform-plugin-sdk => github.com/openshift/hashicorp-terraform-plugin-sdk v1.6.0-openshift // Pin to fork with public rpc types
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20200206190020-71b826cc0f0a // Use OpenShift fork
	github.com/metal3-io/cluster-api-provider-baremetal => github.com/openshift/cluster-api-provider-baremetal v0.0.0-20190821174549-a2a477909c1d // Pin OpenShift fork
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200413201024-c6e8c9b6eb9a // Pin API
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240 // Pin client-go
	github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20200130220348-e5685c0cf530 // Pin MCO so it doesn't get downgraded
	github.com/ovirt/go-ovirt => github.com/ovirt/go-ovirt v0.0.0-20200320082526-4e97a11ff083
	github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.41.1-openshift-3 // Pin to openshift fork with IPv6 fixes
	github.com/terraform-providers/terraform-provider-vsphere => github.com/openshift/terraform-provider-vsphere v1.18.1-openshift-1
	github.com/vmware/govmomi => github.com/vmware/govmomi v0.22.2-0.20200420222347-5fceac570f29
	google.golang.org/api => google.golang.org/api v0.13.0 // Pin to version required by tf-provider-google
	k8s.io/api => k8s.io/api v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/apiserver => k8s.io/apiserver v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/client-go => k8s.io/client-go v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/code-generator => k8s.io/code-generator v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/component-base => k8s.io/component-base v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/cri-api => k8s.io/cri-api v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kubectl => k8s.io/kubectl v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kubelet => k8s.io/kubelet v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/kubernetes => k8s.io/kubernetes v1.17.1 // Replaced by MCO/CRI-O
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/metrics => k8s.io/metrics v0.17.1 // Replaced by MCO/CRI-O
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.1 // Replaced by MCO/CRI-O
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20200506073438-9d49428ff837 // Pin OpenShift fork
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20200120114645-8a9592f1f87b // Pin OpenShift fork
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20200323110431-3311de91e078 // Pin OpenShift fork
	sigs.k8s.io/controller-tools => github.com/abhinavdahiya/controller-tools v0.3.1-0.20200430222905-6fdf2d5fc069 // Using fork for sigs.k8s.io/controller-tools#427
)
