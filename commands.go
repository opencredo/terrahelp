package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/opencredo/terrahelp/terrahelp"
)

const errorExitCode = 1

func encryptSubCommand(f func(provider string) *terrahelp.Tfstate) cli.Command {

	ctxOpts := &terrahelp.TfstateOpts{}

	return cli.Command{
		Name:  "encrypt",
		Usage: "Uses configured provider to encrypt the local .tfstate files",
		Description: "Encrypt currently supports the 'simple' and 'vault' providers to do the encryption / decryption. \n" +
			"Encryption is performed by default on the local terraform.tfstate and its backup file \n" +
			"   terraform.tfstate.backup. The default is to encrypt the entire tfstate files themselves, however  \n" +
			"   by setting inline=true, encryption can be done inline where only sensitive variables within the  \n" +
			"   files are replaced as opposed to the entire file. \n" +
			"   In both cases, the approach is generally to do everything normally as before, i.e. terraform plan\n" +
			"   and terraform apply, and then when happy, to apply the encryption and checkin to version control \n" +
			"   after that if required. Encrypted values will always conform to the following format: \n" +
			"   @terrahelp-encrypted(ENCRYPTED_CONTENT). Decryption essentially operates \n" +
			"   in reverse. \n\n" +

			"   Vault provider (provider=vault) \n" +
			"   -------------------------------  \n" +
			"   NOTE: Any arguments specifically related to the vault provider will always have a vault-xxx prefix.   \n" +
			"   This provider expects a running and accessible version of Vault to be available and makes use of the \n" +
			"   standard Vault environment variables i.e. VAULT_URL and VAULT_TOKEN to configure access to it. \n" +
			"   The transit aka 'encryption as a service' feature is then used to offload and perform the actual \n" +
			"   encryption. (https://www.vaultproject.io/docs/secrets/transit). This command expects you to have \n" +
			"   already setup and registered a named encryption key which will be used here. If not explicitly\n" +
			"   specified, then 'terrahelp' i.e. /transit/key/terrahelp is used as default. Note you can use the \n" +
			"   terrahelp vault-autoconfig command to register this default automatically for you if not already done. \n\n" +

			"   Simple provider (provider=simple) \n" +
			"   --------------------------------  \n" +
			"   NOTE: Any arguments specifically related to the simple provider will always have a simple-xxx prefix.   \n" +
			// The key argument should be the AES key, either 16 or 32 bytes
			// to select AES-128 or AES-256.
			// "AES256Key-32Characters0987654321"
			"   Uses aes-gcm  \n" +
			"   Note: If you lose access to this named encryption key, you will NOT be able to encrypt these values!!\n\n" +

			"   Inline encryption  \n" +
			"   -----------------  \n" +
			"   Inline encrypting works by detected sensitive values within the state files and then replaces \n" +
			"   them with encrypted values. It is assumed that all sensitive values are defined as variables \n" +
			"   within the terraform.tfvars file, which just by way of recap should never be checked into version \n" +
			"   control! terrahelp then uses terraform.tfvars to search for these values within the tfstate files \n" +
			"   and essentially does a find and replace of all the sensitive values using encrypted values. These \n" +
			"   state files can then safely be checked into version control, and when required, can have the \n" +
			"   decrypt command run over them to restore them back to their original state. ",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "provider",
				Value:       terrahelp.ThEncryptProviderSimple,
				Usage:       "encryption provider to use (defaults to " + terrahelp.ThEncryptProviderSimple + ")",
				Destination: &ctxOpts.EncProvider,
			},
			cli.BoolFlag{
				Name:        "inline",
				Usage:       fmt.Sprintf("Only encrypts detected sensitive fields in the tfstate files. (defaults to false)"),
				Destination: &ctxOpts.Inline,
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
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before encrypting (defaults to false)",
				Destination: &ctxOpts.NoBackup,
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups (defaults to " + terrahelp.ThBkpExtension + ")",
				Destination: &ctxOpts.BkpExt,
			},
			cli.StringFlag{
				Name:        "simple-key",
				Usage:       "(Simple provider only) the encryption key to use",
				Destination: &ctxOpts.SimpleKey,
			},
			cli.StringFlag{
				Name:        "vault-namedkey",
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
		Usage: "Uses configured provider to decrypt the local .tfstate files",
		Description: "Decrypt currently supports the 'simple' and 'vault' providers to do the encryption / decryption. \n" +
			"   Decryption is performed by default on the local terraform.tfstate and its backup file \n" +
			"   terraform.tfstate.backup. Unless overriden terrahelp will decrypt the entire tfstate files themselves, however  \n" +
			"   by setting inline=true, decryption can be configured to seek out embedded encrypted values within \n" +
			"   within the state files and then replace them inline with their decrypted versions. \n" +
			"   In both cases, the approach is generally to perform the decryption first, and then to simply do \n" +
			"   everything normally as before, i.e. terraform plan and terraform apply. Encrypted values will always\n" +
			"   conform to the following format: @terrahelp-encrypted(ENCRYPTED_CONTENT). \n" +
			"   Encryption essentially operates in reverse. \n\n" +

			"   Vault provider (provider=vault) \n" +
			"   -------------------------------  \n" +
			"   NOTE: Any arguments specifically related to the vault provider will always have a vault-xxx prefix.   \n" +
			"   This provider expects a running and accessible version of Vault to be available and makes use of the \n" +
			"   standard Vault environment variables i.e. VAULT_URL and VAULT_TOKEN to configure access to it. \n" +
			"   The transit aka 'encryption as a service' feature is then used to offload and perform the actual \n" +
			"   decryption. (https://www.vaultproject.io/docs/secrets/transit). This command expects you to have \n" +
			"   already setup and registered a named encryption key. This will be the same named key that will \n" +
			"   have been previously used by the encrypt command to do the encrypting in the first place.\n" +
			"   If not specified then 'terrahelp' i.e. /transit/key/terrahelp will be used as the default. \n" +
			"   Note: If you lose access to this named encryption key, you will NOT be able to decrypt these values!!\n\n" +

			"   Simple provider (provider=simple) \n" +
			"   --------------------------------  \n" +
			"   NOTE: Any arguments specifically related to the simple provider will always have a simple-xxx prefix.   \n" +
			"   Uses aes-gcm  \n" +
			"   Note: If you lose access to this named encryption key, you will NOT be able to decrypt these values!!\n\n" +

			"   Inline decryption  \n" +
			"   -----------------  \n" +
			"   Inline decrypting works by detecting encrypted values within the state files and then replacing \n" +
			"   them with their decrypted values. ",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "provider",
				Value:       terrahelp.ThEncryptProviderSimple,
				Usage:       "encryption provider to use (defaults to " + terrahelp.ThEncryptProviderSimple + ")",
				Destination: &ctxOpts.EncProvider,
			},
			cli.BoolFlag{
				Name:        "inline",
				Usage:       fmt.Sprintf("Only decrypts detected encrypted fields in the tfstate files. (defaults to false)"),
				Destination: &ctxOpts.Inline,
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
			cli.BoolFlag{
				Name:        "nobackup",
				Usage:       "Suppress the creation of backup files before decrypting (defaults to false)",
				Destination: &ctxOpts.NoBackup,
			},
			cli.StringFlag{
				Name:        "bkpext",
				Value:       terrahelp.ThBkpExtension,
				Usage:       "Extension to use when creating backups (defaults to " + terrahelp.ThBkpExtension + ")",
				Destination: &ctxOpts.BkpExt,
			},
			cli.StringFlag{
				Name:        "vault-namedkey",
				Value:       terrahelp.ThNamedEncryptionKey,
				Usage:       "(Vault provider only) Named encryption key to use (defaults to " + terrahelp.ThNamedEncryptionKey + ")",
				Destination: &ctxOpts.NamedEncKey,
			},
			cli.StringFlag{
				Name:        "simple-key",
				Usage:       "(Simple provider only) the encryption key to use",
				Destination: &ctxOpts.SimpleKey,
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
		os.Exit(errorExitCode)
	}
}
