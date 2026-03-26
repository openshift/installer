package v3

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	prismgoclient "github.com/nutanix-cloud-native/prism-go-client"
	"github.com/nutanix-cloud-native/prism-go-client/environment/types"
)

type clientCacheMap map[string]*Client

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

// GetOrCreate returns the client for the given client name and endpoint.
// - If the client is not found in the cache, it creates a new client, adds it to the cache, and returns it
// - If the client is found in the cache, it validates whether the client is still valid by comparing validation hashes
// - If the client is found in the cache and the validation hash is the same, it returns the client
// - If the client is found in the cache and the validation hash is different, it regenerates the client, updates the cache, and returns the client
// func (c *ClientCache) GetOrCreate(clientName string, endpoint types.ManagementEndpoint, opts ...ClientOption) (*Client, error) {
func (c *ClientCache) GetOrCreate(cachedClientParams types.CachedClientParams, opts ...types.ClientOption[Client]) (*Client, error) {
	currentValidationHash, err := validationHashFromEndpoint(cachedClientParams.ManagementEndpoint())
	if err != nil {
		return nil, fmt.Errorf("failed to calculate validation hash for cachedClientParams with key %s: %w", cachedClientParams.Key(), err)
	}

	client, validationHash, err := c.get(cachedClientParams.Key())
	if err != nil {
		if !errors.Is(err, types.ErrorClientNotFound) {
			return nil, fmt.Errorf("failed to get client with key %s from cache: %w", cachedClientParams.Key(), err)
		}
	}

	if validationHash == currentValidationHash {
		// validation hash is the same, return the client
		return client, nil
	}

	// validation hash is different, regenerate the client
	c.Delete(cachedClientParams)

	// Cache the management endpoint to avoid multiple calls and enable defensive checks
	managementEndpoint := cachedClientParams.ManagementEndpoint()

	if err := validateManagementEndpoint(managementEndpoint, cachedClientParams.Key()); err != nil {
		return nil, err
	}

	credentials := prismgoclient.Credentials{
		URL:         managementEndpoint.Address.Host,
		Endpoint:    managementEndpoint.Address.Host,
		Insecure:    managementEndpoint.Insecure,
		Username:    managementEndpoint.Username,
		Password:    managementEndpoint.Password,
		SessionAuth: c.useSessionAuth,
	}

	setDefaultsForCredentials(&credentials)
	if err := validateCredentials(credentials); err != nil {
		return nil, fmt.Errorf("failed to validate credentials for cachedClientParams with key %s: %w", cachedClientParams.Key(), err)
	}

	if managementEndpoint.AdditionalTrustBundle != "" {
		opts = append(opts, WithPEMEncodedCertBundle([]byte(managementEndpoint.AdditionalTrustBundle)))
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

func validateManagementEndpoint(endpoint types.ManagementEndpoint, key string) error {
	if endpoint.Address == nil {
		return fmt.Errorf("management endpoint address is nil for cachedClientParams with key %s", key)
	}

	// Defensive programming: validate required fields
	if endpoint.Address.Host == "" {
		return fmt.Errorf("management endpoint address host is empty for cachedClientParams with key %s", key)
	}

	if endpoint.Username == "" {
		return fmt.Errorf("API credentials username is empty for cachedClientParams with key %s", key)
	}

	if endpoint.Password == "" {
		return fmt.Errorf("API credentials password is empty for cachedClientParams with key %s", key)
	}

	return nil
}

func validateCredentials(credentials prismgoclient.Credentials) error {
	if credentials.Username == "" {
		return types.ErrorPrismUsernameNotSet
	}

	if credentials.Password == "" {
		return types.ErrorPrismPasswordNotSet
	}

	return nil
}

// Get returns the client and the validation hash for the given client name
func (c *ClientCache) get(clientName string) (*Client, string, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	clnt, ok := c.cache[clientName]
	if !ok {
		return nil, "", types.ErrorClientNotFound
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
func (c *ClientCache) Delete(params types.CachedClientParams) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.cache, params.Key())
	delete(c.validationHashes, params.Key())
}
