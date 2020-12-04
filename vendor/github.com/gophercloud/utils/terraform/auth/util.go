package auth

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// This is copied directly from Terraform in order to remove a single legacy
// vendor dependency.

const uaEnvVar = "TF_APPEND_USER_AGENT"

func terraformUserAgent(version, sdkVersion string) string {
	ua := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io)", version)
	if sdkVersion != "" {
		ua += " " + fmt.Sprintf("Terraform Plugin SDK/%s", sdkVersion)
	}

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
