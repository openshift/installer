// Package operators deals with creating assets for all operators to be installed for the cluster
package operators

import (
	"github.com/openshift/installer/pkg/asset"
	"path/filepath"
)

// operators generates the dependent operator config.yaml files
type operators struct {
	assetStock Stock
	directory  string
}

var _ asset.Asset = (*operators)(nil)

type genericData map[string]string

// Dependencies returns all of the dependencies directly needed by an
// operators asset.
func (o *operators) Dependencies() []asset.Asset {
	return []asset.Asset{
		o.assetStock.KubeCoreOperator(),
		o.assetStock.NetworkOperator(),
		o.assetStock.Tnco(),
		o.assetStock.KubeAddonOperator(),
	}
}

// Generate generates the respective operator config.yml files
func (o *operators) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	//cvo := dependencies[o.assetStock.ClusterVersionOperator()].Contents[0]
	kco := dependencies[o.assetStock.KubeCoreOperator()].Contents[0]
	no := dependencies[o.assetStock.NetworkOperator()].Contents[0]
	tnco := dependencies[o.assetStock.Tnco()].Contents[0]
	//ingress := dependencies[o.assetStock.IngressOperator()].Contents[0]
	addon := dependencies[o.assetStock.KubeAddonOperator()].Contents[0]

	// kco+no+tnco go to kube-system config map
	kubeSys, err := configMap("kube-system", "cluster-config-v1", genericData{
		"kco-config":     string(kco.Data),
		"network-config": string(no.Data),
		"tnco-config":    string(tnco.Data),
	})
	if err != nil {
		return nil, err
	}

	// addon goes to openshift system
	openshiftSys, err := configMap("openshift-system", "cluster-config-v1", genericData{
		"addon-config": string(addon.Data),
	})
	if err != nil {
		return nil, err
	}

	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: filepath.Join(o.directory, "cluster-config.yml"),
				Data: []byte(kubeSys),
			},
			{
				Name: filepath.Join(o.directory, "openshift-config.yml"),
				Data: []byte(openshiftSys),
			},
		},
	}
	return state, nil
}
