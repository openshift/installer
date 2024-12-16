//go:build !(okd || fcos || scos)

package rhcos

func getStreamFileName() string {
	return "coreos/rhcos.json"
}

func getMarketplaceStreamFileName() string {
	return "coreos/marketplace-rhcos.json"
}
