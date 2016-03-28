package terrahelp

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func newVaultEncryptableTfstate(t *testing.T) (*Tfstate, *MockVaultClient) {
	vc := NewMockVaultClient()
	cu, err := createVaultEncrypter(vc)
	if err != nil {
		t.Fatalf("Unabled to create test Tfstate : %s", err)
	}
	return &Tfstate{cu}, vc
}

func newInitVaultEncryptableTfstate(t *testing.T, key string) (*Tfstate, *MockVaultClient) {
	tu, vc := newVaultEncryptableTfstate(t)
	ctx := NewDefaultTfstateOpts()
	ctx.EncProvider = ThEncryptProviderVault
	err := tu.Init(ctx)
	if err != nil {
		t.Fatalf("Unabled to initialise test Tfstate : %s", err)
	}
	return tu, vc
}

func defaultVaultEncryptableTfstateOpts(t *testing.T, bkp bool) *TfstateOpts {
	ctx := NewDefaultTfstateOpts()
	ctx.EncProvider = ThEncryptProviderVault
	ctx.NoBackup = bkp
	return ctx
}

func defaultTestInlineTfstateOpts(t *testing.T, bkp bool) *TfstateOpts {
	ctx := NewDefaultTfstateOpts()
	ctx.EncProvider = ThEncryptProviderVault
	ctx.Inline = true
	ctx.NoBackup = bkp
	return ctx
}

func newVaultEncryptableExampleProject(t *testing.T, ver string) (*testProject, *Tfstate, *MockVaultClient) {
	tctx := newTempProject(t)
	tctx.copyExampleProject(ver)
	tfu, mvc := newInitVaultEncryptableTfstate(t, ThNamedEncryptionKey)
	return tctx, tfu, mvc
}

func TestTfstate_VaultEncrypter_Init(t *testing.T) {
	// Given
	tu, vc := newVaultEncryptableTfstate(t)
	// When
	ctx := NewDefaultTfstateOpts()
	ctx.NamedEncKey = "bob"
	err := tu.Init(ctx)
	// Then
	assert.NoError(t, err)
	b1, _ := vc.transitMountExists()
	b2, _ := vc.namedEncryptionKeyExists("bob")
	assert.True(t, b1, "transit mount not registered")
	assert.True(t, b2, "transit key not registered")
}

func TestTfstate_VaultEncrypter_Encrypt_inline(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, false)

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

func TestTfstate_VaultEncrypter_Encrypt_inlineNoBackups(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, true)

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

func TestTfstate_VaultEncrypter_Decrypt_inline(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-inline")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, false)

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

func TestTfstate_VaultEncrypter_Decrypt_inlineNoBackups(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-inline")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, true)

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

func TestTfstate_VaultEncrypter_Encrypt_wholeFile(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultVaultEncryptableTfstateOpts(t, false)

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

func TestTfstate_VaultEncrypter_Encrypt_wholeFileNoBackups(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultVaultEncryptableTfstateOpts(t, true)

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

func TestTfstate_VaultEncrypter_Decrypt_wholefile(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-wholefile")
	defer tp.restore()
	ctx := defaultVaultEncryptableTfstateOpts(t, false)

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

func TestTfstate_VaultEncrypter_Decrypt_wholefileNoBackups(t *testing.T) {
	// Given a known encrypted project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-wholefile")
	defer tp.restore()
	ctx := defaultVaultEncryptableTfstateOpts(t, true)

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

func TestTfstate_VaultEncrypter_Decrypt_wholefile_prevEncryptedInline(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "encrypted-inline")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, false)
	ctx.Inline = false

	// When
	err := tu.Decrypt(ctx)

	// Then
	assert.Error(t, err, "Expected error if named key not valid")
	assert.Equal(t, "Unable to decrypt ciphertext, not wrapped as expected", err.Error())
}

func TestTfstate_VaultEncrypter_Encrypt_invalidNamedKey(t *testing.T) {
	// Given a known original project setup in temp dir
	// and we are in the project dir ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, false)
	ctx.NamedEncKey = "nonexistant-named-key"

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.Error(t, err, "Expected error if named key not valid")
}

func TestTfstate_VaultEncrypter_Encrypt_missingTfstateBkpFile(t *testing.T) {
	// Given a known original project setup in temp dir
	// (with missing bkp file)  ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, false)
	tp.removeProjectFile(TfstateBkpFilename)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.NoError(t, err)
	// 	encrypted
	tp.assertExpectedFileContent(TfstateFilename, "test-data/example-project/encrypted-inline/terraform.tfstate")
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename))
	// 	backups of only main file
	tp.assertExpectedFileContent(TfstateFilename+ThBkpExtension, "test-data/example-project/original/terraform.tfstate")
	assertFileDoesNotExist(t, tp.getProjectFile(TfstateBkpFilename+ThBkpExtension))
}

func TestTfstate_VaultEncrypter_Encrypt_inline_missingTfvarsFile(t *testing.T) {
	// Given a known original project setup in temp dir
	// (with missing bkp file)  ...
	tp, tu, _ := newVaultEncryptableExampleProject(t, "original")
	defer tp.restore()
	ctx := defaultTestInlineTfstateOpts(t, false)
	tp.removeProjectFile(TfvarsFilename)

	// When
	err := tu.Encrypt(ctx)

	// Then
	assert.Error(t, err, "Missing tfvars should result in an error")
}
