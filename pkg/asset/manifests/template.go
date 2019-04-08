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
}

// OpenStackCredsSecretData holds encoded credentials and is used to generate cloud-creds secret
type OpenStackCredsSecretData struct {
	Base64encodeCloudCreds string
}

// VSphereCredsSecretData holds encoded credentials and is used to generated cloud-creds secret
type VSphereCredsSecretData struct {
	VirtualCenters []VSphereVirtualCenterCredsSecretData
}

// VSphereVirtualCenterCredsSecretData holds the encoded credentials for a vSphere vCenter.
type VSphereVirtualCenterCredsSecretData struct {
	Name                 string
	Base64encodeUsername string
	Base64encodePassword string
}

type cloudCredsSecretData struct {
	AWS       *AwsCredsSecretData
	Azure     *AzureCredsSecretData
	OpenStack *OpenStackCredsSecretData
	VSphere   *VSphereCredsSecretData
}

type bootkubeTemplateData struct {
	CVOClusterID          string
	EtcdCaBundle          string
	EtcdCaCert            string
	EtcdClientCaCert      string
	EtcdClientCaKey       string
	EtcdClientCert        string
	EtcdClientKey         string
	EtcdEndpointDNSSuffix string
	EtcdEndpointHostnames []string
	EtcdMetricCaCert      string
	EtcdMetricClientCert  string
	EtcdMetricClientKey   string
	EtcdSignerCert        string
	EtcdSignerClientCert  string
	EtcdSignerClientKey   string
	EtcdSignerKey         string
	McsTLSCert            string
	McsTLSKey             string
	PullSecretBase64      string
	RootCaCert            string
	WorkerIgnConfig       string
}

type openshiftTemplateData struct {
	CloudCreds                   cloudCredsSecretData
	Base64EncodedKubeadminPwHash string
}
