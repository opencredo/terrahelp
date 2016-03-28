package terrahelp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------
//                       VaultEncrypter
// ---------------------------------------------------------------

func TestVaultEncrypter_Init(t *testing.T) {
	vcu := getTestVaultEncrypter(t, "")
	err := vcu.Init("testkey")
	assert.NoError(t, err)
}

func TestVaultEncrypter_EncryptDecrypt(t *testing.T) {
	// Given
	orig := []byte("sample content")
	vcu := getTestVaultEncrypter(t, "testkey")

	// When
	e, err := vcu.Encrypt("testkey", orig)
	assert.NoError(t, err)

	d, err := vcu.Decrypt("testkey", e)
	assert.NoError(t, err)
	assert.Equal(t, orig, d)

}

func TestVaultEncrypter_Encrypt_IsWrappedCorrectly(t *testing.T) {
	// Given
	vcu := getTestVaultEncrypter(t, "testkey")
	c := []byte("sample content")

	// When
	enc, err := vcu.Encrypt("testkey", c)
	assert.NoError(t, err)
	assertHasTerraHelpCryptoWrapper(t, string(enc))
}

func TestVaultEncrypter_Decrypt_InvalidWrapper(t *testing.T) {
	// Given
	vcu := getTestVaultEncrypter(t, "testkey")
	invalids := []string{
		"@notwrappedproperly(SOMETHING)",
		"@terrahelp-encrypted(vault:v1:SOMETHING",
		"@terrahelp-encrypted()",
		"asf asf @terrahelp-encrypted[vault:v1:SOMETHING)",
	}

	// When
	for _, i := range invalids {
		enc, err := vcu.Decrypt("testkey", []byte(i))
		assert.Error(t, err)
		assert.Empty(t, enc)
		assert.Equal(t, thCryptoWrapInvalidMsg, err.Error(),
			fmt.Sprintf("%s not detected as invalid wrapper", i))
	}
}

func TestVaultEncrypter_Decrypt_NotBase64(t *testing.T) {
	// Given
	vcu := getTestVaultEncrypter(t, "testkey")
	invalids := []string{"@terrahelp-encrypted(vault:v1:NOTBASE64)"}

	// When
	for _, i := range invalids {
		enc, err := vcu.Decrypt("testkey", []byte(i))
		assert.Error(t, err)
		assert.Empty(t, enc)
		assert.True(t, strings.HasPrefix(err.Error(), "illegal base64 data"),
			fmt.Sprintf("%s not detected as having invalid base64 encoding", i))
	}
}

// -------------------------------------------------------------
//                   Test helper methods
// -------------------------------------------------------------
func getTestVaultEncrypter(t *testing.T, key string) *VaultEncrypter {
	cu, err := createVaultEncrypter(NewMockVaultClient())
	if err != nil {
		t.Fatalf("Error trying to creat VaultEncrypter %s", err)
	}
	if key != "" {
		err := cu.Init(key)
		if err != nil {
			t.Fatalf("Error trying to register key for test %s", err)
		}
	}
	return cu
}

func getTestSimpleEncrypter(t *testing.T) *SimpleEncrypter {
	return NewSimpleEncrypter()
}

func assertHasTerraHelpCryptoWrapper(t *testing.T, enc string) {
	assert.True(t, strings.HasPrefix(enc, thCryptoWrapPrefix),
		fmt.Sprintf("Encrypted value %s does not start with expected prefix", enc))
	assert.True(t, strings.HasSuffix(enc, thCryptoWrapSuffix),
		fmt.Sprintf("Encrypted value %s does not end with expected suffix", enc))
}

// ---------------------------------------------------------------
//                       SimpleEncrypter
// ---------------------------------------------------------------

func TestSimpleEncrypter_Init(t *testing.T) {
	vcu := SimpleEncrypter{}
	err := vcu.Init("testkey")
	assert.NoError(t, err)
}

func TestSimpleEncrypter_EncryptDecrypt(t *testing.T) {
	// Given
	orig := []byte("sample content")
	encKey := "AES256Key-32Characters0987654321"
	vcu := getTestSimpleEncrypter(t)

	// When
	e, err := vcu.Encrypt(encKey, orig)
	assert.NoError(t, err)

	d, err := vcu.Decrypt(encKey, e)
	assert.NoError(t, err)
	assert.Equal(t, orig, d)

}

func TestSimpleEncrypter_Encrypt_IsWrappedCorrectly(t *testing.T) {
	// Given
	encKey := "AES256Key-32Characters0987654321"
	vcu := getTestSimpleEncrypter(t)
	c := []byte("sample content")

	// When
	enc, err := vcu.Encrypt(encKey, c)
	assert.NoError(t, err)
	assertHasTerraHelpCryptoWrapper(t, string(enc))
}

func TestSimpleEncrypter_Decrypt_InvalidWrapper(t *testing.T) {
	// Given
	encKey := "AES256Key-32Characters0987654321"
	vcu := getTestSimpleEncrypter(t)
	invalids := []string{
		"@notwrappedproperly(SOMETHING)",
		"@terrahelp-encrypted(vault:v1:SOMETHING",
		"@terrahelp-encrypted()",
		"asf asf @terrahelp-encrypted[vault:v1:SOMETHING)",
	}

	// When
	for _, i := range invalids {
		enc, err := vcu.Decrypt(encKey, []byte(i))
		assert.Error(t, err)
		assert.Empty(t, enc)
		assert.Equal(t, thCryptoWrapInvalidMsg, err.Error(),
			fmt.Sprintf("%s not detected as invalid wrapper", i))
	}
}
