
module "container_definition" {
  source  = "cloudposse/ecs-container-definition/aws"
  version = "0.61.1"

  container_name               = local.name
  container_image              = var.service_container_image
  container_memory             = 256
  container_cpu                = 256
  container_memory_reservation = 256
  essential                    = true
  readonly_root_filesystem     = false
  port_mappings = [
    {
      containerPort = var.target_group_port
      hostPort      = var.target_group_port
      protocol      = "tcp"
    }
  ]

}

module "ecs_service" {
  source  = "cloudposse/ecs-alb-service-task/aws"
  version = "0.74.0"

  namespace   = local.namespace
  stage       = local.stage
  name        = local.name
  environment = local.environment

  container_definition_json = module.container_definition.json_map_encoded_list

  ecs_cluster_arn = module.ecs_cluster.arn

  launch_type = "EC2"
  vpc_id      = module.vpc.vpc_id

  # TODO: this should be in a private subnet
  subnet_ids   = module.subnets.private_subnet_ids
  network_mode = "awsvpc"

  security_group_ids = [module.cluster_nodes_sg.id]

  ignore_changes_task_definition = false
  assign_public_ip               = false # only allowed for Fargate deployments
  propagate_tags                 = "SERVICE"

  desired_count = 1

  health_check_grace_period_seconds  = var.health_check_grace_period_seconds
  deployment_minimum_healthy_percent = var.deployment_minimum_healthy_percent
  deployment_maximum_percent         = var.deployment_maximum_percent
  deployment_controller_type         = var.deployment_controller_type
  task_memory                        = var.task_memory
  task_cpu                           = var.task_cpu
  depends_on                         = [module.ecs_cluster, module.alb]
  ecs_load_balancers = [{
    target_group_arn = module.alb.default_target_group_arn
    container_name   = local.name
    container_port   = var.target_group_port
  }]
}

