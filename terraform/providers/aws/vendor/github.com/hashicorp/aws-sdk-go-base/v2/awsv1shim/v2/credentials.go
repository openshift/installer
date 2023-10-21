// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package awsv1shim

import ( // nosemgrep: no-sdkv2-imports-in-awsv1shim
	"context"
	"fmt"
	"sync/atomic"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type v2CredentialsProvider struct {
	provider awsv2.CredentialsProvider

	v2creds atomic.Value
}

// This adapter deals with multiple levels of caching and a slight mismatch between the AWS SDK for Go v1 and v2 credentials models.
// In the SDK v1 model has a root `credentials.Credentials` struct that handles caching. The `credentials.Value` contains only keys.
// The `credentials.Credentials` struct handles expiry information by calling the credentials provider.
// In the SDK v2 model, the SDK returns an `aws.CredentialsCache` which handles caching. The `aws.Credentials` value contains keys
// as well as the expiry information.
//
// The `v2CredentialsProvider` will typically be used with the following layout:
// (v1)`credentials.Credentials` ==> `v2CredentialsProvider` ==> (v2)`aws.CredentialsCache` ==> (v2)<actual credentials provider>
//
// Since the SDK v1 `credentials.Credentials` handles expiry, it has an `Expire` function to explicitly expire credentials. This is
// used, for example, in the SDK v1 default retry handler to catch an expired credentials error. Because of this, the result of
// `RetrieveWithContext` cannot be cached in `v2CredentialsProvider`.
// NOTE: Since the `Expire()` call is not passed up the chain, the (v2)`aws.CredentialsCache` will not have its cache cleared. This
// may cause problems if a credential is revoked early. If this becomes a problem, every call to `RetrieveWithContext` may need to
// call `Invalidate()` on the (v2)`aws.CredentialsCache`. In practice, `RetrieveWithContext` is rarely called, so this is not likely
// to have a significant impact.
//
// The expiry information is cached in `v2CredentialsProvider` because the SDK v1 model handles expiry separately from the credential
// information, and otherwise calling `IsExpired()` and `ExpiresAt()` would potentially call the actual credential provider on each call.

func (p *v2CredentialsProvider) RetrieveWithContext(ctx credentials.Context) (credentials.Value, error) {
	v2creds, err := p.provider.Retrieve(ctx)
	if err != nil {
		return credentials.Value{}, err
	}
	p.v2creds.Store(&v2creds)

	return credentials.Value{
		AccessKeyID:     v2creds.AccessKeyID,
		SecretAccessKey: v2creds.SecretAccessKey,
		SessionToken:    v2creds.SessionToken,
		ProviderName:    fmt.Sprintf("v2Credentials(%s)", v2creds.Source),
	}, nil
}

func (p *v2CredentialsProvider) IsExpired() bool {
	v2creds := p.credentials()
	if v2creds != nil {
		return v2creds.Expired()
	}
	return true
}

func (p *v2CredentialsProvider) ExpiresAt() time.Time {
	v2creds := p.credentials()
	if v2creds != nil {
		return v2creds.Expires
	}
	return time.Time{}
}

func (p *v2CredentialsProvider) Retrieve() (credentials.Value, error) {
	return p.RetrieveWithContext(context.Background())
}

func (p *v2CredentialsProvider) credentials() *awsv2.Credentials {
	v := p.v2creds.Load()
	if v == nil {
		return nil
	}

	c := v.(*awsv2.Credentials)
	if c != nil && c.HasKeys() && !c.Expired() {
		return c
	}

	return nil
}

func newV2Credentials(v2provider awsv2.CredentialsProvider) *credentials.Credentials {
	return credentials.NewCredentials(&v2CredentialsProvider{
		provider: v2provider,
	})
}
