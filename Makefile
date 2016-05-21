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
