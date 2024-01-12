#!/bin/zsh

# 启动主节点
redis-server redis.conf &

# 等待一段时间启动主节点
sleep 2

# 启动从节点1
redis-server redis6380.conf

# 启动从节点2
redis-server redis6381.conf
