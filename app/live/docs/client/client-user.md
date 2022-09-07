# 获取自己的详细信息接口
用户获取自己的详细信息。包括自己私有信息，如：IM用户名，密码。

## 路径
GET /client/user/profile

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
        "user_id": "user_123",   //用户ID
        "nick": "",              //昵称
        "avatar": "",            //头像
        "extends": {             //扩展信息
            "ext-key":"ext-value"
        },
        "im_userid": 12345,      //IM用户ID
        "im_username": "im_user_123", //IM用户名
        "im_password": "*****"        //IM用户密码
    }
}
```

# 获取其他用户信息
查询其他用户的信息。不包含私有信息。

## 路径
GET /client/user/user/{user_id}

路径参数

| 参数      | 类型     | 必填   | 说明   | 举例       |
|---------| ------  |  ----- |------|----------|
| user_id | string  |  是 | 用户ID | user_123 |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## 请求body
无

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "user_id": "user-1",
        "nick":"用户",
        "avatar":"http://xxx.png",
        "im_userid": "im 系统用户ID",
        "extends": {
            "age":"18"
        }
    }
}
```

# 批量获取用户信息
批量获取其他用户的公开信息。

## 路径
GET /client/user/users

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |     |

## query参数
user_ids 数组参数，传递多个用户ID。

```
curl "https://live-api.qiniu.com/client/user/users?user_ids=user_1&user_ids=user_2"
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": [
        {
            "user_id": "user-1",
            "nick":"用户",
            "avatar":"http;//xxx.png",
            "im_userid": "im 系统用户ID",
            "extends": {
                "age":"18"
            }
        }
    ]
}
```

# 更新用户信息
更新自己的用户信息。

## 路径
PUT /client/user/user

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |     |

## Body 参数
```
{
    "nick": "feng",
    "avatar":"beauty",
    "extends": {
        "score":"15",
        "lost":"hehe"
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

# 根据IM用户id 查找用户信息
根据IM 用户ID，查找用户信息。

## 路径
GET /client/user/imusers

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |     |

## query参数
im_user_ids 数组参数，传递多个IM用户ID。

```
curl "https://live-api.qiniu.com/client/user/musers?im_user_ids=123&im_user_ids=456"
```

## 返回结果
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": [
        {
            "user_id": "user-1",
            "nick":"用户",
            "avatar":"http;//xxx.png",
            "im_userid": "im 系统用户ID",
            "extends": {
                "age":"18"
            }
        }
    ]
}
```