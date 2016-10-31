package terrahelp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTfvars_ExtractSensitiveStringVals(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	expected := []string{
		"madeup-aws-access-key-PEJFNS",
		"madeup-aws-secret-key-KGSDGH",
		"sensitive-value-1-AK#%DJGHS*G",
		"sensitive-value-2-prYh57",
		"sensitive-value-3-//dfhs//",
		"sensitive-value-6"}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}

}

func TestTfvars_ExtractSensitiveStringValWithEqualSign(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	expected := []string{
		"sensitive-value-4 with equals sign i.e. ff=yy"}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}

}

func TestTfvars_ExtractSensitiveListVals(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	expected := []string{
		"sensitive-list-val",
		"sensitive-list-val-1",
		"sensitive-list-val-2"}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}

}

func TestTfvars_ExtractSensitiveFlatMapVals(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	expected := []string{
		"sensitive-flatmap-val-foo",
		"sensitive-flatmap-val-bax",
		"sensitive-flatmap-val"}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}

}

func TestTfvars_ExtractSensitiveFlatMapVals_ExcludesKeyName(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	assert.NotContains(t, actual, "bob")
	assert.NotContains(t, actual, "overlap")

}

func TestTfvars_ExtractSensitiveValsWithWhitespace_WithConfigExclWhitespaceTrue(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	expected := []string{
		" sensitive-with-leading-space",
		"sensitive-with-trail-space ",
		" sensitive-with-leadtrail-space "}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}
}

func TestTfvars_ExtractSensitiveValsWithWhitespace_WithConfigExclWhitespaceFalse(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", false)
	expected := []string{
		" sensitive-with-leading-space",
		"sensitive-with-trail-space ",
		" sensitive-with-leadtrail-space "}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}
}

func TestTfvars_ExtractSensitiveValsWithEmptyString_WithConfigExclWhitespaceTrue(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	assert.NotContains(t, actual, "")
}

func TestTfvars_ExtractSensitiveValsWithEmptyString_WithConfigExclWhitespaceFalse(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", false)
	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	assert.NotContains(t, actual, "")
}

func TestTfvars_ExtractSensitiveValsWithOnlyWhitespace_WithConfigExclWhitespaceTrue(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", true)
	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	assert.NotContains(t, actual, "   ")

}

func TestTfvars_ExtractSensitiveValsWithOnlyWhitespace_WithConfigExclWhitespaceFalse(t *testing.T) {
	// Given
	tu := NewTfVars("test-data/example-project/original/terraform.tfvars", false)
	expected := []string{"   "}

	// When
	actual, err := tu.Values()

	// Then
	assert.NoError(t, err)
	for _, v := range expected {
		assert.Contains(t, actual, v)
	}

}
