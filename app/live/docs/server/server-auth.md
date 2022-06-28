# 获取用户访问token
互动直播客户端，需要访问互动直播服务时，需要使用jwt 鉴权token。
该接口用于颁发访问token。

## 路径
GET /server/auth/token

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |


## Query 参数
| 参数           | 类型    | 必填 |说明             | 举例              |
|---------      | ---    |  -- |------          | ---               |
| app_id        | string  |  是 | 互动直播的应用ID  | live_app_1        |
| user_id       | string  |  是 | 客户端用户ID，唯一标识一个用户     | user_1    |
| device_id     | string  |  否 | 客户端设备ID     |  device_1        |
| expires_at    | int64   |  否 | token 过期时间，秒级时间戳。默认7 天过期 | 1655882609 |

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success", //code 非0 时，错误原因描述
    "data": {
        "access_token":"xxx",   //客户端访问token
		"expires_at":1655882609 //过期时间，秒级时间戳
    }
}
```