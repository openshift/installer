package stock

import (
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/metadata"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of assets that can be generated.
type Stock struct {
	installConfigStock
	kubeconfigStock
	tlsStock
	bootstrapStock
	machineStock
	clusterStock
	manifestsStock
	metadataStock
}

type installConfigStock struct {
	installconfig.StockImpl
}

type tlsStock struct {
	tls.StockImpl
}

type kubeconfigStock struct {
	kubeconfig.StockImpl
}

type bootstrapStock struct {
	bootstrap.StockImpl
}

type machineStock struct {
	machine.StockImpl
}

type clusterStock struct {
	cluster.StockImpl
}

type manifestsStock struct {
	manifests.StockImpl
}

type metadataStock struct {
	metadata.StockImpl
}

var _ installconfig.Stock = (*Stock)(nil)

// EstablishStock establishes the stock of assets.
func EstablishStock() *Stock {
	s := &Stock{}
	s.installConfigStock.EstablishStock()
	s.tlsStock.EstablishStock(&s.installConfigStock)
	s.kubeconfigStock.EstablishStock(&s.installConfigStock, &s.tlsStock)
	s.machineStock.EstablishStock(s, s)
	s.manifestsStock.EstablishStock(&s.installConfigStock, s, s, s)
	s.bootstrapStock.EstablishStock(s, s, s, s)
	s.clusterStock.EstablishStock(s, s, s, s)
	s.metadataStock.EstablishStock(s, s)
	return s
}
