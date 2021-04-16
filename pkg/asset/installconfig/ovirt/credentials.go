package ovirt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

var errHTTPNotFound = errors.New("http response 404")

// readFile reads a file provided in the args and return
// the content or in case of failure return an error
func readFile(pathFile string) ([]byte, error) {
	content, err := ioutil.ReadFile(pathFile)
	if err != nil {
		return content, errors.Wrapf(err, "failed to read file: %s", pathFile)
	}
	return content, nil
}

// Add PEM into the System Pool
func (c *clientHTTP) addTrustBundle(pemContent string, engineConfig *Config) error {
	c.certPool, _ = x509.SystemCertPool()
	if c.certPool == nil {
		logrus.Debug("failed to load cert pool.... Creating new cert pool")
		c.certPool = x509.NewCertPool()
	}

	if len(pemContent) != 0 {
		if !c.certPool.AppendCertsFromPEM([]byte(pemContent)) {
			return errors.New("unable to load certificate")
		}
		logrus.Debugf("loaded %s into the system pool: ", engineConfig.CAFile)
		engineConfig.CABundle = strings.TrimSpace(string(pemContent))
	}
	return nil
}

// downloadFile from specificed URL and store via filepath
// Return error in case of failure
func (c *clientHTTP) downloadFile() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.skipVerify,
			RootCAs:            c.certPool,
		},
	}

	if c.saveFilePath == "" {
		return errors.New("saveFilePath must be specified")
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(c.urlAddr)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%s: %w", c.urlAddr, errHTTPNotFound)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

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

// askPassword will ask the password to connect to the Engine API.
// The password provided will be added in the Config struct.
// If an error happens, it will ask again username for users.
func askPassword(c *Config) error {
	origPwdTpl := survey.PasswordQuestionTemplate
	survey.PasswordQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[Press Ctrl+C to switch username, {{ HelpInputRune }} for help]{{color "reset"}} {{end}}`

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Engine password",
				Help:    "Password for the choosen username",
			},
			Validate: survey.ComposeValidators(survey.Required, authenticated(c)),
		},
	}, &c.Password)

	survey.PasswordQuestionTemplate = origPwdTpl

	if err != nil {
		return err
	}

	return nil
}

// askUsername will ask username to connect to the Engine API.
// The username provided will be added in the Config struct.
// Returns Config and error if failure.
func askUsername(c *Config) error {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine username",
				Help:    "The username to connect to the Engine API",
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

// askQuestionTrueOrFalse generic function to ask question to users which
// requires true (Yes) or false (No) as answer
func askQuestionTrueOrFalse(question string, helpMessage string) (bool, error) {
	value := false
	err := survey.AskOne(
		&survey.Confirm{
			Message: question,
			Help:    helpMessage,
		},
		&value, survey.Required)
	if err != nil {
		return value, err
	}

	return value, nil
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

// showPEM will print information about PEM file provided in param or error
// if a failure happens
func showPEM(pemFilePath string) error {
	certpem, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read the cert: %s", pemFilePath)
	}

	block, _ := pem.Decode(certpem)
	if block == nil {
		return errors.New("failed to parse certificate PEM")
	}

	if block.Type != "CERTIFICATE" {
		return errors.New("PEM-block should be CERTIFICATE type")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Wrapf(err, "failed to read the cert: %s", pemFilePath)
	}

	logrus.Info("Loaded the following PEM file:")

	logrus.Info("\tVersion: ", cert.Version)
	logrus.Info("\tSignature Algorithm: ", cert.SignatureAlgorithm.String())
	logrus.Info("\tSerial Number: ", cert.SerialNumber)
	logrus.Info("\tIssuer: ", cert.Issuer.String())
	logrus.Info("\tValidity:")
	logrus.Info("\t\tNot Before: ", cert.NotBefore)
	logrus.Info("\t\tNot After: ", cert.NotAfter)
	logrus.Info("\tSubject: ", cert.Subject.ToRDNSequence())

	return nil

}

// askPEMFile ask users the PEM bundle and returns the bundle string
// or in case of failure returns error
func askPEMFile() (string, error) {
	bundlePEM := ""
	err := survey.AskOne(&survey.Multiline{
		Message: "Certificate bundle",
		Help:    "The certificate bundle to installer be able to communicate with oVirt API",
	},
		&bundlePEM,
		survey.ComposeValidators(survey.Required))
	if err != nil {
		return bundlePEM, err
	}

	return bundlePEM, nil
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
	err = httpResource.downloadFile()
	if errors.Is(err, errHTTPNotFound) {
		return engineConfig, err
	}

	if err != nil {
		logrus.Warning("cannot download PEM file from Engine!", err)
		answer, err := askQuestionTrueOrFalse(
			"Would you like to continue?",
			"By not using a trusted CA, insecure connections can "+
				"cause man-in-the-middle attacks among many others.")
		if err != nil || !answer {
			return engineConfig, err
		}
	} else {
		err = showPEM(httpResource.saveFilePath)
		if err != nil {
			engineConfig.Insecure = true
		} else {
			answer, err := askQuestionTrueOrFalse(
				"Would you like to use the above certificate to connect to the Engine? ",
				"Certificate to connect to the Engine. Make sure this cert CA is trusted locally.")
			if err != nil {
				return engineConfig, err
			}
			if answer {
				pemFile, err := readFile(httpResource.saveFilePath)
				engineConfig.CABundle = string(pemFile)
				if err != nil {
					return engineConfig, err
				}
				if len(engineConfig.CABundle) > 0 {
					engineConfig.Insecure = false
				}
			} else {
				answer, err = askQuestionTrueOrFalse(
					"Would you like to import another PEM bundle?",
					"You can use your own PEM bundle to connect to the Engine API")
				if err != nil {
					return engineConfig, err
				}
				if answer {
					engineConfig.CABundle, _ = askPEMFile()
					if len(engineConfig.CABundle) > 0 {
						engineConfig.Insecure = false
					}
				}

			}
		}
	}

	if !engineConfig.Insecure {
		err = httpResource.addTrustBundle(engineConfig.CABundle, &engineConfig)
		if err != nil {
			engineConfig.Insecure = true
		}
	}

	if engineConfig.Insecure {
		logrus.Error(
			"****************************************************************************\n",
			"* Could not configure secure communication to the oVirt engine.            *\n",
			"* As of 4.7 insecure mode for oVirt is no longer supported in the          *\n",
			"* installer. Please see the help article titled \"Installing OpenShift on   *\n",
			"* RHV/oVirt in insecure mode\" for details how to configure insecure mode   *\n",
			"* manually.                                                                *\n",
			"****************************************************************************",
		)
		return engineConfig,
			errors.New(
				"cannot detect engine ca cert imported in the system",
			)
	}
	return askCredentials(engineConfig)
}
