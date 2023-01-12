##
# Set of rules to manage podman-compose
#
# Requires 'mk/db.mk'
# Requires 'mk/kafka.mk'
# Requires 'mk/pulp.mk'
##

COMPOSE_PROJECT_NAME ?= cs

.PHONY: deps-up
compose-up: $(GO_OUTPUT)/dbmigrate
	$(DATABASE_COMPOSE_OPTIONS) \
	$(KAFKA_COMPOSE_OPTIONS) \
	$(DOCKER)-compose --project-name=$(COMPOSE_PROJECT_NAME) -f deployments/docker-compose.yaml up --detach
	$(MAKE) .db-health-wait
	$(MAKE) db-migrate-up
	@echo "Run 'make db-migrate-seed' to seed the database"

.PHONY: compose-down
compose-down:
	$(DATABASE_COMPOSE_OPTIONS) \
	$(DOCKER)-compose --project-name=$(COMPOSE_PROJECT_NAME) -f deployments/docker-compose.yaml down --volumes

.PHONY: compose-clean
compose-clean: compose-down
	$(DOCKER) volume prune

.PHONY: compose-build
compose-build:
	$(DATABASE_COMPOSE_OPTIONS) \
	$(KAFKA_COMPOSE_OPTIONS) \
	$(DOCKER)-compose --project-name=$(COMPOSE_PROJECT_NAME) -f deployments/docker-compose.yaml build
