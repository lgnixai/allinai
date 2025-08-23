#!/bin/bash

# 聊天历史测试数据生成和执行脚本
# 用于快速生成测试数据并验证

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_NAME=${DB_NAME:-"one_api"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  聊天历史测试数据生成脚本${NC}"
echo -e "${BLUE}================================${NC}"

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}检查依赖...${NC}"
    
    if ! command -v mysql &> /dev/null; then
        echo -e "${RED}错误: 未找到 mysql 客户端${NC}"
        echo "请安装 MySQL 客户端或确保其在 PATH 中"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 依赖检查通过${NC}"
}

# 检查数据库连接
check_database_connection() {
    echo -e "${YELLOW}检查数据库连接...${NC}"
    
    if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT 1;" &> /dev/null; then
        echo -e "${RED}错误: 无法连接到数据库${NC}"
        echo "请检查以下配置："
        echo "  - 数据库主机: $DB_HOST"
        echo "  - 数据库端口: $DB_PORT"
        echo "  - 数据库用户: $DB_USER"
        echo "  - 数据库名称: $DB_NAME"
        echo "  - 数据库密码: 请设置 DB_PASSWORD 环境变量"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 数据库连接成功${NC}"
}

# 检查表是否存在
check_table_exists() {
    echo -e "${YELLOW}检查 chat_histories 表...${NC}"
    
    if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "DESCRIBE chat_histories;" &> /dev/null; then
        echo -e "${RED}错误: chat_histories 表不存在${NC}"
        echo "请先执行数据库迁移脚本: bin/migration_v0.6-v0.7.sql"
        exit 1
    fi
    
    echo -e "${GREEN}✓ chat_histories 表存在${NC}"
}

# 清理现有测试数据
clean_test_data() {
    echo -e "${YELLOW}清理现有测试数据...${NC}"
    
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
        DELETE FROM chat_histories WHERE user_id IN (1, 2, 3);
    " 2>/dev/null || true
    
    echo -e "${GREEN}✓ 测试数据清理完成${NC}"
}

# 生成测试数据
generate_test_data() {
    echo -e "${YELLOW}生成测试数据...${NC}"
    
    if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$SCRIPT_DIR/generate_chat_history_test_data.sql"; then
        echo -e "${RED}错误: 生成测试数据失败${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 测试数据生成完成${NC}"
}

# 验证测试数据
validate_test_data() {
    echo -e "${YELLOW}验证测试数据...${NC}"
    
    if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$SCRIPT_DIR/test_chat_history_data.sql"; then
        echo -e "${RED}错误: 验证测试数据失败${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 测试数据验证完成${NC}"
}

# 显示测试数据统计
show_statistics() {
    echo -e "${YELLOW}测试数据统计:${NC}"
    
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
        SELECT 
            '数据统计' as info,
            COUNT(*) as total_records,
            COUNT(DISTINCT user_id) as total_users,
            COUNT(DISTINCT session_id) as total_sessions,
            COUNT(DISTINCT model) as total_models,
            SUM(tokens) as total_tokens,
            ROUND(SUM(cost), 6) as total_cost
        FROM chat_histories 
        WHERE user_id IN (1, 2, 3);
    " 2>/dev/null || echo -e "${RED}无法获取统计数据${NC}"
}

# 显示使用说明
show_usage() {
    echo -e "${BLUE}使用方法:${NC}"
    echo "  $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示此帮助信息"
    echo "  -c, --clean             仅清理测试数据"
    echo "  -g, --generate          仅生成测试数据"
    echo "  -v, --validate          仅验证测试数据"
    echo "  -f, --full              完整流程（清理+生成+验证）"
    echo ""
    echo "环境变量:"
    echo "  DB_HOST                 数据库主机 (默认: localhost)"
    echo "  DB_PORT                 数据库端口 (默认: 3306)"
    echo "  DB_USER                 数据库用户 (默认: root)"
    echo "  DB_NAME                 数据库名称 (默认: one_api)"
    echo "  DB_PASSWORD             数据库密码 (必需)"
    echo ""
    echo "示例:"
    echo "  DB_PASSWORD=mypassword $0 --full"
    echo "  DB_PASSWORD=mypassword $0 --generate"
}

# 主函数
main() {
    case "${1:-}" in
        -h|--help)
            show_usage
            exit 0
            ;;
        -c|--clean)
            check_dependencies
            check_database_connection
            check_table_exists
            clean_test_data
            ;;
        -g|--generate)
            check_dependencies
            check_database_connection
            check_table_exists
            generate_test_data
            show_statistics
            ;;
        -v|--validate)
            check_dependencies
            check_database_connection
            check_table_exists
            validate_test_data
            ;;
        -f|--full|"")
            check_dependencies
            check_database_connection
            check_table_exists
            clean_test_data
            generate_test_data
            validate_test_data
            show_statistics
            ;;
        *)
            echo -e "${RED}错误: 未知选项 $1${NC}"
            show_usage
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}  操作完成！${NC}"
    echo -e "${GREEN}================================${NC}"
}

# 检查是否提供了数据库密码
if [ -z "$DB_PASSWORD" ]; then
    echo -e "${RED}错误: 请设置 DB_PASSWORD 环境变量${NC}"
    echo "示例: DB_PASSWORD=mypassword $0"
    exit 1
fi

# 执行主函数
main "$@"

