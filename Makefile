COMPOSE = docker compose -f deployments/docker-compose.yml
TEMPL = $(shell go env GOPATH)/bin/templ

.PHONY: run build test lint up down logs ps templ templ-watch

templ: ## Génère le code Go à partir des fichiers .templ
	$(TEMPL) generate

templ-watch: ## Régénère les .templ à la volée pendant le dev
	$(TEMPL) generate --watch

run: templ ## Lance l'API en local (go run)
	go run cmd/api/main.go

build: templ ## Compile le binaire en local
	go build -o bin/api cmd/api/main.go

test: ## Lance les tests
	go test ./...

lint: ## Lint le code
	golangci-lint run

up: ## Build + lance l'API en conteneur (arrière-plan)
	$(COMPOSE) up --build -d

down: ## Arrête et nettoie les conteneurs
	$(COMPOSE) down

logs: ## Suit les logs du conteneur api
	$(COMPOSE) logs -f api

ps: ## Liste les conteneurs du projet
	$(COMPOSE) ps
