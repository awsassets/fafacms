# /user/password/forget

忘记密码

## 请求

不需要前缀。

```
POST /user/password/forget

{
	"email": "gdccmcm14@live.com"
}
```

## 响应

正常：

```
{
  "flag": true,
  "cid": "ad2244297e604903a61964ced346b4bb"
}
```

修改密码验证码将会发送到邮箱中。

用户邮箱未找到：

```
{
  "flag": false,
  "cid": "ec36c1116b6d4bc684c9833dd94924e9",
  "error": {
    "id": 100027,
    "msg": "email not found"
  }
}
```

验证码未过期：

```
{
  "flag": false,
  "cid": "faad0d5126794bcabd6433cdb4d2c0f8",
  "error": {
    "id": 100028,
    "msg": "reset code expired time not reach"
  }
}
```

验证码过期后，才能重新请求忘记密码验证码，验证码有效期五分钟。

参数不对：

```
{
  "flag": false,
  "cid": "c61529f7d0cf4fc3a8c6ad2665ae2e17",
  "error": {
    "id": 100010,
    "msg": "paras input not right:Key: 'ForgetPasswordRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
  }
}
```