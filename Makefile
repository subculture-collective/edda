SHELL := /bin/bash
.SHELLFLAGS := -euo pipefail -c

ENV_FILE ?= .env
COMPOSE_FILE ?= docker-compose.yml
PUBLIC_URL ?= https://edda.subcult.tv
RELEASE_TAG ?= $(shell git rev-parse --short HEAD)
RUN_TIMESTAMP = $(shell date -u +%Y%m%dT%H%M%SZ)

.PHONY: help verify test frontend-build clean-runtime-artifacts \
	compose-config local-db-up local-db-down \
	app-build app-up app-deploy deploy app-status app-logs api-logs web-logs \
	db-backup migrate-prod smoke rollback-sim

help:
	@printf '%s\n' \
		'Common targets:' \
		'  make verify              Run backend tests and frontend production build' \
		'  make local-db-up         Start local Postgres profile and run local migrations' \
		'  make compose-config      Validate the canonical Docker Compose config' \
		'  make app-build           Build api/web images with RELEASE_TAG=<tag>' \
		'  make app-up              Recreate api/web using ENV_FILE=<file> and RELEASE_TAG=<tag>' \
		'  make deploy              Backup DB, run migrations, build, and recreate api/web' \
		'  make app-status          Show api/web container status' \
		'  make app-logs            Follow api/web logs' \
		'  make smoke               Run public production smoke checks' \
		'  make rollback-sim        Simulate rollback using latest deploy artifacts'

verify: test frontend-build

test:
	go test ./...

frontend-build:
	pnpm --dir frontend build

compose-config:
	@EDDA_RELEASE_TAG='$(RELEASE_TAG)' EDDA_ENV_FILE='$(ENV_FILE)' docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' config >/dev/null
	@echo 'compose config OK: $(COMPOSE_FILE)'

local-db-up:
	docker compose -f '$(COMPOSE_FILE)' --profile local up -d postgres
	task migrate

local-db-down:
	docker compose -f '$(COMPOSE_FILE)' --profile local down

app-build:
	EDDA_RELEASE_TAG='$(RELEASE_TAG)' EDDA_ENV_FILE='$(ENV_FILE)' docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' build api web

app-up:
	EDDA_RELEASE_TAG='$(RELEASE_TAG)' EDDA_ENV_FILE='$(ENV_FILE)' docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' up -d --no-build --force-recreate api web

app-deploy deploy:
	RUN_DIR=".sisyphus/evidence/nuc-deploy-$(RUN_TIMESTAMP)"; \
	mkdir -p "$$RUN_DIR"; \
	set -a; . './$(ENV_FILE)'; set +a; \
	export EDDA_RELEASE_TAG='$(RELEASE_TAG)'; \
	export EDDA_ENV_FILE='$(ENV_FILE)'; \
	echo "deploying Edda $$EDDA_RELEASE_TAG with env $(ENV_FILE)"; \
	echo "evidence: $$RUN_DIR"; \
	bash scripts/capture_prod_release.sh "$$EDDA_RELEASE_TAG" "$$RUN_DIR/rollback-manifest.env"; \
	bash scripts/db_backup.sh '$(ENV_FILE)' "$$RUN_DIR/pre-deploy.dump"; \
	bash scripts/run_prod_migrations.sh '$(ENV_FILE)' "$$RUN_DIR/goose-status.txt"; \
	docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' build api web; \
	docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' up -d --no-build --force-recreate api web; \
	docker inspect "$${EDDA_API_CONTAINER_NAME:-edda-api}" "$${EDDA_WEB_CONTAINER_NAME:-edda-web}" > "$$RUN_DIR/post-compose-inspect.json"; \
	docker ps --filter "name=^/$${EDDA_API_CONTAINER_NAME:-edda-api}$$" --filter "name=^/$${EDDA_WEB_CONTAINER_NAME:-edda-web}$$" --format '{{.Names}} {{.Image}} {{.Status}} {{.Ports}}' | tee "$$RUN_DIR/post-cutover-state.txt"

app-status:
	@set -a; . './$(ENV_FILE)'; set +a; \
	docker ps --filter "name=^/$${EDDA_API_CONTAINER_NAME:-edda-api}$$" --filter "name=^/$${EDDA_WEB_CONTAINER_NAME:-edda-web}$$" --format '{{.Names}} {{.Image}} {{.Status}} {{.Ports}}'

app-logs:
	EDDA_RELEASE_TAG='$(RELEASE_TAG)' EDDA_ENV_FILE='$(ENV_FILE)' docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' logs -f api web

api-logs:
	EDDA_RELEASE_TAG='$(RELEASE_TAG)' EDDA_ENV_FILE='$(ENV_FILE)' docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' logs -f api

web-logs:
	EDDA_RELEASE_TAG='$(RELEASE_TAG)' EDDA_ENV_FILE='$(ENV_FILE)' docker compose -f '$(COMPOSE_FILE)' --env-file '$(ENV_FILE)' logs -f web

db-backup:
	mkdir -p .sisyphus/evidence
	bash scripts/db_backup.sh '$(ENV_FILE)' ".sisyphus/evidence/pre-deploy-$(RUN_TIMESTAMP).dump"

migrate-prod:
	mkdir -p .sisyphus/evidence
	bash scripts/run_prod_migrations.sh '$(ENV_FILE)' ".sisyphus/evidence/goose-status-$(RUN_TIMESTAMP).txt"

smoke:
	bash scripts/smoke_prod_deploy.sh '$(PUBLIC_URL)'

rollback-sim:
	LATEST_RUN=$$(ls -1dt .sisyphus/evidence/nuc-deploy-* 2>/dev/null | head -n 1); \
	if [[ -z "$$LATEST_RUN" ]]; then echo 'no .sisyphus/evidence/nuc-deploy-* run found' >&2; exit 1; fi; \
	EDDA_ROLLBACK_MODE=simulate bash scripts/rollback_prod.sh '$(ENV_FILE)' "$$LATEST_RUN/rollback-manifest.env" "$$LATEST_RUN/pre-deploy.dump"

clean-runtime-artifacts:
	@if [ -e backend/app ]; then git checkout -- backend/app; fi
	@rm -f backend/switchyard
