package terrahelp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pairs = map[string]string{
	"sensitive-value-1-AK#%DJGHS*G": "vault:v1:c2Vuc2l0aXZlLXZhbHVlLTEtQUsjJURKR0hTKkc=",
	"sensitive-value-3-//dfhs//":    "vault:v1:c2Vuc2l0aXZlLXZhbHVlLTMtLy9kZmhzLy8=",
}

func TestMockVaultClient_doSimEncryption(t *testing.T) {
	m := NewMockVaultClient()

	for k, v := range pairs {
		enc := m.doSimEncryption(k)
		assert.Equal(t, v, enc)
	}

}

func TestMockVaultClient_doSimDecryption(t *testing.T) {
	m := NewMockVaultClient()

	for k, v := range pairs {
		dec, err := m.doSimDecryption(v)
		assert.NoError(t, err)
		assert.Equal(t, k, dec)
	}

}
