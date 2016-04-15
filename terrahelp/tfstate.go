package terrahelp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// Tfstate holds the details about and exposes actions which can be
// performed on terraform tfstate files
type Tfstate struct {
	Encrypter Encrypter
}

// TfstateOpts holds the options detailing how and on what state files
// to perform the vault based encryption and decryption.
type TfstateOpts struct {
	TfstateFile        string
	TfStateBkpFile     string
	TfvarsFilename     string
	EncProvider        string
	EncMode            string
	NamedEncKey        string
	SimpleKey          string
	BkpExt             string
	NoBackup           bool
	AllowDoubleEncrypt bool
}

// NewDefaultTfstateOpts creates TfstateOpts with all the
// default values set
func NewDefaultTfstateOpts() *TfstateOpts {
	return &TfstateOpts{
		EncProvider:        ThEncryptProviderSimple,
		TfstateFile:        TfstateFilename,
		TfStateBkpFile:     TfstateBkpFilename,
		TfvarsFilename:     TfvarsFilename,
		NamedEncKey:        ThNamedEncryptionKey,
		SimpleKey:          "",
		BkpExt:             ThBkpExtension,
		NoBackup:           false,
		AllowDoubleEncrypt: true,
		EncMode:            ThEncryptModeFull,
	}
}

// Default file related values
const (
	TfstateFilename    = "terraform.tfstate"
	TfstateBkpFilename = "terraform.tfstate.backup"
	TfvarsFilename     = "terraform.tfvars"
	ThBkpExtension     = ".terrahelpbkp"

	// ThNamedEncryptionKey is default Vault named encryption key
	ThNamedEncryptionKey = "terrahelp"

	errMsgAlreadyEncrypted = "Content has already been encrypted, not performing a double encryption!"
)

// Valid encryption providers
const (
	ThEncryptProviderSimple = "simple"
	ThEncryptProviderVault  = "vault"
)

// Valid encryption modes
const (
	ThEncryptModeInline = "inline"
	ThEncryptModeFull   = "full"
)

func (o *TfstateOpts) getEncryptionKey() string {
	switch {
	case (o.EncProvider == ThEncryptProviderSimple):
		return o.SimpleKey
	case (o.EncProvider == ThEncryptProviderVault):
		return o.NamedEncKey
	default:
		return ""
	}
}

// InlineMode returns true if the Encryption mode is 'inline'
func (o *TfstateOpts) InlineMode() bool {
	return o.EncMode == ThEncryptModeInline
}

// ValidateForEncryptDecrypt ensures valid options have been set
// for the encryption / decruption process
func (o *TfstateOpts) ValidateForEncryptDecrypt() error {
	if o.EncProvider == ThEncryptProviderSimple && o.SimpleKey == "" {
		return fmt.Errorf("You must supply a valid simple-key when using the simply provider. " +
			"The simple provider uses AES and so the AES key should be either 16 or 32 byte to select AES-128 or AES-256 encryption")

	}
	if o.EncProvider == ThEncryptProviderVault && o.NamedEncKey == "" {
		return fmt.Errorf("You must supply a vault-namedkey when using the vault provider ")
	}
	return nil
}

// Init provides the opportunity for the Encrypter provider to
// perform any additional config or initialisation which may
// be required before use
func (t *Tfstate) Init(ctx *TfstateOpts) error {
	return t.Encrypter.Init(ctx.NamedEncKey)
}

// Encrypt will ensure appropriate aspects of the tfstate files
// are encrypted as per the configured options supplied
func (t *Tfstate) Encrypt(ctx *TfstateOpts) error {
	for _, f := range []string{ctx.TfstateFile, ctx.TfStateBkpFile} {
		if _, err := os.Stat(f); f != "" && err == nil {
			err := t.encrypt(ctx, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Decrypt will ensure appropriate aspects of the tfstate files
// are decrypted as per the configured options supplied
func (t *Tfstate) Decrypt(ctx *TfstateOpts) error {
	for _, f := range []string{ctx.TfstateFile, ctx.TfStateBkpFile} {
		if _, err := os.Stat(f); f != "" && err == nil {
			err := t.decrypt(ctx, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tfstate) encrypt(ctx *TfstateOpts, f string) error {

	// Backup if required
	err := t.backup(!ctx.NoBackup, f, ctx.BkpExt)
	if err != nil {
		return err
	}

	// Encrypt and write out content
	var b []byte
	if ctx.InlineMode() {
		log.Printf("Encrypting inline: %s ", f)
		b, err = t.encryptInline(f, ctx.getEncryptionKey(), ctx.TfvarsFilename, ctx.AllowDoubleEncrypt)
	} else {
		log.Printf("Encrypting: %s ", f)
		b, err = t.encryptFileContent(f, ctx.getEncryptionKey(), ctx.AllowDoubleEncrypt)
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, b, 0777)
}

func (t *Tfstate) decrypt(ctx *TfstateOpts, f string) error {

	// Backup if required
	err := t.backup(!ctx.NoBackup, f, ctx.BkpExt)
	if err != nil {
		return err
	}

	// Read and decrypt content
	log.Printf("Decrypting %s ", f)
	var plain []byte
	ciphertext, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	if ctx.InlineMode() {
		plain, err = t.decryptInline(ciphertext, ctx.getEncryptionKey())
	} else {
		plain, err = t.Encrypter.Decrypt(ctx.getEncryptionKey(), ciphertext)
	}
	if err != nil {
		return err
	}

	// Write out plaintext values again
	return ioutil.WriteFile(f, plain, 0777)
}

func (t *Tfstate) backup(b bool, f, bkpext string) error {
	if b {
		bkp := f + bkpext
		log.Printf("Backuping up %s --> %s ", f, bkp)
		return CopyFile(f, bkp)
	}
	return nil
}

func (t *Tfstate) decryptInline(b []byte, key string) ([]byte, error) {
	r := regexp.MustCompile(thCryptoWrapRegExp)
	m := r.FindAll(b, -1)
	for _, j := range m {
		dec, err := t.Encrypter.Decrypt(key, j)
		if err != nil {
			return nil, err
		}
		b = bytes.Replace(b, j, dec, -1)
	}
	return b, nil
}

func (t *Tfstate) encryptFileContent(f, key string, dblEncrypt bool) ([]byte, error) {
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if !dblEncrypt {
		r := regexp.MustCompile(thCryptoWrapRegExp)
		m := r.FindSubmatch(c)
		if len(m) >= 1 {
			return nil, fmt.Errorf(errMsgAlreadyEncrypted)
		}
	}

	return t.Encrypter.Encrypt(key, c)
}

// tfsf = tfstate file, tfvf = tfvars file
func (t *Tfstate) encryptInline(tfsf, key, tfvf string, dblEncrypt bool) ([]byte, error) {
	plain, err := ioutil.ReadFile(tfsf)
	if err != nil {
		return nil, err
	}

	if !dblEncrypt {
		r := regexp.MustCompile(thCryptoWrapRegExp)
		m := r.FindSubmatch(plain)
		if len(m) >= 1 {
			return nil, fmt.Errorf(errMsgAlreadyEncrypted)
		}
	}

	tfvu := &Tfvars{}
	inlineCreds, err := tfvu.ExtractSensitiveVals(tfvf)
	if err != nil {
		return nil, err
	}

	inlinedText := string(plain)
	for _, v := range inlineCreds {
		ct, err := t.Encrypter.Encrypt(key, []byte(v))
		if err != nil {
			return nil, err
		}
		inlinedText = strings.Replace(inlinedText, v, string(ct), -1)
	}

	return []byte(inlinedText), nil
}
