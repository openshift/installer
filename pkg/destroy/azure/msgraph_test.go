package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
)

type fakeTokenCredential struct{}

func (f *fakeTokenCredential) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{Token: "fake-token", ExpiresOn: time.Now().Add(time.Hour)}, nil
}

func newTestClient(t *testing.T, handler http.Handler) *MsGraphClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewMSGraphClient(srv.URL, &fakeTokenCredential{})
}

type testItem struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Maybe  bool     `json:"maybe"`
	Lots   bool     `json:"lots"`
	More   int32    `json:"more"`
	Fields []string `json:"fields"`
}

func TestListTheoreticalKindWithFilter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1.0/things", r.URL.Path)
		assert.Equal(t, "Bearer fake-token", r.Header.Get("Authorization"))
		assert.Equal(t, "name eq 'foo'", r.URL.Query().Get("$filter"))

		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(odataCollection[testItem]{
			Value: []testItem{
				{ID: "1", Name: "foo"},
				{ID: "2", Name: "foo-2"},
			},
		}))
	})

	client := newTestClient(t, handler)
	items, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "name eq 'foo'")
	if assert.NoError(t, err) {
		assert.Len(t, items, 2)
		assert.Equal(t, "1", items[0].ID)
		assert.Equal(t, "foo-2", items[1].Name)
	}
}

func TestListTheoreticalKindWithFilterEmpty(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(odataCollection[testItem]{Value: []testItem{}}))
	})

	client := newTestClient(t, handler)
	items, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	if assert.NoError(t, err) {
		assert.Empty(t, items)
	}
}

func TestDeleteTheoreticalKind(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v1.0/things/obj-1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	client := newTestClient(t, handler)
	err := DeleteTheoreticalKind(context.Background(), client, "things", "obj-1")
	assert.NoError(t, err)
}

func TestODataErrorParsing(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		assert.NoError(t, json.NewEncoder(w).Encode(odataErrorResponse{
			Error: odataErrorDetail{
				Code:    "Authorization_RequestDenied",
				Message: "Insufficient privileges to complete the operation.",
			},
		}))
	})

	client := newTestClient(t, handler)
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Authorization_RequestDenied")
		assert.Contains(t, err.Error(), "Insufficient privileges")
	}
}

func TestNonODataErrorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("internal server error"))
		assert.NoError(t, err)
	})

	client := newTestClient(t, handler)
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "unexpected status code 500")
	}
}

func TestDeleteNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		assert.NoError(t, json.NewEncoder(w).Encode(odataErrorResponse{
			Error: odataErrorDetail{
				Code:    "Request_ResourceNotFound",
				Message: "Resource 'obj-gone' does not exist.",
			},
		}))
	})

	client := newTestClient(t, handler)
	err := DeleteTheoreticalKind(context.Background(), client, "things", "obj-gone")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Request_ResourceNotFound")
	}
}

func TestTrailingSlashNormalized(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.False(t, strings.Contains(r.URL.Path, "//"), "path should not contain double slashes: %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(odataCollection[testItem]{}))
	})

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	// Simulate endpoint with trailing slash like Azure environments provide
	client := NewMSGraphClient(srv.URL+"/", &fakeTokenCredential{})
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	assert.NoError(t, err)
}

func TestAuthorizationHeaderSent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer fake-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(odataCollection[testItem]{}))
	})

	client := newTestClient(t, handler)
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	assert.NoError(t, err)
}

func TestNewMSGraphClientNilForUnavailable(t *testing.T) {
	assert.Nil(t, NewMSGraphClient("", &fakeTokenCredential{}))
	assert.Nil(t, NewMSGraphClient("N/A", &fakeTokenCredential{}))
}

func TestIsAvailable(t *testing.T) {
	assert.False(t, (*MsGraphClient)(nil).IsAvailable())
	assert.False(t, NewMSGraphClient("", &fakeTokenCredential{}).IsAvailable())
	assert.False(t, NewMSGraphClient("N/A", &fakeTokenCredential{}).IsAvailable())
	assert.True(t, NewMSGraphClient("https://graph.microsoft.com", &fakeTokenCredential{}).IsAvailable())
}

func TestUnavailableClientReturnsError(t *testing.T) {
	var client *MsGraphClient // nil
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "not available")
	}

	err = DeleteTheoreticalKind(context.Background(), client, "things", "obj-1")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "not available")
	}
}

func TestRetryAfterParsed(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
		assert.NoError(t, json.NewEncoder(w).Encode(odataErrorResponse{
			Error: odataErrorDetail{
				Code:    "activityLimitReached",
				Message: "Too many requests",
			},
		}))
	})

	client := newTestClient(t, handler)
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	var graphErr *MSGraphError
	if assert.ErrorAs(t, err, &graphErr) {
		assert.Equal(t, http.StatusTooManyRequests, graphErr.StatusCode)
		assert.Equal(t, 30*time.Second, graphErr.RetryAfter)
	}
}

func TestRetryAfterZeroWhenAbsent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("error"))
		assert.NoError(t, err)
	})

	client := newTestClient(t, handler)
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	var graphErr *MSGraphError
	if assert.ErrorAs(t, err, &graphErr) {
		assert.Zero(t, graphErr.RetryAfter)
	}
}

func TestMSGraphErrorUnwrapsThroughWrapping(t *testing.T) {
	original := &MSGraphError{
		StatusCode: http.StatusForbidden,
		Code:       "Authorization_RequestDenied",
		Message:    "Insufficient privileges",
	}

	// Simulate how callers wrap the error with fmt.Errorf %w
	wrapped := fmt.Errorf("failed to gather list of Service Principals by tag: %w", original)
	doubleWrapped := fmt.Errorf("outer context: %w", wrapped)

	var graphErr *MSGraphError
	if assert.ErrorAs(t, doubleWrapped, &graphErr) {
		assert.Equal(t, http.StatusForbidden, graphErr.StatusCode)
		assert.Equal(t, "Authorization_RequestDenied", graphErr.Code)
	}
}

func TestWithTimeout(t *testing.T) {
	client := NewMSGraphClient("https://graph.microsoft.com", &fakeTokenCredential{}, WithTimeout(42*time.Second))
	assert.Equal(t, 42*time.Second, client.httpClient.Timeout)
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 99 * time.Second}
	client := NewMSGraphClient("https://graph.microsoft.com", &fakeTokenCredential{}, WithHTTPClient(custom))
	assert.Same(t, custom, client.httpClient)
}

func TestWithHTTPClientOverridesWithTimeout(t *testing.T) {
	custom := &http.Client{Timeout: 99 * time.Second}
	client := NewMSGraphClient("https://graph.microsoft.com", &fakeTokenCredential{},
		WithTimeout(5*time.Second),
		WithHTTPClient(custom),
	)
	assert.Same(t, custom, client.httpClient)
	assert.Equal(t, 99*time.Second, client.httpClient.Timeout)
}

func TestWithTimeoutAfterWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	client := NewMSGraphClient("https://graph.microsoft.com", &fakeTokenCredential{},
		WithHTTPClient(custom),
		WithTimeout(30*time.Second),
	)
	assert.Same(t, custom, client.httpClient)
	assert.Equal(t, 30*time.Second, client.httpClient.Timeout)
}

func TestWithTimeoutUsedInRequests(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(odataCollection[testItem]{}))
	})

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client := NewMSGraphClient(srv.URL, &fakeTokenCredential{}, WithTimeout(5*time.Second))
	_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
	assert.NoError(t, err)
}

func TestDefaultClientHasTimeout(t *testing.T) {
	client := NewMSGraphClient("https://graph.microsoft.com", &fakeTokenCredential{})
	assert.Equal(t, 30*time.Second, client.httpClient.Timeout)
}

func TestMSGraphErrorPreservesStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		odata      bool
		code       string
		message    string
	}{
		{
			name:       "403 with OData error",
			statusCode: http.StatusForbidden,
			odata:      true,
			code:       "Authorization_RequestDenied",
			message:    "Insufficient privileges",
		},
		{
			name:       "401 with OData error",
			statusCode: http.StatusUnauthorized,
			odata:      true,
			code:       "InvalidAuthenticationToken",
			message:    "Access token is empty",
		},
		{
			name:       "500 without OData",
			statusCode: http.StatusInternalServerError,
			odata:      false,
			message:    "something broke",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				if tc.odata {
					assert.NoError(t, json.NewEncoder(w).Encode(odataErrorResponse{
						Error: odataErrorDetail{Code: tc.code, Message: tc.message},
					}))
				} else {
					_, err := w.Write([]byte(tc.message))
					assert.NoError(t, err)
				}
			})

			client := newTestClient(t, handler)
			_, err := ListTheoreticalKindWithFilter[testItem](context.Background(), client, "things", "filter")
			if assert.Error(t, err) {
				var graphErr *MSGraphError
				if assert.ErrorAs(t, err, &graphErr) {
					assert.Equal(t, tc.statusCode, graphErr.StatusCode, "status code should be preserved in error")
				}
			}
		})
	}
}
