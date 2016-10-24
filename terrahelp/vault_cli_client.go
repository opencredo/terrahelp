package terrahelp

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type VaultCliClient struct {
}

type vaultOutput struct {
	Data map[string]string `json:"data"`
}

var _ VaultClient = &VaultCliClient{}

func NewVaultCliClient() (*VaultCliClient, error) {

	return &VaultCliClient{}, nil
}

// MountTransitBackend ensures the transit backend is mounted
func (v *VaultCliClient) MountTransitBackend() error {
	exists, err := v.transitMountExists()
	if err != nil {
		return err
	}

	if !exists {
		log.Println("Mounting transit backend ... ")
		_, err := exec.Command("vault", "mount", "-path=transit", "transit").Output()
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
func (v *VaultCliClient) RegisterNamedEncryptionKey(key string) error {
	exists, err := v.namedEncryptionKeyExists(key)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("Named encryption key '%s' does not exist, creating at %s ... ", key, v.encryptKeyPath(key))

		_, e := exec.Command("vault", "write", v.encryptKeyPath(key)).Output()
		return e
	}

	log.Printf("Named encryption key '%s' already exists at %s ... ", key, v.encryptKeyPath(key))
	return nil
}

func (v *VaultCliClient) transitMountExists() (bool, error) {

	out, err := exec.Command("vault", "mounts").Output()

	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if fields[0] == "transit/" {
			return true, nil
		}
	}

	return false, nil
}

func (v *VaultCliClient) namedEncryptionKeyExists(key string) (bool, error) {
	_, err := exec.Command("vault", "read", key).Output()

	switch err.(type) {
	case nil:
		return true, nil
	case (*exec.ExitError):
		var exitError *exec.ExitError = err.(*exec.ExitError)
		if strings.HasPrefix(string(exitError.Stderr), "No value found at") {
			return false, nil
		}
		return false, err
	default:
		return false, err
	}

}

// Decrypt uses the named encryption key to decrypt the supplied content
func (v *VaultCliClient) Decrypt(key, ciphertext string) (string, error) {
	return v.transit(v.decryptEndpoint(key), "ciphertext", ciphertext, "plaintext")
}

// Encrypt uses the named encryption key to encrypt the supplied content
func (v *VaultCliClient) Encrypt(key, b64text string) (string, error) {
	return v.transit(v.encryptEndpoint(key), "plaintext", b64text, "ciphertext")
}

func (v *VaultCliClient) transit(key, inputField, value, expectedField string) (string, error) {
	cmd := exec.Command("vault", "write", "-format=json", key, inputField+"=-")
	cmd.Stdin = strings.NewReader(value)
	out, err := cmd.Output()

	if err != nil {
		log.Printf("error occurred %s", cmd.Args)
		return "", err
	}

	output := &vaultOutput{}

	json.Unmarshal(out, output)

	return output.Data[expectedField], nil
}

func (v *VaultCliClient) encryptKeyPath(key string) string {
	return "/transit/keys/" + key
}

func (v *VaultCliClient) encryptEndpoint(key string) string {
	return "/transit/encrypt/" + key
}

func (v *VaultCliClient) decryptEndpoint(key string) string {
	return "/transit/decrypt/" + key
}
