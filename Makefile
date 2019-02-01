# Copyright 2026 The MatrixHub Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GIT_TAG ?= $(shell git describe --tags --dirty --always)
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

PLATFORMS ?=
IMAGE_BUILD_CMD ?= docker buildx build

IMAGE_REPO := ghcr.io/matrixhub-ai/matrixhub
VERSION ?= $(GIT_TAG)

BASE_IMAGE_PREFIX ?= docker.io
NPM_CONFIG_REGISTRY ?= https://registry.npmjs.org
GOPROXY ?= https://proxy.golang.org,direct
HTTP_PROXY ?=
HTTPS_PROXY ?=
VITE_APP_API_URL ?= http://127.0.0.1:3001
LOCAL_CONFIG ?= config/config.yaml
MATRIXHUB_BASE_URL ?= http://localhost:3001
# URL by which MatrixHub server components reach MatrixHub itself. Needed by
# E2E tests that register the deployment as its own remote registry; unlike
# MATRIXHUB_BASE_URL it is resolved inside the deployment, not by the test
# process. Tests that need it skip when it is empty.
MATRIXHUB_SELF_URL ?=
# In-cluster Service address for the KIND deployment (see scripts/deploy-matrixhub.sh:
# release and namespace are both "matrixhub", service port is apiserver.service.port).
E2E_KIND_SELF_URL ?= http://matrixhub.matrixhub.svc.cluster.local:9527
UNIT_TEST_COVERAGE_PROFILE ?= coverage.out
GOLANGCI_LINT_VERSION ?= v2.12.2
ACTIONLINT_VERSION ?= latest
GOVULNCHECK_VERSION ?= latest

version_pkg = github.com/matrixhub-ai/matrixhub/pkg/version
LD_FLAGS += -X '$(version_pkg).GitVersion=$(GIT_TAG)'
LD_FLAGS += -X '$(version_pkg).GitCommit=$(GIT_COMMIT)'
LD_FLAGS += -X '$(version_pkg).BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'
PUSH ?= --load

.PHONY: all
all: help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9.-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: image-build
image-build: ## Build the MatrixHub image
	$(IMAGE_BUILD_CMD) \
		-t $(IMAGE_REPO):$(VERSION) \
		-f deploy/docker/matrixhub/Dockerfile \
		$(if $(PLATFORMS),--platform=$(PLATFORMS)) \
		--build-arg GOPROXY="$(GOPROXY)" \
		--build-arg HTTP_PROXY="$(HTTP_PROXY)" \
		--build-arg HTTPS_PROXY="$(HTTPS_PROXY)" \
		--build-arg BASE_IMAGE_PREFIX="$(BASE_IMAGE_PREFIX)" \
		$(PUSH) \
		$(IMAGE_BUILD_EXTRA_OPTS) \
		.

.PHONY: image-push
image-push: ## Build and push the MatrixHub image
	$(MAKE) image-build PUSH=--push

website/build: website
	make -C website build

.PHONY: serve-website
serve-website: ## Serve built documentation website locally
	make -C website serve

.PHONY: clean
clean: ## Clean all build artifacts
	rm -rf bin
	make -C website clean

.PHONY: local-run-ui
local-run-ui: ## Run frontend UI locally
	cd ui && pnpm i && VITE_APP_API_URL="$(VITE_APP_API_URL)" pnpm dev

.PHONY: local-build-ui
local-build-ui: ## Build UI for local API serving
	cd ui && pnpm i && pnpm build

.PHONY: local-run-api
local-run-api: ## Serve the API only
	go run ./cmd/matrixhub apiserver

.PHONY: local-run
local-run: ## Run MatrixHub locally (web + API)
	$(MAKE) -j2 local-run-api local-run-ui

.PHONY: local-run-built-ui
local-run-built-ui: local-build-ui ## Build UI and run API server serving it
	go run ./cmd/matrixhub apiserver -c $(LOCAL_CONFIG)

##@ Verification

.PHONY: verify
verify: verify.go verify.ui verify.workflow verify.govulncheck ## Run locally runnable CI static checks

.PHONY: verify.go
verify.go: lint verify.mocks ## Run Go lint and mock generation checks

.PHONY: lint
lint: ## Run golangci-lint
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with --fix option
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --fix

.PHONY: verify.mocks
verify.mocks: ## Verify generated gomock files
	bash ./scripts/verify-mocks.sh

.PHONY: verify.ui
verify.ui: ## Run UI lint, typecheck, and build
	cd ui && pnpm install --frozen-lockfile
	cd ui && pnpm run lint
	cd ui && pnpm run typecheck
	cd ui && pnpm run build

.PHONY: verify.workflow
verify.workflow: ## Run GitHub Actions workflow and release-note automation checks
	go run github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION) -color -shellcheck=
	bash ./scripts/verify-action-refs.sh
	bash -n scripts/changelog.sh scripts/changelog.test.sh
	bash scripts/changelog.test.sh

.PHONY: verify.govulncheck
verify.govulncheck: ## Run Go vulnerability reachability checks
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

##@ Testing

.PHONY: test.unit
test.unit: ## Run all unit tests
	bash ./scripts/unit-test.sh

.PHONY: test.unit.coverage
test.unit.coverage: ## Run all unit tests with coverage
	COVERAGE=true COVERAGE_PROFILE="$(UNIT_TEST_COVERAGE_PROFILE)" bash ./scripts/unit-test.sh

.PHONY: test.sqlite
test.sqlite: ## Run SQLite migration, repository, and session contract tests
	CGO_ENABLED=1 go test -count=1 ./internal/infra/config ./internal/infra/db ./internal/repo

.PHONY: deploy.kind-cluster
deploy.kind-cluster: ## Setup KIND cluster
	./scripts/setup-kind-cluster.sh

