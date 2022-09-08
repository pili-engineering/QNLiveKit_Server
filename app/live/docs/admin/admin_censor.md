三鉴接口。

# 管理员登录
## 路径
GET /manager/login

## Query 参数

| 参数            | 类型     | 必填  | 说明  | 举例               |
| ---------      |--------|-----|-----|------------------|
| user_name     | string | 是   | 用户名 | |
|password         | string | 是   | 密码  |     |



## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id": "8lemqRjDTkYT2hIX",
    "code": 0,
    "message": "success",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpQQQQ.eyJleHAQQQE2NjMyMzU3ODgsIkFwcElkIjoiIiwiVXNlcklkIjoQQQV6aXdlbiIsIkRldmljZUlkIjoiIiwiUm9sZSI6ImFkbWluIn0.nAP14pT0RJF4UDKaJTZrpmrOyKFKK9kya9xBcZTVDEk",
        "expires_at": 1663235788
    }
}
```

#修改审核设置
修改审核设置信息。

##路径
POST /admin/censor/config


## 请求头
| 参数           | 说明               | 举例              |
|----           |------------------| ---               |
| Authorization | JWT 鉴权token(管理员） |   |

## Body 参数
使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数            | 类型    | 必填  | 说明       | 举例              |
| ---------      |-------|-----|----------| ------           |
| enable   | int   | 是   | 审核结果     |  1 通过；2 违规              |
| pulp    | bool | 是   | 是否开启涉黄|                 |
|   terror      | bool  | 是   | 是否开启暴恐         |                 |
|     politician    |  bool     | 是   |       是否开启敏感人物   |                 |
|ads         |  bool     | 是   |  是否开启广告        |                 |
|    interval     | int   | 是   |  截帧时长，单位秒       |                 |

举例如下
```
{
    "enable": true,    
    "interval": 20,    
    "pulp": true,      
    "terror": true,   
    "politician":true, 
    "ads":true 
}
```

## 返回
```

{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "enable": true,
        "pulp": true,
        "terror": true,
        "politician": true,
        "ads": true,
        "interval": 20
    }
}
```

#获取审核设置
获取审核设置信息

##路径
GET /admin/censor/config


## 请求头
| 参数           | 说明               | 举例              |
|----           |------------------| ---               |
| Authorization | JWT 鉴权token(管理员） |   |


## 返回
```

{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success"  //code 非0 时，错误原因描述
    "data": {
        "enable": true,
        "pulp": true,
        "terror": true,
        "politician": true,
        "ads": true,
        "interval": 20
    }
}
```


# 查看待审核图片/已审核图片详情
## 路径
GET /admin/censor/record

## 请求头
| 参数           | 说明               | 举例              |
|----           |------------------| ---               |
| Authorization | JWT 鉴权token(管理员） |   |

## Query 参数


| 参数            | 类型     | 必填  | 说明     | 举例                 |
| ---------      |--------|-----|--------|--------------------|
| is_review     | string | 否   | 图片是否审核 | 0：没审核 1：审核 2：默认：全部 |
|page_num         | string | 否   | 第几页    | 1                  |
|  page_size              | string | 否   | 每页大小   | 10                 |
|  live_id   | string | 否   | 直播间id  |                    |


## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id": "8lemqfC2LK3w1xIX",
    "code": 200,
    "message": "success",
    "data": {
        "total_count": 955,
        "page_total": 319,
        "end_page": false,
        "list": [
            {
                "ID": 1081,
                "url": "qiniu:///niu-cube/img/p998AVbGDcYMdMpkyZur_r9FO5M=/2022-09-06-06-45-16_1662446716.jpg",
                "job_id": "6316e882000187043001d2562b6c4c63",
                "created_at": 1662466047,
                "suggestion": "review",
                "pulp": "pass",
                "terror": "review",
                "politician": "pass",
                "ads": "pass",
                "live_id": "1559814678429835264",
                "is_review": 0,
                "review_answer": 0,
                "review_user_id": "0",
                "review_time": null
            },
            {
                "ID": 1079,
                "url": "qiniu:///niu-cube/img/7xrxCmkkWAJ46zYyZHiWt6XcoQU=/2022-09-06-06-45-15_1662446715.jpg",
                "job_id": "6316e882000187043001d2562b6c4c63",
                "created_at": 1662466006,
                "suggestion": "pass",
                "pulp": "pass",
                "terror": "pass",
                "politician": "pass",
                "ads": "pass",
                "live_id": "1559814678429835264",
                "is_review": 0,
                "review_answer": 0,
                "review_user_id": "0",
                "review_time": null
            },
            {
                "ID": 1076,
                "url": "qiniu:///niu-cube/img/-0mFoIGNaT_IMEp06eZxUUCEMzU=/2022-09-06-06-45-13_1662446713.jpg",
                "job_id": "6316e882000187043001d2562b6c4c63",
                "created_at": 1662465944,
                "suggestion": "pass",
                "pulp": "pass",
                "terror": "pass",
                "politician": "pass",
                "ads": "pass",
                "live_id": "1559814678429835264",
                "is_review": 0,
                "review_answer": 0,
                "review_user_id": "0",
                "review_time": null
            }
        ]
    }
}
```


# 查看待审核直播间/已审核直播间
## 路径
GET /admin/censor/live

## 请求头
| 参数           | 说明               | 举例              |
|----           |------------------| ---               |
| Authorization | JWT 鉴权token(管理员） |   |

## Query 参数


| 参数        | 类型     | 必填  | 说明     | 举例                      |
|-----------|--------|-----|--------|-------------------------|
| audit     | string | 否   | 图片是否审核 | 0：全部直播间（默认 1：审核 2：默认：全部 |
| page_num  | string | 否   | 第几页    | 1                       |
| page_size | string | 否   | 每页大小   | 10                      |

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id": "8lemqZh5gZEJ0BIX",
    "code": 200,
    "message": "success",
    "data": {
        "total_count": 5,
        "page_total": 2,
        "end_page": false,
        "list": [
            {
                "live_id": "1559814678429835264",
                "title": "XXX直播",
                "anchor_id": "test_user_222",
                "live_status": 2,
                "count": 960,
                "time": 1662453912
            },
            {
                "live_id": "1555033434651369472",
                "title": "XXX直播",
                "anchor_id": "test_user_222",
                "live_status": 2,
                "count": 11,
                "time": 1662540372
            },
            {
                "live_id": "1561926251407482880",
                "title": "AAAAAAAAAAAAAAAAAAA",
                "anchor_id": "test_user_222",
                "live_status": 0,
                "count": 54,
                "time": null
            }
        ]
    }
}
```


# 审核图片
## 路径
GET /admin/censor/audit

## 请求头
| 参数           | 说明               | 举例              |
|----           |------------------| ---               |
| Authorization | JWT 鉴权token(管理员） |   |

## Body 参数
使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数            | 类型    | 必填   | 说明       | 举例              |
| ---------      |-------|  ----- |----------| ------           |
| review_answer   | int   |  是     | 审核结果     |  1 通过；2 违规              |
| image_list    | []int |  是     | 待审核图片的id |                 |

举例如下
```
{
    "image_list": [1077,1078,1070],
    "review_answer": 2
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

