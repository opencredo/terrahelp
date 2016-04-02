## Terrahelp - terraforming, with a little help from your friends

terrahelp is as a command line utility written in [Go](https://github.com/golang/go) and is aimed at 
providing supplementary functionality which can sometimes prove useful when working with 
[Terraform](https://www.terraform.io). 


At present, the functionality offered includes:

* _Encryption & decryption of terraform state files_.
This can be done in either full or inline mode, and provides the ability to leverage either a simple or vault based encryption provider. 
For more details and an example of how to use it please see [the example README](https://github.com/opencredo/terrahelp/tree/master/examples/tfstate-encrypt). 
Additionally this blog post on [Securing Terraform State with Vault](https://www.opencredo.com/securing-terraform-state-with-vault)
provides more details and background as well.

        NAME:
           terrahelp tfstate - Options for performing actions on the local tfstate files.
        
        USAGE:
           terrahelp tfstate command [command options] [arguments...]
        
        COMMANDS:
            vault-autoconfig	(Vault provider only) performs a very basic vault setup to allow Vault 
                                provider to be used out of the box.
            encrypt		        Uses configured provider to encrypt local .tfstate files
            decrypt		        Uses configured provider to decrypt local .tfstate files
        
        OPTIONS:
           --help, -h	show help

## Installation

### Pre-built binary

An initial pre-built terrahelp binary can be found [here](https://github.com/opencredo/terrahelp/releases/).  

#### OSX, Linux & *BSD

Download a binary, set the correct permissions, add to your PATH:

    chmod +x terrahelp
    export PATH=$PATH:/wherever/terrahelp

And run it:

    terrahelp -help

#### Windows

Still coming ...

### Build it yourself  

To set up your Go environment - look [here](https://golang.org/doc/code.html).

You must have Go > 1.6 installed.

    mkdir -p "$GOPATH/src/github.com/opencredo/"
    git clone https://github.com/opencredo/terrahelp.git "$GOPATH/src/github.com/opencredo/terrahelp"
    cd "$GOPATH/src/github.com/opencredo/terrahelp"
    make build

And to run terrahelp:

    ./terrahelp
