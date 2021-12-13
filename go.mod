module github.com/openshift/installer

go 1.17

require (
	cloud.google.com/go v0.79.0
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/Azure/azure-sdk-for-go v51.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.1
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/IBM-Cloud/bluemix-go v0.0.0-20210611051827-cdc80c935c05
	github.com/IBM/go-sdk-core/v5 v5.4.3
	github.com/IBM/networking-go-sdk v0.14.0
	github.com/IBM/platform-services-go-sdk v0.18.16
	github.com/IBM/vpc-go-sdk v1.0.1
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1154
	github.com/aliyun/aliyun-oss-go-sdk v2.1.8+incompatible
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.37.4
	github.com/clarketm/json v1.14.1
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/ignition/v2 v2.9.0
	github.com/coreos/stream-metadata-go v0.1.3
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-playground/validator/v10 v10.2.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.2.0
	github.com/gophercloud/gophercloud v0.19.0
	github.com/gophercloud/utils v0.0.0-20210720165645-8a3ad2ad9e70
	github.com/h2non/filetype v1.0.12
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/metal3-io/baremetal-operator v0.0.0-20211201170610-92ffa60c683d
	github.com/metal3-io/baremetal-operator/apis v0.0.0
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.0.0
	github.com/openshift/api v0.0.0-20211119153416-313e51eab8c8
	github.com/openshift/client-go v0.0.0-20211025111749-96ca2abfc56c
	github.com/openshift/cloud-credential-operator v0.0.0-20200316201045-d10080b52c9e
	github.com/openshift/cluster-api-provider-alibaba v0.0.0-20211123142315-d16daa4e6349
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20210924143856-c2d3ece4da38
	github.com/openshift/cluster-api-provider-ibmcloud v0.0.0-20211008100740-4d7907adbd6b
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20191219173431-2336783d4603
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20211111151530-06177b773958
	github.com/openshift/library-go v0.0.0-20210811133500-5e31383de2a7
	github.com/openshift/machine-config-operator v0.0.0
	github.com/ovirt/go-ovirt v0.0.0-20210308100159-ac0bcbc88d7c
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/ulikunitz/xz v0.5.8
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	github.com/vmware/govmomi v0.24.0
	github.com/zclconf/go-cty v1.9.1 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c
	google.golang.org/api v0.44.0
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.38.0
	gopkg.in/ini.v1 v1.62.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.9.0
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/cluster-api-provider-aws v0.0.0
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0
	sigs.k8s.io/controller-tools v0.7.0
)

require (
	github.com/hashicorp/terraform-exec v0.13.1
	github.com/hashicorp/terraform-json v0.9.0
)

require (
	cloud.google.com/go/storage v1.11.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/IBM/go-sdk-core/v4 v4.10.0 // indirect
	github.com/Netflix/go-expect v0.0.0-20190729225929-0e00d9168667 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/vcontext v0.0.0-20201120045928-b0e13dab675c // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-openapi/errors v0.19.8 // indirect
	github.com/go-openapi/strfmt v0.20.1 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gobuffalo/flect v0.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/renameio v1.0.0 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.5.2 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6 // indirect
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/compress v1.11.8 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2-0.20190823105129-775207bd45b6 // indirect
	github.com/openshift/cluster-api v0.0.0-20190805113604-f8de78af80fc // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.mongodb.org/mongo-driver v1.5.1 // indirect
	go.opencensus.io v0.22.5 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	k8s.io/component-base v0.22.2 // indirect
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e // indirect
	sigs.k8s.io/controller-runtime v0.9.7 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.57.0
	github.com/IBM-Cloud/terraform-provider-ibm => github.com/openshift/terraform-provider-ibm v1.26.2-openshift-2
	github.com/go-log/log => github.com/go-log/log v0.1.1-0.20181211034820-a514cf01a3eb // Pinned by MCO
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.17.0
	github.com/hashicorp/terraform-exec => github.com/hashicorp/terraform-exec v0.0.0-20210211154940-d3d1d84e068a
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20211201170610-92ffa60c683d // Use OpenShift fork
	github.com/metal3-io/baremetal-operator/apis => github.com/openshift/baremetal-operator/apis v0.0.0-20211201170610-92ffa60c683d // Use OpenShift fork
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils => github.com/openshift/baremetal-operator/pkg/hardwareutils v0.0.0-20211201170610-92ffa60c683d // Use OpenShift fork
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200929181438-91d71ef2122c // Pin client-go
	github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20201009041932-4fe8559913b8 // Pin MCO so it doesn't get downgraded
	github.com/spf13/afero => github.com/spf13/afero v1.2.2
	github.com/terraform-providers/terraform-provider-aws => github.com/openshift/terraform-provider-aws v1.60.1-0.20211215220004-24df6d73af46 // Pin to openshift fork with tag v3.1.0-openshift-2
	github.com/terraform-providers/terraform-provider-azurerm => github.com/openshift/terraform-provider-azurerm v1.44.1-0.20210224232508-7509319df0f4 // Pin to 2.48.0-openshift
	github.com/terraform-providers/terraform-provider-azurestack => github.com/openshift/terraform-provider-azurestack v0.10.0-openshift // Use OpenShift fork
	github.com/terraform-providers/terraform-provider-google v1.20.1-0.20200623174414-27107f2ee160 => github.com/openshift/terraform-providers-terraform-provider-google v1.20.1-0.20211201190933-7b79c6d1afc8 // Pin to 3.27.0-openshift
	github.com/terraform-providers/terraform-provider-ignition/v2 => github.com/community-terraform-providers/terraform-provider-ignition/v2 v2.1.0
	k8s.io/client-go => k8s.io/client-go v0.22.0
	k8s.io/kubectl => k8s.io/kubectl v0.21.0-rc.0
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20210121023454-5ffc5f422a80
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20210626224711-5d94c794092f // Indirect dependency through MAO from cluster API providers
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20211111204942-611d320170af
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.0-alpha.1
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.3.1-0.20200617211605-651903477185
)

// Prevent the following modules from upgrading version as result of terraform-provider-kubernetes module
// The following modules need to be locked to compile correctly with
// terraform-provider-azure, terraform-provider-google, and
// terraform-provider-ibm. This also locks in Terraform to v0.14.6.
replace (
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.10 // indirect
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/Azure/go-autorest/autorest/validation => github.com/Azure/go-autorest/autorest/validation v0.3.1
	github.com/IBM/vpc-go-sdk => github.com/IBM/vpc-go-sdk v0.7.0
	github.com/apparentlymart/go-cidr => github.com/apparentlymart/go-cidr v1.1.0
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.35.37
	github.com/go-openapi/errors => github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.4
	github.com/go-openapi/validate => github.com/go-openapi/validate v0.19.8
	github.com/hashicorp/go-getter => github.com/hashicorp/go-getter v1.5.2
	github.com/hashicorp/go-plugin => github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/terraform-plugin-test/v2 => github.com/hashicorp/terraform-plugin-test/v2 v2.2.0
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.7
	google.golang.org/api => google.golang.org/api v0.33.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc => google.golang.org/grpc v1.35.0
)
