直播间相关。

# 1. 创建直播间
创建一个直播间。

## 路径
POST /server/live

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## Body 参数
使用JSON 格式数据

| 参数            | 类型     | 必填   | 说明                | 举例              |
| ---------      | ------  |  ----- |-------------------| ------               |
| anchor_id     | string  |  是 | 主播的用户ID           | user_1 |
| title         | string  |  是 | 直播标题              | XXX直播 |
| notice        | string  |  否 | 直播公告              | 直播公告 |
| cover_url     | string  |  否 | 用户头像，URL 地址       |  https://xxx.com/avator.png |
 | start_at     | string  |  否 | 预计直播开始时间          | 2022-08-01 10:30:00      |
 | end_at       | string  |  否 | 预计直播结束时间          | 2022-08-01 12:30:00    |
| extends       | map[string]string |  否 | 扩展信息，key value 结构 |  |

举例如下
```
{
    "anchor_id":"user_1",
    "title": "XXX直播",
    "notice":"直播预告",
    "cover_url":"http;//xxx.png",
    "start_at": "2022-08-01 10:30:00",
    "end_at": "2022-08-01 12:30:00",
    "extends": {
        "age":"18"
    }
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

# 2. 查询直播间
查询直播间信息。

## 路径
GET /server/live/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| live_id       | string  |  是 | 直播间ID    | user_123    |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

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

# 3. 关闭直播间
关闭直播间。

## 路径
POST /server/live/{live_id}/stop

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| live_id       | string  |  是 | 直播间ID    | user_123    |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

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

# 4. 删除直播间
删除直播间。

## 路径
DELETE /server/live/{live_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| live_id       | string  |  是 | 直播间ID    | user_123    |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

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