package mirror

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/mock"
	"github.com/openshift/installer/pkg/types"
)

func TestCaBundle_Generate(t *testing.T) {

	cases := []struct {
		name           string
		dependencies   []asset.Asset
		expectedError  string
		expectedConfig string
	}{
		{
			name: "missing-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				&agent.OptionalInstallConfig{},
			},
		},
		{
			name: "default",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				&agent.OptionalInstallConfig{
					Supplied: true,
					AssetBase: installconfig.AssetBase{
						Config: &types.InstallConfig{
							ObjectMeta: v1.ObjectMeta{
								Namespace: "cluster-0",
							},
						},
					},
				},
			},
		},
		{
			name: "additional-trust-bundle",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeInstall},
				&joiner.ClusterInfo{},
				&agent.OptionalInstallConfig{
					Supplied: true,
					AssetBase: installconfig.AssetBase{
						Config: &types.InstallConfig{
							ObjectMeta: v1.ObjectMeta{
								Namespace: "cluster-0",
							},
							AdditionalTrustBundle: `
-----BEGIN CERTIFICATE-----
MIIDZTCCAk2gAwIBAgIURbA8lR+5xlJZUoOXK66AHFWd3uswDQYJKoZIhvcNAQEL
BQAwQjELMAkGA1UEBhMCWFgxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UE
CgwTRGVmYXVsdCBDb21wYW55IEx0ZDAeFw0yMjA3MDgxOTUzMTVaFw0yMjA4MDcx
OTUzMTVaMEIxCzAJBgNVBAYTAlhYMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAa
BgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCroH9c2PLWI0O/nBrmKtS2IuReyWaR0DOMJY7C/vc12l9zlH0D
xTOUfEtdqRktjVsUn1vIIiFakxd0QLIPcMyKplmbavIBUQp+MZr0pNVX+lwcctbA
7FVHEnbWYNVepoV7kZkTVvMXAqFylMXU4gDmuZzIxhVMMxjialJNED+3ngqvX4w3
4q4KSk1ytaHGwjREIErwPJjv5PK48KVJL2nlCuA+tbxu1r8eVkOUvZlxAuNNXk/U
mf3QX5EiUlTtsmRAct6fIUT3jkrsHSS/tZ66EYJ9Q0OBoX2lL/Msmi27OQvA7uYn
uqYlwJzU43tCsiip9E9z/UrLcMYyXx3oPJyPAgMBAAGjUzBRMB0GA1UdDgQWBBTI
ahE8DDT4T1vta6cXVVaRjnel0zAfBgNVHSMEGDAWgBTIahE8DDT4T1vta6cXVVaR
jnel0zAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCQbsMtPFkq
PxwOAIds3IoupuyIKmsF32ECEH/OlS+7Sj7MUJnGTQrwgjrsVS5sl8AmnGx4hPdL
VX98nEcKMNkph3Hkvh4EvgjSfmYGUXuJBcYU5jqNQrlrGv37rEf5FnvdHV1F3MG8
A0Mj0TLtcTdtaJFoOrnQuD/k0/1d+cMiYGTSaT5XK/unARqGEMd4BlWPh5P3SflV
/Vy2hHlMpv7OcZ8yaAI3htENZLus+L5kjHWKu6dxlPHKu6ef5k64su2LTNE07Vr9
S655uiFW5AX2wDVUcQEDCOiEn6SI9DTt5oQjWPMxPf+rEyfQ2f1QwVez7cyr6Qc5
OIUk31HnM/Fj
-----END CERTIFICATE-----
`,
						},
					},
				},
			},
			expectedConfig: `-----BEGIN CERTIFICATE-----
MIIDZTCCAk2gAwIBAgIURbA8lR+5xlJZUoOXK66AHFWd3uswDQYJKoZIhvcNAQEL
BQAwQjELMAkGA1UEBhMCWFgxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UE
CgwTRGVmYXVsdCBDb21wYW55IEx0ZDAeFw0yMjA3MDgxOTUzMTVaFw0yMjA4MDcx
OTUzMTVaMEIxCzAJBgNVBAYTAlhYMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAa
BgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCroH9c2PLWI0O/nBrmKtS2IuReyWaR0DOMJY7C/vc12l9zlH0D
xTOUfEtdqRktjVsUn1vIIiFakxd0QLIPcMyKplmbavIBUQp+MZr0pNVX+lwcctbA
7FVHEnbWYNVepoV7kZkTVvMXAqFylMXU4gDmuZzIxhVMMxjialJNED+3ngqvX4w3
4q4KSk1ytaHGwjREIErwPJjv5PK48KVJL2nlCuA+tbxu1r8eVkOUvZlxAuNNXk/U
mf3QX5EiUlTtsmRAct6fIUT3jkrsHSS/tZ66EYJ9Q0OBoX2lL/Msmi27OQvA7uYn
uqYlwJzU43tCsiip9E9z/UrLcMYyXx3oPJyPAgMBAAGjUzBRMB0GA1UdDgQWBBTI
ahE8DDT4T1vta6cXVVaRjnel0zAfBgNVHSMEGDAWgBTIahE8DDT4T1vta6cXVVaR
jnel0zAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCQbsMtPFkq
PxwOAIds3IoupuyIKmsF32ECEH/OlS+7Sj7MUJnGTQrwgjrsVS5sl8AmnGx4hPdL
VX98nEcKMNkph3Hkvh4EvgjSfmYGUXuJBcYU5jqNQrlrGv37rEf5FnvdHV1F3MG8
A0Mj0TLtcTdtaJFoOrnQuD/k0/1d+cMiYGTSaT5XK/unARqGEMd4BlWPh5P3SflV
/Vy2hHlMpv7OcZ8yaAI3htENZLus+L5kjHWKu6dxlPHKu6ef5k64su2LTNE07Vr9
S655uiFW5AX2wDVUcQEDCOiEn6SI9DTt5oQjWPMxPf+rEyfQ2f1QwVez7cyr6Qc5
OIUk31HnM/Fj
-----END CERTIFICATE-----
`,
		},

		{
			name: "add-nodes command - missing-config",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{},
				&agent.OptionalInstallConfig{},
			},
		},
		{
			name: "add-nodes command - default",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{
					Namespace: "cluster-0",
				},
				&agent.OptionalInstallConfig{},
			},
		},
		{
			name: "add-nodes command - additional-trust-bundle",
			dependencies: []asset.Asset{
				&workflow.AgentWorkflow{Workflow: workflow.AgentWorkflowTypeAddNodes},
				&joiner.ClusterInfo{
					Namespace: "cluster-0",
					UserCaBundle: `-----BEGIN CERTIFICATE-----
MIIDZTCCAk2gAwIBAgIURbA8lR+5xlJZUoOXK66AHFWd3uswDQYJKoZIhvcNAQEL
BQAwQjELMAkGA1UEBhMCWFgxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UE
CgwTRGVmYXVsdCBDb21wYW55IEx0ZDAeFw0yMjA3MDgxOTUzMTVaFw0yMjA4MDcx
OTUzMTVaMEIxCzAJBgNVBAYTAlhYMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAa
BgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCroH9c2PLWI0O/nBrmKtS2IuReyWaR0DOMJY7C/vc12l9zlH0D
xTOUfEtdqRktjVsUn1vIIiFakxd0QLIPcMyKplmbavIBUQp+MZr0pNVX+lwcctbA
7FVHEnbWYNVepoV7kZkTVvMXAqFylMXU4gDmuZzIxhVMMxjialJNED+3ngqvX4w3
4q4KSk1ytaHGwjREIErwPJjv5PK48KVJL2nlCuA+tbxu1r8eVkOUvZlxAuNNXk/U
mf3QX5EiUlTtsmRAct6fIUT3jkrsHSS/tZ66EYJ9Q0OBoX2lL/Msmi27OQvA7uYn
uqYlwJzU43tCsiip9E9z/UrLcMYyXx3oPJyPAgMBAAGjUzBRMB0GA1UdDgQWBBTI
ahE8DDT4T1vta6cXVVaRjnel0zAfBgNVHSMEGDAWgBTIahE8DDT4T1vta6cXVVaR
jnel0zAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCQbsMtPFkq
PxwOAIds3IoupuyIKmsF32ECEH/OlS+7Sj7MUJnGTQrwgjrsVS5sl8AmnGx4hPdL
VX98nEcKMNkph3Hkvh4EvgjSfmYGUXuJBcYU5jqNQrlrGv37rEf5FnvdHV1F3MG8
A0Mj0TLtcTdtaJFoOrnQuD/k0/1d+cMiYGTSaT5XK/unARqGEMd4BlWPh5P3SflV
/Vy2hHlMpv7OcZ8yaAI3htENZLus+L5kjHWKu6dxlPHKu6ef5k64su2LTNE07Vr9
S655uiFW5AX2wDVUcQEDCOiEn6SI9DTt5oQjWPMxPf+rEyfQ2f1QwVez7cyr6Qc5
OIUk31HnM/Fj
-----END CERTIFICATE-----
`,
				},
				&agent.OptionalInstallConfig{},
			},
			expectedConfig: `-----BEGIN CERTIFICATE-----
MIIDZTCCAk2gAwIBAgIURbA8lR+5xlJZUoOXK66AHFWd3uswDQYJKoZIhvcNAQEL
BQAwQjELMAkGA1UEBhMCWFgxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UE
CgwTRGVmYXVsdCBDb21wYW55IEx0ZDAeFw0yMjA3MDgxOTUzMTVaFw0yMjA4MDcx
OTUzMTVaMEIxCzAJBgNVBAYTAlhYMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAa
BgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCroH9c2PLWI0O/nBrmKtS2IuReyWaR0DOMJY7C/vc12l9zlH0D
xTOUfEtdqRktjVsUn1vIIiFakxd0QLIPcMyKplmbavIBUQp+MZr0pNVX+lwcctbA
7FVHEnbWYNVepoV7kZkTVvMXAqFylMXU4gDmuZzIxhVMMxjialJNED+3ngqvX4w3
4q4KSk1ytaHGwjREIErwPJjv5PK48KVJL2nlCuA+tbxu1r8eVkOUvZlxAuNNXk/U
mf3QX5EiUlTtsmRAct6fIUT3jkrsHSS/tZ66EYJ9Q0OBoX2lL/Msmi27OQvA7uYn
uqYlwJzU43tCsiip9E9z/UrLcMYyXx3oPJyPAgMBAAGjUzBRMB0GA1UdDgQWBBTI
ahE8DDT4T1vta6cXVVaRjnel0zAfBgNVHSMEGDAWgBTIahE8DDT4T1vta6cXVVaR
jnel0zAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCQbsMtPFkq
PxwOAIds3IoupuyIKmsF32ECEH/OlS+7Sj7MUJnGTQrwgjrsVS5sl8AmnGx4hPdL
VX98nEcKMNkph3Hkvh4EvgjSfmYGUXuJBcYU5jqNQrlrGv37rEf5FnvdHV1F3MG8
A0Mj0TLtcTdtaJFoOrnQuD/k0/1d+cMiYGTSaT5XK/unARqGEMd4BlWPh5P3SflV
/Vy2hHlMpv7OcZ8yaAI3htENZLus+L5kjHWKu6dxlPHKu6ef5k64su2LTNE07Vr9
S655uiFW5AX2wDVUcQEDCOiEn6SI9DTt5oQjWPMxPf+rEyfQ2f1QwVez7cyr6Qc5
OIUk31HnM/Fj
-----END CERTIFICATE-----
`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			parents := asset.Parents{}
			parents.Add(tc.dependencies...)

			asset := &CaBundle{}
			err := asset.Generate(context.Background(), parents)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)

				files := asset.Files()
				if tc.expectedConfig != "" {
					assert.Len(t, files, 1)
					assert.Equal(t, CaBundleFilename, files[0].Filename)
					assert.Equal(t, tc.expectedConfig, string(files[0].Data))
				} else {
					if len(files) == 1 {
						assert.Equal(t, CaBundleFilename, files[0].Filename)
						assert.Equal(t, []byte{}, files[0].Data)
					} else {
						assert.Empty(t, files)
					}
				}
			}
		})
	}
}

func TestCaBundle_LoadedFromDisk(t *testing.T) {

	cases := []struct {
		name          string
		data          string
		fetchError    error
		expectedFound bool
		expectedError string
	}{
		{
			name: "valid-config-file",
			data: `
-----BEGIN CERTIFICATE-----
MIIDZTCCAk2gAwIBAgIURbA8lR+5xlJZUoOXK66AHFWd3uswDQYJKoZIhvcNAQEL
BQAwQjELMAkGA1UEBhMCWFgxFTATBgNVBAcMDERlZmF1bHQgQ2l0eTEcMBoGA1UE
CgwTRGVmYXVsdCBDb21wYW55IEx0ZDAeFw0yMjA3MDgxOTUzMTVaFw0yMjA4MDcx
OTUzMTVaMEIxCzAJBgNVBAYTAlhYMRUwEwYDVQQHDAxEZWZhdWx0IENpdHkxHDAa
BgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCroH9c2PLWI0O/nBrmKtS2IuReyWaR0DOMJY7C/vc12l9zlH0D
xTOUfEtdqRktjVsUn1vIIiFakxd0QLIPcMyKplmbavIBUQp+MZr0pNVX+lwcctbA
7FVHEnbWYNVepoV7kZkTVvMXAqFylMXU4gDmuZzIxhVMMxjialJNED+3ngqvX4w3
4q4KSk1ytaHGwjREIErwPJjv5PK48KVJL2nlCuA+tbxu1r8eVkOUvZlxAuNNXk/U
mf3QX5EiUlTtsmRAct6fIUT3jkrsHSS/tZ66EYJ9Q0OBoX2lL/Msmi27OQvA7uYn
uqYlwJzU43tCsiip9E9z/UrLcMYyXx3oPJyPAgMBAAGjUzBRMB0GA1UdDgQWBBTI
ahE8DDT4T1vta6cXVVaRjnel0zAfBgNVHSMEGDAWgBTIahE8DDT4T1vta6cXVVaR
jnel0zAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCQbsMtPFkq
PxwOAIds3IoupuyIKmsF32ECEH/OlS+7Sj7MUJnGTQrwgjrsVS5sl8AmnGx4hPdL
VX98nEcKMNkph3Hkvh4EvgjSfmYGUXuJBcYU5jqNQrlrGv37rEf5FnvdHV1F3MG8
A0Mj0TLtcTdtaJFoOrnQuD/k0/1d+cMiYGTSaT5XK/unARqGEMd4BlWPh5P3SflV
/Vy2hHlMpv7OcZ8yaAI3htENZLus+L5kjHWKu6dxlPHKu6ef5k64su2LTNE07Vr9
S655uiFW5AX2wDVUcQEDCOiEn6SI9DTt5oQjWPMxPf+rEyfQ2f1QwVez7cyr6Qc5
OIUk31HnM/Fj
-----END CERTIFICATE-----
`,
			expectedFound: true,
			expectedError: "",
		},
		{
			name: "invalid-config-file",
			data: `
-----BEGIN CERTIFICATE-----
nope
-----END CERTIFICATE-----
`,
			expectedFound: true,
			expectedError: "x509: malformed certificate",
		},
		{
			name:          "empty",
			data:          "",
			expectedFound: true,
			expectedError: "",
		},
		{
			name:       "file-not-found",
			fetchError: &os.PathError{Err: os.ErrNotExist},
		},
		{
			name:          "error-fetching-file",
			fetchError:    errors.New("fetch failed"),
			expectedError: "failed to load mirror/ca-bundle.crt file: fetch failed",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			fileFetcher := mock.NewMockFileFetcher(mockCtrl)
			fileFetcher.EXPECT().FetchByName(CaBundleFilename).
				Return(
					&asset.File{
						Filename: CaBundleFilename,
						Data:     []byte(tc.data)},
					tc.fetchError,
				)

			asset := &CaBundle{}
			found, err := asset.Load(fileFetcher)
			assert.Equal(t, tc.expectedFound, found, "unexpected found value returned from Load")
			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
