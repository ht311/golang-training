.PHONY: generate dev-be dev-fe dev up down tf-init tf-plan tf-apply lint

# --- Code generation ---
generate: generate-openapi generate-go

generate-openapi:
	cd api/typespec && npx tsp compile .
	cp api/typespec/tsp-output/@typespec/openapi/openapi.yaml api/openapi/openapi.yaml

generate-go:
	cd backend && ~/go/bin/oapi-codegen --config oapi-codegen.cfg.yaml ../api/openapi/openapi.yaml

# --- Local development ---
dev-be:
	cd backend && DATABASE_URL=$${DATABASE_URL} go run ./cmd/server

dev-fe:
	cd frontend && npm run dev

# --- Docker Compose ---
up:
	docker compose up --build

down:
	docker compose down

# --- Build ---
build-be:
	cd backend && go build -o bin/server ./cmd/server

build-fe:
	cd frontend && npm run build

# --- Lint & Format ---
lint:
	cd backend && go vet ./...
	cd frontend && npm run lint

# --- Terraform ---
tf-init:
	cd infrastructure && terraform init

tf-plan:
	cd infrastructure && terraform plan

tf-apply:
	cd infrastructure && terraform apply
