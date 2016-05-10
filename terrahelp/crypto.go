package terrahelp

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// CryptoHandler defines and exposes cryptographic actions which
// can be performed against terraform related files and output
type CryptoHandler struct {
	Encrypter Encrypter
}

// CryptoHandlerOpts holds the specific options detailing how, and on what
// to perform the cryptographic actions.
type CryptoHandlerOpts struct {
	*TransformOpts
	EncProvider        string
	EncMode            string
	NamedEncKey        string
	SimpleKey          string
	AllowDoubleEncrypt bool
}

// NewDefaultCryptoHandlerOpts creates CryptoHandlerOpts with all the
// default values set
func NewDefaultCryptoHandlerOpts() *CryptoHandlerOpts {
	return &CryptoHandlerOpts{
		TransformOpts: &TransformOpts{TransformItems: []Transformable{
			NewFileTransformable(TfstateFilename, true, ThBkpExtension),
			NewFileTransformable(TfstateBkpFilename, true, ThBkpExtension)},
			TfvarsFilename: TfvarsFilename},
		EncProvider:        ThEncryptProviderSimple,
		NamedEncKey:        ThNamedEncryptionKey,
		SimpleKey:          "",
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
	ThNamedEncryptionKey   = "terrahelp"
	errMsgAlreadyEncrypted = "Content has already been encrypted, and double encryption has been disabled!"
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

type cryptoTransformAction func(*CryptoHandlerOpts, Transformable) error

func (o *CryptoHandlerOpts) getEncryptionKey() string {
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
func (o *CryptoHandlerOpts) InlineMode() bool {
	return o.EncMode == ThEncryptModeInline
}

// ValidateForEncryptDecrypt ensures valid options have been set
// for the encryption / decruption process
func (o *CryptoHandlerOpts) ValidateForEncryptDecrypt() error {
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
func (t *CryptoHandler) Init(ctx *CryptoHandlerOpts) error {
	return t.Encrypter.Init(ctx.NamedEncKey)
}

// Encrypt will ensure the appropriate areas of the input content
// are encrypted as per the configured options supplied
func (t *CryptoHandler) Encrypt(ctx *CryptoHandlerOpts) error {
	return t.applyCryptoAction(ctx,
		func(ctx *CryptoHandlerOpts, ci Transformable) error { return t.encrypt(ctx, ci) })
}

// Decrypt will ensure appropriate aspects of the input content
// are decrypted as per the configured options supplied
func (t *CryptoHandler) Decrypt(ctx *CryptoHandlerOpts) error {
	return t.applyCryptoAction(ctx,
		func(ctx *CryptoHandlerOpts, ci Transformable) error { return t.decrypt(ctx, ci) })
}

// applyCryptoAction will loop through the various items to be encrypted/decrypted, first
// validating them, and if OK, proceeds to send them on for the actual encryption/decryption
func (t *CryptoHandler) applyCryptoAction(ctx *CryptoHandlerOpts, cryptoFunc cryptoTransformAction) error {
	// Do first pass over all items to be encrypted/decrypted to ensure they are
	// all valid before we begin
	for _, ci := range ctx.TransformItems {
		if err := ci.validate(); err != nil {
			log.Printf("Not a valid item for encryption: %v\n", err)
			return err
		}
	}

	for _, ci := range ctx.TransformItems {
		err := cryptoFunc(ctx, ci)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *CryptoHandler) encrypt(ctx *CryptoHandlerOpts, ci Transformable) error {

	// Do any pre encryption actions (e.g. backup)
	// if required
	err := ci.beforeTransform()
	if err != nil {
		return err
	}

	// Read, encrypt, then write out result
	in, err := ci.read()
	if err != nil {
		return err
	}
	b, err := t.encryptBytes(ctx, in)
	if err != nil {
		return err
	}
	return ci.write(b)
}

func (t *CryptoHandler) encryptBytes(ctx *CryptoHandlerOpts, in []byte) ([]byte, error) {
	if ctx.InlineMode() {
		return t.encryptInline(in, ctx.getEncryptionKey(), ctx.TfvarsFilename, ctx.AllowDoubleEncrypt)
	}
	return t.encryptFullContent(in, ctx.getEncryptionKey(), ctx.AllowDoubleEncrypt)
}

func (t *CryptoHandler) decrypt(ctx *CryptoHandlerOpts, ci Transformable) error {

	// Do any pre decryption actions (e.g. backup)
	// if required
	err := ci.beforeTransform()
	if err != nil {
		return err
	}

	// Read, decrypt, then write out result
	ciphertext, err := ci.read()
	if err != nil {
		return err
	}
	plain, err := t.decryptBytes(ctx, ciphertext)
	if err != nil {
		return err
	}

	// Write out plaintext values again
	return ci.write(plain)
}

func (t *CryptoHandler) decryptBytes(ctx *CryptoHandlerOpts, in []byte) ([]byte, error) {
	if ctx.InlineMode() {
		return t.decryptInline(in, ctx.getEncryptionKey())
	}
	return t.Encrypter.Decrypt(ctx.getEncryptionKey(), in)
}

func (t *CryptoHandler) decryptInline(b []byte, key string) ([]byte, error) {
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

func (t *CryptoHandler) encryptFullContent(b []byte, key string, dblEncrypt bool) ([]byte, error) {

	if !dblEncrypt {
		r := regexp.MustCompile(thCryptoWrapRegExp)
		m := r.FindSubmatch(b)
		if len(m) >= 1 {
			return nil, newCryptoWrapError(errMsgAlreadyEncrypted)
		}
	}

	return t.Encrypter.Encrypt(key, b)
}

// tfvf = tfvars file
func (t *CryptoHandler) encryptInline(plain []byte, key, tfvf string, dblEncrypt bool) ([]byte, error) {

	if !dblEncrypt {
		r := regexp.MustCompile(thCryptoWrapRegExp)
		m := r.FindSubmatch(plain)
		if len(m) >= 1 {
			return nil, newCryptoWrapError(errMsgAlreadyEncrypted)
		}
	}

	tfvu := NewTfVars(tfvf)
	inlineCreds, err := tfvu.Values()
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
