# -------------------------------------------------
#      Example terraform.tfvars file based on 0.7.7
#      Note: this is only for testing / example purposes
#            this file should NEVER really be checked into
#            version control
# -------------------------------------------------
# Some comment
pretend_aws_access_key     = "madeup-aws-access-key-PEJFNS"
pretend_aws_secret_key     = "madeup-aws-secret-key-KGSDGH"
tf_sensitive_key_1         = "sensitive-value-1-AK#%DJGHS*G"
tf_sensitive_key_2         = "sensitive-value-2-prYh57"
tf_sensitive_key_3         = "sensitive-value-3-//dfhs//"

# Some more comments
tf_sensitive_key_4         = "sensitive-value-4 with equals sign i.e. ff=yy"
# tf_sensitive_key_5         = "encrypted-value-5"
tf_sensitive_key_6         = "sensitive-value-6"

# new list and maps (terraform 0.7.x and higher)
tf_sensitive_list_vals = [
  "sensitive-list-val-1",
  "sensitive-list-val-2",
  "sensitive-list-val"
]

tf_sensitive_flatmap_vals = {
  foo       = "sensitive-flatmap-val-foo"
  bax       = "sensitive-flatmap-val-bax"
  "bob"     = "sensitive-flatmap-val-bob"
  "overlap" = "sensitive-flatmap-val"
}
