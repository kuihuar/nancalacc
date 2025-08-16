#!/bin/bash

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

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 目录配置
BIN_DIR="bin"
CONFIG_DIR="configs"
CONFIG_FILE="config.yaml"
DOCKERFILE_DIR="."
DOCKERFILE_NAME="Dockerfile"

# 目录配置
BIN_DIR="bin"
CONFIG_DIR="configs"
CONFIG_FILE="config.yaml"
BINARY_NAME="nancalacc-linux-amd64"

# 检查文件是否存在
check_files() {
    log_info "检查必要文件..."
    
    # 检查二进制文件
    if [ ! -f "$BIN_DIR/nancalacc-linux-amd64" ]; then
        log_error "二进制文件不存在: $BIN_DIR/nancalacc-linux-amd64"
        log_info "请先运行: ./build.sh -b --env prod"
        exit 1
    fi
    
    # 检查配置文件
    if [ ! -f "$CONFIG_DIR/$CONFIG_FILE" ]; then
        log_error "配置文件不存在: $CONFIG_DIR/$CONFIG_FILE"
        exit 1
    fi
    
    log_success "所有文件检查通过"
}

# 测试 Docker 构建
test_build() {
    log_info "测试 Docker 构建..."
    
    # 显示文件信息
    log_info "二进制文件大小: $(ls -lh $BIN_DIR/nancalacc-linux-amd64 | awk '{print $5}')"
    log_info "配置文件大小: $(ls -lh $CONFIG_DIR/$CONFIG_FILE | awk '{print $5}')"
    
    # 构建测试镜像
    if docker build --build-arg BINARY_NAME="nancalacc-linux-amd64" -t test-nancalacc .; then
        log_success "Docker 构建成功！"
        
        # 清理测试镜像
        docker rmi test-nancalacc > /dev/null 2>&1
        log_info "测试镜像已清理"
    else
        log_error "Docker 构建失败"
        exit 1
    fi
}

# 主函数
main() {
    log_info "开始 Docker 构建测试..."
    check_files
    test_build
    log_success "测试完成！"
}

main "$@" 