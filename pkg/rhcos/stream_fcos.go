//go:build okd || fcos

package rhcos

func getStreamFileName() string {
	return "coreos/fcos.json"
}

func getMarketplaceStreamFileName() string {
	// There is no need for OKD marketplace images at this time
	// so we can skip reading a marketplace stream file.
	return ""
}
