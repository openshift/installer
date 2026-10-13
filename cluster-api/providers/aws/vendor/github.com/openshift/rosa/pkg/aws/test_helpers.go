package aws

import cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

func MockOidcConfig(id string, issuerUrl string) (*cmv1.OidcConfig, error) {
	mock := cmv1.NewOidcConfig().
		ID(id).IssuerUrl(issuerUrl)

	return mock.Build()
}
