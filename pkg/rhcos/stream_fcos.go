//go:build okd || fcos

package rhcos

func getStreamFileName() string {
	return "coreos/fcos.json"
}

func getMarketplaceStreamFileName() string {
	return "coreos/marketplace-fcos.json"
}
