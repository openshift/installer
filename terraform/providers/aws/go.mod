module github.com/openshift/installer/terraform/providers/aws

go 1.18

require github.com/hashicorp/terraform-provider-aws v1.60.1-0.20230615210323-24881e9f4460 // v5.4.0

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230619160724-3fbb1f12458c // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go v1.44.287 // indirect
	github.com/aws/aws-sdk-go-v2 v1.18.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.18.27 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.26 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/accessanalyzer v1.19.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/account v1.10.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/acm v1.17.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/auditmanager v1.25.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/cleanrooms v1.1.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudcontrol v1.11.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.21.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/comprehend v1.24.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/computeoptimizer v1.24.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/docdbelastic v1.1.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.102.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/finspace v1.10.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/fis v1.14.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/glacier v1.14.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/healthlake v1.16.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.20.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/identitystore v1.16.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/inspector2 v1.15.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.28 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ivschat v1.4.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/kendra v1.40.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.36.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.27.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/medialive v1.31.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/oam v1.1.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/opensearchserverless v1.2.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/pipes v1.2.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/rbin v1.8.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.45.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/resourceexplorer2 v1.2.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/rolesanywhere v1.2.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.15.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3control v1.31.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/scheduler v1.1.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/securitylake v1.4.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.36.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssmcontacts v1.15.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssmincidents v1.21.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.19.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/swf v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/transcribe v1.26.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/vpclattice v1.0.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/xray v1.16.13 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/beevik/etree v1.2.0 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/aws-cloudformation-resource-schema-sdk-go v0.20.0 // indirect
	github.com/hashicorp/aws-sdk-go-base/v2 v2.0.0-beta.29 // indirect
	github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2 v2.0.0-beta.30 // indirect
	github.com/hashicorp/awspolicyequivalence v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200723130312-85980079f637 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.10 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hc-install v0.5.2 // indirect
	github.com/hashicorp/hcl/v2 v2.17.0 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-exec v0.18.1 // indirect
	github.com/hashicorp/terraform-json v0.17.0 // indirect
	github.com/hashicorp/terraform-plugin-framework v1.3.1 // indirect
	github.com/hashicorp/terraform-plugin-framework-timeouts v0.4.0 // indirect
	github.com/hashicorp/terraform-plugin-framework-validators v0.10.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.16.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/hashicorp/terraform-plugin-mux v0.10.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.26.1 // indirect
	github.com/hashicorp/terraform-plugin-testing v1.3.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.2.1 // indirect
	github.com/hashicorp/terraform-svchost v0.1.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20230413205102-771768614e91 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/zclconf/go-cty v1.13.2 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.56.1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// https://bugzilla.redhat.com/show_bug.cgi?id=2064702
replace golang.org/x/crypto => golang.org/x/crypto v0.0.0-20221010152910-d6f0a8c073c2

// https://bugzilla.redhat.com/show_bug.cgi?id=2100495
replace golang.org/x/text => golang.org/x/text v0.3.7

// https://issues.redhat.com/browse/OCPBUGS-6422
replace golang.org/x/net => golang.org/x/net v0.5.0
