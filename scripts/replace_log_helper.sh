#!/bin/bash

# 批量替换 log.NewHelper 为 log.NewEnhancedHelper
# 使用方法: ./scripts/replace_log_helper.sh

echo "开始替换 log.NewHelper 为 log.NewEnhancedHelper..."

# 替换所有 Go 文件中的 log.NewHelper
find . -name "*.go" -type f -exec sed -i '' 's/log\.NewHelper(/log.NewEnhancedHelper(/g' {} \;

echo "替换完成！"
echo ""
echo "请检查以下文件确保替换正确："
echo "1. internal/service/account.go"
echo "2. internal/biz/account.go" 
echo "3. internal/data/account.go"
echo "4. internal/task/cron.go"
echo "5. 其他使用 log.NewHelper 的文件"
echo ""
echo "然后运行 'go build ./...' 检查编译是否正常" 