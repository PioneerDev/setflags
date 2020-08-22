### 部署

#### 使用docker部署

```
目前共有三个服务需要启动
一个是API服务，一个是机器人websocket客户端，一个是定时任务
都依赖Postgres数据库

使用以下命令构建服务镜像, 然后运行即可
配置文件通过挂载方式传递进镜像内, config-example.yaml中说明了各个配置项的作用

也可以使用同级目录下的run.sh执行以下命令
```

1. API
> docker build . -f APIDockerfile  -t set-flags:develop
> docker run -d -v /home/ubuntu/setflags/secrets/config.yaml:/api/secrets/config.yaml -p 8081:8080 set-flags:develop

2. robot
> docker build . -f BlazeDockerfile  -t blaze:develop
> docker run -d -v /home/ubuntu/setflags/secrets/config.yaml:/api/secrets/config.yaml blaze:develop

3. cron job
> docker build . -f RemindDockerfile  -t reminder:develop
> docker run -it -v /home/ubuntu/setflags/secrets/config.yaml:/api/secrets/config.yaml reminder:develop
