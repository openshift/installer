package provider

import (
	"net"
	"strings"

	"github.com/hashicorp/terraform-provider-dns/internal/hashcode"
)

func hashIPString(v interface{}) int {
	addr := v.(string)
	ip := net.ParseIP(addr)
	if ip != nil {
		return hashcode.String(ip.String())
	}
	// Preserve pre-Go 1.17 handling of leading zeroes
	return hashcode.String(stripLeadingZeros(addr))
}

func stripLeadingZeros(input string) string {
	if strings.Contains(input, ".") {
		classes := strings.Split(input, ".")
		if len(classes) != 4 {
			return input
		}
		for classIndex, class := range classes {
			if len(class) <= 1 {
				continue
			}
			classes[classIndex] = strings.TrimLeft(class, "0")
			if classes[classIndex] == "" {
				classes[classIndex] = "0"
			}
		}
		return strings.Join(classes, ".")
	}

	return input
}
