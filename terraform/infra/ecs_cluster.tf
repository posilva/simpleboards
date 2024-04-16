# ECS Cluster using fargate or ec2
module "ecs_cluster" {
  source  = "cloudposse/ecs-cluster/aws"
  version = "0.6.1"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  container_insights_enabled = true
  # TODO: disabled for now as ec2 provider gives more control to debug
  capacity_providers_fargate      = false
  capacity_providers_fargate_spot = false
  capacity_providers_ec2 = {
    default = {
      associate_public_ip_address = false
      instance_type               = "t3.medium"
      security_group_ids          = [module.cluster_nodes_sg.id]
      subnet_ids                  = module.subnets.private_subnet_ids
      min_size                    = 0
      max_size                    = 2
    }
  }
}

module "cluster_nodes_sg" {
  source  = "cloudposse/security-group/aws"
  version = "2.2.0"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  # Security Group names must be unique within a VPC.
  # This module follows Cloud Posse naming conventions and generates the name
  # based on the inputs to the null-label module, which means you cannot
  # reuse the label as-is for more than one security group in the VPC.
  #
  # Here we add an attribute to give the security group a unique name.
  attributes = ["cluster-nodes"]

  # Allow unlimited egress
  allow_all_egress = true

  create_before_destroy = true

  rules_map = {
    ingress = [
      {
        key                      = "intra_cluster"
        type                     = "ingress"
        source_security_group_id = module.alb.security_group_id
        from_port                = 80 # TODO: should be configurable
        to_port                  = 80
        protocol                 = "tcp"
        cidr_blocks              = []
        self                     = null
        description              = "Allow HTTP from inside the security group"
      },
      {
        key         = "http_inside"
        type        = "ingress"
        from_port   = 80
        to_port     = 80
        protocol    = "tcp"
        cidr_blocks = []
        self        = true
        description = "Allow HTTP from inside the security group"
      }
    ]
  }
  vpc_id = module.vpc.vpc_id

}

module "alb" {
  source  = "cloudposse/alb/aws"
  version = "1.11.1"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  security_group_enabled = true

  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.subnets.public_subnet_ids
  internal        = false
  ip_address_type = "ipv4"

  http_enabled  = true
  http2_enabled = false
  http_redirect = false
  https_enabled = false
  https_port    = 443

  cross_zone_load_balancing_enabled = true

  deletion_protection_enabled = false
  deregistration_delay        = 15

  # health check settings
  # TODO: the health check may need to be adjusted
  health_check_path                = "/"
  health_check_timeout             = var.health_check_timeout
  health_check_healthy_threshold   = var.health_check_healthy_threshold
  health_check_unhealthy_threshold = var.health_check_unhealthy_threshold
  health_check_interval            = var.health_check_interval
  health_check_matcher             = var.health_check_matcher

  access_logs_enabled = false
}
