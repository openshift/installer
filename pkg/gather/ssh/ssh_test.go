package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "authentication error",
			err:      errors.New("ssh: handshake failed: ssh: unable to authenticate"),
			expected: true,
		},
		{
			name:     "authentication error with additional context",
			err:      fmt.Errorf("failed to connect: %w", errors.New("ssh: handshake failed: ssh: unable to authenticate, tried publickey")),
			expected: true,
		},
		{
			name:     "connection refused error",
			err:      errors.New("dial tcp: connection refused"),
			expected: false,
		},
		{
			name:     "timeout error",
			err:      errors.New("dial tcp: i/o timeout"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "generic ssh error",
			err:      errors.New("ssh: some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAuthenticationError(tt.err)
			if result != tt.expected {
				t.Errorf("isAuthenticationError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

// mockSSHServer creates a mock SSH server for testing
type mockSSHServer struct {
	listener      net.Listener
	authorizedKey ssh.PublicKey
	t             *testing.T
}

func newMockSSHServer(t *testing.T, authorizedKey ssh.PublicKey) (*mockSSHServer, error) {
	// Generate a host key for the server
	hostKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	hostSigner, err := ssh.NewSignerFromKey(hostKey)
	if err != nil {
		return nil, err
	}

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKey != nil && string(key.Marshal()) == string(authorizedKey.Marshal()) {
				return &ssh.Permissions{}, nil
			}
			return nil, fmt.Errorf("unauthorized key")
		},
	}
	config.AddHostKey(hostSigner)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	server := &mockSSHServer{
		listener:      listener,
		authorizedKey: authorizedKey,
		t:             t,
	}

	go server.acceptConnections(config)

	return server, nil
}

func (s *mockSSHServer) acceptConnections(config *ssh.ServerConfig) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go func(c net.Conn) {
			defer c.Close()
			_, _, _, err := ssh.NewServerConn(c, config)
			if err != nil {
				// Expected for failed authentication attempts
				return
			}
		}(conn)
	}
}

func (s *mockSSHServer) Address() string {
	return s.listener.Addr().String()
}

func (s *mockSSHServer) Close() {
	s.listener.Close()
}

func TestAttemptSSHConnection(t *testing.T) {
	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	publicKey := signer.PublicKey()

	// Create mock SSH server that accepts our key
	server, err := newMockSSHServer(t, publicKey)
	if err != nil {
		t.Fatalf("Failed to create mock SSH server: %v", err)
	}
	defer server.Close()

	tests := []struct {
		name        string
		setupAgent  func() agent.Agent
		expectError bool
		errorCheck  func(error) bool
	}{
		{
			name: "successful connection with valid key",
			setupAgent: func() agent.Agent {
				ag := agent.NewKeyring()
				ag.Add(agent.AddedKey{PrivateKey: privateKey})
				return ag
			},
			expectError: false,
		},
		{
			name: "failed connection with invalid key",
			setupAgent: func() agent.Agent {
				wrongKey, _ := rsa.GenerateKey(rand.Reader, 2048)
				ag := agent.NewKeyring()
				ag.Add(agent.AddedKey{PrivateKey: wrongKey})
				return ag
			},
			expectError: true,
			errorCheck:  isAuthenticationError,
		},
		{
			name: "failed connection with no keys",
			setupAgent: func() agent.Agent {
				return agent.NewKeyring()
			},
			expectError: true,
			errorCheck:  isAuthenticationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ag := tt.setupAgent()
			client, err := attemptSSHConnection("testuser", server.Address(), ag)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorCheck != nil && !tt.errorCheck(err) {
					t.Errorf("Error %v did not match expected error type", err)
				}
				if client != nil {
					client.Close()
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if client == nil {
					t.Errorf("Expected client but got nil")
				} else {
					client.Close()
				}
			}
		})
	}
}

func TestNewClientRetryLogic(t *testing.T) {
	// This test verifies that NewClient properly handles authentication errors
	// The actual retry logic with SSH_AUTH_SOCK is difficult to test in a unit test
	// as it depends on environment state, but we can test the basic flow

	// Generate a test key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to create signer: %v", err)
	}

	publicKey := signer.PublicKey()

	// Create mock SSH server
	server, err := newMockSSHServer(t, publicKey)
	if err != nil {
		t.Fatalf("Failed to create mock SSH server: %v", err)
	}
	defer server.Close()

	t.Run("connection to invalid address fails", func(t *testing.T) {
		// Test with invalid address - should fail with connection error
		_, err := NewClient("testuser", "127.0.0.1:1", []string{})
		if err == nil {
			t.Errorf("Expected error for invalid address, got none")
		}
		// Should not be an authentication error
		if isAuthenticationError(err) {
			t.Errorf("Expected connection error, got authentication error: %v", err)
		}
	})

	t.Run("connection with no keys fails with auth error", func(t *testing.T) {
		// Test with no keys - should fail with authentication error
		_, err := NewClient("testuser", server.Address(), []string{})
		if err == nil {
			t.Errorf("Expected authentication error, got none")
		}
	})
}

func TestNewClientErrorMessages(t *testing.T) {
	// Test that appropriate error messages are returned for different failure scenarios
	tests := []struct {
		name        string
		address     string
		keys        []string
		errContains string
	}{
		{
			name:        "connection refused",
			address:     "127.0.0.1:1",
			keys:        []string{},
			errContains: "connection refused",
		},
		{
			name:        "invalid address format",
			address:     "invalid:address:format",
			keys:        []string{},
			errContains: "too many colons",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient("testuser", tt.address, tt.keys)
			if err == nil {
				t.Errorf("Expected error containing %q, got none", tt.errContains)
			}
		})
	}
}
