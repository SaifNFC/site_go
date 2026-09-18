COMPOSE = docker compose -f deployments/docker-compose.yml
TEMPL = $(shell go env GOPATH)/bin/templ
K8S = deployments/kubernetes

.PHONY: run build test lint up down logs logs-sync ps templ templ-watch \
	k8s-image k8s-image-sync k8s-image-operator k8s-apply k8s-delete k8s-status k8s-logs k8s-port-forward \
	k8s-status-sync k8s-logs-sync k8s-status-operator k8s-logs-operator k8s-fix-coredns \
	ansible-setup ansible-deploy \
	semaphore-setup semaphore-up semaphore-down semaphore-logs

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

up: ## Build + lance l'API et le microservice tmdb-sync en conteneurs (arrière-plan)
	$(COMPOSE) up --build -d

down: ## Arrête et nettoie les conteneurs
	$(COMPOSE) down

logs: ## Suit les logs des conteneurs api + tmdb-sync
	$(COMPOSE) logs -f

logs-sync: ## Suit uniquement les logs du conteneur tmdb-sync
	$(COMPOSE) logs -f tmdb-sync

ps: ## Liste les conteneurs du projet
	$(COMPOSE) ps

k8s-image: ## Build l'image de l'API directement dans le runtime containerd de Minikube
	minikube image build -t letterboxd-clone:latest -f deployments/docker/Dockerfile .

k8s-image-sync: ## Build l'image du microservice tmdb-sync directement dans le runtime containerd de Minikube
	minikube image build -t tmdb-sync:latest -f deployments/docker/Dockerfile.tmdb-sync .

k8s-image-operator: ## Build l'image du filmsync-operator directement dans le runtime containerd de Minikube
	minikube image build -t filmsync-operator:latest -f deployments/docker/Dockerfile.filmsync-operator .

k8s-apply: ## Déploie l'API, tmdb-sync et le filmsync-operator (configmaps/secret/deployments/services/CRD/RBAC) sur le cluster courant
	@test -f $(K8S)/secret.yaml || (echo "Manque $(K8S)/secret.yaml — copie $(K8S)/secret.example.yaml et remplis-le" && exit 1)
	kubectl apply -f $(K8S)/configmap.yaml -f $(K8S)/secret.yaml -f $(K8S)/deployment.yaml -f $(K8S)/service.yaml
	kubectl apply -f $(K8S)/tmdb-sync-configmap.yaml -f $(K8S)/tmdb-sync-deployment.yaml -f $(K8S)/tmdb-sync-service.yaml
	kubectl apply -f $(K8S)/filmsync-crd.yaml
	kubectl apply -f $(K8S)/operator-rbac.yaml -f $(K8S)/filmsync-operator-deployment.yaml

k8s-delete: ## Supprime les ressources déployées (API + tmdb-sync + filmsync-operator)
	kubectl delete -f $(K8S)/deployment.yaml -f $(K8S)/service.yaml -f $(K8S)/configmap.yaml -f $(K8S)/secret.yaml --ignore-not-found
	kubectl delete -f $(K8S)/tmdb-sync-deployment.yaml -f $(K8S)/tmdb-sync-service.yaml -f $(K8S)/tmdb-sync-configmap.yaml --ignore-not-found
	kubectl delete -f $(K8S)/operator-rbac.yaml -f $(K8S)/filmsync-operator-deployment.yaml --ignore-not-found
	kubectl delete -f $(K8S)/filmsync-crd.yaml --ignore-not-found

k8s-status: ## Affiche l'état des pods/services de l'API
	kubectl get deployment,pod,svc -l app=letterboxd-api

k8s-status-sync: ## Affiche l'état des pods/services de tmdb-sync
	kubectl get deployment,pod,svc -l app=tmdb-sync

k8s-status-operator: ## Affiche l'état du pod/deployment de filmsync-operator, et les FilmSync existants
	kubectl get deployment,pod -l app=filmsync-operator
	kubectl get filmsyncs

k8s-logs: ## Suit les logs du pod de l'API
	kubectl logs -l app=letterboxd-api -f

k8s-logs-sync: ## Suit les logs du pod tmdb-sync
	kubectl logs -l app=tmdb-sync -f

k8s-logs-operator: ## Suit les logs du pod filmsync-operator
	kubectl logs -l app=filmsync-operator -f

k8s-port-forward: ## Expose le service API en local sur http://localhost:8081
	kubectl port-forward svc/letterboxd-api 8081:80

k8s-fix-coredns: ## Corrige CoreDNS sur Minikube (proxy DNS Docker/VPN qui timeout, cf. CLAUDE.md) — à relancer après un minikube delete
	kubectl apply -f $(K8S)/coredns-patch.yaml
	kubectl rollout restart deployment -n kube-system coredns
	kubectl rollout status deployment -n kube-system coredns --timeout=60s

ansible-setup: ## Installe ansible/python3-kubernetes (sudo, 1 fois) + la collection kubernetes.core
	sudo apt install -y ansible python3-kubernetes
	ansible-galaxy collection install -r ansible/requirements.yml

ansible-deploy: ## Déploie l'API + tmdb-sync via Ansible (kubernetes.core.k8s)
	ansible-playbook ansible/deploy.yml

semaphore-setup: ## Installe le binaire Semaphore (.deb, sudo, 1 fois) — lance ensuite le wizard toi-même
	curl -fsSL -o /tmp/semaphore.deb https://github.com/semaphoreui/semaphore/releases/download/v2.19.12/semaphore_community_2.19.12_linux_amd64.deb
	sudo apt install -y /tmp/semaphore.deb
	rm /tmp/semaphore.deb
	mkdir -p ansible/semaphore
	@echo "Lance maintenant : cd ansible/semaphore && semaphore setup"
	@echo "(pas via make — le wizard interactif bug avec le buffering stdin de make)"
	@echo "Choisis SQLite (option 4) comme base. Puis crée l'admin avec :"
	@echo "  semaphore user add --config=config.json --admin --login=<toi> --email=<toi> --name=<toi> --password=<toi>"

semaphore-up: ## Démarre Semaphore comme service utilisateur sur http://localhost:3000
	systemctl --user daemon-reload
	systemctl --user enable --now semaphore.service

semaphore-down: ## Arrête le service Semaphore
	systemctl --user disable --now semaphore.service

semaphore-logs: ## Suit les logs du service Semaphore
	journalctl --user -u semaphore.service -f
