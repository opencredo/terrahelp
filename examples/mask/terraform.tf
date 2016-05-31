# Something basic
provider "aws" {
  access_key = "${var.pretend_aws_access_key}"
  secret_key = "${var.pretend_aws_secret_key}"
  region = "us-east-1"
}

resource "template_file" "example" {
  template = "\nmsg1 = ${msg1}\nmsg2 = ${msg2}\nmsg3 = ${msg3}"
  vars {
    msg1 = "${var.tf_sensitive_key_1}"
    msg2 = "${var.tf_normal_key_1}"
    msg3 = "${var.tf_sensitive_key_3}"
  }
}

output "rendered" {
  value = "${template_file.example.rendered}"
}

output "normal_val_2" {
  value = "${var.tf_normal_key_2}"
}
