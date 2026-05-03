# --- Neon (PostgreSQL) ---
resource "neon_project" "main" {
  name      = var.app_name
  region_id = "aws-ap-northeast-1"
}

resource "neon_branch" "main" {
  project_id = neon_project.main.id
  name       = "main"
}

resource "neon_database" "blog" {
  project_id = neon_project.main.id
  branch_id  = neon_branch.main.id
  name       = "blog"
  owner_name = "blog_owner"
}

resource "neon_role" "blog_owner" {
  project_id = neon_project.main.id
  branch_id  = neon_branch.main.id
  name       = "blog_owner"
}

locals {
  database_url = "postgresql://${neon_role.blog_owner.name}:${neon_role.blog_owner.password}@${neon_project.main.database_host}/${neon_database.blog.name}?sslmode=require"
}

# --- Fly.io (Go backend) ---
resource "fly_app" "backend" {
  name = "${var.app_name}-backend"
  org  = "personal"
}

resource "fly_machine" "backend" {
  app    = fly_app.backend.name
  region = var.region
  name   = "${var.app_name}-backend-vm"

  image = "registry.fly.io/${fly_app.backend.name}:latest"

  services = [
    {
      ports = [
        { port = 443, handlers = ["tls", "http"] },
        { port = 80, handlers = ["http"] },
      ]
      protocol      = "tcp"
      internal_port = 8080
    }
  ]

  env = {
    PORT = "8080"
  }
}

resource "fly_secret" "database_url" {
  app   = fly_app.backend.name
  name  = "DATABASE_URL"
  value = local.database_url
}

# --- Vercel (Next.js frontend) ---
resource "vercel_project" "frontend" {
  name      = "${var.app_name}-frontend"
  framework = "nextjs"

  git_repository = {
    type = "github"
    repo = var.github_repo
  }

  root_directory = "frontend"
}

resource "vercel_project_environment_variable" "api_url" {
  project_id = vercel_project.frontend.id
  key        = "NEXT_PUBLIC_API_URL"
  value      = "https://${fly_app.backend.name}.fly.dev"
  target     = ["production", "preview"]
}
