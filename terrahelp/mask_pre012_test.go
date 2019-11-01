package terrahelp

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func defaultTestMaskOpts_enablePre012(t *testing.T) (*MaskOpts, *stdinSim, *stdoutSim) {
	ctx := NewDefaultMaskOpts()
	ctx.EnablePre012 = true
	stdinSim := newStdinSim(t)
	stdinSim.start()
	stdoutSim := newStdoutSim(t)
	stdoutSim.start()
	ctx.TransformItems = []Transformable{
		NewStreamTransformable(stdinSim.simReadFile, stdoutSim.simWriteFile)}
	return ctx, stdinSim, stdoutSim
}

func TestMasker_Mask_enablePre012_StreamedNonSensitiveData(t *testing.T) {
	// Given some input content ...
	var data = `hello there
                    I am some data
                    to be piped in`

	ctx, stdinSim, stdoutSim := defaultTestMaskOpts_enablePre012(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, nil)

	// When we simulate piping this content (which contains
	// NO sensitive data) into stdIn
	stdinSim.write(data)
	err := m.Mask()

	// Then the output should be exactly the same as the original
	// content passed in
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t, data, b)
}

func TestMasker_Mask_enablePre012_StreamedSensitiveData(t *testing.T) {
	// Given some input content
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts_enablePre012(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// only the known sensitive data) into stdIn
	stdinSim.write(
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "sensitive-value-1-AK#%DJGHS*G"`)
	err := m.Mask()

	// Then the output should have that sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "******"`, b)
}

func TestMasker_Mask_enablePre012_StreamedNothing2KnownSensitiveDataTransform(t *testing.T) {
	// Given some input content
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts_enablePre012(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// a transition from nothing to the known sensitive data) into stdIn
	stdinSim.write(
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "\n blah blah"
    vars.#:    "3" => "3"
    vars.msg1: "" => "sensitive-value-1-AK#%DJGHS*G" (forces new resource)`)
	err := m.Mask()

	// Then the output should have that sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "\n blah blah"
    vars.#:    "3" => "3"
    vars.msg1: "" => "******" (forces new resource)`, b)
}

func TestMasker_Mask_enablePre012_StreamedPrevVal2KnownSensitiveDataTransform(t *testing.T) {
	// Given some input content ...
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts_enablePre012(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// a transition from a previous sensitive value to new
	// known sensitive data) into stdIn
	stdinSim.write(
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "\n blah blah"
    vars.#:    "3" => "3"
    vars.msg1: "some-prev-sensitive-value-1\"" => "sensitive-value-1-AK#%DJGHS*G" (forces new resource)`)
	err := m.Mask()

	// Then the output should have both the previous and
	// known sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "\n blah blah"
    vars.#:    "3" => "3"
    vars.msg1: "******" => "******" (forces new resource)`, b)
}

func TestMasker_Mask_enablePre012_StreamedPrevVal2KnownSensitiveDataTransform_IgnorePrev(t *testing.T) {
	// Given some input content and an explicit directive NOT
	// to attempt to detect previous sensitive values ...
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts_enablePre012(t)
	ctx.ReplacePrevVals = false
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// a transition from a previous sensitive value to a new
	// known sensitive data) into stdIn
	stdinSim.write(
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "\n blah blah"
    vars.#:    "3" => "3"
    vars.msg1: "some-prev-sensitive-value-1\"" => "sensitive-value-1-AK#%DJGHS*G" (forces new resource)`)
	err := m.Mask()

	// Then only the new known sensitive values should be masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
The Terraform execution plan has been generated and is shown below ...
-/+ template_file.example
    rendered:  "\n blah blah"
    vars.#:    "3" => "3"
    vars.msg1: "some-prev-sensitive-value-1\"" => "******" (forces new resource)`, b)
}
