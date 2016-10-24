## Terrahelp example - masking sensitive data 

This example contains a very simple terraform setup composed entirely of local resources (e.g. template resource) and exists in order to demonstrate how you can do masking of sensitive data which may be output from varius terraform commands.
 
This example is completely safe to run and will not land up costing you any money in a cloud provider! It currently demonstrates a terraform 0.7.7 based setup which includes the new lists and maps functionality.
 
The CLI itself offers a more comprehensive view of the various options available, so please use this if you need more info.

### Simple inline masking of terraform output

This example will demonstrate how you can use the `mask` command in order to mask sensitive data which may be exposed when performing terraform actions.

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
                


* Run the same command, but pipe the output through the `terrahelp mask` command. 

        terraform plan | terrahelp mask  

* The result should now look something like below:

        Refreshing Terraform state in-memory prior to plan...
        The refreshed state will be used to calculate this plan, but
        will not be persisted to local or remote state storage.
        
        ...
        
        <= data.template_file.example
            rendered:  "<computed>"
            template:  "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}\nmsg4 = ${msg4}\nmsg5 = ${msg5}\nmsg6 = ${msg6}\nmsg7 = ${msg7}"
            vars.%:    "7"
            vars.msg1: "******"
            vars.msg2: "normal value 1"
            vars.msg3: "******"
            vars.msg4: "******"
            vars.msg5: "******"
            vars.msg6: "******"
            vars.msg7: "******"
        
        
        Plan: 0 to add, 0 to change, 0 to destroy.

To change the mask character and/or length, you can use the `-maskchar` and `-numchars` flags, e.g. `terraform plan | terrahelp mask -maskchar=# -numchars=3`

By default, the mask command will also attempt to detect whether any previous sensitive data may be exposed, and if so will mask this as well. This may happen for example when changing the value of one sensitive value to another e.g.

        + template_file.example
            rendered:  "" => "<computed>"
            template:  "" => "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}"
            vars.#:    "" => "3"
            vars.msg1: "old-sensitive-value" => "sensitive-value-1-AK#%DJGHS*G"
            vars.msg2: "" => "normal value 1"
            vars.msg3: "" => "sensitive-value-3-//dfhs//"

In which case the resulting mask will look as follows 

        + template_file.example
            rendered:  "" => "<computed>"
            template:  "" => "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}"
            vars.#:    "" => "3"
            vars.msg1: "******" => "******"
            vars.msg2: "" => "normal value 1"
            vars.msg3: "" => "******"

If you want to suppress this default behaviour you can use the `-prev=false`
            
### How does it work out what is considered sensitive?
At present, `terrahelp` relies on using the `terraform.tfvars` file as the mechanism to indicate which values should be considered sensitive, and thus masked out when detected.  

