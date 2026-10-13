package rosa

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	restclient "k8s.io/client-go/rest"
)

// TokenResponse contains the access token and the duration until it expires.
type TokenResponse struct {
	AccessToken string
	ExpiresIn   time.Duration
}

// RequestToken requests an OAuth access token for the specified API server using username/password credentials.
// returns a TokenResponse which contains the AccessToken and the ExpiresIn duration.
func RequestToken(ctx context.Context, apiURL, username, password string, config *restclient.Config) (*TokenResponse, error) {
	clientID := "openshift-challenging-client"
	oauthURL, err := buildOauthURL(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build oauth url: %w", err)
	}

	tokenReqURL := fmt.Sprintf("%s/oauth/authorize?response_type=token&client_id=%s", oauthURL, clientID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenReqURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", getBasicHeader(username, password))
	// this is a required header by the openshift oauth server to prevent CORS errors.
	// see https://access.redhat.com/documentation/en-us/openshift_container_platform/4.14/html/authentication_and_authorization/understanding-authentication#oauth-token-requests_understanding-authentication
	request.Header.Set("X-CSRF-Token", "1")

	transport, err := restclient.TransportFor(config)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Transport: transport}
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// don't resolve redirects and return the response instead
		return http.ErrUseLastResponse
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return nil, fmt.Errorf("expected status code %d, but got %d", http.StatusFound, resp.StatusCode)
	}

	// extract access_token & expires_in from redirect URL
	tokenResponse, err := extractTokenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to extract access token from redirect url")
	}

	return tokenResponse, nil
}

func getBasicHeader(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

func buildOauthURL(apiURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(apiURL)
	if err != nil {
		return "", err
	}
	host, _, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return "", err
	}
	parsedURL.Host = host

	oauthURL := strings.Replace(parsedURL.String(), "api", "oauth", 1)
	return oauthURL, nil
}

func extractTokenResponse(resp *http.Response) (*TokenResponse, error) {
	location, err := resp.Location()
	if err != nil {
		return nil, err
	}

	fragments, err := url.ParseQuery(location.Fragment)
	if err != nil {
		return nil, err
	}
	if len(fragments["access_token"]) == 0 {
		return nil, fmt.Errorf("access_token not found")
	}

	expiresIn, err := strconv.Atoi(fragments.Get("expires_in"))
	if err != nil || expiresIn == 0 {
		expiresIn = 86400 // default to 1 day
	}

	return &TokenResponse{
		AccessToken: fragments.Get("access_token"),
		ExpiresIn:   time.Second * time.Duration(expiresIn),
	}, nil
}
