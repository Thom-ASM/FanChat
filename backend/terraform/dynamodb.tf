
resource "aws_dynamodb_table" "youtube_stream_table" {
  name         = "youtube_stream_table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "showId"

  attribute {
    name = "showId"
    type = "S"
  }

  tags = {
    Name        = "dynamodb-table-1"
    Environment = "production"
  }
}
