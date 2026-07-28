locals {
  manifest_origin_id = "fan-chat-manifest-s3"
}

resource "aws_cloudfront_origin_access_control" "manifest" {
  name                              = "fan-chat-manifest-oac"
  description                       = "OAC for the FanChat manifest bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_cache_policy" "manifest" {
  name        = "fan-chat-manifest-cache"
  comment     = "Short-lived cache for FanChat manifests"
  min_ttl     = 0
  default_ttl = 60
  max_ttl     = 300

  parameters_in_cache_key_and_forwarded_to_origin {
    enable_accept_encoding_brotli = true
    enable_accept_encoding_gzip   = true

    cookies_config {
      cookie_behavior = "none"
    }

    headers_config {
      header_behavior = "none"
    }

    query_strings_config {
      query_string_behavior = "none"
    }
  }
}

resource "aws_cloudfront_response_headers_policy" "manifest" {
  name    = "fan-chat-manifest-cors"
  comment = "Allow the browser extension to fetch manifests"

  cors_config {
    access_control_allow_credentials = false
    origin_override                  = true

    access_control_allow_headers {
      items = ["*"]
    }

    access_control_allow_methods {
      items = ["GET", "HEAD"]
    }

    access_control_allow_origins {
      items = ["*"]
    }
  }
}

resource "aws_cloudfront_distribution" "manifest" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "FanChat public manifest"
  default_root_object = "manifest.json"
  price_class         = "PriceClass_100"

  origin {
    domain_name              = aws_s3_bucket.manifest.bucket_regional_domain_name
    origin_id                = local.manifest_origin_id
    origin_access_control_id = aws_cloudfront_origin_access_control.manifest.id
  }

  default_cache_behavior {
    target_origin_id = local.manifest_origin_id

    allowed_methods = [
      "GET",
      "HEAD",
    ]

    cached_methods = [
      "GET",
      "HEAD",
    ]

    cache_policy_id            = aws_cloudfront_cache_policy.manifest.id
    response_headers_policy_id = aws_cloudfront_response_headers_policy.manifest.id

    viewer_protocol_policy = "redirect-to-https"
    compress               = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
