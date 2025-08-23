#!/bin/bash

# 聊天历史测试数据生成和执行脚本 (SQLite版本)
# 用于快速生成测试数据并验证

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
DB_FILE=${DB_FILE:-"one-api.db"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}  聊天历史测试数据生成脚本 (SQLite)${NC}"
echo -e "${BLUE}================================${NC}"

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}检查依赖...${NC}"
    
    if ! command -v sqlite3 &> /dev/null; then
        echo -e "${RED}错误: 未找到 sqlite3 客户端${NC}"
        echo "请安装 SQLite3 或确保其在 PATH 中"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 依赖检查通过${NC}"
}

# 检查数据库文件是否存在
check_database_file() {
    echo -e "${YELLOW}检查数据库文件...${NC}"
    
    if [ ! -f "$DB_FILE" ]; then
        echo -e "${RED}错误: 数据库文件 $DB_FILE 不存在${NC}"
        echo "请确保数据库文件存在或设置正确的 DB_FILE 环境变量"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 数据库文件存在: $DB_FILE${NC}"
}

# 检查表是否存在
check_table_exists() {
    echo -e "${YELLOW}检查 chat_histories 表...${NC}"
    
    if ! sqlite3 "$DB_FILE" "SELECT name FROM sqlite_master WHERE type='table' AND name='chat_histories';" | grep -q "chat_histories"; then
        echo -e "${RED}错误: chat_histories 表不存在${NC}"
        echo "请先执行数据库迁移脚本: bin/migration_v0.6-v0.7.sql"
        exit 1
    fi
    
    echo -e "${GREEN}✓ chat_histories 表存在${NC}"
}

# 清理现有测试数据
clean_test_data() {
    echo -e "${YELLOW}清理现有测试数据...${NC}"
    
    sqlite3 "$DB_FILE" "DELETE FROM chat_histories WHERE user_id IN (1, 2, 3);" 2>/dev/null || true
    
    echo -e "${GREEN}✓ 测试数据清理完成${NC}"
}

# 生成测试数据
generate_test_data() {
    echo -e "${YELLOW}生成测试数据...${NC}"
    
    if ! sqlite3 "$DB_FILE" < "$SCRIPT_DIR/generate_chat_history_test_data_sqlite.sql"; then
        echo -e "${RED}错误: 生成测试数据失败${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 测试数据生成完成${NC}"
}

# 验证测试数据
validate_test_data() {
    echo -e "${YELLOW}验证测试数据...${NC}"
    
    if ! sqlite3 "$DB_FILE" < "$SCRIPT_DIR/test_chat_history_data_sqlite.sql"; then
        echo -e "${RED}错误: 验证测试数据失败${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 测试数据验证完成${NC}"
}

# 显示测试数据统计
show_statistics() {
    echo -e "${YELLOW}测试数据统计:${NC}"
    
    sqlite3 "$DB_FILE" "
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
    echo "  DB_FILE                 数据库文件路径 (默认: one-api.db)"
    echo ""
    echo "示例:"
    echo "  $0 --full"
    echo "  DB_FILE=/path/to/database.db $0 --generate"
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
            check_database_file
            check_table_exists
            clean_test_data
            ;;
        -g|--generate)
            check_dependencies
            check_database_file
            check_table_exists
            generate_test_data
            show_statistics
            ;;
        -v|--validate)
            check_dependencies
            check_database_file
            check_table_exists
            validate_test_data
            ;;
        -f|--full|"")
            check_dependencies
            check_database_file
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

# 执行主函数
main "$@"

