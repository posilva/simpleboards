module "ecr_github_role" {
  source  = "cloudposse/iam-role/aws"
  version = "0.17.0"

  enabled     = true
  namespace   = local.namespace
  environment = local.environment
  stage       = local.stage
  name        = "${local.name}_github"

  policy_description = "${local.name} Github Actions Setup Role "
  role_description   = "Used by github repos to manage infra"

  principals = {
    Federated = [
      "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/token.actions.githubusercontent.com"
    ]
  }

  assume_role_actions = ["sts:AssumeRoleWithWebIdentity", "sts:TagSession", "sts:AssumeRole"]

  assume_role_conditions = [
    {
      test     = "ForAllValues:StringEquals"
      variable = "token.actions.githubusercontent.com:iss"
      values   = ["https://token.actions.githubusercontent.com"]
    },
    {
      test     = "ForAllValues:StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    },
    {
      test     = "ForAnyValue:StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:posilva/${local.name}:*"]
    }
  ]

  policy_documents = [
    data.aws_iam_policy_document.ecr_policy.json
  ]

  tags_enabled = true
}

