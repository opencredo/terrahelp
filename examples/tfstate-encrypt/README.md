## Terrahelp Example - Encrypt / decrypt

This example contains a very simple terraform setup composed entirely of 
local resources (e.g. template resource) and exists in order to demonstrate how
you can do basic encryption and decryption functionality in the absence of a 
formal solution (ref https://github.com/hashicorp/terraform/issues/516).
 
This example is completely safe to run and will not land up costing you any 
money in a cloud provider!
 
The CLI itself offers a more comprehensive view of the various options available, 
so please use this if you need more info.
Additionally you can read this corresponding blog which gives a more detailed explanation
of this functionality and its usage: [Securing Terraform State with Vault](https://www.opencredo.com/securing-terraform-state-with-vault).

### Simple inline encryption of terraform output

This example will demonstrate _inline_ encryption and decryption using the _simple_ encryption provider
where we will pipe the content in directly. This specific example uses the basic command line arguments
as opposed to environment variables to control the process

* Run a `terraform plan` as normal

        terraform plan
        
* Inspect the result which should look something like below:        
        
        Refreshing Terraform state prior to plan...
        
        The Terraform execution plan has been generated and is shown below.
        Resources are shown in alphabetical order for quick scanning. Green resources
        will be created (or destroyed and then created if an existing resource
        exists), yellow resources are being changed in-place, and red resources
        will be destroyed.
        
        Note: You didn't specify an "-out" parameter to save this plan, so when
        "apply" is called, Terraform can't guarantee this is what will execute.
        
        + template_file.example
            rendered:  "" => "<computed>"
            template:  "" => "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}"
            vars.#:    "" => "3"
            vars.msg1: "" => "sensitive-value-1-AK#%DJGHS*G"
            vars.msg2: "" => "normal value 1"
            vars.msg3: "" => "sensitive-value-3-//dfhs//"
        
        
        Plan: 1 to add, 0 to change, 0 to destroy.

* Run the same command, but pipe the output into the terrahelp encrypt command. 

        terraform plan | terrahelp encrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline 

* The result should now look something like below:

        Refreshing Terraform state prior to plan...
        
        
        The Terraform execution plan has been generated and is shown below.
        Resources are shown in alphabetical order for quick scanning. Green resources
        will be created (or destroyed and then created if an existing resource
        exists), yellow resources are being changed in-place, and red resources
        will be destroyed.
        
        Note: You didn't specify an "-out" parameter to save this plan, so when
        "apply" is called, Terraform can't guarantee this is what will execute.
        
        + template_file.example
            rendered:  "" => "<computed>"
            template:  "" => "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}"
            vars.#:    "" => "3"
            vars.msg1: "" => "@terrahelp-encrypted(Qwkbpytfx0br0Wh7fH6tfA6T1H+0cWnUKSkYnYrn8yd9lmZUt1vMp90uojPe)"
            vars.msg2: "" => "normal value 1"
            vars.msg3: "" => "@terrahelp-encrypted(FulrH6WwiDTNqPtVaWtHtJip+igD2hUSxWsabShX4NrF6tuKARD/R5fG)"
        
        
        Plan: 1 to add, 0 to change, 0 to destroy.

* For decryption, you could pipe the output again into the decrypt command, however more than
  likely, you will probably want to save the results into a file and then decrypt that. The
  sequence of commands to do that would be something as follows:
  
        terraform plan -out=my-infra.plan
        
        terrahelp encrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline -file=my-infra.plan
        
        terrahelp decrypt -simple-key=AES256Key-32Characters0987654321 -mode=inline -file=my-infra.plan        

### Simple inline encryption of tfstate files

This example will demonstrate _inline_ encryption and decryption using the _simple_ encryption provider
and will use explicit command line arguments (an example using environment variables is shown with the
Vault provider example).

* Run terraform as normal

        terraform plan
        terraform apply

* Verify `terraform.tfstate` contents before encryption (e.g. by doing a `cat terraform.tfstate`).
This should look something like below:
    
        {
            "version": 1,
            "serial": 1,
            "modules": [
                {
                    "path": [
                        "root"
                    ],
                    "outputs": {
                        "normal_val_2": "normal value 2",
                        "rendered": "\nmsg1 = sensitive-value-1-AK#%DJGHS*G\nmsg2 = normal value 1\nmsg3 = sensitive-value-3-//dfhs//"
                    },
                    "resources": {
                        "template_file.example": {
                            "type": "template_file",
                            "primary": {
                                "id": "b2cc7afb65fe7b6ac21328905d82e28fcbcdad1992cefce82cfa91691af24b91",
                                "attributes": {
                                    "id": "b2cc7afb65fe7b6ac21328905d82e28fcbcdad1992cefce82cfa91691af24b91",
                                    "rendered": "\nmsg1 = sensitive-value-1-AK#%DJGHS*G\nmsg2 = normal value 1\nmsg3 = sensitive-value-3-//dfhs//",
                                    "template": "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}",
                                    "vars.#": "3",
                                    "vars.msg1": "sensitive-value-1-AK#%DJGHS*G",
                                    "vars.msg2": "normal value 1",
                                    "vars.msg3": "sensitive-value-3-//dfhs//"
                                }
                            }
                        }
                    }
                }
            ]
        }


* Encrypt

        terrahelp encrypt -inline=true -simple-key="AES256Key-32Characters0987654321" 

* Inspect `terraform.tfstate` content after encryption. Note how all the sensitive values, as 
detected in the `terraform.tfvars` file, have now been replaced with encrypted versions. The
content should look something like that below: 

        {
            "version": 1,
            "serial": 1,
            "modules": [
                {
                    "path": [
                        "root"
                    ],
                    "outputs": {
                        "normal_val_2": "normal value 2",
                        "rendered": "\nmsg1 = @terrahelp-encrypted(43ZtxgMU7gxF5ZaV171iVypFe+Pam1Oev7TNCfklw2KZ2KBE6TJBiPpErYfB)\nmsg2 = normal value 1\nmsg3 = @terrahelp-encrypted(QJCFNaHuas+ZEFEI99qi9tp4z5MEZIVOFTcBCwCzbMj70vtXoO757KDd)"
                    },
                    "resources": {
                        "template_file.example": {
                            "type": "template_file",
                            "primary": {
                                "id": "b2cc7afb65fe7b6ac21328905d82e28fcbcdad1992cefce82cfa91691af24b91",
                                "attributes": {
                                    "id": "b2cc7afb65fe7b6ac21328905d82e28fcbcdad1992cefce82cfa91691af24b91",
                                    "rendered": "\nmsg1 = @terrahelp-encrypted(43ZtxgMU7gxF5ZaV171iVypFe+Pam1Oev7TNCfklw2KZ2KBE6TJBiPpErYfB)\nmsg2 = normal value 1\nmsg3 = @terrahelp-encrypted(QJCFNaHuas+ZEFEI99qi9tp4z5MEZIVOFTcBCwCzbMj70vtXoO757KDd)",
                                    "template": "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}",
                                    "vars.#": "3",
                                    "vars.msg1": "@terrahelp-encrypted(43ZtxgMU7gxF5ZaV171iVypFe+Pam1Oev7TNCfklw2KZ2KBE6TJBiPpErYfB)",
                                    "vars.msg2": "normal value 1",
                                    "vars.msg3": "@terrahelp-encrypted(QJCFNaHuas+ZEFEI99qi9tp4z5MEZIVOFTcBCwCzbMj70vtXoO757KDd)"
                                }
                            }
                        }
                    }
                }
            ]
        }

* To get your normal `terraform.tfstate` content back, decrypt

        terrahelp decrypt -inline=true -simple-key="AES256Key-32Characters0987654321" 

* Again verify `terraform.tfstate` content after decryption. This should now look exactly the same
as it did before doing the encryption


### Vault full encryption of tfstate files

This example will demonstrate _full_ encryption and decryption using the _vault_ encryption provider
and will use environment variables rather than explicit command line arguments to control the process.

* First, ensure you have a running Vault server available. You can quite easily download the latest version from 
here, then open up a new terminal, and for experimentation purposes, simply run the server in dev mode i.e.

        vault server -dev -dev-root-token-id="terrahelp-devonly-vault-root-token"

* In a separate terminal, change into this example project folder, and setup the necessary environment
  variables required for us to talk to our dev Vault server, as well as run the next set of terrahelp
  commands. Specifically will also run the `vault-autoconfig` command to configure Vault with the
  named encryption key we wnat to use. i.e.

        export VAULT_TOKEN="terrahelp-devonly-vault-root-token"
        export VAULT_ADDR="http://127.0.0.1:8200"
        export VAULT_SKIP_VERIFY="true"
        
        export TH_ENCRYPTION_PROVIDER="vault"
        export TH_ENCRYPTION_MODE="full"
        export TH_VAULT_NAMED_KEY="examplekey"
        terrahelp vault-autoconfig

* Run terraform as normal and inspect the terraform.tfstate content before encryption is applied

        terraform plan
        terraform apply

* Verify `terraform.tfstate` contents before encryption (e.g. by doing a `cat terraform.tfstate`).
This should look something like below:
    
        {
            "version": 1,
            "serial": 1,
            "modules": [
                {
                    "path": [
                        "root"
                    ],
                    "outputs": {
                        "normal_val_2": "normal value 2",
                        "rendered": "\nmsg1 = sensitive-value-1-AK#%DJGHS*G\nmsg2 = normal value 1\nmsg3 = sensitive-value-3-//dfhs//"
                    },
                    "resources": {
                        "template_file.example": {
                            "type": "template_file",
                            "primary": {
                                "id": "b2cc7afb65fe7b6ac21328905d82e28fcbcdad1992cefce82cfa91691af24b91",
                                "attributes": {
                                    "id": "b2cc7afb65fe7b6ac21328905d82e28fcbcdad1992cefce82cfa91691af24b91",
                                    "rendered": "\nmsg1 = sensitive-value-1-AK#%DJGHS*G\nmsg2 = normal value 1\nmsg3 = sensitive-value-3-//dfhs//",
                                    "template": "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}",
                                    "vars.#": "3",
                                    "vars.msg1": "sensitive-value-1-AK#%DJGHS*G",
                                    "vars.msg2": "normal value 1",
                                    "vars.msg3": "sensitive-value-3-//dfhs//"
                                }
                            }
                        }
                    }
                }
            ]
        }


* Encrypt

         terrahelp encrypt  

* Inspect `terraform.tfstate` content after encryption. Note how all the sensitive values, as 
detected in the terraform.tfvars file, have now been replaced with encrypted versions, and will
look something like below: 

        @terrahelp-encrypted(vault:v1:h7Yx1VAYvd2pyW0dd/iWifSe6yFB8QI7Zv2KjlW5USM5AyT9o3g3U2bU3
        vbDweRCGUXq2P8qpNcp8LUXDUon2Q6ee8I20X6yJyj5I2AS9V9ec4YcFOS9odqG+6dFqdlgWUkvEXPsH6puL0rX
        depvR17dvK1QTID0iE14HS7b4UnwI0Ti+f2VX4GvKHhnfKwCejKVu3g2bXdjn35h+EH9cHonSTx24SI6mM5k9Uy
        L7ht7AfPtPkdiUW7XSiW69UsZ+ZWrz8  ...  Ej3NYiY71Z/B2Rfm3M3V22BjfCsoUAHR1gL8acb5xQryuk+B/
        zQdLx7fXgxS8rMPKFwrJVRVtdcJtLFtLLf42AV1oUCqYvvusyNiGkQ6p3/2cgbkWsm/gN2lc26AuD6wVtd44qi
        CKK5iBZU4HQH6P5dycL0Sjgg4vJvcve85fQOLtfrr+UnQP0hdTSfSUl5cjPZlW2s9AX3Y1UCdAhsJ2pajJHdRp
        rhpbhTC+E/tlm3ndCeT/nxj8w==)

* To get your normal tfstate content back, decrypt

        terrahelp decrypt  

* Verify `terraform.tfstate` contents after decryption. This should now look exactly the same
as it did before doing the encryption



