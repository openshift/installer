package capiutils

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/openshift/installer/pkg/asset"
)

const (
	// ManifestDir defines the directory name for Cluster API manifests.
	ManifestDir = "cluster-api"
)

// GenerateClusterAssetsOutput is the output of GenerateClusterAssets.
type GenerateClusterAssetsOutput struct {
	Manifests          []*asset.RuntimeFile
	InfrastructureRefs []*corev1.ObjectReference
}

// GenerateMachinesOutput is the output of GenerateMachines.
type GenerateMachinesOutput struct {
	Manifests []*asset.RuntimeFile
}
