[![Travis CI][Travis-Image]][Travis-Url]
[![ReportCard][ReportCard-Image]][ReportCard-Url]

# Terrahelp
##### Terraforming, with a little help from your friends

Terrahelp is as a command line utility written in [Go](https://golang.org) and is aimed at
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
           X.X.X

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

### macOS

Install using [Homebrew](https://brew.sh/):

    brew install terrahelp
    terrahelp -v

### Manual Installation Using the Pre-Built Binaries

Available from the Terrahelp repository's [releases page](https://github.com/opencredo/terrahelp/releases)

The community has also made it available as a [Terrahelp AUR package](https://aur.archlinux.org/packages/terrahelp)

#### macOS, Linux & *BSD

Download a binary, set the correct permissions, add to your PATH:

    chmod +x terrahelp
    export PATH=$PATH:/wherever/terrahelp

And run it:

    terrahelp --help

##### macOS Additional Step

`terrahelp` may be prevented from running if you downloaded it using a web browser. To fix this, remove the quarantine attribute before running again:

    xattr -d com.apple.quarantine terrahelp

#### Windows

Not yet supported

## Build from source

### Prerequisites

Install Go (Terrahelp is currently built against 1.13.x).  The following official resources will guide you through your environment setup.

* [Getting Started](https://golang.org/doc/install)
* [Go Documentation](https://golang.org/doc)

Clone the Terrahelp repository.

```bash
mkdir -p "$GOPATH/src/github.com/opencredo/"
git clone https://github.com/opencredo/terrahelp.git "$GOPATH/src/github.com/opencredo/terrahelp"
cd "$GOPATH/src/github.com/opencredo/terrahelp"
```

### Dependencies

Terrahelp uses Go modules to manage it's dependencies.  During Go's transition to switching on modules by default, Terrahelp is setup to buildusing the vendor directory.
Supportive targets are prvoided to allow the vendor directory to be recreated if required.

### Building and Executing

After a build has completed successfully a binary will be built and placed into a local bin directory.  The following commands build and execute terrahelp.

    make build
    ./bin/terrahelp -v

### Testing

    make test

### Installing and Executing

Installation places the binary in the `$GOPATH/bin` directory. Assuming that the directory has been added to your `PATH`, the following commands will install and execute Terrahelp.

    make install
    terrahelp -v

### Want to cross compile it?

The make file allows both OSX and Linux binaries to be created at the same time or individually.
The following commands show joint creation followed by OSX, (darwin) then Linux creation.  All cross compiled binaries will be placed in a `dist` directory.

    make dist
    make darwin
    make linux

### Clean your project

A number of work directories will have been created through the previous build steps. The local `bin` and `dist` directories will contain binaries.
The following command can be used to return the project back to a pre build state.

    make clean

### Dependency management

The following targets have been created to allow dependencies to be managed through Go modules.  As mentioned before Terrahelp builds using the vendor directory.

* `make dependencies`
  * Downloads the dependencies to the Go modules cache.
* `make tidy-dependencies`
  * Adds missing and removes unused modules.
* `make vendor-dependencies`
  * Copies the dependencies into the local vendor directory.
* `make clean-dependencies`
  * Removes the local vendor directory.

**NOTE:**  The Makefile defines a variable called `BUILDARGS` and this is currently set with `-mod=vendor`.  This instructs various go commands to use the vendor directory.  This can be overridden to build to project using standard go module flows.

    BUILDARGS='' make build

## Releasing

### Brew
***NOTE:*** This step should be performed *after* a new version of `terrahelp` has been released.

Follow the instructions outlined in [Submit a new version of an existing formula][Homebrew-Update-Formula] to update the version of `terrahelp` installed by Brew.

For reference, the formula can be viewed in the homebrew-core repository [here][Terrahelp-Formula].

[Travis-Image]: https://travis-ci.org/opencredo/terrahelp.svg?branch=master
[Travis-Url]: https://travis-ci.org/opencredo/terrahelp
[ReportCard-Url]: http://goreportcard.com/report/opencredo/terrahelp
[ReportCard-Image]: http://goreportcard.com/badge/opencredo/terrahelp
[Homebrew-Update-Formula]: https://docs.brew.sh/How-To-Open-a-Homebrew-Pull-Request#submit-a-new-version-of-an-existing-formula
[Terrahelp-Formula]: https://github.com/Homebrew/homebrew-core/blob/master/Formula/terrahelp.rb
