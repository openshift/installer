module github.com/openshift/installer

go 1.23.2

toolchain go1.23.5

require (
	cloud.google.com/go/kms v1.20.2
	cloud.google.com/go/monitoring v1.21.2
	cloud.google.com/go/storage v1.43.0
	github.com/AlecAivazis/survey/v2 v2.3.5
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.17.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.2
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3 v3.0.0-beta.2
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4 v4.2.1
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.7.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2 v2.2.1
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns v1.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph v0.8.2
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage v1.6.0
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.2.1
	github.com/Azure/go-autorest/autorest v0.11.30
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/IBM-Cloud/bluemix-go v0.0.0-20211102075456-ffc4e11dfb16
	github.com/IBM-Cloud/power-go-client v1.11.0
	github.com/IBM/go-sdk-core/v5 v5.19.0
	github.com/IBM/ibm-cos-sdk-go v1.12.2
	github.com/IBM/keyprotect-go-client v0.12.2
	github.com/IBM/networking-go-sdk v0.51.2
	github.com/IBM/platform-services-go-sdk v0.79.0
	github.com/IBM/vpc-go-sdk v0.65.0
	github.com/apparentlymart/go-cidr v1.1.0
	github.com/awalterschulze/gographviz v0.0.0-20190522210029-fa59802746ab
	github.com/aws/aws-sdk-go v1.55.5
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/config v1.29.14
	github.com/aws/aws-sdk-go-v2/credentials v1.17.67
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.159.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.43.11
	github.com/aws/aws-sdk-go-v2/service/iam v1.32.0
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.25.17
	github.com/aws/aws-sdk-go-v2/service/route53 v1.48.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.80.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.19
	github.com/aws/smithy-go v1.22.3
	github.com/cavaliercoder/go-cpio v0.0.0-20180626203310-925f9528c45e
	github.com/clarketm/json v1.14.1
	github.com/containers/image/v5 v5.31.0
	github.com/coreos/ignition/v2 v2.20.0
	github.com/coreos/stream-metadata-go v0.1.8
	github.com/daixiang0/gci v0.13.4
	github.com/digitalocean/go-libvirt v0.0.0-20240220204746-fcabe97a6eed
	github.com/diskfs/go-diskfs v1.4.0
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible
	github.com/go-logr/logr v1.4.2
	github.com/go-openapi/errors v0.22.0
	github.com/go-openapi/runtime v0.28.0
	github.com/go-openapi/strfmt v0.23.0
	github.com/go-openapi/swag v0.23.0
	github.com/go-openapi/validate v0.24.0
	github.com/go-playground/validator/v10 v10.24.0
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/googleapis/gax-go/v2 v2.14.0
	github.com/gophercloud/gophercloud/v2 v2.5.0
	github.com/gophercloud/utils/v2 v2.0.0-20250212084022-725b94822eeb
	github.com/h2non/filetype v1.0.12
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.7.0
	github.com/hashicorp/terraform-exec v0.17.3
	github.com/jarcoal/httpmock v1.3.1
	github.com/jongio/azidext/go/azidext v0.5.0
	github.com/kdomanski/iso9660 v0.2.1
	github.com/metal3-io/baremetal-operator/apis v0.4.0
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.4.0
	github.com/microsoft/kiota-authentication-azure-go v0.6.0
	github.com/microsoftgraph/msgraph-sdk-go v0.59.0
	github.com/nutanix-cloud-native/cluster-api-provider-nutanix v1.5.4-0.20250116153252-296a5347104c
	github.com/nutanix-cloud-native/prism-go-client v0.5.0
	github.com/onsi/gomega v1.36.2
	github.com/openshift/api v0.0.0-20250527072845-f5e205b58365
	github.com/openshift/assisted-image-service v0.0.0-20240607085136-02df2e56dde6
	github.com/openshift/assisted-service/api v0.0.0
	github.com/openshift/assisted-service/client v0.0.0
	github.com/openshift/assisted-service/models v0.0.0
	github.com/openshift/client-go v0.0.0-20241203091221-452dfb8fa071
	github.com/openshift/cloud-credential-operator v0.0.0-20240404165937-5e8812d64187
	github.com/openshift/cluster-api-provider-baremetal v0.0.0-20220408122422-7a548effc26e
	github.com/openshift/cluster-api-provider-libvirt v0.2.1-0.20230308152226-83c0473d4429
	github.com/openshift/cluster-api-provider-ovirt v0.1.1-0.20220323121149-e3f2850dd519
	github.com/openshift/hive/apis v0.0.0-20231220215202-ad99b9e52d27
	github.com/openshift/library-go v0.0.0-20250114132252-af5b21ebad2f
	github.com/openshift/machine-api-provider-gcp v0.0.1-0.20241021180644-0eca0846914a
	github.com/openshift/machine-api-provider-ibmcloud v0.0.0-20231207164151-6b0b8ea7b16d
	github.com/ovirt/go-ovirt v0.0.0-20210809163552-d4276e35d3db
	github.com/pelletier/go-toml v1.9.5
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.7
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/common v0.62.0
	github.com/rogpeppe/go-internal v1.13.1
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.9.1
	github.com/stretchr/testify v1.10.0
	github.com/thedevsaddam/retry v0.0.0-20200324223450-9769a859cc6d
	github.com/thoas/go-funk v0.9.3
	github.com/ulikunitz/xz v0.5.12
	github.com/vincent-petithory/dataurl v1.0.0
	github.com/vmware/govmomi v0.47.1
	go.uber.org/mock v0.5.0
	golang.org/x/crypto v0.36.0
	golang.org/x/oauth2 v0.27.0
	golang.org/x/sync v0.12.0
	golang.org/x/sys v0.31.0
	golang.org/x/term v0.30.0
	golang.org/x/text v0.23.0
	google.golang.org/api v0.214.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250218202821-56aae31c358a
	google.golang.org/grpc v1.71.0
	gopkg.in/ini.v1 v1.67.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.32.3
	k8s.io/apiextensions-apiserver v0.32.1
	k8s.io/apimachinery v0.32.3
	k8s.io/client-go v0.32.1
	k8s.io/cloud-provider-vsphere v1.31.0
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.130.1
	k8s.io/utils v0.0.0-20241210054802-24370beab758
	libvirt.org/go/libvirtxml v1.10002.0
	sigs.k8s.io/cluster-api v1.9.6
	sigs.k8s.io/cluster-api-provider-aws/v2 v2.7.1-0.20250314180547-17a09f59176c
	sigs.k8s.io/cluster-api-provider-azure v1.19.3
	sigs.k8s.io/cluster-api-provider-gcp v1.8.1-0.20250225090028-d80bfabadd3f
	sigs.k8s.io/cluster-api-provider-ibmcloud v0.11.0-alpha.0.0.20250319131234-d3cc59096981
	sigs.k8s.io/cluster-api-provider-openstack v0.11.1
	sigs.k8s.io/cluster-api-provider-vsphere v1.12.0-rc.0.0.20250203113257-38f6b67f9b7d
	sigs.k8s.io/controller-runtime v0.20.1
	sigs.k8s.io/controller-tools v0.16.3
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8
	sigs.k8s.io/yaml v1.4.0
)

require (
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/auth v0.13.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.6 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/iam v1.2.2 // indirect
	cloud.google.com/go/longrunning v0.6.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.24 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.3.3 // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/PaesslerAG/gval v1.0.0 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/asaskevich/govalidator/v11 v11.0.2-0.20250122183457-e11347878e23 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/containers/storage v1.54.0 // indirect
	github.com/coreos/go-json v0.0.0-20230131223807-18775e0fb4fb // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/coreos/vcontext v0.0.0-20230201181013-d72178a18687 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/elliotwutingfeng/asciiset v0.0.0-20230602022725-51bbb787efab // indirect
	github.com/emicklei/go-restful/v3 v3.12.1 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	github.com/go-openapi/analysis v0.23.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.22.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gobuffalo/flect v1.0.3 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gofrs/uuid/v5 v5.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.7.0-rc.1 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itchyny/gojq v0.12.8 // indirect
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/microsoft/kiota-abstractions-go v0.18.0 // indirect
	github.com/microsoft/kiota-http-go v0.16.0 // indirect
	github.com/microsoft/kiota-serialization-form-go v0.9.0 // indirect
	github.com/microsoft/kiota-serialization-json-go v0.9.0 // indirect
	github.com/microsoft/kiota-serialization-text-go v0.7.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go-core v0.34.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/moby/sys/mountinfo v0.7.1 // indirect
	github.com/moby/sys/user v0.1.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/openshift/assisted-service v1.0.10-0.20230830164851-6573b5d7021d // indirect
	github.com/openshift/custom-resource-status v1.1.3-0.20220503160415-f2fdb4999d87 // indirect
	github.com/openshift/machine-api-operator v0.2.1-0.20240930121047-57b7917e6140 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/ppc64le-cloud/powervs-utils v0.0.0-20240610070307-1c0d75a5c247 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/zclconf/go-cty v1.11.0 // indirect
	go.mongodb.org/mongo-driver v1.17.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.starlark.net v0.0.0-20230525235612-a134d8f9ddca // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20250210185358-939b2ce775ac // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/genproto v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/djherbis/times.v1 v1.3.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/gorm v1.24.5 // indirect
	k8s.io/apiserver v0.32.1 // indirect
	k8s.io/cli-runtime v0.31.3 // indirect
	k8s.io/cluster-bootstrap v0.32.1 // indirect
	k8s.io/component-base v0.32.1 // indirect
	k8s.io/kube-openapi v0.0.0-20241105132330-32ad38e42d3f // indirect
	k8s.io/kubectl v0.31.3 // indirect
	sigs.k8s.io/kustomize/api v0.17.3 // indirect
	sigs.k8s.io/kustomize/kyaml v0.17.2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
)

// OpenShift Forks
replace (
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20231128154154-6736c9b9c6c8
	github.com/metal3-io/baremetal-operator/apis => github.com/openshift/baremetal-operator/apis v0.0.0-20231128154154-6736c9b9c6c8
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils => github.com/openshift/baremetal-operator/pkg/hardwareutils v0.0.0-20231128154154-6736c9b9c6c8
	k8s.io/cloud-provider-vsphere => github.com/openshift/cloud-provider-vsphere v1.19.1-0.20240626105621-6464d0bb4928
// sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.5
// sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20200929152424-eab2e087f366 // Indirect dependency through MAO from cluster API providers
)

replace github.com/openshift/assisted-service/api => github.com/openshift/assisted-service/api v0.0.0-20250625193139-474abcbea19a

replace github.com/openshift/assisted-service/client => github.com/openshift/assisted-service/client v0.0.0-20241001055825-63e8b0d3ad63

replace github.com/openshift/assisted-service/models => github.com/openshift/assisted-service/models v0.0.0-20241001055825-63e8b0d3ad63

// https://issues.redhat.com/browse/OCPBUGS-8119
// https://issues.redhat.com/browse/OCPBUGS-27507
replace github.com/containerd/containerd => github.com/containerd/containerd v1.6.26

replace github.com/vmware-tanzu/vm-operator/pkg/constants/testlabels => github.com/vmware-tanzu/vm-operator/pkg/constants/testlabels v0.0.0-20240404200847-de75746a9505

// This is to force capi back for the older provider version
replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.19.3
