#!/bin/bash

# Nancalacc Optimization Build Script
# 用于构建和部署 Go 微服务项目

set -e  # 遇到错误立即退出
set -x  # 显示执行的命令
set -v  # 显示执行的命令和输出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help         显示帮助信息"
    echo "  -c, --clean        清理构建产物"
    echo "  -t, --test         运行测试"
    echo "  -b, --build        构建项目"
    echo "  -d, --docker       构建Docker镜像"
    echo "  -p, --push         推送Docker镜像"
    echo "  -a, --all          执行完整构建流程"
    echo "  --env ENV          指定环境 (dev/test/prod)"
    echo "  --version VERSION  指定版本号"
    echo ""
    echo "Examples:"
    echo "  $0 -a                    # 完整构建流程 (包含 CGC pack/and/bundle)"
    echo "  $0 -b --env prod         # 构建生产环境版本 (包含 CGC 打包)"
    echo "  $0 -d -p --version 1.0.0 # 构建并推送Docker镜像"
    echo "  $0 -t                    # 仅运行测试"
    echo "  $0 -c                    # 仅清理构建产物"
}

# 默认配置
ENV=${ENV:-"dev"}
VERSION=${VERSION:-"1.0.0"}
PROJECT_NAME="nancalacc"
DOCKER_REGISTRY=${DOCKER_REGISTRY:-"your-registry.com"}
DOCKER_IMAGE="${DOCKER_REGISTRY}/${PROJECT_NAME}:${VERSION}"

# 目录配置 - 根据 cgc-build1.sh 简化
PROJECT_DIR=$(pwd)
BINARY_DIR="$PROJECT_DIR/bin"
BUILD_DIR_CGC="$PROJECT_DIR/buildcgc"
BUILD_DIR="$BUILD_DIR_CGC/build"
PACK_DIR="$BUILD_DIR_CGC/pack"
BUNDLE_DIR="$BUILD_DIR_CGC/bundle"
CONFIG_DIR="configs"
CONFIG_FILE="config.yaml"
README_FILE="README.md"
CMD_DIR="cmd"
MAIN_PACKAGE="nancalacc"
DOCKERFILE_DIR="."
DOCKERFILE_NAME="Dockerfile"
BINARY_NAME="${PROJECT_NAME}"

# 生成二进制文件名的统一函数
generate_binary_name() {
    local target_os=${1:-"linux"}
    local target_arch=${2:-"amd64"}
    
    # 根据目标平台确定 OS 和 ARCH
    case $target_os in
        "linux")
            GOOS="linux"
            ;;
        "darwin")
            GOOS="darwin"
            ;;
        *)
            GOOS="linux"
            ;;
    esac
    
    case $target_arch in
        "amd64")
            GOARCH="amd64"
            ;;
        "arm64")
            GOARCH="arm64"
            ;;
        *)
            GOARCH="amd64"
            ;;
    esac
    
    echo "${BINARY_NAME}-${GOOS}-${GOARCH}"
}

# 检查依赖
check_dependencies() {
    log_info "检查构建依赖..."
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go 1.19+"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_info "Go 版本: $GO_VERSION"
    
    # 检查 Docker (如果需要)
    if [[ "$*" == *"--docker"* ]] || [[ "$*" == *"-d"* ]]; then
        if ! command -v docker &> /dev/null; then
            log_error "Docker 未安装，请先安装 Docker"
            exit 1
        fi
        log_info "Docker 已安装"
    fi
    
    # 检查 CGC 工具 (如果需要)
    if [[ "$*" == *"--build"* ]] || [[ "$*" == *"-b"* ]] || [[ "$*" == *"--all"* ]] || [[ "$*" == *"-a"* ]]; then
        if ! command -v cgc &> /dev/null; then
            log_error "CGC 工具未安装，请先安装 CGC"
            exit 1
        fi
        log_info "CGC 工具已安装"
    fi
    
    log_success "依赖检查完成"
}

# 清理构建产物
clean_build() {
    log_info "清理构建产物..."
    
    # 清理构建目录
    if [ -d "$BUILD_DIR" ]; then
        rm -rf "$BUILD_DIR"
        log_info "已清理构建目录: $BUILD_DIR"
    fi
    
    # 清理 Go 缓存
    # log_info "已清理 Go 缓存"
    
    # 清理 Docker 镜像 (如果存在)
    if docker images | grep -q "$PROJECT_NAME"; then
        docker rmi $(docker images | grep "$PROJECT_NAME" | awk '{print $3}') 2>/dev/null || true
        log_info "已清理 Docker 镜像"
    fi
    
    log_success "清理完成"
}

# 下载依赖
download_deps() {
    log_info "下载 Go 依赖..."
    
    # 设置 Go 代理 (可选)
    # export GOPROXY=https://goproxy.cn,direct
    
    go mod download
    go mod tidy
    
    log_success "依赖下载完成"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    
    # 运行单元测试
    go test -v ./...
    
    # 运行测试并生成覆盖率报告
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    
    log_success "测试完成，覆盖率报告已生成: coverage.html"
}

# 代码质量检查
code_quality_check() {
    log_info "执行代码质量检查..."
    
    # 检查代码格式
    if ! go fmt ./...; then
        log_error "代码格式检查失败"
        exit 1
    fi
    
    # 运行 golint (如果安装)
    if command -v golint &> /dev/null; then
        golint ./... || log_warning "golint 检查发现问题"
    fi
    
    # 运行 staticcheck (如果安装)
    if command -v staticcheck &> /dev/null; then
        staticcheck ./... || log_warning "staticcheck 检查发现问题"
    fi
    
    log_success "代码质量检查完成"
}

# 构建项目
build_project() {
    log_info "构建项目 (环境: $ENV, 版本: $VERSION)..."
    
    # 创建构建目录
    mkdir -p "$BINARY_DIR"
    
    # 设置构建参数
    LDFLAGS="-X main.Version=$VERSION -X main.Environment=$ENV -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')"
    
    # 生成包含平台信息的文件名
    PLATFORM_BINARY_NAME=$(generate_binary_name "linux" "amd64")
    
    # 根据环境设置不同的构建参数
    case $ENV in
        "prod")
            CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BINARY_DIR/$PLATFORM_BINARY_NAME" ./$CMD_DIR/$MAIN_PACKAGE
            ;;
        "test")
            CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -race -o "$BINARY_DIR/$PLATFORM_BINARY_NAME" ./$CMD_DIR/$MAIN_PACKAGE
            ;;
        *)
            GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o "$BINARY_DIR/$PLATFORM_BINARY_NAME" ./$CMD_DIR/$MAIN_PACKAGE
            ;;
    esac
    
    # 复制配置文件
    # if [ -d "$CONFIG_DIR" ]; then
    #     cp -r "$CONFIG_DIR" "$BUILD_DIR/"
    #     log_info "配置文件已复制到构建目录: $BUILD_DIR/$CONFIG_DIR"
    # fi
    
    # # 复制其他必要文件
    # if [ -f "$README_FILE" ]; then
    #     cp "$README_FILE" "$BUILD_DIR/"
    # fi
    
    log_success "项目构建完成: $BINARY_DIR/$PLATFORM_BINARY_NAME"
}

# 执行 CGC 打包操作 - 根据 cgc-build1.sh 修正
execute_cgc_pack() {
    log_info "执行 CGC 打包操作..."
    
    # 检查 cgc 命令是否存在
    if ! command -v cgc &> /dev/null; then
        log_error "cgc 命令未找到，请确保已安装 CGC 工具"
        exit 1
    fi
    
    # 复制 Dockerfile 到构建目录
    # if [ -f "$DOCKERFILE_DIR/$DOCKERFILE_NAME" ]; then
    #     cp "$DOCKERFILE_DIR/$DOCKERFILE_NAME" "$BUILD_DIR/"
    #     log_info "Dockerfile 已复制到构建目录: $BUILD_DIR/$DOCKERFILE_NAME"
    # else
    #     log_error "Dockerfile 不存在: $DOCKERFILE_DIR/$DOCKERFILE_NAME"
    #     exit 1
    # fi
    
    # 进入构建目录
    #cd "$BUILD_DIR"

    # 执行 cgc build -f -a amd64 -o $BUILD_DIR -e Dockerfile
    #log_info "执行: cgc build -f -a amd64 -o $BUILD_DIR -e Dockerfile"
    #ls -al $BUILD_DIR

    BINARY_NAME="nancalacc-linux-amd64-$(date +%Y%m%d%H%M%S)"

    if ! cgc build -f -a amd64 --build-arg BINARY_NAME=nancalacc-linux-amd64 -o $BUILD_DIR -e Dockerfile; then
        log_error "cgc build 执行失败"
        exit 1
    fi
    
    # 执行 cgc pack -i $BUILD_DIR -o $PACK_DIR
    log_info "执行: cgc pack -i $BUILD_DIR -o $PACK_DIR"
    if ! cgc pack -i $BUILD_DIR -o $PACK_DIR; then
        log_error "cgc pack 执行失败"
        exit 1
    fi
    
    # 执行 cgc bundle --dir $PACK_DIR -o $BUNDLE_DIR
    log_info "执行: cgc bundle --dir $PACK_DIR -o $BUILD_DIR"
    if ! cgc bundle --dir $PACK_DIR -o ~/Downloads; then
        log_error "cgc bundle 执行失败"
        exit 1
    fi
    
    # 返回原目录
    cd - > /dev/null
    
    log_success "CGC 打包操作完成"
}

# 构建 Docker 镜像
build_docker() {
    log_info "构建 Docker 镜像: $DOCKER_IMAGE"
    
    # 检查 Dockerfile 是否存在
    if [ ! -f "$DOCKERFILE_DIR/$DOCKERFILE_NAME" ]; then
        log_error "Dockerfile 不存在: $DOCKERFILE_DIR/$DOCKERFILE_NAME"
        exit 1
    fi
    
    # 计算平台信息，Docker 镜像需要 Linux 版本
    PLATFORM_BINARY_NAME=$(generate_binary_name "linux" "amd64")
    
    # 检查二进制文件是否存在
    if [ ! -f "$BUILD_DIR/$PLATFORM_BINARY_NAME" ]; then
        log_error "Linux 二进制文件不存在: $BUILD_DIR/$PLATFORM_BINARY_NAME"
        log_info "请先构建 Linux 版本: $0 -b --env prod"
        exit 1
    fi
    
    # 构建镜像，传递二进制文件名参数
    docker build --build-arg BINARY_NAME="$PLATFORM_BINARY_NAME" -t "$DOCKER_IMAGE" .
    
    # 创建 latest 标签
    docker tag "$DOCKER_IMAGE" "${DOCKER_REGISTRY}/${PROJECT_NAME}:latest"
    
    log_success "Docker 镜像构建完成: $DOCKER_IMAGE"
}

# 推送 Docker 镜像
push_docker() {
    log_info "推送 Docker 镜像到仓库..."
    
    # 推送指定版本
    docker push "$DOCKER_IMAGE"
    
    # 推送 latest 标签
    docker push "${DOCKER_REGISTRY}/${PROJECT_NAME}:latest"
    
    log_success "Docker 镜像推送完成"
}

# 部署检查
deploy_check() {
    log_info "执行部署前检查..."
    
    # 计算平台信息
    PLATFORM_BINARY_NAME=$(generate_binary_name "linux" "amd64")
    
    # 检查二进制文件
    if [ ! -f "$BUILD_DIR/$PLATFORM_BINARY_NAME" ]; then
        log_error "二进制文件不存在: $BUILD_DIR/$PLATFORM_BINARY_NAME"
        exit 1
    fi
    
    # 检查配置文件
    if [ ! -d "$BUILD_DIR/$CONFIG_DIR" ]; then
        log_warning "配置文件目录不存在: $BUILD_DIR/$CONFIG_DIR"
    fi
    
    # 检查文件权限
    chmod +x "$BUILD_DIR/$PLATFORM_BINARY_NAME"
    
    log_success "部署检查完成"
}

# 显示构建信息
show_build_info() {
    # 计算平台信息
    PLATFORM_BINARY_NAME=$(generate_binary_name "linux" "amd64")
    
    log_info "构建信息:"
    echo "  项目名称: $PROJECT_NAME"
    echo "  版本: $VERSION"
    echo "  环境: $ENV"
    echo "  项目目录: $PROJECT_DIR"
    echo "  构建目录: $BUILD_DIR"
    echo "  打包目录: $PACK_DIR"
    echo "  打包输出目录: $BUNDLE_DIR"
    echo "  二进制文件: $BUILD_DIR/$PLATFORM_BINARY_NAME"
    if [[ "$*" == *"--docker"* ]] || [[ "$*" == *"-d"* ]]; then
        echo "  Docker 镜像: $DOCKER_IMAGE"
    fi
}

# 主函数
main() {
    local CLEAN=false
    local TEST=false
    local BUILD=false
    local DOCKER=false
    local PUSH=false
    local ALL=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--clean)
                CLEAN=true
                shift
                ;;
            -t|--test)
                TEST=true
                shift
                ;;
            -b|--build)
                BUILD=true
                shift
                ;;
            -d|--docker)
                DOCKER=true
                shift
                ;;
            -p|--push)
                PUSH=true
                shift
                ;;
            -a|--all)
                ALL=true
                shift
                ;;
            --env)
                ENV="$2"
                shift 2
                ;;
            --version)
                VERSION="$2"
                shift 2
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定任何操作，显示帮助
    if [[ "$CLEAN" == false && "$TEST" == false && "$BUILD" == false && "$DOCKER" == false && "$PUSH" == false && "$ALL" == false ]]; then
        show_help
        exit 0
    fi
    
    # 显示构建信息
    show_build_info "$@"
    
    # 检查依赖
    check_dependencies "$@"
    
    # 执行清理
    if [[ "$CLEAN" == true ]]; then
        clean_build
    fi
    
    # 执行完整流程
    if [[ "$ALL" == true ]]; then
        log_info "开始执行完整构建流程..."
        clean_build
        download_deps
        code_quality_check
        run_tests
        build_project
        execute_cgc_pack
        if [[ "$DOCKER" == true ]]; then
            build_docker
            if [[ "$PUSH" == true ]]; then
                push_docker
            fi
        fi
        deploy_check
        log_success "完整构建流程完成！"
        exit 0
    fi
    
    # 执行单独的操作
    if [[ "$TEST" == true ]]; then
        download_deps
        run_tests
    fi
    
    if [[ "$BUILD" == true ]]; then
        download_deps
        code_quality_check
        build_project
        execute_cgc_pack
        deploy_check
    fi
    
    if [[ "$DOCKER" == true ]]; then
        if [[ "$BUILD" == false ]]; then
            download_deps
            build_project
            execute_cgc_pack
        fi
        build_docker
    fi
    
    if [[ "$PUSH" == true ]]; then
        if [[ "$DOCKER" == false ]]; then
            log_error "推送镜像前需要先构建 Docker 镜像"
            exit 1
        fi
        push_docker
    fi
    
    log_success "构建完成！"
}

# 执行主函数
main "$@"
