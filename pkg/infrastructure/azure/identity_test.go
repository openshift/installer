package azure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoleAssignmentWithRetry(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "retries principal not found",
			run: func(t *testing.T) {
				calls := 0
				err := createRoleAssignmentWithRetry(context.Background(), time.Millisecond, time.Second, func(context.Context) error {
					calls++
					if calls == 1 {
						return &azcore.ResponseError{StatusCode: http.StatusBadRequest, ErrorCode: "PrincipalNotFound"}
					}
					return nil
				})

				assert.NoError(t, err)
				assert.Equal(t, 2, calls)
			},
		},
		{
			name: "does not retry 429 outside the azure sdk",
			run: func(t *testing.T) {
				calls := 0
				tooManyRequestsErr := &azcore.ResponseError{StatusCode: http.StatusTooManyRequests}
				err := createRoleAssignmentWithRetry(context.Background(), time.Millisecond, time.Second, func(context.Context) error {
					calls++
					return tooManyRequestsErr
				})

				assert.Same(t, tooManyRequestsErr, err)
				assert.Equal(t, 1, calls)
			},
		},
		{
			name: "returns forbidden without retrying",
			run: func(t *testing.T) {
				calls := 0
				forbiddenErr := &azcore.ResponseError{StatusCode: http.StatusForbidden, ErrorCode: "AuthorizationFailed"}
				err := createRoleAssignmentWithRetry(context.Background(), time.Millisecond, time.Second, func(context.Context) error {
					calls++
					return forbiddenErr
				})

				assert.Same(t, forbiddenErr, err)
				assert.Equal(t, 1, calls)
			},
		},
		{
			name: "cancelled context stops retrying",
			run: func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				calls := 0
				err := createRoleAssignmentWithRetry(ctx, time.Hour, time.Minute, func(context.Context) error {
					calls++
					cancel()
					return &azcore.ResponseError{StatusCode: http.StatusBadRequest, ErrorCode: "PrincipalNotFound"}
				})

				assert.ErrorIs(t, err, context.Canceled)
				assert.Equal(t, 1, calls)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestRoleAssignmentAzureSDKRetriesHTTP429(t *testing.T) {
	requests := 0
	client := newTestRoleAssignmentsClient(t, roleAssignmentTestTransport(func(req *http.Request) (*http.Response, error) {
		requests++
		if requests == 1 {
			return testRoleAssignmentResponse(req, http.StatusTooManyRequests), nil
		}
		return testRoleAssignmentResponse(req, http.StatusCreated), nil
	}))

	err := createRoleAssignmentWithRetry(context.Background(), time.Millisecond, 5*time.Second, func(ctx context.Context) error {
		_, err := client.Create(ctx, "/subscriptions/sub/resourceGroups/rg", "role-assignment-id", armauthorization.RoleAssignmentCreateParameters{
			Properties: &armauthorization.RoleAssignmentProperties{},
		}, nil)
		return err
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, requests)
}

func TestRoleAssignmentAzureSDKRetriesNetworkTimeout(t *testing.T) {
	requests := 0
	client := newTestRoleAssignmentsClient(t, roleAssignmentTestTransport(func(req *http.Request) (*http.Response, error) {
		requests++
		if requests == 1 {
			return nil, testTimeoutError{}
		}
		return testRoleAssignmentResponse(req, http.StatusCreated), nil
	}))

	err := createRoleAssignmentWithRetry(context.Background(), time.Millisecond, 5*time.Second, func(ctx context.Context) error {
		_, err := client.Create(ctx, "/subscriptions/sub/resourceGroups/rg", "role-assignment-id", armauthorization.RoleAssignmentCreateParameters{
			Properties: &armauthorization.RoleAssignmentProperties{},
		}, nil)
		return err
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, requests)
}

func TestCreateRoleAssignmentWithRetryPreservesLastErrorAtDeadline(t *testing.T) {
	calls := 0
	principalNotFoundErr := &azcore.ResponseError{StatusCode: http.StatusBadRequest, ErrorCode: "PrincipalNotFound"}
	err := createRoleAssignmentWithRetry(context.Background(), time.Millisecond, 250*time.Millisecond, func(ctx context.Context) error {
		calls++
		if calls == 1 {
			return principalNotFoundErr
		}
		<-ctx.Done()
		return ctx.Err()
	})

	var responseErr *azcore.ResponseError
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.ErrorAs(t, err, &responseErr)
	assert.Same(t, principalNotFoundErr, responseErr)
	assert.Equal(t, 2, calls)
}

func TestCreateRoleAssignmentWithRetryPreservesLastErrorOnParentDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	calls := 0
	principalNotFoundErr := &azcore.ResponseError{StatusCode: http.StatusBadRequest, ErrorCode: "PrincipalNotFound"}
	err := createRoleAssignmentWithRetry(ctx, time.Millisecond, time.Second, func(ctx context.Context) error {
		calls++
		if calls == 1 {
			return principalNotFoundErr
		}
		<-ctx.Done()
		return ctx.Err()
	})

	var responseErr *azcore.ResponseError
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.ErrorAs(t, err, &responseErr)
	assert.Same(t, principalNotFoundErr, responseErr)
	assert.Equal(t, 2, calls)
}

func TestShouldRetryRoleAssignmentError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "principal not found",
			err:  &azcore.ResponseError{ErrorCode: "PrincipalNotFound"},
			want: true,
		},
		{
			name: "too many requests code is handled by the sdk",
			err:  &azcore.ResponseError{ErrorCode: "TooManyRequests"},
			want: false,
		},
		{
			name: "request throttled code is handled by the sdk",
			err:  &azcore.ResponseError{ErrorCode: "RequestThrottled"},
			want: false,
		},
		{
			name: "server error is handled by the sdk",
			err:  &azcore.ResponseError{StatusCode: http.StatusInternalServerError},
			want: false,
		},
		{
			name: "not implemented is not retried",
			err:  &azcore.ResponseError{StatusCode: http.StatusNotImplemented},
			want: false,
		},
		{
			name: "wrapped principal not found",
			err:  fmt.Errorf("role assignment failed: %w", &azcore.ResponseError{ErrorCode: "PrincipalNotFound"}),
			want: true,
		},
		{
			name: "forbidden",
			err:  &azcore.ResponseError{StatusCode: http.StatusForbidden},
			want: false,
		},
		{
			name: "network timeout is handled by the sdk",
			err:  testTimeoutError{},
			want: false,
		},
		{
			name: "context cancelled",
			err:  context.Canceled,
			want: false,
		},
		{
			name: "context deadline exceeded",
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "other error",
			err:  errors.New("unexpected error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, shouldRetryRoleAssignmentError(tt.err))
		})
	}
}

type testTimeoutError struct{}

func (testTimeoutError) Error() string   { return "request timed out" }
func (testTimeoutError) Timeout() bool   { return true }
func (testTimeoutError) Temporary() bool { return true }

type roleAssignmentTestTransport func(*http.Request) (*http.Response, error)

func (transport roleAssignmentTestTransport) Do(req *http.Request) (*http.Response, error) {
	return transport(req)
}

func newTestRoleAssignmentsClient(t *testing.T, transport roleAssignmentTestTransport) *armauthorization.RoleAssignmentsClient {
	client, err := armauthorization.NewRoleAssignmentsClient("sub", testTokenCredential{}, &arm.ClientOptions{
		ClientOptions: azpolicy.ClientOptions{
			Retry:     azpolicy.RetryOptions{RetryDelay: time.Millisecond},
			Transport: transport,
		},
		DisableRPRegistration: true,
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return client
}

func testRoleAssignmentResponse(req *http.Request, statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{}`)),
		Request:    req,
	}
}

type testTokenCredential struct{}

func (testTokenCredential) GetToken(context.Context, azpolicy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "test-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}
