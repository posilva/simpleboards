provider "aws" {
  region = local.region

  #--Standard Provider Configurations
  default_tags {
    tags = {
      Environment = "dev"
    }
  }

  #   #--Mocked Credentials
  #   access_key = "test"
  #   secret_key = "test"


  #   #--These settings allow for authentication and other validations which are enforced
  #   #--in the AWS provider to be bypassed by Localstack.
  #   skip_credentials_validation = true
  #   skip_metadata_api_check     = true
  #   skip_requesting_account_id  = true

  #   #--Redirect Service Endpoints to Localstack. Whilst we won't be any of these it's good
  #   #--to see how they work and one should be specified to avoid rogue creations.
  #   #--See the Localstack docs for a full list of suitable endpoints):
  #   #--https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/custom-service-endpoints#available-endpoint-customizations
  #   endpoints {
  #     dynamodb = "http://localhost:4566"
  #   }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.45.0"
    }
  }
}

// Used by get the current aws number account.
data "aws_caller_identity" "current" {}
