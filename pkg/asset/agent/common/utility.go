package common

import (
	"github.com/google/uuid"
)

// GenerateInfraEnvID returns a random infra env ID.
func GenerateInfraEnvID() string {
	return uuid.New().String()
}
