variable "project_name" {
  type        = string
  description = "application project name"
  default     = "fan-chat"
}

variable "aws_region" {
  type        = string
  description = "aws region to deploy to"
  default     = "eu-west-2"
}
