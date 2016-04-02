package terrahelp

import (
	"fmt"

	"github.com/hashicorp/vault/api"

	"errors"
	"log"
	"os"
)

// VaultClient defines the basic functionality required by terrahelp
// when interacting with Vault
type VaultClient interface {

	// RegisterNamedEncryptionKey registers the named encryption key
	// within Vault's transit backend
	RegisterNamedEncryptionKey(key string) error
	// MountTransitBackend ensures the transit backend is mounted
	MountTransitBackend() error
	// Encrypt uses the named encryption key to encrypt the supplied content
	Encrypt(key, text string) (string, error)
	// Decrypt uses the named encryption key to decrypt the supplied content
	Decrypt(key, ciphertext string) (string, error)

	transitMountExists() (bool, error)
	namedEncryptionKeyExists(key string) (bool, error)
}

// DefaultVaultClient provides a wrapper around the core Vault
// client and uses it to provide the required functionality
type DefaultVaultClient struct {
	*api.Client
}

// NewDefaultVaultClient creates a new DefaultVaultClient
func NewDefaultVaultClient() (*DefaultVaultClient, error) {

	if os.Getenv("VAULT_TOKEN") == "" || os.Getenv("VAULT_ADDR") == "" {
		return nil, errors.New("\n  This CLI relies on the standard Vault environment variables (VAULT_TOKEN, VAULT_ADDR etc)" +
			"\n  for obtaining the configuration and authentication details required to connect to the Vault server" +
			"\n  please configure these before continuing.")
	}
	vc := api.DefaultConfig()
	vc.ReadEnvironment()
	vclient, err := api.NewClient(vc)
	if err != nil {
		return nil, fmt.Errorf("Issue getting client : %s", err)
	}

	return &DefaultVaultClient{vclient}, nil
}

// MountTransitBackend ensures the transit backend is mounted
func (v *DefaultVaultClient) MountTransitBackend() error {
	exists, err := v.transitMountExists()
	if err != nil {
		return err
	}

	if !exists {
		log.Println("Mounting transit backend ... ")
		err := v.Sys().Mount("transit", &api.MountInput{
			Type:   "transit",
			Config: api.MountConfigInput{},
		})
		if err != nil {
			return err
		}
	} else {
		log.Println("transit backend already exists ... ")
	}
	return nil
}

// RegisterNamedEncryptionKey registers the named encryption key
// within Vault's transit backend
func (v *DefaultVaultClient) RegisterNamedEncryptionKey(key string) error {
	exists, err := v.namedEncryptionKeyExists(key)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("Named encryption key '%s' does not exist, creating at %s ... ", key, v.encryptKeyPath(key))
		_, e := v.Logical().Write(v.encryptKeyPath(key), map[string]interface{}{})
		return e
	}

	log.Printf("Named encryption key '%s' already exists at %s ... ", key, v.encryptKeyPath(key))
	return nil
}

func (v *DefaultVaultClient) transitMountExists() (bool, error) {
	mp, err := v.Sys().ListMounts()
	if err != nil {
		return false, err
	}
	for key := range mp {
		if key == "transit/" {
			return true, nil
		}
	}
	return false, nil
}

func (v *DefaultVaultClient) namedEncryptionKeyExists(key string) (bool, error) {
	s, err := v.Logical().Read(v.encryptKeyPath(key))
	if err != nil {
		return false, err
	}
	if s == nil {
		return false, nil
	}
	return true, nil

}

// Decrypt uses the named encryption key to decrypt the supplied content
func (v *DefaultVaultClient) Decrypt(key, ciphertext string) (string, error) {
	kv := map[string]interface{}{"ciphertext": ciphertext}
	s, err := v.Logical().Write(v.decryptEndpoint(key), kv)
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", fmt.Errorf("Unable to get decryped value using encryption key %s ", key)
	}
	return s.Data["plaintext"].(string), nil
}

// Encrypt uses the named encryption key to encrypt the supplied content
func (v *DefaultVaultClient) Encrypt(key, b64text string) (string, error) {
	kv := map[string]interface{}{"plaintext": b64text}
	s, err := v.Logical().Write(v.encryptEndpoint(key), kv)
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", fmt.Errorf("Unable to get encryption value using encryption key %s ", key)
	}
	return s.Data["ciphertext"].(string), nil
}

func (v *DefaultVaultClient) encryptKeyPath(key string) string {
	return "/transit/keys/" + key
}

func (v *DefaultVaultClient) encryptEndpoint(key string) string {
	return "/transit/encrypt/" + key
}

func (v *DefaultVaultClient) decryptEndpoint(key string) string {
	return "/transit/decrypt/" + key
}
