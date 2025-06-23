![introduce](https://socialify.git.ci/gucooing/BaPs/image?description=1&font=Source+Code+Pro&forks=1&issues=1&language=1&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

# BaPs 🎮 此存储库已过期，无法使用

#### [English](README_EN.md)
  
> ⚠️ 项目仅供学习用途，严禁用于商业用途，请于24小时内删除。

# 仅供学习用途，严禁用于商业用途，请于24小时内删除！！！

> 🌟 由于是无状态设计,所以对内存的要求会略高
  
> 📅 当前支持版本：Japan 1.57.342681

## 📍Discord

[![Discord](https://img.shields.io/badge/Join-Discord-blue?logo=discord&logoSize=auto)](https://discord.gg/222yVp6pUq)


---
## 🚀 已实现功能
```
- 登录  
- 新手教程  
- 队伍管理  
- 抽卡  
- 剧情 待测试  
- 账号基础管理  
- MomoTalk  
- 邮件 全局/私人 收发管理  
- 角色养成管理  
- 背包管理  
- 副本 - 悬赏通缉 / 特别依赖 / 学院交流会 / 综合战术考试  
- 可恢复品自动恢复  
- 咖啡厅  
- 好友管理  
- 课程表  
- 社团  
- 战斗援助  
- 总力战  
- 彩奈登录奖励  
- 制约解除决战  
- 大决战  
- 商店
- 角色好感系统
- 竞技场
- 贴纸
```
---
## 🛠️ 使用方法

#### 前置准备 (此步骤非常重要！！！)

1. 前往[Releases](https://github.com/gucooing/BaPs/releases/latest)下载最新的发行版本并拷贝到运行目录（请根据自己的系统进行下载）
2. 拷贝仓库的data文件夹到运行目录
3. 下载[Releases](https://github.com/gucooing/BaPs/releases/latest)中的Excel.bin文件,并替换到data文件夹中
4. 使用参数```-g true```运行一次将会自动生成config.json文件,打开并编辑config.json文件
5. 运行
##### tips 也可以前往[Build Action](https://github.com/gucooing/BaPs/actions/workflows/Build.yml)下载最新的版本

>若Excel.bin找不到请前往源代码中data文件夹下载
---

### 🐳 Docker部署
#### 请修改http:127.0.0.1:5000为实际服务器地址
```bash
docker run -d \
  -p 5000:5000 \
  -v /data/baps/config:/usr/ba/config \
  -v /data/baps/sqlite:/usr/ba/sqlite \
  -e Config.HttpNet.OuterAddr=http:127.0.0.1:5000 \
  ghcr.io/gucooing/baps:latest
``` 
## 更多的镜像站
``` 
docker hub: gucooing/baps:latest "由gitea编译"
ghcr: ghcr.io/gucooing/baps:latest
南京大学镜像站:
    ghcr: ghcr.nju.edu.cn/gucooing/baps:latest
```
---

## ⚙️ 配置说明
>需要注意的是,实际的json文件中不能存在注释
```json
{
  "LogLevel": "info",
  "ResourcesPath": "./resources", // 发行版无用
  "DataPath": "./data",
  "ExcelUrl": "https://github.com/gucooing/BaPs/raw/refs/heads/main/data/Excel.bin?download=",// 当Excel.bin不存在时，自动下载的url
  "GucooingApiKey": "123456", // 使用api时验证身份的key
  "AutoRegistration": true, // 是否自动注册
  "Tutorial": false, // 是否开启教程-不完善
  "OtherAddr": {
    "ServerInfoUrl": "https://yostar-serverinfo.bluearchiveyostar.com", // 上游服务器地址 如果值为: 'local' 时使用本地文件
    "ManagementDataUrl": "https://prod-noticeindex.bluearchiveyostar.com/prod/index.json" // 公告地址
  },
  "HttpNet": {
    "InnerIp": "0.0.0.0", // 监听IP
    "InnerPort": "5000", // 监听端口
    "OuterAddr": "http://127.0.0.1:5000", // 外网地址
    "Tls": false, // 是否启用ssl
    "CertFile": "./data/cert.pem",
    "KeyFile":   "./data/key.pem",
    "Encoding": true // 是否压缩数据包
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
---

## 🌐 代理设置
转代以下地址:其中 http://127.0.0.1:5000 为服务器地址
```plaintext
https://ba-jp-sdk.bluearchive.jp  →  http://127.0.0.1:5000
https://yostar-serverinfo.bluearchiveyostar.com → http://127.0.0.1:5000
```

### ⛓️代理方案

可前往以下docs查看
- [Android_MitmProxy代理方案](Android_Mitmproxy_Readme_ZH.md)

---

## ⌨️ GM工具

1. 需要注意的是此GM已近乎不可用 [BlueArchiveGM](https://github.com/PrimeStudentCouncil/BlueArchiveGM/releases/latest);web版: [BlueArchiveGM Web](https://gm.bluearchive.cc)
2. 我们欢迎更多开发者开发适用于BaPs的GM

---
## 🤝 参与贡献
我们欢迎所有想帮助我们的人加入，可通过以下方式进行帮助我们：
- 🐛 提交Issue报告问题
- 💡 提交Pull Request改进代码
- 📖 完善项目文档
- 🚀 加入Discord频道为我们提供建议
---

## ⚠️ 注意事项
1. 由于版权原因，dev使用的resources我们不会公开
2. 由于版权原因，部分源代码将不会被公开，但我们可以保证非公开部分代码无任何恶意内容
3. 玩家数据并不会实时保存到数据库中,如果有最新数据的需求,可通过api进行访问玩家数据
4. 该项目不支持,也不会适配32位系统

---
## 🤜 感谢名单

- 感谢 [zset](https://github.com/liyiheng/zset) 以此为基础实现排行榜

## 🍬 名人堂
### 此名人堂专为某些🍬人开启，各位可在此查看来自我们标注的招笑消息，以及前往discord查看本名人堂最新消息

> 名人堂永久位置：
```
@honest_quokka_88413

1. 把开源项目当成你家祖坟？
此人 fork 项目之后不学无术、半瓶醋摇晃，连代码都没看明白，就敢动后门逻辑。结果改成谁都能进，连阿猫阿狗都能调用——这怕不是把安全性当笑话，看门密码都敢贴门口。他自己手贱改出来的玩意儿，事后还说原项目就那样。

2. 协议说废就废？
项目明明挂着协议，他直接口嗨一句“协议无效”，仿佛法律是他家开的，GitHub 是他亲戚开的。协议是你说没用就没用吗？你是 GitHub CEO 还是哪路天庭来的钦差大臣？

3. 全体作者被你喂狗了？
最不要脸的操作来了：fork 回去之后，把作者信息一个不留全删了，仿佛这项目是他五年前在梦里写出来的。你删得快，咱原项目还存在呢，真以为自己删的掉吗？

详情请查看：https://discord.gg/222yVp6pUq 名人堂板块
```
> 名人堂：
```
欢迎来自各位投稿
```