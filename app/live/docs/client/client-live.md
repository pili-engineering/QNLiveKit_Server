直播间相关。

# 创建直播间
创建一个直播间。创建直播间的用户，作为直播间的主播。

## 路径
POST /client/live/room/instance

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
使用JSON 格式数据

| 参数        | 类型               | 必填  | 说明                | 举例    |
|-----------|------------------|-----|-------------------| ------               |
| title     | string           | 是   | 直播标题              | XXX直播 |
| notice    | string           | 否   | 直播公告              | 直播公告 |
| cover_url | string           | 否   | 用户头像，URL 地址       |  https://xxx.com/avator.png |
| extends   | map[string]string | 否   | 扩展信息，key value 结构 |  |
| start_at  | string         | 否   | 直播开始时间            |"2021-12-02 11:11:11"  |
| end_at    | string       | 否   | 直播结束时间            |   "2021-12-02 11:11:11"                     |
| publish_expire_at| string   |   否   | 自定义推流地址Token过期时间  |      "2021-12-02 11:11:11"                  |

举例如下
```
{
    "title": "XXX直播",
    "notice":"直播预告",
    "cover_url":"http;//xxx.png",
    "extends": {
        "age":"18"
    },
    "start_at": "2021-12-02 11:11:11"  , 
    "end_at":"2021-12-02 11:11:11"  ,
    "publish_expire_at":"2021-12-02 11:11:11"
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
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
}
```

# 查询直播间
查询直播间信息。

## 路径
GET /client/live/room/info/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
无

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
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
}
```

# 开始直播
开始直播间的直播。开始直播后，其他用户可以看到该直播间。

只有直播间的主播，才能开始直播。

直播间关闭之后，不能调用该接口开始直播。

## 路径
PUT /client/live/room/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
无

## 返回
如果开启成功，返回直播间详情信息。

```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success",  //code 非0 时，错误原因描述
    "data": {
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
}
```

# 停止直播
停止直播。

只有直播间的主播才能停止直播。

只有直播间在直播状态，才能停止直播。

## 路径
POST /client/live/room/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
无

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 删除直播间
删除直播间。

只有直播间的主播才能删除直播间。

## 路径
DELETE /client/live/room/instance/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
无

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 搜索直播
直播间搜索。

## 路径
GET /client/live/room

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Query 参数
```
{
    "page_num": 1,    //页码
    "page_size": 20,  //分页数量
    "keyword":""      //搜索关键字
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success",  //code 非0 时，错误原因描述
    "data": {
        "total_count":100, //总数
        "page_total": 5,   //总页数
        "end_page": false, //当前是否最后一页
        "list": [
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
        ]
    }
}
```

# 直播列表
分页查看直播间列表。

只会看到正在直播的直播间。

## 路径
GET /client/live/room/list

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Query 参数
```
{
    "page_num": 1,    //页码
    "page_size": 20  //分页数量
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success",  //code 非0 时，错误原因描述
    "data": {
        "total_count":100, //总数
        "page_total": 5,   //总页数
        "end_page": false, //当前是否最后一页
        "list": [
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
        ]
    }
}
```

# 加入直播
用户加入一个直播间。

一个用户在同一个时间，只能加入一个直播间。

加入第二个直播间，会主动从第一个直播间退出。

## 路径
POST /client/live/room/user/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success",  //code 非0 时，错误原因描述
    "data": {
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
}
```

# 离开直播
直播间内的用户离开直播间。

## 路径
DELETE /client/live/room/user/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
无。

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 心跳
用户加入直播间后，需要通过心跳，确认自己还在直播间。

用户加入直播间后，每5秒需要发送一个心跳。服务端在连续丢失3个心跳后，会认为用户已经离开直播间。

如果主播心跳丢失，会在直播间内提示主播离线。主播离线10分钟，会关闭直播间。

## 路径
GET /client/live/room/heartbeat/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例       |
| ---------      | ------  |  ----- |   ------          |----------|
| live_id       | string  |  是 | 直播间ID    | live_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
无。

## 返回
```
{
    "request_id": "xxxxx",
    "code": 0,
    "message": "success",
    "data": {
        "live_id": "live_123",
        "live_status": 1
    }
}
```

# 更新直播扩展
更新直播间的扩展信息。

## 路径
PUT /client/live/room/extends

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body参数
```
{
        "live_id":"live_1",  //直播间ID
	    "extends": {
            "age":"18"
        }
}
```

## 返回
```
{
    "request_id": "xxxxx",
    "code": 0,
    "message": "success"
}
```

# 房间用户列表
获取直播间内的用户列表。

## 路径
GET /client/live/room/user_list

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |


## Query参数
```
{
    "live_id":"live_123"      //直播ID
    "page_num": 1,            //页码
    "page_size": 20          //分页数量
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success",  //code 非0 时，错误原因描述
    "data": {
        "total_count":100, //总数
        "page_total": 5,   //总页数
        "end_page": false, //当前是否最后一页
        "list": [
            {
                "user_id":"user_1",  //用户ID
                "im_userid":9000990, //用户的的IM 账号ID
                "im_username":"",    //IM用户名
                "nick":"主播1", //昵称
                "avatar":"",   //用户头像
                "extends": {   //用户扩展信息
                    "age":"18"
                 },
            } 
        ]
    }
}
