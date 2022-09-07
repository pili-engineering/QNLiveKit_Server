# 申请上麦
用户在直播间内，请求上麦。请求上麦成功后，返回上麦token，用户使用该token 发布视频语音流。

用户必须先加入直播间。

## 路径
POST /client/mic

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body参数
```
{
    "live_id":"live_123",   //直播间ID
    "mic": true,            //是否开启麦克风
    "camera": true,         //是否开启摄像头
    "extend":{
        "ext-key":"ext-value"
    }
}
```

## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "rtc_token": ""   //返回连麦token，使用该token 完成连麦。
    }
}
```

# 申请下麦
已上麦的用户申请下麦。

## 路径
DELETE /client/mic

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body参数
```
{
    "live_id":"live_123",   //直播间ID
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 房间麦位列表
获取直播间内的麦位列表信息。

## 路径
GET /client/mic/room/list/{live_id}

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
    "message": "success"  //code 非0 时，错误原因描述
    "data": [{
        "user": {
            "user_id":"user_1",  //用户ID
            "im_userid":9000990, //用户的的IM 账号ID
            "im_username":"",    //IM用户名
            "nick":"主播1", //昵称
            "avatar":"",   //用户头像
            "extends": {   //用户扩展信息
                "ext-key":"ext-value"
             }
        },
        "mic": true,    //是否开启麦克风
        "camera": true, //是否开启摄像头
        "status": 1,    //麦位状态：0, 在线；1，离开；2, 禁用
        "extends": {    //扩展信息
            "ext-key":"ext-value"
        } 
    }]
}
```

# 更新麦位扩展信息
更新麦位的扩展信息。

## 路径
PUT /client/mic/extension

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body参数
```
{
    "live_id":"live_123",   //直播ID
    "extends": {    //扩展信息
        "ext-key":"ext-value"
    } 
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 打开/关闭麦克风/摄像头
上麦用户更改自己的推流配置。

可以打开或者关闭麦克风；打开或者关闭摄像头。

## 路径
PUT /client/mic/switch

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body参数
```
{
    "live_id":"live_123",   //直播ID
    "user-id":"user_123",   //上麦用户ID
    "type":"mic",           //类型：mic 麦克风；
    "status": true          //true 打开; false 关闭
}
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```