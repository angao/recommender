# Recommender

**`Recommender` 通过 `Prometheus` 获取应用 `Container` 监控数据（`CPU`，`Memory`，`Disk IO`, `Network IO`），并按照 `Container` 聚合数据，最后持久化到 `MySQL` 数据库。并以 `RESTful` 的方式提供 `HTTP API`  接口。**

## 一、编译部署

1、编译

```bash
make build
```
在 `./bin/` 目录生成可执行二进制文件

2、部署

- 创建 MySQL 数据库，如 `recommender`，执行 `deploy/create_tables.sql` 文件创建数据库表
- 如果以二进制部署，直接将二进制 `scp` 到服务器执行即可
- 支持部署到 Kubernetes 集群

```bash
# 打包成 image
make docker-build
# 创建 Deployment
kubectl create -f deploy/recommender.deploy.yaml
```

## 二、配置

`recommender` 通过配置文件的方式支持可配置，通过 `recommender --config-file="/etc/config.yaml"` 方式配置，详细配置如下：

> 默认使用 `MySQL` 数据库

```yaml
databaseConfig:
  # 数据库用户名
  username: "root"
  # 数据库密码
  password: "123456"
  # 数据库地址
  url: "192.168.99.100"
  # 数据库端口
  port: 30006
  # 数据库名
  name: "recommender"
  # 连接池，默认为 2
  maxIdleConns: 2
  # 最大连接数，默认 10
  maxOpenConns: 10
prometheusConfig:
  # Prometheus 服务地址
  address: "http://192.168.19.0:32100"
extraConfig:
  # Prometheus 默认查询历史时长，默认 90d
  history: "90d"
  # 对外 HTTP API 端口，默认 9098
  apiPort: 9098
```

## 三、`API` 接口

1、创建应用
```
method: POST
url: /api/v1/application
param: 
{
    "name: "test"
}

return 
{
    "code": 200,
    "message": "success"
}
```
2、获取应用
```
获取指定名称应用:
method: GET
url: /api/v1/application/:name

return 
{
    "code": 200,
    "data": {
        "id": 162,
        "name": "test",
        "created": "2018-10-15T14:02:25+08:00",
        "updated": "2018-10-15T14:02:25+08:00",
        "deleted": "0001-01-01T00:00:00Z"
    },
    "message": "success"
}

3、获取全部应用:
method: GET
url: /api/v1/applications

return 
{
    "code": 200,
    "data": [
        {
            "id": 162,
            "name": "web",
            "created": "2018-10-15T14:02:25+08:00",
            "updated": "2018-10-15T14:02:25+08:00",
            "deleted": "0001-01-01T00:00:00Z"
        }
    ],
    "message": "success"
}
```
4、删除应用
```
method: DELETE
url: /api/v1/application/:name

return 
{
    "code": 200,
    "message": "success"
}
```
5、获取指定应用的资源推荐值
```
method: GET
url: /api/v1/resource/:name

return 
{
    "code": 200,
    "data": {
        "id": 162,
        "name": "web",
        "container_resource": [
            {
                "id": 23,
                "name": "test",
                "application_id": 162,
                "timeframe_id": 0,
                "cpu_limit": 909,
                "memory_limit": 844,
                "disk_read_io_limit": 923,
                "disk_write_io_limit": 978,
                "network_receive_io_limit": 959,
                "network_transmit_io_limit": 997,
                "created": "2018-10-16T10:25:55+08:00",
                "updated": "2018-10-16T10:30:15+08:00"
            }
        ]
    },
    "message": "success"
}
```
6、获取全部应用的资源推荐
```
method: GET
url: /api/v1/resources

return 
{
    "code": 200,
    "data": [
        {
            "id": 162,
            "name": "app-3103947192",
            "container_resource": [
                {
                    "id": 23,
                    "name": "test",
                    "application_id": 162,
                    "timeframe_id": 0,
                    "cpu_limit": 909,
                    "memory_limit": 844,
                    "disk_read_io_limit": 923,
                    "disk_write_io_limit": 978,
                    "network_receive_io_limit": 959,
                    "network_transmit_io_limit": 997,
                    "created": "2018-10-16T10:25:55+08:00",
                    "updated": "2018-10-16T10:30:15+08:00"
                }
            ]
        }
    ],
    "message": "success"
}
```
7、获取指定时间段指定应用资源推荐
```
method: GET
url: /api/v1/resources/timeframe/:name/:appName
name: 指定时间段名称
appName: 指定应用名称

return
{
    "code": 200,
    "data": {
        "id": 162,
        "name": "web",
        "container_resource": [
            {
                "id": 23,
                "name": "test",
                "application_id": 162,
                "timeframe_id": 0,
                "cpu_limit": 909,
                "memory_limit": 844,
                "disk_read_io_limit": 923,
                "disk_write_io_limit": 978,
                "network_receive_io_limit": 959,
                "network_transmit_io_limit": 997,
                "created": "2018-10-16T10:25:55+08:00",
                "updated": "2018-10-16T10:30:15+08:00"
            }
        ]
    },
    "message": "success"
}
```
8、删除指定应用的推荐值
```
method: DELETE
url: /api/v1/resource/:name

return
{
    "code": 200,
    "message": "success"
}
```
9、删除指定时间段应用的推荐值
```
method: DELETE
url: /api/v1/resources/timeframe/:name

return
{
    "code": 200,
    "message": "success"
}
```
10、创建指定时间段
```
method: POST
url: /api/v1/timeframe
param:
{
    "name": "double11",
    "start": "2017-11-10 23:00:00", // 开始时间
    "end": "2017-11-11 01:00:00",  // 结束时间
    "status": "on", // status 有 on 和 off 两个状态，只有 on 的状态才会从 Prometheus 拉取数据计算，计算完成后将 on 更新成 off
    "description":"" //  描述
}
return
{
    "code": 200,
    "message": "success"
}
```
11、获取全部指定时间段
```
method: GET
url: /api/v1/timeframes

return
{
    "code": 200,
    "data": [
        {
            "id": 6,
            "name": "double11",
            "start": "2017-11-10T23:00:00+08:00",
            "end": "2017-11-11T01:00:00+08:00",
            "status": "off",
            "description": "",
            "created": "2018-10-15T14:35:24+08:00",
            "updated": "2018-10-16T10:28:04+08:00",
            "deleted": "0001-01-01T00:00:00Z"
        }
    ],
    "message": "success"
}
```
12、获取指定时间段
```
method: GET
url: /api/v1/timeframe/:name

return
{
    "code": 200,
    "data": {
        "id": 6,
        "name": "double11",
        "start": "2017-11-10T23:00:00+08:00",
        "end": "2017-11-11T01:00:00+08:00",
        "status": "off",
        "description": "",
        "created": "2018-10-15T14:35:24+08:00",
        "updated": "2018-10-16T10:28:04+08:00",
        "deleted": "0001-01-01T00:00:00Z"
    },
    "message": "success"
}
```
13、更新指定时间段
```
method: PUT
url: /api/v1/timeframe

param
{
    "id": 6,
    "name": "double11",
    "start": "2017-11-10 23:00:00", 
    "end": "2017-11-11 01:00:00",  
    "status": "on",
    "description":""
}
return 
{
    "code": 200,
    "message": "success"
}
```
14、删除指定时间段
```
method: DELETE
url: /api/v1/timeframe/:name

return 
{
    "code": 200,
    "message": "success"
}
```