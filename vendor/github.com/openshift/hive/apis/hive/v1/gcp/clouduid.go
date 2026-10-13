package gcp

import (
	"crypto/md5"
	"fmt"
)

// CloudControllerUID generates a UID used by the GCP cloud controller provider
// to generate certain load balancing resources
func CloudControllerUID(infraID string) string {
	hash := md5.Sum([]byte(infraID))
	return fmt.Sprintf("%x", hash)[:16]
}
