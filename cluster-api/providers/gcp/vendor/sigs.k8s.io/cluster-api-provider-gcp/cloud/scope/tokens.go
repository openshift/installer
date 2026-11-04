/*
Copyright 2025 The Kubernetes Authors.

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
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/compute/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// defaultTokenRequestTimeout bounds every HTTP call a production
// TokenClient's underlying client makes to resolve credentials and refresh
// tokens. oauth2.TokenSource.Token() takes no context of its own, so
// without a bounded client here a stalled request to Google's token
// endpoint (or the GKE metadata server under WIF) can block the calling
// reconcile indefinitely.
const defaultTokenRequestTimeout = 30 * time.Second

// TokenClient mints short-lived GCP access tokens for a configured
// identity, hiding GCP's credential-resolution and oauth2 token-source
// machinery behind a single Token method - callers never see an
// oauth2.TokenSource or the credentials used to build one.
type TokenClient struct {
	source oauth2.TokenSource
}

// NewTokenClient resolves credentialsRef (an explicit service account key)
// or, when nil, implicit ADC (including Workload Identity Federation via
// the GKE metadata server), and returns a TokenClient that mints tokens
// from the result, bounding every HTTP call it makes to
// defaultTokenRequestTimeout. Use NewTokenClientWithTimeout for a different
// timeout (e.g. a shorter one in tests).
func NewTokenClient(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client) (TokenClient, error) {
	return NewTokenClientWithTimeout(ctx, credentialsRef, crClient, defaultTokenRequestTimeout)
}

// NewTokenClientWithTimeout is NewTokenClient with an explicit timeout
// bounding every HTTP call made while resolving credentials. This is
// reused for every future token refresh made through the returned
// TokenClient, not just this construction call - oauth2.TokenSource.Token()
// takes no context of its own, so without a bounded client here a stalled
// request to Google's token endpoint (or the GKE metadata server under WIF)
// can block the calling reconcile indefinitely.
//
// This mints tokens directly from the identity's own credentials rather than
// via the IAM Credentials API's self-impersonation endpoint, which requires
// the identity to hold roles/iam.serviceAccountTokenCreator on itself.
func NewTokenClientWithTimeout(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, timeout time.Duration) (TokenClient, error) {
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Timeout: timeout})

	resolved, err := CredentialsClient{}.Resolve(ctx, credentialsRef, crClient, compute.CloudPlatformScope)
	if err != nil {
		return TokenClient{}, fmt.Errorf("resolving credentials: %w", err)
	}
	return TokenClient{source: resolved.Credentials.TokenSource}, nil
}

// Token mints (or returns a cached, not-yet-expired) short-lived access
// token for the identity NewTokenClient resolved.
func (t TokenClient) Token() (string, error) {
	tok, err := t.source.Token()
	if err != nil {
		return "", fmt.Errorf("minting access token: %w", err)
	}
	return tok.AccessToken, nil
}
