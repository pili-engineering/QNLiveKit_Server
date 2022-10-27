
# 查看礼物配置
查看所有的礼物配置。

## 路径
GET /client/gift/config/{type}

路径参数

| 参数   | 类型  | 必填 | 说明                    | 举例 |
|------|-----| ---- |-----------------------| ---- |
| type | int | 是   | 类型 ID，type==-1时表示全部类型 |      |

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
            "gift_id":3,        //礼物id
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

# 查看直播间礼物列表
查看直播间的礼物记录。只有主播才能查看。

## 路径
GET  /client/gift/list/live/{live_id}

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
                  "live_id":""// 直播间ID
                  "type":1,     //礼物类型，用户在礼物配置接口配置
                 "gift_id":1, //礼物id，用户在礼物配置接口配置
                 "amount":99,  //礼物金额
                 "created_at":"2022-09-01 00:00:00", //发送时间
            }
        ]
    }
}
```


# 查看主播礼物列表
查看直播间的礼物记录。只有主播才能查看。

## 路径
GET  /client/gift/list/anchor

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
                  "live_id":""// 直播间ID
                  "type":1,     //礼物类型，用户在礼物配置接口配置
                 "gift_id":1, //礼物id，用户在礼物配置接口配置
                 "amount":99,  //礼物金额
                 "created_at":"2022-09-01 00:00:00", //发送时间
            }
        ]
    }
}
```


# 查看用户打赏礼物列表
查看直播间的礼物记录。

## 路径
GET  /client/gift/list/user

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
                  "live_id":""// 直播间ID
                  "type":1,     //礼物类型，用户在礼物配置接口配置
                 "gift_id":1, //礼物id，用户在礼物配置接口配置
                 "amount":99,  //礼物金额
                 "created_at":"2022-09-01 00:00:00", //发送时间
            }
        ]
    }
}
```


# 发送礼物
直播间内用户，在直播间发送礼物。

## 路径
POST   /client/gift/send

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
    "live_id":"",  //直播间ID
    "user_id":"", //发送礼物的用户ID
    "amount":99,  //礼物金额
}
```

## 返回
该接口正确处理请求时返回如下 JSON 数据
```
{
    "request_id": "XXX",
    "code": 0,
    "message": "success",
    "data": {
        "biz_id": "XXXX",
        "user_id": "test_XXX",
        "live_id": "XXXXX",
        "anchor_id": "XXXX",
        "gift_id": 3,
        "amount": 99,
        "status": 1
    }
```
