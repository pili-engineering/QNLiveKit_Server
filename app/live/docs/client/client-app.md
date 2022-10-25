描述应用相关的接口。

# 获取应用配置信息
客户端请求互动直播服务相关的应用信息。

## 路径
GET /client/app/config

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
    "data": {
        "im_app_id": "<im app id>"   //对应的IM 应用ID
    }
}
```


