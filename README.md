## Terrahelp - terraforming, with a little help from your friends

terrahelp is as a command line utility written in [Go](https://github.com/golang/go). 
It aims to provide you with a utility to help fill in some of the gaps, and perform some additional tasks 
which you always seem to find yourself doing when working with [Terraform](https://www.terraform.io). 

This is an ongoing project, and at present the following capabilities are provided:

* Encryption and decryption of local terraform state files (using Vault's "encryption as a service" functionality)

The above encryption can either be done against the entire file, or inline, i.e. only the sensitive data within
the tfstate file is encrypted.

The following blog post also provides an intial overview of how the encryption
functionality can be used.

## Installation

### Pre-built binary

An initial pre-built terrahelp binary can be found [here](https://github.com/opencredo/terrahelp/releases/).  

#### OSX, Linux & *BSD

Download a binary, set the correct permissions:

    chmod +x terrahelp

And run it:

    ./terrahelp

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
