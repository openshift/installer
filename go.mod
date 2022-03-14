module github.com/openshift-agent-team/fleeting

go 1.16

require (
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/ignition/v2 v2.9.0
	github.com/go-openapi/strfmt v0.21.1
	github.com/go-openapi/swag v0.19.15
	github.com/openshift/assisted-image-service v0.0.0-20220307202600-054a1afa8d28
	github.com/openshift/assisted-service v1.0.10-0.20220116113517-db25501e204a
	github.com/openshift/hive/apis v0.0.0-20210506000654-5c038fb05190
	github.com/thoas/go-funk v0.9.1
	github.com/vincent-petithory/dataurl v1.0.0
	k8s.io/api v0.21.1
	sigs.k8s.io/yaml v1.3.0
)

replace (
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20201022175424-d30c7a274820
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20201016155852-4090a6970205
)
