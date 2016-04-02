# ----------------------------------------------------------------
# sensitive variables (should be passed in via tfvars)
# ----------------------------------------------------------------
variable "pretend_aws_access_key" {}
variable "pretend_aws_secret_key" {}
variable "tf_sensitive_key_1"     {}
variable "tf_sensitive_key_2"     {}
variable "tf_sensitive_key_3"     {}
variable "tf_sensitive_key_4"     {}
variable "tf_sensitive_key_6"     {}


# ----------------------------------------------------------------
# Non sensitive variable defaults
# ----------------------------------------------------------------
variable "tf_normal_key_1"     { default = "normal value 1" }
variable "tf_normal_key_2"     { default = "normal value 2" }
