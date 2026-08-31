package azure

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/wait"
)

type blockBlobUploadFunc func(context.Context, io.ReadSeekCloser, *blockblob.UploadOptions) (blockblob.UploadResponse, error)

func (f blockBlobUploadFunc) Upload(ctx context.Context, body io.ReadSeekCloser, options *blockblob.UploadOptions) (blockblob.UploadResponse, error) {
	return f(ctx, body, options)
}

func TestUploadBlockBlobRetriesAuthorizationPropagationErrors(t *testing.T) {
	for _, errorCode := range []string{"AuthorizationPermissionMismatch", "AuthorizationFailure"} {
		t.Run(errorCode, func(t *testing.T) {
			attempts := 0
			uploader := blockBlobUploadFunc(func(_ context.Context, body io.ReadSeekCloser, _ *blockblob.UploadOptions) (blockblob.UploadResponse, error) {
				attempts++
				got, err := io.ReadAll(body)
				assert.NoError(t, err)
				assert.Equal(t, []byte("ignition"), got)
				if attempts == 1 {
					return blockblob.UploadResponse{}, &azcore.ResponseError{StatusCode: http.StatusForbidden, ErrorCode: errorCode}
				}
				return blockblob.UploadResponse{}, nil
			})

			err := uploadBlockBlobWithRetry(context.Background(), uploader, []byte("ignition"), nil, true, wait.Backoff{Steps: 3})

			assert.NoError(t, err)
			assert.Equal(t, 2, attempts)
		})
	}
}

func TestUploadBlockBlobRetryExhaustion(t *testing.T) {
	attempts := 0
	uploadErr := &azcore.ResponseError{StatusCode: http.StatusForbidden, ErrorCode: "AuthorizationPermissionMismatch"}
	uploader := blockBlobUploadFunc(func(context.Context, io.ReadSeekCloser, *blockblob.UploadOptions) (blockblob.UploadResponse, error) {
		attempts++
		return blockblob.UploadResponse{}, uploadErr
	})

	err := uploadBlockBlobWithRetry(context.Background(), uploader, nil, nil, true, wait.Backoff{Steps: 3})

	assert.ErrorIs(t, err, uploadErr)
	assert.Equal(t, 3, attempts)
	var responseErr *azcore.ResponseError
	if assert.ErrorAs(t, err, &responseErr) {
		assert.Equal(t, "AuthorizationPermissionMismatch", responseErr.ErrorCode)
	}
}

func TestUploadBlockBlobDoesNotRetryNonQualifyingErrors(t *testing.T) {
	tests := []struct {
		name            string
		uploadErr       error
		tokenCredential bool
	}{
		{
			name:            "unrelated forbidden response",
			uploadErr:       &azcore.ResponseError{StatusCode: http.StatusForbidden, ErrorCode: "AuthenticationFailed"},
			tokenCredential: true,
		},
		{
			name:            "qualifying code with unrelated status",
			uploadErr:       &azcore.ResponseError{StatusCode: http.StatusUnauthorized, ErrorCode: "AuthorizationFailure"},
			tokenCredential: true,
		},
		{
			name:            "non-response error",
			uploadErr:       errors.New("upload failed"),
			tokenCredential: true,
		},
		{
			name:            "shared-key authorization response",
			uploadErr:       &azcore.ResponseError{StatusCode: http.StatusForbidden, ErrorCode: "AuthorizationPermissionMismatch"},
			tokenCredential: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			uploader := blockBlobUploadFunc(func(context.Context, io.ReadSeekCloser, *blockblob.UploadOptions) (blockblob.UploadResponse, error) {
				attempts++
				return blockblob.UploadResponse{}, tt.uploadErr
			})

			err := uploadBlockBlobWithRetry(context.Background(), uploader, nil, nil, tt.tokenCredential, wait.Backoff{Steps: 3})

			assert.ErrorIs(t, err, tt.uploadErr)
			assert.Equal(t, 1, attempts)
		})
	}
}

func TestUploadBlockBlobRetryRespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0
	uploadErr := &azcore.ResponseError{StatusCode: http.StatusForbidden, ErrorCode: "AuthorizationFailure"}
	uploader := blockBlobUploadFunc(func(context.Context, io.ReadSeekCloser, *blockblob.UploadOptions) (blockblob.UploadResponse, error) {
		attempts++
		cancel()
		return blockblob.UploadResponse{}, uploadErr
	})

	err := uploadBlockBlobWithRetry(ctx, uploader, nil, nil, true, wait.Backoff{Duration: retryTime, Steps: 3})

	assert.ErrorIs(t, err, context.Canceled)
	assert.ErrorIs(t, err, uploadErr)
	assert.Equal(t, 1, attempts)
}
