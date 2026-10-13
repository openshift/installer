package gcp

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	// uidConfigMapName is the name of the config map that contains the cloud controller UID
	uidConfigMapName = "ingress-uid"
	// uidNamespace is the namespace which contains the above config map
	uidNamespace = metav1.NamespaceSystem
	// uidCluster is the data key for the clusters uid
	uidCluster = "uid"
	// uidProvider is the data key for the providers uid
	uidProvider = "provider-uid"
)

// CloudControllerUID returns a configmap with a unique UID
// per cluster used by the GCP cloud controller provider to name
// load balancer resources.
// This is based on GCP provider code that manages this configmap:
// https://github.com/openshift/kubernetes/blob/a45281f7de40f996f67d0ee7b886add59e7b5e8d/pkg/cloudprovider/providers/gce/gce_clusterid.go#L38-L57
func CloudControllerUID(infraID string) *corev1.ConfigMap {
	uid := gcptypes.CloudControllerUID(infraID)
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      uidConfigMapName,
			Namespace: uidNamespace,
		},
		Data: map[string]string{
			uidCluster:  uid,
			uidProvider: uid,
		},
	}
}
