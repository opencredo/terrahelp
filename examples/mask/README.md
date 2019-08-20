## Terrahelp example - masking sensitive data 

This example contains a very simple terraform setup composed entirely of local resources (e.g. template resource) and exists in order to demonstrate how you can do masking of sensitive data which may be output from various terraform commands.
 
This example is completely safe to run and will not land up costing you any money in a cloud provider! It currently demonstrates a terraform 0.12.x based setup which includes new syntax.
 
The CLI itself offers a more comprehensive view of the various options available, so please use this if you need more info.

### Simple inline masking of terraform output

This example will demonstrate how you can use the `mask` command in order to mask sensitive data which may be exposed when performing terraform actions.

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
              + "msg1" = "sensitive-value-1-AK#"
              + "msg2" = "normal value 1"
              + "msg3" = "sensitive-value-3"
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

* Run the same command, but pipe the output through the `terrahelp mask` command. 
    ```bash
    terraform plan | terrahelp mask
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
              + "msg1" = "******"
              + "msg2" = "normal value 1"
              + "msg3" = "******"
              + "msg4" = "******"
              + "msg5" = "******"
              + "msg6" = "******"
              + "msg7" = "******"
            }
        }
    
    Plan: 1 to add, 0 to change, 0 to destroy.
    
    ------------------------------------------------------------------------
    
    Note: You didn't specify an "-out" parameter to save this plan, so Terraform
    can't guarantee that exactly these actions will be performed if
    "terraform apply" is subsequently run.
    ```

To change the mask character and/or length, you can use the `-maskchar` and `-numchars` flags.
    ```bash
    terraform plan | terrahelp mask -maskchar=# -numchars=3
    ```

By default, the mask command will also attempt to detect whether any previous sensitive data may be exposed, and if so will mask this as well. This may happen for example when changing the value of one sensitive value to another e.g.

```bash
  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "46739816783755f5ca0127ca0b16863effd53533" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
          ~ "msg1" = "sensitive-value-1-AK#" -> "sensitive-value-1"
            "msg2" = "normal value 1"
            "msg3" = "sensitive-value-3"
            "msg4" = "sensitive-value-4 with equals sign i.e. ff=yy"
            "msg5" = "sensitive-list-val-1"
            "msg6" = "sensitive-flatmap-val-foo"
            "msg7" = "sensitive-flatmap-val"
        }
    }
```

In which case the resulting mask will look as follows 

```bash
  # template_dir.config must be replaced
-/+ resource "template_dir" "config" {
        destination_dir = "./renders"
      ~ id              = "46739816783755f5ca0127ca0b16863effd53533" -> (known after apply)
        source_dir      = "./templates"
      ~ vars            = { # forces replacement
          ~ "msg1" = "******" -> "******"
            "msg2" = "normal value 1"
            "msg3" = "******"
            "msg4" = "******"
            "msg5" = "******"
            "msg6" = "******"
            "msg7" = "******"
        }
    }
```

If you want to suppress this default behaviour you can use the `-prev=false`
            
### How does it work out what is considered sensitive?
At present, `terrahelp` relies on using the `terraform.tfvars` file as the mechanism to indicate which values should be considered sensitive, and thus masked out when detected.  

