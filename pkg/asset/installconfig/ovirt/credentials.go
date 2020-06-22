package ovirt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

// Add PEM into the System Pool
func (c *clientHTTP) addTrustBundle(pemFilePath string) error {
	c.certPool, _ = x509.SystemCertPool()
	if c.certPool == nil {
		logrus.Debug("failed to load cert pool.... Creating new cert pool")
		c.certPool = x509.NewCertPool()
	}

	pem, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read the cert: %s", pemFilePath)
	}

	if len(pem) != 0 {
		if !c.certPool.AppendCertsFromPEM(pem) {
			return errors.Wrapf(err, "unable to load local certificate: %s", pemFilePath)
		}
		logrus.Debugf("loaded %s into the system pool: ", pemFilePath)
	}
	return nil
}

// downloadFile from specificed URL and store via filepath
// Return error in case of failure
func (c *clientHTTP) downloadFile() (int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.skipVerify,
			RootCAs:            c.certPool,
		},
	}

	if c.saveFilePath == "" {
		return http.StatusNotFound, errors.New("saveFilePath must be specificed")
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(c.urlAddr)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return resp.StatusCode, errors.Errorf("http response 404 for: %s", c.urlAddr)
	}

	if err != nil {
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	out, err := os.Create(c.saveFilePath)
	if err != nil {
		return resp.StatusCode, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return resp.StatusCode, err
}

// checkURLResponse performs a GET on the provided urlAddr to ensure that
// the url actually exists. Users can set skipVerify as true or false to
// avoid cert validation. In case of failure, returns error.
func (c *clientHTTP) checkURLResponse() error {

	logrus.Debugf("checking URL response... urlAddr: %s skipVerify: %s", c.urlAddr, strconv.FormatBool(c.skipVerify))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.skipVerify,
			RootCAs:            c.certPool,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(c.urlAddr)
	if err != nil {
		return errors.Wrapf(err, "error checking URL response")
	}
	defer resp.Body.Close()

	return nil
}

// askPassword will ask the password to connect to Engine API.
// The password provided will be added in the Config struct.
// If an error happens, it will ask again username for users.
func askPassword(c *Config) error {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Engine password",
				Help:    "",
			},
			Validate: survey.ComposeValidators(survey.Required, authenticated(c)),
		},
	}, &c.Password)
	if err != nil {
		return err
	}

	return nil
}

// askUsername will ask username to connect to Engine API.
// The username provided will be added in the Config struct.
// Returns Config and error if failure.
func askUsername(c *Config) error {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine username",
				Help:    "The username to connect to Engine API",
				Default: "admin@internal",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &c.Username)
	if err != nil {
		return err
	}

	return nil
}

// askCredentials will handle username and password for connecting with Engine
func askCredentials(c Config) (Config, error) {
	loginAttempts := 3
	logrus.Debugf("login attempts available: %d", loginAttempts)
	for loginAttempts > 0 {
		err := askUsername(&c)
		if err != nil {
			return c, err
		}

		err = askPassword(&c)
		if err != nil {
			loginAttempts = loginAttempts - 1
			logrus.Debugf("login attempts now: %d", loginAttempts)
			if loginAttempts == 0 {
				return c, err
			}
		} else {
			break
		}
	}
	return c, nil
}

// engineSetup will ask users: FQDN, execute validations and about
// the credentials. In case of failure, returns Config and error
func engineSetup() (Config, error) {
	engineConfig := Config{}
	httpResource := clientHTTP{}

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine FQDN[:PORT]",
				Help:    "The Engine FQDN[:PORT] (engine.example.com:443)",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &engineConfig.FQDN)
	if err != nil {
		return engineConfig, err
	}
	logrus.Debug("engine FQDN: ", engineConfig.FQDN)

	// By default, we set Insecure true
	engineConfig.Insecure = true

	// Set c.URL with the API endpoint
	engineConfig.URL = fmt.Sprintf("https://%s/ovirt-engine/api", engineConfig.FQDN)
	logrus.Debug("Engine URL: ", engineConfig.URL)

	// Start creating clientHTTP struct for checking if Engine FQDN is responding
	httpResource.skipVerify = true
	httpResource.urlAddr = engineConfig.URL
	err = httpResource.checkURLResponse()
	if err != nil {
		return engineConfig, err
	}

	// Set Engine PEM URL for Download
	engineConfig.PemURL = fmt.Sprintf(
		"https://%s/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA",
		engineConfig.FQDN)
	logrus.Debug("PEM URL: ", engineConfig.PemURL)

	// Create tmpFile to store the Engine PEM file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "engine-")
	if err != nil {
		return engineConfig, err
	}
	defer os.Remove(tmpFile.Name())

	// Download PEM
	httpResource.saveFilePath = tmpFile.Name()
	httpResource.skipVerify = true
	httpResource.urlAddr = engineConfig.PemURL
	resp, err := httpResource.downloadFile()
	if resp == http.StatusNotFound {
		return engineConfig, err
	}

	if err != nil {
		logrus.Warning("cannot download PEM file from Engine!", err)
		engineConfig.Insecure = true
	} else {
		err = httpResource.addTrustBundle(httpResource.saveFilePath)
		if err != nil {
			engineConfig.Insecure = true
		} else {
			engineConfig.Insecure = false
		}
		logrus.Debugf("engine PEM temporary stored: %s", httpResource.saveFilePath)
	}

	if engineConfig.Insecure == true {
		logrus.Warning("cannot detect Engine CA cert imported in the system. Communication with the Engine will be insecure.")
	}

	return askCredentials(engineConfig)
}
