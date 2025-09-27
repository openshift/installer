package validations

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	baseDomainRegex     = `^[a-z\d]([\-]*[a-z\d]+)+$`
	dnsNameRegex        = `^([a-z\d]([\-]*[a-z\d]+)*\.)+[a-z\d]+([\-]*[a-z\d]+)+$`
	wildCardDomainRegex = `^(validateNoWildcardDNS\.).+\.?$`
)

// ValidateDomainNameFormat validates that the provided domain or FQDN conforms to DNS naming rules.
// Returns 0 and nil on success, or an HTTP status code and error on failure.
func ValidateDomainNameFormat(dnsDomainName string) (int32, error) {
	domainName := dnsDomainName
	wildCardMatched, wildCardMatchErr := regexp.MatchString(wildCardDomainRegex, dnsDomainName)
	if wildCardMatchErr == nil && wildCardMatched {
		trimmedDomain := strings.TrimPrefix(dnsDomainName, "validateNoWildcardDNS.")
		domainName = strings.TrimSuffix(trimmedDomain, ".")
	}

	matched, err := regexp.MatchString(baseDomainRegex, domainName)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Single DNS base domain validation for %s: %w", dnsDomainName, err)
	}
	if matched && len(domainName) > 1 && len(domainName) < 63 {
		return 0, nil
	}

	matched, err = regexp.MatchString(dnsNameRegex, domainName)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("DNS name validation for %s: %w", dnsDomainName, err)
	}

	if !matched || isDottedDecimalDomain(domainName) || len(domainName) > 255 {
		return http.StatusBadRequest, fmt.Errorf("DNS format mismatch: %s domain name is not valid. Must match regex [%s], be no more than 255 characters, and not be in dotted decimal format (##.##.##.##)", dnsDomainName, dnsNameRegex)
	}
	return 0, nil
}

// RFC 1123 (https://datatracker.ietf.org/doc/html/rfc1123#page-13)
// states that domains cannot resemble the format ##.##.##.##
func isDottedDecimalDomain(domain string) bool {
	regex := `([\d]+\.){3}[\d]+`
	return regexp.MustCompile(regex).MatchString(domain)
}
