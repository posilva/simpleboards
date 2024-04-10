module "dynamodb_table" {
  source  = "cloudposse/dynamodb/aws"
  version = "0.32.0"

  namespace    = var.namespace
  stage        = var.stage
  environment  = var.environment
  name         = "leaderboards"
  hash_key     = "pk"
  range_key    = "sk"
  billing_mode = var.dynamodb_billing_mode

  dynamodb_attributes = [
    {
      name = "pk"
      type = "S"
    },
    {
      name = "sk"
      type = "S"
    }
  ]
}
