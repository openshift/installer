package tls

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/apparentlymart/go-cidr/cidr"

	"github.com/openshift/installer/pkg/types"
)

const (
	tlsDir = "tls"
)

func assetFilePath(filename string) string {
	return filepath.Join(tlsDir, filename)
}

func apiAddress(cfg *types.InstallConfig) string {
	return fmt.Sprintf("%s-api.%s", cfg.ObjectMeta.Name, cfg.BaseDomain)
}

func cidrhost(network net.IPNet, hostNum int) (string, error) {
	ip, err := cidr.Host(&network, hostNum)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
