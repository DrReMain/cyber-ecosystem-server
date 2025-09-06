# ============================================================================
# Cyber-Ecosystem Microservice Build System
# Supports API Gateway and RPC Services with Ent ORM Integration
# ============================================================================

# ================================ Version Information =======================
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "v0.0.0-dev")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')
BUILD_USER := $(shell whoami)
BUILD_HOST := $(shell hostname)
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# ================================ System Information ========================
LOCAL_GOHOSTOS := $(shell go env GOHOSTOS)
LOCAL_GOARCH := $(shell go env GOARCH)
LOCAL_GOPATH := $(shell go env GOPATH)

# ================================ Feature Configuration =====================
ENT_FEATURE := sql/execquery,sql/modifier,intercept
MYSQL_VERSION := 8.4-oracle
DEFAULT_LANG := en-US

# ================================ Color Definitions =========================
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
WHITE := \033[37m
RESET := \033[0m

# ================================ Helper Functions ==========================
define print_header
	echo "$(CYAN)============================================================$(RESET)"
	echo "$(CYAN) $(1)$(RESET)"
	echo "$(CYAN)============================================================$(RESET)"
endef

define print_success
	echo "$(GREEN)✅ $(1)$(RESET)"
endef

define print_info
	echo "$(BLUE)ℹ️  $(1)$(RESET)"
endef

define print_warning
	echo "$(YELLOW)⚠️  $(1)$(RESET)"
endef

define print_error
	echo "$(RED)❌ $(1)$(RESET)"
endef

# ================================ Parameter Validation ======================
define check_target
	if [ -z "$${target}" ]; then \
		$(call print_error,Parameter 'target' is required); \
		echo "$${YELLOW}Example: make ent_new target=system entity=User$${RESET}"; \
		exit 1; \
	fi
endef

define check_entity
	if [ -z "$${entity}" ]; then \
		$(call print_error,Parameter 'entity' is required); \
		echo "$${YELLOW}Example: make ent_new target=system entity=User$${RESET}"; \
		exit 1; \
	fi
endef

# ================================ Default Target ============================
.DEFAULT_GOAL := help
.PHONY: init tidy clean ent_new ent_gen ent_visualize ent_migrate_diff ent_migrate_apply \
        rpc_pb rpc_gen api_locales api_format api_gen \
        rpc_build api_build start stop restart \
        image run help version lint test

# ================================ Development Environment ===================

# Initialize development environment
init:
	@$(call print_header,Initializing Development Environment)
	@$(call print_info,Checking Go environment...)
	@go version || ($(call print_error,Go is not installed or not in PATH) && exit 1)
	@$(call print_info,Checking Docker environment...)
	@docker --version >/dev/null 2>&1 || $(call print_warning,Docker not installed - some features will be unavailable)
	@$(call print_info,Installing goctl toolchain...)
	@go install github.com/zeromicro/go-zero/tools/goctl@latest
	@$(call print_info,Checking and installing goctl dependencies...)
	@goctl env check --install --verbose --force
	@$(call print_info,Installing development tools...)
	@go install entgo.io/ent/cmd/ent@latest
	@$(call print_success,Development environment initialized successfully)

# Update project dependencies
tidy:
	@$(call print_header,Updating Project Dependencies)
	@$(call print_info,Cleaning module cache...)
	@go clean -modcache
	@$(call print_info,Downloading and organizing dependencies...)
	@go mod tidy -v
	@$(call print_info,Verifying dependency integrity...)
	@go mod verify
	@$(call print_success,Dependencies updated successfully)

# Clean build artifacts
clean:
	@$(call print_header,Cleaning Build Artifacts)
	@$(call print_info,Cleaning Go build cache...)
	@go clean -cache -testcache -modcache
	@$(call print_success,Cleanup completed)

# Run code quality checks
lint:
	@$(call print_header,Running Code Quality Checks)
	@$(call print_info,Running go vet...)
	@go vet ./...
	@$(call print_info,Checking go fmt...)
	@if [ -n "$$(gofmt -l .)" ]; then \
		$(call print_error,Code is not properly formatted. Run 'go fmt ./...'); \
		gofmt -l .; \
		exit 1; \
	fi
	@$(call print_success,Code quality checks passed)

# Run tests
test:
	@$(call print_header,Running Tests)
	@$(call print_info,Running unit tests...)
	@go test -race -cover -v ./...
	@$(call print_success,All tests passed)

# ================================ Database ORM (Ent) ========================

# Create new Ent entity
# Usage: make ent_new target=system entity=User
ent_new:
	@$(call check_target)
	@$(call check_entity)
	@$(call print_header,Creating Ent Entity)
	@$(call print_info,Target service: $(target))
	@$(call print_info,Entity name: $(entity))
	@if [ ! -d "rpc/$(target)" ]; then \
		$(call print_error,RPC service directory 'rpc/$(target)' does not exist); \
		$(call print_info,Please create the RPC service structure first); \
		exit 1; \
	fi
	@mkdir -p rpc/$(target)/ent/schema
	@go run -mod=mod entgo.io/ent/cmd/ent new --target=rpc/$(target)/ent/schema $(entity)
	@$(call print_success,Entity '$(entity)' created successfully)
	@$(call print_info,Schema file: rpc/$(target)/ent/schema/$(shell echo $(entity) | tr A-Z a-z).go)
	@$(call print_warning,Don not forget to define your schema fields and relationships!)

# Generate Ent code
# Usage: make ent_gen target=system
ent_gen:
	@$(call check_target)
	@$(call print_header,Generating Ent Code)
	@$(call print_info,Target service: $(target))
	@if [ ! -d "rpc/$(target)/ent/schema" ]; then \
		$(call print_error,Schema directory 'rpc/$(target)/ent/schema' does not exist); \
		$(call print_info,Please create entities first using: make ent_new target=$(target) entity=YourEntity); \
		exit 1; \
	fi
	@if [ ! -n "$$(find rpc/$(target)/ent/schema -name '*.go' -not -path '*/.*' 2>/dev/null)" ]; then \
		$(call print_error,No entity schemas found in rpc/$(target)/ent/schema/); \
		$(call print_info,Please create entities first using: make ent_new target=$(target) entity=YourEntity); \
		exit 1; \
	fi
	@mkdir -p rpc/$(target)/ent/template
	@go run -mod=mod entgo.io/ent/cmd/ent generate \
		--template glob="./rpc/$(target)/ent/template/*.tmpl" \
		./rpc/$(target)/ent/schema \
		--feature $(ENT_FEATURE)
	@$(call print_success,Ent code generation completed)
	@$(call print_info,Generated files are in: rpc/$(target)/ent/)

# Visualize database schema
# Usage: make ent_visualize target=system
ent_visualize:
	@$(call check_target)
	@$(call print_header,Generating Database Schema Visualization)
	@$(call print_info,Target service: $(target))
	@$(call print_warning,This requires Docker to be running)
	@command -v atlas >/dev/null 2>&1 || ($(call print_error,Atlas CLI is not installed. Please install it from https://atlasgo.io/getting-started\#installation) && exit 1)
	@atlas schema inspect \
		-u "ent://rpc/$(target)/ent/schema" \
		--dev-url "docker+mysql://_/mysql:$(MYSQL_VERSION)/dev" \
		-w
	@$(call print_success,Database schema visualization opened in browser)

# Generate database migration files
# Usage: make ent_migrate_diff target=system
ent_migrate_diff:
	@$(call check_target)
	@$(call print_header,Generating Database Migration)
	@$(call print_info,Target service: $(target))
	@command -v atlas >/dev/null 2>&1 || ($(call print_error,Atlas CLI is not installed. Please install it from https://atlasgo.io/getting-started\#installation) && exit 1)
	@mkdir -p rpc/$(target)/ent/migrate/migrations
	@atlas migrate diff \
		--dir "file://rpc/$(target)/ent/migrate/migrations" \
		--to "ent://rpc/$(target)/ent/schema" \
		--dev-url "docker+mysql://_/mysql:$(MYSQL_VERSION)/dev"
	@$(call print_success,Migration files generated successfully)
	@$(call print_info,Migration directory: rpc/$(target)/ent/migrate/migrations/)
	@$(call print_warning,Please review the generated migration files before applying)

# Apply database migrations
# Usage: make ent_migrate_apply target=system dsn=root:pass@localhost:3306/dbname
ent_migrate_apply:
	@$(call check_target)
	@if [ -z "$(dsn)" ]; then \
		$(call print_error,Parameter 'dsn' is required); \
		echo "$(YELLOW)Format: username:password@host:port/database$(RESET)"; \
		echo "$(YELLOW)Example: make ent_migrate_apply target=system dsn=root:pass@localhost:3306/mydb$(RESET)"; \
		exit 1; \
	fi
	@$(call print_header,Applying Database Migration)
	@$(call print_info,Target service: $(target))
	@$(call print_info,Database DSN: mysql://$(dsn))
	@if [ ! -d "rpc/$(target)/ent/migrate/migrations" ]; then \
		$(call print_error,No migration directory found. Generate migrations first using: make ent_migrate_diff target=$(target)); \
		exit 1; \
	fi
	@$(call print_warning,This operation will modify your database schema!)
	@$(call print_warning,Make sure you have a backup of your database)
	@printf "$(YELLOW)Are you sure you want to continue? [y/N]: $(RESET)"; \
	read -r confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "Operation cancelled"; \
		exit 1; \
	fi
	@atlas migrate apply \
		--dir "file://rpc/$(target)/ent/migrate/migrations" \
		--url "mysql://$(dsn)"
	@$(call print_success,Database migration applied successfully)

# ================================ RPC Services ==============================

# Combine Proto files
# Usage: make rpc_pb target=system
rpc_pb:
	@$(call check_target)
	@$(call print_header,Combining Proto Files)
	@$(call print_info,Target service: $(target))
	@if [ ! -f "rpc/$(target)/desc/main.go" ]; then \
		$(call print_error,Proto combination tool not found: rpc/$(target)/desc/main.go); \
		$(call print_info,Please ensure the proto description tool exists); \
		exit 1; \
	fi
	@go run ./rpc/$(target)/desc/main.go
	@$(call print_success,Proto files combined successfully)
	@$(call print_info,Output file: rpc/$(target)/$(target).proto)

# Generate RPC service code
# Usage: make rpc_gen target=system
rpc_gen:
	@$(call check_target)
	@$(call print_header,Generating RPC Service Code)
	@$(call print_info,Target service: $(target))
	@if [ ! -f "rpc/$(target)/$(target).proto" ]; then \
		$(call print_error,Proto file not found. Please run: make rpc_pb target=$(target)); \
		exit 1; \
	fi
	@goctl rpc protoc ./rpc/$(target)/$(target).proto \
		--go_out=./rpc/$(target)/ \
		--go-grpc_out=./rpc/$(target)/ \
		--zrpc_out=./rpc/$(target)/ \
		-m --style=go_zero
	@$(call print_success,RPC service code generated successfully)
	@$(call print_info,Generated files are in: rpc/$(target)/)

# ================================ API Services ==============================

# Combine localization files
# Usage: make api_locales target=admin default=zh-CN
api_locales:
	@$(call check_target)
	@$(call print_header,Combining Localization Files)
	@$(call print_info,Target service: $(target))
	@$(call print_info,Default language: $(if $(default),$(default),$(DEFAULT_LANG)))
	@if [ ! -f "api/$(target)/lang/main.go" ]; then \
		$(call print_error,Localization tool not found: api/$(target)/lang/main.go); \
		$(call print_info,Please ensure the localization tool exists); \
		exit 1; \
	fi
	@if [ ! -d "api/$(target)/lang" ]; then \
		$(call print_error,Language directory not found: api/$(target)/lang); \
		exit 1; \
	fi
	@go run ./api/$(target)/lang/main.go \
		-messages api/$(target)/lang \
		-default $(if $(default),$(default),$(DEFAULT_LANG)) \
		--strict
	@$(call print_success,Localization files combined successfully)

# Format API files
# Usage: make api_format target=admin
api_format:
	@$(call check_target)
	@$(call print_header,Formatting API Files)
	@$(call print_info,Target service: $(target))
	@if [ ! -d "api/$(target)/desc" ]; then \
		$(call print_error,API description directory not found: api/$(target)/desc); \
		exit 1; \
	fi
	@goctl api format --dir api/$(target)/desc/
	@$(call print_success,API files formatted successfully)

# Generate API service code
# Usage: make api_gen target=admin
api_gen:
	@$(call check_target)
	@$(call print_header,Generating API Service Code)
	@$(call print_info,Target service: $(target))
	@if [ ! -f "api/$(target)/desc/$(target).api" ]; then \
		$(call print_error,API description file not found: api/$(target)/desc/$(target).api); \
		$(call print_info,Please ensure the API description file exists); \
		exit 1; \
	fi
	@goctl api go \
		-api ./api/$(target)/desc/$(target).api \
		-dir ./api/$(target)/ \
		--style=go_zero
	@$(call print_success,API service code generated successfully)
	@$(call print_info,Generated files are in: api/$(target)/)

# ================================ Build & Deploy ============================

# Build RPC service
# Usage: make rpc_build target=system os=linux arch=amd64
rpc_build:
	@$(call check_target)
	@$(call print_header,Building RPC Service)
	@if [ ! -f "rpc/$(target)/$(target).go" ]; then \
		$(call print_error,RPC service main file not found: rpc/$(target)/$(target).go); \
		$(call print_info,Please ensure the service code has been generated); \
		exit 1; \
	fi
	@$(eval OS := $(if $(os),$(os),$(LOCAL_GOHOSTOS)))
	@$(eval ARCH := $(if $(arch),$(arch),$(LOCAL_GOARCH)))
	@$(eval EXT := $(if $(findstring windows,$(OS)),.exe,))
	@$(call print_info,Target platform: $(OS)/$(ARCH))
	@$(call print_info,Version: $(VERSION))
	@$(call print_info,Build time: $(BUILD_TIME))
	@$(call print_info,Commit hash: $(COMMIT_HASH))
	@mkdir -p target/rpc_$(target)
	@env CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-s -w \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildTime=$(BUILD_TIME)' \
		-X 'main.BuildUser=$(BUILD_USER)' \
		-X 'main.BuildHost=$(BUILD_HOST)' \
		-X 'main.CommitHash=$(COMMIT_HASH)'" \
		-trimpath \
		-o target/rpc_$(target)/rpc_$(target)$(EXT) \
		-v ./rpc/$(target)/$(target).go
	@if [ -f "rpc/$(target)/etc/$(target).yaml" ]; then \
		cp rpc/$(target)/etc/$(target).yaml target/rpc_$(target)/rpc_$(target).yaml; \
	else \
		$(call print_warning,Configuration file not found - service might not start properly); \
	fi
	@$(call print_success,RPC service built successfully)
	@$(call print_info,Binary: target/rpc_$(target)/rpc_$(target)$(EXT))
	@if command -v du >/dev/null 2>&1; then \
		$(call print_info,Size: $$(du -h target/rpc_$(target)/rpc_$(target)$(EXT) | cut -f1)); \
	fi

# Build API service
# Usage: make api_build target=admin os=linux arch=amd64
api_build:
	@$(call check_target)
	@$(call print_header,Building API Service)
	@if [ ! -f "api/$(target)/$(target).go" ]; then \
		$(call print_error,API service main file not found: api/$(target)/$(target).go); \
		$(call print_info,Please ensure the service code has been generated); \
		exit 1; \
	fi
	@$(eval OS := $(if $(os),$(os),$(LOCAL_GOHOSTOS)))
	@$(eval ARCH := $(if $(arch),$(arch),$(LOCAL_GOARCH)))
	@$(eval EXT := $(if $(findstring windows,$(OS)),.exe,))
	@$(call print_info,Target platform: $(OS)/$(ARCH))
	@$(call print_info,Version: $(VERSION))
	@$(call print_info,Build time: $(BUILD_TIME))
	@$(call print_info,Commit hash: $(COMMIT_HASH))
	@mkdir -p target/api_$(target)
	@env CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-s -w \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildTime=$(BUILD_TIME)' \
		-X 'main.BuildUser=$(BUILD_USER)' \
		-X 'main.BuildHost=$(BUILD_HOST)' \
		-X 'main.CommitHash=$(COMMIT_HASH)'" \
		-trimpath \
		-o target/api_$(target)/api_$(target)$(EXT) \
		-v ./api/$(target)/$(target).go
	@if [ -f "api/$(target)/etc/$(target).yaml" ]; then \
		cp api/$(target)/etc/$(target).yaml target/api_$(target)/api_$(target).yaml; \
	else \
		$(call print_warning,Configuration file not found - service might not start properly); \
	fi
	@$(call print_success,API service built successfully)
	@$(call print_info,Binary: target/api_$(target)/api_$(target)$(EXT))
	@if command -v du >/dev/null 2>&1; then \
		$(call print_info,Size: $$(du -h target/api_$(target)/api_$(target)$(EXT) | cut -f1)); \
	fi

# ================================ Service Management ========================

# Start service
# Usage: make start target=rpc_system OR make start target=api_admin
start:
	@$(call check_target)
	@$(call print_header,Starting Service)
	@SERVICE_BINARY="target/$(target)/$(target)"; \
	if [ -f "$${SERVICE_BINARY}.exe" ]; then \
		SERVICE_BINARY="$${SERVICE_BINARY}.exe"; \
	fi; \
	if [ ! -f "$$SERVICE_BINARY" ]; then \
		$(call print_error,Service binary not found. Please build the service first); \
		$(call print_info,Run: make rpc_build target=<service> OR make api_build target=<service>); \
		exit 1; \
	fi; \
	CONFIG_FILE="target/$(target)/$(target).yaml"; \
	if [ ! -f "$$CONFIG_FILE" ]; then \
		$(call print_warning,Configuration file not found. Service might not start properly); \
	fi; \
	$(call print_info,Starting service: $(target)); \
	nohup "$$SERVICE_BINARY" -f "$$CONFIG_FILE" > target/$(target)/service.log 2>&1 & echo $$! > target/$(target)/service.pid; \
	sleep 2; \
	SERVICE_PID=$$(pgrep -f "$(target)" | head -n1); \
	if [ -n "$$SERVICE_PID" ]; then \
		$(call print_success,Service started successfully - PID: $$SERVICE_PID); \
		$(call print_info,Log file: target/$(target)/service.log); \
		$(call print_info,Use 'tail -f target/$(target)/service.log' to view logs); \
	else \
		$(call print_error,Service failed to start. Check the log file for details); \
		$(call print_info,Log content:); \
		if [ -f "target/$(target)/service.log" ]; then tail -10 "target/$(target)/service.log"; fi; \
		exit 1; \
	fi

# Stop service
# Usage: make stop target=rpc_system
stop:
	@$(call check_target)
	@$(call print_header,Stopping Service)
	@PID_FILE="target/$(target)/service.pid"; \
	if [ ! -f "$$PID_FILE" ]; then \
		$(call print_warning,PID file not found. Checking if service is running by name); \
		SERVICE_PID=$$(pgrep -f "^target/$(target)/$(target)$$"); \
	else \
		SERVICE_PID=$$(cat "$$PID_FILE"); \
	fi; \
	if [ -z "$$SERVICE_PID" ]; then \
		$(call print_success,Service '$(target)' is not running); \
		exit 0; \
	fi; \
	$(call print_info,Stopping service: $(target) - PID: $$SERVICE_PID); \
	kill -TERM "$$SERVICE_PID" 2>/dev/null || true; \
	for i in 5 4 3 2 1; do \
		if ! kill -0 "$$SERVICE_PID" 2>/dev/null; then \
			break; \
		fi; \
		printf "$(YELLOW)Waiting for graceful shutdown... $$i$(RESET)\n"; \
		sleep 1; \
	done; \
	if kill -0 "$$SERVICE_PID" 2>/dev/null; then \
		$(call print_warning,Service did not stop gracefully. Forcing termination...); \
		kill -KILL "$$SERVICE_PID" 2>/dev/null || true; \
		sleep 1; \
	fi; \
	if kill -0 "$$SERVICE_PID" 2>/dev/null; then \
		$(call print_error,Failed to stop service); \
		exit 1; \
	else \
		rm -f "$$PID_FILE"; \
		$(call print_success,Service stopped successfully); \
	fi

# Restart service
# Usage: make restart target=rpc_system
restart: stop start

# ================================ Containerization ==========================

# Build Docker image
# Usage: make image target=rpc_system
image:
	@$(call check_target)
	@$(call print_header,Building Docker Image)
	@if [ ! -f "target/$(target)/Dockerfile" ]; then \
		$(call print_error,Dockerfile not found: target/$(target)/Dockerfile); \
		$(call print_info,Please create a Dockerfile for your service); \
		exit 1; \
	fi
	@$(call print_info,Image name: $(target):$(VERSION))
	@$(call print_info,Also tagged as: $(target):latest)
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg BUILD_HOST=$(BUILD_HOST) \
		--build-arg COMMIT_HASH=$(COMMIT_HASH) \
		-t $(target):$(VERSION) \
		-t $(target):latest \
		-f target/$(target)/Dockerfile .
	@$(call print_success,Docker image built successfully)
	@$(call print_info,Image tags: $(target):$(VERSION) and $(target):latest)
	@IMAGE_SIZE=$$(docker images $(target):$(VERSION) --format 'table {{.Size}}' | tail -1); \
	$(call print_info,Image size: $$IMAGE_SIZE)

# Run Docker container
# Usage: make run target=rpc_system p=8001
run:
	@$(call check_target)
	@if [ -z "$(p)" ]; then \
		$(call print_error,Parameter 'p' (port) is required); \
		echo "$(YELLOW)Example: make run target=rpc_system p=8001$(RESET)"; \
		exit 1; \
	fi
	@$(call print_header,Starting Docker Container)
	@if ! docker images $(target) --format '{{.Repository}}' | grep -q "^$(target)$$"; then \
		$(call print_error,Docker image '$(target)' not found. Build it first using: make image target=$(target)); \
		exit 1; \
	fi
	@if docker ps -a --format 'table {{.Names}}' | grep -q "^$(target)$$"; then \
		$(call print_warning,Container '$(target)' already exists. Removing...); \
		docker rm -f $(target) 2>/dev/null || true; \
	fi
	@$(call print_info,Container name: $(target))
	@$(call print_info,Port mapping: $(p):$(p))
	@docker run -itd \
		-p $(p):$(p) \
		--name=$(target) \
		--restart=unless-stopped \
		$(target):$(VERSION)
	@$(call print_success,Container started successfully)
	@$(call print_info,Access URL: http://localhost:$(p))
	@$(call print_info,View logs: docker logs -f $(target))
	@$(call print_info,Container status: docker ps | grep $(target))

# ================================ Utility Commands ==========================

# Display version information
version:
	@$(call print_header,Version Information)
	@echo "$(CYAN)Project Version:$(RESET)  $(VERSION)"
	@echo "$(CYAN)Build Time:$(RESET)       $(BUILD_TIME)"
	@echo "$(CYAN)Build User:$(RESET)       $(BUILD_USER)"
	@echo "$(CYAN)Build Host:$(RESET)       $(BUILD_HOST)"
	@echo "$(CYAN)Commit Hash:$(RESET)      $(COMMIT_HASH)"
	@echo "$(CYAN)System Info:$(RESET)      $(LOCAL_GOHOSTOS)/$(LOCAL_GOARCH)"
	@echo "$(CYAN)Go Path:$(RESET)          $(LOCAL_GOPATH)"
	@echo "$(CYAN)Go Version:$(RESET)       $(shell go version 2>/dev/null || echo 'Not installed')"
	@echo "$(CYAN)Goctl Version:$(RESET)    $(shell goctl -v 2>/dev/null || echo 'Not installed')"
	@echo "$(CYAN)Atlas Version:$(RESET)    $(shell atlas version 2>/dev/null || echo 'Not installed')"

# Display comprehensive help information
help:
	$(call print_header,Cyber-Ecosystem Microservice Build System Help)
	@echo ""
	@echo "This Makefile provides a complete build system for Cyber-Ecosystem microservices"
	@echo "with Ent ORM integration, supporting both API gateways and RPC services."
	@echo ""
	@echo "$(CYAN)🚀 QUICK START:$(RESET)"
	@echo "  1. $(GREEN)make init$(RESET)                                # Initialize development environment"
	@echo "  2. $(GREEN)make ent_new target=system entity=User$(RESET)   # Create database entity"
	@echo "  3. $(GREEN)make ent_gen target=system$(RESET)               # Generate ORM code"
	@echo "  4. $(GREEN)make rpc_gen target=system$(RESET)               # Generate RPC service"
	@echo "  5. $(GREEN)make rpc_build target=system$(RESET)             # Build the service"
	@echo "  6. $(GREEN)make start target=rpc_system$(RESET)             # Start the service"
	@echo ""
	@echo "$(CYAN)📁 PROJECT STRUCTURE:$(RESET)"
	@echo "  ├── api/                     # API Gateway services"
	@echo "  │   └── [service]/           # Individual API service"
	@echo "  │       ├── desc/            # API description files (.api)"
	@echo "  │       ├── etc/             # Configuration files"
	@echo "  │       └── lang/            # Localization files"
	@echo "  ├── rpc/                     # RPC services"
	@echo "  │   └── [service]/           # Individual RPC service"
	@echo "  │       ├── desc/            # Proto description files"
	@echo "  │       ├── ent/             # Ent ORM schemas and generated code"
	@echo "  │       └── etc/             # Configuration files"
	@echo "  └── target/                  # Build output directory"
	@echo ""
	@echo "$(CYAN)🛠️  DEVELOPMENT ENVIRONMENT:$(RESET)"
	@echo "  $(GREEN)init$(RESET)                         Initialize development environment Install Go tools: goctl, ent, Check Docker availability"
	@echo "  $(GREEN)tidy$(RESET)                         Update and organize Go module dependencies, Clean module cache and verify integrity"
	@echo "  $(GREEN)clean$(RESET)                        Clean all caches, Remove Go build cache"
	@echo "  $(GREEN)lint$(RESET)                         Run code quality checks, Execute go vet, and go fmt checks"
	@echo "  $(GREEN)test$(RESET)                         Run all unit tests with race detection, Include coverage reports"
	@echo ""
	@echo "$(CYAN)🗄️  DATABASE ORM (ENT):$(RESET)"
	@echo "  $(GREEN)ent_new$(RESET)                      Create new database entity schema"
	@echo "                               Usage: make ent_new target=system entity=User"
	@echo "                               Creates: rpc/system/ent/schema/user.go"
	@echo ""
	@echo "  $(GREEN)ent_gen$(RESET)                      Generate Ent ORM code from schemas"
	@echo "                               Usage: make ent_gen target=system"
	@echo "                               Generates client, mutations, queries, and types"
	@echo ""
	@echo "  $(GREEN)ent_visualize$(RESET)                Generate database schema visualization"
	@echo "                               Requires Docker and Atlas CLI"
	@echo "                               Opens browser with ER diagram"
	@echo ""
	@echo "  $(GREEN)ent_migrate_diff$(RESET)             Generate database migration files"
	@echo "                               Compare current schema with database"
	@echo "                               Create versioned migration files"
	@echo ""
	@echo "  $(GREEN)ent_migrate_apply$(RESET)            Apply database migrations"
	@echo "                               Usage: make ent_migrate_apply target=system dsn=user:pass@host:port/db"
	@echo "                               ⚠️  This modifies your database!"
	@echo ""
	@echo "$(CYAN)🌐 RPC SERVICES:$(RESET)"
	@echo "  $(GREEN)rpc_pb$(RESET)                       Combine multiple .proto files"
	@echo "                               Usage: make rpc_pb target=system"
	@echo "                               Merges desc/*.proto into single system.proto"
	@echo ""
	@echo "  $(GREEN)rpc_gen$(RESET)                      Generate RPC service code from proto"
	@echo "                               Creates gRPC client/server code"
	@echo "                               Generates Go-Zero RPC service structure"
	@echo ""
	@echo "$(CYAN)🔌 API SERVICES:$(RESET)"
	@echo "  $(GREEN)api_locales$(RESET)                  Combine localization files"
	@echo "                               Usage: make api_locales target=admin default=en-US"
	@echo "                               Merge multiple language files into one"
	@echo ""
	@echo "  $(GREEN)api_format$(RESET)                   Format .api description files"
	@echo "                               Standardize API definition formatting"
	@echo ""
	@echo "  $(GREEN)api_gen$(RESET)                      Generate API service code"
	@echo "                               Usage: make api_gen target=admin"
	@echo "                               Creates HTTP handlers and routing"
	@echo ""
	@echo "$(CYAN)🔨 BUILD & COMPILE:$(RESET)"
	@echo "  $(GREEN)rpc_build$(RESET)                    Build RPC service binary"
	@echo "                               Usage: make rpc_build target=system [os=linux|darwin|windows] [arch=amd64|arm64]"
	@echo "                               Cross-compilation supported, If it is used for Docker, the operating system should be Linux."
	@echo "                               Embeds version info and build metadata"
	@echo ""
	@echo "  $(GREEN)api_build$(RESET)                    Build API service binary"
	@echo "                               Usage: make api_build target=admin [os=linux|darwin|windows] [arch=amd64|arm64]"
	@echo "                               Cross-compilation supported, If it is used for Docker, the operating system should be Linux."
	@echo "                               Optimized binary with stripped symbols"
	@echo ""
	@echo "$(CYAN)⚙️  SERVICE MANAGEMENT:$(RESET)"
	@echo "  $(GREEN)start$(RESET)                        Start service in background"
	@echo "                               Usage: make start target=rpc_system"
	@echo "                               Logs output to target/[service]/service.log"
	@echo ""
	@echo "  $(GREEN)stop$(RESET)                         Stop running service"
	@echo "                               Graceful shutdown with SIGTERM"
	@echo "                               Force kill after 5 seconds if needed"
	@echo ""
	@echo "  $(GREEN)restart$(RESET)                      Stop and start service"
	@echo "                               Equivalent to: make stop start"
	@echo ""
	@echo "$(CYAN)🐳 CONTAINERIZATION:$(RESET)"
	@echo "  $(GREEN)image$(RESET)                        Build Docker image"
	@echo "                               Usage: make image target=rpc_system"
	@echo "                               Tags: [service]:$(VERSION) and [service]:latest"
	@echo "                               Includes build args and metadata"
	@echo ""
	@echo "  $(GREEN)run$(RESET)                          Run service in Docker container"
	@echo "                               Usage: make run target=rpc_system p=8001"
	@echo "                               Auto-restart enabled"
	@echo "                               Port mapping and health checks"
	@echo ""
	@echo "$(CYAN)🔧 UTILITY COMMANDS:$(RESET)"
	@echo "  $(GREEN)version$(RESET)                      Display version and build information"
	@echo "                               Show Git info, Go version, system details"
	@echo ""
	@echo "  $(GREEN)help$(RESET)                         Display this comprehensive help"
	@echo ""
	@echo "$(CYAN)📝 PARAMETER REFERENCE:$(RESET)"
	@echo "  $(YELLOW)target$(RESET)    Service name (e.g., system, admin, order_service)"
	@echo "  $(YELLOW)entity$(RESET)    Database entity name (e.g., User, Product, Order)"
	@echo "  $(YELLOW)os$(RESET)        Target OS: linux, darwin, windows"
	@echo "  $(YELLOW)arch$(RESET)      Target architecture: amd64, arm64"
	@echo "  $(YELLOW)dsn$(RESET)       Database connection: username:password@host:port/database"
	@echo "  $(YELLOW)default$(RESET)   Default language: zh-CN, ar-EG, en-US, etc."
	@echo "  $(YELLOW)p$(RESET)         Port number for container exposure"
	@echo ""
	@echo "$(CYAN)💡 COMMON WORKFLOWS:$(RESET)"
	@echo ""
	@echo "  $(MAGENTA)Creating a new RPC service:$(RESET)"
	@echo "    make init"
	@echo "    make ent_new target=payment_system entity=Payment"
	@echo "    make ent_new target=payment_system entity=Transaction"
	@echo "    # Edit schema files to add fields and relationships"
	@echo "    make ent_gen target=payment_system"
	@echo "    make rpc_pb target=payment_system"
	@echo "    make rpc_gen target=payment_system"
	@echo "    make rpc_build target=payment_system"
	@echo "    make start target=rpc_payment_system"
	@echo ""
	@echo "  $(MAGENTA)Creating a new API service:$(RESET)"
	@echo "    make init"
	@echo "    # Create API definition files in api/admin/desc/"
	@echo "    make api_format target=admin"
	@echo "    make api_gen target=admin"
	@echo "    make api_build target=admin"
	@echo "    make start target=api_admin"
	@echo ""
	@echo "  $(MAGENTA)Database migration workflow:$(RESET)"
	@echo "    # After modifying entity schemas"
	@echo "    make ent_gen target=system"
	@echo "    make ent_migrate_diff target=system"
	@echo "    # Review generated migration files"
	@echo "    make ent_migrate_apply target=system dsn=root:password@localhost:3306/mydb"
	@echo ""
	@echo "  $(MAGENTA)Production deployment:$(RESET)"
	@echo "    make clean"
	@echo "    make tidy"
	@echo "    make lint"
	@echo "    make test"
	@echo "    make rpc_build target=system os=linux arch=amd64"
	@echo "    make image target=rpc_system"
	@echo "    make run target=rpc_system p=8080"
	@echo ""
	@echo "$(CYAN)🚨 TROUBLESHOOTING:$(RESET)"
	@echo ""
	@echo "  $(RED)Service won't start:$(RESET)"
	@echo "    • Check logs: tail -f target/[service]/service.log"
	@echo "    • Verify config: cat target/[service]/[service].yaml"
	@echo "    • Check ports: make status target=[service]"
	@echo ""
	@echo "  $(RED)Build failures:$(RESET)"
	@echo "    • Update dependencies: make tidy"
	@echo "    • Check Go version: go version"
	@echo "    • Regenerate code: make ent_gen rpc_gen api_gen"
	@echo ""
	@echo "  $(RED)Docker issues:$(RESET)"
	@echo "    • Check Docker: docker --version"
	@echo "    • Rebuild image: make clean && make image"
	@echo ""
	@echo "$(GREEN)📚 For more information:$(RESET)"
	@echo "  • Go-Zero: https://go-zero.dev/"
	@echo "  • Ent ORM: https://entgo.io/"
	@echo "  • Atlas CLI: https://atlasgo.io/"
	@echo ""
