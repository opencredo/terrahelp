package terrahelp

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func defaultTestMaskOpts(t *testing.T) (*MaskOpts, *stdinSim, *stdoutSim) {
	ctx := NewDefaultMaskOpts()
	stdinSim := newStdinSim(t)
	stdinSim.start()
	stdoutSim := newStdoutSim(t)
	stdoutSim.start()
	ctx.TransformItems = []Transformable{
		NewStreamTransformable(stdinSim.simReadFile, stdoutSim.simWriteFile)}
	return ctx, stdinSim, stdoutSim
}

func TestMasker_Mask_StreamedNonSensitiveData(t *testing.T) {
	// Given some input content ...
	var data = `hello there
                    I am some data
                    to be piped in`

	ctx, stdinSim, stdoutSim := defaultTestMaskOpts(t)
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

func TestMasker_Mask_StreamedSensitiveData(t *testing.T) {
	// Given some input content
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// only the known sensitive data) into stdIn
	stdinSim.write(
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
  + resource "template_dir" "config" {
      + destination_dir = "./renders"
      + id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
      + source_dir      = "./templates"
      + vars            = {
          "msg1" = "sensitive-value-1-AK#%DJGHS*G"
        }
    }`)
	err := m.Mask()

	// Then the output should have that sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
  + resource "template_dir" "config" {
      + destination_dir = "./renders"
      + id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
      + source_dir      = "./templates"
      + vars            = {
          "msg1" = "******"
        }
    }`, b)
}

func TestMasker_Mask_StreamedSensitiveDataWithANSIEscapeCodes(t *testing.T) {
	// Given some input content
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// only the known sensitive data) into stdIn
	stdinSim.write(
		`
Terraform will perform the following actions:

[1m  # template_dir.config[0m will be created[0m[0m
[0m[32m  +[0m [0mresource "template_dir" "config" {
      [32m+[0m [0m[1m[0mdestination_dir[0m[0m = "./renders"
      [32m+[0m [0m[1m[0mid[0m[0m              = (known after apply)
      [32m+[0m [0m[1m[0msource_dir[0m[0m      = "./templates"
      [32m+[0m [0m[1m[0mvars[0m[0m            = {
          [32m+[0m [0m"msg1" = "sensitive-value-1-AK#%DJGHS*G"
        }
    }`)
	err := m.Mask()

	// Then the output should have that sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
Terraform will perform the following actions:

  # template_dir.config will be created
  + resource "template_dir" "config" {
      + destination_dir = "./renders"
      + id              = (known after apply)
      + source_dir      = "./templates"
      + vars            = {
          + "msg1" = "******"
        }
    }`, b)
}

func TestMasker_Mask_StreamedNothing2KnownSensitiveDataTransform(t *testing.T) {
	// Given some input content
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-1-AK#%DJGHS*G"}})

	// When we simulate piping this content (which contains
	// a transition from nothing to the known sensitive data) into stdIn
	stdinSim.write(
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
            "msg1" = "normal value 1"
          ~ "msg2" = "" -> "sensitive-value-1-AK#%DJGHS*G"
        }
    }`)
	err := m.Mask()

	// Then the output should have that sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
            "msg1" = "normal value 1"
          ~ "msg2" = "" -> "******"
        }
    }`, b)
}

func TestMasker_Mask_StreamedPrevVal2KnownSensitiveDataTransform(t *testing.T) {
	// Given some input content ...
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts(t)
	defer stdinSim.end()
	defer stdoutSim.end()
	m := NewMasker(ctx, &DefaultReplaceables{
		[]string{"sensitive-value-2"}})

	// When we simulate piping this content (which contains
	// a transition from a previous sensitive value to new
	// known sensitive data) into stdIn
	stdinSim.write(
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
            "msg1" = "normal value 1"
          ~ "msg2" = "sensitive-value-2-//dfhs//" -> "sensitive-value-2"
        }
    }`)

	err := m.Mask()

	// Then the output should have both the previous and
	// known sensitive data masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
            "msg1" = "normal value 1"
          ~ "msg2" = "******" -> "******"
        }
    }`, b)
}

func TestMasker_Mask_StreamedPrevVal2KnownSensitiveDataTransform_IgnorePrev(t *testing.T) {
	// Given some input content and an explicit directive NOT
	// to attempt to detect previous sensitive values ...
	ctx, stdinSim, stdoutSim := defaultTestMaskOpts(t)
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
Terraform will perform the following actions:

  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
            "msg1" = "normal value 1"
          ~ "msg2" = "sensitive-value-2-//dfhs//" -> "sensitive-value-2"
        }
    }`)
	err := m.Mask()

	// Then only the new known sensitive values should be masked
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t,
		`
Terraform will perform the following actions:

  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "4278f6895f67aa77cfbdac8ce2c7342275116eec" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
            "msg1" = "normal value 1"
          ~ "msg2" = "sensitive-value-2-//dfhs//" -> "sensitive-value-2"
        }
    }`, b)
}
