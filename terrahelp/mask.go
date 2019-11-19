package terrahelp

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"log"
	"regexp"
	"strings"
)

// Masker exposes the ability to obfuscate sensitive data found within certain
// content by replacing it with a masked value
type Masker struct {
	ctx         *MaskOpts
	replacables Replaceables
}

// NewMasker creates a new NewMasker with the specified options
// and Replaceables
func NewMasker(ctx *MaskOpts, svh Replaceables) *Masker {
	if svh == nil {
		svh = &DefaultReplaceables{Vals: []string{}}
	}
	return &Masker{ctx: ctx, replacables: svh}
}

// MaskOpts holds the specific options detailing how, and on what
// to perform the masking action.
type MaskOpts struct {
	*TransformOpts
	MaskChar              string
	MaskNumChar           int
	ReplacePrevVals       bool
	ExcludeWhitespaceOnly bool
	EnablePre012          bool
}

func (m *MaskOpts) getMask() string {
	return strings.Repeat(m.MaskChar, m.MaskNumChar)
}

// NewDefaultMaskOpts creates MaskOpts with all the
// default values set
func NewDefaultMaskOpts() *MaskOpts {
	return &MaskOpts{
		TransformOpts:   &TransformOpts{TfvarsFilename: TfvarsFilename},
		MaskChar:        MaskChar,
		MaskNumChar:     NumberOfMaskChar,
		ReplacePrevVals: true,
	}
}

// Default mask related values
const (
	MaskChar         = "*"
	NumberOfMaskChar = 6

	PrevVal2CurrentValConsolePatternPre012 = "\"(.+)\"\\s*=>\\s*\"(\\%s*)\""
	PrevVal2MaskedValConsolePatternPre012  = "\"%s\" => \"%s\""

	PrevVal2CurrentValConsolePattern = "= \"(.+)\"\\s*->\\s*\"(\\%s*)\""
	PrevVal2MaskedValConsolePattern  = "= \"%s\" -> \"%s\""
)

// Mask will ensure the appropriate areas of the input content
// are replaced with the mask pattern as per the configured options
func (m *Masker) Mask() error {
	if len(m.ctx.TransformItems) == 0 {
		log.Printf("No piped input detected, nor any files provided to apply mask over\n")
		return nil
	}

	for _, ci := range m.ctx.TransformItems {
		if err := ci.validate(); err != nil {
			log.Printf("Not a valid item for masking: %v\n", err)
			return err
		}
	}

	for _, ci := range m.ctx.TransformItems {
		if err := m.mask(ci); err != nil {
			return err
		}
	}
	return nil
}

func (m *Masker) mask(t Transformable) error {

	// Do any pre transformation actions (e.g. backup)
	// if required
	err := t.beforeTransform()
	if err != nil {
		return err
	}

	// Read, mask, then write out result
	in, err := t.read()
	if err != nil {
		return err
	}

	b, err := m.maskBytes(in)
	if err != nil {
		return err
	}
	return t.write(b)
}

func (m *Masker) maskBytes(plain []byte) ([]byte, error) {

	// Convert and strip out the ascii colours.
	inlinedText := stripansi.Strip(string(plain))
	sensitiveVals, err := m.replacables.Values()
	if err != nil {
		return nil, err
	}

	for _, v := range sensitiveVals {
		inlinedText = strings.Replace(inlinedText, v, m.ctx.getMask(), -1)
	}

	if m.ctx.ReplacePrevVals {
		// Additionally there are some patterns (specifically when doing terraform plans
		// and apply where previous sensitive values may also be exposed. We try to catch
		// these too

		selectPattern := PrevVal2CurrentValConsolePattern
		replacePattern := PrevVal2MaskedValConsolePattern

		// If the enable 0.11 flag has been set then use the tf 0.11 patterns
		if m.ctx.EnablePre012 {
			selectPattern = PrevVal2CurrentValConsolePatternPre012
			replacePattern = PrevVal2MaskedValConsolePatternPre012
		}

		r := regexp.MustCompile(fmt.Sprintf(selectPattern, m.ctx.MaskChar))
		inlinedText = r.ReplaceAllString(inlinedText,
			fmt.Sprintf(replacePattern, m.ctx.getMask(), m.ctx.getMask()))
	}
	return []byte(inlinedText), nil
}
