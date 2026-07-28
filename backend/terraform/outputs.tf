output "manifest_url" {
  description = "Public URL for the FanChat manifest"
  value       = "https://${aws_cloudfront_distribution.manifest.domain_name}/manifest.json"
}
