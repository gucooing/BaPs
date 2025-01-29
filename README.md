# BaPs

## [Discord](https://discord.gg/mmvZbCUKAG)

## 由于是无状态设计,所以对内存的要求会略高
## 现阶段仍然有许多bug且只支持jp客户端
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
11.副本 - 悬赏通缉/特别依赖/学院交流会
12.可恢复品自动恢复
13.咖啡厅
14.好友管理
15.课程表
```
## 代理方法:
转代以下地址:其中 http://127.0.0.1:5000 为服务器地址
```
https://ba-jp-sdk.bluearchive.jp/account/yostar_auth_request http://127.0.0.1:5000/account/yostar_auth_request
https://ba-jp-sdk.bluearchive.jp/account/yostar_auth_submit http://127.0.0.1:5000/account/yostar_auth_submit
https://ba-jp-sdk.bluearchive.jp/user/yostar_createlogin http://127.0.0.1:5000/user/yostar_createlogin
https://ba-jp-sdk.bluearchive.jp/user/login http://127.0.0.1:5000/user/login
https://prod-gateway.bluearchiveyostar.com:5100/api/gateway http://127.0.0.1:5000/getEnterTicket/gateway
https://prod-game.bluearchiveyostar.com:5000/api/gateway http://127.0.0.1:5000/api/gateway
```
如果你无法转代上面的地址可以添加下面的转代规则:
```
https://yostar-serverinfo.bluearchiveyostar.com http://127.0.0.1:5000
```

## 我们欢迎所有想帮助我们的人加入
## 注意:玩家数据并不会实时保存到数据库中,如果有最新数据的需求,可通过api进行访问玩家数据
## Api的使用,过于复杂,没时间写docs自己研究