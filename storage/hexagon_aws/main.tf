variable "aws_region" {
  description = "AWS region."
  type        = string
}

variable "bucket_subdomain" {
  description = "subdomain of the name of s3 bucket to create"
  type        = string
}

variable "bucket_domain" {
  description = "domain of the name of s3 bucket to create"
  type        = string
}

provider "aws" {
  region = var.aws_region
}

resource "aws_s3_bucket" "bucket" {
  bucket = "${var.bucket_subdomain}.${var.bucket_domain}"

  tags = {
    Name        = var.bucket_subdomain
    Description = "Bucket to store results of metrics processor job."
  }
}

output "bucket_url" {
  description = "The URL of the bucket"
  value       = "s3://${aws_s3_bucket.bucket.bucket_domain_name}"
}
