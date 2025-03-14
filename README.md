# BaPs

## [Discord](https://discord.gg/mmvZbCUKAG)

## 由于是无状态设计,所以对内存的要求会略高
## 现阶段只支持jp客户端
## 当前data支持版本: 1.54.327262 (不同版本不是不可用,只是会少资源)
## 由于版权原因，部分源代码将不会被公开，但我们可以保证非公开部分代码无任何恶意内容
## 由于版权原因，dev使用的resources我们不会公开

## 已实现功能
```
1.登录
2.新手教程
3.队伍管理
4.抽卡
5.剧情 待测试
6.账号基础管理
7.MomoTalk
8.邮件 全局/私人 收发管理
9.角色养成管理
10.背包管理
11.副本 - 悬赏通缉/特别依赖/学院交流会/综合战术考试
12.可恢复品自动恢复
13.咖啡厅
14.好友管理
15.课程表
16.社团
17.战斗援助
18.总力战
19.彩奈登录奖励
20.制约解除决战
21.大决战
```
## 代理方法:
转代以下地址:其中 http://127.0.0.1:5000 为服务器地址
```
https://ba-jp-sdk.bluearchive.jp http://127.0.0.1:5000
https://yostar-serverinfo.bluearchiveyostar.com http://127.0.0.1:5000
```

## 使用方法
#### 1.前往[Releases](https://github.com/gucooing/BaPs/releases/latest)下载最新的发行版本并拷贝到运行目录
#### 2.拷贝仓库的data文件夹到运行目录
#### 3.直接运行一次将会自动生成config.json文件,打开并编辑config.json文件
#### 4.运行

## Docker运行

```
docker run -d \
  -p 5000:5000 \
  -v /data/baps/config.json:/usr/ba/config.json \
  -v /data/baps/sqlite/BaPs.db:/usr/ba/BaPs.db \
  -v /data/baps/sqlite/Rank.db:/usr/ba/Rank.db \
  ghcr.io/gucooing/baps:latest
```

<details>
你展开了一个可用的镜像加速,这个镜像加速站来源于网络

```
docker run -d \
  -p 5000:5000 \
  -v /data/baps/config.json:/usr/ba/config.json \
  -v /data/baps/sqlite/BaPs.db:/usr/ba/BaPs.db \
  -v /data/baps/sqlite/Rank.db:/usr/ba/Rank.db \
  ghcr.nju.edu.cn/gucooing/baps:latest
```
</details>


# config.json
需要注意的是,实际的json文件中不能存在注释
```
{
  "LogLevel": "info",
  "ResourcesPath": "./resources", // 发行版无用
  "DataPath": "./data",
  "GucooingApiKey": "123456", // 使用api时验证身份的key
  "AutoRegistration": true, // 是否自动注册
  "Tutorial": false, // 是否开启教程-不完善
  "HttpNet": {
    "InnerAddr": "0.0.0.0", // 监听地址
    "InnerPort": "5000", // 监听端口
    "OuterAddr": "10.0.0.3", // 外网地址
    "OuterPort": "5000", // 外网端口
    "Tls": false, // 是否启用ssl
    "CertFile": "./data/cert.pem",
    "KeyFile":   "./data/key.pem"
  },
  "GateWay": {
    "MaxPlayerNum": 0, // 最大在线玩家数
    "MaxCachePlayerTime": 720, // 最大玩家缓存时间
    "BlackCmd": {}, // 发行版无用
    "IsLogMsgPlayer": true // 发行版无用
  },
  "DB": {
    "dbType": "sqlite", // 使用的数据库类型,支持sqlite和mysql
    "dsn": "BaPs.db" // 数据库地址,如果是mysql请填写mysql url
  },
  "RaidRankDB": {
    "dbType": "sqlite", // 使用的数据库类型,支持sqlite和mysql
    "dsn": "Rank.db" // 数据库地址,如果是mysql请填写mysql url
  },
  "Irc": { // 可使用通用irc服务器地址
    "HostAddress": "127.0.0.1", // 社团聊天服务器irc地址
    "Port": 16666, // 社团聊天服务器irc端口
    "Password": "mx123" // 社团聊天服务器irc密码
  }
}
```

## 我们欢迎所有想帮助我们的人加入
## 注意:玩家数据并不会实时保存到数据库中,如果有最新数据的需求,可通过api进行访问玩家数据
## Api的使用,过于复杂,没时间写docs自己研究

## 感谢[zset](https://github.com/liyiheng/zset),以此为基础实现排行榜