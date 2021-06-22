# 给后端开发详细的部署说明

环境要求：类Unix系统，Linux/Mac。

最快速部署请安装 `docker`,`docker-compose`, 然后可以执行一键脚本：

```
git clone github.com/hunterhug/fafacms
cd fafacms

# 先打包镜像
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh

# Linux使用install.sh
# Mac请使用install_mac.sh
cd install

# 修改邮箱地址等
vim config.yaml

chmod 777 install.sh
sudo ./install.sh
```

主要集成了`mysql:5.7.27`，`phpmyadmin:edge-4.9`和`redis:5.0.5`，端口分别为`3306`，`8000`，`6379`，

`MYSQL` 账号密码：`root/123456789`,`Redis` 密码：`123456789`，打开 `IP:8000` 登录数据库进行查看。

持久卷将会挂载在 `/opt/mydocker` 中。 具体配置和挂载卷可修改 `docker-compose.yaml` 和 `config.yaml` 文件。

运行后，请打开 `IP:8080` 进行API对接，超级管理员账户密码：`admin/admin` 。

# 详细说明

## 后端部署(常规)

你也可以裸机部署。

获取代码:

```
go get -v github.com/hunterhug/fafacms
```

代码就会保存在 `Golang GOPATH` 目录下.

运行:

```
fafacms -config=./config.yaml
```

其中`config.yaml`说明如下（具体参考实际配置）:

```
# 全局配置
DefaultConfig:
  # 服务监听端口
  WebPort: :8080
  # 是否使用对象存储，否时存储在本地
  StorageOss: false
  # 本地文件存储位置
  StoragePath: ./data/storage
  # 关闭注册功能
  CloseRegister: false
  # 打开调试日志
  LogDebug: true
  # 日志地址
  LogPath: ./data/log/fafacms_log.log

# 对象存储配置
OssConfig:
  # 区域
  Endpoint: oss-cn-shenzhen.aliyuncs.com
  # 桶名
  BucketName: syoss
  # 密钥对
  AccessKeyId: 1111
  AccessKeySecret: 11111

# 数据库配置
DbConfig:
  # 数据库名
  Name: fafa
  # 数据库地址，*必填
  Host: 127.0.0.1
  # 数据库用户名，*必填
  User: root
  # 数据库密码，*必填
  Pass: 123456789
  # 数据库端口，默认值：3306
  Port: 3306
  # 数据库最大空闲连接数
  MaxIdleConns: 10
  # 数据库最大打开连接数
  MaxOpenConns: 20
  # 数据库日志调试
  Debug: true
  # 数据库调试日志是否打印到文件中
  # 当 debug = true 时有效，false 时打印到终端
  DebugToFile: true
  # 数据库调试日志打印到的文件路径
  # 当 DebugToFile = true 时有效
  DebugToFileName: ./data/log/fafacms_db.log

# 邮件配置
MailConfig:
  # 忘记密码，激活用户时发邮件服务器
  Host: smtp-mail.outlook.com
  # 邮件服务器端口
  Port: 587
  # 账户密码
  Account: gdccmcm14@live.com
  # 账户密码
  Password: dqwngvtplopdrjjda
  # 邮件验证码内容
  From: FaFaCMS
  # 邮箱发送主题
  Subject: "FaFa CMS Code"
  # 邮箱内容，两个占位符，第二个%s为验证码，第一个是字符串功能。
  Body: "%s Code is <br/> <p style='text-align:center'>%s</p> <br/>Valid in 5 minute."

# Session设置
SessionConfig:
  # Redis地址(可改)
  RedisHost: 127.0.0.1:6379
  RedisMaxIdle: 64
  RedisMaxActive: 0
  RedisIdleTimeout: 120
  # Redis默认连接数据库(默认保持)
  RedisDB: 0
  # Redis密码(可为空,可改)
  RedisPass": 123456789
```

具体命令参数如下：

```
  -auth_skip_debug
        Auth skip debug
  -auto_ban
        Auto ban the content or comment
  -ban_time int
        Content or comment will be ban in how much bad's time (default 10)
  -can_scale
        Can scale the picture auto (default true)
  -config string
        Config file (default "./config.json")
  -email_debug
        Email debug
  -history_record
        Content history can be record (default true)
  -init_db
        Init create db table (default true)
  -scale_width int
        The width of scale size of picture (default 500)
  -session_expire_time int
        Login session expire second time, token will destroy after this time (default 604800)
  -single_login
        User can only single point login
  -time_zone int
        Time zone offset the utc (default 8)
```

正常启动如下：

```
./fafacms config=/root/fafacms/config.yaml -history_record=true -init_db=false
```

表示文章内容开启历史记录功能，开启自动违禁评论和内容，并且关闭数据库数据填充（第二次启动时可设置为false）。

## 后端部署(Docker)

你也可以单独使用`docker`进行部署, 构建镜像(Docker版本必须大于17.06):

```
sudo chmod 777 ./docker_build.sh
sudo ./docker_build.sh
````

先新建数据卷, 并且移动配置并修改:

```
mkdir /root/fafacms
cp docker_config.yaml /root/fafacms/config.yaml
```

启动容器:

```
sudo docker run -d --name fafacms -p 8080:8080 -v /root/fafacms:/root/fafacms --env RUN_OPTS="-config=/root/fafacms/config.yaml -history_record=true -init_db=true" hunterhug/fafacms

sudo docker logs -f --tail 10 fafacms
```

其中`/root/fafacms`是挂载的持久化卷, 配置`config.yaml`放置在该文件夹下.

开发中`Debug`:

```
sudo docker run -d --name fafacms -p 8080:8080 -v /root/fafacms:/root/fafacms --env RUN_OPTS="-config=/root/fafacms/config.yaml -email_debug=true -auth_skip_debug=true -history_record=true -init_db=true" hunterhug/fafacms
```
