package common

import "github.com/google/uuid"

// InfraEnvID holds an uniuqe identifier for infra env resource.
var InfraEnvID string

func init() {
	// Generate a UUID during initialization
	InfraEnvID = uuid.New().String()
}
