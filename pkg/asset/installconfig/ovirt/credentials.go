package ovirt

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

// commands definitions
const (
	SudoCommand  = "/usr/bin/sudo"
	CpCommand    = "/usr/bin/cp"
	TrustCommand = "/usr/bin/trust"
	CurlCommand  = "/usr/bin/curl"
	ChmodCommand = "/usr/bin/chmod"
)

// EngineConfig struct - Hold all info about user environment
var EngineConfig = Config{}

// checkURLResponse performs a GET on the provided urlAddr to ensure that
// the url actually exists. Users can set skipVerify as true or false to
// avoid cert validation. In case of failure, returns error.
func checkURLResponse(urlAddr string, skipVerify bool) error {

	logrus.Debugf("Checking URL Response... urlAddr: %s skipVerify: %s", urlAddr, strconv.FormatBool(skipVerify))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(urlAddr)
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

// checkCATrust executes trust list command to validate if fqdn
// provided is trusted locally. Returns true/false and error in case of
// failure
func checkCATrust(fqdn string) (bool, error) {
	logrus.Infof("Checking if %s CA is trusted locally...", fqdn)

	stdout, err := execCommand(TrustCommand, "list")
	if err != nil {
		return false, errors.Wrapf(err, "Cannot execute command: %s", TrustCommand)
	}

	if strings.Contains(string(stdout), fqdn) {
		logrus.Info("Detected: Engine CA as trusted locally...")
		return true, nil
	}

	logrus.Warningf("Engine: %s CA is NOT trusted locally...", fqdn)

	return false, nil
}

// fileContent receives filename as argument and return the content as []byte
func fileContent(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Error reading file: %s", filename)
	}

	return content, nil
}

// askPassword will ask the password to connect to Engine API.
// The password provided will be added in the Config struct.
// If an error happens, it will ask again username for users.
func askPassword() error {
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "oVirt engine password",
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
				Message: "oVirt engine username",
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

	// Set c.URL with the API endpoint
	EngineConfig.URL = fmt.Sprintf(
		"https://%s/ovirt-engine/api",
		EngineConfig.FQDN)

	// Set PEM URL for Download
	EngineConfig.PemURL = fmt.Sprintf(
		"https://%s/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA",
		EngineConfig.FQDN)
	logrus.Debug("PEM URL: ", EngineConfig.PemURL)

	// Simple http get to check if FQDN/URL are valid
	err = checkURLResponse(EngineConfig.URL, true)
	if err != nil {
		return EngineConfig, err
	}

	// Check if CA is trusted locally
	var importCACert bool = false

	// Store if connection with Engine is secure
	var ConnectionSecure bool = false

	ConnectionSecure, err = checkCATrust(EngineConfig.FQDN)
	if err != nil {
		return EngineConfig, err
	}

	if ConnectionSecure {
		EngineConfig.Insecure = false
	} else {
		EngineConfig.Insecure = true
		message := fmt.Sprintf("Would you like to import Engine CA cert locally from: %s ?", EngineConfig.PemURL)
		err = survey.AskOne(
			&survey.Confirm{
				Message: message,
				Default: false,
				Help:    "In order to securly communicate with the oVirt engine, the certificate authority must be trusted by the local system.",
			},
			&importCACert,
			nil)
		if err != nil {
			return EngineConfig, err
		}
	}

	// If users request, let's work to import Engine CA locally
	if importCACert == true {
		logrus.Info("Downloading Engine CA cert from Engine...")
		tmpFile, err := ioutil.TempFile(os.TempDir(), "engine-")
		if err != nil {
			fmt.Println("Cannot create temporary file", err)
		}
		defer os.Remove(tmpFile.Name())

		logrus.Debugf("CA cert temporary stored: %s", tmpFile.Name())

		/* curl command */
		_, err = execCommand(CurlCommand, "-k", EngineConfig.PemURL, "-o", tmpFile.Name())
		if err != nil {
			return EngineConfig, err
		}
		logrus.Debugf("PEM file: %s", tmpFile.Name())

		var content []byte
		content, err = fileContent(tmpFile.Name())
		if err != nil {
			return EngineConfig, err
		}
		logrus.Debug(string(content))

		/* Check if CA file already exists */
		CaPath := fmt.Sprintf("/etc/pki/ca-trust/source/anchors/%s.pem", strings.ReplaceAll(EngineConfig.FQDN, ".", "-"))
		if _, err := os.Stat(CaPath); err == nil {
			return EngineConfig, err
		}

		/* Copy the CA to anchors */
		_, err = execCommand(SudoCommand, CpCommand, tmpFile.Name(), CaPath)
		if err != nil {
			return EngineConfig, err
		}

		/* chmod to allow non root users to read the cert */
		_, err = execCommand(SudoCommand, ChmodCommand, "0644", CaPath)
		if err != nil {
			return EngineConfig, err
		}

		/* Update CA database */
		_, err = execCommand(SudoCommand, "/usr/bin/update-ca-trust")
		if err != nil {
			return EngineConfig, err
		}
		EngineConfig.Insecure = false
		logrus.Infof("%s imported with success!", CaPath)
	}

	if EngineConfig.Insecure == true {
		logrus.Warning("Communication with the Engine will be insecure.")
	}

	err = askCredentials()
	if err != nil {
		return EngineConfig, err
	}

	return EngineConfig, nil
}
