locals {
  orchestrator_function_name = "${var.project_name}-orchestrator"
  orchestrator_zip           = "${path.module}/../orchestrator/dist/orchestrator.zip"
}

data "aws_iam_policy_document" "orchestrator_assume_role" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

resource "aws_iam_role" "orchestrator" {
  name               = "${var.project_name}-orchestrator-lambda"
  assume_role_policy = data.aws_iam_policy_document.orchestrator_assume_role.json
}

resource "aws_iam_role_policy_attachment" "orchestrator_basic_execution" {
  role       = aws_iam_role.orchestrator.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "orchestrator_permissions" {
  statement {
    sid    = "ReadYouTubeOAuthSecret"
    effect = "Allow"

    actions = [
      "secretsmanager:GetSecretValue",
    ]

    resources = [
      aws_secretsmanager_secret.youtube_oauth.arn,
    ]
  }
}

resource "aws_iam_role_policy" "orchestrator_permissions" {
  name   = "${var.project_name}-orchestrator-permissions"
  role   = aws_iam_role.orchestrator.id
  policy = data.aws_iam_policy_document.orchestrator_permissions.json
}

resource "aws_cloudwatch_log_group" "orchestrator" {
  name              = "/aws/lambda/${local.orchestrator_function_name}"
  retention_in_days = 14
}

resource "aws_lambda_function" "orchestrator" {
  function_name = local.orchestrator_function_name
  description   = "Orchestrates FanChat YouTube broadcasts and publishers"

  filename         = local.orchestrator_zip
  source_code_hash = filebase64sha256(local.orchestrator_zip)

  role    = aws_iam_role.orchestrator.arn
  runtime = "provided.al2023"
  handler = "bootstrap"

  architectures = [
    "arm64",
  ]

  memory_size = 256
  timeout     = 60

  environment {
    variables = {
      YOUTUBE_OAUTH_SECRET_ARN = aws_secretsmanager_secret.youtube_oauth.arn
      ECS_CLUSTER_ARN          = aws_ecs_cluster.encoder_ecs_cluster.arn
      ECS_TASK_DEFINITION_ARN  = aws_ecs_task_definition.encoder_service.arn
      SCHEDULER_ROLE_ARN       = aws_iam_role.scheduler_role.arn
      YOUTUBE_STREAM_TABLE_NAME = aws_dynamodb_table.youtube_stream_table.name
    }
  }

  depends_on = [
    aws_iam_role_policy_attachment.orchestrator_basic_execution,
    aws_iam_role_policy.orchestrator_permissions,
    aws_cloudwatch_log_group.orchestrator,
  ]
}
