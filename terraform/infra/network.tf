# Network configuration 

data "aws_availability_zones" "available" {
  state = "available"
}

module "vpc" {
  source  = "cloudposse/vpc/aws"
  version = "2.2.0"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment


  ipv4_primary_cidr_block = var.vpc_cidr

  assign_generated_ipv6_cidr_block = false
}

module "subnets" {
  source  = "cloudposse/dynamic-subnets/aws"
  version = "2.4.2"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment


  availability_zones      = slice(data.aws_availability_zones.available.names, 0, 2)
  vpc_id                  = module.vpc.vpc_id
  igw_id                  = [module.vpc.igw_id]
  ipv4_cidr_block         = [module.vpc.vpc_cidr_block]
  ipv4_enabled            = true
  ipv6_enabled            = false
  ipv6_egress_only_igw_id = [module.vpc.ipv6_egress_only_igw_id]
  ipv6_cidr_block         = [module.vpc.vpc_ipv6_cidr_block]
  nat_gateway_enabled     = true
  nat_instance_enabled    = false
  route_create_timeout    = "5m"
  route_delete_timeout    = "10m"
  max_subnet_count        = 3
}
