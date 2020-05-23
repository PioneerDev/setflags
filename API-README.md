## set-flags后端api

### 数据库 postgres
```shell script
docker run --name docker-postgres -e POSTGRES_PASSWORD=123456 -p 5432:5432 -d postgres:10.3

docker exec -it a4bd54d0e28b psql -U postgres -d postgres -h localhost -p 5432
```


### web框架 gin
```shell script
go run main.go
```
> 项目启动后会执行数据库迁移

### 项目概览

#### 配置文件
> conf/app.ini

#### 模型
> models/

#### 路由文件
> routers/