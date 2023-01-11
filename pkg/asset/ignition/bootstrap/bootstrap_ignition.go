package bootstrap

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
)

// parseCertificateBundle loads each certificate in the bundle to the Ingition
// carrier type, ignoring any invisible character before, after and in between
// certificates.
func parseCertificateBundle(userCA []byte) ([]igntypes.Resource, error) {
	userCA = bytes.TrimSpace(userCA)

	var carefs []igntypes.Resource
	for len(userCA) > 0 {
		var block *pem.Block
		block, userCA = pem.Decode(userCA)
		if block == nil {
			return nil, fmt.Errorf("unable to parse certificate, please check the certificates")
		}

		carefs = append(carefs, igntypes.Resource{Source: ignutil.StrToPtr(dataurl.EncodeBytes(pem.EncodeToMemory(block)))})

		userCA = bytes.TrimSpace(userCA)
	}

	return carefs, nil
}

// GenerateIgnitionShimWithCertBundleAndProxy is used to generate an ignition file that contains both a user ca bundle
// in its Security section and proxy settings (if any).
func GenerateIgnitionShimWithCertBundleAndProxy(bootstrapConfigURL string, userCA string, proxy *types.Proxy) ([]byte, error) {
	ign := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Replace: igntypes.Resource{
					Source: ignutil.StrToPtr(bootstrapConfigURL),
				},
			},
		},
	}

	carefs, err := parseCertificateBundle([]byte(userCA))
	if err != nil {
		return nil, err
	}
	if len(carefs) > 0 {
		ign.Ignition.Security = igntypes.Security{
			TLS: igntypes.TLS{
				CertificateAuthorities: carefs,
			},
		}
	}

	if proxy != nil {
		ign.Ignition.Proxy = ignitionProxy(proxy)
	}

	data, err := ignition.Marshal(ign)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ignitionProxy(proxy *types.Proxy) igntypes.Proxy {
	var ignProxy igntypes.Proxy
	if proxy == nil {
		return ignProxy
	}
	if httpProxy := proxy.HTTPProxy; httpProxy != "" {
		ignProxy.HTTPProxy = &httpProxy
	}
	if httpsProxy := proxy.HTTPSProxy; httpsProxy != "" {
		ignProxy.HTTPSProxy = &httpsProxy
	}
	ignProxy.NoProxy = make([]igntypes.NoProxyItem, 0, len(proxy.NoProxy))
	if noProxy := proxy.NoProxy; noProxy != "" {
		noProxySplit := strings.Split(noProxy, ",")
		for _, p := range noProxySplit {
			ignProxy.NoProxy = append(ignProxy.NoProxy, igntypes.NoProxyItem(p))
		}
	}
	return ignProxy
}
