#1. 项目编译运行
## 1.1 项目编译
* 本项目使用go 1.16 编译
```
1. clone 项目
git clone https://github.com/pili-engineering/QNLiveKit_Server 

2. 编译
cd QNLiveKit_Server/app/live
go build
```
* 在当前目录下，生成可执行文件  live


## 1.2 项目运行
* 项目的运行，需要提供一个配置文件  config.yaml。配置文件的说明，见后面。
```
./live -f config.yaml 2>&1 > live.log &
```

## 1.3 使用docker-compose部署
* 安装 docker、docker-compose
* 新建 .env 文件 ，内容示例如下
```
PORT_LIVE=8099
PORT_PROMETHEUS=29090
PORT_GRAFANA=23000
QINIU_ACCESS_KEY = {{ ak }}
QINIU_SECRET_KEY = {{ sk }}
IM_ADMIN_TOKEN = {{ im_admin_token }}
PILI_PUB_KEY = {{ publish_key }}
```

| 字段名    | 说明                                                                                                                                                           |   |
|-----|--------------------------------------------------------------------------------------------------------------------------------------------------------------|-----|
| PORT_LIVE  | qnlivekit_server-qnliv中， "${PORT_LIVE}:8099" 将容器的8099端口映射到宿主机PORT_LIVE端口                                                                                     |     |
|  PORT_PROMETHEUS   | prometheus    "${PORT_PROMETHEUS}:9090"  将容器的9090端口映射到宿主机PORT_PROMETHEUS端口                                                                                   |     |
| PORT_GRAFANA    | grafana       "${PORT_GRAFANA}:3000"   将容器的3000端口映射到宿主机PORT_GRAFANA端口                                                                                        |     |
| QINIU_ACCESS_KEY    | ak/sk鉴权    {{ .qiniu_key_pair_access_key }}                                                                                                                  |     |
| QINIU_SECRET_KEY    | ak/sk鉴权    {{ .qiniu_key_pair_secret_key }}                                                                                                                  |     |
| IM_ADMIN_TOKEN    | IM服务 获取管理员token ，详情参考[IM服务接入指南](https://developer.qiniu.com/IM/8332/startim) 与 [IM服务端接入指南](https://developer.qiniu.com/IM/8102/im-service-integration-guide) |     |
| PILI_PUB_KEY    | Pili服务直播推流鉴权，当采用限时鉴权，security_type为expiry时， 从【直播空间设置】获取publish_key ；当security_type为无鉴权（none）和 限时鉴权SK（expiry_sk），该位置留空                                        |     |


* 执行  `docker-compose --env-file live.env up`
* 构建（重新构建）项目中的服务容器 `docker-compose build`
 


#2. 配置文件
* 互动直播使用 yaml 格式的配置文件，文件内容如下所示：

```yaml
node:
  node_id: {{.node_id}}

account:
  access_key: {{ .qiniu_key_pair_access_key }}
  secret_key: {{ .qiniu_key_pair_secret_key }}

httpq:
  addr: ":8099"

auth:
  jwt_key: jwtkey
  server:
    enable: true
    access_key: ak
    secret_key: sk


cache:
  type: {{type}}
  addr: {{host:ip}}

mysql:
  databases:
    - host: 127.0.0.1
      port: 3306
      username: {{ user }}
      password: {{ password }}
      max_idle_conns: 10
      max_open_conns: 100
      database: live
      conn_max_life_time: 5
      default: live


    - host: 127.0.0.1
      port: 3306
      username: {{ user }}
      password: {{ password }}
      max_idle_conns: 10
      max_open_conns: 100
      database: live
      conn_max_life_time: 5
      default: live
      read_only: true



im:
  app_id: {{ im_app_id }}
  endpoint: {{ im_endpoint }}
  token: {{ im_admin_token }}

rtc:
  app_id: {{ app_id }}

pili:
  hub: {{ 直播hub }}
  room_token_expire_s: 3600
  playback_url: https://{{xxx}}-playback.com
  security_type: { { security_type } }
  publish_key: {{ publish_key }}
  publish_expire_s: 3600
  stream_pattern: qn_live_kit-%s
  publish_url:  rtmp://{{xxx}}-publish.com
  publish_domain: {{xxx}}-publish.com
  rtmp_play_url: rtmp://{{xxx}}-rtmp.com
  flv_play_url: https://{{xxx}}-hdl.com
  hls_play_url: https://{{xxx}}-hls.com
  media_url: https://{{xxx}}

cron:
  single_task_node: 1

prome:
  client_mode: pusher
  exporter_config:
    listen_addr: ":9200"
  pusher_config:
    url: "https://{{xxx}}"
    job: live
    instance: live_{{.node}}
    interval_s: 10

callback:
  addr: https://niucube-api.qiniu.com/v1/live/callback

kodo:
  callback: https://{{xxx}}
  bucket: {{ censor_bucket }}
  addr: https://{{xxx}}

gift:
  gift_addr: http://127.0.0.1:8099/manager/gift/test


```


## 配置文件详细介绍

<br>

###  node配置
```
node_id: {{ nodeId }}
```

| 字段名    | 类型  | 示例  | 说明                                                                   |
|-----|-----|-----|----------------------------------------------------------------------|
| node_id    | int | 1   | 互动直播支持分布式部署，每个节点指定一个nodeId。所有节点的nodeId 不能重复。<br/>nodeId 会用于内部的ID 生成。 |


<br>

###  面向服务的鉴权配置 account：

```yaml
account:
  access_key: {{ .qiniu_key_pair_access_key }}
  secret_key: {{ .qiniu_key_pair_secret_key }}
```

| 字段名        | 示例                               | 说明                                                                                      |
|------------|----------------------------------|-----------------------------------------------------------------------------------------|
| access_key | {{ .qiniu_key_pair_access_key }} | 七牛对API请求鉴权，用户需要使用密钥（AccessKey/SecretKey）进行身份验证。此处填入AccessKey                            |
| secret_key | {{ .qiniu_key_pair_secret_key }} | 此处填入SecretKey，第一次使用七牛API之前，您需要在七牛开发者平台的 [密钥管理](https://portal.qiniu.com/user/key) 中取得密钥 |

<br>

###  监听地址
```yaml
httpq:
  addr: ":8099"
```

| 字段名    | 类型     | 示例  | 说明                                                                                 |
|-----|--------|-----|------------------------------------------------------------------------------------|
| addr    | string | 1   | 互动直播服务提供HTTP RESTful API 接口。service 配置服务的监听地址`{host:ip}`<br/>host：监听的IP<br/>port: 监听的端口 |


<br>

### 互动直播服务端鉴权
```yaml
auth:
  jwt_key: jwtkey
  server:
    enable: true
    access_key: ak
    secret_key: sk

```

| 字段名                | 类型     | 示例      | 说明                                                                               |
|--------------------|--------|---------|----------------------------------------------------------------------------------|
| jwt_key            | string | jwt-key | 用于客户端鉴权token 的加解密                                                                |
| server:<br/>&nbsp;&nbsp;&nbsp;enable| bool   | true    | enable：true，开启鉴权；false，关闭鉴权。如果互动直播服务，与需要调用互动直播服务接口的其他服务，在同一个安全局域网内，可以关闭服务端接口的鉴权。 |
|server:<br/>&nbsp;&nbsp;&nbsp;access_key        | string | ak      | 互动直播服务，提供面向服务端的接口，与面向客户端的接口。面向服务端的接口，使用 ak/sk 的鉴权方式。这里为ak                        |
| server:<br/>&nbsp;&nbsp;&nbsp;secret_key       | string   | sk      | secret key                                                                       |


<br>

### Redis配置
```yaml
cache_config:
    type: node
    addr: {{ip:port}}
```
或者
```yaml
cache_config:
    type: cluster
    addrs:
    	- {{ip:port}}
    	- {{ip:port}}
```


| 字段名   | 类型     | 示例           |  说明   |
|-------|--------|--------------|-----|
| type  | string | node/cluster | 低代码服务，使用redis 作为缓存。支持redis 单节点模式，与集群模式。<br/>node：单节点模式 <br/>cluster：集群模式   |
| addr  | string | addr: 127.0.0.1:6379             | 配置type为node时，必须配置单节点redis 的地址。    |
| addrs |        | 如上所示         | 配置type为cluster时，必须配置的redis集群地址列表。    |

<br>

### Mysql配置
```yaml
mysql:
  databases:
    - host: 127.0.0.1
      port: 3306
      username: {{ user }}
      password: {{ password }}
      max_idle_conns: 10
      max_open_conns: 100
      database: live
      conn_max_life_time: 5
      default: live
    - host: 127.0.0.1
      port: 3306
      username: {{ user }}
      password: {{ password }}
      max_idle_conns: 10
      max_open_conns: 100
      database: live
      conn_max_life_time: 5       
      default: live
      read_only: true
```
* 互动直播服务使用mysql 进行业务数据存储。
* 支持配置多个数据库，使用读写分离模式。


| 字段名                | 说明                                                                                                                                                                  |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| host/port          |   MySQL 实例的 IP 地址/主机名和可选端口。                                                                                                  |
| username           |  数据库用户的登录名/用户名  |
| password           |   数据库用户密码                                                                                                                                                                  |
| max_idle_conns     |      空闲连接池中的最大连接数，默认2（Grafana v5.4+）                                                                                                                                                              |
| max_open_conns     |   数据库的最大打开连接数，默认unlimited（Grafana v5.4+）                                                                                                                                                                 |
| database           |  MySQL 数据库的名称。                                                                                                                                                                  |
| conn_max_life_time | 可以重用连接的最长时间（以秒为单位），默认为14400/4 小时。这应该始终低于 MySQL (Grafana v5.4+) 中配置的wait_timeout。                                                                                                                                                                   |
| default            |  默认数据源意味着它将被预先选用                                                                                                                                                                |


<br>

### IM配置

```yaml
im:
  app_id: {{ im_app_id }}
  endpoint: {{ im_endpoint }}
  token: {{ im_admin_token }}
```

| 字段名    | 类型     | 示例  | 说明                                                                                                                                                                  |     |
|-----|--------|-----|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----|
| app_id    | string |     | 互动直播服务，使用七牛IM 作为即时通信消息服务。<br/>app_id：创建IM App时，为App生成的唯一标识                                                                                                          |     |
| endpoint    | string       |     | endpoint：App所在API服务的地址  ，参考 [IM服务接入指南](https://developer.qiniu.com/IM/8332/startim) 与 [IM服务端接入指南](https://developer.qiniu.com/IM/8102/im-service-integration-guide) |     |
| token   |   string     |     | token：IM中获取的管理员token，参考 [IM服务接入指南](https://developer.qiniu.com/IM/8332/startim) 与 [IM服务端接入指南](https://developer.qiniu.com/IM/8102/im-service-integration-guide)     |     |

<br>

### RTC PILI配置
RTC
```yaml
rtc:
  app_id: {{ app_id }}
```
| 字段名    | 类型     | 示例           | 说明                                                                                                                                                      |
|-----|--------|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| app_id    | string | c9n9509o5e70 | 互动直播服务，使用七牛的RTC 服务，作为连麦服务。<br/>RTC 应用的配置，参考：[管理实时音视频](https://developer.qiniu.com/rtc/9858/applist)|



Pili
```yaml
pili:
  hub: {{ 直播hub }}
  room_token_expire_s: 3600
  playback_url: https://{{xxx}}-playback.com
  security_type: { { security_type } }
  publish_key: {{ publish_key }}
  publish_expire_s: 3600
  stream_pattern: qn_live_kit-%s
  publish_url:  rtmp://{{xxx}}-publish.com
  publish_domain: {{xxx}}-publish.com
  rtmp_play_url: rtmp://{{xxx}}-rtmp.com
  flv_play_url: https://{{xxx}}-hdl.com
  hls_play_url: https://{{xxx}}-hls.com
  media_url: https://{{xxx}}
```

| 字段名    | 类型     | 示例            | 说明                                              |
|-----|--------|---------------|-------------------------------------------------|
| hub    | string | live          | 直播空间信息，参考：[直播云](https://developer.qiniu.com/pili) |
|  room_token_expire_s      | int    | 3600          |                                                 |
|  playback_url      | string |               | 点播域名                                            |
|    security_type    | string | none/expiry/expiry_sk | 推流鉴权方式,详情见下表                                    |
|    publish_key    | string |               |    详情见下表                                               |
|    publish_expire_s               | int    | 3600          |                                                 |
|    stream_pattern               |   string     | qn_live_kit-%s |                                                 |
|     publish_url              |  string      |               | 推流域名                                            |
|      publish_domain                        |  string      |               | 推流域名                                            |
|    rtmp_play_url                          |   string     |               | RTMP播放域名                                        |
|    flv_play_url             |  string      |               | FLV播放域名                                         |
|     hls_play_url          |  string      |               | HLS播放域名                                         |
|      media_url               |  string      |               |                                                 |


* 直播推流鉴权说明：目前直播推流鉴权，支持三种方式：无鉴权，限时鉴权，限时鉴权SK，不同鉴权模式的配置方式如下

| 鉴权模式   | security_type | publish_key                | publish_expire_s                |
|--------|---------------|----------------------------|---------------------------------|
| 无鉴权    | none          | 无需指定，留空                    | 过期时间秒。如：3600 表示 一小时后过期，推流URL 过期 |
| 限时鉴权   | expiry        | 使用配置的key 鉴权。从【直播空间设置】获取key | 同上                              |
| 限时鉴权SK | expiry_sk     | 使用RTC 用户的SK 鉴权。无需配置，留空。    | 同上                              |

<br>

### cron配置

```yaml
cron:
  single_task_node: 1
```

| 字段名    | 类型  | 示例  | 说明                                                                                                                                                      |
|-----|-----|-----|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| single_task_node    | int | 1   |单任务运行节点的ID。有些任务，需要全局单线程执行，只有节点ID 等于 single_task_node 的节点，才会运行单线程的任务。|

<br>

###  Prometheus系统监控配置
```yaml
prome_config:
  client_mode: exporter
  exporter_config:
    listen_addr: ":9200"
```
或
```yaml
prome_config:
  client_mode: pusher
  pusher_config:
    url: "https://{{xxx}}"
    job: live
    instance: live_{{.node}}
    interval_s: 10
```

| 字段名                                                                                                          | 类型     | 示例              | 说明                                                                                                                                   |
|--------------------------------------------------------------------------------------------------------------|--------|-----------------|--------------------------------------------------------------------------------------------------------------------------------------|
| client_mode                                                                                                  | string | exporter/pusher | client_mode：指标导出模式配置，支持的模式如下<br/>  exporter：exporter 模式，开启http 监听，由prometheus 服务来拉取。<br/>pusher：pusher模式，主动向prometheus-gateway 推送指标。 |
| exporter_config <br/> &nbsp;&nbsp;listen_addr                                                                | string | ":9200"         | 在client_mode 为 exporter 时，需要配置。<br/> listen_addr： 监听的端口                                                                              |
| pusher_config<br/>  &nbsp;&nbsp;url<br/>&nbsp;&nbsp;job<br/>&nbsp;&nbsp;instance <br/>&nbsp;&nbsp;interval_s | string |                 | 在client_mode 为 pusher_config 时，需要配置。url：prometheus-gateway 的指标收集地址。<br/> job：任务名称 <br/>instance：live 服务的实例ID。 <br/>interval_s：推送间隔。单位：秒。                                 |
|

<br>

### callback
```yaml
callback:
  addr: https://niucube-api.qiniu.com/v1/live/callback
```

| 字段名    | 类型     | 示例  | 说明                                            |
|-----|--------|-----|-----------------------------------------------|
| addr    | string | https://niucube-api.qiniu.com/v1/live/callback   | 配置低代码服务的回调地址,该地址由业务服务来实现.<br/> 低代码服务会将状态信息，通过回调的方式，通知给业务服务 |

<br>

### kodo配置（三鉴模块）- censor
```yaml
kodo:
  callback: https://{{xxx}}
  bucket: {{ censor_bucket }}
  addr: https://{{xxx}}
```

| 字段名    | 类型     | 示例   | 说明                                                                                                       |
|-----|--------|------|----------------------------------------------------------------------------------------------------------|
| callback    | string |      | AI审核结果回调地址，配置为项目的域名。                                                                                     |
| bucket    | string       |  | AI审核疑似违规照片的存储bucket。                                                                                     |
| addr   |   string     |      | bucket内存储文件的外链域名。<br/>存储的bucket，与外链域名配置，参考 [对象存储使用](https://developer.qiniu.com/kodo/8452/kodo-console-guide) |

<br>

### 礼物模块配置 - gift
```yaml
gift:
  gift_addr: http://127.0.0.1:8099/manager/gift/test
```

| 字段名    | 类型     | 示例  | 说明                                                                                                                                                      |
|-----|--------|-----|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| gift_addr    | string | http://127.0.0.1:8099/manager/gift/test    | 礼物支付的回调接口。礼物发送之前，需要业务服务提供礼物支付功能。<br/>   live 服务在发送礼物之前，回调该接口，完成实际支付。<br/> 只有支付成功的礼物，才会被发送。<br/>暂时通过 http://127.0.0.1:8099/manager/gift/test 接口返回模拟的支付成功 |


#### 礼物支付的回调接口要求

 Body 参数
使用JSON 格式数据


| 参数       | 类型     | 必填  | 说明        | 举例                  |
|----------|--------|-----|-----------|---------------------|
| LiveId   | string | 是   | 直播间ID     | 1582200377771377896 |
| UserId   | string | 是   | 用户ID      | test_user           |
| BizId    | string | 是   | 交易ID      | 1574597777432780800 |
| AnchorId | string | 是   | 当前直播间主播ID | test_anchor         |
| GiftId   | int    | 是   | 礼物ID      | 3                   |
| Amount   | int    | 是   | 礼物数量      | 99                  |

举例如下
```
{
    "biz_id": "1574597777432780800",
    "user_id":"test_user_001",
    "live_id":"1772215041443373056",
    "anchor_id": "test_user_002"  , 
    "gift_id":1  ,
    "amount": 99
}
```
返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "status":1 //支付状态，1，成功；2，失败
    }
}
```
