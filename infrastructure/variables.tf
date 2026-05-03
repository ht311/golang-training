variable "app_name" {
  description = "Application name used across all services"
  type        = string
  default     = "golang-training-blog"
}

variable "region" {
  description = "Primary deployment region (Fly.io region code)"
  type        = string
  default     = "nrt"
}

variable "github_repo" {
  description = "GitHub repository for Vercel deployment (owner/repo)"
  type        = string
}
