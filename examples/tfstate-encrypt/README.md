## Terrahelp example - encryption & decryption 

This example contains a very simple terraform setup composed entirely of local resources (e.g. template resource) and exists in order to demonstrate how you can do basic encryption and decryption functionality in the absence of a formal solution (ref https://github.com/hashicorp/terraform/issues/516).
 
This example is completely safe to run and will not land up costing you any money in a cloud provider! It currently demonstrates a terraform 0.7.7 based setup which includes the new lists and maps functionality.
 
The CLI itself offers a more comprehensive view of the various options available, so please use this if you need more info.
Additionally you can read this corresponding blog which gives a more detailed explanation of this functionality and its usage: [Securing Terraform State with Vault](https://www.opencredo.com/securing-terraform-state-with-vault).

### Simple inline encryption of terraform output

This example will demonstrate _inline_ encryption and decryption using the _simple_ encryption provider where we will pipe the content (`terraform plan` in this case) directly into it. This specific example uses the basic command line arguments as opposed to environment variables to control the process, and assumes you have opened a terminal window in this directory and have the terraform binary available on your path.

* Run a `terraform plan` as normal

        terraform plan
        
* Inspect the result which should look something like below:        
        
        Refreshing Terraform state in-memory prior to plan...
        The refreshed state will be used to calculate this plan, but
        will not be persisted to local or remote state storage.
        
        ...
                
        <= data.template_file.example
            rendered:  "<computed>"
            template:  "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}\nmsg4 = ${msg4}\nmsg5 = ${msg5}\nmsg6 = ${msg6}\nmsg7 = ${msg7}"
            vars.%:    "7"
            vars.msg1: "sensitive-value-1-AK#%DJGHS*G"
            vars.msg2: "normal value 1"
            vars.msg3: "sensitive-value-3-//dfhs//"
            vars.msg4: "sensitive-value-4 with equals sign i.e. ff=yy"
            vars.msg5: "sensitive-list-val-1"
            vars.msg6: "sensitive-flatmap-val-foo"
            vars.msg7: "sensitive-flatmap-val"
        
        
        Plan: 0 to add, 0 to change, 0 to destroy.

* Run the same command, but pipe the output into the `terrahelp encrypt` command. The default provider is the simple provider so you do not need to explicitly set this, although you do need to provide and encryption key.

        terraform plan | terrahelp encrypt -mode=inline -simple-key=AES256Key-32Characters0987654321 

* The result should now look something like below:

        Refreshing Terraform state in-memory prior to plan...
        The refreshed state will be used to calculate this plan, but
        will not be persisted to local or remote state storage.
        
        ...
        
        <= data.template_file.example
            rendered:  "<computed>"
            template:  "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}\nmsg4 = ${msg4}\nmsg5 = ${msg5}\nmsg6 = ${msg6}\nmsg7 = ${msg7}"
            vars.%:    "7"
            vars.msg1: "@terrahelp-encrypted(xufN6OOCI2TWDp793/zlba4nt3dUnsbbQpB64HTykYPr3+ZUKgze+fgbj2zW)"
            vars.msg2: "normal value 1"
            vars.msg3: "@terrahelp-encrypted(v6Wt2f1w2xvjHI8bsTXK51hrLOtQPswvTzWv+kGj7ojZAJcgf5POFT08)"
            vars.msg4: "@terrahelp-encrypted(iobUvjF5d4rc3q4GCrED3vUSz7gpCNnXM/Taah9OuVV5WDXEMRgxCGIxIiN5Die/JFkCgt+IoOiEL7nOcQ==)"
            vars.msg5: "@terrahelp-encrypted(7cZ2iwc00eLcDOBrP9pVtdlZErRHGr6hl++UynU1jnhRVjwV)"
            vars.msg6: "@terrahelp-encrypted(/kYPdcP3ROpchiHjGv7fysPIZfCnTYpR4XX841jAz2R317QYO/A+nf0=)"
            vars.msg7: "@terrahelp-encrypted(oIdMOgF6Wzg/s6KpRmYTZCDP7RiHw3EZwyc2+A4PSouEkD07GA==)"

* For decryption, you could pipe the output again into the `terrahelp decrypt` command, however more than likely, you will probably want to save the results into a file and then decrypt that. The sequence of commands to do that would be something as follows:
  
        terraform plan -out=my-infra.plan
        
        terrahelp encrypt -mode=inline -simple-key=AES256Key-32Characters0987654321 -file=my-infra.plan
        
        terrahelp decrypt -mode=inline -simple-key=AES256Key-32Characters0987654321 -file=my-infra.plan        

### Simple inline encryption of tfstate files

This example will also demonstrate _inline_ encryption and decryption using the _simple_ encryption provider using explicit command line arguments (an example using environment variables is shown with the Vault provider example), however unlike the example above, will operate over the main terraform tfstate files.

* Run terraform as normal, including apply

        terraform plan
        terraform apply

* Verify `terraform.tfstate` contents before encryption (e.g. by doing a `cat terraform.tfstate`).
This should look something like below:
    
        {
            "version": 3,
            "terraform_version": "0.7.7",
            "serial": 0,
            "lineage": "dfb415f5-07c5-478b-945e-a592f1cf09b6",
            "modules": [
                {
                    "path": [
                        "root"
                    ],
                    "outputs": {
                        "normal_val_2": {
                            "sensitive": false,
                            "type": "string",
                            "value": "normal value 2"
                        },
                        "rendered": {
                            "sensitive": false,
                            "type": "string",
                            "value": "\nmsg1 = sensitive-value-1-AK#%DJGHS*G\nmsg2 = normal value 1\nmsg3 = sensitive-value-3-//dfhs//\nmsg4 = sensitive-value-4 with equals sign i.e. ff=yy\nmsg5 = sensitive-list-val-1\nmsg6 = sensitive-flatmap-val-foo\nmsg7 = sensitive-flatmap-val"
                        }
                    },
                    "resources": {
                        "data.template_file.example": {
                            "type": "template_file",
                            "depends_on": [],
                            "primary": {
                                "id": "88dafb613a33265583a1ba802edb4d9ffafe604602d764c780b1db8c76c6c7fe",
                                "attributes": {
                                    "id": "88dafb613a33265583a1ba802edb4d9ffafe604602d764c780b1db8c76c6c7fe",
                                    "rendered": "\nmsg1 = sensitive-value-1-AK#%DJGHS*G\nmsg2 = normal value 1\nmsg3 = sensitive-value-3-//dfhs//\nmsg4 = sensitive-value-4 with equals sign i.e. ff=yy\nmsg5 = sensitive-list-val-1\nmsg6 = sensitive-flatmap-val-foo\nmsg7 = sensitive-flatmap-val",
                                    "template": "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}\nmsg4 = ${msg4}\nmsg5 = ${msg5}\nmsg6 = ${msg6}\nmsg7 = ${msg7}",
                                    "vars.%": "7",
                                    "vars.msg1": "sensitive-value-1-AK#%DJGHS*G",
                                    "vars.msg2": "normal value 1",
                                    "vars.msg3": "sensitive-value-3-//dfhs//",
                                    "vars.msg4": "sensitive-value-4 with equals sign i.e. ff=yy",
                                    "vars.msg5": "sensitive-list-val-1",
                                    "vars.msg6": "sensitive-flatmap-val-foo",
                                    "vars.msg7": "sensitive-flatmap-val"
                                },
                                "meta": {},
                                "tainted": false
                            },
                            "deposed": [],
                            "provider": ""
                        }
                    },
                    "depends_on": []
                }
            ]
        }

* Run the `terrahelp encrypt` command, except this time you will explicitly specify the file to encrypt (i.e. `terraform.tfstate`). In the absence of a `-file` argument, terrahelp will assume you are piping the input stream in.

        terrahelp encrypt -mode=inline -simple-key=AES256Key-32Characters0987654321 -file=terraform.tfstate

* Inspect `terraform.tfstate` content after encryption. Note how all the sensitive values, as detected in the `terraform.tfvars` file, have now been replaced with encrypted versions. The content should look something like that below: 

        {
            "version": 3,
            "terraform_version": "0.7.7",
            "serial": 0,
            "lineage": "07ecb11e-8b77-41c5-a07b-ca924adaf6bb",
            "modules": [
                {
                    "path": [
                        "root"
                    ],
                    "outputs": {
                        "normal_val_2": {
                            "sensitive": false,
                            "type": "string",
                            "value": "normal value 2"
                        },
                        "rendered": {
                            "sensitive": false,
                            "type": "string",
                            "value": "\nmsg1 = @terrahelp-encrypted(mWlSsFQFNaK0pubo17cx8ruVMmldpUYyhx83nDtRvASKReXeQhVQEXaWCsRg)\nmsg2 = normal value 1\nmsg3 = @terrahelp-encrypted(W+236mopyxfruRnmXqo5tMhCjS1h7Al4kAp5yZ+vT9wb2VS35js4T4NJ)\nmsg4 = @terrahelp-encrypted(qqscl1+HMyCxnSw9sSLbVE05CZ+LIOvpNFtpRflp5H7HTp0NK1SLRhsjG775KdB3mrLB5yYJ9uv0fjPbpw==)\nmsg5 = @terrahelp-encrypted(5mubnPPs0P7wndVFh6H0wG20V+ljzqg7+ZNv5jWcWyTEHk54)\nmsg6 = @terrahelp-encrypted(oX+He8/1SlU6vFIyWbklTqxOAoiQrcLUaEcsqjkkFH1kcjDTOncZNLc=)\nmsg7 = @terrahelp-encrypted(kNEhI+mAkHVaGczjIYH4peO0CLsNDZsZIAFH6jE9ReJqynE4Dw==)"
                        }
                    },
                    "resources": {
                        "data.template_file.example": {
                            "type": "template_file",
                            "depends_on": [],
                            "primary": {
                                "id": "88dafb613a33265583a1ba802edb4d9ffafe604602d764c780b1db8c76c6c7fe",
                                "attributes": {
                                    "id": "88dafb613a33265583a1ba802edb4d9ffafe604602d764c780b1db8c76c6c7fe",
                                    "rendered": "\nmsg1 = @terrahelp-encrypted(mWlSsFQFNaK0pubo17cx8ruVMmldpUYyhx83nDtRvASKReXeQhVQEXaWCsRg)\nmsg2 = normal value 1\nmsg3 = @terrahelp-encrypted(W+236mopyxfruRnmXqo5tMhCjS1h7Al4kAp5yZ+vT9wb2VS35js4T4NJ)\nmsg4 = @terrahelp-encrypted(qqscl1+HMyCxnSw9sSLbVE05CZ+LIOvpNFtpRflp5H7HTp0NK1SLRhsjG775KdB3mrLB5yYJ9uv0fjPbpw==)\nmsg5 = @terrahelp-encrypted(5mubnPPs0P7wndVFh6H0wG20V+ljzqg7+ZNv5jWcWyTEHk54)\nmsg6 = @terrahelp-encrypted(oX+He8/1SlU6vFIyWbklTqxOAoiQrcLUaEcsqjkkFH1kcjDTOncZNLc=)\nmsg7 = @terrahelp-encrypted(kNEhI+mAkHVaGczjIYH4peO0CLsNDZsZIAFH6jE9ReJqynE4Dw==)",
                                    "template": "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}\nmsg4 = ${msg4}\nmsg5 = ${msg5}\nmsg6 = ${msg6}\nmsg7 = ${msg7}",
                                    "vars.%": "7",
                                    "vars.msg1": "@terrahelp-encrypted(mWlSsFQFNaK0pubo17cx8ruVMmldpUYyhx83nDtRvASKReXeQhVQEXaWCsRg)",
                                    "vars.msg2": "normal value 1",
                                    "vars.msg3": "@terrahelp-encrypted(W+236mopyxfruRnmXqo5tMhCjS1h7Al4kAp5yZ+vT9wb2VS35js4T4NJ)",
                                    "vars.msg4": "@terrahelp-encrypted(qqscl1+HMyCxnSw9sSLbVE05CZ+LIOvpNFtpRflp5H7HTp0NK1SLRhsjG775KdB3mrLB5yYJ9uv0fjPbpw==)",
                                    "vars.msg5": "@terrahelp-encrypted(5mubnPPs0P7wndVFh6H0wG20V+ljzqg7+ZNv5jWcWyTEHk54)",
                                    "vars.msg6": "@terrahelp-encrypted(oX+He8/1SlU6vFIyWbklTqxOAoiQrcLUaEcsqjkkFH1kcjDTOncZNLc=)",
                                    "vars.msg7": "@terrahelp-encrypted(kNEhI+mAkHVaGczjIYH4peO0CLsNDZsZIAFH6jE9ReJqynE4Dw==)"
                                },
                                "meta": {},
                                "tainted": false
                            },
                            "deposed": [],
                            "provider": ""
                        }
                    },
                    "depends_on": []
                }
            ]
        }

* To get your normal `terraform.tfstate` content back, simply run the `decrypt` command with the same arguments as above.

        terrahelp decrypt -mode=inline -simple-key=AES256Key-32Characters0987654321 -file=terraform.tfstate 

* Again verify `terraform.tfstate` content after decryption. This should now look exactly the same as it did before doing the encryption


### Vault full encryption of tfstate files

This example will demonstrate _full_ encryption and decryption using the _vault_ encryption provider (tested against Vault Server 0.5.2). On this occasion, instead of explicitly configuring the process via command line arguments, we will use environment variables.

* First, you will need to ensure you have a running Vault server available. You can quite easily download the latest version from [here](https://www.vaultproject.io/downloads.html), then open up a new terminal, and for experimentation purposes, simply run the server in dev mode as below:

        vault server -dev -dev-root-token-id="terrahelp-devonly-vault-root-token"

* In a separate terminal, ensure you change into this example project folder again, and setup the necessary environment variables required for us to talk to our dev Vault server, as well as run the next set of `terrahelp` commands. Specifically will also run the `vault-autoconfig` command to configure Vault with the named encryption key we wnat to use. i.e.

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


* Run the `terrahelp encrypt` command, specifying the file to encrypt (i.e. `terraform.tfstate`). In the absence of a `-file` argument, terrahelp will assume you are piping the input stream in. The other arguments should be picked up from the environment variables.

         terrahelp encrypt -file=terraform.tfstate   

* Inspect `terraform.tfstate` content after encryption. Note how all the sensitive values, as detected in the terraform.tfvars file, have now been replaced with encrypted versions, and will look something like below: 

        @terrahelp-encrypted(vault:v1:h7Yx1VAYvd2pyW0dd/iWifSe6yFB8QI7Zv2KjlW5USM5AyT9o3g3U2bU3
        vbDweRCGUXq2P8qpNcp8LUXDUon2Q6ee8I20X6yJyj5I2AS9V9ec4YcFOS9odqG+6dFqdlgWUkvEXPsH6puL0rX
        depvR17dvK1QTID0iE14HS7b4UnwI0Ti+f2VX4GvKHhnfKwCejKVu3g2bXdjn35h+EH9cHonSTx24SI6mM5k9Uy
        L7ht7AfPtPkdiUW7XSiW69UsZ+ZWrz8  ...  Ej3NYiY71Z/B2Rfm3M3V22BjfCsoUAHR1gL8acb5xQryuk+B/
        zQdLx7fXgxS8rMPKFwrJVRVtdcJtLFtLLf42AV1oUCqYvvusyNiGkQ6p3/2cgbkWsm/gN2lc26AuD6wVtd44qi
        CKK5iBZU4HQH6P5dycL0Sjgg4vJvcve85fQOLtfrr+UnQP0hdTSfSUl5cjPZlW2s9AX3Y1UCdAhsJ2pajJHdRp
        rhpbhTC+E/tlm3ndCeT/nxj8w==)

* To get your normal tfstate content back, simply run the `decrypt` version of the above command i.e.

        terrahelp decrypt -file=terraform.tfstate  

* Verify `terraform.tfstate` contents after decryption. This should now look exactly the same as it did before doing the encryption


### How does it work out what is considered sensitive?
At present, `terrahelp` relies on using the `terraform.tfvars` file as the mechanism to indicate which values should be considered sensitive, and thus encrypted when detected. 
