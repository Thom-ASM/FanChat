provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "fan-chat"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
