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

#2. 配置文件
* 互动直播使用 yaml 格式的配置文件，文件内容如下所示：
```
node_id: {{ nodeId }}

service:
  host: {{ host }}
  port: {{ port }}

jwt_key: {{ jwt_key }}

mysqls:
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

mac_config:
  enable: true
  access_key: ak
  secret_key: sk

im_config:
  app_id: {{ im_app_id }}
  endpoint: {{ im_endpoint }}
  token: {{ im_admin_token }}
  
rtc_config:
  app_id: {{ app_id }}
  access_key: {{ ak }}
  secret_key: {{ ak }}
  room_token_expire_s: 3600
  publish_key: {{ publish_key }}
  playback_url: https://{{xxx}}-playback.com
  hub: {{ 直播hub }} 
  stream_pattern: qn_live_kit-%s
  publish_url: rtmp://{{xxx}}-publish.com
  publish_domain: {{xxx}}-publish.com
  rtmp_play_url: rtmp://{{xxx}}-rtmp.com
  flv_play_url: https://{{xxx}}-hdl.com
  hls_play_url: https://{{xxx}}-hls.com

```
## 2.1 基本配置
### nodeId
```
node_id: {{ nodeId }}
```

* nodeId 是一个整数。
* 互动直播支持分布式部署，每个节点指定一个nodeId。所有节点的nodeId 不能重复。
* nodeId 会用于内部的ID 生成。


### service
```
service:
  host: {{ host }}
  port: {{ port }}
```
* 互动直播服务提供HTTP RESTful API 接口。service 配置服务的监听地址。
* host：监听的IP
* port: 监听的端口

### jwt_key
```
jwt_key: {{ jwt_key }}
```
* 字符串
* 用于客户端鉴权token 的加解密

## 2.2 数据库配置
```
mysqls:
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

## 2.3 面向服务的鉴权配置
```
mac_config:
  enable: true
  access_key: {{ ak }}
  secret_key: {{ sk }}
```
* 互动直播服务，提供面向服务端的接口，与面向客户端的接口。
* 面向服务端的接口，使用 ak/sk 的鉴权方式。
* enable：true，开启鉴权；false，关闭鉴权。如果互动直播服务，与需要调用互动直播服务接口的其他服务，在同一个安全局域网内，可以关闭服务端接口的鉴权。

## 2.4 im_config
```
im_config:
  app_id: {{ im_app_id }}
  endpoint: {{ im_endpoint }}
  token: {{ im_admin_token }}

```
* 互动直播服务，使用七牛IM 作为即时通信消息服务。
* 参考 [IM服务接入指南](https://developer.qiniu.com/IM/8332/startim) 与 [IM服务端接入指南](https://developer.qiniu.com/IM/8102/im-service-integration-guide)

## 2.5 rtc_config
```
rtc_config:
  app_id: {{ app_id }}
  access_key: {{ ak }}
  secret_key: {{ ak }}
  room_token_expire_s: 3600
  publish_key: {{ publish_key }}
  playback_url: https://{{xxx}}-playback.com
  hub: {{ 直播hub }} 
  stream_pattern: qn_live_kit-%s
  publish_url: rtmp://{{xxx}}-publish.com
  publish_domain: {{xxx}}-publish.com
  rtmp_play_url: rtmp://{{xxx}}-rtmp.com
  flv_play_url: https://{{xxx}}-hdl.com
  hls_play_url: https://{{xxx}}-hls.com
```
* 互动直播服务，使用七牛的RTC 服务，作为连麦服务。
* rtc_config 段落配置rtc 应用相关的配置。
* RTC 应用的配置，参考：[管理实时音视频](https://developer.qiniu.com/rtc/9858/applist)
* 直播相关地址配置，参考：[直播云](https://developer.qiniu.com/pili)