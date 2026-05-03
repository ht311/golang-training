output "backend_url" {
  description = "Fly.io backend URL"
  value       = "https://${fly_app.backend.name}.fly.dev"
}

output "frontend_url" {
  description = "Vercel frontend URL"
  value       = "https://${vercel_project.frontend.name}.vercel.app"
}

output "neon_database_host" {
  description = "Neon PostgreSQL host"
  value       = neon_project.main.database_host
  sensitive   = true
}
