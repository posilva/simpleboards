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

data "aws_iam_policy_document" "dynamodb-full-access" {
  statement {
    effect  = "Allow"
    actions = ["dynamodb:*"]
    resources = [
      "${module.dynamodb_table.table_arn}"
    ]
  }
}
resource "aws_iam_policy" "dynamodb-full-access" {
  name        = "dynamodb-all-access"
  description = "Dynamodb All Access"
  policy      = data.aws_iam_policy_document.dynamodb-full-access.json
}
