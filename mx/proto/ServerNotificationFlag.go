package proto

type ServerNotificationFlag int32

const (
	ServerNotificationFlag_None                           = 0
	ServerNotificationFlag_NewMailArrived                 = 4
	ServerNotificationFlag_HasUnreadMail                  = 8
	ServerNotificationFlag_NewToastDetected               = 16
	ServerNotificationFlag_CanReceiveArenaDailyReward     = 32
	ServerNotificationFlag_CanReceiveRaidReward           = 64
	ServerNotificationFlag_ServerMaintenance              = 256
	ServerNotificationFlag_CannotReceiveMail              = 512
	ServerNotificationFlag_InventoryFullRewardMail        = 1024
	ServerNotificationFlag_CanReceiveClanAttendanceReward = 2048
	ServerNotificationFlag_HasClanApplicant               = 4096
	ServerNotificationFlag_HasFriendRequest               = 8192
	ServerNotificationFlag_CheckConquest                  = 16384
	ServerNotificationFlag_CanReceiveEliminateRaidReward  = 32768
	ServerNotificationFlag_CanReceiveMultiFloorRaidReward = 65536
)
