package model

// ChatRoomMemberMessageStats 群聊成员发言统计
type ChatRoomMemberMessageStats struct {
	UserName     string `json:"userName"`     // 用户微信ID
	DisplayName  string `json:"displayName"`  // 用户在群里的显示名称
	NickName     string `json:"nickName"`     // 用户昵称
	SmallHeadUrl string `json:"smallHeadUrl"` // 用户头像URL
	MessageCount int    `json:"messageCount"` // 发言数量
}

// ChatRoomMemberInviteStats 群聊成员邀请统计
type ChatRoomMemberInviteStats struct {
	UserName     string `json:"userName"`     // 用户微信ID
	DisplayName  string `json:"displayName"`  // 用户在群里的显示名称
	NickName     string `json:"nickName"`     // 用户昵称
	SmallHeadUrl string `json:"smallHeadUrl"` // 用户头像URL
	Inviter      string `json:"inviter"`      // 邀请人
	InviteCount  int    `json:"inviteCount"`  // 邀请人数（该用户邀请了多少人）
}

// GetChatRoomMemberMessageStatsResp 群聊成员发言统计响应
type GetChatRoomMemberMessageStatsResp struct {
	ChatRoomName  string                         `json:"chatRoomName"`  // 群聊名称
	TotalMembers  int                            `json:"totalMembers"`  // 总成员数
	ActiveMembers int                            `json:"activeMembers"` // 活跃成员数（有发言的成员）
	TotalMessages int                            `json:"totalMessages"` // 总消息数
	TimeRange     string                         `json:"timeRange"`     // 统计时间范围
	Stats         []*ChatRoomMemberMessageStats `json:"stats"`         // 成员发言统计列表
}

// GetChatRoomMemberInviteStatsResp 群聊成员邀请统计响应
type GetChatRoomMemberInviteStatsResp struct {
	ChatRoomName  string                        `json:"chatRoomName"`  // 群聊名称
	TotalMembers  int                           `json:"totalMembers"`  // 总成员数
	TimeRange     string                        `json:"timeRange"`     // 统计时间范围
	Stats         []*ChatRoomMemberInviteStats `json:"stats"`         // 成员邀请统计列表
}

// 为了向后兼容，保留原有结构
type ChatRoomMemberStats = ChatRoomMemberMessageStats
type GetChatRoomMemberStatsResp = GetChatRoomMemberMessageStatsResp
