# -------------------------------------------------
#      Example terraform file based on 0.7.7
# -------------------------------------------------
provider "aws" {
  access_key = "${var.pretend_aws_access_key}"
  secret_key = "${var.pretend_aws_secret_key}"
  region = "us-east-1"
}

data "template_file" "example" {
  template = "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}\nmsg4 = ${msg4}\nmsg5 = ${msg5}\nmsg6 = ${msg6}\nmsg7 = ${msg7}"
  vars {
    msg1 = "${var.tf_sensitive_key_1}"
    msg2 = "${var.tf_normal_key_1}"
    msg3 = "${var.tf_sensitive_key_3}"
    msg4 = "${var.tf_sensitive_key_4}"
    msg5 = "${var.tf_sensitive_list_vals[0]}"
    msg6 = "${var.tf_sensitive_flatmap_vals["foo"]}"
    msg7 = "${var.tf_sensitive_flatmap_vals["overlap"]}"
  }
}

output "rendered" {
  value = "${data.template_file.example.rendered}"
}

output "normal_val_2" {
  value = "${var.tf_normal_key_2}"
}
