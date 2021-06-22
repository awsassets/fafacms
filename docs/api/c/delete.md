# /content/delete

内容真正删除，从垃圾箱中删除内容，只能删除垃圾箱中的内容。会关联删除内容历史。这个是真的删除，一旦删除，包括历史全部都从物理中消失。

删除内容时，内容下面的所有评论也会物理删除。这个API我建议不使用。

## 请求

```
POST /content/delete

{
	"id": 1
}
```

参数说明：

| 字段   |      含义   |类型  |   参数 |  必选 |
|----------|--------|------|------|------|
| id | 内容ID | int | | 是 |


## 响应

正常：

```
{
  "flag": true,
  "cid": "cfe0d18c74154f21a56169ee9b313d01"
}
```

内容不存在:

```
{
  "flag": false,
  "cid": "d65f8089f6674713981eb461a91156d2",
  "error": {
    "id": 110000,
    "msg": "content not found"
  }
}
```

不能删除内容：

```
{
  "flag": false,
  "cid": "4a6a3f70b7fc42659d7006cd83a1e56c",
  "error": {
    "id": 110007,
    "msg": "content can not delete for content not in rubbish"
  }
}
```

参数不对：

```
{
  "flag": false,
  "cid": "a591ca7626c84897b3e2ebf8f8c91d0f",
  "error": {
    "id": 100010,
    "msg": "paras input not right:Key: 'ReallyDeleteContentRequest.Id' Error:Field validation for 'Id' failed on the 'required' tag"
  }
}
```