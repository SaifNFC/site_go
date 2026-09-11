COMPOSE = docker compose -f deployments/docker-compose.yml
TEMPL = $(shell go env GOPATH)/bin/templ
K8S = deployments/kubernetes

.PHONY: run build test lint up down logs ps templ templ-watch \
	k8s-image k8s-apply k8s-delete k8s-status k8s-logs k8s-port-forward

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

k8s-image: ## Build l'image dans le registre Docker de Minikube
	eval $$(minikube docker-env) && docker build -t letterboxd-clone:latest -f deployments/docker/Dockerfile .

k8s-apply: ## Déploie configmap/secret/deployment/service sur le cluster courant
	@test -f $(K8S)/secret.yaml || (echo "Manque $(K8S)/secret.yaml — copie $(K8S)/secret.example.yaml et remplis-le" && exit 1)
	kubectl apply -f $(K8S)/configmap.yaml -f $(K8S)/secret.yaml -f $(K8S)/deployment.yaml -f $(K8S)/service.yaml

k8s-delete: ## Supprime les ressources déployées
	kubectl delete -f $(K8S)/deployment.yaml -f $(K8S)/service.yaml -f $(K8S)/configmap.yaml -f $(K8S)/secret.yaml --ignore-not-found

k8s-status: ## Affiche l'état des pods/services
	kubectl get deployment,pod,svc -l app=letterboxd-api

k8s-logs: ## Suit les logs du pod
	kubectl logs -l app=letterboxd-api -f

k8s-port-forward: ## Expose le service en local sur http://localhost:8081
	kubectl port-forward svc/letterboxd-api 8081:80
