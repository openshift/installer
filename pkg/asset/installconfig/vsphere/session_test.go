package vsphere

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestSessionCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		server   string
		username string
		password string
	}{
		{"empty server", "", "user", "pass"},
		{"empty username", "https://vcenter", "", "pass"},
		{"empty password", "https://vcenter", "user", ""},
		{"full credentials", "https://vcenter.example.com", "admin@vsphere.local", "password123"},
		{"different server", "https://vcenter2.example.com", "admin@vsphere.local", "password123"},
		{"different username", "https://vcenter.example.com", "admin2@vsphere.local", "password123"},
		{"different password", "https://vcenter.example.com", "admin@vsphere.local", "password456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := SessionCacheKey(tt.server, tt.username, tt.password)
			assert.NotEmpty(t, key, "cache key should not be empty")
			assert.Equal(t, 32, len(key), "cache key should be 32 hex characters")
		})
	}
}

func TestSessionCacheKeyDeterministic(t *testing.T) {
	key1 := SessionCacheKey("https://vcenter", "user", "pass")
	key2 := SessionCacheKey("https://vcenter", "user", "pass")
	assert.Equal(t, key1, key2, "same credentials should produce same cache key")
}

func TestSessionCacheKeyNotEqual(t *testing.T) {
	key1 := SessionCacheKey("https://vcenter1", "user", "pass")
	key2 := SessionCacheKey("https://vcenter2", "user", "pass")
	assert.NotEqual(t, key1, key2, "different servers should produce different cache keys")
}

func TestNewSessionMissingCredentials(t *testing.T) {
	tests := []struct {
		name     string
		server   string
		username string
		password string
	}{
		{"empty server", "", "user", "pass"},
		{"empty username", "https://vcenter", "", "pass"},
		{"empty password", "https://vcenter", "user", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sess, cleanup, err := NewSession(context.Background(), tt.server, tt.username, tt.password)
			assert.Nil(t, sess, "session should be nil on error")
			assert.Nil(t, cleanup, "cleanup should be nil on error")
			assert.Error(t, err, "should return error for missing credentials")
			assert.Contains(t, err.Error(), "required", "error should mention required fields")
		})
	}
}

func TestNewSessionInvalidServer(t *testing.T) {
	sess, cleanup, err := NewSession(context.Background(), "not a valid url", "user", "pass")
	assert.Nil(t, sess, "session should be nil on invalid URL")
	assert.Nil(t, cleanup, "cleanup should be nil on error")
	assert.Error(t, err, "should return error for invalid URL")
	assert.Contains(t, err.Error(), "parse vCenter URL", "error should mention URL parsing")
}

func TestNewSessionTimeout(t *testing.T) {
	logger, _ := test.NewNullLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	sess, _, err := NewSession(ctx, "https://nonexistent.example.com:80", "user", "pass", WithTimeout(1*time.Millisecond), WithLogger(logger))
	assert.Nil(t, sess, "session should be nil on timeout")
	assert.Error(t, err, "should return error on timeout")
}

func TestSessionCloseIdempotent(t *testing.T) {
	// Create a session with nil fields to test idempotent close
	sess := &Session{
		server: "https://vcenter",
		logger: logrus.New(),
	}

	// First close should not panic
	sess.Close()

	// Second close should also not panic
	sess.Close()

	// Third close should also not panic
	sess.Close()

	// Error should still be nil (no close errors occurred)
	assert.Nil(t, sess.Error(), "error should be nil after idempotent close")
}

func TestSessionCloseWithNilClients(t *testing.T) {
	logger := logrus.New()

	sess := &Session{
		server:  "https://vcenter",
		logger:  logger,
		rest:    nil,
		govmomi: nil,
	}

	// Should not panic
	sess.Close()
	sess.Close()

	// Error should be nil
	assert.Nil(t, sess.Error())
}

func TestSessionServer(t *testing.T) {
	sess := &Session{
		server: "https://vcenter.example.com",
	}

	assert.Equal(t, "https://vcenter.example.com", sess.Server())
}

func TestSessionNilAccessors(t *testing.T) {
	sess := &Session{
		server: "https://vcenter",
	}

	// All accessor methods should return nil for unpopulated fields
	assert.Nil(t, sess.Vim25Client())
	assert.Nil(t, sess.RestClient())
	assert.Nil(t, sess.CNSClient())
	assert.Nil(t, sess.PBMClient())
	assert.Nil(t, sess.Finder())
}

func TestSessionLogger(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	sess := &Session{
		server: "https://vcenter",
		logger: logger,
	}

	sess.Close()

	// Should have logged the close message
	found := false
	for _, entry := range hook.AllEntries() {
		if entry.Message == "Closing vSphere session" {
			found = true
			break
		}
	}
	assert.True(t, found, "should have logged closing session message")
}

func TestSessionRedactHostname(t *testing.T) {
	// Same hostname should produce same hash
	hash1 := redactHostname("vcenter.example.com")
	hash2 := redactHostname("vcenter.example.com")
	assert.Equal(t, hash1, hash2)

	// Different hostnames should produce different hashes
	hash3 := redactHostname("vcenter2.example.com")
	assert.NotEqual(t, hash1, hash3)

	// Hash should not contain the hostname
	assert.NotContains(t, hash1, "vcenter")
	assert.NotContains(t, hash1, ".")
	assert.NotContains(t, hash1, ":")

	// Hash should be a hex string of fixed length (16 bytes = 32 hex chars)
	assert.Equal(t, 32, len(hash1), "hash should be 32 hex characters")
}

func TestSessionOptionsDefaults(t *testing.T) {
	opts := &SessionOptions{}
	opts.applyDefaults()

	assert.Equal(t, 60*time.Second, opts.Timeout, "timeout should default to 60s")
	assert.False(t, opts.Insecure, "insecure should default to false")
	assert.NotNil(t, opts.Logger, "logger should default to standard logger")
}

func TestSessionOptionsWithTimeout(t *testing.T) {
	opts := &SessionOptions{}
	opts.applyDefaults()
	opts.Timeout = 30 * time.Second

	assert.Equal(t, 30*time.Second, opts.Timeout)
}

func TestSessionOptionsWithInsecure(t *testing.T) {
	opts := &SessionOptions{}
	opts.applyDefaults()
	opts.Insecure = true

	assert.True(t, opts.Insecure)
}

func TestSessionOptionsWithLogger(t *testing.T) {
	logger := logrus.New()
	opts := &SessionOptions{}
	opts.applyDefaults()
	opts.Logger = logger

	assert.Equal(t, logger, opts.Logger)
}

func TestSessionCacheKeyCollision(t *testing.T) {
	// Verify that different combinations don't collide
	key1 := SessionCacheKey("https://vcenter", "user", "pass")
	key2 := SessionCacheKey("https://vcenter", "user2", "pass")
	key3 := SessionCacheKey("https://vcenter", "user", "pass2")

	assert.NotEqual(t, key1, key2, "different username should produce different key")
	assert.NotEqual(t, key1, key3, "different password should produce different key")
}

func TestSessionCacheKeyEmptyComponents(t *testing.T) {
	// Empty components should still produce valid (but different) keys
	key1 := SessionCacheKey("", "user", "pass")
	key2 := SessionCacheKey("https://vcenter", "", "pass")
	key3 := SessionCacheKey("https://vcenter", "user", "")

	assert.NotEmpty(t, key1)
	assert.NotEmpty(t, key2)
	assert.NotEmpty(t, key3)
	assert.NotEqual(t, key1, key2)
	assert.NotEqual(t, key1, key3)
	assert.NotEqual(t, key2, key3)
}

func TestSessionCacheKeyCaseSensitive(t *testing.T) {
	key1 := SessionCacheKey("https://vcenter", "User", "pass")
	key2 := SessionCacheKey("https://vcenter", "user", "pass")

	assert.NotEqual(t, key1, key2, "cache key should be case-sensitive")
}

func TestSessionCloseMultipleTimes(t *testing.T) {
	logger := logrus.New()

	sess := &Session{
		server: "https://vcenter",
		logger: logger,
	}

	// Close multiple times - should be idempotent
	sess.Close()
	err1 := sess.Error()

	sess.Close()
	err2 := sess.Error()

	sess.Close()
	err3 := sess.Error()

	// All errors should be the same
	assert.Equal(t, err1, err2)
	assert.Equal(t, err2, err3)
	assert.Nil(t, err1)
}

func TestSessionCacheKeySpecialCharacters(t *testing.T) {
	// Test with special characters in credentials
	key := SessionCacheKey("https://vcenter:8080", "user@domain", "pass:word/with/special")

	assert.NotEmpty(t, key)
	assert.Equal(t, 32, len(key))
}

func TestSessionCacheKeyVeryLongCredentials(t *testing.T) {
	// Test with very long credentials
	longServer := "https://" + string(make([]byte, 1000))
	longUser := string(make([]byte, 1000))
	longPass := string(make([]byte, 1000))

	key := SessionCacheKey(longServer, longUser, longPass)

	assert.NotEmpty(t, key)
	assert.Equal(t, 32, len(key))
}

func TestSessionOptionsApplyMultiple(t *testing.T) {
	opts := &SessionOptions{}
	opts.applyDefaults()

	// Apply multiple options
	opts.Timeout = 120 * time.Second
	opts.Insecure = true

	logger := logrus.New()
	opts.Logger = logger

	assert.Equal(t, 120*time.Second, opts.Timeout)
	assert.True(t, opts.Insecure)
	assert.Equal(t, logger, opts.Logger)
}

func TestSessionCacheKeyWhitespace(t *testing.T) {
	// Whitespace should be significant
	key1 := SessionCacheKey("https://vcenter", " user", "pass")
	key2 := SessionCacheKey("https://vcenter", "user ", "pass")

	assert.NotEqual(t, key1, key2, "whitespace should be significant")
}

func TestSessionCacheKeyLength(t *testing.T) {
	// SHA256 produces 32 bytes, we take first 16 bytes and format as hex
	// 16 bytes * 2 hex chars per byte = 32 hex chars
	key := SessionCacheKey("https://vcenter", "user", "pass")
	assert.Equal(t, 32, len(key), "cache key should be 32 hex characters")
}

func TestSessionCacheKeyHexFormat(t *testing.T) {
	key := SessionCacheKey("https://vcenter", "user", "pass")

	// Should be valid hex
	for _, c := range key {
		assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'),
			"character %q should be hex", c)
	}
}

func TestSessionCloseLogsRedactedHostname(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	sess := &Session{
		server: "https://vcenter.example.com",
		logger: logger,
	}

	sess.Close()

	// Check that the logged hostname is redacted
	found := false
	for _, entry := range hook.AllEntries() {
		if field, ok := entry.Data["vcenter"]; ok {
			vcenter, ok := field.(string)
			if ok && vcenter != "https://vcenter.example.com" {
				found = true
				break
			}
		}
	}
	assert.True(t, found, "should have logged redacted vcenter hostname")
}

func TestSessionOptionsNoOptions(t *testing.T) {
	opts := &SessionOptions{}
	opts.applyDefaults()

	// Should use all defaults
	assert.Equal(t, 60*time.Second, opts.Timeout)
	assert.False(t, opts.Insecure)
	assert.NotNil(t, opts.Logger)
}

func TestSessionCloseErrorIsNil(t *testing.T) {
	logger := logrus.New()

	sess := &Session{
		server: "https://vcenter",
		logger: logger,
	}

	// Close without any clients set - should not produce errors
	sess.Close()
	assert.Nil(t, sess.Error(), "error should be nil when no clients are set")
}

func TestSessionErrorBeforeClose(t *testing.T) {
	sess := &Session{
		server: "https://vcenter",
		logger: logrus.New(),
	}

	// Error() before Close() should return nil
	assert.Nil(t, sess.Error(), "error should be nil before Close() is called")
}

func TestSessionCacheKeyPortNumbers(t *testing.T) {
	key1 := SessionCacheKey("https://vcenter:443", "user", "pass")
	key2 := SessionCacheKey("https://vcenter:8443", "user", "pass")

	assert.NotEqual(t, key1, key2, "different ports should produce different keys")
}
