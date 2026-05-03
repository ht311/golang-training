terraform {
  required_version = ">= 1.5"

  required_providers {
    fly = {
      source  = "fly-apps/fly"
      version = "~> 0.2"
    }
    neon = {
      source  = "kislerdm/neon"
      version = "~> 0.6"
    }
    vercel = {
      source  = "vercel/vercel"
      version = "~> 2.0"
    }
  }
}

provider "fly" {
  # FLY_API_TOKEN env var
}

provider "neon" {
  # NEON_API_KEY env var
}

provider "vercel" {
  # VERCEL_API_TOKEN env var
}
