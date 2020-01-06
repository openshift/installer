package manifests

// AwsCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type AwsCredsSecretData struct {
	Base64encodeAccessKeyID     string
	Base64encodeSecretAccessKey string
}

// AzureCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type AzureCredsSecretData struct {
	Base64encodeSubscriptionID string
	Base64encodeClientID       string
	Base64encodeClientSecret   string
	Base64encodeTenantID       string
	Base64encodeResourcePrefix string
	Base64encodeResourceGroup  string
	Base64encodeRegion         string
	Base64encodeEnvironment	   string
}

// GCPCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type GCPCredsSecretData struct {
	Base64encodeServiceAccount string
}

// OpenStackCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type OpenStackCredsSecretData struct {
	Base64encodeCloudCreds    string
	Base64encodeCloudCredsINI string
}

// VSphereCredsSecretData holds encoded credentials and is used to generated cloud-creds secret
type VSphereCredsSecretData struct {
	VCenter              string
	Base64encodeUsername string
	Base64encodePassword string
}

type cloudCredsSecretData struct {
	AWS       *AwsCredsSecretData
	Azure     *AzureCredsSecretData
	GCP       *GCPCredsSecretData
	OpenStack *OpenStackCredsSecretData
	VSphere   *VSphereCredsSecretData
}

type bootkubeTemplateData struct {
	CVOClusterID               string
	EtcdCaBundle               string
	EtcdEndpointDNSSuffix      string
	EtcdEndpointHostnames      []string
	EtcdMetricCaCert           string
	EtcdMetricSignerCert       string
	EtcdMetricSignerClientCert string
	EtcdMetricSignerClientKey  string
	EtcdMetricSignerKey        string
	EtcdSignerCert             string
	EtcdSignerClientCert       string
	EtcdSignerClientKey        string
	EtcdSignerKey              string
	McsTLSCert                 string
	McsTLSKey                  string
	PullSecretBase64           string
	RootCaCert                 string
	WorkerIgnConfig            string
}

type openshiftTemplateData struct {
	CloudCreds                   cloudCredsSecretData
	Base64EncodedKubeadminPwHash string
}
