package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	serviceusage "google.golang.org/api/serviceusage/v1beta1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestQuotaLimitRetryBackoffUsesAllSteps(t *testing.T) {
	backoff := quotaLimitRetryBackoff
	configuredSteps := backoff.Steps
	backoff.Duration = time.Nanosecond
	backoff.Jitter = 0

	attempts := 0
	err := wait.ExponentialBackoffWithContext(context.Background(), backoff, func(context.Context) (bool, error) {
		attempts++
		return false, nil
	})

	assert.True(t, wait.Interrupted(err))
	assert.Equal(t, configuredSteps, attempts)
}

func TestLoadLimitsRetriesTransientErrors(t *testing.T) {
	tests := []struct {
		name          string
		responseCodes []int
		expectedCalls int32
		expectedCode  int
	}{
		{
			name:          "service unavailable then success",
			responseCodes: []int{http.StatusServiceUnavailable, http.StatusOK},
			expectedCalls: 2,
		},
		{
			name:          "rate limited then success",
			responseCodes: []int{http.StatusTooManyRequests, http.StatusOK},
			expectedCalls: 2,
		},
		{
			name:          "bad request is not retried",
			responseCodes: []int{http.StatusBadRequest},
			expectedCalls: 1,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name: "transient error is returned after retry limit",
			responseCodes: []int{
				http.StatusServiceUnavailable,
				http.StatusServiceUnavailable,
				http.StatusServiceUnavailable,
				http.StatusServiceUnavailable,
			},
			expectedCalls: 4,
			expectedCode:  http.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				call := calls.Add(1)
				index := min(int(call)-1, len(test.responseCodes)-1)
				code := test.responseCodes[index]
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(code)
				if code == http.StatusOK {
					_, _ = w.Write([]byte(`{"metrics":[]}`))
					return
				}
				_, _ = fmt.Fprintf(w, `{"error":{"code":%d,"message":"request failed"}}`, code)
			}))
			defer server.Close()

			client, err := serviceusage.NewService(
				context.Background(),
				option.WithEndpoint(server.URL+"/"),
				option.WithoutAuthentication(),
			)
			if !assert.NoError(t, err) {
				return
			}

			backoff := wait.Backoff{Steps: 4}
			limits, err := loadLimitsWithBackoff(context.Background(), client.Services, backoff, "test-project", "compute.googleapis.com")
			assert.Equal(t, test.expectedCalls, calls.Load())
			if test.expectedCode == 0 {
				assert.NoError(t, err)
				assert.Empty(t, limits)
				return
			}

			var gErr *googleapi.Error
			if assert.ErrorAs(t, err, &gErr) {
				assert.Equal(t, test.expectedCode, gErr.Code)
			}
		})
	}
}

func TestLoadLimitsHonorsCanceledContext(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":{"code":503,"message":"request failed"}}`))
	}))
	defer server.Close()

	client, err := serviceusage.NewService(
		context.Background(),
		option.WithEndpoint(server.URL+"/"),
		option.WithoutAuthentication(),
	)
	if !assert.NoError(t, err) {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = loadLimitsWithBackoff(ctx, client.Services, wait.Backoff{Steps: 4}, "test-project", "compute.googleapis.com")
	assert.ErrorIs(t, err, context.Canceled)
	assert.Zero(t, calls.Load())
}

func TestLoadLimitsDoesNotKeepPartialResultsWhenRetrying(t *testing.T) {
	var calls atomic.Int32
	var pageTokens []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		call := calls.Add(1)
		pageTokens = append(pageTokens, request.URL.Query().Get("pageToken"))
		w.Header().Set("Content-Type", "application/json")
		if call == 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":{"code":503,"message":"request failed"}}`))
			return
		}

		region := "us-central1"
		nextPageToken := "next-page"
		if call == 4 {
			region = "us-east1"
			nextPageToken = ""
		}
		response := &serviceusage.ListConsumerQuotaMetricsResponse{
			Metrics: []*serviceusage.ConsumerQuotaMetric{{
				ConsumerQuotaLimits: []*serviceusage.ConsumerQuotaLimit{{
					Metric: "compute.googleapis.com/cpus",
					QuotaBuckets: []*serviceusage.QuotaBucket{{
						Dimensions:     map[string]string{"region": region},
						EffectiveLimit: 24,
					}},
				}},
			}},
			NextPageToken: nextPageToken,
		}
		assert.NoError(t, json.NewEncoder(w).Encode(response))
	}))
	defer server.Close()

	client, err := serviceusage.NewService(
		context.Background(),
		option.WithEndpoint(server.URL+"/"),
		option.WithoutAuthentication(),
	)
	if !assert.NoError(t, err) {
		return
	}

	limits, err := loadLimitsWithBackoff(
		context.Background(),
		client.Services,
		wait.Backoff{Steps: 4},
		"test-project",
		"compute.googleapis.com",
	)
	assert.NoError(t, err)
	assert.Equal(t, int32(4), calls.Load())
	assert.Equal(t, []string{"", "next-page", "", "next-page"}, pageTokens)
	assert.ElementsMatch(t, []record{
		{
			Service:  "compute.googleapis.com",
			Name:     "compute.googleapis.com/cpus",
			Location: "us-central1",
			Value:    24,
		},
		{
			Service:  "compute.googleapis.com",
			Name:     "compute.googleapis.com/cpus",
			Location: "us-east1",
			Value:    24,
		},
	}, limits)
}
