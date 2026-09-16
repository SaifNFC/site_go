# CLAUDE.md
Ce fichier fournit du contexte à Claude Code pour travailler efficacement sur ce projet.

## Vue d'ensemble du projet
Clone de Letterboxd (films uniquement) : catalogue de films, notes, watchlist, listes personnalisées. 
Projet personnel dans une optique de montée en compétences Go.

Objectifs pédagogiques du projet (à garder en tête pour les suggestions) :
- Pratiquer Go/Gin/GORM en dehors du contexte pro (API interop chez Free)
- Découvrir PostgreSQL en profondeur (au-delà de ce que GORM abstrait)
- Construire un cycle DevOps complet : Docker, CI, Kubernetes, observabilité

## Stack technique

| Couche | Techno |
|---|---|
| Langage | Go |
| Framework web | Gin |
| Front-end | `templ` (composants Go compilés, pas de JS framework) + Pico.css (classless, auto-hébergé, pas de CDN) |
| ORM | GORM (driver PostgreSQL — `gorm.io/driver/postgres`) |
| Base de données | PostgreSQL (hébergé sur Neon, pas de conteneur DB local) |
| Auth | JWT |
| API externe | TMDB (peuplement du catalogue de films) |
| Conteneurisation | Docker (multi-stage build, image finale distroless non-root) |
| Architecture | API monolithique Gin + microservice `tmdb-sync` (test du pattern microservices sur K8s : service indépendant, appelé via DNS interne `http://tmdb-sync`) |
| Orchestration | Kubernetes (Minikube en local, driver Docker / runtime containerd) |
| Automatisation déploiement | Ansible (module `kubernetes.core.k8s`, applique les manifests K8s) + UI web Ansible Semaphore |
| CI | GitHub Actions ou GitLab CI + `golangci-lint` (pas encore fait) |
| Observabilité | Logs structurés, healthcheck `/health` ; métriques Prometheus/Grafana (pas encore fait) |

## État d'avancement (à tenir à jour)

**Fait :**
- API Gin complète : auth JWT, CRUD films, notes, watchlist (couches repository/services/handlers)
- Front HTML server-side (`templ` + Pico.css) : liste des films (`/`) et détail (`/films/:id/view`)
- Docker multi-stage (build `templ generate` + binaire + assets statiques dans l'image finale)
- Manifests Kubernetes (`deployments/kubernetes/`) testés sur Minikube : Deployment (probes, resources, securityContext non-root avec `runAsUser: 65532`), Service ClusterIP, ConfigMap, Secret (pattern `secret.example.yaml` committé / `secret.yaml` gitignoré, comme `.env`/`.env.example`)
- Microservice `tmdb-sync` (`cmd/tmdb-sync`) : extraction du sync TMDB en second service Go/Gin autonome (`internal/services/sync_service.go`, upsert TMDB→DB), déployé indépendamment (Dockerfile dédié, ConfigMap/Deployment/Service K8s dédiés, réutilise le Secret existant), appelé depuis l'API principale (`POST /admin/films/:tmdb_id/sync`, protégé JWT) via un client HTTP interne (`internal/tmdbsync`) pointant sur `TMDB_SYNC_URL` — testé en local (binaires natifs), via docker-compose (service `tmdb-sync` dédié, DNS interne du réseau compose) et **sur Minikube** (DNS interne K8s `http://tmdb-sync`, confirmé par les logs des deux pods)
- `make k8s-image`/`k8s-image-sync` utilisent `minikube image build` (pas `docker build` via `minikube docker-env`, incompatible avec le runtime containerd de Minikube sur cette machine)
- Ansible (`ansible/deploy.yml`) : playbook idempotent qui applique les 7 manifests K8s via `kubernetes.core.k8s` (paquets système `ansible`/`python3-kubernetes` + collection galaxy `kubernetes.core`, cf. `ansible/requirements.yml`) — équivalent de `make k8s-apply`, sans le build d'image (volontairement hors scope)
- Ansible Semaphore (UI web pour lancer/suivre les playbooks) : binaire installé via `.deb` officiel, config/DB SQLite dans `ansible/semaphore/` (gitignoré), tourne en service `systemd --user` sur `http://localhost:3000` — Project/Repository (chemin local)/Inventory/Template configurés et testés avec succès depuis le navigateur
- Repo poussé sur GitHub (`SaifNFC/site_go`, compte perso — identité git configurée en local au repo uniquement, cf. machine pro avec config GitLab globale)

**Pas fait :**
- Résilience de l'appel API principale → `tmdb-sync` : pas de retry/timeout configurable ni de circuit breaker (une panne de `tmdb-sync` remonte en 500 sèche)
- NetworkPolicy restreignant l'accès à `tmdb-sync` (aujourd'hui joignable par n'importe quel pod du cluster)
- Tests unitaires sur la couche services (seulement un test sur le client TMDB, rien sur `SyncService`)
- Pages HTML login/register (JSON endpoints déjà là, pas de vue)
- CI (lint/test/build automatisés)
- Ansible : build d'image et démarrage Minikube restent hors playbook (manuel via `make k8s-image*`), pas de rôle dédié pour templater `secret.yaml`
- Prometheus/Grafana

## Structure du repo
```
.
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── tmdb-sync/         # microservice de sync TMDB, service Gin autonome
│       └── main.go
├── internal/
│   ├── config/          # chargement config (env vars, .env), partagé par les deux binaires
│   ├── database/        # connexion GORM, migrations (migrations lancées uniquement par l'API principale)
│   ├── models/           # entités GORM (User, Film, Note, Watchlist)
│   ├── handlers/          # handlers Gin (controllers) + page_handler.go (vues HTML) + sync_handler.go
│   ├── middleware/        # auth JWT, logging, recovery
│   ├── repository/        # couche accès données (si séparée des handlers)
│   ├── services/           # logique métier (ex: sync_service.go, utilisé par tmdb-sync)
│   ├── tmdb/               # client API TMDB (utilisé par tmdb-sync)
│   ├── tmdbsync/           # client HTTP interne API principale -> microservice tmdb-sync
│   └── views/              # composants templ (.templ + .go générés, gitignorés)
├── web/
│   └── static/            # CSS (Pico.css auto-hébergé) servi via router.Static
├── migrations/            # si golang-migrate utilisé en plus d'AutoMigrate (pas encore fait)
├── deployments/
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── Dockerfile.tmdb-sync
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── configmap.yaml
│   │   ├── secret.example.yaml  # template committé
│   │   ├── secret.yaml          # vraies valeurs, gitignoré
│   │   ├── tmdb-sync-deployment.yaml
│   │   ├── tmdb-sync-service.yaml
│   │   └── tmdb-sync-configmap.yaml
│   └── docker-compose.yml   # service api uniquement ; tmdb-sync se lance en local (host.docker.internal)
├── ansible/
│   ├── requirements.yml   # dépendance collection galaxy kubernetes.core
│   ├── deploy.yml         # playbook : applique les manifests K8s (équivalent make k8s-apply)
│   └── semaphore/         # binaire/config/DB Semaphore (gitignoré, machine-local)
├── http/                  # fichiers .http pour tester manuellement chaque domaine (dont sync.http)
├── .github/workflows/ (ou .gitlab-ci.yml)  # pas encore fait
├── go.mod
├── go.sum
└── CLAUDE.md
```

## Modèle de données (v1)
- **users** : id, email, password_hash, created_at
- **films** : id, tmdb_id, titre, annee, poster_url, synopsis
- **notes** : user_id, film_id, note, commentaire, created_at (many-to-many enrichi)
- **watchlists** : user_id, film_id, added_at

## Conventions de code
- Style Go idiomatique standard (`gofmt`, `golangci-lint` avant tout commit)
- Gestion d'erreurs explicite (pas de `panic` en dehors de l'init), erreurs wrappées avec `fmt.Errorf("...: %w", err)`
- Handlers Gin fins : la logique métier ne doit pas vivre dans les handlers mais dans une couche services/repository
- Variables d'environnement pour toute config sensible (DSN PostgreSQL, clé API TMDB, secret JWT) — jamais en dur dans le code
- Tests unitaires pour la logique métier (services), tests d'intégration pour les endpoints critiques si le temps le permet


## Commandes utiles (voir aussi le Makefile, source de vérité)
```bash
# Lancer l'API en local (génère les .templ automatiquement)
make run

# Lint / Tests
make lint
make test

# Docker Compose (API, DB distante sur Neon)
make up / make down / make logs

# Kubernetes (sur un cluster Minikube déjà démarré : `minikube start --driver=docker`)
make k8s-image          # build l'image de l'API dans le docker daemon de Minikube
make k8s-image-sync     # build l'image du microservice tmdb-sync
make k8s-apply          # applique configmaps/secret/deployments/services (API + tmdb-sync)
make k8s-status         # état des pods/deployment/service de l'API
make k8s-status-sync    # état des pods/deployment/service de tmdb-sync
make k8s-logs           # logs du pod de l'API
make k8s-logs-sync      # logs du pod tmdb-sync
make k8s-port-forward   # expose le service API sur http://localhost:8081

# Ansible (déploiement K8s déclaratif, alternative à k8s-apply)
make ansible-setup      # sudo apt install ansible python3-kubernetes + collection kubernetes.core
make ansible-deploy     # ansible-playbook ansible/deploy.yml

# Semaphore (UI web pour lancer/suivre les playbooks Ansible)
make semaphore-setup    # installe le binaire (.deb, sudo) — puis lancer `semaphore setup` soi-même
make semaphore-up       # démarre le service systemd --user sur http://localhost:3000
make semaphore-down     # arrête le service
make semaphore-logs     # logs du service
```

## Notes pour Claude Code
- Le développeur a une bonne autonomie sur l'architecture Go/Gin/GORM au quotidien (contexte pro), donc privilégier des explications concises sur le "pourquoi" plutôt que sur les bases du langage.
- PostgreSQL et Kubernetes sont des découvertes volontaires (le stack pro utilise MariaDB, pas de K8s au quotidien) — être explicite sur le "pourquoi" de chaque décision/manifeste, ne pas supposer de connaissance implicite des concepts K8s (probes, securityContext, selectors/labels, etc.).
- Privilégier des étapes progressives et incrémentales plutôt que de générer de gros blocs de code d'un coup — l'objectif est la montée en compétence autant que la livraison.
- Gin reste le framework de référence, ne pas dévier vers Echo/Fiber/Chi sans raison explicite. Pas de framework JS pour le front — templ + CSS classless (Pico.css) uniquement, JS vanilla minimal si besoin.
- Machine de dev = poste pro, config git/SSH globale liée au GitLab de l'employeur. Ce repo pousse vers un compte GitHub personnel : toute config d'identité/auth git doit rester locale au repo (`git config --local`), jamais toucher au global. Ne jamais faire transiter un token/secret par le chat — le faire saisir par le développeur dans son propre terminal.
- Outils CLI (templ, minikube, kubectl) installés en user-space (`~/go/bin`, `~/.local/bin`). Le développeur a bien les droits `sudo` sur cette machine et `apt install`/paquets `.deb` sont OK pour des outils autonomes (ex: `ansible`, `python3-kubernetes`, Semaphore) — la vraie limite est de ne jamais toucher à la config GitLab/Docker/VPN pro (config git/SSH globale, conteneurs Docker d'autres projets pro qui tournent sur la même machine).