package clients

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// AuthType is the method of authenticating requests to the server
type AuthType string

const (
	// NoAuth uses no authentication
	NoAuth AuthType = "noauth"
	// HTTPBasicAuth uses HTTP Basic Authentication
	HTTPBasicAuth AuthType = "http_basic"
)

// AuthConfig contains data needed to configure authentication in the client
type AuthConfig struct {
	Type     AuthType
	Username string
	Password string
}

func authRoot() string {
	env := os.Getenv("METAL3_AUTH_ROOT_DIR")
	if env != "" {
		return filepath.Clean(env)
	}
	return "/opt/metal3/auth"
}

func readAuthFile(filename string) (string, error) {
	content, err := os.ReadFile(filepath.Clean(filename))
	return strings.TrimSpace(string(content)), err
}

func load(clientType string) (auth AuthConfig, err error) {
	authPath := path.Join(authRoot(), clientType)

	if _, err := os.Stat(authPath); err != nil {
		if os.IsNotExist(err) {
			auth.Type = NoAuth
			return auth, nil
		}
		return auth, err
	}
	auth.Type = HTTPBasicAuth

	auth.Username, err = readAuthFile(path.Join(authPath, "username"))
	if err != nil {
		return
	}

	auth.Password, err = readAuthFile(path.Join(authPath, "password"))
	if err != nil {
		return
	}

	if auth.Username == "" {
		err = fmt.Errorf("Empty HTTP Basic Auth username")
	} else if auth.Password == "" {
		err = fmt.Errorf("Empty HTTP Basic Auth password")
	}
	return
}

// LoadAuth loads the Ironic and Inspector configuration from the environment
func LoadAuth() (ironicAuth, inspectorAuth AuthConfig, err error) {
	ironicAuth, err = load("ironic")
	if err != nil {
		return
	}
	inspectorAuth, err = load("ironic-inspector")
	return
}

// ConfigFromEndpointURL returns an endpoint and an auth config from an
// endpoint URL that may contain HTTP basic auth credentials.
func ConfigFromEndpointURL(endpointURL string) (endpoint string, auth AuthConfig, err error) {
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return
	}

	if parsedURL.User != nil {
		var hasPasswd bool
		auth.Type = HTTPBasicAuth
		auth.Username = parsedURL.User.Username()
		auth.Password, hasPasswd = parsedURL.User.Password()
		if !hasPasswd {
			err = fmt.Errorf("No password supplied for HTTP Basic Auth")
		}
		parsedURL.User = nil
	} else {
		auth.Type = NoAuth
	}

	endpoint = parsedURL.String()
	return
}
