package dns

import (
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

func hashIPString(v interface{}) int {
	addr := v.(string)
	ip := net.ParseIP(addr)
	if ip == nil {
		return hashcode.String(addr)
	}
	return hashcode.String(ip.String())
}
