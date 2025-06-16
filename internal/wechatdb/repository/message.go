package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sjzar/chatlog/internal/model"
	"github.com/sjzar/chatlog/pkg/util"

	"github.com/rs/zerolog/log"
)

// GetMessages 实现 Repository 接口的 GetMessages 方法
func (r *Repository) GetMessages(ctx context.Context, startTime, endTime time.Time, talker string, sender string, keyword string, limit, offset int) ([]*model.Message, error) {

	talker, sender = r.parseTalkerAndSender(ctx, talker, sender)
	messages, err := r.ds.GetMessages(ctx, startTime, endTime, talker, sender, keyword, limit, offset)
	if err != nil {
		return nil, err
	}

	// 补充消息信息
	if err := r.EnrichMessages(ctx, messages); err != nil {
		log.Debug().Msgf("EnrichMessages failed: %v", err)
	}

	return messages, nil
}

// EnrichMessages 补充消息的额外信息
func (r *Repository) EnrichMessages(ctx context.Context, messages []*model.Message) error {
	for _, msg := range messages {
		r.enrichMessage(msg)
	}
	return nil
}

// GetChatRoomMemberStats 获取群聊成员发言统计
// GetChatRoomMemberMessageStats 获取群聊成员发言统计
func (r *Repository) GetChatRoomMemberMessageStats(ctx context.Context, chatRoomName string, startTime, endTime time.Time) (*model.GetChatRoomMemberMessageStatsResp, error) {
	// 调用数据源获取原始统计数据
	stats, err := r.ds.GetChatRoomMemberMessageStats(ctx, chatRoomName, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// 补充用户信息
	for _, stat := range stats {
		// 从联系人缓存中获取用户信息
		if contact, ok := r.contactCache[stat.UserName]; ok {
			if contact.NickName != "" {
				stat.NickName = contact.NickName
			}
		}

		// 从群聊缓存中获取显示名称
		if chatRoom, ok := r.chatRoomCache[chatRoomName]; ok {
			if displayName, ok := chatRoom.User2DisplayName[stat.UserName]; ok {
				stat.DisplayName = displayName
			}
		}
	}

	// 计算统计信息
	totalMessages := 0
	activeMembers := len(stats)
	for _, stat := range stats {
		totalMessages += stat.MessageCount
	}

	// 获取群聊信息
	chatRoomDisplayName := chatRoomName
	if chatRoom, ok := r.chatRoomCache[chatRoomName]; ok {
		chatRoomDisplayName = chatRoom.DisplayName()
	}

	// 构建响应
	resp := &model.GetChatRoomMemberMessageStatsResp{
		ChatRoomName:  chatRoomDisplayName,
		TotalMembers:  0, // 需要从群聊信息中获取
		ActiveMembers: activeMembers,
		TotalMessages: totalMessages,
		TimeRange:     fmt.Sprintf("%s - %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
		Stats:         stats,
	}

	// 获取群聊总成员数
	if chatRoom, ok := r.chatRoomCache[chatRoomName]; ok {
		resp.TotalMembers = len(chatRoom.Users)
	}

	return resp, nil
}

// GetChatRoomMemberInviteStats 获取群聊成员邀请统计
func (r *Repository) GetChatRoomMemberInviteStats(ctx context.Context, chatRoomName string) (*model.GetChatRoomMemberInviteStatsResp, error) {
	// 调用数据源获取原始统计数据
	stats, err := r.ds.GetChatRoomMemberInviteStats(ctx, chatRoomName)
	if err != nil {
		return nil, err
	}

	// 补充用户信息
	for _, stat := range stats {
		// 从联系人缓存中获取用户信息
		if contact, ok := r.contactCache[stat.UserName]; ok {
			if contact.NickName != "" {
				stat.NickName = contact.NickName
			}
		}

		// 从群聊缓存中获取显示名称
		if chatRoom, ok := r.chatRoomCache[chatRoomName]; ok {
			if displayName, ok := chatRoom.User2DisplayName[stat.UserName]; ok {
				stat.DisplayName = displayName
			}
		}
	}

	// 获取群聊信息
	chatRoomDisplayName := chatRoomName
	if chatRoom, ok := r.chatRoomCache[chatRoomName]; ok {
		chatRoomDisplayName = chatRoom.DisplayName()
	}

	// 过滤掉邀请数为0的记录
	filteredStats := make([]*model.ChatRoomMemberInviteStats, 0)
	for _, stat := range stats {
		if stat.InviteCount > 0 {
			filteredStats = append(filteredStats, stat)
		}
	}

	// 构建响应
	resp := &model.GetChatRoomMemberInviteStatsResp{
		ChatRoomName: chatRoomDisplayName,
		TotalMembers: 0, // 需要从群聊信息中获取
		TimeRange:    "全部时间", // 邀请统计不依赖时间范围
		Stats:        filteredStats,
	}

	// 获取群聊总成员数
	if chatRoom, ok := r.chatRoomCache[chatRoomName]; ok {
		resp.TotalMembers = len(chatRoom.Users)
	}

	return resp, nil
}

// GetChatRoomMemberStats 为了向后兼容，保留原有方法名
func (r *Repository) GetChatRoomMemberStats(ctx context.Context, chatRoomName string, startTime, endTime time.Time) (*model.GetChatRoomMemberStatsResp, error) {
	return r.GetChatRoomMemberMessageStats(ctx, chatRoomName, startTime, endTime)
}

// enrichMessage 补充单条消息的额外信息
func (r *Repository) enrichMessage(msg *model.Message) {
	// 处理群聊消息
	if msg.IsChatRoom {
		// 补充群聊名称
		if chatRoom, ok := r.chatRoomCache[msg.Talker]; ok {
			msg.TalkerName = chatRoom.DisplayName()

			// 补充发送者在群里的显示名称
			if displayName, ok := chatRoom.User2DisplayName[msg.Sender]; ok {
				msg.SenderName = displayName
			}
		}
	}

	// 如果不是自己发送的消息且还没有显示名称，尝试补充发送者信息
	if msg.SenderName == "" && !msg.IsSelf {
		contact := r.getFullContact(msg.Sender)
		if contact != nil {
			msg.SenderName = contact.DisplayName()
		}
	}
}

func (r *Repository) parseTalkerAndSender(ctx context.Context, talker, sender string) (string, string) {
	displayName2User := make(map[string]string)
	users := make(map[string]bool)

	talkers := util.Str2List(talker, ",")
	if len(talkers) > 0 {
		for i := 0; i < len(talkers); i++ {
			if contact, _ := r.GetContact(ctx, talkers[i]); contact != nil {
				talkers[i] = contact.UserName
			} else if chatRoom, _ := r.GetChatRoom(ctx, talker); chatRoom != nil {
				talkers[i] = chatRoom.Name
			}
		}
		// 获取群聊的用户列表
		for i := 0; i < len(talkers); i++ {
			if chatRoom, _ := r.GetChatRoom(ctx, talkers[i]); chatRoom != nil {
				for user, displayName := range chatRoom.User2DisplayName {
					displayName2User[displayName] = user
				}
				for _, user := range chatRoom.Users {
					users[user.UserName] = true
				}
			}
		}
		talker = strings.Join(talkers, ",")
	}

	senders := util.Str2List(sender, ",")
	if len(senders) > 0 {
		for i := 0; i < len(senders); i++ {
			if user, ok := displayName2User[senders[i]]; ok {
				senders[i] = user
			} else {
				// FIXME 大量群聊用户名称重复，无法直接通过 GetContact 获取 ID，后续再优化
				for user := range users {
					if contact := r.getFullContact(user); contact != nil {
						if contact.DisplayName() == senders[i] {
							senders[i] = user
							break
						}
					}
				}
			}
		}
		sender = strings.Join(senders, ",")
	}

	return talker, sender
}
