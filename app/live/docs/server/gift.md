
# 查看礼物配置
查看所有的礼物配置。

## 路径
GET /client/gift/config

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