package manifests

// AwsCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type AwsCredsSecretData struct {
	Base64encodeAccessKeyID     string
	Base64encodeSecretAccessKey string
}

// OpenStackCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type OpenStackCredsSecretData struct {
	Base64encodeCloudCreds string
}

type cloudCredsSecretData struct {
	AWS       *AwsCredsSecretData
	OpenStack *OpenStackCredsSecretData
}

type bootkubeTemplateData struct {
	Base64encodeCloudProviderConfig string
	EtcdClientCert                  string
	EtcdClientKey                   string
	KubeCaCert                      string
	KubeCaKey                       string
	McsTLSCert                      string
	McsTLSKey                       string
	PullSecret                      string
	RootCaCert                      string
	ServiceServingCaCert            string
	ServiceServingCaKey             string
	TectonicNetworkOperatorImage    string
	WorkerIgnConfig                 string
	CVOClusterID                    string
	EtcdEndpointHostnames           []string
	EtcdEndpointDNSSuffix           string
}

type tectonicTemplateData struct {
	KubeAddonOperatorImage string
	PullSecret             string
	CloudCreds             cloudCredsSecretData
}
