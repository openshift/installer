//go:build !(okd || fcos || scos)

package rhcos

func getStreamFileName() string {
	return "coreos/rhcos.json"
}
