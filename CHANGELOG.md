## 0.7.6 (unreleased)

## 0.7.5 (2021-10-04)
* [PR-37](https://github.com/opencredo/terrahelp/pull/37) Update Terrahelp build pipeline to user GitHub Actions, (includes update to go 1.17))
* [PR-36](https://github.com/opencredo/terrahelp/pull/36) Remove quoted types in variables for tf > 0.11 support
* [PR-32](https://github.com/opencredo/terrahelp/pull/32) Add Brew installation instructions

## 0.7.4
* [PR-31](https://github.com/opencredo/terrahelp/pull/31) Introduces sha256sum file covering each distribution binary (resolves [#28](https://github.com/opencredo/terrahelp/issues/28))

## 0.7.3
* [PR-30](https://github.com/opencredo/terrahelp/pull/30) Updates Travis CI file to replace missed TRAGET reference with NAME to allow release uploads (resolves [#29](https://github.com/opencredo/terrahelp/issues/29))

## 0.7.2
* [PR-23](https://github.com/opencredo/terrahelp/pull/23) Update Terrahelp to process HCL2 syntax, (including tests and examples) (resolves [#22](https://github.com/opencredo/terrahelp/issues/22))
* [PR-25](https://github.com/opencredo/terrahelp/pull/25) Updates Makefile and README.md to introduce new targets (resolves [#24](https://github.com/opencredo/terrahelp/issues/24))

## 0.7.1
* [PR-21](https://github.com/opencredo/terrahelp/pull/21) Updates Travis deploy credentials for Github releases (resolved [#20](https://github.com/opencredo/terrahelp/issues/20))

## 0.7.0
* [PR-19](https://github.com/opencredo/terrahelp/pull/19) Build against Go 1.13 and manages dependencies through GO Modules (resolved [#18](https://github.com/opencredo/terrahelp/issues/18))
* [PR-19](https://github.com/opencredo/terrahelp/pull/19) Builds against Vault API v1.0.4 (resolves [#16](https://github.com/opencredo/terrahelp/issues/16))
* [PR-19](https://github.com/opencredo/terrahelp/pull/19) Makefile targets expanded to simplify Travis file.

## 0.4.3
* [PR-12](https://github.com/opencredo/terrahelp/pull/12) Updated to be compatible with Vault 0.6.2
* [PR-11](https://github.com/opencredo/terrahelp/pull/11) Exclude empty strings from detection, and provide config flag for handling whitespace only values (resolves [#10](https://github.com/opencredo/terrahelp/issues/10))

## 0.4.2
* [PR-9](https://github.com/opencredo/terrahelp/pull/9) Cater for terraform 0.7.x list and map variables (resolves [#8](https://github.com/opencredo/terrahelp/issues/8))
* Updated examples and command line docs
* Builds against Go 1.7.3
* Confirmed testing against Vault 0.5.2

## 0.4.1
* Add new `vault-cli` provider to use the `vault` command line tool rather than talking to the vault API.

## 0.4.0
**Note: This release contains breaking changes!!**

* [PR-6](https://github.com/opencredo/terrahelp/pull/6) Terrahelp will now ignore stdin input if a `-file` flag is present. The `-file` flag will no longer default to terraform.tfstate and terraform.tfstate.backup (part of resolving [#5](https://github.com/opencredo/terrahelp/issues/5))

## 0.3.1

FEATURES:

* **mask command**: Provide ability to mask sensitive input from terraform commands

## 0.3.0

**Note: This release contains breaking changes!!** 

The core functionality introduced is to expand the encryption/decryption functionality to 
be used on more than just the terraform .tfstate files (pipes and alternate files). [PR-3](https://github.com/opencredo/terrahelp/pull/3)
      
In version 0.2.1 and earlier, the command to apply these crypto functions was previously
exposed as subcommands under the main `tfstate` command, i.e. 
`terrahelp tfstate encrypt -mode=xx`

Moving forward, these commands have been promoted as top level commands i.e.

        terrahelp encrypt -mode=xx 
        terrahelp decrypt -mode=xx 
  

## 0.2.1

IMPROVEMENTS:

* Add flag `-dblencrypt` to control whether double encryption is allowed [PR-2](https://github.com/opencredo/terrahelp/pull/2)

 
## 0.2.0

* First public release of terrahelp  
