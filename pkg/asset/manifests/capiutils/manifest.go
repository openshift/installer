package capiutils

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Manifests is a list of Manifests sortable by filename.
type Manifests []*Manifest

func (m Manifests) Len() int {
	return len(m)
}

func (m Manifests) Less(i, j int) bool {
	return m[i].Filename < m[j].Filename
}

func (m Manifests) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Manifest is a wrapper for a CAPI manifest.
type Manifest struct {
	Object   client.Object
	Filename string
}

// GenerateClusterAssetsOutput is the output of GenerateClusterAssets.
type GenerateClusterAssetsOutput struct {
	Manifests         Manifests
	InfrastructureRef *corev1.ObjectReference
}
