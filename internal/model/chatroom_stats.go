package model

// ChatRoomMemberStats 群聊成员发言统计
type ChatRoomMemberStats struct {
	UserName     string `json:"userName"`     // 用户微信ID
	DisplayName  string `json:"displayName"`  // 用户在群里的显示名称
	NickName     string `json:"nickName"`     // 用户昵称
	SmallHeadUrl string `json:"smallHeadUrl"` // 用户头像URL
	MessageCount int    `json:"messageCount"` // 发言数量
}

// GetChatRoomMemberStatsResp 群聊成员统计响应
type GetChatRoomMemberStatsResp struct {
	ChatRoomName  string                 `json:"chatRoomName"`  // 群聊名称
	TotalMembers  int                    `json:"totalMembers"`  // 总成员数
	ActiveMembers int                    `json:"activeMembers"` // 活跃成员数（有发言的成员）
	TotalMessages int                    `json:"totalMessages"` // 总消息数
	TimeRange     string                 `json:"timeRange"`     // 统计时间范围
	Stats         []*ChatRoomMemberStats `json:"stats"`         // 成员统计列表
}
