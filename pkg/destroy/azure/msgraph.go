package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// MsGraphClient is a lightweight client for Microsoft Graph API calls.
type MsGraphClient struct {
	endpoint   string // e.g. "https://graph.microsoft.com/v1.0"
	scope      string // e.g. "https://graph.microsoft.com/.default"
	cred       azcore.TokenCredential
	httpClient *http.Client
}

// MSGraphClientOption configures the MsGraphClient.
type MSGraphClientOption func(*MsGraphClient)

// WithHTTPClient sets a custom http.Client, replacing the default entirely.
// Mutually exclusive with WithTimeout by intent; if both are applied, the
// last one wins.
func WithHTTPClient(c *http.Client) MSGraphClientOption {
	return func(m *MsGraphClient) {
		m.httpClient = c
	}
}

// WithTimeout sets the timeout on the default http.Client.
// For more control over the client, use WithHTTPClient instead.
func WithTimeout(d time.Duration) MSGraphClientOption {
	return func(m *MsGraphClient) {
		m.httpClient.Timeout = d
	}
}

// NewMSGraphClient constructs a new client with default scope. You should wrap this client in a
// struct and use the generic methods below to implement the graph client you need on demand. If the
// endpoint is empty or "N/A", a nil client will be returned.
func NewMSGraphClient(endpoint string, cred azcore.TokenCredential, opts ...MSGraphClientOption) *MsGraphClient {
	// Set the Microsoft Graph endpoint for the appropriate cloud
	// (e.g., GovCloud). This can be empty for StackCloud or "N/A"
	// for clouds where Microsoft Graph is not available (e.g., German Cloud).
	// See https://issues.redhat.com/browse/OCPBUGS-4549
	// See https://learn.microsoft.com/en-us/graph/sdks/national-clouds?tabs=go
	if endpoint == "" || endpoint == "N/A" {
		return nil
	}
	trimmed := strings.TrimRight(endpoint, "/")
	fullEndpoint := trimmed + "/v1.0"
	defaultScope := trimmed + "/.default"

	c := &MsGraphClient{
		endpoint:   fullEndpoint,
		scope:      defaultScope,
		cred:       cred,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// IsAvailable returns true when the MSGraphClient was configured with a valid endpoint.
func (c *MsGraphClient) IsAvailable() bool {
	return c != nil
}

// ListTheoreticalKindWithFilter lists the given kind and unmarshals it into the given type.
// You are responsible for making sure both the kind and type are correct for the graph SDK.
// No paging support.
func ListTheoreticalKindWithFilter[Kind any](ctx context.Context, c *MsGraphClient, kind, filter string) ([]Kind, error) {
	path := "/" + kind + "?$filter=" + url.QueryEscape(filter)
	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var result odataCollection[Kind]
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", kind, err)
	}
	return result.Value, nil
}

type odataCollection[Kind any] struct {
	Value []Kind `json:"value"`
}

// DeleteTheoreticalKind deletes the given kind with the given id.
// You are responsible for the kind value is correct for the graph SDK.
func DeleteTheoreticalKind(ctx context.Context, c *MsGraphClient, kind, id string) error {
	path := "/" + kind + "/" + url.PathEscape(id)
	_, err := c.doRequest(ctx, http.MethodDelete, path)
	return err
}

// Only List and Delete are implemented because that is all the installer needs today.
// Add new methods following the same pattern if needed.

func (c *MsGraphClient) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	if !c.IsAvailable() {
		return nil, fmt.Errorf("MSGraph client is not available")
	}

	token, err := c.cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{c.scope}})
	if err != nil {
		return nil, fmt.Errorf("failed to get Graph API token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req) //nolint:gosec // endpoint is set at construction from Azure environment config, not user input
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseODataError(body, resp.StatusCode, resp.Header)
	}

	return body, nil
}

func parseODataError(body []byte, statusCode int, headers http.Header) error {
	retryAfter := parseRetryAfter(headers.Get("Retry-After"))

	var odataErr odataErrorResponse
	if err := json.Unmarshal(body, &odataErr); err == nil && odataErr.Error.Code != "" {
		return &MSGraphError{
			StatusCode: statusCode,
			Code:       odataErr.Error.Code,
			Message:    odataErr.Error.Message,
			RetryAfter: retryAfter,
		}
	}
	return &MSGraphError{
		StatusCode: statusCode,
		Code:       "",
		Message:    string(body),
		RetryAfter: retryAfter,
	}
}

func parseRetryAfter(val string) time.Duration {
	if val == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(val); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return 0
}

// MSGraphError represents an error response from the Microsoft Graph API,
// preserving the HTTP status code for classification by callers.
type MSGraphError struct {
	StatusCode int
	Code       string
	Message    string
	RetryAfter time.Duration // from Retry-After header, zero if absent
}

func (e *MSGraphError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("unexpected status code %d: %s", e.StatusCode, e.Message)
}

type odataErrorResponse struct {
	Error odataErrorDetail `json:"error"`
}

type odataErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Methods and types needed specifically by the Azure destroyer:
// Move the MsGraphClient and all private symbols above this doc comment into
// their own package, and keep everything below this doc comment here if you
// ever need to extend the client or use it elsewhere.

type azureDestroyerClient struct {
	*MsGraphClient
}

type msGraphServicePrincipal struct {
	ID          string `json:"id"`
	AppID       string `json:"appId"`
	DisplayName string `json:"displayName"`
}

type msGraphApplication struct {
	ID string `json:"id"`
}

func (c *azureDestroyerClient) listServicePrincipals(ctx context.Context, filter string) ([]msGraphServicePrincipal, error) {
	return ListTheoreticalKindWithFilter[msGraphServicePrincipal](ctx, c.MsGraphClient, "servicePrincipals", filter)
}

func (c *azureDestroyerClient) listApplications(ctx context.Context, filter string) ([]msGraphApplication, error) {
	return ListTheoreticalKindWithFilter[msGraphApplication](ctx, c.MsGraphClient, "applications", filter)
}

func (c *azureDestroyerClient) deleteApplication(ctx context.Context, id string) error {
	return DeleteTheoreticalKind(ctx, c.MsGraphClient, "applications", id)
}
