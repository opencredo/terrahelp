
## 0.3.0 (Unreleased)

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
