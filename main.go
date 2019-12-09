package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/opencredo/terrahelp/terrahelp"
)

var (
	name    = "terrahelp"
	usage   = "Provides additional functions helpful with terraform development"
	version = "0.7.2-dev"
	author  = "https://github.com/opencredo OpenCredo - Nicki Watt"
	commit  string
)

func main() {

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "GitVersion=%s\nGitCommit=%s\n", c.App.Version, commit)
	}

	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version
	app.Author = author
	app.Commands = []cli.Command{
		vaultAutoConfigCommand(newTerraHelperFunc()),
		encryptCommand(newTerraHelperFunc()),
		decryptCommand(newTerraHelperFunc()),
		maskCommand(),
	}
	app.Run(os.Args)
}

func newTerraHelperFunc() func(provider string) *terrahelp.CryptoHandler {
	return func(provider string) *terrahelp.CryptoHandler {

		switch {
		case (provider == terrahelp.ThEncryptProviderSimple):
			e := terrahelp.NewSimpleEncrypter()
			return &terrahelp.CryptoHandler{Encrypter: e}
		case (provider == terrahelp.ThEncryptProviderVault):
			e, err := terrahelp.NewVaultEncrypter()
			if err != nil {
				exitIfError(err)
			}
			return &terrahelp.CryptoHandler{Encrypter: e}
		case (provider == terrahelp.ThEncryptProviderVaultCli):
			e, err := terrahelp.NewVaultCliEncrypter()
			if err != nil {
				exitIfError(err)
			}
			return &terrahelp.CryptoHandler{Encrypter: e}
		}

		exitIfError(fmt.Errorf("Invalid provider %s specified ", provider))
		return nil
	}
}
