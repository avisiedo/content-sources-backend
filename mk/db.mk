##
# Set of rules to interact with a local database
# from a container and database initialization.
#
# Requires 'mk/docker.mk'
##



.PHONY: .db-health
.db-health:
	@$(DOCKER_COMPOSE) -f "$(DOCKER_COMPOSE_FILE)" exec postgres pg_isready &>/dev/null

.PHONY: .db-health-wait
.db-health-wait:
	@while ! $(MAKE) .db-health &>/dev/null; do printf "."; sleep 1; done

.PHONY: db-migrate-up
db-migrate-up: $(GO_OUTPUT)/dbmigrate ## Run dbmigrate up
	$(GO_OUTPUT)/dbmigrate up

.PHONY: db-migrate-seed
db-migrate-seed: $(GO_OUTPUT)/dbmigrate ## Run dbmigrate seed
	$(GO_OUTPUT)/dbmigrate seed

.PHONY: db-cli-connect
db-cli-connect: ## Open a postgres cli in the container (it requires infra-up)
	$(DOCKER_COMPOSE) -f "$(DOCKER_COMPOSE_FILE)" exec -it postgres psql "sslmode=disable dbname=$(DATABASE_NAME) user=$(DATABASE_USER) host=$(DATABASE_HOST) port=$(DATABASE_PORT) password=$(DATABASE_PASSWORD)"

.PHONY: db-dump-table
db-dump-table:
	$(DOCKER_COMPOSE) -f "$(DOCKER_COMPOSE_FILE)" exec -it postgres pg_dump --table "$(DATABASE_TABLE)" --schema-only --dbname=$(DATABASE_NAME) --host=$(DATABASE_HOST) --port=$(DATABASE_PORT) --username=$(DATABASE_USER)

.PHONY: db-shell
db-shell:
	$(DOCKER_COMPOSE) -f "$(DOCKER_COMPOSE_FILE)" exec -it postgres bash
