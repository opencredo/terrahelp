package terrahelp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"fmt"
	"io"

	"bytes"
	"errors"
	"regexp"
	"strings"
)

const (
	thCryptoWrapRegExp = "\\@terrahelp\\-encrypted\\((.*?)\\)"
	//thCryptoWrapRegExp = "\\@terrahelp\\-encrypted\\(.*?\\)"
	thCryptoWrapPrefix = "@terrahelp-encrypted("
	thCryptoWrapSuffix = ")"

	thCryptoWrapInvalidMsg = "Unable to decrypt ciphertext, not wrapped as expected"
)

// Encrypter defines the functionality required to be supported
// by crypto backends which are to be used for encrypting and
// decrypting tfstate files
type Encrypter interface {
	Init(key string) error
	Decrypt(key string, b []byte) ([]byte, error)
	Encrypt(key string, b []byte) ([]byte, error)
}

func applyTHCryptoWrap(b []byte) []byte {
	return bytes.Join([][]byte{
		[]byte(thCryptoWrapPrefix),
		b,
		[]byte(thCryptoWrapSuffix)}, []byte{})
}

func extractFromTHCryptoWrap(b []byte) (string, error) {
	ciphertxt := string(b)
	if ciphertxt != "" && !strings.HasPrefix(ciphertxt, thCryptoWrapPrefix) {
		return "", fmt.Errorf(thCryptoWrapInvalidMsg)
	}

	r := regexp.MustCompile(thCryptoWrapRegExp)
	m := r.FindStringSubmatch(ciphertxt)
	if m == nil || len(m) < 1 || m[1] == "" {
		return "", fmt.Errorf(thCryptoWrapInvalidMsg)
	}

	return m[1], nil
}

// ---------------------------------------------------------------
//                       VaultEncrypter
// ---------------------------------------------------------------

// VaultEncrypter wraps the real core Vault client exposing convenient methods
// required to interact with Vault in order to perform encrypting and
// decrypting of tfstate files
type VaultEncrypter struct {
	vault VaultClient
}

// NewVaultEncrypter creates a new VaultEncrypter
func NewVaultEncrypter() (*VaultEncrypter, error) {
	var vc VaultClient
	vc, err := NewDefaultVaultClient()
	if err != nil {
		return nil, err
	}
	return createVaultEncrypter(vc)
}

func createVaultEncrypter(vc VaultClient) (*VaultEncrypter, error) {
	return &VaultEncrypter{vc}, nil
}

// Init is used to initialise Vault for the purposes of using
// its encryption as a service functionality
func (cu *VaultEncrypter) Init(key string) error {
	err := cu.vault.MountTransitBackend()
	if err != nil {
		return err
	}

	err = cu.vault.RegisterNamedEncryptionKey(key)
	if err != nil {
		return err
	}
	return nil
}

// Encrypt uses the named encryption key to encrypt the
// provided plaintext
func (cu *VaultEncrypter) Encrypt(key string, plaintext []byte) ([]byte, error) {
	enc, err := cu.vault.Encrypt(key, base64.StdEncoding.EncodeToString(plaintext))
	if err != nil {
		return nil, err
	}
	return applyTHCryptoWrap([]byte(enc)), nil
}

// Decrypt uses the named encryption key to decrypt the
// provided ciphertext
func (cu *VaultEncrypter) Decrypt(key string, ciphertext []byte) ([]byte, error) {

	actualContent, err := extractFromTHCryptoWrap(ciphertext)
	if err != nil {
		return nil, err
	}

	pt, err := cu.vault.Decrypt(key, actualContent)
	if err != nil {
		return nil, err
	}
	ptb, err := base64.StdEncoding.DecodeString(pt)
	return ptb, err
}

// ---------------------------------------------------------------
//                       SimpleEncrypter
// ---------------------------------------------------------------

type SimpleEncrypter struct{}

func NewSimpleEncrypter() *SimpleEncrypter {
	return &SimpleEncrypter{}
}

// Init is used to initialise Vault for the purposes of using
// its encryption as a service functionality
func (s *SimpleEncrypter) Init(key string) error {
	return nil
}

// The key argument should be the AES key, either 16 or 32 bytes
// to select AES-128 or AES-256.
// "AES256Key-32Characters0987654321"
func (s *SimpleEncrypter) Encrypt(skey string, text []byte) ([]byte, error) {

	key := []byte(skey)

	var block cipher.Block

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))

	// iv =  initialization vector
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	sc := base64.StdEncoding.EncodeToString(ciphertext)

	return applyTHCryptoWrap([]byte(sc)), nil
}

func (s *SimpleEncrypter) Decrypt(skey string, stext []byte) ([]byte, error) {

	actualContent, err := extractFromTHCryptoWrap(stext)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(actualContent)
	if err != nil {
		return nil, err
	}

	key := []byte(skey)
	// should be plaintext by now

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
