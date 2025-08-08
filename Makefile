# 基础变量定义
VERSION:=$(shell git describe --tags --always 2>/dev/null || echo "0.0.0")
LOCAL_GOHOSTOS:=$(shell go env GOHOSTOS)
LOCAL_GOARCH:=$(shell go env GOARCH)
LOCAL_GOPATH:=$(shell go env GOPATH)

# 构建相关变量
BUILD_TIME:=$(shell date +%Y-%m-%d_%H:%M:%S)
BUILD_USER:=$(shell whoami)
BUILD_HOST:=$(shell hostname)

# Ent 特性配置
ENT_FEATURE:=sql/execquery,sql/modifier,intercept

# 默认目标
.DEFAULT_GOAL := help

.PHONY: init tidy format atlas ent_new ent_gen combine gen_rpc gen_api build_rpc build_api start stop restart image run help

# 初始化环境
init:
	@echo "开始初始化环境..."
	@go install github.com/zeromicro/go-zero/tools/goctl@latest
	@goctl env check --install --verbose --force
	@echo "环境初始化完成"

# 更新依赖
tidy:
	@echo "开始更新依赖..."
	@export GOPROXY=https://goproxy.cn,direct
	@go mod tidy -v
	@echo "依赖更新完成"

# 格式化 *.api target=admin
format:
	@echo "开始格式化 API 文件..."
	@goctl api format --dir api/$(target)/desc/
	@echo "格式化完成"

############################################# Ent #################################################

# visualize the schema target=admin_system
atlas:
	@echo "开始生成数据库可视化..."
	@atlas schema inspect -u "ent://rpc/$(target)/ent/schema" --dev-url "docker+mysql://_/mysql:8.4-oracle/dev" -w
	@echo "可视化完成"

# ent new target=admin_system entity=User
ent_new:
	@echo "开始创建新的 Ent 实体..."
	@go run -mod=mod entgo.io/ent/cmd/ent new --target=rpc/$(target)/ent/schema $(entity)
	@echo "实体创建完成"

# ent generate target=admin_system
ent_gen:
	@echo "开始生成 Ent 代码..."
	@go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./rpc/$(target)/ent/template/*.tmpl" ./rpc/$(target)/ent/schema --feature $(ENT_FEATURE)
	@echo "代码生成完成"

############################################# GEN ################################################

# combine *.proto target=admin_system
combine:
	@echo "开始合并 Proto 文件..."
	@go run ./rpc/$(target)/desc/main.go
	@echo "合并完成"

# 项目生成 rpc target=admin_system
gen_rpc:
	@echo "开始生成 RPC 服务..."
	@goctl rpc protoc ./rpc/$(target)/$(target).proto --go_out=./rpc/$(target)/ --go-grpc_out=./rpc/$(target)/ --zrpc_out=./rpc/$(target)/ -m --style=go_zero
	@echo "RPC 服务生成完成"

# 项目生成 api target=admin
gen_api:
	@echo "开始生成 API 服务..."
	@goctl api go -api ./api/$(target)/desc/$(target).api -dir ./api/$(target)/ --style=go_zero
	@echo "API 服务生成完成"

############################################# Build #############################################

# 构建 rpc os=windows|darwin|linux arch=amd64|arm64 ext=.exe target=admin_system
build_rpc:
	@echo "开始构建 RPC 服务..."
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
	@echo "RPC 服务构建完成"

# 构建 api os=windows|darwin|linux arch=amd64|arm64 ext=.exe target=admin
build_api:
	@echo "开始构建 API 服务..."
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
	@echo "API 服务构建完成"

############################################# RUN ################################################

# 裸机运行 target=rpc_admin_system
start:
	@echo "开始启动服务..."
	@nohup ./target/$(target)/$(target) -f ./target/$(target)/$(target).yaml > /dev/null 2>&1 &
	@echo "服务启动完成"

# 裸机停止 target=rpc_admin_system
stop:
	@echo "开始停止服务..."
	@-pkill -f $(target)
	@for i in 5 4 3 2 1; do \
		echo -n "停止中 $$i"; \
		sleep 1; \
		echo " "; \
	done
	@echo "服务停止完成"

# 裸机重启
restart: stop start

# 构建容器镜像 target=rpc_admin_system
image:
	@echo "开始构建 Docker 镜像..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg BUILD_HOST=$(BUILD_HOST) \
		-t $(target):$(VERSION) \
		-f target/$(target)/Dockerfile .
	@echo "Docker 镜像构建完成"

# 启动容器 target=rpc_admin_system p=8001
run:
	@echo "开始运行容器..."
	@docker run -itd -p $(p):$(p) --name=$(target) $(target):$(VERSION)
	@echo "容器运行完成"

# SHOW Help
help:
	@echo ''
	@echo '使用方法:'
	@echo ' make [target]'
	@echo ''
	@echo '目标列表:'
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
	@echo '参数说明:'
	@echo '  VERSION: 项目版本号'
	@echo '  BUILD_TIME: 构建时间'
	@echo '  BUILD_USER: 构建用户'
	@echo '  BUILD_HOST: 构建主机'
	@echo '  ENT_FEATURE: Ent 特性配置'
