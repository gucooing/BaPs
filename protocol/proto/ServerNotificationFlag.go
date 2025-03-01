package proto

type ServerNotificationFlag int32

const (
	ServerNotificationFlag_None                           ServerNotificationFlag = 0
	ServerNotificationFlag_NewMailArrived                 ServerNotificationFlag = 4   // 新邮件
	ServerNotificationFlag_HasUnreadMail                  ServerNotificationFlag = 8   // 邮件未读
	ServerNotificationFlag_NewToastDetected               ServerNotificationFlag = 16  // 新通知
	ServerNotificationFlag_CanReceiveArenaDailyReward     ServerNotificationFlag = 32  // 竞技场每日奖励
	ServerNotificationFlag_CanReceiveRaidReward           ServerNotificationFlag = 64  // 总力战总分奖励
	ServerNotificationFlag_ServerMaintenance              ServerNotificationFlag = 256 // 服务器维护
	ServerNotificationFlag_CannotReceiveMail              ServerNotificationFlag = 512
	ServerNotificationFlag_InventoryFullRewardMail        ServerNotificationFlag = 1024 // 库存满发邮件
	ServerNotificationFlag_CanReceiveClanAttendanceReward ServerNotificationFlag = 2048 // 社团出席奖励
	ServerNotificationFlag_HasClanApplicant               ServerNotificationFlag = 4096 // 有社团申请人
	ServerNotificationFlag_HasFriendRequest               ServerNotificationFlag = 8192 // 有好友申请
	ServerNotificationFlag_CheckConquest                  ServerNotificationFlag = 16384
	ServerNotificationFlag_CanReceiveEliminateRaidReward  ServerNotificationFlag = 32768 // 大决战总分奖励
	ServerNotificationFlag_CanReceiveMultiFloorRaidReward ServerNotificationFlag = 65536
)
