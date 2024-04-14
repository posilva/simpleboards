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
  default     = "sgs"
}

variable "stage" {
  type        = string
  description = "The name or role of the account the resource is for, such as prod or dev"
  default     = "dev"
}

variable "environment" {
  type        = string
  description = "A short abbreviation for the AWS region hosting the resource, or gbl for resources like IAM roles that have no region"
  default     = "gbl"
}

# service's name
variable "name" {
  type        = string
  description = "Service name"
  default     = "simpleboards"
}

# DynamoDB
#
variable "dynamodb_billing_mode" {
  type        = string
  description = "DynamoDB billing mode"
  default     = "PAY_PER_REQUEST"
}


# Network 
# 
variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR"
  default     = "10.1.0.0/16"
}

# Service 
# 

variable "health_check_grace_period_seconds" {
  type        = number
  description = "Health check grace period seconds"
  default     = 39
}

variable "deployment_minimum_healthy_percent" {
  type        = number
  description = "Deployment minimum healthy percent"
  default     = 100
}

variable "deployment_maximum_percent" {
  type        = number
  description = "Deployment maximum percent"
  default     = 200
}

variable "deployment_controller_type" {
  type        = string
  description = "Deployment controller type"
  default     = "ECS"
}

variable "task_memory" {
  type        = number
  description = "Task memory"
  default     = 512
}

variable "task_cpu" {
  type        = number
  description = "Task cpu"
  default     = 256
}

variable "health_check_timeout" {
  type        = number
  description = "Health check timeout"
  default     = 10
}

variable "health_check_healthy_threshold" {
  type        = number
  description = "Health check healthy threshold"
  default     = 3
}

variable "health_check_unhealthy_threshold" {
  type        = number
  description = "Health check unhealthy threshold"
  default     = 3
}

variable "health_check_interval" {
  type        = number
  description = "Health check interval"
  default     = 15
}

variable "health_check_matcher" {
  type        = string
  description = "Health check matcher"
  default     = "200-399"
}

variable "target_group_port" {
  type        = number
  description = "Target group port"
  default     = 80
}
