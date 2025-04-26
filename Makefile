MAIN_PACKAGE := ./cmd/portkey-mcp-server
BINARY_NAME ?= portkey-mcp-server

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.appVersion=$(shell git rev-parse --short HEAD)" -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BINARY_NAME)"

IMAGE_NAME ?= portkey-mcp-server
IMAGE_TAG ?= $(shell git rev-parse --short HEAD)
PLATFORMS ?= linux/amd64,linux/arm64

# Build a Docker image using Docker Buildx
# Usage:
#   make docker-build                                 # Build with default settings
#   make docker-build PLATFORMS=linux/amd64           # Build for specific platforms
#   make docker-build IMAGE_TAG=v1.0.0                # Set a specific tag
#   make docker-build GITHUB_PAT=/path/to/pat.txt     # Specify GitHub PAT file
#   make docker-build PUSH=true                       # Build and push to registry
.PHONY: docker-build
docker-build:
	@echo "Building Docker image $(IMAGE_NAME):$(IMAGE_TAG) for platforms: $(PLATFORMS)"
	@if [ -z "$(GITHUB_PAT)" ] && [ ! -f "$(HOME)/.github/token" ]; then \
		echo "Warning: GITHUB_PAT not provided and default location not found."; \
		echo "Private dependencies may fail to download."; \
		echo "Set GITHUB_PAT to path of file containing your GitHub token."; \
	fi
	
	@if [ -n "$(GITHUB_PAT)" ]; then \
		GITHUB_PAT_FILE="$(GITHUB_PAT)"; \
	elif [ -f "$(HOME)/.github/token" ]; then \
		GITHUB_PAT_FILE="$(HOME)/.github/token"; \
	fi
	
	@if [ -n "$(GITHUB_PAT_FILE)" ]; then \
		GITHUB_SECRET_ARGS="--secret id=GITHUB_PAT,src=$(GITHUB_PAT_FILE)"; \
	else \
		GITHUB_SECRET_ARGS=""; \
	fi
	
	@docker buildx create --use --name portkey-builder --driver docker-container --bootstrap || true
	@docker buildx build \
		--platform $(PLATFORMS) \
		$${GITHUB_SECRET_ARGS} \
		--build-arg APP_VERSION=$(IMAGE_TAG) \
		$(if $(PUSH),--push,--load) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		-t $(IMAGE_NAME):latest \
		.
	@echo "Docker image build complete: $(IMAGE_NAME):$(IMAGE_TAG)"

# Run the Docker container with appropriate port mapping and environment variables
# Usage:
#   make docker-run                                    # Run with default settings
#   make docker-run PORT=9000                          # Use a different host port
#   make docker-run IMAGE_TAG=v1.0.0                   # Run a specific image tag
#   make docker-run TRANSPORT=stdio                    # Use stdio transport instead of SSE
#   make docker-run PORTKEY_API_KEY=your-api-key       # Set Portkey API key
.PHONY: docker-run
docker-run:
	@echo "Running Docker container $(IMAGE_NAME):$(IMAGE_TAG)"
	@PORT=$${PORT:-8080}; \
	TRANSPORT=$${TRANSPORT:-sse}; \
	PORTKEY_API_KEY=$${PORTKEY_API_KEY:-dummy-key}; \
	CONTAINER_PORT=8080; \
	PORTKEY_BASE_URL=$${PORTKEY_BASE_URL:-https://api.portkey.ai/v1}; \
	PORT_MAPPING=""; \
	if [ "$${TRANSPORT}" = "sse" ]; then \
		PORT_MAPPING="-p $${PORT}:$${CONTAINER_PORT}"; \
	fi; \
	docker run --rm -it \
		$${PORT_MAPPING} \
		-e TRANSPORT=$${TRANSPORT} \
		-e TRANSPORT_SSE_ADDRESS=:$${CONTAINER_PORT} \
		-e PORTKEY_BASE_URL=$${PORTKEY_BASE_URL} \
		-e PORTKEY_API_KEY=$${PORTKEY_API_KEY} \
		$(IMAGE_NAME):$(IMAGE_TAG)

LINT_TIMEOUT ?= 5m
LINT_REPORT_FILE ?= lint-report.json
LINT_ERRORS_FILE ?= lint-errors.log

.PHONY: lint
lint:
	@echo "Running linter (timeout: $(LINT_TIMEOUT))..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5 # Fixed, so check for upgrades, occasionally
	@bash -c "golangci-lint run --out-format=json --timeout=$(LINT_TIMEOUT) > $(LINT_REPORT_FILE) 2> >(tee $(LINT_ERRORS_FILE) >&2)"
	@echo "Lint complete. Results in $(LINT_REPORT_FILE), errors in $(LINT_ERRORS_FILE)"

# golangci-lint maintains a cache to speed up subsequent runs. Sometimes this cache can
# retain stale results even after code has been fixed. Use this target to clear the cache
# if you're seeing persistent lint errors that you believe you've fixed.
.PHONY: lint-clear-cache
lint-clear-cache:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5 # Fixed, so check for upgrades, occasionally
	@echo "Clearing golangci-lint cache..."
	@golangci-lint cache clean
	@echo "Golangci-lint cache cleared successfully"

GOSEC_REPORT_FILE ?= gosec-report.json
GOSEC_ERRORS_FILE ?= gosec-errors.log

GOVULNCHECK_REPORT_FILE ?= govulncheck-report.json
GOVULNCHECK_ERRORS_FILE ?= govulncheck-errors.log

.PHONY: security
security:
	@echo "Running security checks..."
	@go install github.com/securego/gosec/v2/cmd/gosec@v2.21.4 # Fixed, so check for upgrades, occasionally
	@go install golang.org/x/vuln/cmd/govulncheck@v1.1.3 # Fixed, so check for upgrades, occasionally
	@echo "Running gosec..."
	@bash -c "gosec -confidence=high -exclude-dir=docs -fmt=json -out=$(GOSEC_REPORT_FILE) -severity=medium ./... 2> >(tee $(GOSEC_ERRORS_FILE) >&2)"
	@echo "Running govulncheck..."
	@bash -c "govulncheck -json ./... > $(GOVULNCHECK_REPORT_FILE) 2> >(tee $(GOVULNCHECK_ERRORS_FILE) >&2)"
	@echo "Security checks complete. Results in $(GOSEC_REPORT_FILE) and $(GOVULNCHECK_REPORT_FILE)"

GOVET_REPORT_FILE ?= vet-report.txt
GOVET_ERRORS_FILE ?= vet-errors.log

.PHONY: semantic-analysis
semantic-analysis:
	@echo "Running go vet..."
	@bash -c "go vet -all ./... > $(GOVET_REPORT_FILE) 2> >(tee $(GOVET_ERRORS_FILE) >&2)"
	@echo "Semantic analysis complete. Results in $(GOVET_REPORT_FILE)"

GOTEST_REPORT_FILE ?= test-report.json
GOTEST_ERRORS_FILE ?= test-errors.log

.PHONY: test
test:
	@$(MAKE) mocks
	@echo "Running tests..."
	@bash -c "go test -count=1 -covermode=atomic -json -race -run=Test -tags=unit,integration -v ./... > $(GOTEST_REPORT_FILE) 2> >(tee $(GOTEST_ERRORS_FILE) >&2)"
	@echo "Tests complete. Results in $(GOTEST_REPORT_FILE)"

BENCHMARK_REPORT_FILE ?= benchmark-report.txt
BENCHMARK_ERRORS_FILE ?= benchmark-errors.log

.PHONY: benchmark
benchmark:
	@$(MAKE) mocks
	@echo "Running benchmarks..."
	@bash -c "go test -bench=. -benchtime=1s -tags=benchmark -v ./... > $(BENCHMARK_REPORT_FILE) 2> >(tee $(BENCHMARK_ERRORS_FILE) >&2)"
	@echo "Benchmarks complete. Results in $(BENCHMARK_REPORT_FILE)"

.PHONY: mocks
mocks:
	@echo "Installing mockgen..."
	@go install go.uber.org/mock/mockgen@v0.5.0 # Fixed, so check for upgrades, occasionally
	@echo "Generating mocks..."
	@go generate ./... # generate mocks for tests
	@echo "Mock generation complete"

HOOKS_SOURCE_DIR ?= scripts/git-hooks
HOOKS_TARGET_DIR ?= .git/hooks
HOOKS = pre-push

.PHONY: install-hooks
install-hooks:
	@if [ ! -d ".git" ]; then \
		echo "Error: This is not a Git repository."; \
		exit 1; \
	fi
	@echo "Installing Git hooks..."
	@mkdir -p $(HOOKS_TARGET_DIR)
	@for hook in $(HOOKS); do \
		cp $(HOOKS_SOURCE_DIR)/$$hook $(HOOKS_TARGET_DIR)/$$hook; \
		chmod +x $(HOOKS_TARGET_DIR)/$$hook; \
		echo "Installed $$hook"; \
	done
	@echo "All hooks installed successfully."

.PHONY: clean-hooks
clean-hooks:
	@echo "Cleaning up Git hooks..."
	@for hook in $(HOOKS); do \
		rm -f $(HOOKS_TARGET_DIR)/$$hook; \
		echo "Removed $$hook"; \
	done
	@echo "Git hooks cleaned."
