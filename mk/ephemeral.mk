
# .PHONY: ephemeral-setup
# ephemeral-setup: ## Configure bonfire to run locally
# 	bonfire config write-default > $(PROJECT_DIR)/config/bonfire-config.yaml

APP ?= $(APP_NAME)
NAMESPACE ?= $(shell oc project -q 2>/dev/null)
# POOL could be:
#   default
#   minimal
#   managed-kafka
#   real-managed-kafka
POOL ?= default
export NAMESPACE
export APP
export POOL


# CLIENTS_RBAC_BASE_URL ?= http://localhost:8801/api/rbac/v1  # For local workstation
# CLIENTS_RBAC_BASE_URL ?= http://rbac-service:8080/api/rbac/v1
# export CLIENTS_RBAC_BASE_URL


ifneq (default,$(POOL))
EPHEMERAL_OPTS += --no-single-replicas
else
EPHEMERAL_OPTS += --single-replicas
endif

ifeq (False,$(CLIENTS_RBAC_ENABLED))
EPHEMERAL_OPTS += --set-parameter "$(APP_COMPONENT)/CLIENTS_RBAC_ENABLED=False"
else
ifneq (,$(CLIENTS_RBAC_BASE_URL))
EPHEMERAL_OPTS += --set-parameter "$(APP_COMPONENT)/CLIENTS_RBAC_BASE_URL=$(CLIENTS_RBAC_BASE_URL)"
endif
endif

# https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/onboarding.html
# Checkout this: https://github.com/RedHatInsights/bonfire/commit/15ac80bfcf9c386eabce33cb219b015a58b756c8
.PHONY: ephemeral-login
ephemeral-login: .old-ephemeral-login ## Help in login to the ephemeral cluster
	@#if [ "$(GH_SESSION_COOKIE)" != "" ]; then python3 $(GO_OUTPUT)/get-token.py; else $(MAKE) .old-ephemeral-login; fi

.PHONY: .old-ephemeral-login
.old-ephemeral-login:
	xdg-open "https://oauth-openshift.apps.c-rh-c-eph.8p0c.p1.openshiftapps.com/oauth/token/request"
	@echo "- Login with github"
	@echo "- Do click on 'Display Token'"
	@echo "- Copy 'Log in with this token' command"
	@echo "- Paste the command in your terminal"
	@echo ""
	@echo "Now you should have access to the cluster, remember to use bonfire to manage namespace lifecycle:"
	@echo '# make ephemeral-namespace-reserve'
	@echo ""
	@echo "Check the namespaces reserved to you by:"
	@echo '# make ephemeral-namespace-list'
	@echo ""
	@echo "If you need to extend 1hour the time for the namespace reservation"
	@echo '# make ephemeral-namespace-extend-1h'
	@echo ""
	@echo "Finally if you don't need the reserved namespace or just you want to cleanup and restart with a fresh namespace you run:"
	@echo '# make ephemeral-namespace-release-all'

# Download https://gitlab.cee.redhat.com/klape/get-token/-/blob/main/get-token.py
$(GO_OUTPUT/get-token.py):
	curl -Ls -o "$(GO_OUTPUT/get-token.py)" "https://gitlab.cee.redhat.com/klape/get-token/-/raw/main/get-token.py"

# Changes to config/bonfire-local.yaml could impact to this rule
.PHONY: ephemeral-deploy
ephemeral-deploy:  ## Deploy application using 'config/bonfire-local.yaml' file
	[ "$(EPHEMERAL_NO_BUILD)" == "y" ] || $(MAKE) ephemeral-build-deploy
	source .venv/bin/activate && \
	bonfire deploy \
	    --source appsre \
		--local-config-path configs/bonfire-local.yaml \
		--secrets-dir "$(PROJECT_DIR)/configs/secrets" \
		--import-secrets \
		--namespace "$(NAMESPACE)" \
		--set-parameter "$(APP_COMPONENT)/IMAGE=$(DOCKER_IMAGE_BASE)" \
		--set-parameter "$(APP_COMPONENT)/IMAGE_TAG=$(DOCKER_IMAGE_TAG)" \
		--frontends true \
		$(EPHEMERAL_OPTS) \
		$(APP) rbac

.PHONY: ephemeral-undeploy
ephemeral-undeploy: ## Undeploy application from the current namespace
	source .venv/bin/activate && \
	bonfire process \
	    --source appsre \
		--local-config-path configs/bonfire-local.yaml \
		--namespace "$(NAMESPACE)" \
		--set-parameter "$(APP_COMPONENT)/IMAGE=$(DOCKER_IMAGE_BASE)" \
		--set-parameter "$(APP_COMPONENT)/IMAGE_TAG=$(DOCKER_IMAGE_TAG)" \
		$(EPHEMERAL_OPTS) \
		$(APP) rbac 2>/dev/null | json2yaml | oc delete -f -
	! oc get secrets/content-sources-certs &>/dev/null || oc delete secrets/content-sources-certs

.PHONY: ephemeral-process
ephemeral-process: ## Process application from the current namespace
	source .venv/bin/activate && \
	bonfire process \
	    --source appsre \
		--local-config-path configs/bonfire-local.yaml \
		--namespace "$(NAMESPACE)" \
		--set-parameter "$(APP_COMPONENT)/IMAGE=$(DOCKER_IMAGE_BASE)" \
		--set-parameter "$(APP_COMPONENT)/IMAGE_TAG=$(DOCKER_IMAGE_TAG)" \
		$(EPHEMERAL_OPTS) \
		$(APP) 2>/dev/null | json2yaml

# TODO Add command to specify to bonfire the clowdenv template to be used
.PHONY: ephemeral-namespace-create
ephemeral-namespace-create:  ## Create a namespace (requires ephemeral environment)
	@command -v bonfire 1>/dev/null 2>/dev/null || { echo "error:bonfire is not available in this environment"; exit 1; }
	# bonfire namespace reserve --force --pool "$(POOL)" 2>/dev/null
	oc project "$(shell bonfire namespace reserve --force --pool "$(POOL)" 2>/dev/null)"

.PHONY: ephemeral-namespace-delete-all
ephemeral-namespace-delete-all: ## Delete all namespace created by us (requires ephemeral environment)
	source .venv/bin/activate && \
	for item in $$( bonfire namespace list --mine --output json | jq -r '. | to_entries | map(select(.key | match("ephemeral-*";"i"))) | map(.key) | .[]' ); do \
	  bonfire namespace release --force $$item ; \
	done

.PHONY: ephemeral-namespace-list
ephemeral-namespace-list: ## List all the namespaces reserved to the current user (requires ephemeral environment)
	source .venv/bin/activate && \
	bonfire namespace list --mine

.PHONY: ephemeral-namespace-extend-1h
ephemeral-namespace-extend-1h: ## Extend for 1 hour the usage of the current ephemeral environment
	source .venv/bin/activate && \
	bonfire namespace extend --duration 1h "$(NAMESPACE)"

.PHONY: ephemeral-namespace-describe
ephemeral-namespace-describe: ## Display information about the current namespace
	@source .venv/bin/activate && \
	bonfire namespace describe "$(NAMESPACE)"


# DOCKER_IMAGE_BASE should be a public image
# Tested by 'make ephemeral-build-deploy DOCKER_IMAGE_BASE=quay.io/avisied0/content-sources-backend'
.PHONY: ephemeral-build-deploy
ephemeral-build-deploy:  ## Build and deploy image using 'build_deploy.sh' scripts; It requires to pass DOCKER_IMAGE_BASE
	IMAGE="$(DOCKER_IMAGE_BASE)" IMAGE_TAG="$(DOCKER_IMAGE_TAG)" ./build_deploy.sh &> build_deploy.log

.PHONY: ephemeral-pr-checks
ephemeral-pr-checks:
	IMAGE="$(DOCKER_IMAGE_BASE)" bash ./pr_checks.sh

# https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/
.PHONY: ephemeral-run-dnsutil
ephemeral-run-dnsutil:  ## Run a shell in a new pod to debug dns situations
	oc run dnsutil --rm --image=registry.k8s.io/e2e-test-images/jessie-dnsutils:1.3 -it -- bash
