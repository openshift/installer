package alibabacloud

// CloudProviderConfig is the alibabacloud cloud provider config
type CloudProviderConfig struct {
	AccessKeyId       string
	AccessKeySecret   string
	ResourceGroupName string
}

// JSON generates the cloud provider json config for the alibabacloud platform.
// managed resource names are matching the convention defined by capz
func (params CloudProviderConfig) JSON() (string, error) {

	// TODO AlibabaCloud
	return "", nil
}
