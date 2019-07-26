package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	googleoauth "golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	authEnvs            = []string{"GOOGLE_CREDENTIALS", "GOOGLE_CLOUD_KEYFILE_JSON", "GCLOUD_KEYFILE_JSON"}
	defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".gcp", "osServiceAccount.json")
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
	var loaders []credLoader
	for _, env := range authEnvs {
		loaders = append(loaders, &envLoader{env: env})
	}
	loaders = append(loaders, &fileLoader{path: defaultAuthFilePath})
	loaders = append(loaders, &cliLoader{})

	for _, l := range loaders {
		creds, err := l.Load(ctx)
		if err != nil {
			continue
		}
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
	env string
}

func (e *envLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	if val := os.Getenv(e.env); len(val) > 0 {
		return (&fileOrContentLoader{pathOrContent: val}).Load(ctx)
	}
	return nil, errors.New("empty environment variable")
}

func (e *envLoader) String() string {
	return fmt.Sprintf("loading from environment variable %q", e.env)
}

type fileOrContentLoader struct {
	pathOrContent string
}

func (fc *fileOrContentLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(fc.pathOrContent); err == nil {
		return (&fileLoader{path: fc.pathOrContent}).Load(ctx)
	}

	return (&contentLoader{content: fc.pathOrContent}).Load(ctx)
}

func (fc *fileOrContentLoader) String() string {
	return fmt.Sprintf("loading from file or content %q", fc.pathOrContent)
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
	return fmt.Sprintf("loading from file %q", f.path)
}

type contentLoader struct {
	content string
}

func (f *contentLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	return googleoauth.CredentialsFromJSON(ctx, []byte(f.content), compute.CloudPlatformScope)
}

func (f *contentLoader) String() string {
	return fmt.Sprintf("loading from content %q", f.content)
}

type cliLoader struct{}

func (c *cliLoader) Load(ctx context.Context) (*googleoauth.Credentials, error) {
	return googleoauth.FindDefaultCredentials(ctx, compute.CloudPlatformScope)
}

func (c *cliLoader) String() string {
	return fmt.Sprintf("loading from gcloud defaults")
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
