package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

var (
	// ErrorClientNotFound is returned when the client is not found in the cache
	ErrorClientNotFound = errors.New("client not found in client cache")
	// ErrorPrismAddressNotSet is returned when the address is not set for Nutanix Prism Central
	ErrorPrismAddressNotSet = errors.New("address not set for Nutanix Prism Central")
	// ErrorPrismPortNotSet is returned when the port is not set for Nutanix Prism Central
	ErrorPrismPortNotSet = errors.New("port not set for Nutanix Prism Central")
	// ErrorPrismUsernameNotSet is returned when the username is not set for Nutanix Prism Central
	ErrorPrismUsernameNotSet = errors.New("username not set for Nutanix Prism Central")
	// ErrorPrismPasswordNotSet is returned when the password is not set for Nutanix Prism Central
	ErrorPrismPasswordNotSet = errors.New("password not set for Nutanix Prism Central")
)

// CachedClientParams define the interface that needs to be implemented by an object that will be used to create
// a cached client.
type CachedClientParams interface {
	// ManagementEndpoint returns the struct containing all information needed to construct a new client
	// and is used to calculate the validation hash for the client for the purpose of cache invalidation.
	// The validation hash is calculated based on the serialized version of the ManagementEndpoint.
	ManagementEndpoint() ManagementEndpoint
	// Key returns a unique key for the client that is used to store the client in the cache
	Key() string
}

type CacheOpts[T any] func(*T)
type ClientOption[T any] func(*T) error

func (ep ManagementEndpoint) GetHash() (string, error) {
	// Note: this will only work reliably as long as types.ManagementEndpoint is predictably serializable i.e. does
	// not contain a map. Due to randomized ordering of map keys in Go, we would constantly invalidate caches
	// if the ManagementEndpoint has a map.
	serializedEndpoint, err := json.Marshal(ep)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(serializedEndpoint)
	hashedBytes := hasher.Sum(nil)
	currentValidationHash := hex.EncodeToString(hashedBytes)

	return currentValidationHash, nil
}
