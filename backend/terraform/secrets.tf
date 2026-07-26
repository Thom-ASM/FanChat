resource "aws_secretsmanager_secret" "youtube_oauth" {
  name = "${var.project_name}/youtube-oauth"
}

resource "aws_secretsmanager_secret" "youtube_ingest" {
  name = "${var.project_name}/youtube-ingest"
}
