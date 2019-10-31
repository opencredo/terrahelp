[![Travis CI][Travis-Image]][Travis-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

# Terrahelp
##### terraforming, with a little help from your friends

`terrahelp` is as a command line utility written in [Go](https://github.com/golang/go) and is aimed at 
providing supplementary functionality which can sometimes prove useful when working with 
[Terraform](https://www.terraform.io). 

At present, it offers the following:

* _Encryption & decryption functionality_.
Run in either full or inline mode, and leveraging either a simple or [Vault](https://www.vaultproject.io) based encryption provider, this
functionality provides the ability to encrypt and decrypt files such as terraform.tfstate files, as well as piped in 
output from commands such as terraform apply etc. 

* _Masking functionality_.
If you don't want to encrypt sensitive data, but rather just mask it out with something like ***** then you can use
the mask command instead. This can either be run over a file, or have the content piped into it.

For more details, and some examples of how to use it please see [the example READMEs](https://github.com/opencredo/terrahelp/tree/master/examples). 

Additionally the blog post [Securing Terraform State with Vault](https://www.opencredo.com/securing-terraform-state-with-vault) also provides more details and background as well.

        NAME:
           terrahelp - Provides additional functions helpful with terraform development

        USAGE:
           terrahelp [global options] command [command options] [arguments...]

        VERSION:
           0.4.3

        AUTHOR(S):
           https://github.com/opencredo OpenCredo - Nicki Watt

        COMMANDS:
            vault-autoconfig	Auto configures Vault with a basic setup to support encrypt and decrypt actions.
            encrypt		        Uses configured provider to encrypt specified content
            decrypt		        Uses configured provider to decrypt specified content
            mask                    Mask will overwrite sensitive data in output or files with a masked value (eg. ******).
            help, h                 Shows a list of commands or help for one command
            
        GLOBAL OPTIONS:
           --help, -h		show help
           --version, -v	print the version


## Installation

### Pre-built binaries

Available from this repository's releases page [here](https://github.com/opencredo/terrahelp/releases/).

The community has also made it available as an AUR package via https://aur.archlinux.org/packages/terrahelp 

#### OSX, Linux & *BSD

Download a binary, set the correct permissions, add to your PATH:

    chmod +x terrahelp
    export PATH=$PATH:/wherever/terrahelp

And run it:

    terrahelp -help

##### OSX Additional Step

`terrahelp` may be prevented from running if you downloaded it using a web browser. To fix this, remove the quarantine attribute before running again:
 
    xattr -d com.apple.quarantine terrahelp

#### Windows

Not here yet ...

### Build it yourself  

To set up your Go environment - look [here](https://golang.org/doc/code.html).

Install Go (Terrahelp is currently built against 1.7.3)

    mkdir -p "$GOPATH/src/github.com/opencredo/"
    git clone https://github.com/opencredo/terrahelp.git "$GOPATH/src/github.com/opencredo/terrahelp"
    cd "$GOPATH/src/github.com/opencredo/terrahelp"

*Dependencies*

Terrahelp uses [govendor](https://github.com/kardianos/govendor) to manage it's dependencies, and currently also checks them into Git to enable a seamless build experience. However should you wish to change / upgrade any of these, you can get govendor, and then run the appropriate commands (e.g sync, fetch etc)

    go get -u github.com/kardianos/govendor
    govendor sync

*Build it*

    go install
    
*Test it*
    
    go test -v ./...

*Run it:*

    terrahelp -v 
    
*Want to cross compile it?*

    env GOOS=darwin GOARCH=amd64 go build -o=terrahelp-darwin-amd64
    env GOOS=linux GOARCH=amd64 go build -o=terrahelp-linux-amd64

[Travis-Image]: https://travis-ci.org/opencredo/terrahelp.svg?branch=master
[Travis-Url]: https://travis-ci.org/opencredo/terrahelp
[ReportCard-Url]: http://goreportcard.com/report/opencredo/terrahelp
[ReportCard-Image]: http://goreportcard.com/badge/opencredo/terrahelp
