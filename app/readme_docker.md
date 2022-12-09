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

* 执行  `docker-compose --env-file live.env up`



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

censor:
  callback: https://{{xxx}}
  bucket: {{ censor_bucket }}
  addr: https://{{xxx}}

gift:
  gift_addr: http://127.0.0.1:8099/manager/gift/test


```

## 配置文件详细介绍

###  node配置
```
node_id: {{ nodeId }}
```

* nodeId 是一个整数。
* 互动直播支持分布式部署，每个节点指定一个nodeId。所有节点的nodeId 不能重复。
* nodeId 会用于内部的ID 生成。

###  面向服务的鉴权配置

```yaml
account:
  access_key: {{ .qiniu_key_pair_access_key }}
  secret_key: {{ .qiniu_key_pair_secret_key }}
```
* 互动直播服务，提供面向服务端的接口，与面向客户端的接口。
* 面向服务端的接口，使用 ak/sk 的鉴权方式。
* enable：true，开启鉴权；false，关闭鉴权。如果互动直播服务，与需要调用互动直播服务接口的其他服务，在同一个安全局域网内，可以关闭服务端接口的鉴权。

###  监听地址
```yaml
httpq:
  addr: ":8099"
```
* 互动直播服务提供HTTP RESTful API 接口。service 配置服务的监听地址`{host:ip}`。
* host：监听的IP
* port: 监听的端口

### 互动直播服务端鉴权
```yaml
auth:
  jwt_key: jwtkey
  server:
    enable: true
    access_key: ak
    secret_key: sk


```
`jwt_key`
* 字符串
* 用于客户端鉴权token 的加解密

`server`
* 互动直播服务，提供面向服务端的接口，与面向客户端的接口。
* 面向服务端的接口，使用 ak/sk 的鉴权方式。
* enable：true，开启鉴权；false，关闭鉴权。如果互动直播服务，与需要调用互动直播服务接口的其他服务，在同一个安全局域网内，可以关闭服务端接口的鉴权。

### Redis配置
```yaml
cache:
  type: {{type}}
  addr: {{host:ip}}
```
低代码服务，使用redis 作为缓存。支持redis 单节点模式，与集群模式。

* type，redis 模式，可选值如下
    * node：单节点模式。
    * cluster：集群模式

* addr：配置type为node时，必须配置单节点redis 的地址。
* addrs：配置type为cluster时，必须配置的redis集群地址列表。

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

### IM配置
```yaml
im:
  app_id: {{ im_app_id }}
  endpoint: {{ im_endpoint }}
  token: {{ im_admin_token }}
```
* 互动直播服务，使用七牛IM 作为即时通信消息服务。
* 参考 [IM服务接入指南](https://developer.qiniu.com/IM/8332/startim) 与 [IM服务端接入指南](https://developer.qiniu.com/IM/8102/im-service-integration-guide)


### RTC PILI配置
RTC
```yaml
rtc:
  app_id: {{ app_id }}
```
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

* 互动直播服务，使用七牛的RTC 服务，作为连麦服务。
* rtc_config 段落配置rtc 应用相关的配置。
* RTC 应用的配置，参考：[管理实时音视频](https://developer.qiniu.com/rtc/9858/applist)
* 直播相关地址配置，参考：[直播云](https://developer.qiniu.com/pili)
* 直播推流鉴权说明：目前直播推流鉴权，支持三种方式：无鉴权，限时鉴权，限时鉴权SK，不同鉴权模式的配置方式如下

| 鉴权模式   | security_type | publish_key                | publish_expire_s                |
|--------|---------------|----------------------------|---------------------------------|
| 无鉴权    | none          | 无需指定，留空                    | 过期时间秒。如：3600 表示 一小时后过期，推流URL 过期 |
| 限时鉴权   | expiry        | 使用配置的key 鉴权。从【直播空间设置】获取key | 同上                              |
| 限时鉴权SK | expiry_sk     | 使用RTC 用户的SK 鉴权。无需配置，留空。    | 同上                              |

### cron配置

```yaml
cron:
  single_task_node: 1
```
live 服务，有一些后台任务运行。
* single_task_node：单任务运行节点的ID。有些任务，需要全局单线程执行，只有节点ID 等于 single_task_node 的节点，才会运行单线程的任务。

###  Prometheus系统监控配置
```yaml
prome:
  client_mode: pusher
  exporter_config:
    listen_addr: ":9200"
  pusher_config:
    url: "https://{{xxx}}"
    job: live
    instance: live_{{.node}}
    interval_s: 10
```

低代码服务运行过程中，生成prometheus 监控数据指标。

支持两种方式的指标导出。

* client_mode：指标导出模式配置，支持的模式如下
    * exporter：exporter 模式，开启http 监听，由prometheus 服务来拉取。
    * pusher：pusher模式，主动向prometheus-gateway 推送指标。
* exporter_config：在client_mode 为 exporter 时，需要配置。
    * listen_addr：监听的端口。
* pusher_config:
    * url：prometheus-gateway 的指标收集地址。
    * job：任务名称
    * instance：live 服务的实例ID。
    * interval_s：推送间隔。单位：秒。


### Callback
```yaml
callback:
  addr: https://niucube-api.qiniu.com/v1/live/callback
```

* 配置低代码服务的回调地址，该地址由业务服务来实现。
* 低代码服务会将状态信息，通过回调的方式，通知给业务服务。

### kodo配置（三鉴模块）
```yaml
censor:
  callback: https://{{xxx}}
  bucket: {{ censor_bucket }}
  addr: https://{{xxx}}
```
live 服务，使用七牛的AI 审核功能，对直播内容进行内容审核。
* censor_callback：AI审核结果回调地址，配置为项目的域名。
* censor_bucket: AI审核疑似违规照片的存储bucket。
* censor_addr: bucket内存储文件的外链域名。

* 存储的bucket，与外链域名配置，参考 [对象存储使用](https://developer.qiniu.com/kodo/8452/kodo-console-guide)


### 礼物模块配置
```yaml
gift:
  gift_addr: http://127.0.0.1:8099/manager/gift/test
```
live 服务，支持直播间礼物发送功能。
* gift_addr：礼物支付的回调接口。礼物发送之前，需要业务服务提供礼物支付功能。live 服务在发送礼物之前，回调该接口，完成实际支付。只有支付成功的礼物，才会被发送。


