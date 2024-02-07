
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd"
# Docker image name for the mkdocs based local development setup
IMAGE=metal-api/documentation

GOPRIVATE ?= "github.com/onmetal/*"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GITHUB_PAT_PATH ?=
ifeq (,$(GITHUB_PAT_PATH))
GITHUB_PAT_MOUNT ?=
else
GITHUB_PAT_MOUNT ?= --secret id=github_pat,src=$(GITHUB_PAT_PATH)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

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
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: vgopath openapi-gen models-schema applyconfiguration-gen deepcopy-gen client-gen ## Generate DeepCopy, DeepCopyInto, and DeepCopyObject method implementations and applyconfiguration.
	@VGOPATH=$(VGOPATH) \
   	MODELS_SCHEMA=$(MODELS_SCHEMA) \
   	DEEPCOPY_GEN=$(DEEPCOPY_GEN) \
   	CLIENT_GEN=$(CLIENT_GEN) \
   	OPENAPI_GEN=$(OPENAPI_GEN) \
   	APPLYCONFIGURATION_GEN=$(APPLYCONFIGURATION_GEN) \
	hack/generate.sh

.PHONY: fmt
fmt: goimports
	go fmt ./...
	$(GOIMPORTS) -w .

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint
	$(GOLANGCI_LINT) run --fix

.PHONY: add-license
add-license: addlicense ## Add license header to all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -f hack/license-header.txt {} +

.PHONY: check-license
check-license: addlicense ## Check license header presence in all .go files in project
	@find . -name '*.go' -exec $(ADDLICENSE) -check -c 'IronCore authors' {} +

.PHONY: check
check: manifests generate add-license lint test # Generate manifests, code, lint, add licenses, test

.PHONY: docs
docs: gen-crd-api-reference-docs ## Run go generate to generate API reference documentation.
	$(GEN_CRD_API_REFERENCE_DOCS) -api-dir ./apis/metal/v1alpha4 -config ./hack/api-reference/template.json -template-dir ./hack/api-reference/template -out-file ./docs/api-reference/metal.md

.PHONY: start-docs
start-docs: ## Start the local mkdocs based development environment.
	docker build -t $(IMAGE) -f docs/Dockerfile .
	docker run -p 8000:8000 -v `pwd`/:/docs $(IMAGE)

.PHONY: clean-docs
clean-docs: ## Remove all local mkdocs Docker images (cleanup).
	docker container prune --force --filter "label=project=metal_api_documentation"

.PHONY: test
test: envtest
	@KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" \
	KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=600s KUBEBUILDER_CONTROLPLANE_STOP_TIMEOUT=600s go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} $(GITHUB_PAT_MOUNT) .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

##@ Deployment

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/default | kubectl delete -f -

### AUXILIARY ###
LOCAL_BIN ?= $(shell pwd)/bin
$(LOCAL_BIN):
	mkdir -p $(LOCAL_BIN)

## Tools locations
ADDLICENSE ?= $(LOCAL_BIN)/addlicense
CONTROLLER_GEN ?= $(LOCAL_BIN)/controller-gen
GOLANGCI_LINT ?= $(LOCAL_BIN)/golangci-lint
GOIMPORTS ?= $(LOCAL_BIN)/goimports
ENVTEST ?= $(LOCAL_BIN)/setup-envtest
DEEPCOPY_GEN ?= $(LOCAL_BIN)/deepcopy-gen
CLIENT_GEN ?= $(LOCAL_BIN)/client-gen
LISTER_GEN ?= $(LOCAL_BIN)/lister-gen
INFORMER_GEN ?= $(LOCAL_BIN)/informer-gen
DEFAULTER_GEN ?= $(LOCAL_BIN)/defaulter-gen
CONVERSION_GEN ?= $(LOCAL_BIN)/conversion-gen
OPENAPI_GEN ?= $(LOCAL_BIN)/openapi-gen
APPLYCONFIGURATION_GEN ?= $(LOCAL_BIN)/applyconfiguration-gen
MODELS_SCHEMA ?= $(LOCAL_BIN)/models-schema
VGOPATH ?= $(LOCAL_BIN)/vgopath
GEN_CRD_API_REFERENCE_DOCS ?= $(LOCAL_BIN)/gen-crd-api-reference-docs
KUSTOMIZE ?= $(LOCAL_BIN)/kustomize

## Tools versions
ADDLICENSE_VERSION ?= v1.1.1
CONTROLLER_GEN_VERSION ?= v0.13.0
GOLANGCI_LINT_VERSION ?= v1.55.2
GOIMPORTS_VERSION ?= v0.16.1
ENVTEST_K8S_VERSION ?= 1.28.3
CODE_GENERATOR_VERSION ?= v0.28.3
VGOPATH_VERSION ?= v0.1.3
MODELS_SCHEMA_VERSION ?= main
GEN_CRD_API_REFERENCE_DOCS_VERSION ?= v0.3.0
KUSTOMIZE_VERSION ?= v4.5.4

.PHONY: addlicense
addlicense: $(ADDLICENSE)
$(ADDLICENSE): $(LOCAL_BIN)
	@test -s $(ADDLICENSE) || GOBIN=$(LOCAL_BIN) go install github.com/google/addlicense@$(ADDLICENSE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN)
$(CONTROLLER_GEN): $(LOCAL_BIN)
	@test -s $(CONTROLLER_GEN) && $(CONTROLLER_GEN) --version | grep -q $(CONTROLLER_GEN_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCAL_BIN)
	@test -s $(GOLANGCI_LINT) && $(GOLANGCI_LINT) --version | grep -q $(GOLANGCI_LINT_VERSION) || \
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: goimports
goimports: $(GOIMPORTS)
$(GOIMPORTS): $(LOCAL_BIN)
	@test -s $(GOIMPORTS) || GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCAL_BIN)
	@test -s $(ENVTEST) || GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: vgopath
vgopath: $(VGOPATH)
$(VGOPATH): $(LOCAL_BIN)
	@test -s $(VGOPATH) || GOBIN=$(LOCAL_BIN) go install github.com/ironcore-dev/vgopath@$(VGOPATH_VERSION)

.PHONY: deepcopy-gen
deepcopy-gen: $(DEEPCOPY_GEN)
$(DEEPCOPY_GEN): $(LOCAL_BIN)
	@test -s $(DEEPCOPY_GEN) || GOBIN=$(LOCAL_BIN) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GENERATOR_VERSION)

.PHONY: openapi-gen
openapi-gen: $(OPENAPI_GEN)
$(OPENAPI_GEN): $(LOCAL_BIN)
	@test -s $(OPENAPI_GEN) || GOBIN=$(LOCAL_BIN) go install k8s.io/code-generator/cmd/openapi-gen@$(CODE_GENERATOR_VERSION)

.PHONY: models-schema
models-schema: $(MODELS_SCHEMA)
$(MODELS_SCHEMA): $(LOCALBIN)
	@test -s $(MODELS_SCHEMA) || GOBIN=$(LOCAL_BIN) go install github.com/ironcore-dev/ironcore/models-schema@$(MODELS_SCHEMA_VERSION)

.PHONY: gen-crd-api-reference-docs
gen-crd-api-reference-docs: $(GEN_CRD_API_REFERENCE_DOCS) ## Download gen-crd-api-reference-docs locally if necessary.
$(GEN_CRD_API_REFERENCE_DOCS): $(LOCAL_BIN)
	@test -s $(GEN_CRD_API_REFERENCE_DOCS) || GOBIN=$(LOCAL_BIN) go install github.com/ahmetb/gen-crd-api-reference-docs@$(GEN_CRD_API_REFERENCE_DOCS_VERSION)

.PHONY: applyconfiguration-gen
applyconfiguration-gen: $(APPLYCONFIGURATION_GEN) ## Download applyconfiguration-gen locally if necessary.
$(APPLYCONFIGURATION_GEN): $(LOCALBIN)
	@test -s $(APPLYCONFIGURATION_GEN) || GOBIN=$(LOCAL_BIN) go install k8s.io/code-generator/cmd/applyconfiguration-gen@$(CODE_GENERATOR_VERSION)

.PHONY: client-gen
client-gen: $(CLIENT_GEN) ## Download client-gen locally if necessary.
$(CLIENT_GEN): $(LOCALBIN)
	@test -s $(CLIENT_GEN) || GOBIN=$(LOCAL_BIN) go install k8s.io/code-generator/cmd/client-gen@$(CODE_GENERATOR_VERSION)

.PHONY: kustomize
kustomize: $(KUSTOMIZE)
$(KUSTOMIZE): $(LOCAL_BIN)
	@test -s $(KUSTOMIZE) || GOBIN=$(LOCAL_BIN) go install sigs.k8s.io/kustomize/kustomize/v4@$(KUSTOMIZE_VERSION)
