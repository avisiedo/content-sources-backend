##
# Set of rules to interact with a local database
# from a container and database initialization.
#
# Requires 'mk/docker.mk'
##

.PHONY: db-up
db-up: DOCKER_IMAGE=docker.io/postgres:14
db-up: $(GO_OUTPUT)/dbmigrate ## Start postgres database
	$(DOCKER) volume exists postgres || $(DOCKER) volume create postgres
	$(DOCKER) container exists postgres || $(DOCKER) run \
	  -d \
	  --rm \
	  --name postgres \
	  -p 5432:5432 \
	  -e POSTGRES_PASSWORD=$(DATABASE_PASSWORD) \
	  -e POSTGRES_USER=$(DATABASE_USER) \
	  -e POSTGRES_DB=$(DATABASE_NAME) \
	  -v postgres:/var/lib/postgresql/data \
	  $(DOCKER_IMAGE)
	$(MAKE) db-migrate-up
	@echo "Now run db-migrate-seed to populate the database"

.PHONY: db-migrate-up
db-migrate-up: $(GO_OUTPUT)/dbmigrate ## Run dbmigrate up
	$(GO_OUTPUT)/dbmigrate up

.PHONY: db-migrate-seed
db-migrate-seed: ## Run dbmigrate seed
	$(GO_OUTPUT)/dbmigrate seed

.PHONY: db-down
db-down: ## Stop postgres database
	! $(DOCKER) container exists postgres || $(DOCKER) container stop postgres

.PHONY: db-clean
db-clean: db-down ## Clean database volume
	! $(DOCKER) volume exists postgres || $(DOCKER) volume rm postgres
