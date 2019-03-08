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
	CVOClusterID                    string
	EtcdCaBundle                    string
	EtcdCaCert                      string
	EtcdClientCaCert                string
	EtcdClientCaKey                 string
	EtcdClientCert                  string
	EtcdClientKey                   string
	EtcdEndpointDNSSuffix           string
	EtcdEndpointHostnames           []string
	EtcdMetricsCaCert               string
	EtcdMetricsClientCert           string
	EtcdMetricsClientKey            string
	EtcdSignerCert                  string
	EtcdSignerClientCert            string
	EtcdSignerClientKey             string
	EtcdSignerKey                   string
	EtcdMetricsServerCert           string
	EtcdMetricsServerKey            string
	McsTLSCert                      string
	McsTLSKey                       string
	PullSecretBase64                string
	RootCaCert                      string
	WorkerIgnConfig                 string
}

type openshiftTemplateData struct {
	CloudCreds                   cloudCredsSecretData
	Base64EncodedKubeadminPwHash string
}
