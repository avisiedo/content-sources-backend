##
# Set of rules to manage podman-compose
#
# Requires 'mk/db.mk'
# Requires 'mk/kafka.mk'
# Requires 'mk/pulp.mk'
##

.PHONY: deps-up
compose-up: $(GO_OUTPUT)/dbmigrate
	$(DATABASE_COMPOSE_OPTIONS) \
	$(KAFKA_COMPOSE_OPTIONS) \
	$(DOCKER)-compose --project-name=$(COMPOSE_PROJECT_NAME) -f compose/docker-compose.yml up --detach
	$(MAKE) .db-health-wait
	$(MAKE) db-migrate-up
	@echo "Run 'make db-migrate-seed' to seed the database"

.PHONY: deps-down
compose-down:
	$(DATABASE_COMPOSE_OPTIONS) \
	$(DOCKER)-compose --project-name=$(COMPOSE_PROJECT_NAME) -f compose/docker-compose.yml down --volumes
