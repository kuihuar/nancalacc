#!/bin/bash

# Nancalacc Account API Examples
# 使用curl命令测试账户管理API

# 设置基础URL
BASE_URL="http://localhost:8000"
API_VERSION="v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查服务是否运行
check_service() {
    print_info "检查服务状态..."
    if curl -s "$BASE_URL/health" > /dev/null; then
        print_success "服务运行正常"
        return 0
    else
        print_error "服务未运行，请先启动服务"
        return 1
    fi
}

# 创建同步账户
create_sync_account() {
    print_info "创建同步账户..."
    
    local name="${1:-test-account}"
    local email="${2:-test@example.com}"
    
    response=$(curl -s -X POST "$BASE_URL/$API_VERSION/account" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{
            \"name\": \"$name\",
            \"email\": \"$email\"
        }")
    
    if [ $? -eq 0 ]; then
        print_success "账户创建成功"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "账户创建失败"
        echo "$response"
    fi
}

# 获取账户列表
get_accounts() {
    print_info "获取账户列表..."
    
    response=$(curl -s -X GET "$BASE_URL/$API_VERSION/accounts" \
        -H "Authorization: Bearer $TOKEN")
    
    if [ $? -eq 0 ]; then
        print_success "获取账户列表成功"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "获取账户列表失败"
        echo "$response"
    fi
}

# 获取单个账户
get_account() {
    local account_id="$1"
    
    if [ -z "$account_id" ]; then
        print_error "请提供账户ID"
        return 1
    fi
    
    print_info "获取账户信息: $account_id"
    
    response=$(curl -s -X GET "$BASE_URL/$API_VERSION/account/$account_id" \
        -H "Authorization: Bearer $TOKEN")
    
    if [ $? -eq 0 ]; then
        print_success "获取账户信息成功"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "获取账户信息失败"
        echo "$response"
    fi
}

# 更新账户
update_account() {
    local account_id="$1"
    local name="$2"
    local email="$3"
    
    if [ -z "$account_id" ]; then
        print_error "请提供账户ID"
        return 1
    fi
    
    print_info "更新账户: $account_id"
    
    local data="{}"
    if [ -n "$name" ]; then
        data=$(echo "$data" | jq ".name = \"$name\"")
    fi
    if [ -n "$email" ]; then
        data=$(echo "$data" | jq ".email = \"$email\"")
    fi
    
    response=$(curl -s -X PUT "$BASE_URL/$API_VERSION/account/$account_id" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$data")
    
    if [ $? -eq 0 ]; then
        print_success "账户更新成功"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "账户更新失败"
        echo "$response"
    fi
}

# 删除账户
delete_account() {
    local account_id="$1"
    
    if [ -z "$account_id" ]; then
        print_error "请提供账户ID"
        return 1
    fi
    
    print_warning "删除账户: $account_id"
    read -p "确认删除账户 $account_id? (y/N): " confirm
    
    if [[ $confirm =~ ^[Yy]$ ]]; then
        response=$(curl -s -X DELETE "$BASE_URL/$API_VERSION/account/$account_id" \
            -H "Authorization: Bearer $TOKEN")
        
        if [ $? -eq 0 ]; then
            print_success "账户删除成功"
            echo "$response" | jq '.' 2>/dev/null || echo "$response"
        else
            print_error "账户删除失败"
            echo "$response"
        fi
    else
        print_info "取消删除操作"
    fi
}

# 启动同步任务
start_sync() {
    local account_id="$1"
    
    if [ -z "$account_id" ]; then
        print_error "请提供账户ID"
        return 1
    fi
    
    print_info "启动同步任务: $account_id"
    
    response=$(curl -s -X POST "$BASE_URL/$API_VERSION/account/$account_id/sync" \
        -H "Authorization: Bearer $TOKEN")
    
    if [ $? -eq 0 ]; then
        print_success "同步任务启动成功"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "同步任务启动失败"
        echo "$response"
    fi
}

# 获取同步状态
get_sync_status() {
    local account_id="$1"
    
    if [ -z "$account_id" ]; then
        print_error "请提供账户ID"
        return 1
    fi
    
    print_info "获取同步状态: $account_id"
    
    response=$(curl -s -X GET "$BASE_URL/$API_VERSION/account/$account_id/sync/status" \
        -H "Authorization: Bearer $TOKEN")
    
    if [ $? -eq 0 ]; then
        print_success "获取同步状态成功"
        echo "$response" | jq '.' 2>/dev/null || echo "$response"
    else
        print_error "获取同步状态失败"
        echo "$response"
    fi
}

# 显示帮助信息
show_help() {
    echo "Nancalacc Account API 测试脚本"
    echo ""
    echo "用法: $0 [命令] [参数]"
    echo ""
    echo "命令:"
    echo "  check                   检查服务状态"
    echo "  create [name] [email]   创建同步账户"
    echo "  list                    获取账户列表"
    echo "  get <id>                获取单个账户"
    echo "  update <id> [name] [email] 更新账户"
    echo "  delete <id>             删除账户"
    echo "  sync <id>               启动同步任务"
    echo "  status <id>             获取同步状态"
    echo "  help                    显示此帮助信息"
    echo ""
    echo "环境变量:"
    echo "  TOKEN                   认证令牌"
    echo "  BASE_URL                服务基础URL (默认: http://localhost:8000)"
    echo ""
    echo "示例:"
    echo "  $0 create test-account test@example.com"
    echo "  $0 get 12345"
    echo "  $0 sync 12345"
}

# 主函数
main() {
    # 检查jq是否安装
    if ! command -v jq &> /dev/null; then
        print_warning "jq未安装，JSON输出将不会格式化"
    fi
    
    # 设置默认值
    TOKEN="${TOKEN:-}"
    BASE_URL="${BASE_URL:-http://localhost:8000}"
    
    case "$1" in
        "check")
            check_service
            ;;
        "create")
            create_sync_account "$2" "$3"
            ;;
        "list")
            get_accounts
            ;;
        "get")
            get_account "$2"
            ;;
        "update")
            update_account "$2" "$3" "$4"
            ;;
        "delete")
            delete_account "$2"
            ;;
        "sync")
            start_sync "$2"
            ;;
        "status")
            get_sync_status "$2"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        "")
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@" 