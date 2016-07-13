package terrahelp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"io"

	"bytes"
	"errors"
	"regexp"
	"strings"
)

const (
	thCryptoWrapRegExp = "\\@terrahelp\\-encrypted\\((.*?)\\)"
	thCryptoWrapPrefix = "@terrahelp-encrypted("
	thCryptoWrapSuffix = ")"

	thCryptoWrapInvalidMsg = "Unable to decrypt ciphertext, not wrapped as expected"
)

// A CryptoWrapError describes an error where a missing or invalid use of the
// terrahelp wrapper value (i.e. @terrahelp-encrypted() ) prevents the encryption
// or decryption being performed
type CryptoWrapError struct {
	msg string
}

func (e *CryptoWrapError) Error() string {
	return "terrahelp encryption error : " + e.msg
}

func newCryptoWrapError(msg string) *CryptoWrapError {
	return &CryptoWrapError{msg: msg}
}

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
		return "", newCryptoWrapError(thCryptoWrapInvalidMsg)
	}

	r := regexp.MustCompile(thCryptoWrapRegExp)
	m := r.FindStringSubmatch(ciphertxt)
	if m == nil || len(m) < 1 || m[1] == "" {
		return "", newCryptoWrapError(thCryptoWrapInvalidMsg)
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

func NewVaultCliEncrypter() (*VaultEncrypter, error) {
	var vc VaultClient
	vc, err := NewVaultCliClient()
	if err != nil {
		return nil, err
	}
	return &VaultEncrypter{vc}, nil
}

func createVaultEncrypter(vc VaultClient) (*VaultEncrypter, error) {
	return &VaultEncrypter{vc}, nil
}

// Init is used to initialise the VaultEncrypter for the purposes of using
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

// SimpleEncrypter provides basic AES based encryption
type SimpleEncrypter struct{}

// NewSimpleEncrypter creates a new SimpleEncrypter with
// default configuration
func NewSimpleEncrypter() *SimpleEncrypter {
	return &SimpleEncrypter{}
}

// Init is used to initialise Vault for the purposes of using
// its encryption as a service functionality
func (s *SimpleEncrypter) Init(key string) error {
	return nil
}

// Encrypt will perform AES based encryption on the byte content provided.
// The key should be an AES key, of either either 16 or 32 characters
// which then informs whether AES-128 or AES-256 encryption is applied.
func (s *SimpleEncrypter) Encrypt(key string, b []byte) ([]byte, error) {

	bkey := []byte(key)

	var block cipher.Block

	block, err := aes.NewCipher(bkey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(b))

	// iv =  initialization vector
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	sc := base64.StdEncoding.EncodeToString(ciphertext)

	return applyTHCryptoWrap([]byte(sc)), nil
}

// Decrypt will use the supplied AES key to decrypt the byte content provided.
func (s *SimpleEncrypter) Decrypt(key string, b []byte) ([]byte, error) {

	actualContent, err := extractFromTHCryptoWrap(b)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(actualContent)
	if err != nil {
		return nil, err
	}

	bkey := []byte(key)
	// should be plaintext by now

	block, err := aes.NewCipher(bkey)

	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("decryption failed: ciphertext is too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
