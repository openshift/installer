package powervs

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	cosHostPrefix = "s3."
	cosHostSuffix = ".cloud-object-storage.appdomain.cloud"
)

// ParseOSImageURL splits an IBM Cloud Object Storage image URL of the form
// https://s3.<region>.cloud-object-storage.appdomain.cloud/<bucket>/<object>
// into its bucket, object and region components.
func ParseOSImageURL(rawURL string) (bucket, object, region string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "https" {
		return "", "", "", fmt.Errorf("URL scheme must be https")
	}

	host := u.Hostname()
	if !strings.HasPrefix(host, cosHostPrefix) || !strings.HasSuffix(host, cosHostSuffix) {
		return "", "", "", fmt.Errorf("host must be of the form s3.<region>%s", cosHostSuffix)
	}
	region = strings.TrimSuffix(strings.TrimPrefix(host, cosHostPrefix), cosHostSuffix)
	if region == "" {
		return "", "", "", fmt.Errorf("region could not be derived from host %q", host)
	}

	parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", fmt.Errorf("URL path must be of the form /<bucket>/<object>")
	}

	return parts[0], parts[1], region, nil
}
