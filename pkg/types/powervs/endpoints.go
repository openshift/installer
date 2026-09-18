package powervs

import configv1 "github.com/openshift/api/config/v1"

// EndpointURLForService returns the URL override for a named IBM Cloud service
// from the ServiceEndpoints slice, or an empty string if no override is set.
// When the SDK receives an empty string URL it uses its own built-in default.
func EndpointURLForService(name string, endpoints []configv1.PowerVSServiceEndpoint) string {
	for _, ep := range endpoints {
		if ep.Name == name {
			return ep.URL
		}
	}
	return ""
}
