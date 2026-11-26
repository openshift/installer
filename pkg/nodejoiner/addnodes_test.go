package nodejoiner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	faketesting "k8s.io/client-go/testing"

	configv1 "github.com/openshift/api/config/v1"
	fakeclientconfig "github.com/openshift/client-go/config/clientset/versioned/fake"
)

const fakeProxyCACert = "-----BEGIN CERTIFICATE-----\nfake-proxy-ca-data\n-----END CERTIFICATE-----\n"

func TestSetupProxyCACert(t *testing.T) {
	testCases := []struct {
		name            string
		configObjects   []runtime.Object
		kubeObjects     []runtime.Object
		systemCAContent string
		expectedSSLFile bool
		expectedContent string
		expectedError   string
	}{
		{
			name:            "no proxy configured",
			expectedSSLFile: false,
		},
		{
			name: "proxy without trusted CA",
			configObjects: []runtime.Object{
				&configv1.Proxy{
					ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
				},
			},
			expectedSSLFile: false,
		},
		{
			name: "proxy with trusted CA, no system bundle",
			configObjects: []runtime.Object{
				proxyWithTrustedCA("user-ca-bundle"),
			},
			kubeObjects: []runtime.Object{
				caConfigMap("user-ca-bundle", fakeProxyCACert),
			},
			expectedSSLFile: true,
			expectedContent: fakeProxyCACert,
		},
		{
			name: "proxy with trusted CA, system bundle concatenated",
			configObjects: []runtime.Object{
				proxyWithTrustedCA("user-ca-bundle"),
			},
			kubeObjects: []runtime.Object{
				caConfigMap("user-ca-bundle", fakeProxyCACert),
			},
			systemCAContent: "system-ca-data\n",
			expectedSSLFile: true,
			expectedContent: "system-ca-data\n" + fakeProxyCACert,
		},
		{
			name: "proxy with trusted CA, system bundle missing trailing newline",
			configObjects: []runtime.Object{
				proxyWithTrustedCA("user-ca-bundle"),
			},
			kubeObjects: []runtime.Object{
				caConfigMap("user-ca-bundle", fakeProxyCACert),
			},
			systemCAContent: "system-ca-data-no-newline",
			expectedSSLFile: true,
			expectedContent: "system-ca-data-no-newline\n" + fakeProxyCACert,
		},
		{
			name: "proxy with trusted CA, configmap not found",
			configObjects: []runtime.Object{
				proxyWithTrustedCA("user-ca-bundle"),
			},
			expectedSSLFile: false,
		},
		{
			name: "proxy with trusted CA, ca-bundle.crt key missing",
			configObjects: []runtime.Object{
				proxyWithTrustedCA("user-ca-bundle"),
			},
			kubeObjects: []runtime.Object{
				&corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{Name: "user-ca-bundle", Namespace: "openshift-config"},
					Data:       map[string]string{},
				},
			},
			expectedSSLFile: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			// Point systemCACertBundle at a temp file so tests don't read the real system bundle.
			fakeSystemBundle := filepath.Join(dir, "system-ca-bundle.pem")
			original := systemCACertBundle
			systemCACertBundle = fakeSystemBundle
			defer func() { systemCACertBundle = original }()

			if tc.systemCAContent != "" {
				if err := os.WriteFile(fakeSystemBundle, []byte(tc.systemCAContent), 0600); !assert.NoError(t, err) {
					return
				}
			}

			fakeConfigClient := fakeclientconfig.NewClientset(tc.configObjects...)
			fakeK8sClient := fake.NewClientset(tc.kubeObjects...)

			// Return a proper not-found error for configmaps not in the test's kubeObjects.
			fakeK8sClient.PrependReactor("get", "configmaps", func(action faketesting.Action) (bool, runtime.Object, error) {
				name := action.(faketesting.GetAction).GetName()
				for _, obj := range tc.kubeObjects {
					if cm, ok := obj.(*corev1.ConfigMap); ok && cm.Name == name {
						return false, nil, nil // let the default handler serve it
					}
				}
				return true, nil, k8serrors.NewNotFound(schema.GroupResource{Resource: "configmaps"}, name)
			})

			os.Unsetenv("SSL_CERT_FILE")

			err := setupProxyCACertWithClients(dir, fakeConfigClient, fakeK8sClient)

			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
				return
			}
			if !assert.NoError(t, err) {
				return
			}

			caFile := filepath.Join(dir, "proxy-ca-bundle.crt")
			if tc.expectedSSLFile {
				assert.Equal(t, caFile, os.Getenv("SSL_CERT_FILE"))
				content, readErr := os.ReadFile(caFile)
				if assert.NoError(t, readErr) {
					assert.Equal(t, tc.expectedContent, string(content))
				}
			} else {
				assert.Empty(t, os.Getenv("SSL_CERT_FILE"))
				_, statErr := os.Stat(caFile)
				assert.True(t, os.IsNotExist(statErr), "expected no CA bundle file to be written")
			}
		})
	}
}

func proxyWithTrustedCA(configMapName string) *configv1.Proxy {
	return &configv1.Proxy{
		ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
		Spec: configv1.ProxySpec{
			TrustedCA: configv1.ConfigMapNameReference{Name: configMapName},
		},
	}
}

func caConfigMap(name, cert string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "openshift-config"},
		Data:       map[string]string{"ca-bundle.crt": cert},
	}
}
