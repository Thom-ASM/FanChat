resource "aws_s3_bucket" "manifest" {
  bucket = "${var.project_name}-manifest"
}


data "aws_iam_policy_document" "manifest_cloudfront" {
  statement {
    sid    = "AllowCloudFrontRead"
    effect = "Allow"

    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.manifest.arn}/*",
    ]

    principals {
      type = "Service"

      identifiers = [
        "cloudfront.amazonaws.com",
      ]
    }

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"

      values = [
        aws_cloudfront_distribution.manifest.arn,
      ]
    }
  }
}

resource "aws_s3_bucket_policy" "manifest" {
  bucket = aws_s3_bucket.manifest.id
  policy = data.aws_iam_policy_document.manifest_cloudfront.json
}
