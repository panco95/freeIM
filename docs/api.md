
# FreeIM-server api doc

> v1.0.0

Base URLs:

* <a href="http://127.0.0.1:8081/api/v1">服务端本机测试: http://127.0.0.1:8081/api/v1</a>

# 登录/注册

## GET 获取图形验证码

GET /captcha/image

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|type|query|string| 是 |登录：login | 注册：register|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "key": "774bd7f0-f4b1-46dc-91ee-ce49d5e8d8b2",
    "captcha": "iVBORw0KGgoAAAANSUhEUgAAAFIAAAAgCAIAAAAOpLgGAAAHd0lEQVR4nOxZWUxT2xpeq8PedIINtKUyFIgRsJShIKjEFwzGByVFog+k8cUBDHFKjCExJEbjS41TwAdjgqhJw4vRQAgGMGjgQdOAkWCLA4NgW9vuEguF0u5N976RdcI5V7rbXs65x6vnfk+l+advre//1+pCwLIs+OeB96ML+DH4P+1/EsLThhD+7ZX8rQhDG3H+tZmHoY1m+6894cOL/NfmDAAQ/OgCfsP6nvqvLn102n6/X7iKqJahUMjlcpnNZpIkCYLYvn27Wq2OsY41kqFQiM/n/0fDhWVZmqZjqXAN4WkHg0GLxfLu3TuHw8Hn81NTU7OzswsKCkQiEVegQCDw7Nmz1tbW6elpv98PAFAoFCdPnjx69GjUIhiGmZqampiYWFxcTElJEQqFOTk5UXd7bm7u/fv3U1NTXAsUKQK7DhRFdXR06PX6LVu2SKVSCCGO45WVlV1dXUtLS+vtWZZdWVl5/fp1SUkJn88Xi8VFRUUYhvF4vKysrBcvXoR1WcPS0lJ/f39jY+OuXbvS0tJwHC8vLzcajSRJRvByuVytra3V1dV5eXmJiYkAgPz8/JaWltnZ2cjpEMLQ9ng8dXV1Uqk0Ly9Po9EolUoMwyCEBoNhbGwsbBS/33/u3DkIoUgkunr16ujo6OXLlwEAQqHw0KFDEdIHg0GTySSTySQSSXp6OkEQPN63KZudnf3kyRMuL4fDcfz4cckqFAqFRCJBfVFQUHD//v0N0vZ6vW1tbc3NzUNDQw6Ho6ura+vWrQCA8vLyoaGhsFHm5+fVajWEcNOmTR6Ph2EYm82GOiIhIcHv93OlDwQCT58+FQqFer3+7t27AwMDFRUVAsG31mtpaVlYWAjr9fHjx/r6eq1W29TU1NnZ2dPTYzAYAABSqfTIkSO/cwunZU7aNE2TJOl0OtGfk5OTu3fvhhDm5+fbbLawUSYmJtAu7dixA33DMExhYSHqo9HRUa70LMs6nc7Ozk7UPgzD3Lp1KzU1FcfxGzduUBQV1sXtdj98+LCjowN5LS8vP3/+HImrpqbmj5y5mIcZaQKBQC6Xo880TY+NjQ0PDyuVytra2rXvv8OnT58EAgFFUampqRRFoaZQq9UWiyUUCtnt9vz8fKTD9UhKStq7dy+O44iPw+Hw+XxlZWXbtm3jGs5yubyurg6JAgAQFxfHsqxMJvP5fKWlpYuLi1KplGVZCCHXVOM8wBwOh9FoHBgYePv2LYQwNzdXp9PRNI3j+Hpjr9eLYRhFUUKhkKZpDMOQ5Hg8XigUWlpaYhiGizbiFggEuru7Ozs7u7u7aZpOSEhA9qj671wghGuc0Zn3+PFjAEBGRkZhYeFahREmOecvsLm5uZGREYvFAgAQiUQkSd65c8dkMrlcrvXGEEJElaZpiqJ+j76qfAzDIp9GDMMMDw83NDSYTCav18uy7Pj4+OnTp69fv26z2SI4IvT39w8ODlIUVVVVpdVqYznAOXdboVA0NDSgVhkeHjabzX19fR6PRywWHz58+DvjxMTEuLg4AIDH41leXkYnyvz8PGKrUCi4thqBx+PFx8fr9fqVlZXk5GSr1frhw4eRkRGfz4dh2KlTpyL4TkxMdHR0WK1WjUZz4MCB7OzsqJxjAk3TNpvtwoULAACCIBobG9dPiNnZWY1GAyFUqVRjY2MMw/h8Pp1Oh2a71WqNMNIQ/H7/9PT0+Pi40+m02+1XrlxBKjt48GAEL6/Xe+3aNaVSSRCE0WgMhUJRE/020mK5+kokkoSEBNSBDMOsN1CpVDqdbmZmhiRJk8lUX1/f29vrcrlYlq2qqpLJZFFTiESirKws9JlhmPT09Li4uOXl5fj4+Lm5ueTk5PUuDMOMjIzcu3fP7XafPXv2xIkTnz9/drvdNE2XlJQg9XEhvMhZljWbzYFAQK1Wy+Vys9lsMpl4PB6fz6+url5ZWfnjREFjyWAwvHnzxmq1tre39/T0kKtQqVR1dXUpKSlRaQMAvn79GgwGRSKR0+l89OgRhFAsFhcXF3PdiN1u982bNycnJzMyMkQi0e3bt202m91ux3H8/PnzZWVlkZJxyeDSpUt79uwpLS1FKw0hTEtLa2pqstvtYe0XFhYePHhQUVGBlhlCqNVqTSYT15XjO/T29jY3N585c2b//v1IHQRBHDt27MuXL1wuFy9eTEpKAgCIxWIkRoTMzMz29vYoIudajoyMDJIkLRYLTdMEQWi12n379tXW1iqVyrD2MplMr9cXFBTMzMwsLi7K5fLNmzcjrcay1QRBvHz5cnBwkKZpAEBubm5NTY3BYFCpVFwuGIZ5vV70GxG1iU6nKy4u3rlzZ2VlZeR0nAf65OQkuoEnrSIlJSU5OXltIEe4CVAUhbpAKBTG/uMxEAi0tbW9evVKo9FotdqcnJzMzEwcxyNE6OvrMxqNRUVFulVotdoYc0WijfTPMMw3Sfx7J6+V8udfAr5bPoZhgsEghDBGjWw878ZKj7DbPzBU7NjgO/lfWOifDLWxF96f+98DG37b/l95QowFG3k84sDPRPsv7KyfW+Qbxr8CAAD//0OtPKokoDW4AAAAAElFTkSuQmCC"
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» key|string|true|none|验证码key|none|
|»» captcha|string|true|none|验证码base64|none|
|» message|string|true|none||none|

## GET 获取邮箱验证码

GET /captcha/email

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|email|query|string| 是 |邮箱地址|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "发送成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 账号登陆

POST /login/basic

> Body 请求参数

```yaml
account: test123
password: test123
captchaKey: string
captcha: string
platform: android

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» account|body|string| 是 |用户名/邮箱/手机号|
|» password|body|string| 是 |账号密码|
|» captchaKey|body|string| 是 |验证码key|
|» captcha|body|string| 是 |验证码|
|» platform|body|string| 是 |平台：android  ios  h5|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiZXhwIjoxNjY5NDU5MjA5LCJpYXQiOjE2Njg4NTQ0MDksImlzcyI6IlNaS0oifQ.3NTC29HU7tftXWVoU6NxtPvNAc8LdqC316LkdNiDxtQ",
    "needUpdatePassword": false
  },
  "message": "登陆成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» token|string|true|none||用户Token|
|»» needUpdatePassword|boolean|true|none||是否需要设置密码|
|» message|string|true|none||none|

## GET 获取手机验证码

GET /captcha/mobile

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|mobile|query|string| 是 |none|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "发送成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 账号注册

POST /register/basic

> Body 请求参数

```yaml
account: test123
password: test123
captchaKey: string
captcha: string
platform: android

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» account|body|string| 是 |用户名|
|» password|body|string| 是 |账号密码|
|» captchaKey|body|string| 是 |验证码key|
|» captcha|body|string| 是 |验证码|
|» platform|body|string| 是 |平台：android  ios  h5|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNjY5NDU3OTQ0LCJpYXQiOjE2Njg4NTMxNDQsImlzcyI6IlNaS0oifQ.bvX-YxbTSYx6XCGamEPOquNjP-qGZvFhz7r_jmZw4zQ",
    "needUpdatePassword": false
  },
  "message": "注册成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» token|string|true|none||用户token|
|»» needUpdatePassword|boolean|true|none||是否需要设置密码|
|» message|string|true|none||none|

## POST 邮箱登录

POST /login/email

> Body 请求参数

```yaml
email: 1129443982@qq.com
captcha: test123
platform: android

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» email|body|string| 是 |邮箱地址|
|» captcha|body|string| 是 |验证码|
|» platform|body|string| 是 |平台：android  ios  h5|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MywiZXhwIjoxNjY5NDU4NjQxLCJpYXQiOjE2Njg4NTM4NDEsImlzcyI6IlNaS0oifQ.Y_Y9CgRqKuSX6RvCmvmFLuNq2zZp4c3CM8QsVFure-U",
    "needUpdatePassword": true
  },
  "message": "登陆成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» token|string|true|none||用户Token|
|»» needUpdatePassword|boolean|true|none||是否需要设置密码|
|» message|string|true|none||none|

## POST 手机号登录

POST /login/mobile

> Body 请求参数

```yaml
mobile: "15712345678"
captcha: string
platform: android

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» mobile|body|string| 是 |手机号|
|» captcha|body|string| 是 |验证码|
|» platform|body|string| 是 |平台：android  ios  h5|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NCwiZXhwIjoxNjY5NDU4NzU0LCJpYXQiOjE2Njg4NTM5NTQsImlzcyI6IlNaS0oifQ.--68agPJb7-4fzhdnk14z_2c2mYA6grFdXrkLhxA1Gw",
    "needUpdatePassword": true
  },
  "message": "登陆成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» token|string|true|none||用户Token|
|»» needUpdatePassword|boolean|true|none||是否需要设置密码|
|» message|string|true|none||none|

## POST 通过邮箱重置密码

POST /reset/password/email

> Body 请求参数

```yaml
email: string
captcha: string
password: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» email|body|string| 是 |邮箱地址|
|» captcha|body|string| 是 |验证码|
|» password|body|string| 是 |新密码|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "设置成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 通过手机号重置密码

POST /reset/password/mobile

> Body 请求参数

```yaml
mobile: string
captcha: string
password: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» mobile|body|string| 是 |手机号码|
|» captcha|body|string| 是 |验证码|
|» password|body|string| 是 |新密码|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "设置成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

# 个人信息

## GET 我的信息

GET /me/info

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 否 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "account": {
      "id": 10000039,
      "username": "p00004",
      "email": "",
      "mobile": "",
      "nickname": "p00004",
      "avatar": "",
      "gender": "",
      "birth": "",
      "age": 0,
      "intro": "",
      "longitude": 0,
      "latitude": 0,
      "country": "",
      "province": "",
      "city": "",
      "district": ""
    }
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» account|object|true|none||none|
|»»» id|integer|true|none|账号ID|none|
|»»» username|string|true|none|用户名|none|
|»»» email|string|true|none|邮箱|none|
|»»» mobile|string|true|none|手机号|none|
|»»» nickname|string|true|none|昵称|none|
|»»» avatar|string|true|none|头像|none|
|»»» gender|string|true|none|性别|none|
|»»» birth|string|true|none|生日|none|
|»»» age|integer|true|none|年龄|none|
|»»» intro|string|true|none|个性签名|none|
|»»» longitude|integer|true|none|经度|none|
|»»» latitude|integer|true|none|纬度|none|
|»»» country|string|true|none|国家|none|
|»»» province|string|true|none|省份|none|
|»»» city|string|true|none|城市|none|
|»»» district|string|true|none|区县|none|
|» message|string|true|none||none|

## PUT 修改我的密码

PUT /me/update/password

> Body 请求参数

```yaml
password: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» password|body|string| 是 |新密码|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## PUT 修改个人信息

PUT /me/update/info

> Body 请求参数

```yaml
nickname: string
avatar: string
longitude: 0
latitude: 0
intro: string
gender: string
country: string
province: string
city: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» nickname|body|string| 否 |昵称（传空字符串不更新）|
|» avatar|body|string| 否 |头像（传空字符串不更新）|
|» longitude|body|number| 否 |经度（传0串不更新）|
|» latitude|body|number| 否 |纬度（传0不更新）|
|» intro|body|string| 否 |个人介绍（传空字符串不更新）|
|» gender|body|string| 否 |性别：男 女（传空字符串不更新）|
|» country|body|string| 否 |国家（传空字符串不更新）|
|» province|body|string| 否 |省（传空字符串不更新）|
|» city|body|string| 否 |城市（传空字符串不更新）|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

# 好友

## GET 搜索好友

GET /friends/search

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|account|query|string| 是 |账号ID/用户名/昵称/邮箱/手机号|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 10000001,
        "username": "panco002",
        "email": "",
        "mobile": "",
        "nickname": "panco002",
        "avatar": "",
        "gender": "",
        "birth": "",
        "age": 0,
        "intro": ""
      },
      {
        "id": 10000003,
        "username": "panco002",
        "email": "",
        "mobile": "",
        "nickname": "panco002",
        "avatar": "",
        "gender": "",
        "birth": "",
        "age": 0,
        "intro": ""
      }
    ],
    "total": 2
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||通用数组|
|»»» id|integer|true|none||用户ID|
|»»» username|string|true|none||用户名|
|»»» email|string|true|none||电子邮件|
|»»» mobile|string|true|none||手机号|
|»»» nickname|string|true|none||昵称|
|»»» avatar|string|true|none||头像|
|»»» gender|string|true|none||性别|
|»»» birth|string|true|none||生日|
|»»» age|integer|true|none||年龄|
|»»» intro|string|true|none||个人介绍|
|»» total|integer|true|none||总数|
|» message|string|true|none||none|

## POST 添加好友

POST /friends/add

> Body 请求参数

```yaml
toId: "10000001"
reason: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |对方ID|
|» reason|body|string| 否 |申请添加理由|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 同意/拒绝好友请求

POST /friends/add/reply

> Body 请求参数

```yaml
toId: "1000001"
status: pass
reason: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |对方ID|
|» status|body|string| 是 |通过：pass | 拒绝：deny|
|» reason|body|string| 否 |拒绝原因|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## GET 好友请求列表

GET /friends/add/applies

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "createdAt": "2022-11-20T16:03:10.614+08:00",
        "fromAccount": {
          "id": 10000000,
          "username": "panco001",
          "email": "",
          "mobile": "",
          "nickname": "panco001",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "applyReason": "",
        "denyReason": "",
        "status": "wait",
        "replyTime": null
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||通用数组|
|»»» createdAt|string|false|none||请求添加时间|
|»»» fromAccount|object|false|none||请求用户|
|»»»» id|integer|true|none||用户ID|
|»»»» username|string|true|none||用户名|
|»»»» email|string|true|none||邮箱|
|»»»» mobile|string|true|none||手机号|
|»»»» nickname|string|true|none||昵称|
|»»»» avatar|string|true|none||头像|
|»»»» gender|string|true|none||性别|
|»»»» birth|string|true|none||生日|
|»»»» age|integer|true|none||年龄|
|»»»» intro|string|true|none||个人介绍|
|»»» applyReason|string|false|none||添加理由|
|»»» denyReason|string|false|none||none|
|»»» status|string|false|none||none|
|»»» replyTime|null|false|none||none|
|»» total|integer|true|none||总数|
|» message|string|true|none||none|

## GET 好友列表(或黑名单)

GET /friends

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|blacklist|query|integer| 是 |是否黑名单。0不是 | 1是|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "account": {
          "id": 10000000,
          "username": "panco001",
          "email": "",
          "mobile": "",
          "nickname": "panco001",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "friendGroups": [
          {
            "id": 17,
            "friends": null,
            "name": "朋友q1"
          }
        ],
        "remark": "小美",
        "label": "客户,经理,老板"
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» account|object|false|none||好友信息|
|»»»» id|integer|true|none||好友ID|
|»»»» username|string|true|none||用户名|
|»»»» email|string|true|none||邮箱|
|»»»» mobile|string|true|none||手机号|
|»»»» nickname|string|true|none||昵称|
|»»»» avatar|string|true|none||头像|
|»»»» gender|string|true|none||性别|
|»»»» birth|string|true|none||生日|
|»»»» age|integer|true|none||年龄|
|»»»» intro|string|true|none||介绍|
|»»» friendGroups|[object]|false|none||好友分组数组|
|»»»» id|integer|false|none||分组ID|
|»»»» name|string|false|none||分组名称|
|»»» remark|string|false|none||备注|
|»»» label|string|false|none||自定义字段(标签)|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## DELETE 删除好友

DELETE /friends

> Body 请求参数

```yaml
toId: "1000001"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |好友ID|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "删除成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 添加黑名单

POST /friends/blaklist

> Body 请求参数

```yaml
toId: "1000001"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |用户ID|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## DELETE 移除黑名单

DELETE /friends/blaklist

> Body 请求参数

```yaml
toId: "1000001"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |用户ID|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## PUT 设置好友备注

PUT /friends/remark

> Body 请求参数

```yaml
toId: "1000001"
remark: 小美

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |用户ID|
|» remark|body|string| 是 |备注|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## PUT 设置好友自定义字段(标签)

PUT /friends/label

> Body 请求参数

```yaml
toId: "1000001"
label: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |用户ID|
|» label|body|string| 是 |自定义字段(标签)，格式前端定 (json或分隔符)|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 创建好友分组

POST /friends/groups

> Body 请求参数

```yaml
name: 朋友
members:
  - "1"
  - "2"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» name|body|string| 是 |分组名称|
|» members|body|array| 否 |成员数组（好友ID）|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## DELETE 删除好友分组

DELETE /friends/groups

> Body 请求参数

```yaml
groupId: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |分组ID|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "删除成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## GET 获取好友分组列表

GET /friends/groups

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 17,
        "friends": [
          {
            "id": 2,
            "account": {
              "id": 10000000,
              "username": "panco001",
              "email": "",
              "mobile": "",
              "nickname": "panco001",
              "avatar": "",
              "gender": "",
              "birth": "",
              "age": 0,
              "intro": ""
            },
            "remark": "小美",
            "label": "客户,经理,老板"
          }
        ],
        "name": "朋友q1"
      },
      {
        "id": 18,
        "friends": [],
        "name": "朋友q11"
      },
      {
        "id": 21,
        "friends": [],
        "name": "朋友q1122"
      }
    ],
    "total": 3
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» id|integer|true|none||none|
|»»» friends|[object]|true|none||none|
|»»»» id|integer|false|none||none|
|»»»» account|object|false|none||none|
|»»»»» id|integer|true|none||none|
|»»»»» username|string|true|none||none|
|»»»»» email|string|true|none||none|
|»»»»» mobile|string|true|none||none|
|»»»»» nickname|string|true|none||none|
|»»»»» avatar|string|true|none||none|
|»»»»» gender|string|true|none||none|
|»»»»» birth|string|true|none||none|
|»»»»» age|integer|true|none||none|
|»»»»» intro|string|true|none||none|
|»»»» remark|string|false|none||none|
|»»»» label|string|false|none||none|
|»»» name|string|true|none||none|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## POST 好友分组添加成员

POST /friends/groups/members

> Body 请求参数

```yaml
groupId: "1"
members:
  - ""

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |好友分组ID|
|» members|body|array| 是 |成员数组（好友ID）|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "添加成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## DELETE 好友分组删除成员

DELETE /friends/groups/members

> Body 请求参数

```yaml
groupId: "1"
members:
  - ""

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |好友分组ID|
|» members|body|array| 是 |成员数组（好友ID）|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "删除成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## GET 通过ID获取好友分组

GET /friends/group

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|groupId|query|integer| 是 |好友分组ID|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "id": 17,
    "friends": [
      {
        "id": 2,
        "account": {
          "id": 10000000,
          "username": "panco001",
          "email": "",
          "mobile": "",
          "nickname": "panco001",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "remark": "小美",
        "label": "客户,经理,老板"
      }
    ],
    "name": "朋友q1"
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» id|integer|true|none||none|
|»» friends|[object]|true|none||none|
|»»» id|integer|false|none||none|
|»»» account|object|false|none||none|
|»»»» id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» email|string|true|none||none|
|»»»» mobile|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» avatar|string|true|none||none|
|»»»» gender|string|true|none||none|
|»»»» birth|string|true|none||none|
|»»»» age|integer|true|none||none|
|»»»» intro|string|true|none||none|
|»»» remark|string|false|none||none|
|»»» label|string|false|none||none|
|»» name|string|true|none||none|
|» message|string|true|none||none|

## PUT 重命名好友分组

PUT /friends/groups/name

> Body 请求参数

```yaml
groupId: "1"
name: 朋友

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |好友分组ID|
|» name|body|string| 是 |分组名称|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "id": 17,
    "friends": [
      {
        "id": 2,
        "account": {
          "id": 10000000,
          "username": "panco001",
          "email": "",
          "mobile": "",
          "nickname": "panco001",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "remark": "小美",
        "label": "客户,经理,老板"
      }
    ],
    "name": "朋友q1"
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» id|integer|true|none||none|
|»» friends|[object]|true|none||none|
|»»» id|integer|false|none||none|
|»»» account|object|false|none||none|
|»»»» id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» email|string|true|none||none|
|»»»» mobile|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» avatar|string|true|none||none|
|»»»» gender|string|true|none||none|
|»»»» birth|string|true|none||none|
|»»»» age|integer|true|none||none|
|»»»» intro|string|true|none||none|
|»»» remark|string|false|none||none|
|»»» label|string|false|none||none|
|»» name|string|true|none||none|
|» message|string|true|none||none|

## GET 单个好友信息

GET /friends/info

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|toId|query|integer| 是 |好友信息|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "account": {
          "id": 10000000,
          "username": "panco001",
          "email": "",
          "mobile": "",
          "nickname": "panco001",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "friendGroups": [
          {
            "id": 17,
            "friends": null,
            "name": "朋友q1"
          }
        ],
        "remark": "小美",
        "label": "客户,经理,老板"
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» account|object|false|none||好友信息|
|»»»» id|integer|true|none||好友ID|
|»»»» username|string|true|none||用户名|
|»»»» email|string|true|none||邮箱|
|»»»» mobile|string|true|none||手机号|
|»»»» nickname|string|true|none||昵称|
|»»»» avatar|string|true|none||头像|
|»»»» gender|string|true|none||性别|
|»»»» birth|string|true|none||生日|
|»»»» age|integer|true|none||年龄|
|»»»» intro|string|true|none||介绍|
|»»» friendGroups|[object]|false|none||好友分组|
|»»»» id|integer|false|none||none|
|»»»» friends|null|false|none||none|
|»»»» name|string|false|none||none|
|»»» remark|string|false|none||备注|
|»»» label|string|false|none||自定义字段(标签)|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## GET 校验好友（是否你是他的好友或者黑名单）

GET /friends/verify

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|toId|query|integer| 是 |好友ID|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "isFriend": true,
    "isBlacklist": false
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» isFriend|boolean|true|none||你是否对方的好友|
|»» isBlacklist|boolean|true|none||你是否对方的黑名单|
|» message|string|true|none||none|

## GET 附近的人

GET /friends/near

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|longitude|query|number| 是 |经度|
|latitude|query|number| 是 |纬度|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 10000036,
        "username": "a123456",
        "email": "",
        "mobile": "",
        "nickname": "a123456",
        "avatar": "",
        "gender": "",
        "birth": "",
        "age": 0,
        "intro": "",
        "longitude": 113.88761132955551,
        "latitude": 22.802167782129153,
        "country": "",
        "province": "",
        "city": "",
        "district": "",
        "distance": "245948m"
      },
      {
        "id": 10000021,
        "username": "15000000001",
        "email": "",
        "mobile": "",
        "nickname": "15000000001",
        "avatar": "",
        "gender": "",
        "birth": "",
        "age": 0,
        "intro": "",
        "longitude": 113.91609638929367,
        "latitude": 22.578574959221548,
        "country": "",
        "province": "",
        "city": "",
        "district": "",
        "distance": "255891m"
      },
      {
        "id": 10000029,
        "username": "15000000003",
        "email": "",
        "mobile": "",
        "nickname": "15000000003",
        "avatar": "",
        "gender": "",
        "birth": "",
        "age": 0,
        "intro": "",
        "longitude": 113.91647189855576,
        "latitude": 22.578661139740966,
        "country": "",
        "province": "",
        "city": "",
        "district": "",
        "distance": "255924m"
      }
    ],
    "total": 3
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||通用数组|
|»»» id|integer|true|none||用户ID|
|»»» username|string|true|none||用户名|
|»»» email|string|true|none||电子邮件|
|»»» mobile|string|true|none||手机号|
|»»» nickname|string|true|none||昵称|
|»»» avatar|string|true|none||头像|
|»»» gender|string|true|none||性别|
|»»» birth|string|true|none||生日|
|»»» age|integer|true|none||年龄|
|»»» intro|string|true|none||个人介绍|
|»»» longitude|number|true|none||经度|
|»»» latitude|number|true|none||纬度|
|»»» country|string|true|none||国家|
|»»» province|string|true|none||省份|
|»»» city|string|true|none||城市|
|»»» district|string|true|none||区县|
|»»» distance|string|true|none||距离|
|»» total|integer|true|none||总数|
|» message|string|true|none||none|

# 群组

## POST 创建群

POST /chatGroups

> Body 请求参数

```yaml
name: string
intro: string
avatar: string
disableAddMember: "0"
disableViewMember: "0"
disbaleAddGroup: "0"
enbaleBeforeMsg: "0"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» name|body|string| 是 |群组名称|
|» intro|body|string| 是 |介绍|
|» avatar|body|string| 否 |群头像|
|» disableAddMember|body|integer| 否 |禁止加成员好友|
|» disableViewMember|body|integer| 否 |禁用查看成员资料|
|» disbaleAddGroup|body|integer| 否 |禁用主动申请入群|
|» enbaleBeforeMsg|body|integer| 否 |是否开启加群之前的漫游消息|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "创建成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## PUT 修改群资料

PUT /chatGroups

> Body 请求参数

```yaml
groupId: "1"
name: string
intro: string
avatar: string
disableAddMember: "0"
disableViewMember: "0"
disbaleAddGroup: "0"
enbaleBeforeMsg: "0"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» name|body|string| 是 |群组名称|
|» intro|body|string| 是 |介绍|
|» avatar|body|string| 否 |群头像（传空字符串不更新）|
|» disableAddMember|body|integer| 否 |禁止加成员好友|
|» disableViewMember|body|integer| 否 |禁用查看成员资料|
|» disbaleAddGroup|body|integer| 否 |禁用主动申请入群|
|» enbaleBeforeMsg|body|integer| 否 |是否开启加群之前的漫游消息|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "修改成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## GET 我的群列表

GET /chatGroups

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 1,
        "name": "02群",
        "avatar": "",
        "intro": "02群",
        "members": 2,
        "members_limit": 500,
        "disableAddMember": false,
        "disableViewMember": false,
        "disbaleAddGroup": false,
        "enbaleBeforeMsg": false
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» id|integer|false|none||none|
|»»» name|string|false|none||none|
|»»» avatar|string|false|none||none|
|»»» intro|string|false|none||none|
|»»» members|integer|false|none||none|
|»»» members_limit|integer|false|none||none|
|»»» disableAddMember|boolean|false|none||none|
|»»» disableViewMember|boolean|false|none||none|
|»»» disbaleAddGroup|boolean|false|none||none|
|»»» enbaleBeforeMsg|boolean|false|none||none|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## POST 加群申请

POST /chatGroups/join

> Body 请求参数

```yaml
groupId: "1"
reason: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» reason|body|string| 是 |加群原因|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "申请成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## GET 加群审批列表

GET /chatGroups/join

> Body 请求参数

```yaml
{}

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "createdAt": "2022-11-25T09:39:52.15+08:00",
        "account": {
          "id": 10000005,
          "username": "panco0002",
          "email": "",
          "mobile": "",
          "nickname": "panco0002",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "chatGroup": {
          "id": 1,
          "name": "02群",
          "avatar": "",
          "intro": "02群",
          "members": 1,
          "members_limit": 500,
          "disableAddMember": false,
          "disableViewMember": false,
          "disbaleAddGroup": false,
          "enbaleBeforeMsg": false
        },
        "applyReason": "我想加群",
        "denyReason": "",
        "status": "wait",
        "replyTime": null
      },
      {
        "createdAt": "2022-11-25T09:40:16.306+08:00",
        "account": {
          "id": 10000006,
          "username": "panco0003",
          "email": "",
          "mobile": "",
          "nickname": "panco0003",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "chatGroup": {
          "id": 1,
          "name": "02群",
          "avatar": "",
          "intro": "02群",
          "members": 1,
          "members_limit": 500,
          "disableAddMember": false,
          "disableViewMember": false,
          "disbaleAddGroup": false,
          "enbaleBeforeMsg": false
        },
        "applyReason": "我想加群",
        "denyReason": "",
        "status": "wait",
        "replyTime": null
      }
    ],
    "total": 2
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» createdAt|string|true|none||申请时间|
|»»» account|object|true|none||申请用户信息|
|»»»» id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» email|string|true|none||none|
|»»»» mobile|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» avatar|string|true|none||none|
|»»»» gender|string|true|none||none|
|»»»» birth|string|true|none||none|
|»»»» age|integer|true|none||none|
|»»»» intro|string|true|none||none|
|»»» chatGroup|object|true|none||申请目标群信息|
|»»»» id|integer|true|none||none|
|»»»» name|string|true|none||none|
|»»»» avatar|string|true|none||none|
|»»»» intro|string|true|none||none|
|»»»» members|integer|true|none||none|
|»»»» members_limit|integer|true|none||none|
|»»»» disableAddMember|boolean|true|none||none|
|»»»» disableViewMember|boolean|true|none||none|
|»»»» disbaleAddGroup|boolean|true|none||none|
|»»»» enbaleBeforeMsg|boolean|true|none||none|
|»»» applyReason|string|true|none||申请理由|
|»»» denyReason|string|true|none||none|
|»»» status|string|true|none||none|
|»»» replyTime|null|true|none||none|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## POST 加群审批

POST /chatGroups/join/reply

> Body 请求参数

```yaml
groupId: "1"
accountId: "1"
status: pass
reason: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» accountId|body|integer| 是 |加群用户ID|
|» status|body|string| 是 |通过pass  拒绝deny|
|» reason|body|string| 否 |拒绝理由|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "处理成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## GET 群信息(包括成员列表)

GET /chatGroups/info

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|groupId|query|integer| 是 |群ID|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "id": 2,
    "name": "测试群1",
    "avatar": "",
    "intro": "测试群1",
    "members": 2,
    "members_limit": 500,
    "membersList": [
      {
        "account": {
          "id": 10000004,
          "username": "panco0001",
          "email": "",
          "mobile": "",
          "nickname": "panco0001",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "role": "owner",
        "remark": "",
        "isBanned": false
      },
      {
        "account": {
          "id": 10000005,
          "username": "panco0002",
          "email": "",
          "mobile": "",
          "nickname": "panco0002",
          "avatar": "",
          "gender": "",
          "birth": "",
          "age": 0,
          "intro": ""
        },
        "role": "general",
        "remark": "",
        "isBanned": true
      }
    ],
    "selfInfo": {
      "chatGroupMemberRole": "general",
      "isBanned": true
    },
    "disableAddMember": false,
    "disableViewMember": false,
    "disbaleAddGroup": false,
    "enbaleBeforeMsg": false
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» id|integer|true|none||群ID|
|»» name|string|true|none||群名称|
|»» avatar|string|true|none||群头像|
|»» intro|string|true|none||群介绍|
|»» members|integer|true|none||群员数|
|»» members_limit|integer|true|none||最大群员数|
|»» membersList|[object]|true|none||群员列表|
|»»» account|object|true|none||用户信息|
|»»»» id|integer|true|none||none|
|»»»» username|string|true|none||none|
|»»»» email|string|true|none||none|
|»»»» mobile|string|true|none||none|
|»»»» nickname|string|true|none||none|
|»»»» avatar|string|true|none||none|
|»»»» gender|string|true|none||none|
|»»»» birth|string|true|none||none|
|»»»» age|integer|true|none||none|
|»»»» intro|string|true|none||none|
|»»» role|string|true|none||群角色：owner群主 manager管理员 general普通用户|
|»»» remark|string|true|none||群昵称|
|»»» isBanned|boolean|true|none||是否禁言中|
|»» selfInfo|object|true|none||我的群信息|
|»»» chatGroupMemberRole|string|true|none||我的群角色|
|»»» isBanned|boolean|true|none||是否禁言中|
|»» disableAddMember|boolean|true|none||是否禁止加成员好友|
|»» disableViewMember|boolean|true|none||是否禁用查看成员资料|
|»» disbaleAddGroup|boolean|true|none||是否禁用主动申请入群|
|»» enbaleBeforeMsg|boolean|true|none||是否开启加群之前的漫游消息|
|» message|string|true|none||none|

## POST 解散群

POST /chatGroups/dissolve

> Body 请求参数

```yaml
groupId: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "处理成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 退出群

POST /chatGroups/exit

> Body 请求参数

```yaml
groupId: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "退出成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 转让群

POST /chatGroups/transfer

> Body 请求参数

```yaml
groupId: "1"
toId: "2"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» toId|body|integer| 是 |被装让用户ID|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 踢出群员

POST /chatGroups/kick

> Body 请求参数

```yaml
groupId: "1"
toId: "2"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» toId|body|integer| 是 |踢出用户ID|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 设置管理员

POST /chatGroups/manager

> Body 请求参数

```yaml
groupId: "1"
toId: "2"
isManager: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» toId|body|integer| 是 |被设置用户ID|
|» isManager|body|integer| 是 |1设置，0取消|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 禁言群员

POST /chatGroups/banned

> Body 请求参数

```yaml
groupId: "1"
toId: "2"
minute: "10"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» groupId|body|integer| 是 |群组ID|
|» toId|body|integer| 是 |被禁言用户ID|
|» minute|body|integer| 是 |禁言时长(分钟)|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": null,
  "message": "处理成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|null|true|none||none|
|» message|string|true|none||none|

## POST 搜索群

POST /chatGroups/search

> Body 请求参数

```yaml
name: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» name|body|string| 是 |群组名称|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 1,
        "name": "01群",
        "avatar": "",
        "intro": "01群开放加入了",
        "members": 1,
        "members_limit": 500,
        "disableAddMember": false,
        "disableViewMember": false,
        "disbaleAddGroup": false,
        "enbaleBeforeMsg": false
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» id|integer|false|none||群ID|
|»»» name|string|false|none||群名称|
|»»» avatar|string|false|none||群头像|
|»»» intro|string|false|none||群介绍|
|»»» members|integer|false|none||群员数|
|»»» members_limit|integer|false|none||最大群员数|
|»»» disableAddMember|boolean|false|none||是否禁止加成员好友|
|»»» disableViewMember|boolean|false|none||是否禁用查看成员资料|
|»»» disbaleAddGroup|boolean|false|none||是否禁用主动申请入群|
|»»» enbaleBeforeMsg|boolean|false|none||是否开启加群之前的漫游消息|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## GET 获取好友共同群

GET /chatGroups/common

> Body 请求参数

```yaml
toId: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |none|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 1,
        "name": "02群",
        "avatar": "",
        "intro": "02群",
        "members": 2,
        "members_limit": 500,
        "disableAddMember": false,
        "disableViewMember": false,
        "disbaleAddGroup": false,
        "enbaleBeforeMsg": false
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» id|integer|false|none||none|
|»»» name|string|false|none||none|
|»»» avatar|string|false|none||none|
|»»» intro|string|false|none||none|
|»»» members|integer|false|none||none|
|»»» members_limit|integer|false|none||none|
|»»» disableAddMember|boolean|false|none||none|
|»»» disableViewMember|boolean|false|none||none|
|»»» disbaleAddGroup|boolean|false|none||none|
|»»» enbaleBeforeMsg|boolean|false|none||none|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

## GET 附近的群

GET /chatGroups/near

> Body 请求参数

```yaml
name: string

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|longitude|query|number| 否 |经度|
|latitude|query|number| 否 |none|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» name|body|string| 是 |群组名称|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "items": [
      {
        "id": 1,
        "name": "01群",
        "avatar": "",
        "intro": "01群开放加入了",
        "members": 1,
        "members_limit": 500,
        "disableAddMember": false,
        "disableViewMember": false,
        "disbaleAddGroup": false,
        "enbaleBeforeMsg": false
      }
    ],
    "total": 1
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» items|[object]|true|none||none|
|»»» id|integer|false|none||群ID|
|»»» name|string|false|none||群名称|
|»»» avatar|string|false|none||群头像|
|»»» intro|string|false|none||群介绍|
|»»» members|integer|false|none||群员数|
|»»» members_limit|integer|false|none||最大群员数|
|»»» disableAddMember|boolean|false|none||是否禁止加成员好友|
|»»» disableViewMember|boolean|false|none||是否禁用查看成员资料|
|»»» disbaleAddGroup|boolean|false|none||是否禁用主动申请入群|
|»»» enbaleBeforeMsg|boolean|false|none||是否开启加群之前的漫游消息|
|»» total|integer|true|none||none|
|» message|string|true|none||none|

# 聊天

## CONNECT 发送正在输入

CONNECT /api/v1/ws

通过websocket发送消息给对方自己正在输入和完成输入：
{
"toId": 10000012,   //好友ID
"ope": "friend",
"type": "input",      //消息类型：input输入
"status": "normal" //状态：normal正在输入  finish完成输入
}
对方收到类型为input的消息通知，fromId为好友ID

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 发送消息

POST /chat/send

body字段说明：
{
 "content": "普通主文本消息内容",   //当type==1时（文字消息），此字段必须有值
 "thumb": "图片缩略图，为base64编码",//当type==2（图片消息） 和type==4（视频消息），此字段必须有值
 "fileUrl": "文件的云端URL地址",//当type==2这里是图片的url，type==3这里是语音文件url，type==4是视频url，type==6文件url
 "filePath": "文本的本地路径",//图片/语音/视频/文件 的本地url，用在发送者查看自己发送的消息时（消息发送中或发送失败时，查看的是本地的文件）
 "fileName": "文件的名称"//type==6时，这里是文件的名称，
}

> Body 请求参数

```yaml
toId: "10000012"
ope: friend
type: text
isPrivate: "1"
body: hello

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |接收方ID(好友或群组)|
|» ope|body|string| 是 |friend好友  group群|
|» type|body|string| 是 |text文字 pic图片 voice语音 video视频 geo地理 file文件|
|» isPrivate|body|string| 是 |1加密 0不加密|
|» body|body|string| 是 |消息内容|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "id": 1,
    "createdAt": "2022-12-01T00:06:09.58+08:00",
    "updatedAt": "2022-12-01T00:06:09.58+08:00",
    "fromId": 10000011,
    "toId": 10000012,
    "ope": "friend",
    "type": "text",
    "body": "hello",
    "isPrivate": true,
    "status": "normal"
  },
  "message": "发送成功"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» id|integer|true|none||消息ID|
|»» createdAt|string|true|none||创建时间|
|»» fromId|integer|true|none||发送方ID|
|»» toId|integer|true|none||接收方ID|
|»» ope|string|true|none||消息通道|
|»» type|string|true|none||消息类型|
|»» body|string|true|none||消息内容|
|»» isPrivate|boolean|true|none||是否私密消息|
|»» status|string|true|none||消息状态 normal正常 revocation已撤回|
|» message|string|true|none||none|

## POST 撤回消息

POST /chat/revocation

撤回成功会给对方发送消息：
{
        "id": 2,
        "createdAt": "2022-12-01T00:12:24.799+08:00",
        "updatedAt": "2022-12-01T00:12:24.799+08:00",
        "fromId": 10000011,
        "toId": 10000012,
        "ope": "friend",
        "type": "text",
        "body": "hello",
        "isPrivate": true,
        "status": "revocation"
    }

> Body 请求参数

```yaml
toId: "1"
id: "1"
ope: friend

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 是 |接收方ID|
|» id|body|integer| 是 |消息ID|
|» ope|body|string| 是 |消息通道 friend好友  group群组|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## POST 已读消息

POST /chat/read

已读成功会给对方发送消息：
{
"id": 0,
"createdAt": "2022-12-01T00:12:24.799+08:00",
"updatedAt": "2022-12-01T00:12:24.799+08:00",
"fromId": 10000011,
"toId": 10000012,
"ope": "friend",
"type": "read",  //主要看这个字段 
"body": "hello",
"read": true,
"isPrivate": true,
"status": "normal"
}

> Body 请求参数

```yaml
toId: "1"

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|
|body|body|object| 否 |none|
|» toId|body|integer| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

## GET 获取聊天记录

GET /chat/messages

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|toId|query|integer| 是 |对方ID（好友获取群组）|
|messageType|query|string| 否 |消息类型，跟发消息时一致|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

# 其他

## GET 获取七牛云参数

GET /upload/qiniu/params

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Authorization|header|string| 是 |用户Token|

> 返回示例

> 成功

```json
{
  "code": 0,
  "result": {
    "accessDomain": "http://rlw127vpd.hn-bkt.clouddn.com",
    "uploadToken": "U7nSFEZVJ7gnNNqTrsvN5BeRzDrbQMTlHChzCybC:JQfGwQ41NgKxhyb3OWrm_zeDVuY=:eyJzY29wZSI6ImZyZWVpbSIsImRlYWRsaW5lIjoxNjcwNTcwMTczfQ=="
  },
  "message": ""
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» result|object|true|none||none|
|»» accessDomain|string|true|none||资源访问域名|
|»» uploadToken|string|true|none||资源上传Token|
|» message|string|true|none||none|

# 数据模型

