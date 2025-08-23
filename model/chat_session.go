package model

import (
	"time"
)

// ChatSession 聊天会话表
type ChatSession struct {
	Id              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId          int        `json:"user_id" gorm:"type:int;not null;index"`
	SessionId       string     `json:"session_id" gorm:"type:varchar(64);not null;uniqueIndex"`
	Topic           string     `json:"topic" gorm:"type:varchar(255);not null"`
	Model           string     `json:"model" gorm:"type:varchar(64);not null"`
	TotalMessages   int        `json:"total_messages" gorm:"type:int;default:0"`
	TotalTokens     int        `json:"total_tokens" gorm:"type:int;default:0"`
	TotalCost       float64    `json:"total_cost" gorm:"type:real;default:0"`
	ChannelId       int        `json:"channel_id" gorm:"type:int;default:0"`
	Status          int        `json:"status" gorm:"type:int;default:1"` // 1: 活跃, 2: 已结束, 3: 已删除
	CreatedTime     int64      `json:"created_time" gorm:"type:int;not null"`
	UpdatedTime     int64      `json:"updated_time" gorm:"type:int;not null"`
	LastMessageTime int64      `json:"last_message_time" gorm:"type:int;not null"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName 返回表名
func (ChatSession) TableName() string {
	return "chat_sessions"
}

// BeforeCreate 创建前的钩子
func (s *ChatSession) BeforeCreate() error {
	if s.CreatedTime == 0 {
		s.CreatedTime = time.Now().Unix()
	}
	if s.UpdatedTime == 0 {
		s.UpdatedTime = time.Now().Unix()
	}
	if s.LastMessageTime == 0 {
		s.LastMessageTime = time.Now().Unix()
	}
	return nil
}

// BeforeUpdate 更新前的钩子
func (s *ChatSession) BeforeUpdate() error {
	s.UpdatedTime = time.Now().Unix()
	return nil
}

// Insert 插入会话
func (s *ChatSession) Insert() error {
	return DB.Create(s).Error
}

// Update 更新会话
func (s *ChatSession) Update() error {
	return DB.Save(s).Error
}

// Delete 删除会话
func (s *ChatSession) Delete() error {
	return DB.Delete(s).Error
}

// GetById 根据ID获取会话
func (s *ChatSession) GetById(id int) error {
	return DB.Where("id = ?", id).First(s).Error
}

// GetBySessionId 根据SessionId获取会话
func (s *ChatSession) GetBySessionId(sessionId string) error {
	return DB.Where("session_id = ?", sessionId).First(s).Error
}

// GetUserSessions 获取用户的会话列表
func GetUserSessions(userId int, page, pageSize int) ([]ChatSession, int64, error) {
	var sessions []ChatSession
	var total int64

	// 获取总数
	err := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("user_id = ? AND deleted_at IS NULL", userId).
		Order("last_message_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error

	return sessions, total, err
}

// GetUserActiveSessions 获取用户的活跃会话
func GetUserActiveSessions(userId int) ([]ChatSession, error) {
	var sessions []ChatSession
	err := DB.Where("user_id = ? AND status = 1 AND deleted_at IS NULL", userId).
		Order("last_message_time DESC").
		Find(&sessions).Error
	return sessions, err
}

// DeleteUserSession 删除用户的会话
func DeleteUserSession(userId int, sessionId string) error {
	return DB.Where("user_id = ? AND session_id = ?", userId, sessionId).
		Delete(&ChatSession{}).Error
}

// DeleteUserAllSessions 删除用户的所有会话
func DeleteUserAllSessions(userId int) error {
	return DB.Where("user_id = ?", userId).Delete(&ChatSession{}).Error
}

// UpdateSessionStats 更新会话统计信息
func UpdateSessionStats(sessionId string, messageCount, tokenCount int, cost float64) error {
	return DB.Model(&ChatSession{}).
		Where("session_id = ?", sessionId).
		Updates(map[string]interface{}{
			"total_messages":    messageCount,
			"total_tokens":      tokenCount,
			"total_cost":        cost,
			"last_message_time": time.Now().Unix(),
			"updated_time":      time.Now().Unix(),
		}).Error
}

// GetUserSessionStats 获取用户会话统计
func GetUserSessionStats(userId int) (map[string]interface{}, error) {
	var stats struct {
		TotalSessions         int64   `json:"total_sessions"`
		ActiveSessions        int64   `json:"active_sessions"`
		TotalMessages         int64   `json:"total_messages"`
		TotalTokens           int64   `json:"total_tokens"`
		TotalCost             float64 `json:"total_cost"`
		AvgMessagesPerSession float64 `json:"avg_messages_per_session"`
	}

	// 总会话数
	err := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Count(&stats.TotalSessions).Error
	if err != nil {
		return nil, err
	}

	// 活跃会话数
	err = DB.Model(&ChatSession{}).Where("user_id = ? AND status = 1 AND deleted_at IS NULL", userId).Count(&stats.ActiveSessions).Error
	if err != nil {
		return nil, err
	}

	// 总消息数、总token数、总费用
	err = DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).
		Select("SUM(total_messages) as total_messages, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// 计算平均消息数
	if stats.TotalSessions > 0 {
		stats.AvgMessagesPerSession = float64(stats.TotalMessages) / float64(stats.TotalSessions)
	}

	return map[string]interface{}{
		"total_sessions":           stats.TotalSessions,
		"active_sessions":          stats.ActiveSessions,
		"total_messages":           stats.TotalMessages,
		"total_tokens":             stats.TotalTokens,
		"total_cost":               stats.TotalCost,
		"avg_messages_per_session": stats.AvgMessagesPerSession,
	}, nil
}

// SearchUserSessions 搜索用户的会话
func SearchUserSessions(userId int, keyword string, page, pageSize int) ([]ChatSession, int64, error) {
	var sessions []ChatSession
	var total int64

	// 构建搜索条件
	query := DB.Where("user_id = ? AND deleted_at IS NULL", userId)
	if keyword != "" {
		query = query.Where("topic LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	err := query.Model(&ChatSession{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = query.Order("last_message_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error

	return sessions, total, err
}

// GetModelStats 获取模型使用统计
func GetModelStats(userId int) ([]map[string]interface{}, error) {
	var stats []struct {
		Model        string  `json:"model"`
		SessionCount int64   `json:"session_count"`
		MessageCount int64   `json:"message_count"`
		TokenCount   int64   `json:"token_count"`
		TotalCost    float64 `json:"total_cost"`
	}

	err := DB.Model(&ChatSession{}).
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Select("model, COUNT(*) as session_count, SUM(total_messages) as message_count, SUM(total_tokens) as token_count, SUM(total_cost) as total_cost").
		Group("model").
		Order("total_cost DESC").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(stats))
	for i, stat := range stats {
		result[i] = map[string]interface{}{
			"model":         stat.Model,
			"session_count": stat.SessionCount,
			"message_count": stat.MessageCount,
			"token_count":   stat.TokenCount,
			"total_cost":    stat.TotalCost,
		}
	}

	return result, nil
}


