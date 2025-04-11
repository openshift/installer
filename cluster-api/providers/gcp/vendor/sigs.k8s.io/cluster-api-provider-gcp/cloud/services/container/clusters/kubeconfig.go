/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clusters

import (
	"context"
	"encoding/base64"
	"fmt"

	"cloud.google.com/go/container/apiv1/containerpb"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	infrav1exp "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/secret"
)

const (
	// GkeScope is the scope to request when generating access token.
	GkeScope = "https://www.googleapis.com/auth/cloud-platform"
)

func (s *Service) reconcileKubeconfig(ctx context.Context, cluster *containerpb.Cluster, log *logr.Logger) error {
	log.Info("Reconciling kubeconfig")
	clusterRef := types.NamespacedName{
		Name:      s.scope.Cluster.Name,
		Namespace: s.scope.Cluster.Namespace,
	}

	configSecret, err := secret.GetFromNamespacedName(ctx, s.scope.Client(), clusterRef, secret.Kubeconfig)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "getting kubeconfig secret", "name", clusterRef)
			return fmt.Errorf("getting kubeconfig secret %s: %w", clusterRef, err)
		}
		log.Info("kubeconfig secret not found, creating")

		if createErr := s.createCAPIKubeconfigSecret(
			ctx,
			cluster,
			&clusterRef,
			log,
		); createErr != nil {
			return fmt.Errorf("creating kubeconfig secret: %w", createErr)
		}
	} else if updateErr := s.updateCAPIKubeconfigSecret(ctx, configSecret); updateErr != nil {
		return fmt.Errorf("updating kubeconfig secret: %w", err)
	}

	return nil
}

func (s *Service) reconcileAdditionalKubeconfigs(ctx context.Context, cluster *containerpb.Cluster, log *logr.Logger) error {
	log.Info("Reconciling additional kubeconfig")
	clusterRef := types.NamespacedName{
		Name:      s.scope.Cluster.Name + "-user",
		Namespace: s.scope.Cluster.Namespace,
	}

	// Create the additional kubeconfig for users. This doesn't need updating on every sync
	_, err := secret.GetFromNamespacedName(ctx, s.scope.Client(), clusterRef, secret.Kubeconfig)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("getting kubeconfig (user) secret %s: %w", clusterRef, err)
		}

		createErr := s.createUserKubeconfigSecret(
			ctx,
			cluster,
			&clusterRef,
		)
		if createErr != nil {
			return fmt.Errorf("creating additional kubeconfig secret: %w", err)
		}
	}

	return nil
}

func (s *Service) createUserKubeconfigSecret(ctx context.Context, cluster *containerpb.Cluster, clusterRef *types.NamespacedName) error {
	controllerOwnerRef := *metav1.NewControllerRef(s.scope.GCPManagedControlPlane, infrav1exp.GroupVersion.WithKind("GCPManagedControlPlane"))

	contextName := s.getKubeConfigContextName(false)

	cfg, err := s.createBaseKubeConfig(contextName, cluster)
	if err != nil {
		return fmt.Errorf("creating base kubeconfig: %w", err)
	}

	execConfig := &api.ExecConfig{
		APIVersion:         "client.authentication.k8s.io/v1beta1",
		Command:            "gke-gcloud-auth-plugin",
		InstallHint:        "Install gke-gcloud-auth-plugin for use with kubectl by following\n		https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke",
		ProvideClusterInfo: true,
	}
	cfg.AuthInfos = map[string]*api.AuthInfo{
		contextName: {
			Exec: execConfig,
		},
	}

	out, err := clientcmd.Write(*cfg)
	if err != nil {
		return fmt.Errorf("serialize kubeconfig to yaml: %w", err)
	}

	kubeconfigSecret := kubeconfig.GenerateSecretWithOwner(*clusterRef, out, controllerOwnerRef)
	if err := s.scope.Client().Create(ctx, kubeconfigSecret); err != nil {
		return fmt.Errorf("creating secret: %w", err)
	}

	return nil
}

func (s *Service) createCAPIKubeconfigSecret(ctx context.Context, cluster *containerpb.Cluster, clusterRef *types.NamespacedName, log *logr.Logger) error {
	controllerOwnerRef := *metav1.NewControllerRef(s.scope.GCPManagedControlPlane, infrav1exp.GroupVersion.WithKind("GCPManagedControlPlane"))

	contextName := s.getKubeConfigContextName(false)

	cfg, err := s.createBaseKubeConfig(contextName, cluster)
	if err != nil {
		log.Error(err, "failed creating base config")
		return fmt.Errorf("creating base kubeconfig: %w", err)
	}

	token, err := s.generateToken(ctx)
	if err != nil {
		log.Error(err, "failed generating token")
		return err
	}
	cfg.AuthInfos = map[string]*api.AuthInfo{
		contextName: {
			Token: token,
		},
	}

	out, err := clientcmd.Write(*cfg)
	if err != nil {
		log.Error(err, "failed serializing kubeconfig to yaml")
		return fmt.Errorf("serialize kubeconfig to yaml: %w", err)
	}

	kubeconfigSecret := kubeconfig.GenerateSecretWithOwner(*clusterRef, out, controllerOwnerRef)
	if err := s.scope.Client().Create(ctx, kubeconfigSecret); err != nil {
		log.Error(err, "failed creating secret")
		return fmt.Errorf("creating secret: %w", err)
	}

	return nil
}

func (s *Service) updateCAPIKubeconfigSecret(ctx context.Context, configSecret *corev1.Secret) error {
	data, ok := configSecret.Data[secret.KubeconfigDataName]
	if !ok {
		return errors.Errorf("missing key %q in secret data", secret.KubeconfigDataName)
	}

	config, err := clientcmd.Load(data)
	if err != nil {
		return errors.Wrap(err, "failed to convert kubeconfig Secret into a clientcmdapi.Config")
	}

	token, err := s.generateToken(ctx)
	if err != nil {
		return err
	}

	contextName := s.getKubeConfigContextName(false)
	config.AuthInfos[contextName].Token = token

	out, err := clientcmd.Write(*config)
	if err != nil {
		return errors.Wrap(err, "failed to serialize config to yaml")
	}

	configSecret.Data[secret.KubeconfigDataName] = out

	err = s.scope.Client().Update(ctx, configSecret)
	if err != nil {
		return fmt.Errorf("updating kubeconfig secret: %w", err)
	}

	return nil
}

func (s *Service) getKubeConfigContextName(isUser bool) string {
	contextName := fmt.Sprintf("gke_%s_%s_%s", s.scope.GCPManagedControlPlane.Spec.Project, s.scope.GCPManagedControlPlane.Spec.Location, s.scope.ClusterName())
	if isUser {
		contextName += "-user"
	}
	return contextName
}

func (s *Service) createBaseKubeConfig(contextName string, cluster *containerpb.Cluster) (*api.Config, error) {
	certData, err := base64.StdEncoding.DecodeString(cluster.GetMasterAuth().GetClusterCaCertificate())
	if err != nil {
		return nil, fmt.Errorf("decoding cluster CA cert: %w", err)
	}
	cfg := &api.Config{
		APIVersion: api.SchemeGroupVersion.Version,
		Clusters: map[string]*api.Cluster{
			contextName: {
				Server:                   "https://" + cluster.GetEndpoint(),
				CertificateAuthorityData: certData,
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  contextName,
				AuthInfo: contextName,
			},
		},
		CurrentContext: contextName,
	}

	return cfg, nil
}

func (s *Service) generateToken(ctx context.Context) (string, error) {
	req := &credentialspb.GenerateAccessTokenRequest{
		Name: "projects/-/serviceAccounts/" + s.scope.GetCredential().ClientEmail,
		Scope: []string{
			GkeScope,
		},
	}
	resp, err := s.scope.CredentialsClient().GenerateAccessToken(ctx, req)
	if err != nil {
		return "", errors.Errorf("error generating access token: %v", err)
	}

	return resp.GetAccessToken(), nil
}
