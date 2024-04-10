# Global
#
variable "region" {
  type        = string
  description = "AWS region"
  default     = "us-east-1"
}

variable "namespace" {
  type        = string
  description = "Namespace (e.g. `local`)"
  default     = "local"
}

variable "stage" {
  type        = string
  description = "Stage (e.g. `prod`, `dev`, `staging`)"
  default     = "dev"
}

variable "environment" {
  type        = string
  description = "Environment (e.g. `prod`, `dev`, `staging`)"
  default     = "dev"
}

# service's name
variable "service_name" {
  type        = string
  description = "Service name"
  default     = "leaderboards"
}

# DynamoDB
#
variable "dynamodb_billing_mode" {
  type        = string
  description = "DynamoDB billing mode"
  default     = "PAY_PER_REQUEST"
}
