package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/opencredo/terrahelp/terrahelp"
)

const (
	cryptoWrapErrorExitCode = 1
	otherErrorExitCode      = 2
)

func encryptCommand(f func(provider string) *terrahelp.CryptoHandler) cli.Command {

	var noBackup bool
	var bkpExt string
	ctxOpts := &terrahelp.CryptoHandlerOpts{TransformOpts: &terrahelp.TransformOpts{}}

	return cli.Command{
		Name:  "encrypt",
		Usage: "Uses configured provider to encrypt specified content",
		Description: "Using either the 'simple' (default) or 'vault' provider, encrypt will ensure that the relevant content\n" +
			"   is encrypted in either full, or inline mode. You may want to filter out certain sensitive content from the\n" +
			"   output of terraform commands such as a 'terraform plan' or 'terraform apply'. Terrahelp always assumes \n" +
			"   you are piping your content in (for example from stdin) unless you explicitly specify file(s) to read in from \n" +
			"   (and thus also write out to). Either way the encryption process is run over the appropriate content. \n" +
			"   In general it is expected that the following process will occur: \n" +
			"   * The standard terraform commands will be performed i.e. \n" +
			"       terraform plan\n" +
			"       terraform apply\n" +
			"   * Then when happy, the encryption can be applied (for example against the tfstate files) i.e.\n" +
			"       terrahelp encrypt -file=terraform.tfstate -file=terraform.tfstate.backup  \n" +
			"   * Finally, if required, these files can then be checked in to version control.\n\n" +
			"   The desired 'provider' and well as 'mode' can be supplied as CLI arguments, or via the TH_ENCRYPTION_PROVIDER \n" +
			"   and TH_ENCRYPTION_MODE environment variables. Encrypted values will always conform to the following format: \n" +
			"   @terrahelp-encrypted(ENCRYPTED_CONTENT). Decryption essentially operates in reverse. \n\n" +

			"   Encryption modes: full,inline  \n" +
			"   -----------------------------  \n" +
			"   The default encryption mode is 'full' which means the full content (of for example the .tfstate files) \n" +
			"   will be encrypted. Quite often however it is desirable to only encrypt certain sensitive fragments within \n" +
			"   the data, and not simply the whole content as a single unit and this is where the inline mode comes in. \n" +
			"   Inline encrypting works by detecting sensitive values within the state files and then replaces them with \n" +
			"   encrypted values. Terrahelp assumes all sensitive values are defined as values within the terraform.tfvars \n" +
			"   file, which just by way of recap should NEVER be checked into version control! Terrahelp then uses the  \n" +
			"   terraform.tfvars file to identify which values are considered sensitive, searches for any occurrence \n" +
			"   of these values within the provided content and essentially does a find and replace of all the sensitive values \n" +
			"   with appropriately encrypted ones. \n\n" +

			"   Providers: simple,vault,vault-cli \n" +
			"   ---------------------------------  \n" +
			"     * 'simple' provider: \n" +
			"         NOTE: Any arguments specifically related to the simple provider will always have a simple-xxx prefix.   \n" +
			"         This provider performs simple in memory AES encryption, and requires that you pass in either a 16 or 32 \n" +
			"         character key to indicate if you want 128 or 256 bit AES encryption. For example a key with value \n" +
			"         'AES256Key-32Characters0987654321' will result in your content being encrypted with 256 bit AES encryption. \n" +
			"         Note: If you lose access to this encryption key, you will NOT be able to decrypt these values!!\n\n" +

			"     * 'vault' provider (https://www.vaultproject.io): \n" +
			"         NOTE: Any arguments specifically related to the vault provider will always have a vault-xxx prefix.   \n" +
			"         This version of the vault provider expects a running and accessible instance of Vault to be available \n" +
			"         and to have access to it via its http API. Currently validated against Vault 0.5.2 \n\n" +

			"         It is configured via standard Vault environment variables i.e.  \n" +
			"           export VAULT_TOKEN=your-vault-root-token \n" +
			"           export VAULT_ADDR=http://127.0.0.1:8200\n" +
			"           export VAULT_SKIP_VERIFY=true\n\n" +

			"         Vault's transit aka 'encryption as a service' feature is then used to offload and perform the actual \n" +
			"         encryption. (https://www.vaultproject.io/docs/secrets/transit). The Vault transit backend makes use \n" +
			"         of a registered named encryption key to gain access to the underlying encryption key itself, as well \n" +
			"         details about what encryption algorithm to use. This encrypt command expects you to have \n" +
			"         already setup and registered a named encryption key which will be used here. If not explicitly\n" +
			"         specified, then 'terrahelp' i.e. /transit/key/terrahelp is assumed as default. Note you can use the \n" +
			"         terrahelp vault-autoconfig command to auto register and generate a new key against this default \n" +
			"         named key for you if not already done. \n\n" +

			"     * 'vault-cli' provider: \n" +
			"         This version of the Vault provider relies on direct access to the Vault CLI (as opposed to accessing \n" +
			"         Vault via its HTTP API alone). It is however configured, and works, in the same way the vault provider \n" +
			"         described above, with the exception you will also need the Vault CLI to be available in the PATH\n" +

			"   EXIT STATUS \n" +
			"   ----------- \n" +
			"   The encrypt command exits with one of the following values: \n" +
			"        0     Encryption succeeded. \n" +
			"        1     Error encrypting due to terrahelp wrapper issues or violations (e.g. double encrypt not permitted). \n" +
			"        >1    Any other error occurred. \n\n" +

			"   EXAMPLES \n" +
			"   ----------- \n" +
			"   To fully encrypt the terraform.tfstate & terraform.tfstate.backup files using simple encryption:\n\n" +

			"        $  terrahelp encrypt -simple-key=AES256Key-32Characters0987654321 -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"   To inline encrypt the terraform.tfstate & terraform.tfstate.backup files using simple encryption:\n\n" +

			"        $  terrahelp encrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"   To inline encrypt the output of a terraform plan using simple encryption:\n\n" +

			"        $  terraform plan | terrahelp encrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline\n\n" +

			"   To fully encrypt the terraform.tfstate & terraform.tfstate.backup files using vault encryption:\n\n" +

			"        $  terrahelp encrypt -provider=vault vault-namedkey=my-vault-named-key -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"   To inline encrypt the terraform.tfstate & terraform.tfstate.backup files using vault encryption:\n\n" +

			"        $  terrahelp encrypt -provider=vault vault-namedkey=my-vault-named-key -mode=inline -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"\n",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "provider",
				Value:       terrahelp.ThEncryptProviderSimple,
				EnvVar:      "TH_ENCRYPTION_PROVIDER",
				Usage:       "Encryption provider (simple|vault|vault-cli) to use",
				Destination: &ctxOpts.EncProvider,
			},
			cli.StringFlag{
				Name:        "mode",
				Value:       terrahelp.ThEncryptModeFull,
				EnvVar:      "TH_ENCRYPTION_MODE",
				Usage:       fmt.Sprintf("Encryption mode (inline|full) to use"),
				Destination: &ctxOpts.EncMode,
			},
			cli.StringSliceFlag{
				Name:  "file",
				Usage: fmt.Sprintf("File(s) to encrypt - can be specified multiple times"),
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups",
				Destination: &bkpExt,
			},
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before encrypting (defaults to false)",
				Destination: &noBackup,
			},
			cli.StringFlag{
				Name:        "tfvars",
				Value:       terrahelp.TfvarsFilename,
				Usage:       "Terraform tfvars filename",
				Destination: &ctxOpts.TfvarsFilename,
			},
			cli.BoolTFlag{
				Name:        "dblencrypt",
				Usage:       "Permits the double encryption of the content in a file (defaults to true)",
				Destination: &ctxOpts.AllowDoubleEncrypt,
			},
			cli.BoolTFlag{
				Name:        "exclwhitespace",
				Usage:       "Excludes the encryption of whitespace only values (defaults to true)",
				Destination: &ctxOpts.ExcludeWhitespaceOnly,
			},
			cli.StringFlag{
				Name:        "simple-key",
				EnvVar:      "TH_SIMPLE_KEY",
				Usage:       "(Simple provider only) the encryption key to use",
				Destination: &ctxOpts.SimpleKey,
			},
			cli.StringFlag{
				Name:        "vault-namedkey",
				EnvVar:      "TH_VAULT_NAMED_KEY",
				Value:       terrahelp.ThNamedEncryptionKey,
				Usage:       "(Vault provider only) Named encryption key to use",
				Destination: &ctxOpts.NamedEncKey,
			},
		},
		Action: func(c *cli.Context) {
			th := f(ctxOpts.EncProvider)
			err := ctxOpts.ValidateForEncryptDecrypt()
			exitIfError(err)
			setupTransformableItems(c, ctxOpts.TransformOpts, noBackup, bkpExt)
			err = th.Encrypt(ctxOpts)
			exitIfError(err)
		},
	}
}

func decryptCommand(f func(provider string) *terrahelp.CryptoHandler) cli.Command {

	var noBackup bool
	var bkpExt string
	ctxOpts := &terrahelp.CryptoHandlerOpts{TransformOpts: &terrahelp.TransformOpts{}}

	return cli.Command{
		Name:  "decrypt",
		Usage: "Uses configured provider to decrypt specified content",
		Description: "Using either the 'simple' (default) or 'vault' provider, decrypt will ensure that the against content\n" +
			"   is decrypted and restored back to its pre-encrypted state. You may previously have filtered, and possibly saved\n" +
			"   as a file, the output of terraform commands such as a 'terraform plan' or 'terraform apply' which may have contained\n" +
			"   certain sensitive content. Terrahelp always assumes you are piping your content in (for example from stdin) unless \n" +
			"   you explicitly specify file(s) to read in from (and thus also write out to). Either way the decryption process is run \n" +
			"   over the appropriate content. This command essentially does the reverse of the encrypt command, and in " +
			"   general it is expected that:  \n" +

			"   * The encrypted versions of the terraform.tfstate files will be obtained \n" +
			"   * The appropriate decryption applied i.e.\n" +
			"       terrahelp decrypt -file=terraform.tfstate -file=terraform.tfstate.backup \n" +
			"   * The standard terraform commands can then be performed i.e. \n" +
			"       terraform plan\n" +
			"       terraform apply\n\n" +

			"   Note: You need to ensure you use the same 'provider' and 'mode' that was used when the files were \n" +
			"   encrypted in the first place. The 'provider' and well as 'mode' can be supplied as CLI arguments,  \n" +
			"   or via the TH_ENCRYPTION_PROVIDER and TH_ENCRYPTION_MODE environment variables. \n\n" +

			"   EXIT STATUS \n" +
			"   ----------- \n" +
			"   The decrypt command exits with one of the following values: \n" +
			"        0     Decryption succeeded. \n" +
			"        1     Error decrypting due to terrahelp wrapper issues or violations (e.g. not prev encrypted). \n" +
			"        >1    Any other error occurred. \n\n" +

			"   EXAMPLES \n" +
			"   ----------- \n" +
			"   To fully decrypt the terraform.tfstate & terraform.tfstate.backup files using simple encryption:\n\n" +

			"        $  terrahelp decrypt -simple-key=AES256Key-32Characters0987654321 -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"   To inline decrypt the terraform.tfstate & terraform.tfstate.backup files using simple encryption:\n\n" +

			"        $  terrahelp decrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"   To inline decrypt the previously saved output of a terraform plan using simple encryption:\n\n" +

			"        $  cat plan-out.tfplan | terrahelp decrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline\n\n" +

			"   To fully decrypt the terraform.tfstate & terraform.tfstate.backup files using vault encryption:\n\n" +

			"        $  terrahelp decrypt -provider=vault vault-namedkey=my-vault-named-key  -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"   To inline decrypt the terraform.tfstate & terraform.tfstate.backup files using vault encryption:\n\n" +

			"        $  terrahelp decrypt -provider=vault vault-namedkey=my-vault-named-key -mode=inline -file=terraform.tfstate -file=terraform.tfstate.backup \n\n" +

			"\n",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "provider",
				Value:       terrahelp.ThEncryptProviderSimple,
				EnvVar:      "TH_ENCRYPTION_PROVIDER",
				Usage:       "Encryption provider (simple|vault|vault-cli) to use",
				Destination: &ctxOpts.EncProvider,
			},
			cli.StringFlag{
				Name:        "mode",
				Value:       terrahelp.ThEncryptModeFull,
				EnvVar:      "TH_ENCRYPTION_MODE",
				Usage:       fmt.Sprintf("Encryption mode (inline|full) to use"),
				Destination: &ctxOpts.EncMode,
			},
			cli.StringSliceFlag{
				Name:  "file",
				Usage: fmt.Sprintf("File(s) to encrypt - can be specified multiple times"),
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups",
				Destination: &bkpExt,
			},
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before decrypting (defaults to false)",
				Destination: &noBackup,
			},
			cli.StringFlag{
				Name:        "simple-key",
				EnvVar:      "TH_SIMPLE_KEY",
				Usage:       "(Simple provider only) the encryption key to use",
				Destination: &ctxOpts.SimpleKey,
			},
			cli.StringFlag{
				Name:        "vault-namedkey",
				Value:       terrahelp.ThNamedEncryptionKey,
				EnvVar:      "TH_VAULT_NAMED_KEY",
				Usage:       "(Vault provider only) Named encryption key to use",
				Destination: &ctxOpts.NamedEncKey,
			},
		},
		Action: func(c *cli.Context) {
			th := f(ctxOpts.EncProvider)
			err := ctxOpts.ValidateForEncryptDecrypt()
			exitIfError(err)
			setupTransformableItems(c, ctxOpts.TransformOpts, noBackup, bkpExt)
			err = th.Decrypt(ctxOpts)
			exitIfError(err)
		},
	}

}

func vaultAutoConfigCommand(f func(provider string) *terrahelp.CryptoHandler) cli.Command {

	ctxOpts := &terrahelp.CryptoHandlerOpts{TransformOpts: &terrahelp.TransformOpts{}}

	return cli.Command{
		Name:  "vault-autoconfig",
		Usage: "Auto configures Vault with a basic setup to support encrypt and decrypt actions.",
		Description: "This is really a one off helper command to help get off the ground and running quickly with the Vault provider. \n" +
			"   Essentially it ensures that the transit backend is mounted (if not\n" +
			"   already), and that the named encryption key (defaults to /transit/key/terrahelp) is generated and registered \n" +
			"   The vault provider uses this named key as part of the encrypt and decrypt functionality. \n",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "namedkey",
				Value:       terrahelp.ThNamedEncryptionKey,
				EnvVar:      "TH_VAULT_NAMED_KEY",
				Usage:       "Named Vault encryption key to use",
				Destination: &ctxOpts.NamedEncKey,
			},
		},
		Action: func(c *cli.Context) {
			th := f(terrahelp.ThEncryptProviderVault)
			err := th.Init(ctxOpts)
			exitIfError(err)
		},
	}
}

func maskCommand() cli.Command {

	ctxOpts := &terrahelp.MaskOpts{TransformOpts: &terrahelp.TransformOpts{}}
	var noBackup bool
	var bkpExt string

	return cli.Command{
		Name:  "mask",
		Usage: "Mask will overwrite sensitive data in output or files with a masked value (eg. ******).",
		Description: "Given a configured mask pattern (numchars x maskchar) to replace, when this command is run, any sensitive \n" +
			"   data detected (by its presence in the terraform.tfvars file, i.e. the same mechanism used by \n" +
			"   encrypt/decrypt) will be overwritten with the masked value. \n\n" +

			"   The typical use case for mask is with streamed or piped data (.e.g from terraform plan \n" +
			"   or terraform apply). If the sensitive input values change from one run to another, quite often the output\n" +
			"   from these functions will also include details of the changes including the previous (but\n" +
			"   no doubt still sensitive) value. The prev flag is here to try and help in these cases. \n" +
			"   Specifically, if true (default is true), the function will try to detect common patterns typically \n" +
			"   used by terraform commands to denote a sensitive input change, and additionally mask the previous value as well. \n\n" +

			"   EXAMPLES \n" +
			"   ----------- \n" +
			"   To inline mask the output of a terraform plan (with default mask of ******):\n\n" +

			"        $  terraform plan | terrahelp mask\n\n" +

			"   To inline mask the output of a terraform plan (with mask of ###):\n\n" +

			"        $  terraform plan | terrahelp mask -maskchar=# -numchars=3 \n\n" +

			"   To suppress the attempted detection of previous sensitive values when masking the output of a terraform plan:\n\n" +

			"        $  terraform plan | terrahelp mask -prev=false \n\n",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "maskchar",
				Value:       terrahelp.MaskChar,
				Usage:       "Forms mask pattern (numchars x maskchar) to replace sensitive data with",
				Destination: &ctxOpts.MaskChar,
			},
			cli.IntFlag{
				Name:        "numchars",
				Value:       terrahelp.NumberOfMaskChar,
				Usage:       fmt.Sprintf("Forms mask pattern (numchars x maskchar) to replace sensitive data with"),
				Destination: &ctxOpts.MaskNumChar,
			},
			cli.BoolTFlag{
				Name:        "prev",
				Usage:       "Include the attempted detection, and masking of previous sensitive values (defaults to true)",
				Destination: &ctxOpts.ReplacePrevVals,
			},
			cli.StringSliceFlag{
				Name:  "file",
				Usage: "File(s) to have sensitive data replaced with mask - can be specified multiple times",
			},
			cli.StringFlag{
				Name:        "tfvars",
				Value:       terrahelp.TfvarsFilename,
				Usage:       "Terraform tfvars filename, used to detect sensitive vals",
				Destination: &ctxOpts.TfvarsFilename,
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups",
				Destination: &bkpExt,
			},
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before masking (defaults to false)",
				Destination: &noBackup,
			},
			cli.BoolTFlag{
				Name:        "exclwhitespace",
				Usage:       "Excludes the masking of whitespace only values (defaults to true)",
				Destination: &ctxOpts.ExcludeWhitespaceOnly,
			},
			cli.BoolFlag{
				Name:        "enablepre012",
				Usage:       "Configures Terrahelp to process pre 0.12 formated console output, (defaults to false)",
				Destination: &ctxOpts.EnablePre012,
			},
		},
		Action: func(c *cli.Context) {
			setupTransformableItems(c, ctxOpts.TransformOpts, noBackup, bkpExt)
			m := terrahelp.NewMasker(ctxOpts, terrahelp.NewTfVars(ctxOpts.TfvarsFilename, ctxOpts.ExcludeWhitespaceOnly))
			err := m.Mask()
			exitIfError(err)
		},
	}
}

// At present code required to add decent exit code support in the cli library
// is awaiting a 2.0. release (https://github.com/codegangsta/cli/pull/266)
// so until then we have to do a bit of an ugly emergency exit ourselves
func exitIfError(e error) {
	if e != nil {
		fmt.Printf("ERROR occurred : %s\n", e)
		switch e.(type) {
		case *terrahelp.CryptoWrapError:
			os.Exit(cryptoWrapErrorExitCode)
		default:
			os.Exit(otherErrorExitCode)
		}
	}
}

// Determines whether the source of items to transform (encrypt/decrypt/mask)
// should be based on stdIn (no files specified) or specific files as provided
// via the command line. Create and set up the appropriate Transformables
// as options for further processing.
func setupTransformableItems(c *cli.Context,
	ctxOpts *terrahelp.TransformOpts,
	noBackup bool, bkpExt string) {
	files := c.StringSlice("file")

	if files == nil || len(files) == 0 {
		ctxOpts.TransformItems = []terrahelp.Transformable{terrahelp.NewStdStreamTransformable()}
		return
	}
	ctxOpts.TransformItems = []terrahelp.Transformable{}
	for _, f := range files {
		ctxOpts.TransformItems = append(ctxOpts.TransformItems,
			terrahelp.NewFileTransformable(f, !noBackup, bkpExt))
	}
}
