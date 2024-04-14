locals {
  region      = var.region
  name        = var.name
  namespace   = var.namespace
  stage       = var.stage
  environment = var.environment

  tags = {
    Name       = local.name
    Owner      = "posilva@gmail.com"
    Repository = "https://github.com/posilva/${local.name}"
  }
}

