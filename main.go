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
	app.Version = "0.7.6-dev"
	app.Author = "https://github.com/opencredo OpenCredo - Nicki Watt"
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
