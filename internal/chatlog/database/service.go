package database

import (
	"time"

	"github.com/sjzar/chatlog/internal/chatlog/ctx"
	"github.com/sjzar/chatlog/internal/model"
	"github.com/sjzar/chatlog/internal/wechatdb"
)

type Service struct {
	ctx *ctx.Context
	db  *wechatdb.DB
}

func NewService(ctx *ctx.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Start() error {
	db, err := wechatdb.New(s.ctx.WorkDir, s.ctx.Platform, s.ctx.Version)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Service) Stop() error {
	if s.db != nil {
		s.db.Close()
	}
	s.db = nil
	return nil
}

func (s *Service) GetDB() *wechatdb.DB {
	return s.db
}

func (s *Service) GetMessages(start, end time.Time, talker string, sender string, keyword string, limit, offset int) ([]*model.Message, error) {
	return s.db.GetMessages(start, end, talker, sender, keyword, limit, offset)
}

func (s *Service) GetContacts(key string, limit, offset int) (*wechatdb.GetContactsResp, error) {
	return s.db.GetContacts(key, limit, offset)
}

func (s *Service) GetChatRooms(key string, limit, offset int) (*wechatdb.GetChatRoomsResp, error) {
	return s.db.GetChatRooms(key, limit, offset)
}

// GetSession retrieves session information
func (s *Service) GetSessions(key string, limit, offset int) (*wechatdb.GetSessionsResp, error) {
	return s.db.GetSessions(key, limit, offset)
}

func (s *Service) GetMedia(_type string, key string) (*model.Media, error) {
	return s.db.GetMedia(_type, key)
}

// GetChatRoomMemberStats 获取群聊成员发言统计
func (s *Service) GetChatRoomMemberStats(chatRoomName string, start, end time.Time) (*model.GetChatRoomMemberStatsResp, error) {
	return s.db.GetChatRoomMemberStats(chatRoomName, start, end)
}

// GetChatRoomMemberMessageStats 获取群聊成员发言统计
func (s *Service) GetChatRoomMemberMessageStats(chatRoomName string, start, end time.Time) (*model.GetChatRoomMemberMessageStatsResp, error) {
	return s.db.GetChatRoomMemberMessageStats(chatRoomName, start, end)
}

// GetChatRoomMemberInviteStats 获取群聊成员邀请统计
func (s *Service) GetChatRoomMemberInviteStats(chatRoomName string) (*model.GetChatRoomMemberInviteStatsResp, error) {
	return s.db.GetChatRoomMemberInviteStats(chatRoomName)
}

// Close closes the database connection
func (s *Service) Close() {
	// Add cleanup code if needed
	s.db.Close()
}
