terraform {
  backend "s3" {
    bucket       = "your-terraform-state-bucket"
    key          = "${var.project_name}/terraform.tfstate"
    region       = "eu-west-2"
    encrypt      = true
    use_lockfile = true
  }
}
