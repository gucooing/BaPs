BA登录流程:

```
通过api获取登录通行证
protocol:50000
YostarUID YostarToken -----> Api Get EnterTicket

通过登录通行证向GateWay进行登录验证，并且在服务端进行玩家初始化
protocol:1009
EnterTicket ------> GateWay Account_CheckYostar GetSessionKey

*后续GateWay交互均携带GetSessionKey标识身份

最后一道登录手续,获取玩家基本信息
protocol:1002
设备信息,接入地址 ------> GateWay Account_Auth