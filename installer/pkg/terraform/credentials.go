package terraform

import "errors"

// Credentials holds all the information that are necessary to
// make TerraForm authenticate against various providers.
type Credentials struct {
	*AWSCredentials `json:",inline"`
}

// ToEnvironment returns the environment variables that are expected by
// TerraForm's provider.
func (c *Credentials) ToEnvironment() (map[string]string, error) {
	env := make(map[string]string)

	if c.AWSCredentials != nil {
		if err := c.AWSCredentials.Validate(); err != nil {
			return nil, err
		}
		mergeEnvMap(env, c.AWSCredentials.ToEnvironment())
	}

	return env, nil
}

// AWSCredentials represents the credentials required by TerraForm's AWS
// provider.
type AWSCredentials struct {
	AWSAccessKeyID     string `json:"AWSAccessKeyID"`
	AWSSecretAccessKey string `json:"AWSSecretAccessKey"`
	AWSSessionToken    string `json:"AWSSessionToken"`
}

// Validate verifies that the given credentials are valid.
func (a *AWSCredentials) Validate() error {
	if a.AWSAccessKeyID == "" || a.AWSSecretAccessKey == "" {
		return errors.New("AWSAccessKeyID & AWSSecretAccessKey must be specified")
	}
	return nil
}

// ToEnvironment returns the environment variables that are expected by
// TerraForm's provider.
func (a *AWSCredentials) ToEnvironment() map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     a.AWSAccessKeyID,
		"AWS_SECRET_ACCESS_KEY": a.AWSSecretAccessKey,
		"AWS_SESSION_TOKEN":     a.AWSSessionToken,
	}
}

// mergeEnvMap is an utility function that write all keys from the src string
// map to the destination one, overriding any existing one.
func mergeEnvMap(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}
