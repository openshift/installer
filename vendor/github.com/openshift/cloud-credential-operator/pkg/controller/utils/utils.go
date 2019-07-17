package utils

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/sirupsen/logrus"

	configv1 "github.com/openshift/api/config/v1"
)

const (
	awsCredsSecretIDKey     = "aws_access_key_id"
	awsCredsSecretAccessKey = "aws_secret_access_key"
)

func LoadCredsFromSecret(kubeClient client.Client, namespace, secretName string) ([]byte, []byte, error) {

	secret := &corev1.Secret{}
	err := kubeClient.Get(context.TODO(),
		types.NamespacedName{
			Name:      secretName,
			Namespace: namespace,
		},
		secret)
	if err != nil {
		return nil, nil, err
	}
	accessKeyID, ok := secret.Data[awsCredsSecretIDKey]
	if !ok {
		return nil, nil, fmt.Errorf("AWS credentials secret %v did not contain key %v",
			secretName, awsCredsSecretIDKey)
	}
	secretAccessKey, ok := secret.Data[awsCredsSecretAccessKey]
	if !ok {
		return nil, nil, fmt.Errorf("AWS credentials secret %v did not contain key %v",
			secretName, awsCredsSecretAccessKey)
	}
	return accessKeyID, secretAccessKey, nil
}

// LoadInfrastructureName loads the cluster Infrastructure config and returns the infra name
// used to identify this cluster, and tag some cloud objects.
func LoadInfrastructureName(c client.Client, logger log.FieldLogger) (string, error) {
	infra := &configv1.Infrastructure{}
	err := c.Get(context.Background(), types.NamespacedName{Name: "cluster"}, infra)
	if err != nil {
		logger.WithError(err).Error("error loading Infrastructure config 'cluster'")
		return "", err
	}

	logger.Debugf("Loaded infrastructure name: %s", infra.Status.InfrastructureName)
	return infra.Status.InfrastructureName, nil

}
