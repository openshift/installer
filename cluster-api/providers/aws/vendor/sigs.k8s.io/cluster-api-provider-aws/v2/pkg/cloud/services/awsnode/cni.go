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
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kustomize/api/konfig"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

const (
	awsNodeName      = "aws-node"
	awsNodeNamespace = "kube-system"
)

// ReconcileCNI will reconcile the CNI of a service.
func (s *Service) ReconcileCNI(ctx context.Context) error {
	s.scope.Info("Reconciling aws-node DaemonSet in cluster", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	remoteClient, err := s.scope.RemoteClient()
	if err != nil {
		s.scope.Error(err, "getting client for remote cluster")
		return fmt.Errorf("getting client for remote cluster: %w", err)
	}

	if s.scope.DisableVPCCNI() {
		if err := s.deleteCNI(ctx, remoteClient); err != nil {
			return fmt.Errorf("disabling aws vpc cni: %w", err)
		}
		return nil
	}

	var ds appsv1.DaemonSet
	if err := remoteClient.Get(ctx, types.NamespacedName{Namespace: awsNodeNamespace, Name: awsNodeName}, &ds); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		return ErrCNIMissing
	}

	var needsUpdate bool
	if len(s.scope.VpcCni().Env) > 0 {
		s.scope.Info("updating aws-node daemonset environment variables", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

		for i := range ds.Spec.Template.Spec.Containers {
			container := &ds.Spec.Template.Spec.Containers[i]
			if container.Name == "aws-node" {
				container.Env, needsUpdate = s.applyUserProvidedEnvironmentProperties(container.Env)
			}
		}
	}

	secondarySubnets := s.secondarySubnets()
	if len(secondarySubnets) == 0 {
		if needsUpdate {
			s.scope.Info("adding environment properties to vpc-cni", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))
			if err = remoteClient.Update(ctx, &ds, &client.UpdateOptions{}); err != nil {
				return err
			}
		}

		// with no secondary subnets there is no need for eni configs
		return nil
	}

	sgs, err := s.getSecurityGroups()
	if err != nil {
		return err
	}

	metaLabels := map[string]string{
		"app.kubernetes.io/managed-by": "cluster-api-provider-aws",
		"app.kubernetes.io/part-of":    s.scope.Name(),
	}

	s.scope.Info("for each subnet", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))
	for _, subnet := range secondarySubnets {
		var eniConfig amazoncni.ENIConfig
		if err := remoteClient.Get(ctx, types.NamespacedName{Namespace: metav1.NamespaceSystem, Name: subnet.AvailabilityZone}, &eniConfig); err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}
			s.scope.Info("Creating ENIConfig", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()), "subnet", subnet.ID, "availability-zone", subnet.AvailabilityZone)
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

		s.scope.Info("Updating ENIConfig", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()), "subnet", subnet.ID, "availability-zone", subnet.AvailabilityZone)
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
			s.scope.Info("Removing old ENIConfig", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()), "eniConfig", oldEniConfig.Name)
			if err := remoteClient.Delete(ctx, &oldEniConfig, &client.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	s.scope.Info("updating containers", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
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

// applyUserProvidedEnvironmentProperties takes a container environment and applies user provided values to it.
func (s *Service) applyUserProvidedEnvironmentProperties(containerEnv []corev1.EnvVar) ([]corev1.EnvVar, bool) {
	var (
		envVars     = make(map[string]corev1.EnvVar)
		needsUpdate = false
	)
	for _, e := range s.scope.VpcCni().Env {
		envVars[e.Name] = e
	}
	// Handle the case where we overwrite an existing value if it's not already the desired value.
	// This will prevent continuously updating the DaemonSet even though there are no changes.
	for i, e := range containerEnv {
		if v, ok := envVars[e.Name]; ok {
			// Take care of comparing secret ref with Stringer.
			if containerEnv[i].String() != v.String() {
				needsUpdate = true
				containerEnv[i] = v
			}
			delete(envVars, e.Name)
		}
	}
	// Handle case when there are values that aren't in the list of environment properties
	// of aws-node.
	for _, v := range envVars {
		needsUpdate = true
		containerEnv = append(containerEnv, v)
	}
	return containerEnv, needsUpdate
}

func (s *Service) deleteCNI(ctx context.Context, remoteClient client.Client) error {
	// EKS has a tendency to pre-install the vpc-cni automagically even if you don't specify it as an addon
	// and looks like a kubectl apply from a script of a manifest that looks like this
	// https://github.com/aws/amazon-vpc-cni-k8s/blob/master/config/master/aws-k8s-cni.yaml
	// and removing these pieces will enable someone to install and alternative CNI. There is also another use
	// case where someone would want to remove the vpc-cni and reinstall it via the helm chart located here
	// https://github.com/aws/amazon-vpc-cni-k8s/tree/master/charts/aws-vpc-cni meaning we need to account for
	// managed-by: Helm label, or we will delete the helm chart resources every reconcile loop. EKS does make
	// a CRD for eniconfigs but the default env var on the vpc-cni pod is ENABLE_POD_ENI=false. We will make an
	// assumption no CRs are ever created and leave the CRD to reduce complexity of this operation.

	s.scope.Info("Ensuring all resources for AWS VPC CNI in cluster are deleted", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())

	s.scope.Info("Trying to delete AWS VPC CNI DaemonSet", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
	if err := s.deleteResource(ctx, remoteClient, types.NamespacedName{
		Namespace: awsNodeNamespace,
		Name:      awsNodeName,
	}, &appsv1.DaemonSet{}); err != nil {
		return err
	}

	s.scope.Info("Trying to delete AWS VPC CNI ServiceAccount", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
	if err := s.deleteResource(ctx, remoteClient, types.NamespacedName{
		Namespace: awsNodeNamespace,
		Name:      awsNodeName,
	}, &corev1.ServiceAccount{}); err != nil {
		return err
	}

	s.scope.Info("Trying to delete AWS VPC CNI ClusterRoleBinding", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
	if err := s.deleteResource(ctx, remoteClient, types.NamespacedName{
		Namespace: string(meta.RESTScopeNameRoot),
		Name:      awsNodeName,
	}, &rbacv1.ClusterRoleBinding{}); err != nil {
		return err
	}

	s.scope.Info("Trying to delete AWS VPC CNI ClusterRole", "cluster-name", s.scope.Name(), "cluster-namespace", s.scope.Namespace())
	if err := s.deleteResource(ctx, remoteClient, types.NamespacedName{
		Namespace: string(meta.RESTScopeNameRoot),
		Name:      awsNodeName,
	}, &rbacv1.ClusterRole{}); err != nil {
		return err
	}

	record.Eventf(s.scope.InfraCluster(), "DeletedVPCCNI", "The AWS VPC CNI has been removed from the cluster. Ensure you enable a CNI via another mechanism")

	return nil
}

func (s *Service) deleteResource(ctx context.Context, remoteClient client.Client, key client.ObjectKey, obj client.Object) error {
	if err := remoteClient.Get(ctx, key, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("deleting resource %s: %w", key, err)
		}
		s.scope.Debug(fmt.Sprintf("resource %s was not found, no action", key))
	} else {
		// resource found, delete if no label or not managed by helm
		if val, ok := obj.GetLabels()[konfig.ManagedbyLabelKey]; !ok || val != "Helm" {
			if err := remoteClient.Delete(ctx, obj, &client.DeleteOptions{}); err != nil {
				if !apierrors.IsNotFound(err) {
					return fmt.Errorf("deleting %s: %w", key, err)
				}
				s.scope.Debug(fmt.Sprintf(
					"resource %s was not found, not deleted", key))
			} else {
				s.scope.Debug(fmt.Sprintf("resource %s was deleted", key))
			}
		} else {
			s.scope.Debug(fmt.Sprintf("resource %s is managed by helm, not deleted", key))
		}
	}

	return nil
}
