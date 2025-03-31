![introduce](https://socialify.git.ci/gucooing/BaPs/image?description=1&font=Source+Code+Pro&forks=1&issues=1&language=1&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

# BaPs 🎮

#### [English](README_EN.md)
  
> ⚠️ 项目仅供学习用途，严禁用于商业用途，请于24小时内删除。

# 仅供学习用途，严禁用于商业用途，请于24小时内删除！！！

> 🌟 由于是无状态设计,所以对内存的要求会略高
  
> 📅 当前支持版本：Japan 1.55.331706

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
```
---
## 🛠️ 使用方法

#### 前置准备 (此步骤非常重要！！！)

1. 前往[Releases](https://github.com/gucooing/BaPs/releases/latest)下载最新的发行版本并拷贝到运行目录（请根据自己的系统进行下载）
2. 拷贝仓库的data文件夹到运行目录
3. 下载[Releases](https://github.com/gucooing/BaPs/releases/latest)中的Excel.bin文件,并替换到data文件夹中
4. 直接运行一次将会自动生成config.json文件,打开并编辑config.json文件
5. 运行

>若Excel.bin找不到请前往源代码中data文件夹下载
---

### 🐳 Docker部署
```bash
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

---

## ⚙️ 配置说明
>需要注意的是,实际的json文件中不能存在注释
```
{
  "LogLevel": "info",
  "ResourcesPath": "./resources", // 发行版无用
  "DataPath": "./data",
  "GucooingApiKey": "123456", // 使用api时验证身份的key
  "AutoRegistration": true, // 是否自动注册
  "Tutorial": false, // 是否开启教程-不完善
  "OtherAddr": {
    "ServerInfoUrl": "https://yostar-serverinfo.bluearchiveyostar.com", // 上游服务器地址
    "ManagementDataUrl": "https://prod-noticeindex.bluearchiveyostar.com/prod/index.json" // 公告地址
  },
  "HttpNet": {
    "InnerIp": "0.0.0.0", // 监听IP
    "InnerPort": "5000", // 监听端口
    "OuterAddr": "http://127.0.0.1:5000", // 外网地址
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
前往下方仓库下载GM工具以方便使用

- [BlueArchiveGM](https://github.com/PrimeStudentCouncil/BlueArchiveGM/releases/latest)

免下载在线版GM菜单

- [BlueArchiveGM Web](https://gm.bluearchive.cc)

默认连接地址：http://127.0.0.1:5000 
默认密钥：123456 (可前往config.json进行更改)

**GM更新可能会不及时，若需要在本地使用请使用本地版**
**GM在线版支持本地使用哦~**

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

---
## 🤜 感谢名单

- 感谢 [zset](https://github.com/liyiheng/zset) 以此为基础实现排行榜
