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

package iamauth

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
)

const (
	configMapName = "aws-auth"
	configMapNS   = metav1.NamespaceSystem

	roleKey  = "mapRoles"
	usersKey = "mapUsers"
)

type configMapBackend struct {
	client crclient.Client
}

func (b *configMapBackend) MapRole(mapping ekscontrolplanev1.RoleMapping) error {
	if errs := mapping.Validate(); errs != nil {
		return kerrors.NewAggregate(errs)
	}

	authConfig, err := b.getAuthConfig()
	if err != nil {
		return fmt.Errorf("getting auth config: %w", err)
	}

	for _, existingMapping := range authConfig.RoleMappings {
		if cmp.Equal(existingMapping, mapping) {
			// A mapping already exists that matches, so ignore
			return nil
		}
	}

	authConfig.RoleMappings = append(authConfig.RoleMappings, mapping)

	return b.saveAuthConfig(authConfig)
}

func (b *configMapBackend) MapUser(mapping ekscontrolplanev1.UserMapping) error {
	if errs := mapping.Validate(); errs != nil {
		return kerrors.NewAggregate(errs)
	}

	authConfig, err := b.getAuthConfig()
	if err != nil {
		return fmt.Errorf("getting auth config: %w", err)
	}

	for _, existingMapping := range authConfig.UserMappings {
		if cmp.Equal(existingMapping, mapping) {
			// A mapping already exists that matches, so ignore
			return nil
		}
	}

	authConfig.UserMappings = append(authConfig.UserMappings, mapping)

	return b.saveAuthConfig(authConfig)
}

func (b *configMapBackend) getAuthConfig() (*ekscontrolplanev1.IAMAuthenticatorConfig, error) {
	ctx := context.Background()

	configMapRef := types.NamespacedName{
		Name:      configMapName,
		Namespace: configMapNS,
	}

	authConfigMap := &corev1.ConfigMap{}

	err := b.client.Get(ctx, configMapRef, authConfigMap)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, fmt.Errorf("getting %s/%s config map: %w", configMapName, configMapNS, err)
	}

	authConfig := &ekscontrolplanev1.IAMAuthenticatorConfig{
		RoleMappings: []ekscontrolplanev1.RoleMapping{},
		UserMappings: []ekscontrolplanev1.UserMapping{},
	}
	if authConfigMap.Data == nil {
		return authConfig, nil
	}

	mappedRoles, err := b.getMappedRoles(authConfigMap)
	if err != nil {
		return nil, fmt.Errorf("getting mapped roles: %w", err)
	}
	authConfig.RoleMappings = mappedRoles

	mappedUsers, err := b.getMappedUsers(authConfigMap)
	if err != nil {
		return nil, fmt.Errorf("getting mapped users: %w", err)
	}
	authConfig.UserMappings = mappedUsers

	return authConfig, nil
}

func (b *configMapBackend) saveAuthConfig(authConfig *ekscontrolplanev1.IAMAuthenticatorConfig) error {
	ctx := context.Background()

	configMapRef := types.NamespacedName{
		Name:      configMapName,
		Namespace: configMapNS,
	}

	authConfigMap := &corev1.ConfigMap{}

	err := b.client.Get(ctx, configMapRef, authConfigMap)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("getting %s/%s config map: %w", configMapName, configMapNS, err)
	}

	if authConfigMap.Data == nil {
		authConfigMap.Data = make(map[string]string)
	}
	authConfigMap = authConfigMap.DeepCopy()

	delete(authConfigMap.Data, roleKey)
	delete(authConfigMap.Data, usersKey)

	if len(authConfig.RoleMappings) > 0 {
		roleMappings, err := yaml.Marshal(authConfig.RoleMappings)
		if err != nil {
			return fmt.Errorf("marshalling auth config roles: %w", err)
		}
		authConfigMap.Data[roleKey] = string(roleMappings)
	}

	if len(authConfig.UserMappings) > 0 {
		userMappings, err := yaml.Marshal(authConfig.UserMappings)
		if err != nil {
			return fmt.Errorf("marshalling auth config users: %w", err)
		}
		authConfigMap.Data[usersKey] = string(userMappings)
	}

	if authConfigMap.UID == "" {
		authConfigMap.Name = configMapName
		authConfigMap.Namespace = configMapNS
		return b.client.Create(ctx, authConfigMap)
	}

	return b.client.Update(ctx, authConfigMap)
}

func (b *configMapBackend) getMappedRoles(cm *corev1.ConfigMap) ([]ekscontrolplanev1.RoleMapping, error) {
	mappedRoles := []ekscontrolplanev1.RoleMapping{}

	rolesSection, ok := cm.Data[roleKey]
	if !ok {
		return mappedRoles, nil
	}

	if err := yaml.Unmarshal([]byte(rolesSection), &mappedRoles); err != nil {
		return nil, fmt.Errorf("unmarshalling mapped roles: %w", err)
	}

	return mappedRoles, nil
}

func (b *configMapBackend) getMappedUsers(cm *corev1.ConfigMap) ([]ekscontrolplanev1.UserMapping, error) {
	mappedUsers := []ekscontrolplanev1.UserMapping{}

	usersSection, ok := cm.Data[usersKey]
	if !ok {
		return mappedUsers, nil
	}

	if err := yaml.Unmarshal([]byte(usersSection), &mappedUsers); err != nil {
		return nil, fmt.Errorf("unmarshalling mapped users: %w", err)
	}

	return mappedUsers, nil
}
