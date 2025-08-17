#!/bin/bash

# 检查服务是否已在运行
if [ -f .pid ]; then
  echo "服务已在运行中，PID: $(cat .pid)。如果需要重启，请先运行 stop.sh。"
  exit 1
fi

echo "正在启动服务..."
# 在后台运行 Go 程序，并将输出重定向到日志文件
nohup go run . > server.log 2>&1 &

# 获取后台进程的 PID 并保存
PID=$!
echo $PID > .pid

echo "服务已成功启动，PID: $PID"
echo "日志文件位于: server.log"
