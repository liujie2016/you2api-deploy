#!/bin/bash

# 检查 .pid 文件是否存在
if [ ! -f .pid ]; then
  echo "服务未在运行 (找不到 .pid 文件)。"
  exit 1
fi

# 从文件中读取 PID
PID=$(cat .pid)

echo "正在停止服务，PID: $PID"
# 杀死进程
kill $PID

# 检查进程是否已成功停止
if [ $? -eq 0 ]; then
  echo "服务已成功停止。"
  # 删除 .pid 文件
  rm .pid
else
  echo "停止服务失败。请手动检查进程: ps aux | grep 'go run'"
fi
