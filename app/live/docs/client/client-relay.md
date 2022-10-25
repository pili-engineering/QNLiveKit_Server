# 开始跨房连麦
开始跨房连麦。

只有直播间的主播，才能发起跨房连麦请求。

## 路径
POST /client/relay/start

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |


## Body参数
```
{
    "init_room_id":"live_1",    //发起方的直播ID
    "recv_room_id":"live_2",    //接收方的直播ID
    "recv_user_id":"user_2"     //接收方的主播ID
}
```

## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "relay_id": "relay_id_1",   //跨房PK 会话ID
        "relay_status": 1,          //跨房状态：0, 等待接收方同意; 1，接收方已同意；2，发起方已完成跨房；3，接收方已完成跨房；4，两方都完成跨房；5，接收方拒绝；6，跨房已结束
        "relay_token": ""           //跨房token
    }
}
```

# 获取跨房会话信息
查询跨房连麦的会话信息。

## 路径
GET /client/relay/{relay_id}

路径参数

| 参数       | 类型     | 必填   | 说明     | 举例        |
|----------| ------  |  ----- |--------|-----------|
| relay_id | string  |  是 | 跨房会话ID | relay_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |


## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "sid": "relay_id_1",      //跨房PK 会话ID
        "init_user_id":"user_1",  //发起方的主播ID
        "init_room_id":"live_1",  //发起方的直播ID
        "recv_room_id":"live_2",  //接收方的直播ID
        "recv_user_id":"user_2"   //接收方的主播ID
        "status": 1,              //跨房状态：0, 等待接收方同意; 1，接收方已同意；2，发起方已完成跨房；3，接收方已完成跨房；4，两方都完成跨房；5，接收方拒绝；6，跨房已结束
        "start_at":"2022-08-01 15:00:00", //跨房开始时间
        "stop_at":"2022-08-01 16:00:00",  //跨房结束时间
        "extends": {
            "ext-key":"ext-value"
        }
    }
}
```

# 获取跨房会话Token
获取跨房连麦的会话Token，使用该token 可以向目的房间进行推流。

## 路径
GET /client/relay/{relay_id}/token

路径参数

| 参数       | 类型     | 必填   | 说明     | 举例        |
|----------| ------  |  ----- |--------|-----------|
| relay_id | string  |  是 | 跨房会话ID | relay_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "relay_id": "relay_id_1",   //跨房PK 会话ID
        "relay_status": 1,          //跨房状态：0, 等待接收方同意; 1，接收方已同意；2，发起方已完成跨房；3，接收方已完成跨房；4，两方都完成跨房；5，接收方拒绝；6，跨房已结束
        "relay_token": ""           //跨房token
    }
}
```

# 停止跨房
停止跨房连麦。

## 路径
POST /client/relay/{relay_id}/stop

路径参数

| 参数       | 类型     | 必填   | 说明     | 举例        |
|----------| ------  |  ----- |--------|-----------|
| relay_id | string  |  是 | 跨房会话ID | relay_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 上报跨房完成
本地推流到跨房的目的房间后，上报状态完成。更新跨房会话的状态。

## 路径
POST /client/relay/{relay_id}/started

路径参数

| 参数       | 类型     | 必填   | 说明     | 举例        |
|----------| ------  |  ----- |--------|-----------|
| relay_id | string  |  是 | 跨房会话ID | relay_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 更新跨房扩展信息
更新跨房连麦会话的扩展信息。

## 路径
POST /client/relay/{relay_id}/extends

路径参数

| 参数       | 类型     | 必填   | 说明     | 举例        |
|----------| ------  |  ----- |--------|-----------|
| relay_id | string  |  是 | 跨房会话ID | relay_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
```
{
    "extends": {
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
}
```
