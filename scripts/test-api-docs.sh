#!/bin/bash

# API文档生成测试脚本
# 用于验证API文档生成功能是否正常工作

set -e

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

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 未安装"
        return 1
    else
        print_success "$1 已安装"
        return 0
    fi
}

# 检查文件是否存在
check_file() {
    if [ -f "$1" ]; then
        print_success "文件存在: $1"
        return 0
    else
        print_error "文件不存在: $1"
        return 1
    fi
}

# 检查目录是否存在
check_directory() {
    if [ -d "$1" ]; then
        print_success "目录存在: $1"
        return 0
    else
        print_error "目录不存在: $1"
        return 1
    fi
}

# 测试环境检查
test_environment() {
    print_info "检查测试环境..."
    
    # 检查必要的命令
    check_command "protoc"
    check_command "go"
    check_command "make"
    check_command "curl"
    
    # 检查项目结构
    check_directory "api"
    check_directory "docs"
    check_file "Makefile"
}

# 测试API文档生成
test_api_generation() {
    print_info "测试API文档生成..."
    
    # 清理之前的生成文件
    if [ -d "docs/api/openapi" ]; then
        rm -rf docs/api/openapi/*
    fi
    
    # 生成OpenAPI文档
    print_info "生成OpenAPI文档..."
    if make generate-openapi; then
        print_success "OpenAPI文档生成成功"
    else
        print_error "OpenAPI文档生成失败"
        return 1
    fi
    
    # 检查生成的文件
    if [ -f "docs/api/openapi/account-v1.yaml" ]; then
        print_success "OpenAPI文件生成成功"
    else
        print_warning "未找到account-v1.yaml文件"
    fi
}

# 测试Swagger UI生成
test_swagger_ui() {
    print_info "测试Swagger UI生成..."
    
    # 生成Swagger UI
    print_info "生成Swagger UI..."
    if make generate-swagger-ui; then
        print_success "Swagger UI生成成功"
    else
        print_error "Swagger UI生成失败"
        return 1
    fi
    
    # 检查生成的文件
    check_file "docs/api/swagger/index.html"
    check_directory "docs/api/swagger/swagger-ui"
}

# 测试示例文件
test_examples() {
    print_info "测试示例文件..."
    
    # 检查curl示例
    check_file "docs/api/examples/curl/account-api.sh"
    
    # 检查Postman集合
    check_file "docs/api/examples/postman/nancalacc-api.postman_collection.json"
    check_file "docs/api/examples/postman/nancalacc-local.postman_environment.json"
    
    # 测试curl脚本
    print_info "测试curl脚本..."
    if [ -x "docs/api/examples/curl/account-api.sh" ]; then
        print_success "curl脚本可执行"
        
        # 测试帮助命令
        if ./docs/api/examples/curl/account-api.sh help &> /dev/null; then
            print_success "curl脚本帮助命令正常"
        else
            print_warning "curl脚本帮助命令异常"
        fi
    else
        print_error "curl脚本不可执行"
    fi
}

# 测试文档服务器
test_docs_server() {
    print_info "测试文档服务器..."
    
    # 检查是否可以启动服务器
    if command -v python3 &> /dev/null; then
        print_info "启动文档服务器进行测试..."
        
        # 在后台启动服务器
        cd docs/api && python3 -m http.server 8081 &
        SERVER_PID=$!
        
        # 等待服务器启动
        sleep 2
        
        # 测试服务器响应
        if curl -s http://localhost:8081/swagger/ > /dev/null; then
            print_success "文档服务器启动成功"
        else
            print_warning "文档服务器启动失败"
        fi
        
        # 停止服务器
        kill $SERVER_PID 2>/dev/null || true
    else
        print_warning "Python3未安装，跳过服务器测试"
    fi
}

# 验证OpenAPI文件
validate_openapi() {
    print_info "验证OpenAPI文件..."
    
    if [ -f "docs/api/openapi/account-v1.yaml" ]; then
        # 检查文件格式
        if python3 -c "import yaml; yaml.safe_load(open('docs/api/openapi/account-v1.yaml'))" 2>/dev/null; then
            print_success "OpenAPI文件格式正确"
        else
            print_error "OpenAPI文件格式错误"
            return 1
        fi
        
        # 检查必要字段
        if grep -q "openapi:" docs/api/openapi/account-v1.yaml; then
            print_success "OpenAPI版本字段存在"
        else
            print_error "OpenAPI版本字段缺失"
            return 1
        fi
        
        if grep -q "info:" docs/api/openapi/account-v1.yaml; then
            print_success "OpenAPI信息字段存在"
        else
            print_error "OpenAPI信息字段缺失"
            return 1
        fi
    else
        print_warning "未找到OpenAPI文件，跳过验证"
    fi
}

# 生成测试报告
generate_report() {
    print_info "生成测试报告..."
    
    local report_file="docs/api/test-report-$(date +%Y%m%d-%H%M%S).md"
    
    cat > "$report_file" << EOF
# API文档生成测试报告

**测试时间**: $(date)
**测试环境**: $(uname -s) $(uname -r)

## 测试结果

### 环境检查
- protoc: $(command -v protoc >/dev/null && echo "✅ 已安装" || echo "❌ 未安装")
- go: $(command -v go >/dev/null && echo "✅ 已安装" || echo "❌ 未安装")
- make: $(command -v make >/dev/null && echo "✅ 已安装" || echo "❌ 未安装")
- curl: $(command -v curl >/dev/null && echo "✅ 已安装" || echo "❌ 未安装")

### 文件检查
- Makefile: $(test -f Makefile && echo "✅ 存在" || echo "❌ 不存在")
- api目录: $(test -d api && echo "✅ 存在" || echo "❌ 不存在")
- docs目录: $(test -d docs && echo "✅ 存在" || echo "❌ 不存在")

### 生成文件
- OpenAPI文件: $(test -f docs/api/openapi/account-v1.yaml && echo "✅ 生成成功" || echo "❌ 生成失败")
- Swagger UI: $(test -f docs/api/swagger/index.html && echo "✅ 生成成功" || echo "❌ 生成失败")

### 示例文件
- curl脚本: $(test -f docs/api/examples/curl/account-api.sh && echo "✅ 存在" || echo "❌ 不存在")
- Postman集合: $(test -f docs/api/examples/postman/nancalacc-api.postman_collection.json && echo "✅ 存在" || echo "❌ 不存在")

## 建议

$(if [ -f docs/api/openapi/account-v1.yaml ]; then
    echo "- ✅ API文档生成功能正常"
else
    echo "- ❌ 需要检查API文档生成配置"
fi)

$(if [ -f docs/api/swagger/index.html ]; then
    echo "- ✅ Swagger UI生成功能正常"
else
    echo "- ❌ 需要检查Swagger UI配置"
fi)

EOF

    print_success "测试报告已生成: $report_file"
}

# 主函数
main() {
    print_info "开始API文档生成测试..."
    
    local exit_code=0
    
    # 运行各项测试
    test_environment || exit_code=1
    test_api_generation || exit_code=1
    test_swagger_ui || exit_code=1
    test_examples || exit_code=1
    test_docs_server || exit_code=1
    validate_openapi || exit_code=1
    
    # 生成报告
    generate_report
    
    # 输出结果
    if [ $exit_code -eq 0 ]; then
        print_success "所有测试通过！"
    else
        print_error "部分测试失败，请检查上述错误信息"
    fi
    
    return $exit_code
}

# 运行主函数
main "$@" 