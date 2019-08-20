# -------------------------------------------------
#      Example terraform file based on 0.12.x
# -------------------------------------------------
provider "aws" {
  access_key = var.pretend_aws_access_key
  secret_key = var.pretend_aws_secret_key
  region = "us-east-1"
}

resource "template_dir" "config" {
  source_dir      = "${path.module}/templates"
  destination_dir = "${path.module}/renders"

  vars = {
      msg1 = var.tf_sensitive_key_1
      msg2 = var.tf_normal_key_1
      msg3 = var.tf_sensitive_key_3
      msg4 = var.tf_sensitive_key_4
      msg5 = var.tf_sensitive_list_vals[0]
      msg6 = var.tf_sensitive_flatmap_vals["foo"]
      msg7 = var.tf_sensitive_flatmap_vals["overlap"]
    }
}

output "sensitive_key_1" {
  value = var.tf_sensitive_key_1
}

output "normal_val_2" {
  value = var.tf_normal_key_2
}