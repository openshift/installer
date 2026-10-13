package securestore

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/99designs/keyring"
	gokeyring "github.com/zalando/go-keyring"
)

const (
	KindInternetPassword = "Internet password" // MacOS Keychain item kind
	ItemKey              = "RedHatSSO"
	CollectionName       = "login" // Common OS default collection name
	MaxWindowsByteSize   = 2500    // Windows Credential Manager has a 2500 byte limit
)

var (
	ErrKeyringUnavailable = fmt.Errorf("keyring is valid but is not available on the current OS")
	ErrKeyringInvalid     = fmt.Errorf("keyring is invalid, expected one of: [%v]", strings.Join(AllowedBackends, ", "))
	AllowedBackends       = []string{
		string(keyring.WinCredBackend),
		string(keyring.KeychainBackend),
		string(keyring.SecretServiceBackend),
		string(keyring.PassBackend),
	}
)

func getKeyringConfig(backend string) keyring.Config {
	return keyring.Config{
		AllowedBackends: []keyring.BackendType{keyring.BackendType(backend)},
		// Generic
		ServiceName: ItemKey,
		// MacOS
		KeychainName:                   CollectionName,
		KeychainTrustApplication:       true,
		KeychainSynchronizable:         false,
		KeychainAccessibleWhenUnlocked: false,
		// Windows
		WinCredPrefix: ItemKey,
		// Secret Service
		LibSecretCollectionName: CollectionName,
	}
}

// IsBackendAvailable provides validation that the desired backend is available on the current OS.
func IsBackendAvailable(backend string) (isAvailable bool) {
	if backend == "" {
		return false
	}

	for _, avail := range AvailableBackends() {
		if avail == backend {
			isAvailable = true
			break
		}
	}

	return isAvailable
}

// AvailableBackends provides a slice of all available backend keys on the current OS.
func AvailableBackends() []string {
	b := []string{}

	if isDarwin() {
		// Assume Keychain is always available on Darwin. It will not return from keyring.AvailableBackends()
		b = append(b, "keychain")
	}

	// Intersection between available backends from OS and allowed backends
	for _, avail := range keyring.AvailableBackends() {
		for _, allowed := range AllowedBackends {
			if string(avail) == allowed {
				b = append(b, allowed)
			}
		}
	}

	return b
}

// UpsertConfigToKeyring will upsert the provided credentials to the desired OS secure store.
func UpsertConfigToKeyring(backend string, creds []byte) error {
	if err := ValidateBackend(backend); err != nil {
		return err
	}

	if isDarwin() && isKeychain(backend) {
		return keychainUpsert(creds)
	}

	ring, err := keyring.Open(getKeyringConfig(backend))
	if err != nil {
		return err
	}

	compressed, err := compressConfig(creds)
	if err != nil {
		return err
	}

	// check if available backend contains windows credential manager and exceeds the byte limit
	if len(compressed) > MaxWindowsByteSize &&
		backend == string(keyring.WinCredBackend) {
		return fmt.Errorf("credentials are too large for Windows Credential Manager: %d bytes (max %d)", len(compressed), MaxWindowsByteSize)
	}

	err = ring.Set(keyring.Item{
		Label:       ItemKey,
		Key:         ItemKey,
		Description: KindInternetPassword,
		Data:        compressed,
	})

	return err
}

// RemoveConfigFromKeyring will remove the credentials from the first priority OS secure store.
func RemoveConfigFromKeyring(backend string) error {
	if err := ValidateBackend(backend); err != nil {
		return err
	}

	if isDarwin() && isKeychain(backend) {
		return keychainRemove()
	}

	ring, err := keyring.Open(getKeyringConfig(backend))
	if err != nil {
		return err
	}

	err = ring.Remove(ItemKey)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			// Ignore not found errors, key is already removed
			return nil
		}
	}

	return err
}

// GetConfigFromKeyring will retrieve the credentials from the first priority OS secure store.
func GetConfigFromKeyring(backend string) ([]byte, error) {
	if err := ValidateBackend(backend); err != nil {
		return nil, err
	}

	if isDarwin() && isKeychain(backend) {
		return keychainGet()
	}

	credentials := []byte("")

	ring, err := keyring.Open(getKeyringConfig(backend))
	if err != nil {
		return nil, err
	}

	i, err := ring.Get(ItemKey)
	if err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
		return credentials, err
	} else if errors.Is(err, keyring.ErrKeyNotFound) {
		// Not found, continue
	} else {
		credentials = i.Data
	}

	if len(credentials) == 0 {
		// No creds to decompress, return early
		return credentials, nil
	}

	creds, err := decompressConfig(credentials)
	if err != nil {
		return nil, err
	}

	return creds, nil

}

// Validates that the requested backend is valid and available, returns an error if not.
func ValidateBackend(backend string) error {
	if backend == "" {
		return ErrKeyringInvalid
	} else {
		isAllowedBackend := false
		for _, allowed := range AllowedBackends {
			if allowed == backend {
				isAllowedBackend = true
				break
			}
		}
		if !isAllowedBackend {
			return ErrKeyringInvalid
		}
	}

	if !IsBackendAvailable(backend) {
		return ErrKeyringUnavailable
	}

	return nil
}

func keychainGet() ([]byte, error) {
	credentials, err := gokeyring.Get(ItemKey, ItemKey)
	if err != nil && !errors.Is(err, gokeyring.ErrNotFound) {
		return []byte(credentials), err
	} else if errors.Is(err, gokeyring.ErrNotFound) {
		return []byte(""), nil
	}

	if len(credentials) == 0 {
		// No creds to decompress, return early
		return []byte(""), nil
	}

	creds, err := decompressConfig([]byte(credentials))
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func keychainUpsert(creds []byte) error {
	compressed, err := compressConfig(creds)
	if err != nil {
		return err
	}

	err = gokeyring.Set(ItemKey, ItemKey, string(compressed))
	if err != nil {
		return err
	}

	return nil
}

func keychainRemove() error {
	err := gokeyring.Delete(ItemKey, ItemKey)
	if err != nil {
		if errors.Is(err, gokeyring.ErrNotFound) {
			// Ignore not found errors, key is already removed
			return nil
		}
		if strings.Contains(err.Error(), "Keychain Error. (-25244)") {
			return fmt.Errorf("%s\nThis application may not have permission to delete from the Keychain. Please check the permissions in the Keychain and try again", err.Error())
		}
	}

	return err
}

// Compresses credential bytes to help ensure all OS secure stores can store the data.
// Windows Credential Manager has a 2500 byte limit.
func compressConfig(creds []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err := gz.Write(creds)
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Decompresses credential bytes
func decompressConfig(creds []byte) ([]byte, error) {
	reader := bytes.NewReader(creds)
	gzreader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(gzreader)
	if err != nil {
		return nil, err
	}

	return output, err
}

// isDarwin returns true if the current OS runtime is "darwin"
func isDarwin() bool {
	return runtime.GOOS == "darwin"
}

// isKeychain returns true if the backend is "keychain"
func isKeychain(backend string) bool {
	return backend == "keychain"
}
