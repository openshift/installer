package v3

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/nutanix-cloud-native/prism-go-client"
	"github.com/nutanix-cloud-native/prism-go-client/environment/types"
)

type clientCacheMap map[string]*Client

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

// ClientCache is a cache for prism clients
type ClientCache struct {
	cache            clientCacheMap
	validationHashes map[string]string
	mtx              sync.RWMutex

	useSessionAuth bool
}

// CacheOpts is a functional option for the ClientCache
type CacheOpts func(*ClientCache)

// WithSessionAuth sets the session auth for the ClientCache
// If sessionAuth is true, the client will use session auth instead of basic auth for authentication of requests
// If sessionAuth is false, the client will use basic auth for authentication of requests
func WithSessionAuth(sessionAuth bool) CacheOpts {
	return func(c *ClientCache) {
		c.useSessionAuth = sessionAuth
	}
}

// NewClientCache returns a new ClientCache
func NewClientCache(opts ...CacheOpts) *ClientCache {
	cache := &ClientCache{
		cache:            make(clientCacheMap),
		validationHashes: make(map[string]string),
		mtx:              sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(cache)
	}

	return cache
}

// CachedClientParams define the interface that needs to be implemented by an object that will be used to create
// a cached client.
type CachedClientParams interface {
	// ManagementEndpoint returns the struct containing all information needed to construct a new client
	// and is used to calculate the validation hash for the client for the purpose of cache invalidation.
	// The validation hash is calculated based on the serialized version of the ManagementEndpoint.
	ManagementEndpoint() types.ManagementEndpoint
	// Key returns a unique key for the client that is used to store the client in the cache
	Key() string
}

// GetOrCreate returns the client for the given client name and endpoint.
// - If the client is not found in the cache, it creates a new client, adds it to the cache, and returns it
// - If the client is found in the cache, it validates whether the client is still valid by comparing validation hashes
// - If the client is found in the cache and the validation hash is the same, it returns the client
// - If the client is found in the cache and the validation hash is different, it regenerates the client, updates the cache, and returns the client
// func (c *ClientCache) GetOrCreate(clientName string, endpoint types.ManagementEndpoint, opts ...ClientOption) (*Client, error) {
func (c *ClientCache) GetOrCreate(cachedClientParams CachedClientParams, opts ...ClientOption) (*Client, error) {
	currentValidationHash, err := validationHashFromEndpoint(cachedClientParams.ManagementEndpoint())
	if err != nil {
		return nil, fmt.Errorf("failed to calculate validation hash for cachedClientParams with key %s: %w", cachedClientParams.Key(), err)
	}

	client, validationHash, err := c.get(cachedClientParams.Key())
	if err != nil {
		if !errors.Is(err, ErrorClientNotFound) {
			return nil, fmt.Errorf("failed to get client with key %s from cache: %w", cachedClientParams.Key(), err)
		}
	}

	if validationHash == currentValidationHash {
		// validation hash is the same, return the client
		return client, nil
	}

	// validation hash is different, regenerate the client
	c.Delete(cachedClientParams)

	credentials := prismgoclient.Credentials{
		URL:         cachedClientParams.ManagementEndpoint().Address.Host,
		Endpoint:    cachedClientParams.ManagementEndpoint().Address.Host,
		Insecure:    cachedClientParams.ManagementEndpoint().Insecure,
		Username:    cachedClientParams.ManagementEndpoint().ApiCredentials.Username,
		Password:    cachedClientParams.ManagementEndpoint().ApiCredentials.Password,
		SessionAuth: c.useSessionAuth,
	}

	setDefaultsForCredentials(&credentials)
	if err := validateCredentials(credentials); err != nil {
		return nil, fmt.Errorf("failed to validate credentials for cachedClientParams with key %s: %w", cachedClientParams.Key(), err)
	}

	if cachedClientParams.ManagementEndpoint().AdditionalTrustBundle != "" {
		opts = append(opts, WithPEMEncodedCertBundle([]byte(cachedClientParams.ManagementEndpoint().AdditionalTrustBundle)))
	}

	client, err = NewV3Client(credentials, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for cachedClientParams with key %s: %w", cachedClientParams.Key(), err)
	}

	c.set(cachedClientParams.Key(), currentValidationHash, client)
	return client, nil
}

func validationHashFromEndpoint(endpoint types.ManagementEndpoint) (string, error) {
	// Note: this will only work reliably as long as types.ManagementEndpoint is predictably serializable i.e. does
	// not contain a map. Due to randomized ordering of map keys in Go, we would constantly invalidate caches
	// if the ManagementEndpoint has a map.
	serializedEndpoint, err := json.Marshal(endpoint)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(serializedEndpoint)
	hashedBytes := hasher.Sum(nil)
	currentValidationHash := hex.EncodeToString(hashedBytes)

	return currentValidationHash, nil
}

func setDefaultsForCredentials(credentials *prismgoclient.Credentials) {
	if credentials.Port == "" {
		credentials.Port = "9440"
	}

	if credentials.URL == "" {
		credentials.URL = fmt.Sprintf("%s:%s", credentials.Endpoint, credentials.Port)
	}
}

func validateCredentials(credentials prismgoclient.Credentials) error {
	if credentials.Username == "" {
		return ErrorPrismUsernameNotSet
	}

	if credentials.Password == "" {
		return ErrorPrismPasswordNotSet
	}

	return nil
}

// Get returns the client and the validation hash for the given client name
func (c *ClientCache) get(clientName string) (*Client, string, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	clnt, ok := c.cache[clientName]
	if !ok {
		return nil, "", ErrorClientNotFound
	}

	validationHash, ok := c.validationHashes[clientName]
	if !ok {
		return clnt, "", nil
	}

	return clnt, validationHash, nil
}

// Set adds the client to the cache
func (c *ClientCache) set(clientName string, validationHash string, client *Client) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.cache[clientName] = client
	c.validationHashes[clientName] = validationHash
}

// Delete removes the client from the cache
func (c *ClientCache) Delete(params CachedClientParams) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.cache, params.Key())
	delete(c.validationHashes, params.Key())
}
