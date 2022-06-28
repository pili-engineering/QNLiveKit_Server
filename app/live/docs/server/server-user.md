用户相关接口。
# 1.用户注册
注册一个用户。

## 路径
POST /server/user/register

## 请求头 
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## Body 参数
使用JSON 格式数据

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| user_id       | string  |  是 | 客户端用户ID，唯一标识一个用户     | user_1    |
| nick          | string  |  否 | 用户昵称         | 风清扬            |
| avatar        | string  |  否 | 用户头像，URL 地址     |  https://xxx.com/avator.png |
| extends       | map[string]string |  否 | 扩展信息，key value 结构 |  |


举例如下
```
{
    "user_id": "user-1",
    "nick":"用户",
    "avatar":"http;//xxx.png",
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
}
```

# 2. 批量用户注册
批量注册用户。

## 路径
POST /server/user/register/batch

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## Body参数
使用Json 格式数据。

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| user_id       | string  |  是 | 客户端用户ID，唯一标识一个用户     | user_1    |
| nick          | string  |  否 | 用户昵称         | 风清扬            |
| avatar        | string  |  否 | 用户头像，URL 地址     |  https://xxx.com/avator.png |
| extends       | map[string]string |  否 | 扩展信息，key value 结构 |  |

举例如下
```
[
    {
        "user_id": "user-1",
        "nick":"用户",
        "avatar":"http;//xxx.png",
        "extends": {
            "age":"18"
        }
    }
]
```

## 返回
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
}
```

# 3. 更新用户信息
更新用户信息。

## 路径
PUT /server/user/{user_id}

路径参数

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| user_id       | string  |  是 | 客户端用户ID，唯一标识一个用户     | user_1    |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## Body 参数
使用JSON 格式数据

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| nick          | string  |  否 | 用户昵称         | 风清扬            |
| avatar        | string  |  否 | 用户头像，URL 地址     |  https://xxx.com/avator.png |
| extends       | map[string]string |  否 | 扩展信息，key value 结构 |  |


举例如下
```
{
    "nick":"用户",
    "avatar":"http;//xxx.png",
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
}
```


# 4. 查询用户信息
查询用户信息。

## 路径
GET /server/user/info/{user_id}

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| user_id       | string  |  是 | 客户端用户ID，唯一标识一个用户     | user_1    |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

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
        "avatar":"http;//xxx.png",
        "im_userid": "im 系统用户ID",
        "extends": {
            "age":"18"
        }
    }
}
```

# 5. 批量查询用户信息
批量查询用户信息。

## 路径
GET /server/user/infos

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## query参数
user_ids 数组参数，传递多个用户ID。

```
curl "https://live-api.qiniu.com?user_ids=user_1&user_ids=user2"
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

# 6. 查询用户详细信息
查询用户详细信息。

## 路径
GET /server/user/profile/{user_id}

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------          | ------               |
| user_id       | string  |  是 | 客户端用户ID，唯一标识一个用户     | user_1    |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

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
        "avatar":"http;//xxx.png",
        "im_userid": "im 系统用户ID",
        "im_username":"用于登录IM 的用户名",
        "im_password":"用于登录IM 用户密码",
        "extends": {
            "age":"18"
        }
    }
}
```

# 7. 批量用户详细信息
批量查询用户详细信息。

## 路径
GET /server/user/profiles

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## 请求body
Json 格式的数组，传递用户ID 列表。

```
[
    "user_1",
    "user_2"
]
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
            "im_username":"用于登录IM 的用户名",
            "im_password":"用于登录IM 用户密码",
            "extends": {
                "age":"18"
            }
        }
    ]
}
```
