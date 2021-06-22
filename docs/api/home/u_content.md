# /u/content

列出文章内容。只有已发布的正常内容才会被列出。

## 请求

不需要授权和前缀。

```
POST /u/content

{
    "user_id": 0,
    "user_name": "",
    "node_id": 0,
    "node_seo": "",
    "first_publish_time_begin": 0,
    "first_publish_time_end": 0,
    "publish_time_begin": 0,
    "publish_time_end": 0,
    "limit": 1,
    "page": 1,
    "sort": [
		"=id",
		"-user_id",
		"-top",
		"+sort_num",
		"-first_publish_time",
		"-publish_time",
		"-views",
		"=comment_num",
		"=bad",
		"=cool",
		"=seo"
}
```

参数说明：

| 字段   |      含义   | 类型  |   参数 |  必选 |
|----------|--------|------|------|------|
| user_id | 用户ID | int |  | 否 |
| user_name |    用户名  |  string |  | 否 |
| node_id |    内容所属节点ID |   int | 可作筛选条件 | 否 |
| node_seo |    内容所属节点SEO |   string | 可作筛选条件 | 否 |
| first_publish_begin | 在此时间后第一次发布的内容 | int | 秒，留空不筛选 | 否 |
| first_publish_end |    在此时间前第一次发布创建的内容  |  int | | 否 |
| publish_time_begin | 在此时间后最后一次发布过的内容 | int | 秒，留空不筛选| |
| publish_time_end |    在此时间前最后一次发布过的内容   |   int | | |
| sort |    内容排序   |   []string | 默认按照以上进行多列排序.=表示不排序，-表示降序，+表示升序 | 否 |
| limit |  列表每页数量 |   int | 最大页数：100，默认：20 | |
| page |  列表第几页 |   int | 默认：1 | 否 |

## 响应

正常：

```
{
  "flag": true,
  "cid": "4e20a95a9280450283e410bb240cece6",
  "data": {
    "contents": [
      {
        "id": 2,
        "seo": "1243",
        "title": "sss",
        "user_id": 1,
        "user_name": "admin",
        "node_id": 1,
        "node_seo": "qqq",
        "top": 0,
        "first_publish_time": "2019-05-27 00:41:28",
        "publish_time": "2019-09-14 15:08:07",
        "first_publish_time_int": 1558888888,
        "publish_time_int": 1568444887,
        "image_path": "",
        "views": 0,
        "is_lock": false,
        "describe": "ddddd"
      }
    ],
    "limit": 1,
    "page": 1,
    "total_pages": 2
  }
}
```


`is_lock` 表示是否有密码。

`first_publish_time`是内容第一次发布的时间，`publish_time`是内容最新一次发布的时间。