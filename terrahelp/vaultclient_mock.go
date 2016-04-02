package terrahelp

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const vaultV1EncryptionPrefix = "vault:v1:"

// MockVaultClient provides a mock implementation of the
// VaultClient interface for testing purposes
type MockVaultClient struct {
	key            string
	transitMounted bool
}

// NewMockVaultClient creates a new MockVaultClient
func NewMockVaultClient() *MockVaultClient {
	return &MockVaultClient{}
}

// MountTransitBackend mocks the mounting of the transit backend
func (m *MockVaultClient) MountTransitBackend() error {
	m.transitMounted = true
	return nil
}

func (m *MockVaultClient) transitMountExists() (bool, error) {
	return m.transitMounted, nil
}

func (m *MockVaultClient) namedEncryptionKeyExists(key string) (bool, error) {
	return m.key != "", nil
}

// RegisterNamedEncryptionKey registers the named encryption key
// within the mock Vault service
func (m *MockVaultClient) RegisterNamedEncryptionKey(key string) error {
	m.transitMounted = true
	m.key = key
	return nil
}

// Encrypt uses the named encryption key to mock encrypt the supplied content
func (m *MockVaultClient) Encrypt(key, s string) (string, error) {
	if !m.transitMounted {
		return "", errors.New("Mock client has not had transit backend mounted")
	}
	if m.key != key {
		return "", fmt.Errorf("Unknown encryption key %s", key)
	}

	return m.doSimEncryption(s), nil
}

// Decrypt uses the named encryption key to mock decrypt the supplied content
func (m *MockVaultClient) Decrypt(key, s string) (string, error) {
	if !m.transitMounted {
		return "", errors.New("Mock client has not had transit backend mounted")
	}
	if m.key != key {
		return "", fmt.Errorf("Unknown encryption key %s", key)
	}
	return m.doSimDecryption(s)
}

// NB this is not proper encryption, it is ONLY for testing and is merely
//    base64 encoding the plaintext value, so not really encryption at all
//    but as its only use for testing its OK
func (m *MockVaultClient) doSimEncryption(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		b[i] = s[i]
	}

	return vaultV1EncryptionPrefix + base64.StdEncoding.EncodeToString(b)
}

// NB this is not proper encryption, it is ONLY for testing and is merely
//    base64 encoding the plaintext value, so not really encryption at all
//    but as its only use for testing its OK
func (m *MockVaultClient) doSimDecryption(s string) (string, error) {
	if !strings.HasPrefix(s, vaultV1EncryptionPrefix) {
		return "", fmt.Errorf("Unable to sim decrypt string %s, does not have expected starting prefix %s)", s, vaultV1EncryptionPrefix)
	}
	b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, vaultV1EncryptionPrefix))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
