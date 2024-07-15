package configimage

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/tls"
)

func TestCaBundle_Generate(t *testing.T) {
	expectedBundleRaw := bytes.Join([][]byte{
		lbCABundle().BundleRaw,
		localhostCABundle().BundleRaw,
		serviceNetworkCABundle().BundleRaw,
		ingressCABundle().BundleRaw,
	}, []byte{})

	cases := []struct {
		name         string
		dependencies []asset.Asset
		expected     *tls.CertBundle
	}{
		{
			name: "valid dependencies",
			dependencies: []asset.Asset{
				lbCABundle(),
				localhostCABundle(),
				serviceNetworkCABundle(),
				ingressCABundle(),
			},
			expected: &tls.CertBundle{
				BundleRaw: expectedBundleRaw,
				FileList: []*asset.File{
					{
						Filename: "tls/kube-apiserver-complete-server-ca-bundle.crt",
						Data:     expectedBundleRaw,
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &ImageBasedKubeAPIServerCompleteCABundle{}
			err := asset.Generate(context.TODO(), parents)
			assert.NoError(t, err)
			assert.Equal(t, string(tc.expected.BundleRaw), string(asset.CertBundle.BundleRaw))
			assert.Equal(t, tc.expected.FileList, asset.CertBundle.FileList)
		})
	}
}

func lbCABundle() *tls.KubeAPIServerLBCABundle {
	return &tls.KubeAPIServerLBCABundle{
		CertBundle: tls.CertBundle{
			BundleRaw: []byte(testCert),
		},
	}
}

func localhostCABundle() *tls.KubeAPIServerLocalhostCABundle {
	return &tls.KubeAPIServerLocalhostCABundle{
		CertBundle: tls.CertBundle{
			BundleRaw: []byte(testCert),
		},
	}
}

func serviceNetworkCABundle() *tls.KubeAPIServerServiceNetworkCABundle {
	return &tls.KubeAPIServerServiceNetworkCABundle{
		CertBundle: tls.CertBundle{
			BundleRaw: []byte(testCert),
		},
	}
}

func ingressCABundle() *IngressOperatorCABundle {
	return &IngressOperatorCABundle{
		CertBundle: tls.CertBundle{
			BundleRaw: []byte(testCert),
		},
	}
}
