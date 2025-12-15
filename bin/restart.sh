#!/bin/bash

# IPTV服务重启脚本

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "正在重启IPTV服务..."
echo ""

# 先停止服务
"$SCRIPT_DIR/stop.sh"

# 等待一下，确保进程完全停止
sleep 2

# 启动服务
"$SCRIPT_DIR/start.sh"

