## Terrahelp example - encryption & decryption 

This example contains a very simple terraform setup composed entirely of local resources (e.g. template resource) and exists in order to demonstrate how you can do basic encryption and decryption functionality in the absence of a formal solution (ref https://github.com/hashicorp/terraform/issues/516).
 
This example is completely safe to run and will not land up costing you any money in a cloud provider! It currently demonstrates a terraform 0.12.x based setup which includes new syntax.
 
The CLI itself offers a more comprehensive view of the various options available, so please use this if you need more info.
Additionally you can read this corresponding blog which gives a more detailed explanation of this functionality and its usage: [Securing Terraform State with Vault](https://www.opencredo.com/securing-terraform-state-with-vault).

### Simple inline encryption of terraform output

This example will demonstrate _inline_ encryption and decryption using the _simple_ encryption provider where we will pipe the content (`terraform plan` in this case) directly into it. This specific example uses the basic command line arguments as opposed to environment variables to control the process, and assumes you have opened a terminal window in this directory and have the terraform binary available on your path.

* Initialise the root Terraform module
    ```bash
    terraform init
    ```

* Run a plan over the existing example.
    ```bash
    terraform plan
            
    Refreshing Terraform state in-memory prior to plan...
    The refreshed state will be used to calculate this plan, but will not be
    persisted to local or remote state storage.
    
    
    ------------------------------------------------------------------------
    
    An execution plan has been generated and is shown below.
    Resource actions are indicated with the following symbols:
      + create
    
    Terraform will perform the following actions:
    
      # template_dir.config will be created
      + resource "template_dir" "config" {
          + destination_dir = "./renders"
          + id              = (known after apply)
          + source_dir      = "./templates"
          + vars            = {
              + "msg1" = "sensitive-value-1-AK#%DJGHS*G"
              + "msg2" = "normal value 1"
              + "msg3" = "sensitive-value-3-//dfhs//"
              + "msg4" = "sensitive-value-4 with equals sign i.e. ff=yy"
              + "msg5" = "sensitive-list-val-1"
              + "msg6" = "sensitive-flatmap-val-foo"
              + "msg7" = "sensitive-flatmap-val"
            }
        }
    
    Plan: 1 to add, 0 to change, 0 to destroy.
    
    ------------------------------------------------------------------------
    
    Note: You didn't specify an "-out" parameter to save this plan, so Terraform
    can't guarantee that exactly these actions will be performed if
    "terraform apply" is subsequently run.
    ```

* Run the same command, but pipe the output into the `terrahelp encrypt` command. The default provider is the simple provider so you do not need to explicitly set this, although you do need to provide and encryption key.
  ```bash
  terraform plan | terrahelp encrypt -mode=inline -simple-key=057EFE8CF0F15DE86876F9E313E3D0D6 
  ``` 

* The result should now look something like below:
    ```bash
    Refreshing Terraform state in-memory prior to plan...
    The refreshed state will be used to calculate this plan, but will not be
    persisted to local or remote state storage.
    
    
    ------------------------------------------------------------------------
    
    An execution plan has been generated and is shown below.
    Resource actions are indicated with the following symbols:
      + create
    
    Terraform will perform the following actions:
    
      # template_dir.config will be created
      + resource "template_dir" "config" {
          + destination_dir = "./renders"
          + id              = (known after apply)
          + source_dir      = "./templates"
          + vars            = {
              + "msg1" = "@terrahelp-encrypted(u2DrCMPUTbj0ge6mlHpCxvsuJmEua5FFRElZY86vVh0Ya7AU02jOOE+qbcAy)"
              + "msg2" = "normal value 1"
              + "msg3" = "@terrahelp-encrypted(tgJCKQc/L8xYAk6pe7Zo7O5nvp1V1FWQbdhtzBTgopB4KtDJCAboulen)"
              + "msg4" = "@terrahelp-encrypted(LIxWr4Xdwg3oAwcfoKSfoi+pxltF87XLksiCsy+t/GE9B39uRTF4Dz2zPWn3kU0jB83ZzXBFjfjkF8f03g==)"
              + "msg5" = "@terrahelp-encrypted(s9tqpkHbSNsmuR0/pBdQKblqB7kqaxTihT0P18vAo29p4zbc)"
              + "msg6" = "@terrahelp-encrypted(cSsGZUjLJ/uW+exPysK7YeogqbmONy7mXmjiajhnGTyJtjF1/Qpi4jo=)"
              + "msg7" = "@terrahelp-encrypted(lEd5lrvy419VXqAxvbpCYn6gD8FAvafH/Xlqx7U0fv2CJo9gfQ==)"
            }
        }
    
    Plan: 1 to add, 0 to change, 0 to destroy.
    
    ------------------------------------------------------------------------
    
    Note: You didn't specify an "-out" parameter to save this plan, so Terraform
    can't guarantee that exactly these actions will be performed if
    "terraform apply" is subsequently run.
    ```

* For decryption, you could pipe the output again into the `terrahelp decrypt` command, however more than likely, you will probably want to save the results into a file and then decrypt that. The sequence of commands to do that would be something as follows:
  ```bash
  terraform plan -out=my-infra.plan
  
  terrahelp encrypt -mode=inline -simple-key=057EFE8CF0F15DE86876F9E313E3D0D6 -file=my-infra.plan
  
  terrahelp decrypt -mode=inline -simple-key=057EFE8CF0F15DE86876F9E313E3D0D6 -file=my-infra.plan
  ```

### Simple inline encryption of tfstate files

This example will also demonstrate _inline_ encryption and decryption using the _simple_ encryption provider using explicit command line arguments (an example using environment variables is shown with the Vault provider example), however unlike the example above, will operate over the main terraform tfstate files.

* Run terraform as normal, including apply
    ```bash
    terraform plan
    terraform apply`
    ```

* Verify `terraform.tfstate` contents before encryption (e.g. by doing a `cat terraform.tfstate`).
This should look something like below:
    
   ```bash
   {
     "version": 4,
     "terraform_version": "0.12.6",
     "serial": 1,
     "lineage": "6c7602e5-f319-09d1-ac47-2f1841d04b1c",
     "outputs": {
       "normal_val_2": {
         "value": "normal value 2",
         "type": "string"
       }
     },
     "resources": [
       {
         "mode": "managed",
         "type": "template_dir",
         "name": "config",
         "provider": "provider.template",
         "instances": [
           {
             "schema_version": 0,
             "attributes": {
               "destination_dir": "./renders",
               "id": "17a2b4f0fb66d778755329436bf09f1a5f96f1bf",
               "source_dir": "./templates",
               "vars": {
                 "msg1": "sensitive-value-1-AK#",
                 "msg2": "normal value 1",
                 "msg3": "sensitive-value-3",
                 "msg4": "sensitive-value-4 with equals sign i.e. ff=yy",
                 "msg5": "sensitive-list-val-1",
                 "msg6": "sensitive-flatmap-val-foo",
                 "msg7": "sensitive-flatmap-val"
               }
             },
             "private": "bnVsbA=="
           }
         ]
       }
     ]
   }
 
   ```

* Run the `terrahelp encrypt` command, except this time you will explicitly specify the file to encrypt (i.e. `terraform.tfstate`). In the absence of a `-file` argument, terrahelp will assume you are piping the input stream in.

        terrahelp encrypt -mode=inline -simple-key=057EFE8CF0F15DE86876F9E313E3D0D6 -file=terraform.tfstate

* Inspect `terraform.tfstate` content after encryption. Note how all the sensitive values, as detected in the `terraform.tfvars` file, have now been replaced with encrypted versions. The content should look something like that below: 
    ```bash
    {
      "version": 4,
      "terraform_version": "0.12.6",
      "serial": 1,
      "lineage": "0989590c-ee5d-4a58-907f-24e469af36d4",
      "outputs": {
        "normal_val_2": {
          "value": "normal value 2",
          "type": "string"
        }
      },
      "resources": [
        {
          "mode": "managed",
          "type": "template_dir",
          "name": "config",
          "provider": "provider.template",
          "instances": [
            {
              "schema_version": 0,
              "attributes": {
                "destination_dir": "./renders",
                "id": "0e22fe020e90ea9f409194fb6e7cf0663ec02196",
                "source_dir": "./templates",
                "vars": {
                  "msg1": "@terrahelp-encrypted(2RSF9c56qwQVjCvREZiCCdgj8+5kdNC8Up/2m6bL5APFYA0+QMZ2evjIkSCR)",
                  "msg2": "normal value 1",
                  "msg3": "@terrahelp-encrypted(PxjWuGX6DsIx2HU8HGhaph59xVBlsHDqvX9K6dAtkG44rZjPgNOu/Ju2)",
                  "msg4": "@terrahelp-encrypted(x2ubbE+69/o7AFFqx3oI2mZhyFOKrejsJR+xa9wTq+8YUrT4RgAZ7CC6xGbXrux1D3+/G+iEFA/BrWplvw==)",
                  "msg5": "@terrahelp-encrypted(nywjFDzFLiUMNu9pIZN+Jy/J8bbVilPi5IZ9rzt9Ituzr3j0)",
                  "msg6": "@terrahelp-encrypted(Zkuwq7KR6Y5eylLwMFqoWGxshOVBy3EVgRRBUKhfc2pW8ASujiu5TI4=)",
                  "msg7": "@terrahelp-encrypted(0h/U4n/9CmwhOh7Fiw/1zA3I+5pNSJ2y2WZX8v3ztvDTsHFLmg==)"
                }
              },
              "private": "bnVsbA=="
            }
          ]
        }
      ]
    }
    
    ```

* To get your normal `terraform.tfstate` content back, simply run the `decrypt` command with the same arguments as above.
    ```bash
    terrahelp decrypt -mode=inline -simple-key=057EFE8CF0F15DE86876F9E313E3D0D6 -file=terraform.tfstate
    ```

* Again verify `terraform.tfstate` content after decryption. This should now look exactly the same as it did before doing the encryption

### Vault full encryption of tfstate files

This example will demonstrate _full_ encryption and decryption using the _vault_ encryption provider (tested against Vault Server v1.2.1). On this occasion, instead of explicitly configuring the process via command line arguments, we will use environment variables.

* First, you will need to ensure you have a running Vault server available. You can quite easily download the latest version from [here](https://www.vaultproject.io/downloads.html), then open up a new terminal, and for experimentation purposes, simply run the server in dev mode as below:
    ```bash
    vault server -dev -dev-root-token-id="terrahelp-devonly-vault-root-token"
    ```

* In a separate terminal, ensure you change into this example project folder again, and setup the necessary environment variables required for us to talk to our dev Vault server, as well as run the next set of `terrahelp` commands. Specifically will also run the `vault-autoconfig` command to configure Vault with the named encryption key we wnat to use. i.e.
    ```bash
    export VAULT_TOKEN="terrahelp-devonly-vault-root-token"
    export VAULT_ADDR="http://127.0.0.1:8200"
    export VAULT_SKIP_VERIFY="true"
    
    export TH_ENCRYPTION_PROVIDER="vault"
    export TH_ENCRYPTION_MODE="full"
    export TH_VAULT_NAMED_KEY="examplekey"
    
    terrahelp vault-autoconfig
    ```

* Run terraform as normal and inspect the terraform.tfstate content before encryption is applied
    ```bash
    terraform plan
    terraform apply`
    ```

* Verify `terraform.tfstate` contents before encryption (e.g. by doing a `cat terraform.tfstate`).
This should look something like below:
    
    ```bash
    {
      "version": 4,
      "terraform_version": "0.12.6",
      "serial": 2,
      "lineage": "0989590c-ee5d-4a58-907f-24e469af36d4",
      "outputs": {
        "normal_val_2": {
          "value": "normal value 2",
          "type": "string"
        }
      },
      "resources": [
        {
          "mode": "managed",
          "type": "template_dir",
          "name": "config",
          "provider": "provider.template",
          "instances": [
            {
              "schema_version": 0,
              "attributes": {
                "destination_dir": "./renders",
                "id": "0e22fe020e90ea9f409194fb6e7cf0663ec02196",
                "source_dir": "./templates",
                "vars": {
                  "msg1": "sensitive-value-1-AK#%DJGHS*G",
                  "msg2": "normal value 1",
                  "msg3": "sensitive-value-3-//dfhs//",
                  "msg4": "sensitive-value-4 with equals sign i.e. ff=yy",
                  "msg5": "sensitive-list-val-1",
                  "msg6": "sensitive-flatmap-val-foo",
                  "msg7": "sensitive-flatmap-val"
                }
              }
            }
          ]
        }
      ]
    }
    ```

* Run the `terrahelp encrypt` command, specifying the file to encrypt (i.e. `terraform.tfstate`). In the absence of a `-file` argument, terrahelp will assume you are piping the input stream in. The other arguments should be picked up from the environment variables.

         terrahelp encrypt -file=terraform.tfstate   

* Inspect `terraform.tfstate` content after encryption. Note how all the sensitive values, as detected in the terraform.tfvars file, have now been replaced with encrypted versions, and will look something like below: 
    ```
    @terrahelp-encrypted(vault:v1:tS3hHe1pmzigzTim0PPlxBRv7tlJ/o7Us/jFRakyjwnqg+OqoVfONZuIXzGAb7VPd8YixZDuNDxu2hVo8T0
    Edebv4UmaHR9zIg2WUJ2dWONkK2sWzJ1xHdr9nEdfLVtQRF7NLhMh7Rw08HdnenkWNy3eo/l2U1gZNrJxGEVGnaSa0tEhrhRMK/aYONQps0up5e9r
    UpsOnFbqC7i1cnzQELn/xCfsbjaIwmlsr2swXvqzlZ2prMg6LEACIWRraXxiONn5ql5svcAQwY7lt/QKazzHz7976R2XA9h6AWBZqeKtC5m0ovqW
    lixEHetroT11HAWop4dp5AS4Mr6jjiFDCHv5vwmkbKKUsfjDR+aaWzXyr7nwSd8sHl1WLGgFQqZssXd6Yj7SGJkxbLiEy9CZS6Fv+Lq8SSeUtkMd        
    ...                
    K3zra1QEk2hh1aWy9EsuK/TOqO7sS928BBYfsoUGLrj3M29WhyZx1sF+aSYy+rAMpz0vMeA/fRUWNnuI6lRWdevyJ6gcMFFh53/DBecZRtKOO0vH
    w6bzpQFYFtbQQk0RYxWlMf0qiskHA65xt9YXat+Yjf+euyVhrxGajwuhSZ8RX3VEgU2or7kAdM4L7+s4h69ObNZ2mToLbXw3JnIsRRWYdePqTii4
    QMmAk4fVwvFZuKilc1aLhMQkF1vTFhxVuXMrEChQHaTN/GH1Mna1jrkTF9u8KjyawypYbuj7M1QJjCXFOASnJm1tqYHuZ2YxTvvmZH7cjftUpzyv
    6hIOjLSNhzthE4eCid/j4mDu1GSEm3s/8JhV1C9NuO3AJPdGVYPRiYr/0DsKB+7zy238nHKCMICLUnYZa4JLrmTkzv2c/Ko+6wBEkW/mbiLb9/Ew
    cjosfE00F1ro2Kvyrsac7bAUjEaBWZAHeesPL9FBEjDOxqluK8caABlWCDSLog1s/dEbb8YbKLuQbRg5amtGvSOZ27Zv8Avs2m41PMj7PA==)
    ```

* To get your normal tfstate content back, simply run the `decrypt` version of the above command i.e.
    ```bash
    terrahelp decrypt -file=terraform.tfstate  
    ```

* Verify `terraform.tfstate` contents after decryption. This should now look exactly the same as it did before doing the encryption

### How does it work out what is considered sensitive?
At present, `terrahelp` relies on using the `terraform.tfvars` file as the mechanism to indicate which values should be considered sensitive, and thus encrypted when detected. 
