module github.com/openshift/installer

go 1.18

require (
	cloud.google.com/go v0.81.0
	github.com/AlecAivazis/survey/v2 v2.3.5
	github.com/Azure/azure-sdk-for-go v51.2.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.3.0
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.1
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/IBM-Cloud/bluemix-go v0.0.0-20211102075456-ffc4e11dfb16
	github.com/IBM-Cloud/power-go-client v1.1.5
	github.com/IBM/go-sdk-core/v4 v4.9.0
	github.com/IBM/go-sdk-core/v5 v5.9.5
	github.com/IBM/networking-go-sdk v0.14.0
	github.com/IBM/platform-services-go-sdk v0.18.16
	github.com/IBM/vpc-go-sdk v0.20.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1264
	github.com/aliyun/aliyun-oss-go-sdk v2.1.8+incompatible
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.43.19
	github.com/clarketm/json v1.14.1
	github.com/containers/image v3.0.2+incompatible
	github.com/coreos/ignition/v2 v2.13.0
	github.com/coreos/stream-metadata-go v0.1.8
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.8
	github.com/google/uuid v1.3.0
	github.com/gophercloud/gophercloud v0.24.0
	github.com/gophercloud/utils v0.0.0-20220307143606-8e7800759d16
	github.com/h2non/filetype v1.0.12
	github.com/hashicorp/terraform-exec v0.16.1
	github.com/kdomanski/iso9660 v0.2.1
	github.com/libvirt/libvirt-go v5.10.0+incompatible
	github.com/metal3-io/baremetal-operator v0.0.0-20220128094204-28771f489634
	github.com/metal3-io/baremetal-operator/apis v0.0.0
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.0.0
	github.com/nutanix-cloud-native/prism-go-client v0.2.1-0.20220804130801-c8a253627c64
	github.com/openshift/api v3.9.1-0.20191111211345-a27ff30ebf09+incompatible
	github.com/openshift/assisted-image-service v0.0.0-20220307202600-054a1afa8d28
	github.com/openshift/assisted-service v1.0.10-0.20220223093655-7ada9949bf1d
	github.com/openshift/client-go v0.0.0-20211209144617-7385dd6338e3
	github.com/openshift/cloud-credential-operator v0.0.0-20200316201045-d10080b52c9e
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20220408122422-7a548effc26e
	github.com/openshift/cluster-api-provider-ibmcloud v0.0.1-0.20220201105455-8014e5e894b0
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20191219173431-2336783d4603
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20220323121149-e3f2850dd519
	github.com/openshift/hive/apis v0.0.0-20210506000654-5c038fb05190
	github.com/openshift/library-go v0.0.0-20220121154930-b7889002d63e
	github.com/openshift/machine-config-operator v0.0.0
	github.com/ovirt/go-ovirt v0.0.0-20210809163552-d4276e35d3db
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.9.3
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.32.1
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.7.2
	github.com/thedevsaddam/retry v0.0.0-20200324223450-9769a859cc6d
	github.com/ulikunitz/xz v0.5.10
	github.com/vincent-petithory/dataurl v1.0.0
	github.com/vmware/govmomi v0.27.4
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
	google.golang.org/api v0.44.0
	google.golang.org/genproto v0.0.0-20220107163113-42d7afdf6368
	google.golang.org/grpc v1.46.0
	gopkg.in/ini.v1 v1.66.4
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.24.3
	k8s.io/apiextensions-apiserver v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cloud-provider-vsphere v0.0.0
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.70.1
	k8s.io/utils v0.0.0-20220812165043-ad590609e2e5
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0
	sigs.k8s.io/controller-tools v0.9.2
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/shurcooL/httpfs v0.0.0-20171119174359-809beceb2371 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	go.mongodb.org/mongo-driver v1.8.3 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.24.3 // indirect
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42 // indirect
	sigs.k8s.io/controller-runtime v0.11.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.2 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/PaesslerAG/gval v1.0.0 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cavaliercoder/go-cpio v0.0.0-20180626203310-925f9528c45e // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/vcontext v0.0.0-20211021162308-f1dbbca7bef4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dimchansky/utfbom v1.1.0 // indirect
	github.com/diskfs/go-diskfs v1.2.1-0.20210727185522-a769efacd235 // indirect
	github.com/emicklei/go-restful v2.14.2+incompatible // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/loads v0.21.1 // indirect
	github.com/go-openapi/runtime v0.23.0 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-openapi/validate v0.20.3 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gobuffalo/flect v0.2.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.0 // indirect
	github.com/hashicorp/go-version v1.5.0 // indirect
	github.com/hashicorp/hc-install v0.3.2 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/openshift/cluster-api v0.0.0-20190805113604-f8de78af80fc // indirect
	github.com/openshift/custom-resource-status v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/xattr v0.4.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/net v0.0.0-20220812174116-3211cb980234 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gorm.io/gorm v1.22.3 // indirect
)

// OpenShift Forks
replace (
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20220128094204-28771f489634
	github.com/metal3-io/baremetal-operator/apis => github.com/openshift/baremetal-operator/apis v0.0.0-20220128094204-28771f489634
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils => github.com/openshift/baremetal-operator/pkg/hardwareutils v0.0.0-20220128094204-28771f489634
	k8s.io/cloud-provider-vsphere => github.com/openshift/cloud-provider-vsphere v1.19.1-0.20211222185833-7829863d0558
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.5
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20200929152424-eab2e087f366 // Indirect dependency through MAO from cluster API providers
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20210626224711-5d94c794092f // Indirect dependency through MAO from cluster API providers
	sigs.k8s.io/cluster-api-provider-openstack => github.com/openshift/cluster-api-provider-openstack v0.0.0-20211111204942-611d320170af
)

// Pin MCO so it doesn't get downgraded
replace github.com/openshift/machine-config-operator => github.com/openshift/machine-config-operator v0.0.1-0.20201009041932-4fe8559913b8

// Needed because machine-api-operator uses a "later" v12 version, which is actually an earlier version.
// This should be kept in line with the k8s version used.
replace k8s.io/client-go => k8s.io/client-go v0.24.0

// Needed so that the InstallConfig CRD can be created. Later versions of controller-gen balk at using IPNet as a field.
replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.3.1-0.20200617211605-651903477185

// Override the OpenShift API version in hive
replace github.com/openshift/api => github.com/openshift/api v0.0.0-20220823143838-5768cc618ba0

replace github.com/terraform-providers/terraform-provider-nutanix => github.com/nutanix/terraform-provider-nutanix v1.5.0

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.10.0
