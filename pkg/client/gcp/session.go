package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
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
	authEnvs            = []string{"GOOGLE_CREDENTIALS", "GOOGLE_CLOUD_KEYFILE_JSON", "GCLOUD_KEYFILE_JSON"}
	defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".gcp", "osServiceAccount.json")
	credLoaders         = []credLoader{}
	onceLoggers         = map[credLoader]*sync.Once{}
)

// Session is an object representing session for GCP API.
type Session struct {
	Credentials *googleoauth.Credentials
}

// GetSession returns a GCP session by using credentials found in default locations in order:
// env GOOGLE_CREDENTIALS,
// env GOOGLE_CLOUD_KEYFILE_JSON,
// env GCLOUD_KEYFILE_JSON,
// file ~/.gcp/osServiceAccount.json, and
// gcloud cli defaults
// and, if no creds are found, asks for them and stores them on disk in a config file
func GetSession(ctx context.Context) (*Session, error) {
	creds, err := loadCredentials(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}

	return &Session{
		Credentials: creds,
	}, nil
}

func loadCredentials(ctx context.Context) (*googleoauth.Credentials, error) {
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
		return creds, nil
	}
	return getCredentials(ctx)
}

func getCredentials(ctx context.Context) (*googleoauth.Credentials, error) {
	creds, err := (&userLoader{}).Load(ctx)
	if err != nil {
		return nil, err
	}

	filePath := defaultAuthFilePath
	logrus.Infof("Saving the credentials to %q", filePath)
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(filePath, creds.JSON, 0600); err != nil {
		return nil, err
	}
	return creds, nil
}

type credLoader interface {
	Load(context.Context) (*googleoauth.Credentials, error)
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

type fileOrContentLoader struct {
	pathOrContent string
	delegate      credLoader
}

func (fc *fileOrContentLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(fc.pathOrContent); err == nil {
		fc.delegate = &fileLoader{path: fc.pathOrContent}
	} else {
		fc.delegate = &contentLoader{content: fc.pathOrContent}
	}

	return fc.delegate.Load(ctx)
}

func (fc *fileOrContentLoader) String() string {
	if fc.delegate != nil {
		return fmt.Sprintf("%s", fc.delegate)
	}
	return "file or content"
}

type fileLoader struct {
	path string
}

func (f *fileLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	content, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	return (&contentLoader{content: string(content)}).Load(ctx)
}

func (f *fileLoader) String() string {
	return fmt.Sprintf("file %q", f.path)
}

type contentLoader struct {
	content string
}

func (f *contentLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	return googleoauth.CredentialsFromJSON(ctx, []byte(f.content), compute.CloudPlatformScope)
}

func (f *contentLoader) String() string {
	return "content <redacted>"
}

type cliLoader struct{}

func (c *cliLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	return googleoauth.FindDefaultCredentials(ctx, compute.CloudPlatformScope)
}

func (c *cliLoader) String() string {
	return "gcloud CLI defaults"
}

type userLoader struct{}

func (u *userLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	var content string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Multiline{
				Message: "Service Account (absolute path to file or JSON content)",
				// Due to a bug in survey pkg, help message is not rendered
				Help: "The location to file that contains the service account in JSON, or the service account in JSON format",
			},
		},
	}, &content)
	if err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)
	return (&fileOrContentLoader{pathOrContent: content}).Load(ctx)
}
