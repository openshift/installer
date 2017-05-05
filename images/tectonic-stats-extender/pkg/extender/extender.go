package extender

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos-inc/tectonic-licensing/license"
)

// Extender can be Run() to write extensions to an output file.
type Extender struct {
	extensions Extensions
	license    string
	output     string
	period     time.Duration
	publicKey  string
}

// New returns a configured Extender that can be run.
func New(extensions Extensions, license, output string, period time.Duration, publicKey string) *Extender {
	return &Extender{
		extensions: extensions,
		license:    license,
		output:     output,
		period:     period,
		publicKey:  publicKey,
	}
}

// Run is what actually executes a extender with a given configuration.
func (e *Extender) Run(l *log.Logger) {
	l.Info("started stats-extender")
	for {
		if err := e.runOnce(); err != nil {
			l.Warnf("failed to generate extensions: %v", err)
			l.Warnf("will attempt again after %v", e.period)
		} else {
			l.Info("successfully generated extensions")
		}
		if e.period == 0 {
			return
		}
		<-time.After(e.period)
	}
}

func (e *Extender) runOnce() error {
	license, err := e.getLicense()
	if err != nil {
		return err
	}

	e.extensions["accountID"] = license.AccountID
	e.extensions["accountSecret"] = license.AccountSecret
	extensions, err := json.Marshal(e.extensions)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(e.output, extensions, 0644)
	if err != nil {
		return err
	}

	return nil
}

// getLicense returns the decoded Tectonic license on disk.
func (e *Extender) getLicense() (*license.License, error) {
	licenseBytes, err := ioutil.ReadFile(e.license)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve license: %v", err)
	}

	publicKeyBytes := []byte(license.ProductionSigningPublicKey)
	if e.publicKey != "" {
		publicKeyBytes, err = ioutil.ReadFile(e.publicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to retieve public key: %v", err)
		}
	}

	// TODO(squat): use staging public key when not in production.
	publicKey, err := license.LoadPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %v", err)
	}

	decodedLicense, err := license.Decode(publicKey.(*rsa.PublicKey), string(licenseBytes))
	if err != nil {
		return "", fmt.Errorf("failed to decode license: %v", err)
	}

	return decodedLicense.AccountID, nil
}
