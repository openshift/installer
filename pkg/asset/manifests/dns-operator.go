package manifests

import (
	"path/filepath"
	"sort"

	"github.com/ghodss/yaml"

	dnsopapi "github.com/openshift/cluster-dns-operator/pkg/apis/dns/v1alpha1"
	dnsopmanifests "github.com/openshift/cluster-dns-operator/pkg/manifests"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

// dnsOperator generates the dns-operator-*.yml files
type dnsOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
}

var _ asset.Asset = (*dnsOperator)(nil)

// Name returns a human friendly name for the operator
func (doh *dnsOperator) Name() string {
	return "DNS Operator"
}

// Dependencies returns all of the dependencies directly needed by an
// dnsOperator asset.
func (doh *dnsOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		doh.installConfigAsset,
	}
}

// Generate generates the dns-operator-config.yml and dns-operator-manifests.yml files
func (doh *dnsOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(doh.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	doh.installConfig = ic

	// installconfig is ready, we can create the core config from it now
	dnsConfig, err := doh.dnsConfig()
	if err != nil {
		return nil, err
	}

	assetData, err := doh.assetData()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for k := range assetData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	assetContents := make([]asset.Content, 0)
	for _, k := range keys {
		assetContents = append(assetContents, asset.Content{
			Name: filepath.Join("dns-operator", k),
			Data: assetData[k],
		})
	}

	assetContents = append(assetContents, asset.Content{
		Name: "dns-operator-config.yml",
		Data: dnsConfig,
	})

	return &asset.State{Contents: assetContents}, nil
}

func (doh *dnsOperator) dnsOperatorConfig() (*dnsopapi.ClusterDNS, error) {
	clusterIP, err := installconfig.ClusterDNSIP(doh.installConfig)
	if err != nil {
		return nil, err
	}

	return &dnsopapi.ClusterDNS{
		Spec: dnsopapi.ClusterDNSSpec{
			// Check if BaseDomain is correct?
			ClusterIP:     &clusterIP,
			ClusterDomain: &doh.installConfig.BaseDomain,
		},
	}, nil
}

func (doh *dnsOperator) dnsConfig() ([]byte, error) {
	return yaml.Marshal(doh.dnsOperatorConfig())
}

func (doh *dnsOperator) assetData() (map[string][]byte, error) {
	f := dnsopmanifests.NewFactory()
	return f.AssetContent(doh.dnsOperatorConfig())
}
