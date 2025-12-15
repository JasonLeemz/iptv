#!/bin/bash

# IPTV服务停止脚本

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY_NAME="iptv"
PID_FILE="$PROJECT_DIR/iptv.pid"

# 切换到项目目录
cd "$PROJECT_DIR" || exit 1

# 检查PID文件是否存在
if [ ! -f "$PID_FILE" ]; then
    echo "IPTV服务未运行（未找到PID文件）"
    exit 0
fi

# 读取PID
PID=$(cat "$PID_FILE")

# 检查进程是否存在
if ! ps -p "$PID" > /dev/null 2>&1; then
    echo "IPTV服务未运行（进程不存在）"
    rm -f "$PID_FILE"
    exit 0
fi

# 停止进程
echo "正在停止IPTV服务 (PID: $PID)..."
kill "$PID" 2>/dev/null

# 等待进程结束（最多等待10秒）
for i in {1..10}; do
    if ! ps -p "$PID" > /dev/null 2>&1; then
        break
    fi
    sleep 1
done

# 如果进程仍在运行，强制杀死
if ps -p "$PID" > /dev/null 2>&1; then
    echo "进程未正常退出，强制终止..."
    kill -9 "$PID" 2>/dev/null
    sleep 1
fi

# 删除PID文件
rm -f "$PID_FILE"

# 确认进程已停止
if ps -p "$PID" > /dev/null 2>&1; then
    echo "停止IPTV服务失败"
    exit 1
else
    echo "IPTV服务已停止"
fi

