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
	"golang.org/x/oauth2/google"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigFileEnvVar is the name of the environment variable
// that contains the path to the credentials file.
const ConfigFileEnvVar = "GOOGLE_APPLICATION_CREDENTIALS"

// CredentialsClient resolves GCP credentials in a standard form, regardless
// of where they actually come from: an explicit service account key
// referenced by a Secret, a credentials file referenced by
// ConfigFileEnvVar, or implicit ADC (Workload Identity Federation via the
// GKE metadata server, gcloud user credentials, etc). Every caller that
// needs credentials goes through this, rather than each implementing its
// own resolution order.
type CredentialsClient struct{}

// ResolvedCredentials is everything CredentialsClient.Resolve can determine
// about a GCP credential from a single resolution pass, so callers needing
// more than just a working oauth2 credential don't have to re-derive it
// themselves or ask CredentialsClient a second time.
type ResolvedCredentials struct {
	// RawData is the raw credential JSON. Nil when Credentials was resolved
	// via fully implicit ADC (GKE metadata server, gcloud user credentials,
	// etc.) with no backing file - in that case Type is also unset.
	RawData []byte
	// Type is the detected GCP credential type. Unset when RawData is nil.
	Type google.CredentialsType
	// Credentials is the resolved oauth2 Google credentials.
	Credentials *google.Credentials
}

// Resolve resolves credentialsRef (an explicit service account key) or,
// when nil, an ADC credentials file or other implicit ADC source, scoped
// to scopes.
func (c CredentialsClient) Resolve(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, scopes ...string) (ResolvedCredentials, error) {
	rawData, err := c.Data(ctx, credentialsRef, crClient)
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("getting credential data: %w", err)
	}

	if rawData == nil {
		// Neither an explicit Secret ref nor an ADC credentials file - fall
		// back to whatever implicit ADC source applies (GKE metadata server
		// under WIF, gcloud user credentials, etc).
		creds, err := google.FindDefaultCredentials(ctx, scopes...)
		if err != nil {
			return ResolvedCredentials{}, fmt.Errorf("finding default credentials: %w", err)
		}
		return ResolvedCredentials{Credentials: creds}, nil
	}

	credType, err := c.Type(rawData)
	if err != nil {
		return ResolvedCredentials{}, err
	}

	// CredentialsFromJSONWithType (rather than the deprecated
	// CredentialsFromJSON) requires the caller to state the expected
	// credential type up front instead of trusting whatever type the JSON
	// claims to be.
	creds, err := google.CredentialsFromJSONWithType(ctx, rawData, credType, scopes...)
	if err != nil {
		return ResolvedCredentials{}, fmt.Errorf("parsing credentials: %w", err)
	}
	// CredentialsFromJSONWithType never returns (nil, nil), but guard defensively.
	if creds == nil {
		return ResolvedCredentials{}, errors.New("credentials are nil after parsing JSON")
	}

	return ResolvedCredentials{RawData: rawData, Type: credType, Credentials: creds}, nil
}

// Type detects the GCP credential type from raw credential JSON, defaulting
// to google.ServiceAccount when the type is unspecified.
func (CredentialsClient) Type(rawData []byte) (google.CredentialsType, error) {
	header := &credentialHeader{}
	if err := json.Unmarshal(rawData, header); err != nil {
		return "", fmt.Errorf("parsing credential type: %w", err)
	}

	switch header.Type {
	case "external_account":
		return google.ExternalAccount, nil
	case "impersonated_service_account":
		return google.ImpersonatedServiceAccount, nil
	default:
		return google.ServiceAccount, nil
	}
}

// Data returns the raw credential JSON referenced by credentialsRef, or,
// when credentialsRef is nil, the contents of the file referenced by
// ConfigFileEnvVar if set. Returns (nil, nil) when neither applies,
// signalling callers to fall back to implicit ADC (Workload Identity
// Federation, instance metadata, gcloud user credentials, etc.) themselves.
func (c CredentialsClient) Data(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client) ([]byte, error) {
	if credentialsRef != nil {
		return c.fromSecretRef(ctx, credentialsRef, crClient)
	}
	return c.fromADCFile()
}

// fromSecretRef reads the raw credential JSON referenced by credentialsRef.
func (CredentialsClient) fromSecretRef(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client) ([]byte, error) {
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

// fromADCFile reads credential JSON from the file referenced by
// ConfigFileEnvVar, if set.
func (CredentialsClient) fromADCFile() ([]byte, error) {
	credsPath := os.Getenv(ConfigFileEnvVar)
	if credsPath == "" {
		// No explicit credentials file configured; signal to callers to use
		// implicit ADC (Workload Identity Federation, instance metadata, etc.).
		return nil, nil
	}

	byteValue, err := os.ReadFile(credsPath) //nolint:gosec // We need to read a file here
	if err != nil {
		return nil, fmt.Errorf("reading credentials from file %s: %w", credsPath, err)
	}
	return byteValue, nil
}
