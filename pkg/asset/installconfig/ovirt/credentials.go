package ovirt

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"os"
	"os/exec"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

// Global vars for commands definitions
const (
	SUDO_BIN = "/usr/bin/sudo"
	CP_BIN = "/usr/bin/cp"
	TRUST_BIN = "/usr/bin/trust"
	CURL_BIN = "/usr/bin/curl"
	CHMOD_BIN = "/usr/bin/chmod"
)

// Check if URL can be reached before we proceed with the installation
// Params:
//	urlAddr - Full URL
//	skipVerify - Do not try to validate cert
func checkURLResponse(urlAddr string, skipVerify bool) {

	logrus.Debugf("Checking URL Response... urlAddr: %s skipVerify: %s", urlAddr, strconv.FormatBool(skipVerify))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : skipVerify},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(urlAddr)
	if err != nil {
		logrus.Fatalf("Error checking URL response: %s err: %s", urlAddr, err)
	}
	defer resp.Body.Close()
}

// Execute a command
// Params:
//	cmdName - Command name
//	cmdArgs - Arguments for the command
//
// Return - stdout ([]byte)
func execCommand(cmdName string, cmdArgs []string) ([]byte){
	logrus.Debugf("Executing: %s %s ", cmdName, cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		logrus.Fatalf("Error executing the command: %s", err)
	}

	return stdout
}

// Check if CA is trusted locally
// Params:
//	fqdn - Fully qualified domain name
//
// Return:
//	True  - CA is trusted
//	False - CA is NOT trusted
func checkCATrust(fqdn string) (bool){
	logrus.Infof("Checking if %s CA is trusted locally...", fqdn)

	LIST := []string {"list"}
	stdout := execCommand(TRUST_BIN, LIST)

	if strings.Contains(string(stdout), fqdn) {
		logrus.Info("Detected: Engine CA as trusted locally...")
		return true
	}

	logrus.Warningf("Engine: %s CA is NOT trusted locally...", fqdn)

	return false
}

// Return file content
// Params:
//	filename - File name
// Return:
//	file content - []byte
func fileContent(filename string) ([]byte){
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Fatalf("Error reading file: %s Error: %s", filename, err)
	}

	return content
}

func askCredentials() (Config, error) {
	oVirtConfig := Config{}
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine FQDN[:PORT]",
				Help:    "The Engine FQDN[:PORT] (engine.example.com:443)",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &oVirtConfig.FQDN)
	if err != nil {
		return oVirtConfig, err
	}

	// Set c.URL with the API endpoint
	oVirtConfig.URL = fmt.Sprintf(
		"https://%s/ovirt-engine/api",
		oVirtConfig.FQDN)

	// Set PEM URL for Download
	oVirtConfig.PemURL = fmt.Sprintf(
		"https://%s/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA",
		oVirtConfig.FQDN)
	logrus.Debug("PEM URL: ", oVirtConfig.PemURL)

	// Simple http get to check if FQDN/URL are valid
	checkURLResponse(oVirtConfig.URL, true)

	// Check if CA is trusted locally
	var importCACert bool = false
	if checkCATrust(oVirtConfig.FQDN) {
		oVirtConfig.Insecure = false
	} else {
		oVirtConfig.Insecure = true
		message := fmt.Sprintf("Would you like to import Engine CA cert locally from: %s ?", oVirtConfig.PemURL)
		err = survey.AskOne(
		&survey.Confirm{
			Message: message,
			Default: false,
			Help:    "In order to securly communicate with the Engine, the certificate authority must be trusted by the local system.",
		},
		&importCACert,
		nil)
	}
	if err != nil {
		return oVirtConfig, err
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

		/* Build curl command */
		OPTIONS_CURL := []string {"-k", oVirtConfig.PemURL, "-o", tmpFile.Name()}
		execCommand(CURL_BIN, OPTIONS_CURL)
		logrus.Debugf("PEM file: %s", tmpFile.Name())
		logrus.Debug(string(fileContent(tmpFile.Name())))

		/* Check if CA file already exists */
		CA_FULL_PATH := fmt.Sprintf("/etc/pki/ca-trust/source/anchors/%s.pem", strings.ReplaceAll(oVirtConfig.FQDN, ".", "-"))
		if _, err := os.Stat(CA_FULL_PATH); err == nil {
			logrus.Fatalf("%s already exists, cannot import CA!", CA_FULL_PATH)
		}

		/* Copy the CA to anchors */
		CA_TO_ANCHORS := []string {CP_BIN, tmpFile.Name(), CA_FULL_PATH}
		execCommand(SUDO_BIN, CA_TO_ANCHORS)

		/* chmod to allow non root users to read the cert */
		PERMS_PEMFILE := []string {CHMOD_BIN, "0644", CA_FULL_PATH}
		execCommand(SUDO_BIN, PERMS_PEMFILE)

		/* Update CA database */
		UPDATE_CA_TRUST := []string {"/usr/bin/update-ca-trust"}
		execCommand(SUDO_BIN, UPDATE_CA_TRUST)
		oVirtConfig.Insecure = false
		logrus.Infof("%s imported with success!", CA_FULL_PATH)
	}

	if oVirtConfig.Insecure == true {
		logrus.Warning("Communication with the Engine will be insecure.")
	}

	// Ask Engine credentials
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Engine username",
				Help:    "The user must have permissions to create VMs and disks on the Storage Domain with the same name as the OpenShift cluster.",
				Default: "admin@internal",
			},
			Validate: survey.ComposeValidators(survey.Required),
		},
	}, &oVirtConfig.Username)
	if err != nil {
		return oVirtConfig, err
	}

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Engine password",
				Help:    "",
			},
			Validate: survey.ComposeValidators(survey.Required, authenticated(&oVirtConfig)),
		},
	}, &oVirtConfig.Password)
	if err != nil {
		return oVirtConfig, err
	}

	return oVirtConfig, nil
}
