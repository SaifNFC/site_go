# CLAUDE.md
Ce fichier fournit du contexte à Claude Code pour travailler efficacement sur ce projet.

## Vue d'ensemble du projet
Clone de Letterboxd (films uniquement) : catalogue de films, notes, watchlist, listes personnalisées. 
Projet personnel dans une optique de montée en compétences Go et de portfolio pour une recherche d'emploi Go développeur senior à Paris.

Objectifs pédagogiques du projet (à garder en tête pour les suggestions) :
- Pratiquer Go/Gin/GORM en dehors du contexte pro (API interop chez Free)
- Découvrir PostgreSQL en profondeur (au-delà de ce que GORM abstrait)
- Construire un cycle DevOps complet : Docker, CI, Kubernetes, observabilité
- Produire un projet démontrable en entretien technique

## Stack technique

| Couche | Techno |
|---|---|
| Langage | Go |
| Framework web | Gin |
| ORM | GORM (driver PostgreSQL — `gorm.io/driver/postgres`) |
| Base de données | PostgreSQL |
| Auth | JWT |
| API externe | TMDB (peuplement du catalogue de films) |
| Conteneurisation | Docker (multi-stage build) |
| Orchestration | Kubernetes (Kind pour CI, Minikube pour dev local) |
| CI | GitHub Actions ou GitLab CI + `golangci-lint` |
| Observabilité | Logs structurés, healthcheck `/health`, métriques Prometheus |

## Structure du repo (cible)
```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/          # chargement config (env vars, .env)
│   ├── database/        # connexion GORM, migrations
│   ├── models/           # entités GORM (User, Film, Note, Watchlist)
│   ├── handlers/          # handlers Gin (controllers)
│   ├── middleware/        # auth JWT, logging, recovery
│   ├── repository/        # couche accès données (si séparée des handlers)
│   ├── services/           # logique métier (ex: sync TMDB)
│   └── tmdb/               # client API TMDB
├── migrations/            # si golang-migrate utilisé en plus d'AutoMigrate
├── deployments/
│   ├── docker/
│   │   └── Dockerfile
│   ├── kubernetes/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── configmap.yaml
│   │   └── secret.yaml
│   └── docker-compose.yml
├── .github/workflows/ (ou .gitlab-ci.yml)
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


## Commandes utiles (à compléter au fur et à mesure)
```bash
# Lancer l'API en local
go run cmd/api/main.go

# Lint
golangci-lint run

# Tests
go test ./...

# Build Docker
docker build -t letterboxd-clone -f deployments/docker/Dockerfile .

# Docker Compose (API + PostgreSQL)
docker-compose -f deployments/docker-compose.yml up
```

## Notes pour Claude Code
- Le développeur a une bonne autonomie sur l'architecture Go/Gin/GORM au quotidien (contexte pro), donc privilégier des explications concises sur le "pourquoi" plutôt que sur les bases du langage.
- PostgreSQL est une découverte volontaire (le stack pro utilise MariaDB) — être explicite sur les différences PostgreSQL/MySQL quand elles impactent une décision technique (types, `RETURNING`, `ON CONFLICT`, etc.).
- Privilégier des étapes progressives et incrémentales plutôt que de générer de gros blocs de code d'un coup — l'objectif est la montée en compétence autant que la livraison.
- Toujours garder la stack alignée sur ce qui est demandé en entretien Go à Paris (Gin reste le framework de référence, ne pas dévier vers Echo/Fiber/Chi sans raison explicite).