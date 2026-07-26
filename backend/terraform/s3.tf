resource "aws_s3_bucket" "manifest" {
  bucket = "${var.project_name}-manifest"
}

