提供给互动直播客户端的接口。

互动直播SDK 使用这部分接口与互动直播服务进行通信。

业务APP可以使用这部分接口，进行客户端业务扩展。

# 认证方式
客户端接口，使用在请求头部，携带鉴权token的方式，来进行认证。

| 参数           | 说明        | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## 鉴权token 获取
鉴权token，由业务服务根据用户身份来生成。

业务服务生成鉴权token，参考 [获取用户访问token](https://github.com/pili-engineering/QNLiveKit_Server/blob/develop/app/live/docs/server/server-auth.md)

# 返回结果
API 的返回结果，格式如下。

```
{
    "request_id":"",  //请求ID
    "code": 0,        //错误码：0，成功；其他，失败
    "message": "",    //错误原因，当 code 非 0 时，错误原因必须携带。
    "data":{          //返回结果，具体参考每个API 定义
    }
}
```