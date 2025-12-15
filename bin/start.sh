#!/bin/bash

# IPTV服务启动脚本

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY_NAME="iptv"
PID_FILE="$PROJECT_DIR/iptv.pid"
LOG_DIR="$PROJECT_DIR/logs"

# 切换到项目目录
cd "$PROJECT_DIR" || exit 1

# 检查程序是否已经在运行
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "IPTV服务已经在运行中 (PID: $PID)"
        exit 1
    else
        # PID文件存在但进程不存在，删除旧的PID文件
        rm -f "$PID_FILE"
    fi
fi

# 检查二进制文件是否存在
if [ ! -f "$PROJECT_DIR/$BINARY_NAME" ]; then
    echo "错误: 找不到可执行文件 $BINARY_NAME"
    echo "请先编译程序: go build -o iptv ."
    exit 1
fi

# 确保日志目录存在
mkdir -p "$LOG_DIR"

# 启动程序（后台运行）
echo "正在启动IPTV服务..."
nohup "$PROJECT_DIR/$BINARY_NAME" > /dev/null 2>&1 &
PID=$!

# 保存PID到文件
echo $PID > "$PID_FILE"

# 等待一下，检查进程是否成功启动
sleep 1
if ps -p "$PID" > /dev/null 2>&1; then
    echo "IPTV服务启动成功 (PID: $PID)"
    echo "日志文件: $LOG_DIR/app-$(date +%Y-%m-%d).log"
else
    echo "IPTV服务启动失败"
    rm -f "$PID_FILE"
    exit 1
fi

