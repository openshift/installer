module github.com/openshift/installer

go 1.18

require (
	cloud.google.com/go/monitoring v1.6.0
	github.com/AlecAivazis/survey/v2 v2.3.5
	github.com/Azure/azure-sdk-for-go v51.2.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.3.0
	github.com/Azure/go-autorest/autorest v0.11.27
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.1
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/IBM-Cloud/bluemix-go v0.0.0-20211102075456-ffc4e11dfb16
	github.com/IBM-Cloud/power-go-client v1.2.0
	github.com/IBM/go-sdk-core/v4 v4.9.0
	github.com/IBM/go-sdk-core/v5 v5.9.5
	github.com/IBM/networking-go-sdk v0.14.0
	github.com/IBM/platform-services-go-sdk v0.18.16
	github.com/IBM/vpc-go-sdk v0.20.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1264
	github.com/aliyun/aliyun-oss-go-sdk v2.1.8+incompatible
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.44.51
	github.com/clarketm/json v1.14.1
	github.com/containers/image v3.0.2+incompatible
	github.com/containers/image/v5 v5.22.1
	github.com/coreos/ignition/v2 v2.14.0
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
	github.com/openshift/assisted-image-service v0.0.0-20220506122314-2f689a1084b8
	github.com/openshift/assisted-service v0.0.0-20220928142635-a40422bdea61
	github.com/openshift/assisted-service/api v0.0.0
	github.com/openshift/assisted-service/models v0.0.0
	github.com/openshift/client-go v0.0.0-20220831193253-4950ae70c8ea
	github.com/openshift/cloud-credential-operator v0.0.0-20200316201045-d10080b52c9e
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20220408122422-7a548effc26e
	github.com/openshift/cluster-api-provider-ibmcloud v0.0.1-0.20220201105455-8014e5e894b0
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20191219173431-2336783d4603
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20220323121149-e3f2850dd519
	github.com/openshift/hive/apis v0.0.0-20220222213051-def9088fdb5a
	github.com/openshift/library-go v0.0.0-20220922140741-7772048e4447
	github.com/openshift/machine-config-operator v0.0.0
	github.com/ovirt/go-ovirt v0.0.0-20210809163552-d4276e35d3db
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.9.5
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/client_golang v1.12.2
	github.com/prometheus/common v0.32.1
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
	github.com/thedevsaddam/retry v0.0.0-20200324223450-9769a859cc6d
	github.com/ulikunitz/xz v0.5.10
	github.com/vincent-petithory/dataurl v1.0.0
	github.com/vmware/govmomi v0.27.4
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/oauth2 v0.0.0-20220622183110-fd043fe589d2
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab
	google.golang.org/api v0.91.0
	google.golang.org/genproto v0.0.0-20220808131553-a91ffa7f803e
	google.golang.org/grpc v1.48.0
	gopkg.in/ini.v1 v1.66.6
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.25.2
	k8s.io/apiextensions-apiserver v0.25.0
	k8s.io/apimachinery v0.25.2
	k8s.io/cli-runtime v0.25.2
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
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/Microsoft/hcsshim v0.9.3 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/containerd/cgroups v1.0.3 // indirect
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/containers/libtrust v0.0.0-20200511145503-9c3a6c22cd9a // indirect
	github.com/containers/ocicrypt v1.1.5 // indirect
	github.com/containers/storage v1.42.0 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/runc v1.1.3 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20171119174359-809beceb2371 // indirect
	github.com/spf13/pflag v1.0.6-0.20210604193023-d5e0c0615ace // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	go.mongodb.org/mongo-driver v1.8.3 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
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
	k8s.io/apiserver v0.25.2 // indirect
	k8s.io/component-base v0.25.2 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/kubectl v0.25.2 // indirect
	sigs.k8s.io/controller-runtime v0.11.2 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/kustomize/api v0.12.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

require (
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.21.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.9.2 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.20 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/PaesslerAG/gval v1.0.0 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
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
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/analysis v0.21.2 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/loads v0.21.1 // indirect
	github.com/go-openapi/runtime v0.23.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-openapi/validate v0.22.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gobuffalo/flect v0.2.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.0 // indirect
	github.com/hashicorp/go-version v1.5.0 // indirect
	github.com/hashicorp/hc-install v0.3.2 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
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
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/openshift/cluster-api v0.0.0-20190805113604-f8de78af80fc // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/openshift/oc v0.0.0-alpha.0.0.20221118061721-3ac1b026860b
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
	gorm.io/gorm v1.23.8 // indirect
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
replace k8s.io/client-go => k8s.io/client-go v0.25.0

// Needed so that the InstallConfig CRD can be created. Later versions of controller-gen balk at using IPNet as a field.
replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.3.1-0.20200617211605-651903477185

// Override the OpenShift API version in hive

replace github.com/openshift/api => github.com/openshift/api v0.0.0-20221004120407-c46852673d03

replace github.com/terraform-providers/terraform-provider-nutanix => github.com/nutanix/terraform-provider-nutanix v1.5.0

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.10.0

replace github.com/openshift/assisted-service/api => github.com/openshift/assisted-service/api v0.0.0-20220928142635-a40422bdea61

replace github.com/openshift/assisted-service/models => github.com/openshift/assisted-service/models v0.0.0-20220928142635-a40422bdea61

// https://bugzilla.redhat.com/show_bug.cgi?id=2064702
replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd

// https://github.com/kubernetes/kubernetes/issues/79384#issuecomment-505627280
// TODO better comment
replace k8s.io/api => k8s.io/api v0.25.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.25.0

replace k8s.io/apiserver => k8s.io/apiserver v0.25.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.0

replace k8s.io/code-generator => k8s.io/code-generator v0.25.0

replace k8s.io/component-base => k8s.io/component-base v0.25.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.25.0

replace k8s.io/controller-manager => k8s.io/controller-manager v0.25.0

replace k8s.io/cri-api => k8s.io/cri-api v0.25.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.0

replace k8s.io/kubectl => k8s.io/kubectl v0.25.0

replace k8s.io/kubelet => k8s.io/kubelet v0.25.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.0

replace k8s.io/metrics => k8s.io/metrics v0.25.0

replace k8s.io/mount-utils => k8s.io/mount-utils v0.25.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.25.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.0

replace github.com/apcera/gssapi => github.com/openshift/gssapi v0.0.0-20161010215902-5fb4217df13b
