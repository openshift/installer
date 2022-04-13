module github.com/openshift-agent-team/fleeting

go 1.16

require (
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/coreos/ignition/v2 v2.13.0
	github.com/diskfs/go-diskfs v1.2.1-0.20210727185522-a769efacd235 // indirect
	github.com/go-openapi/strfmt v0.21.2
	github.com/go-openapi/swag v0.21.1
	github.com/go-openapi/validate v0.20.3 // indirect
	github.com/openshift/assisted-image-service v0.0.0-20220307202600-054a1afa8d28
	github.com/openshift/assisted-service v1.0.10-0.20220223093655-7ada9949bf1d
	github.com/openshift/hive/apis v0.0.0-20210506000654-5c038fb05190
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/thoas/go-funk v0.9.1
	github.com/vincent-petithory/dataurl v1.0.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	sigs.k8s.io/structured-merge-diff/v4 v4.1.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	sigs.k8s.io/cluster-api-provider-aws => github.com/openshift/cluster-api-provider-aws v0.2.1-0.20201022175424-d30c7a274820
	sigs.k8s.io/cluster-api-provider-azure => github.com/openshift/cluster-api-provider-azure v0.1.0-alpha.3.0.20201016155852-4090a6970205
)
