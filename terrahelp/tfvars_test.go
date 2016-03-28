package terrahelp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTfvars_ExtractSensitiveVals(t *testing.T) {
	// Given
	tu := &Tfvars{}

	// When
	actual, err := tu.ExtractSensitiveVals("test-data/example-project/original/terraform.tfvars")
	expected := []string{
		"sensitive-value-1-AK#%DJGHS*G",
		"sensitive-value-2-prYh57",
		"sensitive-value-3-//dfhs//",
		"sensitive-value-4 with equals sign i.e. ff=yy",
		"sensitive-value-6"}

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}

}
