module github.com/openshift/installer/terraform/providers/vsphereprivate

go 1.18

require (
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/hashicorp/terraform-provider-vsphere v1.26.1-0.20220308165248-f78c6d6f1da6
	github.com/pkg/errors v0.9.1
	github.com/vmware/govmomi v0.27.1
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.18.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.12.0 // indirect
	cloud.google.com/go/storage v1.29.0 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/apparentlymart/go-textseg/v12 v12.0.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go v1.44.206 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-checkpoint v0.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 // indirect
	github.com/hashicorp/go-getter v1.7.0 // indirect
	github.com/hashicorp/go-hclog v1.2.1 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/go-plugin v1.4.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/logutils v1.0.0 // indirect
	github.com/hashicorp/terraform-exec v0.14.0 // indirect
	github.com/hashicorp/terraform-json v0.12.0 // indirect
	github.com/hashicorp/terraform-plugin-go v0.3.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.15 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/zclconf/go-cty v1.8.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.5.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.110.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230221151758-ace64dc21148 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

// https://issues.redhat.com/browse/OCPBUGS-7699
replace github.com/hashicorp/go-getter => github.com/hashicorp/go-getter v1.7.0

// https://issues.redhat.com/browse/OCPBUGS-5667
replace github.com/Masterminds/goutils => github.com/Masterminds/goutils v1.1.1
