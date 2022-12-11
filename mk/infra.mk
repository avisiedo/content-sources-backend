##
# Include this part to use docker-compose or podman-compose for your local
# infrastructure easily.
##

ifeq (docker,$(DOCKER))
DOCKER_COMPOSE ?= docker-compose
endif

ifeq (podman,$(DOCKER))
DOCKER_COMPOSE ?= podman-compose
endif

DOCKER_COMPOSE_FILE ?= deployments/docker-compose.yaml

export DATABASE_USER
export DATABASE_PASSWORD
export DATABASE_NAME

#
# Kafka variables
#

# Kafka version to download when building the container
KAFKA_VERSION ?= 3.3.1
# The local image to use
KAFKA_IMAGE ?= localhost/kafka:latest
# Options passed to the jvm invokation for zookeeper container
ZOOKEEPER_OPTS ?= -Dzookeeper.4lw.commands.whitelist=*
# Options passed to the jvm invokation for kafka container
KAFKA_OPTS ?= -Dzookeeper.4lw.commands.whitelist=*
# zookeepr client port; it is not publised but used inter containers
ZOOKEEPER_CLIENT_PORT ?= 2181
# The list of topics to be created; if more than one split them by a space
ifeq (,$(KAFKA_TOPICS))
	$(warning KAFKA_TOPICS is empty; probably missed definition at mk/variables.mk)
endif
KAFKA_TOPICS ?= platform.content-sources.introspect
KAFKA_GROUP_ID ?= content-sources

# The Kafka configuration directory that will be bound inside the containers
KAFKA_CONFIG_DIR ?= $(PROJECT_DIR)/kafka/config
# The Kafka data directory that will be bound inside the containers
# It must belong to the repository base directory
KAFKA_DATA_DIR ?= $(PROJECT_DIR)/kafka/data

# KAFKA_BOOTSTRAP_SERVERS ?= localhost:9092,localhost:9093
KAFKA_BOOTSTRAP_SERVERS ?= localhost:9092


export ZOOKEEPER_CLIENT_PORT
export ZOOKEEPER_OPTS
export KAFKA_DATA_DIR
export KAFKA_CONFIG_DIR
export KAFKA_TOPICS
export KAFKA_BOOTSTRAP_SERVERS



.PHONY: infra-up
infra-up: RUN_MIGRATE ?= $(MAKE) db-migrate-up
infra-up: ## Start local infrastructure
	@[ "$(DOCKER_COMPOSE)" != "" ] || (echo "error:DOCKER_COMPOSE is empty"; exit 1)
	$(DOCKER_COMPOSE) -f "$(DOCKER_COMPOSE_FILE)" up --detach
	@$(MAKE) .db-health-wait
	@$(RUN_MIGRATE)
	@echo "Run 'make db-migrate-seed' to seed the database"

.PHONY: infra-down
infra-down: ## Stop local infrastructure
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

.PHONY: infra-logs
infra-logs: ## Tail logs and follow them
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs --tail 10 --follow

.PHONY: infra-ps
infra-ps:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) ps

.PHONY: infra-build
infra-build:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build

.PHONY: infra-shell
infra-shell: SERVICE ?= dnsutil
infra-shell:
	docker-compose -f deployments/docker-compose.yaml exec -it $(SERVICE) bash









