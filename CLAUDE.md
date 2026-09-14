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
| Orchestration | Kubernetes (Minikube en local, driver Docker / runtime containerd) |
| Automatisation déploiement | Ansible (prévu, pas encore fait — module `kubernetes.core.k8s`) |
| CI | GitHub Actions ou GitLab CI + `golangci-lint` (pas encore fait) |
| Observabilité | Logs structurés, healthcheck `/health` ; métriques Prometheus/Grafana (pas encore fait) |

## État d'avancement (à tenir à jour)

**Fait :**
- API Gin complète : auth JWT, CRUD films, notes, watchlist (couches repository/services/handlers)
- Front HTML server-side (`templ` + Pico.css) : liste des films (`/`) et détail (`/films/:id/view`)
- Docker multi-stage (build `templ generate` + binaire + assets statiques dans l'image finale)
- Manifests Kubernetes (`deployments/kubernetes/`) testés sur Minikube : Deployment (probes, resources, securityContext non-root avec `runAsUser: 65532`), Service ClusterIP, ConfigMap, Secret (pattern `secret.example.yaml` committé / `secret.yaml` gitignoré, comme `.env`/`.env.example`)
- Repo poussé sur GitHub (`SaifNFC/site_go`, compte perso — identité git configurée en local au repo uniquement, cf. machine pro avec config GitLab globale)

**Pas fait :**
- Service de sync TMDB (le client HTTP existe, pas d'orchestration pour peupler le catalogue)
- Tests unitaires sur la couche services (seulement un test sur le client TMDB)
- Pages HTML login/register (JSON endpoints déjà là, pas de vue)
- CI (lint/test/build automatisés)
- Ansible pour automatiser le déploiement K8s
- Prometheus/Grafana

## Structure du repo
```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/          # chargement config (env vars, .env)
│   ├── database/        # connexion GORM, migrations
│   ├── models/           # entités GORM (User, Film, Note, Watchlist)
│   ├── handlers/          # handlers Gin (controllers) + page_handler.go (vues HTML)
│   ├── middleware/        # auth JWT, logging, recovery
│   ├── repository/        # couche accès données (si séparée des handlers)
│   ├── services/           # logique métier (ex: sync TMDB)
│   ├── tmdb/               # client API TMDB
│   └── views/              # composants templ (.templ + .go générés, gitignorés)
├── web/
│   └── static/            # CSS (Pico.css auto-hébergé) servi via router.Static
├── migrations/            # si golang-migrate utilisé en plus d'AutoMigrate (pas encore fait)
├── deployments/
│   ├── docker/
│   │   └── Dockerfile
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── configmap.yaml
│   │   ├── secret.example.yaml  # template committé
│   │   └── secret.yaml          # vraies valeurs, gitignoré
│   └── docker-compose.yml
├── http/                  # fichiers .http pour tester manuellement chaque domaine
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
make k8s-image          # build l'image dans le docker daemon de Minikube
make k8s-apply          # applique configmap/secret/deployment/service
make k8s-status         # état des pods/deployment/service
make k8s-logs           # logs du pod
make k8s-port-forward   # expose le service sur http://localhost:8081
```

## Notes pour Claude Code
- Le développeur a une bonne autonomie sur l'architecture Go/Gin/GORM au quotidien (contexte pro), donc privilégier des explications concises sur le "pourquoi" plutôt que sur les bases du langage.
- PostgreSQL et Kubernetes sont des découvertes volontaires (le stack pro utilise MariaDB, pas de K8s au quotidien) — être explicite sur le "pourquoi" de chaque décision/manifeste, ne pas supposer de connaissance implicite des concepts K8s (probes, securityContext, selectors/labels, etc.).
- Privilégier des étapes progressives et incrémentales plutôt que de générer de gros blocs de code d'un coup — l'objectif est la montée en compétence autant que la livraison.
- Gin reste le framework de référence, ne pas dévier vers Echo/Fiber/Chi sans raison explicite. Pas de framework JS pour le front — templ + CSS classless (Pico.css) uniquement, JS vanilla minimal si besoin.
- Machine de dev = poste pro, config git/SSH globale liée au GitLab de l'employeur. Ce repo pousse vers un compte GitHub personnel : toute config d'identité/auth git doit rester locale au repo (`git config --local`), jamais toucher au global. Ne jamais faire transiter un token/secret par le chat — le faire saisir par le développeur dans son propre terminal.
- Outils CLI (templ, minikube, kubectl) installés en user-space (`~/go/bin`, `~/.local/bin`) sans droits root — pas de `apt install`/`sudo` sur cette machine.