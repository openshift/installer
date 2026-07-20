package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	googleoauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

var (
	authEnvs            = []string{"GOOGLE_APPLICATION_CREDENTIALS", "GOOGLE_CREDENTIALS", "GOOGLE_CLOUD_KEYFILE_JSON", "GCLOUD_KEYFILE_JSON"}
	defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".gcp", "osServiceAccount.json")
	credLoaders         = []credLoader{}
	onceLoggers         = map[credLoader]*sync.Once{}
)

// Session is an object representing session for GCP API.
type Session struct {
	Credentials *googleoauth.Credentials

	// Path contains the filepath for provided credentials. When authenticating with
	// Default Application Credentials, Path will be empty.
	Path string
}

// GetSession returns a GCP session by using credentials found in default locations in order:
// env GOOGLE_CREDENTIALS,
// env GOOGLE_CLOUD_KEYFILE_JSON,
// env GCLOUD_KEYFILE_JSON,
// file ~/.gcp/osServiceAccount.json, and
// gcloud cli defaults
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession(ctx context.Context) (*Session, error) {
	creds, path, err := loadCredentials(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}

	if creds.JSON != nil {
		if err := validateCredentialURLs(creds.JSON); err != nil {
			return nil, err
		}
	}

	return &Session{
		Credentials: creds,
		Path:        path,
	}, nil
}

func loadCredentials(ctx context.Context) (*googleoauth.Credentials, string, error) {
	if len(credLoaders) == 0 {
		for _, authEnv := range authEnvs {
			credLoaders = append(credLoaders, &envLoader{env: authEnv})
		}
		credLoaders = append(credLoaders, &fileLoader{path: defaultAuthFilePath})
		credLoaders = append(credLoaders, &cliLoader{})

		for _, credLoader := range credLoaders {
			onceLoggers[credLoader] = new(sync.Once)
		}
	}

	for _, loader := range credLoaders {
		creds, err := loader.Load(ctx)
		if err != nil {
			continue
		}
		onceLoggers[loader].Do(func() {
			logrus.Infof("Credentials loaded from %s", loader)
		})
		return creds, loader.Content(), nil
	}
	return getCredentials(ctx)
}

func getCredentials(ctx context.Context) (*googleoauth.Credentials, string, error) {
	creds, err := (&userLoader{}).Load(ctx)
	if err != nil {
		return nil, "", err
	}

	filePath := defaultAuthFilePath
	logrus.Infof("Saving the credentials to %q", filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return nil, "", err
	}
	if err := os.WriteFile(filePath, creds.JSON, 0o600); err != nil {
		return nil, "", err
	}
	return creds, filePath, nil
}

type credLoader interface {
	Load(context.Context) (*googleoauth.Credentials, error)
	Content() string
}

type envLoader struct {
	env      string
	delegate credLoader
}

func (e *envLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	if val := os.Getenv(e.env); len(val) > 0 {
		e.delegate = &fileOrContentLoader{pathOrContent: val}
		return e.delegate.Load(ctx)
	}
	return nil, errors.New("empty environment variable")
}

func (e *envLoader) String() string {
	path := []string{
		fmt.Sprintf("environment variable %q", e.env),
	}
	if e.delegate != nil {
		path = append(path, fmt.Sprintf("%s", e.delegate))
	}
	return strings.Join(path, ", ")
}

func (e *envLoader) Content() string {
	envValue, found := os.LookupEnv(e.env)
	if !found {
		return ""
	}
	return envValue
}

type fileOrContentLoader struct {
	pathOrContent string
	delegate      credLoader
}

func (fc *fileOrContentLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(fc.pathOrContent); err != nil {
		return nil, fmt.Errorf("supplied value should be the path to a GCP credentials file: %w", err)
	}
	fc.delegate = &fileLoader{path: fc.pathOrContent}
	return fc.delegate.Load(ctx)
}

func (fc *fileOrContentLoader) String() string {
	if fc.delegate != nil {
		return fmt.Sprintf("%s", fc.delegate)
	}
	return "file or content"
}

func (fc *fileOrContentLoader) Content() string {
	if _, err := os.Stat(fc.pathOrContent); err != nil {
		return ""
	}
	return fc.pathOrContent
}

type fileLoader struct {
	path string
}

func (f *fileLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	content, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	return (&contentLoader{content: string(content)}).Load(ctx)
}

func (f *fileLoader) String() string {
	return fmt.Sprintf("file %q", f.path)
}

func (f *fileLoader) Content() string {
	return f.path
}

type contentLoader struct {
	content string
}

func (f *contentLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	jsonData := []byte(f.content)
	var t struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(jsonData, &t); err != nil {
		return nil, fmt.Errorf("failed to parse credentials JSON: %w", err)
	}
	return googleoauth.CredentialsFromJSONWithType(ctx, jsonData, googleoauth.CredentialsType(t.Type), compute.CloudPlatformScope)
}

func (f *contentLoader) String() string {
	return "content <redacted>"
}

func (f *contentLoader) Content() string {
	return ""
}

type cliLoader struct{}

func (c *cliLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	return googleoauth.FindDefaultCredentials(ctx, compute.CloudPlatformScope)
}

func (c *cliLoader) String() string {
	return "gcloud CLI defaults"
}

func (c *cliLoader) Content() string {
	return ""
}

type userLoader struct{}

func (u *userLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	var content string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Multiline{
				Message: "Service Account (absolute path to file)",
				// Due to a bug in survey pkg, help message is not rendered
				Help: "The location to file that contains the service account in JSON, or the service account in JSON format",
			},
		},
	}, &content)
	if err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)
	return (&fileLoader{path: content}).Load(ctx)
}

func (u *userLoader) Content() string {
	return defaultAuthFilePath
}

// validateCredentialURLs validates external_account (Workload Identity Federation)
// credential URLs per Google's guidance that callers must validate these fields:
// https://google.aip.dev/auth/4117
// https://cloud.google.com/docs/authentication/client-libraries#external-credentials
//
// WIF with custom universe domains is not supported in the installer — all
// Google endpoint URLs are restricted to googleapis.com.
func validateCredentialURLs(credsJSON []byte) error {
	var t struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(credsJSON, &t); err != nil {
		return nil
	}
	if t.Type != "external_account" {
		return nil
	}

	var creds struct {
		TokenURL                       string `json:"token_url"`
		ServiceAccountImpersonationURL string `json:"service_account_impersonation_url"`
		UniverseDomain                 string `json:"universe_domain"`
		CredentialSource               struct {
			URL string `json:"url"`
		} `json:"credential_source"`
	}
	if err := json.Unmarshal(credsJSON, &creds); err != nil {
		return fmt.Errorf("failed to parse external account credentials: %v", err)
	}

	if creds.UniverseDomain != "" && creds.UniverseDomain != "googleapis.com" {
		return fmt.Errorf("Workload Identity Federation (external_account) with custom universe domain %q is not supported. "+
			"If you need this capability, please open an RFE with Red Hat or a GitHub issue at https://github.com/openshift/installer/issues", creds.UniverseDomain)
	}

	const expectedTokenURL = "https://sts.googleapis.com/v1/token"
	const iamPrefix = "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/"
	const wifNote = "Workload Identity Federation (external_account) credentials are only supported with googleapis.com endpoints"

	if creds.TokenURL != "" && creds.TokenURL != expectedTokenURL {
		return fmt.Errorf("token_url %q must equal %s. %s", creds.TokenURL, expectedTokenURL, wifNote)
	}
	if creds.ServiceAccountImpersonationURL != "" && !strings.HasPrefix(creds.ServiceAccountImpersonationURL, iamPrefix) {
		return fmt.Errorf("service_account_impersonation_url %q must begin with %s. %s", creds.ServiceAccountImpersonationURL, iamPrefix, wifNote)
	}
	if creds.CredentialSource.URL != "" {
		if err := validateCredSourceURL(creds.CredentialSource.URL); err != nil {
			return fmt.Errorf("credential_source.url: %v", err)
		}
	}
	return nil
}

func validateCredSourceURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %v", rawURL, err)
	}
	if parsed.Scheme != "https" {
		return fmt.Errorf("%q must use HTTPS", rawURL)
	}
	return nil
}
