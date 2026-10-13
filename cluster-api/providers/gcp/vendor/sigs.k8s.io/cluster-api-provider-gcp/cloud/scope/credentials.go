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

package scope

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ConfigFileEnvVar is the name of the environment variable
	// that contains the path to the credentials file.
	ConfigFileEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"
)

// Credential is a struct to hold GCP credential data.
type Credential struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	ClientEmail string `json:"client_email"`
	ClientID    string `json:"client_id"`
}

func getCredentials(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client) (*Credential, error) {
	var credentialData []byte
	var err error

	if credentialsRef != nil {
		credentialData, err = getCredentialDataFromRef(ctx, credentialsRef, crClient)
	} else {
		credentialData, err = getCredentialDataUsingADC()
	}
	if err != nil {
		return nil, fmt.Errorf("getting credential data: %w", err)
	}

	return parseCredential(credentialData)
}

func getCredentialDataFromRef(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client) ([]byte, error) {
	secretRefName := types.NamespacedName{
		Name:      credentialsRef.Name,
		Namespace: credentialsRef.Namespace,
	}

	credSecret := &corev1.Secret{}
	if err := crClient.Get(ctx, secretRefName, credSecret); err != nil {
		return nil, fmt.Errorf("getting credentials secret %s\\%s: %w", secretRefName.Namespace, secretRefName.Name, err)
	}

	rawData, ok := credSecret.Data["credentials"]
	if !ok {
		return nil, errors.New("no credentials key in secret")
	}

	return rawData, nil
}

func getCredentialDataUsingADC() ([]byte, error) {
	credsPath := os.Getenv(ConfigFileEnvVar)
	if credsPath == "" {
		return nil, fmt.Errorf("no ADC environment variable found for credentials (expect %s)", ConfigFileEnvVar)
	}

	byteValue, err := os.ReadFile(credsPath) //nolint:gosec // We need to read a file here
	if err != nil {
		return nil, fmt.Errorf("reading credentials from file %s: %w", credsPath, err)
	}
	return byteValue, nil
}

func parseCredential(rawData []byte) (*Credential, error) {
	var credential Credential
	err := json.Unmarshal(rawData, &credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}
