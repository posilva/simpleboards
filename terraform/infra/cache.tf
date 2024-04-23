
module "redis" {
  source  = "cloudposse/elasticache-redis/aws"
  version = "1.2.2"
  enabled = false

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  availability_zones         = module.subnets.availability_zones
  vpc_id                     = module.vpc.vpc_id
  allowed_security_group_ids = [module.cluster_nodes_sg.id]
  subnets                    = module.subnets.private_subnet_ids
  cluster_size               = var.cache_cluster_size
  instance_type              = var.cache_instance_type
  apply_immediately          = true
  automatic_failover_enabled = false
  engine_version             = var.cache_engine_version
  family                     = var.cache_family
  at_rest_encryption_enabled = var.cache_at_rest_encryption_enabled
  transit_encryption_enabled = var.cache_transit_encryption_enabled

  parameter = [
    {
      name  = "notify-keyspace-events"
      value = "lK"
    }
  ]
}
