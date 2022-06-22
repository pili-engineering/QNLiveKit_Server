服务端接口，使用AK/SK 方式进行鉴权。

# 请求认证
## API 域名
```
https://live-api.qiniu.com
```
* 私有部署时，使用自己的域名。
* 私有部署时，可以选择不鉴权，具体配置方式参考 [项目说明文档](https://github.com/pili-engineering/QNLiveKit_Server)

每一个请求均需在 HTTP 请求头部增加一个 Authorization 字段，其值为符合管理凭证的字符串，形式如下：
```
Authorization: <QiniuToken>
```

## 生成管理凭证
`<QiniuToken>`: 管理凭证，用于鉴权。以golang 为例，生成方式如下：

```
// golang

// 构造待签名的 Data
// 1. 添加 Path
data = "<Method> <Path>"

// 2. 添加 Query，前提: Query 存在且不为空
if "<RawQuery>" != "" {
        data += "?<RawQuery>"
}

// 3. 添加 Host
data += "\nHost: <Host>"

// 4. 添加 Content-Type，前提: Content-Type 存在且不为空
if "<Content-Type>" != "" {
        data += "\nContent-Type: <Content-Type>"
}

// 5. 添加回车
data += "\n\n"

// 6. 添加 Body，前提: Content-Length 存在且 Body 不为空，同时 Content-Type 存在且不为空或 "application/octet-stream"
bodyOK := "<Content-Length>" != "" && "<Body>" != ""
contentTypeOK := "<Content-Type>" != "" && "<Content-Type>" != "application/octet-stream"
if bodyOK && contentTypeOK {
        data += "<Body>"
}

// 计算 HMAC-SHA1 签名，并对签名结果做 URL 安全的 Base64 编码
sign = hmac_sha1(data, "Your_Secret_Key")
encodedSign = urlsafe_base64_encode(sign)

// 将 Qiniu 标识与 AccessKey、encodedSign 拼接得到管理凭证
<QiniuToken> = "Qiniu " + "Your_Access_Key" + ":" + encodedSign

```

## 请求示例
```
GET /server/auth/token HTTP/1.1
Host: live-api.qiniu.com
Authorization: Qiniu 7O7hf7Ld1RrC_fpZdFvU8aCgOPuhw2K4eapYOdII:PGTUV-oRxWAIl6mdneJPSJieyyQ=
Content-Type: application/json
```