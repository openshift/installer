package ovirt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

// clientHTTP struct - Hold info about http calls
type clientHTTP struct {
	saveFilePath string // Path for saving file (GET method)
	urlAddr      string // URL or Address
	skipVerify   bool   // skipt cert validatin in the http call
	certPool     *x509.CertPool
}

// EngineConfig struct - Hold all info about user environment
var EngineConfig = Config{}

// HTTPResource struct - Hold info for managing http calls
var HTTPResource = clientHTTP{}

// Import PEM into the System Pool
func (c clientHTTP) importCertIntoSystemPool(pemFilePath string) error {
	c.certPool, _ = x509.SystemCertPool()
	if c.certPool == nil {
		logrus.Debug("Failed to load cert pool.... Creating new cert pool")
		c.certPool = x509.NewCertPool()
	}
	c.certPool = x509.NewCertPool()
	logrus.Debugf("Reading file: %s", pemFilePath)
	pem, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read the cert")
	}

	logrus.Debug(string(pem))
	if len(pem) != 0 {
		logrus.Debug("trying to import...")
		if !c.certPool.AppendCertsFromPEM(pem) {
			return errors.Wrapf(err, "unable to load local certificates")
		}
		logrus.Debugf("Loaded %s into the system pool!", pemFilePath)
	}
	return nil
}

// downloadFile from specificed URL and store via filepath
// Return error in case of failure
func (c clientHTTP) downloadFile() error {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.skipVerify,
			RootCAs:            c.certPool,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(c.urlAddr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.saveFilePath == "" {
		return errors.Wrapf(err, "saveFilePath must be specificed")
	}

	out, err := os.Create(c.saveFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// checkURLResponse performs a GET on the provided urlAddr to ensure that
// the url actually exists. Users can set skipVerify as true or false to
// avoid cert validation. In case of failure, returns error.
func (c clientHTTP) checkURLResponse() error {

	logrus.Debugf("Checking URL Response... urlAddr: %s skipVerify: %s", c.urlAddr, strconv.FormatBool(c.skipVerify))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.skipVerify,
			RootCAs:            c.certPool,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(c.urlAddr)
	if err != nil {
		return errors.Wrapf(err, "Error checking URL response")
	}
	defer resp.Body.Close()

	return nil
}

// execCommand executes a command from cmdName with args
// provided from cmdArgs. Returns the stdout in []byte and error or nil
func execCommand(cmdName string, cmdArgs ...string) ([]byte, error) {
	logrus.Debugf("Executing: %s %s ", cmdName, cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "Error executing the command")
	}

	return stdout, nil
}

// askPassword will ask the password to connect to Engine API.
// The password provided will be added in the Config struct.
// If an error happens, it will ask again username for users.
func askPassword() error {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Engine password",
				Help:    "",
			},
			Validate: survey.ComposeValidators(survey.Required, authenticated(&EngineConfig)),
		},
	}, &EngineConfig.Password)
	if err != nil {
		return err
	}

	return nil
}

// askUsername will ask username to connect to Engine API.
// The username provided will be added in the Config struct.
// Returns Config and error if failure.
func askUsername() error {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine username",
				Help:    "The user must have permissions to create VMs and disks on the Storage Domain with the same name as the OpenShift cluster.",
				Default: "admin@internal",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &EngineConfig.Username)
	if err != nil {
		return err
	}

	return nil
}

// askCredentials will handle username and password for connecting with Engine
// In case of error during password, users will be prompted username again.
func askCredentials() error {
	err := askUsername()
	if err != nil {
		return err
	}

	err = askPassword()
	if err != nil {
		return askUsername()
	}

	return nil
}

// engineSetup will ask users: FQDN, execute validations and about
// the credentials. In case of failure, returns Config and error
func engineSetup() (Config, error) {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine FQDN[:PORT]",
				Help:    "The Engine FQDN[:PORT] (engine.example.com:443)",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &EngineConfig.FQDN)
	if err != nil {
		return EngineConfig, err
	}
	logrus.Debug("Engine FQDN: ", EngineConfig.FQDN)

	// By default, we set Insecure true
	EngineConfig.Insecure = true

	// Set c.URL with the API endpoint
	EngineConfig.URL = fmt.Sprintf("https://%s/ovirt-engine/api", EngineConfig.FQDN)
	logrus.Debug("Engine URL: ", EngineConfig.URL)

	// Start creating HTTPResource struct for checking if Engine FQDN is responding
	HTTPResource.skipVerify = true
	HTTPResource.urlAddr = EngineConfig.URL
	err = HTTPResource.checkURLResponse()
	if err != nil {
		return EngineConfig, err
	}

	// Set Engine PEM URL for Download
	EngineConfig.PemURL = fmt.Sprintf(
		"https://%s/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA",
		EngineConfig.FQDN)
	logrus.Debug("PEM URL: ", EngineConfig.PemURL)

	// Create tmpFile to store the Engine PEM file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "engine-")
	if err != nil {
		fmt.Println("Cannot create temporary file", err)
	}
	defer os.Remove(tmpFile.Name())

	// Download PEM
	HTTPResource.saveFilePath = tmpFile.Name()
	HTTPResource.skipVerify = true
	HTTPResource.urlAddr = EngineConfig.PemURL
	err = HTTPResource.downloadFile()
	if err == nil {
		err = HTTPResource.importCertIntoSystemPool(HTTPResource.saveFilePath)
		if err == nil {
			EngineConfig.Insecure = false
		}
	}
	logrus.Debugf("Engine PEM temporary stored: %s", HTTPResource.saveFilePath)

	if EngineConfig.Insecure == true {
		logrus.Warning("Communication with the Engine will be insecure.")
	}

	err = askCredentials()
	if err != nil {
		return EngineConfig, err
	}

	return EngineConfig, nil
}
