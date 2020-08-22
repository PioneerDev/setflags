#!/usr/bin/env bash

# 获取正在运行的目标容器ID
container_id=`docker ps | grep "set-flags" | awk '{print $1 }'`

echo "set-flags容器ID: $container_id"

# 判断容器是否在运行
if [ ! -n "$container_id" ]; then
  echo "当前没有目标容器在运行"
else
  echo "停止当前容器" && docker rm -f $container_id
fi

# 获取正在运行的目标容器ID
container_id=`docker ps | grep "blaze" | awk '{print $1 }'`

echo "blaze容器ID: $container_id"

# 判断容器是否在运行
if [ ! -n "$container_id" ]; then
  echo "当前没有目标容器在运行"
else
  echo "停止当前容器" && docker rm -f $container_id
fi

# 获取正在运行的目标容器ID
container_id=`docker ps | grep "reminder" | awk '{print $1 }'`

echo "reminder容器ID: $container_id"

# 判断容器是否在运行
if [ ! -n "$container_id" ]; then
  echo "当前没有目标容器在运行"
else
  echo "停止当前容器" && docker rm -f $container_id
fi

docker build . -f APIDockerfile  -t set-flags:develop
docker run -d -v /home/ubuntu/setflags/secrets/config.yaml:/api/secrets/config.yaml -p 8080:8080 set-flags:develop

docker build . -f BlazeDockerfile  -t blaze:develop
docker run -d -v /home/ubuntu/setflags/secrets/config.yaml:/api/secrets/config.yaml blaze:develop

docker build . -f RemindDockerfile  -t reminder:develop
docker run -d -v /home/ubuntu/setflags/secrets/config.yaml:/api/secrets/config.yaml reminder:develop