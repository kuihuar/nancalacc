#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 目录配置
BIN_DIR="bin"
CONFIG_DIR="configs"
CONFIG_FILE="config.yaml"

# 默认参数
BINARY_NAME="nancalacc-linux-amd64"
IMAGE_NAME="nancalacc"
TAG="latest"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -b|--binary)
            BINARY_NAME="$2"
            shift 2
            ;;
        -i|--image)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -h|--help)
            echo "用法: $0 [选项]"
            echo "选项:"
            echo "  -b, --binary NAME    指定二进制文件名 (默认: nancalacc-linux-amd64)"
            echo "  -i, --image NAME     指定镜像名 (默认: nancalacc)"
            echo "  -t, --tag TAG        指定标签 (默认: latest)"
            echo "  -h, --help           显示帮助信息"
            exit 0
            ;;
        *)
            log_error "未知参数: $1"
            exit 1
            ;;
    esac
done

# 检查必要文件
log_info "检查构建文件..."
if [ ! -f "$BIN_DIR/$BINARY_NAME" ]; then
    log_error "二进制文件不存在: $BIN_DIR/$BINARY_NAME"
    log_info "请先运行: ./build.sh -b --env prod"
    exit 1
fi

if [ ! -f "$CONFIG_DIR/$CONFIG_FILE" ]; then
    log_error "配置文件不存在: $CONFIG_DIR/$CONFIG_FILE"
    exit 1
fi

log_success "文件检查通过"

# 显示构建信息
log_info "构建信息:"
log_info "  二进制文件: $BIN_DIR/$BINARY_NAME ($(ls -lh $BIN_DIR/$BINARY_NAME | awk '{print $5}'))"
log_info "  配置文件: $CONFIG_DIR/$CONFIG_FILE ($(ls -lh $CONFIG_DIR/$CONFIG_FILE | awk '{print $5}'))"
log_info "  镜像名: $IMAGE_NAME:$TAG"

# 构建 Docker 镜像
log_info "开始构建 Docker 镜像..."
if docker build \
    --build-arg BINARY_NAME="$BINARY_NAME" \
    -t "$IMAGE_NAME:$TAG" \
    .; then
    log_success "Docker 镜像构建成功: $IMAGE_NAME:$TAG"
    
    # 显示镜像信息
    log_info "镜像信息:"
    docker images "$IMAGE_NAME:$TAG" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
else
    log_error "Docker 镜像构建失败"
    exit 1
fi 