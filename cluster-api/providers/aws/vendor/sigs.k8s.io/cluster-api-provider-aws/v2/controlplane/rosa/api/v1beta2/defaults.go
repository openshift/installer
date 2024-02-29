package v1beta2

import "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

// SetDefaults_RosaControlPlaneSpec is used by defaulter-gen.
func SetDefaults_RosaControlPlaneSpec(s *RosaControlPlaneSpec) { //nolint:golint,stylecheck
	if s.IdentityRef == nil {
		s.IdentityRef = &v1beta2.AWSIdentityReference{
			Kind: v1beta2.ControllerIdentityKind,
			Name: v1beta2.AWSClusterControllerIdentityName,
		}
	}
}
