package stock

import (
	"bufio"
	"os"

	"github.com/openshift/installer/pkg/asset/installconfig"
)

// Stock is the stock of assets that can be generated.
type Stock struct {
	installConfigStock
}

type installConfigStock struct {
	installconfig.StockImpl
}

var _ installconfig.Stock = (*Stock)(nil)

// EstablishStock establishes the stock of assets in the specified directory.
func EstablishStock(directory string) *Stock {
	s := &Stock{}
	inputReader := bufio.NewReader(os.Stdin)
	s.installConfigStock.EstablishStock(directory, inputReader)
	return s
}
