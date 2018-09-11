package stock

import (
	"bufio"
	"os"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/tls"
)

// Stock is the stock of assets that can be generated.
type Stock struct {
	installConfigStock
	kubeconfigStock
	tlsStock
	ignitionStock
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

type ignitionStock struct {
	ignition.StockImpl
}

var _ installconfig.Stock = (*Stock)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func EstablishStock(directory string) *Stock {
	s := &Stock{}
	inputReader := bufio.NewReader(os.Stdin)
	s.installConfigStock.EstablishStock(directory, inputReader)
	s.tlsStock.EstablishStock(directory, &s.installConfigStock)
	s.kubeconfigStock.EstablishStock(directory, &s.installConfigStock, &s.tlsStock)
	s.ignitionStock.EstablishStock(directory, s, s, s)

	return s
}
