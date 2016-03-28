package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/opencredo/terrahelp/terrahelp"
)

func main() {

	app := cli.NewApp()
	app.Name = "terrahelp"
	app.Usage = "Provides additional functions helpful with terraform development"
	app.Version = "0.1.2"
	app.Author = "https://github.com/opencredo (Nicki Watt)"
	app.Commands = []cli.Command{
		{
			Name:  "tfstate",
			Usage: "Options for performing actions on the local tfstate files.",
			Subcommands: []cli.Command{
				vaultPrepSubCommand(newTerraHelperFunc()),
				encryptSubCommand(newTerraHelperFunc()),
				decryptSubCommand(newTerraHelperFunc()),
			},
		},
	}
	app.Run(os.Args)
}

func newTerraHelperFunc() func(provider string) *terrahelp.Tfstate {
	return func(provider string) *terrahelp.Tfstate {

		switch {
		case (provider == terrahelp.ThEncryptProviderSimple):
			e := terrahelp.NewSimpleEncrypter()
			return &terrahelp.Tfstate{Encrypter: e}
		case (provider == terrahelp.ThEncryptProviderVault):
			e, err := terrahelp.NewVaultEncrypter()
			if err != nil {
				exitIfError(err)
			}
			return &terrahelp.Tfstate{Encrypter: e}
		}

		exitIfError(fmt.Errorf("Invalid provider %s specified ", provider))
		return nil
	}
}
