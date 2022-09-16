# 保存礼物配置
保存礼物配置。

对于原来没有的礼物类型，新增配置。

对于相同的配置类型，覆盖原有的配置。

## 路径
POST /server/gift/config

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Body 参数
```
[
    {
        "type": 1,          //礼物类型
        "name": "红包",      //礼物名称
        "amount": 5,        //礼物金额，0 表示自定义金额,
        "img": "",          //礼物图片,
        "animation_type":0, //动态效果类型,
        "animation_img":"", //动态效果图片,
        "order": 0          //排序，从小到大排序，相同order 根据创建时间排序',
    }
]
```

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success" //code 非0 时，错误原因描述
}
```

# 删除礼物配置
删除指定的礼物配置。

## 路径
DELETE /server/gift/config/{type}

路径参数

| 参数    | 类型       | 必填   | 说明   | 举例  |
|-------|----------|  ----- |------|-----|
| type  | integer  |  是 | 礼物类型 | 3   |

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |


## Body 参数
无。

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success" //code 非0 时，错误原因描述
}
```

# 查看礼物配置
查看所有的礼物配置。

## 路径
GET /server/gift/config

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |


## Body 参数
无。

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success" //code 非0 时，错误原因描述
    "data": [
        {
            "type": 1,          //礼物类型
            "name": "红包",      //礼物名称
            "amount": 5,        //礼物金额，0 表示自定义金额,
            "img": "",          //礼物图片,
            "animation_type":0, //动态效果类型,
            "animation_img":"", //动态效果图片,
            "order": 0,         //排序，从小到大排序，相同order 根据创建时间排序',
            "created_at":"2022-09-15 13:00:00",    //创建时间
            "updated_at":"2022-09-20 15:00:00"     //更新时间
        }
    ]
}
```

# 发送礼物
直播间内用户，在直播间发送礼物。

## 路径
PUT /server/gift/live/{live_id}

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| liveId | string       | 是   | 直播间 ID |      |


## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |


## Body 参数
```
{
    "biz_id":"",  //交易ID，唯一标识一次礼物发送，业务方生成
    "user_id":"", //发送礼物的用户ID
    "type":1,     //礼物类型，用户在礼物配置接口配置
    "amount":99,  //礼物金额
    "redo":false, //是否是重新发送 
}
```

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success" //code 非0 时，错误原因描述
}
```


# 查看直播间礼物列表
查看直播间的礼物记录。

## 路径
GET /server/gift/live/{live_id}

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| liveId | string       | 是   | 直播间 ID |      |


## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | 鉴权token |      |

## Query 参数

| 参数         | 类型      | 必填  | 说明   | 举例 |
|------------|---------|-----|------| ---- |
| page_num   | integer | 否   | 页码   |      |
| page_size  | integer | 否   | 分页大小 |      |


## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success" //code 非0 时，错误原因描述
    "data":{
        "total_count":100, //总数
        "page_total": 5,   //总页数
        "end_page": false, //当前是否最后一页
        "list":[
            {
                 "biz_id":"",  //交易ID，唯一标识一次礼物发送，业务方生成
                 "user_id":"", //发送礼物的用户ID
                 "type":1,     //礼物类型，用户在礼物配置接口配置
                 "amount":99,  //礼物金额
                 "created_at":"2022-09-01 00:00:00", //发送时间
            }
        ]
    }
}
```
