# FaFa CMS：项目文档

终极目标：实现一个可用的内容社区产品，围绕内容进行互动，知己交流。

## 简介

此项目代号为 `fafacms`。花花拼音 `fafa`，名称来源于广东话发发，花花的谐音，听起来有诙谐，娱乐等感觉，是一个使用 `Golang` 开发的前后端分离 --> 内容管理系统(CMS)。
目前基本的后端API以及配套文档都已经完成，大家可以从这个项目中学会 `Golang` 相关的开发技能，包括数据库操作，图片存储，邮件发送，容器部署等知识，五脏俱全的典型企业级应用，配置对象存储后可分布式多副本部署。

## 产品概述

1. 用户注册，填入相应信息如QQ，微博，邮箱，自我介绍，头像等，然后收到注册邮件，点击进行激活。未激活用户登陆后会显示未激活，无法使用平台。激活后用户可以登录后台，可以进行评论。用户注册后不提供注销功能。用户如果违禁被拉进黑名单不允许任何操作。用户发布内容和创建节点需要联系管理员赋予VIP权限。总结：未激活用户，普通用户，VIP用户，管理员，只有VIP用户可以创建内容，管理员可以操纵特殊权限路由。
2. 用户超级管理员高级权限控制，需要由管理员为用户分配用户组，用户组下有若干超级管理员路由资源，路由资源均为特殊路由，如更改其他用户密码，查看所有用户文章，用户信息，拉黑违禁用户等路由，如果用户不进入特殊资源路由，正常使用后台，即只能操作自己的资源，否则需要具备相应的组权限。该功能为普通用户无感知隐藏功能。
3. 用户信息一般操作，用户登录后台，进入后台后可以随时退出登录以及补充注册时的用户信息，修改密码等。用户忘记密码可以通过邮件找回。用户昵称一个月只能修改两次，且全局唯一。
4. 内容编辑，VIP用户可以创建内容节点，节点下可以有子节点，但最多两层，节点间实现了拖曳排序的功能，智能无比，在节点下可以新建文章，可以更新内容，设置隐藏文章，文章置顶，设置文章密码等，文章设计了特殊的发布机制和历史版本功能，文章内容先保存在预发布字段，点击发布按钮才真正刷新进正式字段，每次更新内容时可以将草稿保存进历史，每次发布时，会相应保存进发布历史，可以从历史内容版本中恢复等。同时可以对文章进行拖曳排序。文章实现二次删除，被删除时会移到回收站，可以从回收站恢复或彻底删除。
5. 首页阅读和内容评论，所有用户可以浏览其他用户文章并进行评论，内容所有者可以设置关闭或者开启评论，评论相对智能仿QQ音乐，评论可以由评论所有者删除。其他用户也可以为内容或者内容的某条评论点赞或者取消点赞，详细记录登陆用户点赞等情况，防止多次点赞。其他用户可以举报文章和评论。服务端可以配置自动违禁，以及举报阈值，开启时当举报超过一定次数会自动将内容或评论违禁。
6. 文件存储功能：用户头像，节点背景图，文章背景图等内部图片均需要通过上传接口保存进数据库，禁止使用不安全外部图片链接，图片存储在本地或者云对象存储服务中。文件有相应的列出，分类打标签等API功能。
7. 服务端可配置关闭用户注册，管理员权限的用户登录后台后，可以将用户加入黑名单，解除用户黑名单，激活用户，创建用户，将内容封禁，为用户赋予VIP等。
8. 互动消息站内信，如评论被点赞，内容被点赞，内容被评论，评论被评论。系统通知站内信，内容被违禁，评论被违禁，管理员通知广播。站内信会通知相应用户。
9. 关注用户，用户间可以互相关注，关注后，当某用户发布内容时，关注他的用户会收到站内信通知。
10. 私信，用户间私聊。（附加）
11. 用户和内容关键字搜索。（附加）

可以使用 [https://insomnia.rest](https://insomnia.rest/) 进行调试，本项目文件夹 `json` 下存放着最新的`调试文件`。

## 目录

* [API文档](http/README.md)
    * [接口说明](http/api.md)
    * [授权规范](http/auth.md)
        * [获取Token-/user/token/get](api/user/token_get.md)
        * [刷新Token-/user/token/refresh](api/user/token_refresh.md)
        * [删除Token-/user/token/delete](api/user/token_delete.md)
    * [文件管理](api/file/README.md)
        * [上传文件-/file/upload](api/file/upload.md)
        * [列出文件-/file/list](api/file/list.md)
        * [管理员列出文件-/file/admin/list](api/file/admin_list.md)
        * [更新文件-/file/update](api/file/update.md)
        * [更新文件(管)-/file/admin/update](api/file/admin_update.md)
    * [用户管理(管理员)](api/user/README.md)
        * [创建用户组-/group/create](api/group/create.md)
        * [更新用户组-/group/update](api/group/update.md)
        * [删除用户组-/group/delete](api/group/delete.md)
        * [获取用户组-/group/take](api/group/take.md)
        * [列出用户组-/group/list](api/group/list.md)
        * [创建用户-/user/create](api/user/create.md)
        * [更新用户(个人)-/user/update](api/user/update.md)
        * [用户绑定组-/user/assign](api/user/assign.md)
        * [列出用户-/user/list](api/user/list.md)
        * [更新用户-/user/admin/update](api/user/admin_update.md)
        * [获取用户信息(个人)-/user/info](api/user/info.md)
        * [列出资源-/resource/list](api/resource/list.md)
        * [资源绑定组-/resource/assign](api/resource/assign.md)
        * [列出组下用户-/group/user/list](api/group/user_list.md)
        * [列出组下资源-/group/resource/list](api/group/resource_list.md)
    * [内容节点管理(个人)](api/node/README.md)
        * [创建节点-/node/create](api/node/create.md)
        * [修改节点SEO-/node/update/seo](api/node/update_seo.md)
        * [修改节点描述-/node/update/info](api/node/update_info.md)
        * [修改节点图片-/node/update/image](api/node/update_image.md)
        * [修改节点状态-/node/update/status](api/node/update_status.md)
        * [修改节点父亲-/node/update/parent](api/node/update_parent.md)
        * [拖曳节点排序-/node/sort](api/node/sort.md)
        * [删除节点-/node/delete](api/node/delete.md)
        * [获取节点-/node/take](api/node/take.md)
        * [列出节点-/node/list](api/node/list.md)
        * [列出节点(管)-/node/admin/list](api/node/admin_list.md)
    * [内容管理(个人)](api/c/README.md)
        * [创建内容-/content/create](api/c/create.md)
        * [修改内容SEO-/content/update/seo](api/c/update_seo.md)
        * [修改内容图片-/content/update/image](api/c/update_image.md)
        * [修改内容状态-/content/update/status](api/c/update_status.md)
        * [修改内容状态(管)-/content/admin/update/status](api/c/admin_update_status.md)
        * [修改内容节点-/content/update/node](api/c/update_node.md)
        * [修改内容置顶-/content/update/top](api/c/update_top.md)
        * [修改内容密码-/content/update/password](api/c/update_password.md)
        * [修改内容评论设置-/content/update/comment](api/c/update_comment.md)
        * [修改内容详情-/content/update/info](api/c/update_info.md)
        * [拖曳内容排序-/content/sort](api/c/sort.md)
        * [发布内容-/content/publish](api/c/publish.md)
        * [点赞内容-/content/cool](api/c/cool.md)
        * [举报内容-/content/bad](api/c/bad.md)
        * [恢复内容-/content/restore](api/c/restore.md)
        * [内容删除-/content/rubbish](api/c/rubbish.md)
        * [内容回收-/content/recycle](api/c/recycle.md)
        * [删除内容-/content/delete](api/c/delete.md)
        * [获取内容-/content/take](api/c/take.md)
        * [获取内容(管)-/content/admin/take](api/c/admin_take.md)
        * [获取内容历史-/content/history/take](api/c/history_take.md)
        * [获取内容历史(管)-/content/history/admin/take](api/c/history_admin_take.md)
        * [列出内容-/content/list](api/c/list.md)
        * [列出内容(管)-/content/admin/list](api/c/admin_list.md)
        * [列出内容历史-/content/history/list](api/c/history_list.md)
        * [列出内容历史(管)-/content/history/admin/list](api/c/history_admin_list.md)
        * [删除内容历史-/content/history/delete](api/c/history_delete.md)
    * [评论接口(个人)](api/comment/README.md)
        * [创建评论-/comment/create](api/comment/create.md)
        * [评论取消匿名-/comment/real/name](api/comment/real_name.md)
        * [获取评论-/comment/take](api/comment/take.md)    
        * [删除评论-/comment/delete](api/comment/delete.md)   
        * [点赞评论-/comment/cool](api/comment/cool.md)
        * [举报评论-/comment/bad](api/comment/bad.md)
        * [列出评论(管)-/comment/admin/list](api/comment/admin_list.md)
        * [修改评论状态(管)-/comment/admin/update/status](api/comment/admin_update_status.md)
    * [关系接口(个人)](api/relation/README.md)
        * [添加关注-/relation/follow/add](api/relation/follow_add.md)
        * [解除关注-/relation/follow/minute](api/relation/follow_minute.md)    
        * [列出关注你的人-/relation/followed/me](api/relation/followed_me.md)
        * [列出你关注的人-/relation/following/me](api/relation/following_me.md)
        * [列出关注B用户的人-/relation/followed/list](api/relation/followed_list.md)
        * [列出A用户关注的人-/relation/following/list](api/relation/following_list.md)
        * [列出关系(管)-/relation/admin/list](api/relation/admin_list.md)
    * [信息接口(个人)](api/message/README.md)
        * [创建管理员通知(管)-/message/admin/global/create](api/message/admin_global_create.md)
        * [列出管理员通知(管)-/message/admin/global/list](api/message/admin_global_list.md)
        * [更新管理员通知状态(管)-/message/admin/global/update/status](api/message/admin_global_update_status.md)
        * [发送私信-/message/private/send](api/message/private_send.md)
        * [删除发出的私信-/message/private/delete](api/message/private_delete.md)
        * [列出站内信-/message/list](api/message/list.md)
        * [列出站内信(管)-//message/admin/list](api/message/admin_list.md)    
        * [信息已读-/message/read](api/message/read.md)
        * [删除收到的信息-/message/delete](api/message/delete.md)
    * [前端用户接口](api/home/README.md)
        * [用户注册-/user/register](api/home/user_register.md)
        * [用户激活-/user/activate](api/home/user_activate.md)
        * [激活码获取-/user/activate/code](api/home/user_activate_code.md)
        * [忘记密码验证码获取-/user/password/forget](api/home/user_pass_forget.md)
        * [用户修改密码-/user/password/change](api/home/user_pass_change.md)
    * [前端内容接口](api/home/README.md)
        * [列出用户-/u](api/home/u.md)
        * [获取用户节点详情-/u/node](api/home/u_node.md)
        * [列出用户节点-/u/nodes](api/home/u_nodes.md)
        * [获取用户信息-/u/info](api/home/u_info.md)
        * [统计用户内容情况-/u/count](api/home/u_count.md)
        * [列出用户内容-/u/content](api/home/u_content.md)
        * [获取内容-/content](api/home/content.md)
        * [获取内容下的评论-/content/comment](api/home/content_coment.md)
    * [错误码](http/errcode.md)

