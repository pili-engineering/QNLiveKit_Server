直播间商品管理接口。

# 批量导入商品到直播间
## 路径
POST /server/item/add

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## Body 参数
使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------        | ------           |
| live_id       | string    |  是     | 直播间 ID          |                 |
| items         | []Item |  是     | 商品信息          |                 |

Item 的结构如下

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------        | ------           |
| item_id       | string  |  是     | 商品 ID           |                 |
| title         | string  |  是     | 商品标题           |                 |
| tags          | string  |  否     | 商品标签，多个以 ","分割。 |            |
| thumbnail     | string  |  否     | 缩略图 url              |            |
| link          | string  |  否     | 商品链接 url           |             |
| current_price | string  |  是     | 当前价格字符串          |       199元       |
| origin_price  | string  |  否     | 原始价格字符串          |       ￥298       |
| status        | int     |  是     | 0，下架(用户不可见)；1，上架(用户可见)；2，上架不能购买  |            |
| extends       | map[string]string |  否 | 扩展信息，key value 结构 |  |

举例如下
```
{
    "live_id":"live_1", 
    "items":[
        {
            "live_id":"live_1",
            "item_id":"item_1",
            "title": "精品皮鞋",
            "tags": "皮具,一口价",
            "thumbnail":"http;//xxx.png",
            "link":"http;//xxx.png",
            "current_price":"199元",
            "origin_price":"298元",
            "status":1,
            "extends": {
                "age":"18"
            }
        }
    ]
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


# 批量删除直播间商品
## 路径
POST /server/item/delete

## 请求头
| 参数           | 说明            | 举例              |
|----           | ---            | ---               |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign     |

## Body 参数
使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------        | ------           |
| live_id       | string  |  是     | 直播间 ID          |                 |
| items       | []string  |  是     | 商品 ID 列表           |                 |

举例如下
```
{
    "live_id":"live_1",
    "items":["item_1", "item_2"]
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

# 批量更新直播间商品状态

## 路径

POST /server/item/status

## 请求头

| 参数          | 说明            | 举例          |
| ------------- | --------------- | ------------- |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign |

## Body 参数

使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| live_id | string       | 是   | 直播间 ID |      |
| items   | []ItemStatus | 是   | 商品状态列表 |     |

ItemStatus 字段定义如下

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| item_id | string | 是   | 商品 ID   |      |
| status  | int    | 是   | 商品状态  |      |

举例如下

```
{
    "live_id":"live_1",
    "items": [
        {
            "item_id":"item_1",
            "status":1
        }
    ]
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

# 批量调整直播间商品顺序

## 路径

POST /server/item/order

## 请求头

| 参数          | 说明            | 举例          |
| ------------- | --------------- | ------------- |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign |

## Body 参数

使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| live_id | string       | 是   | 直播间 ID |      |
| items   | []ItemOrder  | 是   | 商品状态列表 |     |

ItemOrder 字段定义如下

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| item_id | string | 是   | 商品 ID   |      |
| order | int    | 是   | 调整后商品序号  |      |

举例如下

```
{
    "live_id":"live_1",
    "items": [
        {
            "item_id":"item_1",
            "order":1
        }
    ]
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

# 更新商品信息

## 路径

PUT /server/item/{liveId}/{itemId}

## 请求头

| 参数          | 说明            | 举例          |
| ------------- | --------------- | ------------- |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign |


## 路径参数

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| liveId | string       | 是   | 直播间 ID |      |
| itemId | string       | 是   | 商品ID |     |


## Body 参数

使用JSON 格式数据，传递商品信息。商品信息字段如下：

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------        | ------           |
| title         | string  |  否     | 商品标题           |                 |
| tags          | string  |  否     | 商品标签，多个以 ","分割。 |            |
| thumbnail     | string  |  否     | 缩略图 url              |            |
| link          | string  |  否     | 商品链接 url           |             |
| current_price | string  |  否     | 当前价格字符串          |       199元       |
| origin_price  | string  |  否     | 原始价格字符串          |       ￥298       |
| extends       | map[string]string |  否 | 扩展信息，key value 结构 |  |

以上字段未填写则保持原来的信息不更新。

举例如下

```
{
    "title": "精品皮鞋",
    "tags": "皮具,一口价",
    "thumbnail":"http;//xxx.png",
    "link":"http;//xxx.png",
    "current_price":"199元",
    "origin_price":"298元",
    "status":1,
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

# 更新商品扩展信息

## 路径

PUT /server/item/{liveId}/{itemId}/extends

## 请求头

| 参数          | 说明            | 举例          |
| ------------- | --------------- | ------------- |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign |


## 路径参数

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------       | ---- | --------- | ---- |
| liveId | string       | 是   | 直播间 ID |      |
| itemId | string       | 是   | 商品ID |     |


## Body 参数

使用JSON 格式数据，传递商品扩展信息。支持任意自定义的json 字段名，数据类型为字符串即可。

| 参数            | 类型     | 必填   | 说明             | 举例              |
| ---------      | ------  |  ----- |   ------        | ------           |
| 自定义         | string  |  是     | 参数任意自定义，如：key1 |                 |

只更新请求携带的扩展信息。

举例如下

```
{
    "age": "18",
    "key1": "value1",
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

# 查看直播间商品

## 路径

GET /server/items/{live_id}

## 请求头

| 参数          | 说明            | 举例          |
| ------------- | --------------- | ------------- |
| Authorization | AK/SK 鉴权token | Qiniu AK:sign |

## 路径参数

使用JSON 格式数据，传递商品信息数组。商品信息字段如下：

| 参数    | 类型   | 必填 | 说明      | 举例 |
| ------- | ------ | ---- | --------- | ---- |
| live_id | string | 是   | 直播间 ID |      |

举例如下

```
curl "https://live-api.qiniu.com/server/items/live_1"
```

## 返回

返回结果，按照order 倒序排序。

```
{
    "request_id":"xxxxx", //请求ID
    "code": 0,            //错误码：0，成功；其他，失败
    "message": "success",  //code 非0 时，错误原因描述
    "data":[
        {
            "live_id":"live_1",
            "item_id":"item_1",
            "order":1,
            "title": "精品皮鞋",
            "tags": "皮具,一口价",
            "thumbnail":"http;//xxx.png",
            "link":"http;//xxx.png",
            "current_price":"199元",
            "origin_price":"298元",
            "status":1,
            "extends": {
                "age":"18"
            }
        }
    ]
}
```

