介绍互动直播低代码服务，回调业务服务的接口。

# 回调方式
互动直播低代码服务，使用HTTP请求回调业务服务。

### HTTP 方法
POST 

### HTTP 请求URL
业务方定义可以访问的URL。

### 请求格式
以JSON 格式携带请求内容。

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------ | ---- | --------- | ---- |
| type  | string | 是   | 回调消息类型，具体见后面内容。 |  live-created     |
| body  | json   | 是   | 回调消息内容，JSON 结构体。具体定义根据 type 不同。 |       |

```
{
    "type":"<callback-type>",
    "body": {
        //根据type 的不同，有不同的格式。
    }
}
```

### 应答格式
回调接口，应该返回如下JSON 数据。 

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------ | ---- | --------- | ---- |
| code  | int  | 是   | 处理成功返回0，不成功返回非0  |       |
| message | string | 否   | 错误消息内容 |

```bigquery
{
    "code": 0,
    "message": "错误消息"
}
```

# 直播间相关
## 创建直播间回调
* type: live-created
* body 内容

```
{
    "live_id":"live_1",  //直播间ID
    "title":"XXX直播",    //直播间标题
    "notice":"直播间公告", //直播间公告
    "cover_url":"https://xxx.png", //
    "extends": {
        "age":"18"
    },
    "anchor_info":  {
        "user_id":"user_1", //主播用户ID
        "im_userid":9000990, //主播的IM 用户ID
        "nick":"主播1", //主播昵称
        "avatar":"",   //主播用户头像
        "extends": {   //主播用户扩展信息
            "age":"18"
        },
    },
    "anchor_status":1, //主播状态：0, 离线；1，在线
    "pk_id":"",  //跨房会话ID，跨房PK 时有效
    "online_count":  3, //在线用户数
    "start_time": 1656315406, //开始直播时间，秒级时间戳
    "end_time": 1656315406, //结束直播时间，秒级时间戳
    "chat_id": 19903902, //聊天室ID
    "push_url": "", //推流地址
    "hls_url":"", //hls 拉流URL
    "rtmp_url":"", //rtmp 拉流URL
    "flv_url":"",  //flv 拉流URL
    "pv": 1003, //PV
    "uv": 934,  //UV
    "live_status":1 //直播间状态，0，已创建；1，直播中；2，已结束
}
```


## 删除直播间回调
* type: live-deleted
* body 内容

```
{
    "live_id":"live_1",  //直播间ID
}
```


## 开始直播间回调
* type: live-started
* body 内容

```
{
    "live_id":"live_1",  //直播间ID
}
```


## 删除直播间回调
* type: live-stopped
* body 内容

```
{
    "live_id":"live_1",  //直播间ID
}
```
