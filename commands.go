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

func encryptSubCommand(f func(provider string) *terrahelp.Tfstate) cli.Command {

	ctxOpts := &terrahelp.TfstateOpts{}

	return cli.Command{
		Name:  "encrypt",
		Usage: "Uses configured provider to encrypt local .tfstate files",
		Description: "Using either the 'simple' (default) or 'vault' provider, encrypt will ensure that the local terraform.tfstate\n" +
			"   files are encrypted in either full, or inline mode. In general it is expected that: \n" +
			"   * The standard terraform commands will be performed i.e. \n" +
			"       terraform plan\n" +
			"       terraform apply\n" +
			"   * Then when happy, the encryption applied i.e.\n" +
			"       terrahelp tfstate encrypt\n" +
			"   * Finally, if required, these files can then be checked in to version control.\n\n" +
			"   The desired 'provider' and well as 'mode' can be supplied as CLI arguments, or via the TH_ENCRYPTION_PROVIDER \n" +
			"   and TH_ENCRYPTION_MODE environment variables. Encrypted values will always conform to the following format: \n" +
			"   @terrahelp-encrypted(ENCRYPTED_CONTENT). Decryption essentially operates in reverse. \n\n" +

			"   Encryption modes: full,inline  \n" +
			"   -----------------------------  \n" +
			"   The default encryption mode is 'full' which means the full content of the tfstate files will be encrypted. \n" +
			"   Quite often however it is desirable to only encrypt the sensitive aspects within the tfstate files and \n" +
			"   not the whole file, this is where the inline mode comes in. \n" +
			"   Inline encrypting works by detecting sensitive values within the state files and then replaces them with \n" +
			"   encrypted values. Terrahelp assumes all sensitive values are defined as values within the terraform.tfvars \n" +
			"   file, which just by way of recap should NEVER be checked into version control! Terrahelp then uses the  \n" +
			"   terraform.tfvars file to identify which values are considered sensitive, searches for any occurence \n" +
			"   of these values within the tfstate files and essentially does a find and replace of all the sensitive values \n" +
			"   with appropriately encrypted ones. \n\n" +

			"   Simple provider \n" +
			"   ---------------  \n" +
			"   NOTE: Any arguments specifically related to the simple provider will always have a simple-xxx prefix.   \n" +
			"   This provider performs simple in memory AES encryption, and requires that you pass in either a 16 or 32 \n" +
			"   character key to indicate if you want 128 or 256 bit AES encryption. For example a key with value \n" +
			"   'AES256Key-32Characters0987654321' will result in your content being encrypted with 256 bit AES encryption. \n" +
			"   Note: If you lose access to this encryption key, you will NOT be able to decrypt these values!!\n\n" +

			"   Vault provider (https://www.vaultproject.io) \n" +
			"   -------------------------------------------- \n" +
			"   NOTE: Any arguments specifically related to the vault provider will always have a vault-xxx prefix.   \n" +
			"   This provider expects a running and accessible instance of Vault to be available and makes use of the \n" +
			"   standard Vault environment variables to gain access to it e.g.  \n" +
			"       export VAULT_TOKEN=your-vault-root-token \n" +
			"       export VAULT_ADDR=http://127.0.0.1:8200\n" +
			"       export VAULT_SKIP_VERIFY=true\n" +
			"   Vault's transit aka 'encryption as a service' feature is then used to offload and perform the actual \n" +
			"   encryption. (https://www.vaultproject.io/docs/secrets/transit). The Vault transit backend makes use \n" +
			"   of a registered named encryption key to gain access to the underlying encryption key itself, as well \n" +
			"   details about what encryption algorithm to use. This encrypt command expects you to have \n" +
			"   already setup and registered a named encryption key which will be used here. If not explicitly\n" +
			"   specified, then 'terrahelp' i.e. /transit/key/terrahelp is assumed as default. Note you can use the \n" +
			"   terrahelp tfstate vault-autoconfig command to auto register and generate a new key against this default \n" +
			"   named key for you if not already done. \n\n",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "provider",
				Value:       terrahelp.ThEncryptProviderSimple,
				EnvVar:      "TH_ENCRYPTION_PROVIDER",
				Usage:       "Encryption provider (simple|vault) to use (defaults to " + terrahelp.ThEncryptProviderSimple + ")",
				Destination: &ctxOpts.EncProvider,
			},
			cli.StringFlag{
				Name:        "mode",
				Value:       terrahelp.ThEncryptModeFull,
				EnvVar:      "TH_ENCRYPTION_MODE",
				Usage:       fmt.Sprintf("Encryption mode (inline|full) to use (defaults to full)"),
				Destination: &ctxOpts.EncMode,
			},
			cli.StringFlag{
				Name:        "state",
				Value:       terrahelp.TfstateFilename,
				Usage:       "Terraform state file to encrypt (defaults to " + terrahelp.TfstateFilename + ")",
				Destination: &ctxOpts.TfstateFile,
			},
			cli.StringFlag{
				Name:        "statebkp",
				Value:       terrahelp.TfstateBkpFilename,
				Usage:       "Backup terraform state file to encrypt (defaults to " + terrahelp.TfstateBkpFilename + ")",
				Destination: &ctxOpts.TfStateBkpFile,
			},
			cli.StringFlag{
				Name:        "tfvars",
				Value:       terrahelp.TfvarsFilename,
				Usage:       "Terraform tfvars filename (defaults to " + terrahelp.TfvarsFilename + ")",
				Destination: &ctxOpts.TfvarsFilename,
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups (defaults to " + terrahelp.ThBkpExtension + ")",
				Destination: &ctxOpts.BkpExt,
			},
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before encrypting (defaults to false)",
				Destination: &ctxOpts.NoBackup,
			},
			cli.BoolTFlag{
				Name:        "dblencrypt",
				Usage:       "Permits the double encryption of the content in a file (defaults to true)",
				Destination: &ctxOpts.AllowDoubleEncrypt,
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
				Usage:       "(Vault provider only) Named encryption key to use (defaults to " + terrahelp.ThNamedEncryptionKey + ")",
				Destination: &ctxOpts.NamedEncKey,
			},
		},
		Action: func(c *cli.Context) {
			th := f(ctxOpts.EncProvider)
			err := ctxOpts.ValidateForEncryptDecrypt()
			exitIfError(err)
			err = th.Encrypt(ctxOpts)
			exitIfError(err)
		},
	}
}

func decryptSubCommand(f func(provider string) *terrahelp.Tfstate) cli.Command {

	ctxOpts := &terrahelp.TfstateOpts{}

	return cli.Command{
		Name:  "decrypt",
		Usage: "Uses configured provider to decrypt local .tfstate files",
		Description: "Using either the 'simple' (default) or 'vault' provider, decrypt will ensure that the local \n" +
			"   terraform.tfstate files are decrypted and restored back to a state where terraform can operate \n" +
			"   on them again. You will have previously already encryped the files with either the 'simple' or \n" +
			"   'vault' provider, and using 'full' or 'inline' mode (see encrypt help for more info). \n" +
			"   This command essentially does the reverse of the encrypt command, and in general it is expected that: \n" +
			"   * The encrypted versions of the terraform.tfstate files will be obtained \n" +
			"   * The appropriate decryption applied i.e.\n" +
			"       terrahelp tfstate decrypt\n" +
			"   * The standard terraform commands can then be performed i.e. \n" +
			"       terraform plan\n" +
			"       terraform apply\n\n" +

			"   Note: You need to ensure you use the same 'provider' and 'mode' that was used when the files were \n" +
			"   encrypted in the first place. The 'provider' and well as 'mode' can be supplied as CLI arguments,  \n" +
			"   or via the TH_ENCRYPTION_PROVIDER and TH_ENCRYPTION_MODE environment variables. \n\n",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "provider",
				Value:       terrahelp.ThEncryptProviderSimple,
				EnvVar:      "TH_ENCRYPTION_PROVIDER",
				Usage:       "Encryption provider (simple|vault) to use (defaults to " + terrahelp.ThEncryptProviderSimple + ")",
				Destination: &ctxOpts.EncProvider,
			},
			cli.StringFlag{
				Name:        "mode",
				Value:       terrahelp.ThEncryptModeFull,
				EnvVar:      "TH_ENCRYPTION_MODE",
				Usage:       fmt.Sprintf("Encryption mode (inline|full) to use (defaults to full)"),
				Destination: &ctxOpts.EncMode,
			},
			cli.StringFlag{
				Name:        "state",
				Value:       terrahelp.TfstateFilename,
				Usage:       "Terraform state file to decrypt (defaults to " + terrahelp.TfstateFilename + ")",
				Destination: &ctxOpts.TfstateFile,
			},
			cli.StringFlag{
				Name:        "statebkp",
				Value:       terrahelp.TfstateBkpFilename,
				Usage:       "Backup terraform state file to decrypt (defaults to " + terrahelp.TfstateBkpFilename + ")",
				Destination: &ctxOpts.TfStateBkpFile,
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups (defaults to " + terrahelp.ThBkpExtension + ")",
				Destination: &ctxOpts.BkpExt,
			},
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before decrypting (defaults to false)",
				Destination: &ctxOpts.NoBackup,
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
				Usage:       "(Vault provider only) Named encryption key to use (defaults to " + terrahelp.ThNamedEncryptionKey + ")",
				Destination: &ctxOpts.NamedEncKey,
			},
		},
		Action: func(c *cli.Context) {
			th := f(ctxOpts.EncProvider)
			err := ctxOpts.ValidateForEncryptDecrypt()
			exitIfError(err)
			err = th.Decrypt(ctxOpts)
			exitIfError(err)
		},
	}

}

func vaultPrepSubCommand(f func(provider string) *terrahelp.Tfstate) cli.Command {

	ctxOpts := &terrahelp.TfstateOpts{}

	return cli.Command{
		Name:  "vault-autoconfig",
		Usage: "(Vault provider only) performs a very basic vault setup to allow Vault provider to be used out of the box.",
		Description: "This is really a one off helper command to help get off the ground and running quickly with Vault. \n" +
			"   Essentially it ensures that the transit backend is mounted (if not\n" +
			"   already), and that the named encryption key (defaults to /transit/key/terrahelp) is generated and registered \n" +
			"   The vault provider uses this named key as part of the encrypt and decrypt functionality. \n",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "namedkey",
				Value:       terrahelp.ThNamedEncryptionKey,
				EnvVar:      "TH_VAULT_NAMED_KEY",
				Usage:       "Named Vault encryption key to use (defaults to " + terrahelp.ThNamedEncryptionKey + ")",
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

// At present code required to add decent exit code support in the cli library
// is awaiting a 2.0. release (https://github.com/codegangsta/cli/pull/266)
// so until then we have to do a bit of an ugly emergency exit ourselves
func exitIfError(e error) {
	if e != nil {
		fmt.Printf("ERROR occured : %s\n", e)
		switch e.(type) {
		case *terrahelp.CryptoWrapError:
			os.Exit(cryptoWrapErrorExitCode)
		default:
			os.Exit(otherErrorExitCode)
		}
	}
}
