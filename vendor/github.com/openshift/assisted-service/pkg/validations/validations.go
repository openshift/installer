package validations

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	baseDomainRegex          = `^[a-z\d][\-]*[a-z\d]+$`
	dnsNameRegex             = `^([a-z\d]([\-]*[a-z\d]+)*\.)+[a-z\d]+[\-]*[a-z\d]+$`
	hostnameRegex            = `^[a-z0-9][a-z0-9\-\.]{0,61}[a-z0-9]$`
	installerArgsValuesRegex = `^[A-Za-z0-9@!#$%*()_+-=//.,";':{}\[\]]+$`
)

var allowedFlags = []string{"--append-karg", "--delete-karg", "-n", "--copy-network", "--network-dir", "--save-partlabel", "--save-partindex", "--image-url", "--image-file"}

func ValidateInstallerArgs(args []string) error {
	argsRe := regexp.MustCompile("^-+.*")
	valuesRe := regexp.MustCompile(installerArgsValuesRegex)

	for _, arg := range args {
		if argsRe.MatchString(arg) {
			if !funk.ContainsString(allowedFlags, arg) {
				return fmt.Errorf("found unexpected flag %s for installer - allowed flags are %v", arg, allowedFlags)
			}
			continue
		}

		if !valuesRe.MatchString(arg) {
			return fmt.Errorf("found unexpected chars in value %s for installer", arg)
		}
	}

	return nil
}

func ValidateDomainNameFormat(dnsDomainName string) (int32, error) {
	matched, err := regexp.MatchString(baseDomainRegex, dnsDomainName)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrapf(err, "Single DNS base domain validation for %s", dnsDomainName)
	}
	if matched && len(dnsDomainName) > 1 {
		return 0, nil
	}
	matched, err = regexp.MatchString(dnsNameRegex, dnsDomainName)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrapf(err, "DNS name validation for %s", dnsDomainName)
	}
	if !matched {
		return http.StatusBadRequest, errors.Errorf("DNS format mismatch: %s domain name is not valid", dnsDomainName)
	}
	return 0, nil
}

func ValidateHostname(name string) error {
	matched, err := regexp.MatchString(hostnameRegex, name)
	if err != nil {
		return errors.Wrapf(err, "Hostname validation for %s", name)
	}
	if !matched {
		return errors.Errorf(`Hostname format mismatch: %s name is not valid.
			Hostname must have a maximum length of 64 characters,
			start and end with a lowercase alphanumerical character,
			and can only contain lowercase alphanumerical characters, dashes, and periods.`, name)
	}
	return nil
}

func AllStrings(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func ValidateAdditionalNTPSource(commaSeparatedNTPSources string) bool {
	return AllStrings(strings.Split(commaSeparatedNTPSources, ","), ValidateNTPSource)
}

func ValidateNTPSource(ntpSource string) bool {
	if addr := net.ParseIP(ntpSource); addr != nil {
		return true
	}

	if err := ValidateHostname(ntpSource); err == nil {
		return true
	}

	return false
}

// ValidateHTTPFormat validates the HTTP and HTTPS format
func ValidateHTTPFormat(theurl string) error {
	u, err := url.Parse(theurl)
	if err != nil {
		return fmt.Errorf("URL '%s' format is not valid: %w", theurl, err)
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		return errors.Errorf("The URL scheme must be http(s) and specified in the URL: '%s'", theurl)
	}
	return nil
}

// ValidateHTTPProxyFormat validates the HTTP Proxy and HTTPS Proxy format
func ValidateHTTPProxyFormat(proxyURL string) error {
	if !govalidator.IsURL(proxyURL) {
		return errors.Errorf("Proxy URL format is not valid: '%s'", proxyURL)
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		return errors.Errorf("Proxy URL format is not valid: '%s'", proxyURL)
	}
	if u.Scheme == "https" {
		return errors.Errorf("The URL scheme must be http; https is currently not supported: '%s'", proxyURL)
	}
	if u.Scheme != "http" {
		return errors.Errorf("The URL scheme must be http and specified in the URL: '%s'", proxyURL)
	}
	return nil
}

// ValidateNoProxyFormat validates the no-proxy format which should be a comma-separated list
// of destination domain names, domains, IP addresses or other network CIDRs. A domain can be
// prefaced with '.' to include all subdomains of that domain.
func ValidateNoProxyFormat(noProxy string) error {
	if noProxy == "*" {
		return nil
	}
	domains := strings.Split(noProxy, ",")
	for _, s := range domains {
		s = strings.TrimPrefix(s, ".")
		if govalidator.IsIP(s) {
			continue
		}

		if govalidator.IsCIDR(s) {
			continue
		}

		if govalidator.IsDNSName(s) {
			continue
		}
		return errors.Errorf("NO Proxy format is not valid: '%s'. "+
			"NO Proxy is a comma-separated list of destination domain names, domains, IP addresses or other network CIDRs. "+
			"A domain can be prefaced with '.' to include all subdomains of that domain. Use '*' to bypass proxy for all destinations with OpenShift 4.8 or later.", noProxy)
	}
	return nil
}

func ValidateTags(tags string) error {
	if tags == "" {
		return nil
	}
	if !AllStrings(strings.Split(tags, ","), IsValidTag) {
		errMsg := "Invalid format for Tags: %s. Tags should be a comma-separated list (e.g. tag1,tag2,tag3). " +
			"Each tag can consist of the following characters: Alphanumeric (aA-zZ, 0-9), underscore (_) and white-spaces."
		return errors.Errorf(errMsg, tags)
	}
	return nil
}

func IsValidTag(tag string) bool {
	tagRegex := `^\w+( \w+)*$` // word characters and whitespace
	return regexp.MustCompile(tagRegex).MatchString(tag)
}

// ValidateCaCertificate ensures the specified base64 CA certificate
// is valid by trying to decode and parse it.
func ValidateCaCertificate(certificate string) error {
	decodedCaCert, err := base64.StdEncoding.DecodeString(certificate)
	if err != nil {
		return errors.Wrap(err, "failed to decode certificate")
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(decodedCaCert); !ok {
		return errors.Errorf("unable to parse certificate")
	}

	return nil
}
