#!/bin/zsh

# 文件目录准备
dir=$(grep "^dir" "./redis.conf" | awk '{print $2}' | sed 's/"//g')
user=$(whoami)

# 检查文件夹是否存在
if [ ! -d "$dir" ]; then
    # 创建文件夹
    sudo mkdir -p "$dir"
    echo "文件夹已创建，并为本用户设置写入权限。"
else

    echo "文件夹已存在。"
fi

# 设置文件夹权限为本用户可写
sudo chown $user:$user "$dir"
sudo chmod +w "$dir"

# 启动主节点
redis-server redis.conf &
# 等待一段时间启动主节点
sleep 2
# 启动从节点1
redis-server redis6380.conf &
# 启动从节点2
redis-server redis6381.conf &

#启动哨兵
redis-sentinel redis-sentinel1.conf &
redis-sentinel redis-sentinel2.conf &

