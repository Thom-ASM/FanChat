output "manifest_url" {
  description = "Public URL for the FanChat manifest"
  value       = "https://${aws_cloudfront_distribution.manifest.domain_name}/manifest.json"
}

output "orchestrator_function_name" {

  description = "FanChat orchestrator Lambda function name"
  value       = aws_lambda_function.orchestrator.function_name
}

output "orchestrator_function_arn" {
  description = "FanChat orchestrator Lambda ARN"
  value       = aws_lambda_function.orchestrator.arn
}
