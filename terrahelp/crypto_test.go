package terrahelp

import (
	"github.com/stretchr/testify/assert"

	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func newVaultEncryptableCryptoHandler(t *testing.T) (*CryptoHandler, *MockVaultClient) {
	vc := NewMockVaultClient()
	cu, err := createVaultEncrypter(vc)
	if err != nil {
		t.Fatalf("Unabled to create test Tfstate : %s", err)
	}
	return &CryptoHandler{cu}, vc
}

func newInitVaultEncryptableCrytoHandler(t *testing.T, key string) (*CryptoHandler, *MockVaultClient) {
	tu, vc := newVaultEncryptableCryptoHandler(t)
	ctx := NewDefaultCryptoHandlerOpts()
	ctx.EncProvider = ThEncryptProviderVault
	err := tu.Init(ctx)
	if err != nil {
		t.Fatalf("Unabled to initialise test Tfstate : %s", err)
	}
	return tu, vc
}

func defaultVaultEncryptableCryptoHandlerOpts(t *testing.T, noBkp bool) *CryptoHandlerOpts {
	ctx := NewDefaultCryptoHandlerOpts()
	ctx.EncProvider = ThEncryptProviderVault
	ctx.TransformItems[0].(*FileTransformable).bkp = !noBkp
	ctx.TransformItems[1].(*FileTransformable).bkp = !noBkp
	return ctx
}

func defaultTestInlineCryptoHandlerOpts(t *testing.T, noBkp bool) *CryptoHandlerOpts {
	ctx := NewDefaultCryptoHandlerOpts()
	ctx.EncProvider = ThEncryptProviderVault
	ctx.EncMode = ThEncryptModeInline
	ctx.TransformItems[0].(*FileTransformable).bkp = !noBkp
	ctx.TransformItems[1].(*FileTransformable).bkp = !noBkp
	return ctx
}

type stdinSim struct {
	t            *testing.T
	simReadFile  *os.File
	simWriteFile *os.File
	simWriter    *bufio.Writer
}

type stdoutSim struct {
	t            *testing.T
	simWriteFile *os.File
	simReader    *bufio.Reader
}

func (s *stdoutSim) start() {
	f, err := ioutil.TempFile("", "stdout-sim")
	if err != nil {
		s.t.Fatalf("Unabled to create tmp file to sim stdout %s", err)
	}
	s.simWriteFile, err = os.OpenFile(f.Name(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		s.t.Fatalf("(2) Unabled to create writable tmp file to sim stdout %s", err)
	}
}

func (s *stdoutSim) getAllContent() string {
	b, err := ioutil.ReadFile(s.simWriteFile.Name())
	if err != nil {
		s.t.Fatalf("Unabled to get all content from tmp file simulating stdout %s", err)
	}
	return string(b)
}

func (s *stdoutSim) end() {
	if s.simWriteFile == nil {
		s.t.Fatal("Unabled to end stdout sim (it probably wasn't started)")
	}
	err := s.simWriteFile.Close()
	if err != nil {
		s.t.Fatalf("Unabled to close writable tmp file to sim stdout %s", err)
	}
	err = os.Remove(s.simWriteFile.Name())
	if err != nil {
		s.t.Fatalf("Unabled to cleanup and delete writable tmp file simulating stdout %s", err)
	}

}

func newStdinSim(t *testing.T) *stdinSim {
	return &stdinSim{t: t}
}

func newStdoutSim(t *testing.T) *stdoutSim {
	return &stdoutSim{t: t}
}

func (s *stdinSim) start() {
	f, err := ioutil.TempFile("", "stdin-sim")
	if err != nil {
		s.t.Fatalf("Unabled to create tmp file to sim stdin %s", err)
	}
	s.simReadFile, err = os.OpenFile(f.Name(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		s.t.Fatalf("Unabled to create readable file to sim stdin %s", err)
	}
	s.simWriteFile, err = os.OpenFile(f.Name(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		s.t.Fatalf("Unabled to create writable tmp file to sim stdin %s", err)
	}
	s.simWriter = bufio.NewWriter(s.simWriteFile)
}

func (s *stdinSim) end() {
	if s.simReadFile == nil {
		s.t.Fatal("Unabled to end stdin sim (it probably wasn't started)")
	}
	s.simReadFile.Close()
	s.simWriteFile.Close()
}

func (s *stdinSim) write(in string) {
	_, err := s.simWriter.WriteString(in)
	if err != nil {
		s.t.Fatalf("Unabled to write string to tmp file simulating stdin %s", err)
	}
	err = s.simWriter.Flush()
	if err != nil {
		s.t.Fatalf("Unabled to flush tmp file simulating stdin %s", err)
	}
}

func defaultTestInlinePipedCryptoHandlerOpts(t *testing.T) (*CryptoHandlerOpts, *stdinSim, *stdoutSim) {
	ctx := NewDefaultCryptoHandlerOpts()
	ctx.EncProvider = ThEncryptProviderVault
	ctx.EncMode = ThEncryptModeInline
	stdinSim := newStdinSim(t)
	stdinSim.start()
	stdoutSim := newStdoutSim(t)
	stdoutSim.start()
	ctx.TransformItems = []Transformable{
		NewStreamTransformable(stdinSim.simReadFile, stdoutSim.simWriteFile)}
	return ctx, stdinSim, stdoutSim
}

func newVaultEncryptableExampleProject(t *testing.T, ver string) (*testProject, *CryptoHandler, *MockVaultClient) {
	tctx := newTempProject(t)
	tctx.copyExampleProject(ver)
	tfu, mvc := newInitVaultEncryptableCrytoHandler(t, ThNamedEncryptionKey)
	return tctx, tfu, mvc
}

func TestCryptoHandler_VaultEncrypter_Init(t *testing.T) {
	// Given
	tu, vc := newVaultEncryptableCryptoHandler(t)
	// When
	ctx := NewDefaultCryptoHandlerOpts()
	ctx.NamedEncKey = "bob"
	err := tu.Init(ctx)
	// Then
	assert.NoError(t, err)
	b1, _ := vc.transitMountExists()
	b2, _ := vc.namedEncryptionKeyExists("bob")
	assert.True(t, b1, "transit mount not registered")
	assert.True(t, b2, "transit key not registered")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_StreamedNonSensitiveData(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	var data = `hello there
                    I am some data
                    to be piped in`
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()

	ctx, stdinSim, stdoutSim := defaultTestInlinePipedCryptoHandlerOpts(t)
	defer stdinSim.end()
	defer stdoutSim.end()

	// When
	stdinSim.write(data)
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t, data, b)
}

func TestCryptoHandler_VaultEncrypter_Encrypt_StreamedSensitiveData(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()

	ctx, stdinSim, stdoutSim := defaultTestInlinePipedCryptoHandlerOpts(t)
	defer stdinSim.end()
	defer stdoutSim.end()

	// When
	stdinSim.write(`hello there
                         sensitive-value-1-AK#%DJGHS*G
                         the bit above should be
                         encrypted`)
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	b := stdoutSim.getAllContent()
	assert.Equal(t, `hello there
                         @terrahelp-encrypted(vault:v1:YzJWdWMybDBhWFpsTFhaaGJIVmxMVEV0UVVzakpVUktSMGhUS2tjPQ==)
                         the bit above should be
                         encrypted`, b)
}

func TestCryptoHandler_VaultEncrypter_Encrypt_inline(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	primary tfstate and backup files should be encrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/encrypted-inline/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/encrypted-inline/terraform.tfstate.backup")
	// 	backups of originals should exist
	tp.assertExpectedFileContent(TfstateFilename+ThBkpExtension, "test-data/example-project/original/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename+ThBkpExtension, "test-data/example-project/original/terraform.tfstate.backup")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_inlineDoubleEncryptError(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)
	ctx.AllowDoubleEncrypt = false

	// And a first time successful encryption ...
	err := tu.Encrypt(ctx)
	assert.NoError(t, err)

	// When (2nd double encryption applied)
	err = tu.Encrypt(ctx)
	assert.Error(t, err, "Expected error if value already encrypted at least once")
	assert.IsType(t, newCryptoWrapError(errMsgAlreadyEncrypted), err)
	assert.Contains(t, err.Error(), errMsgAlreadyEncrypted, "not detected as invalid wrapper")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_inlineDoubleEncryptAllowed(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)

	// And a first time successful encryption ...
	err := tu.Encrypt(ctx)
	assert.NoError(t, err)

	// When (2nd double encryption applied)
	err = tu.Encrypt(ctx)
	assert.NoError(t, err, "No error expected when double encryption allowed")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_DoubleEncryptError(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultVaultEncryptableCryptoHandlerOpts(t, false)
	ctx.AllowDoubleEncrypt = false

	// And a first time successful encryption ...
	err := tu.Encrypt(ctx)
	assert.NoError(t, err)

	// When (2nd double encryption applied)
	err = tu.Encrypt(ctx)
	assert.Error(t, err, "Expected error if value already encrypted at least once")
	assert.IsType(t, newCryptoWrapError(errMsgAlreadyEncrypted), err)
	assert.Contains(t, err.Error(), errMsgAlreadyEncrypted, "not detected as invalid wrapper")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_DoubleEncryptAllowed(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultVaultEncryptableCryptoHandlerOpts(t, false)

	// And a first time successful encryption ...
	err := tu.Encrypt(ctx)
	assert.NoError(t, err)

	// When (2nd double encryption applied)
	err = tu.Encrypt(ctx)
	assert.NoError(t, err, "No error expected when double encryption allowed")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_inlineNoBackups(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, true)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	primary tfstate and backup files should be encrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/encrypted-inline/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/encrypted-inline/terraform.tfstate.backup")
	// 	no backups should exist
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateFilename+ThBkpExtension))
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename+ThBkpExtension))
}

func TestCryptoHandler_VaultEncrypter_Decrypt_inline(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-inline")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)

	// When
	err := tu.Decrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	backups of originals should exist
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/original/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/original/terraform.tfstate.backup")
	// 	backups of originals should exist
	tp.assertExpectedFileContent(TfstateFilename+ThBkpExtension, "test-data/example-project/encrypted-inline/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename+ThBkpExtension, "test-data/example-project/encrypted-inline/terraform.tfstate.backup")
}

func TestCryptoHandler_VaultEncrypter_Decrypt_inlineNoBackups(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-inline")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, true)

	// When
	err := tu.Decrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	primary tfstate and backup files should be decrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/original/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/original/terraform.tfstate.backup")
	// 	no backups should exist
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateFilename+ThBkpExtension))
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename+ThBkpExtension))
}

func TestCryptoHandler_VaultEncrypter_Encrypt_wholeFile(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultVaultEncryptableCryptoHandlerOpts(t, false)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	primary tfstate and backup files should be encrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/encrypted-wholefile/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/encrypted-wholefile/terraform.tfstate.backup")
	// 	backups of originals should exist
	tp.assertExpectedFileContent(TfstateFilename+ThBkpExtension, "test-data/example-project/original/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename+ThBkpExtension, "test-data/example-project/original/terraform.tfstate.backup")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_wholeFileNoBackups(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultVaultEncryptableCryptoHandlerOpts(t, true)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	primary tfstate and backup files should be encrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/encrypted-wholefile/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/encrypted-wholefile/terraform.tfstate.backup")
	// 	no backups should exist
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateFilename+ThBkpExtension))
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename+ThBkpExtension))
}

func TestCryptoHandler_VaultEncrypter_Decrypt_wholefile(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-wholefile")
	defer tp.restore()
	ctx := defaultVaultEncryptableCryptoHandlerOpts(t, false)

	// When
	err := tu.Decrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	backups of originals should exist
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/original/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/original/terraform.tfstate.backup")
	// 	backups of originals should exist
	tp.assertExpectedFileContent(TfstateFilename+ThBkpExtension, "test-data/example-project/encrypted-wholefile/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename+ThBkpExtension, "test-data/example-project/encrypted-wholefile/terraform.tfstate.backup")
}

func TestCryptoHandler_VaultEncrypter_Decrypt_wholefileNoBackups(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-wholefile")
	defer tp.restore()
	ctx := defaultVaultEncryptableCryptoHandlerOpts(t, true)

	// When
	err := tu.Decrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	primary tfstate and backup files should be decrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/original/terraform.tfstate")
	tp.assertExpectedFileContent(TfstateBkpFilename, "test-data/example-project/original/terraform.tfstate.backup")
	// 	no backups should exist
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateFilename+ThBkpExtension))
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename+ThBkpExtension))
}

func TestCryptoHandler_VaultEncrypter_Decrypt_wholefile_prevEncryptedInline(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-inline")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)
	ctx.EncMode = ThEncryptModeFull

	// When
	err := tu.Decrypt(ctx)

	// Then
	assert.Error(t, err, "Expected error if named key not valid")
	assert.Contains(t, err.Error(), thCryptoWrapInvalidMsg,
		fmt.Sprint("not detected as invalid wrapper"))
}

func TestCryptoHandler_VaultEncrypter_Encrypt_invalidNamedKey(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)
	ctx.NamedEncKey = "nonexistant-named-key"

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.Error(t, err, "Expected error if named key not valid")
}

func TestCryptoHandler_VaultEncrypter_Encrypt_missingTfstateBkpFile(t *testing.T) {
	// Given a known original project setup in temp dir
	// (with missing bkp file)  ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)
	tp.removeProjectFile(TfstateBkpFilename)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("stat %s: no such file or directory", TfstateBkpFilename))
	// And nothing (including main tfstate file should be encrypted)
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/original/terraform.tfstate")
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename))
}

func TestCryptoHandler_VaultEncrypter_Encrypt_inline_missingTfvarsFile(t *testing.T) {
	// Given a known original project setup in temp dir
	// (with missing bkp file)  ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineCryptoHandlerOpts(t, false)
	tp.removeProjectFile(TfvarsFilename)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.Error(t, err, "Missing tfvars should result in an error")
}
