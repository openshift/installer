package tls

import (
	"fmt"
	"net"
	"path"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/types"
)

const (
	tlsDir = "tls"
)

func assetFilePath(filename string) string {
	return path.Join(tlsDir, filename)
}

func apiAddress(cfg *types.InstallConfig) string {
	return fmt.Sprintf("api.%s", cfg.ClusterDomain())
}

func internalAPIAddress(cfg *types.InstallConfig) string {
	return fmt.Sprintf("api-int.%s", cfg.ClusterDomain())
}

func cidrhost(network net.IPNet, hostNum int) (string, error) {
	ip, err := cidr.Host(&network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
