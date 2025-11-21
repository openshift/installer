/*
Copyright 2021 The Kubernetes Authors.

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

// Package system contains utiilities for the system namespace.
package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	// namespaceEnvVarName is the env var coming from DownwardAPI in the manager manifest.
	namespaceEnvVarName = "POD_NAMESPACE"
	// defaultNamespace is the default value from manifest.
	defaultNamespace = "capa-system"
	// inClusterNamespacePath is the file the default namespace to be used for namespaced API operations is placed at.
	inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

// GetManagerNamespace return the namespace where the controller is running.
func GetManagerNamespace() string {
	ns, err := GetNamespaceFromFile(inClusterNamespacePath)
	if err == nil {
		return ns
	}

	// If file is not there, check if env var is set, otherwise return the default namespace.
	managerNamespace := os.Getenv(namespaceEnvVarName)
	if managerNamespace == "" {
		managerNamespace = defaultNamespace
	}
	return managerNamespace
}

// GetNamespaceFromFile returns the namespace from a file.
// This code is copied from controller-runtime, because it is a private method there.
// https://github.com/kubernetes-sigs/controller-runtime/blob/316aea4229158103123166a5e45076f1a86bd807/pkg/leaderelection/leader_election.go#L104
func GetNamespaceFromFile(nsFilePath string) (string, error) {
	// Check whether the namespace file exists.
	// If not, we are not running in cluster so can't guess the namespace.

	if _, err := os.Stat(nsFilePath); os.IsNotExist(err) {
		return "", errors.Wrapf(err, "not running in-cluster, please specify LeaderElectionNamespace")
	} else if err != nil {
		return "", errors.Wrapf(err, "error checking namespace file: %s", nsFilePath)
	}

	// Load the namespace file and return its content
	namespace, err := os.ReadFile(filepath.Clean(nsFilePath))
	if err != nil {
		return "", fmt.Errorf("error reading namespace file: %w", err)
	}
	return string(namespace), nil
}
