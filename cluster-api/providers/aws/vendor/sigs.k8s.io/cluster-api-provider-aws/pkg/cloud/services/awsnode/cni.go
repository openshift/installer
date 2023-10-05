/*
Copyright 2020 The Kubernetes Authors.

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

package awsnode

import (
	"context"
	"fmt"

	amazoncni "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/crd/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/record"
)

const (
	awsNodeName      = "aws-node"
	awsNodeNamespace = "kube-system"
)

// ReconcileCNI will reconcile the CNI of a service.
func (s *Service) ReconcileCNI(ctx context.Context) error {
	s.scope.Info("Reconciling aws-node DaemonSet in cluster", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())

	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		s.scope.Error(err, "getting client for remote cluster")
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	if s.scope.DisableVPCCNI() {
		if err := s.deleteCNI(ctx, remoteClient); err != nil {
			return fmt.Errorf("disabling aws vpc cni: %w", err)
		}
	}

	if s.scope.SecondaryCidrBlock() == nil {
		return nil
	}

	var ds appsv1.DaemonSet
	if err := remoteClient.Get(ctx, types.NamespacedName{Namespace: awsNodeNamespace, Name: awsNodeName}, &ds); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return ErrCNIMissing
	}

	sgs, err := s.getSecurityGroups()
	if err != nil {
		return err
	}

	metaLabels := map[string]string{
		"app.kubernetes.io/managed-by": "cluster-api-provider-aws",
		"app.kubernetes.io/part-of":    s.scope.Name(),
	}

	s.scope.Info("for each subnet", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
	for _, subnet := range s.secondarySubnets() {
		var eniConfig amazoncni.ENIConfig
		if err := remoteClient.Get(ctx, types.NamespacedName{Namespace: metav1.NamespaceSystem, Name: subnet.AvailabilityZone}, &eniConfig); err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			s.scope.Info("Creating ENIConfig", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace(), "subnet", subnet.ID, "availability-zone", subnet.AvailabilityZone)
			eniConfig = amazoncni.ENIConfig{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: metav1.NamespaceSystem,
					Name:      subnet.AvailabilityZone,
					Labels:    metaLabels,
				},
				Spec: amazoncni.ENIConfigSpec{
					Subnet:         subnet.ID,
					SecurityGroups: sgs,
				},
			}

			if err := remoteClient.Create(ctx, &eniConfig, &client.CreateOptions{}); err != nil {
				return err
			}
		}

		s.scope.Info("Updating ENIConfig", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace(), "subnet", subnet.ID, "availability-zone", subnet.AvailabilityZone)
		eniConfig.Spec = amazoncni.ENIConfigSpec{
			Subnet:         subnet.ID,
			SecurityGroups: sgs,
		}

		if err := remoteClient.Update(ctx, &eniConfig, &client.UpdateOptions{}); err != nil {
			return err
		}
	}

	// Removing any ENIConfig no longer needed
	var eniConfigs amazoncni.ENIConfigList
	err = remoteClient.List(ctx, &eniConfigs, &client.ListOptions{
		Namespace:     metav1.NamespaceSystem,
		LabelSelector: labels.SelectorFromSet(metaLabels),
	})
	if err != nil {
		return err
	}
	for _, eniConfig := range eniConfigs.Items {
		matchFound := false
		for _, subnet := range s.secondarySubnets() {
			if eniConfig.Name == subnet.AvailabilityZone {
				matchFound = true
				break
			}
		}

		if !matchFound {
			oldEniConfig := eniConfig
			s.scope.Info("Removing old ENIConfig", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace(), "eniConfig", oldEniConfig.Name)
			if err := remoteClient.Delete(ctx, &oldEniConfig, &client.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	s.scope.Info("updating containers", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
	for _, container := range ds.Spec.Template.Spec.Containers {
		if container.Name == "aws-node" {
			container.Env = append(s.filterEnv(container.Env),
				corev1.EnvVar{
					Name:  "AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG",
					Value: "true",
				},
				corev1.EnvVar{
					Name:  "ENI_CONFIG_LABEL_DEF",
					Value: "failure-domain.beta.kubernetes.io/zone",
				},
			)
		}
	}

	return remoteClient.Update(ctx, &ds, &client.UpdateOptions{})
}

func (s *Service) getSecurityGroups() ([]string, error) {
	sgRoles := []infrav1.SecurityGroupRole{
		infrav1.SecurityGroupNode,
	}

	sgs := make([]string, 0, len(sgRoles))
	for _, sg := range sgRoles {
		if _, ok := s.scope.SecurityGroups()[sg]; !ok {
			return nil, awserrors.NewFailedDependency(fmt.Sprintf("%s security group not available", sg))
		}
		sgs = append(sgs, s.scope.SecurityGroups()[sg].ID)
	}

	return sgs, nil
}

func (s *Service) filterEnv(env []corev1.EnvVar) []corev1.EnvVar {
	var i int
	for _, e := range env {
		if e.Name == "ENI_CONFIG_LABEL_DEF" || e.Name == "AWS_VPC_K8S_CNI_CUSTOM_NETWORK_CFG" {
			continue
		}
		env[i] = e
		i++
	}
	return env[:i]
}

func (s *Service) deleteCNI(ctx context.Context, remoteClient client.Client) error {
	s.scope.Info("Ensuring aws-node DaemonSet in cluster is deleted", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())

	ds := &appsv1.DaemonSet{}
	if err := remoteClient.Get(ctx, types.NamespacedName{Namespace: awsNodeNamespace, Name: awsNodeName}, ds); err != nil {
		if apierrors.IsNotFound(err) {
			s.scope.V(2).Info("The aws-node DaemonSet is not found, not action")
			return nil
		}
		return fmt.Errorf("getting aws-node daemonset: %w", err)
	}

	s.scope.V(2).Info("The aws-node DaemonSet found, deleting")
	if err := remoteClient.Delete(ctx, ds, &client.DeleteOptions{}); err != nil {
		if apierrors.IsNotFound(err) {
			s.scope.V(2).Info("The aws-node DaemonSet is not found, not deleted")
			return nil
		}
		return fmt.Errorf("deleting aws-node DaemonSet: %w", err)
	}
	record.Eventf(s.scope.InfraCluster(), "DeletedVPCCNI", "The AWS VPC CNI has been removed from the cluster. Ensure you enable a CNI via another mechanism")

	return nil
}
