# Basic variable definitions
VERSION:=$(shell git describe --tags --always 2>/dev/null || echo "0.0.0")
LOCAL_GOHOSTOS:=$(shell go env GOHOSTOS)
LOCAL_GOARCH:=$(shell go env GOARCH)
LOCAL_GOPATH:=$(shell go env GOPATH)

# Build-related variables
BUILD_TIME:=$(shell date +%Y-%m-%d_%H:%M:%S)
BUILD_USER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)

# Ent feature configuration
ENT_FEATURE:=sql/execquery,sql/modifier,intercept

# Default target
.DEFAULT_GOAL := help

.PHONY: init tidy format atlas ent_new ent_gen combine gen_rpc gen_api build_rpc build_api start stop restart image run help

# Initialize environment
init:
	@echo "Initializing environment..."
	@go install github.com/zeromicro/go-zero/tools/goctl@latest
	@goctl env check --install --verbose --force
	@echo "Environment initialization completed"

# Update dependencies
tidy:
	@echo "Updating dependencies..."
	@go mod tidy -v
	@echo "Dependency update completed"

# Format *.api files target=admin
format:
	@echo "Formatting API files..."
	@goctl api format --dir api/$(target)/desc/
	@echo "Formatting completed"

############################################# Ent #################################################

# Visualize database schema target=admin_system
atlas:
	@echo "Generating database visualization..."
	@atlas schema inspect -u "ent://rpc/$(target)/ent/schema" --dev-url "docker+mysql://_/mysql:8.4-oracle/dev" -w
	@echo "Visualization completed"

# Create new Ent entity target=admin_system entity=User
ent_new:
	@echo "Creating new Ent entity..."
	@go run -mod=mod entgo.io/ent/cmd/ent new --target=rpc/$(target)/ent/schema $(entity)
	@echo "Entity creation completed"

# Generate Ent code target=admin_system
ent_gen:
	@echo "Generating Ent code..."
	@go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./rpc/$(target)/ent/template/*.tmpl" ./rpc/$(target)/ent/schema --feature $(ENT_FEATURE)
	@echo "Code generation completed"

############################################# GEN ################################################

# Combine *.proto files target=admin_system
combine:
	@echo "Merging Proto files..."
	@go run ./rpc/$(target)/desc/main.go
	@echo "Merge completed"

# Generate RPC service target=admin_system
gen_rpc:
	@echo "Generating RPC service..."
	@goctl rpc protoc ./rpc/$(target)/$(target).proto --go_out=./rpc/$(target)/ --go-grpc_out=./rpc/$(target)/ --zrpc_out=./rpc/$(target)/ -m --style=go_zero
	@echo "RPC service generation completed"

# Generate API service target=admin
gen_api:
	@echo "Generating API service..."
	@goctl api go -api ./api/$(target)/desc/$(target).api -dir ./api/$(target)/ --style=go_zero
	@echo "API service generation completed"

############################################# Build #############################################

# Build RPC service os=windows|darwin|linux arch=amd64|arm64 ext=.exe target=admin_system
build_rpc:
	@echo "Building RPC service..."
	@mkdir -p target/rpc_$(target)
	@env CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build \
		-ldflags "-s -w \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildTime=$(BUILD_TIME)' \
		-X 'main.BuildUser=$(BUILD_USER)' \
		-X 'main.BuildHost=$(BUILD_HOST)'" \
		-trimpath \
		-o target/rpc_$(target)/rpc_$(target)$(ext) \
		-v ./rpc/$(target)/$(target).go
	@echo "RPC service build completed"

# Build API service os=windows|darwin|linux arch=amd64|arm64 ext=.exe target=admin
build_api:
	@echo "Building API service..."
	@mkdir -p target/api_$(target)
	@env CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) go build \
		-ldflags "-s -w \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildTime=$(BUILD_TIME)' \
		-X 'main.BuildUser=$(BUILD_USER)' \
		-X 'main.BuildHost=$(BUILD_HOST)'" \
		-trimpath \
		-o target/api_$(target)/api_$(target)$(ext) \
		-v ./api/$(target)/$(target).go
	@echo "API service build completed"

############################################# RUN ################################################

# Run bare metal service target=rpc_admin_system
start:
	@echo "Starting service..."
	@nohup ./target/$(target)/$(target) -f ./target/$(target)/$(target).yaml > /dev/null 2>&1 &
	@echo "Service started"

# Stop bare metal service target=rpc_admin_system
stop:
	@echo "Stopping service..."
	@-pkill -f $(target)
	@for i in 5 4 3 2 1; do \
		echo -n "Stopping $$i"; \
		sleep 1; \
		echo " "; \
	done
	@echo "Service stopped"

# Restart bare metal service
restart: stop start

# Build container image target=rpc_admin_system
image:
	@echo "Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg BUILD_HOST=$(BUILD_HOST) \
		-t $(target):$(VERSION) \
		-f target/$(target)/Dockerfile .
	@echo "Docker image build completed"

# Run container target=rpc_admin_system p=8001
run:
	@echo "Starting container..."
	@docker run -itd -p $(p):$(p) --name=$(target) $(target):$(VERSION)
	@echo "Container running"

# SHOW Help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
	@echo ''
	@echo 'Parameter Description:'
	@echo '  VERSION: Project version number'
	@echo '  BUILD_TIME: Build timestamp'
	@echo '  BUILD_USER: Build user'
	@echo '  BUILD_HOST: Build hostname'
	@echo '  ENT_FEATURE: Ent feature configuration'
