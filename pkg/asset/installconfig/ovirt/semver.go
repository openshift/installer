package ovirt

import (
	"fmt"
	"strconv"
	"strings"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const engineMinimumVersionRequired = "4.3.9.4"

// engineVersion holds all information about versioning
type engineVersion struct {
	major       int // Major Version
	minor       int // Minor Version
	maintenance int // Maintenance Version
	build       int // Build Version
}

// parseRelease parse the Engine release string from argument.
// It also removes "-el7" from the release and returns
// a engineVersion struct type. In case the error, it will return
// an error message.
func (e *engineVersion) parseRelease(version string) error {
	// Removes "-el7", "-el8" etc.
	if strings.Contains(version, "-") {
		version = strings.Split(version, "-")[0]
	}

	// Split the release by .
	for k, v := range strings.Split(version, ".") {
		intNumber, err := strconv.Atoi(v)
		if err != nil {
			return errors.Wrap(err, "cannot convert release to int()")
		}
		switch k {
		case 0:
			e.major = intNumber
		case 1:
			e.minor = intNumber
		case 2:
			e.maintenance = intNumber
		case 3:
			e.build = intNumber
		default:
			return errors.New("cannot parse Engine version")
		}
	}
	return nil
}

// checkReleaseSupport compare two engineVersion type, the first param is the current
// engine version and the second is the required version to run OCP on oVirt.
// Returns true in case the version is support or false and error message.
func checkReleaseSupport(current engineVersion, required engineVersion) error {

	msg := fmt.Sprintf(
		"version is not supported for deploying OCP. Engine %s or higher is required",
		engineMinimumVersionRequired)

	if current.major < required.major {
		msg = fmt.Sprintf("MAJOR %s", msg)
		return errors.New(msg)
	}
	if current.minor < required.minor {
		msg = fmt.Sprintf("MINOR %s", msg)
		return errors.New(msg)
	}
	if current.maintenance < required.maintenance {
		msg = fmt.Sprintf("MAINTENANCE %s", msg)
		return errors.New(msg)
	}
	if current.build < required.build {
		msg = fmt.Sprintf("BUILD %s", msg)
		return errors.New(msg)
	}
	return nil
}

// engineMinimumVersionRequirement connects to the Engine API collects the current
// version and compare with the minimum version supported to deploy OCP.
// If error occurs, error type will be returned.
func engineMinimumVersionRequirement(c *ovirtsdk4.Connection) error {
	api := c.SystemService().Get().MustSend().MustApi()
	engineVer := api.MustProductInfo().MustVersion().MustFullVersion()
	logrus.Info("Engine version detected: ", engineVer)

	currentVersion := engineVersion{}
	err := currentVersion.parseRelease(engineVer)
	if err != nil {
		return errors.Wrap(err, "cannot parse current Engine version")
	}

	requiredVersion := engineVersion{}
	err = requiredVersion.parseRelease(engineMinimumVersionRequired)
	if err != nil {
		return errors.Wrap(err, "cannot parse required Engine version")
	}

	err = checkReleaseSupport(currentVersion, requiredVersion)
	if err != nil {
		return errors.Wrap(err, "release not supported")
	}
	return nil
}
