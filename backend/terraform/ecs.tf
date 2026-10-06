
resource "aws_ecs_cluster" "encoder_ecs_cluster" {
  name = "encoder_ecs_cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_task_definition" "encoder_service" {
  family                   = "encoder-service"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"

  cpu    = "512"
  memory = "1024"

  container_definitions = jsonencode([
    {
      name      = "encoder"
      image     = aws_ecr_repository.encoder_repository.repository_url
      essential = true
    }
  ])

  execution_role_arn = aws_iam_role.ecs_execution_role.arn
  task_role_arn      = aws_iam_role.ecs_task_role.arn
}
