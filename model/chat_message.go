package model

import (
	"time"

	"gorm.io/gorm"
)

// ChatMessage 聊天消息表
type ChatMessage struct {
	Id          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionId   string     `json:"session_id" gorm:"type:varchar(64);not null;index"`
	MessageId   string     `json:"message_id" gorm:"type:varchar(64);not null;index"`
	Role        string     `json:"role" gorm:"type:varchar(20);not null"` // user, assistant, system
	Content     string     `json:"content" gorm:"type:text;not null"`
	Tokens      int        `json:"tokens" gorm:"type:int;default:0"`
	Cost        float64    `json:"cost" gorm:"type:real;default:0"`
	Status      int        `json:"status" gorm:"type:int;default:1"` // 1: 正常, 2: 错误
	ErrorMsg    string     `json:"error_msg" gorm:"type:text;default:''"`
	CreatedTime int64      `json:"created_time" gorm:"type:int;not null"`
	UpdatedTime int64      `json:"updated_time" gorm:"type:int;not null"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName 返回表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// BeforeCreate 创建前的钩子
func (m *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedTime == 0 {
		m.CreatedTime = time.Now().Unix()
	}
	if m.UpdatedTime == 0 {
		m.UpdatedTime = time.Now().Unix()
	}
	return nil
}

// BeforeUpdate 更新前的钩子
func (m *ChatMessage) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedTime = time.Now().Unix()
	return nil
}

// Insert 插入消息
func (m *ChatMessage) Insert() error {
	return DB.Create(m).Error
}

// Update 更新消息
func (m *ChatMessage) Update() error {
	return DB.Save(m).Error
}

// Delete 删除消息
func (m *ChatMessage) Delete() error {
	return DB.Delete(m).Error
}

// GetById 根据ID获取消息
func (m *ChatMessage) GetById(id int) error {
	return DB.Where("id = ?", id).First(m).Error
}

// GetByMessageId 根据MessageId获取消息
func (m *ChatMessage) GetByMessageId(messageId string) error {
	return DB.Where("message_id = ?", messageId).First(m).Error
}

// GetSessionMessages 获取会话的所有消息
func GetSessionMessages(sessionId string, page, pageSize int) ([]ChatMessage, int64, error) {
	var messages []ChatMessage
	var total int64

	// 获取总数
	err := DB.Model(&ChatMessage{}).Where("session_id = ? AND deleted_at IS NULL", sessionId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("session_id = ? AND deleted_at IS NULL", sessionId).
		Order("created_time ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// GetSessionAllMessages 获取会话的所有消息（不分页）
func GetSessionAllMessages(sessionId string) ([]ChatMessage, error) {
	var messages []ChatMessage
	err := DB.Where("session_id = ? AND deleted_at IS NULL", sessionId).
		Order("created_time ASC").
		Find(&messages).Error
	return messages, err
}

// DeleteSessionMessages 删除会话的所有消息
func DeleteSessionMessages(sessionId string) error {
	return DB.Where("session_id = ?", sessionId).Delete(&ChatMessage{}).Error
}

// GetUserMessages 获取用户的所有消息
func GetUserMessages(userId int, page, pageSize int) ([]ChatMessage, int64, error) {
	var messages []ChatMessage
	var total int64

	// 通过会话表关联获取用户的消息
	subQuery := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Select("session_id")

	// 获取总数
	err := DB.Model(&ChatMessage{}).Where("session_id IN (?) AND deleted_at IS NULL", subQuery).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("session_id IN (?) AND deleted_at IS NULL", subQuery).
		Order("created_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// SearchMessages 搜索消息内容
func SearchMessages(userId int, keyword string, page, pageSize int) ([]ChatMessage, int64, error) {
	var messages []ChatMessage
	var total int64

	// 通过会话表关联获取用户的消息
	subQuery := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Select("session_id")

	// 构建搜索条件
	query := DB.Where("session_id IN (?) AND deleted_at IS NULL", subQuery)
	if keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	err := query.Model(&ChatMessage{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = query.Order("created_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}

// GetMessageStats 获取消息统计
func GetMessageStats(sessionId string) (map[string]interface{}, error) {
	var stats struct {
		TotalMessages       int64   `json:"total_messages"`
		UserMessages        int64   `json:"user_messages"`
		AssistantMessages   int64   `json:"assistant_messages"`
		SystemMessages      int64   `json:"system_messages"`
		TotalTokens         int64   `json:"total_tokens"`
		TotalCost           float64 `json:"total_cost"`
		AvgTokensPerMessage float64 `json:"avg_tokens_per_message"`
	}

	// 总消息数
	err := DB.Model(&ChatMessage{}).Where("session_id = ? AND deleted_at IS NULL", sessionId).Count(&stats.TotalMessages).Error
	if err != nil {
		return nil, err
	}

	// 各角色消息数
	err = DB.Model(&ChatMessage{}).Where("session_id = ? AND role = 'user' AND deleted_at IS NULL", sessionId).Count(&stats.UserMessages).Error
	if err != nil {
		return nil, err
	}

	err = DB.Model(&ChatMessage{}).Where("session_id = ? AND role = 'assistant' AND deleted_at IS NULL", sessionId).Count(&stats.AssistantMessages).Error
	if err != nil {
		return nil, err
	}

	err = DB.Model(&ChatMessage{}).Where("session_id = ? AND role = 'system' AND deleted_at IS NULL", sessionId).Count(&stats.SystemMessages).Error
	if err != nil {
		return nil, err
	}

	// 总token数和总费用
	err = DB.Model(&ChatMessage{}).Where("session_id = ? AND deleted_at IS NULL", sessionId).
		Select("SUM(tokens) as total_tokens, SUM(cost) as total_cost").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// 计算平均token数
	if stats.TotalMessages > 0 {
		stats.AvgTokensPerMessage = float64(stats.TotalTokens) / float64(stats.TotalMessages)
	}

	return map[string]interface{}{
		"total_messages":         stats.TotalMessages,
		"user_messages":          stats.UserMessages,
		"assistant_messages":     stats.AssistantMessages,
		"system_messages":        stats.SystemMessages,
		"total_tokens":           stats.TotalTokens,
		"total_cost":             stats.TotalCost,
		"avg_tokens_per_message": stats.AvgTokensPerMessage,
	}, nil
}

// GetUserMessageStats 获取用户消息统计
func GetUserMessageStats(userId int) (map[string]interface{}, error) {
	var stats struct {
		TotalMessages       int64   `json:"total_messages"`
		UserMessages        int64   `json:"user_messages"`
		AssistantMessages   int64   `json:"assistant_messages"`
		SystemMessages      int64   `json:"system_messages"`
		TotalTokens         int64   `json:"total_tokens"`
		TotalCost           float64 `json:"total_cost"`
		AvgTokensPerMessage float64 `json:"avg_tokens_per_message"`
	}

	// 通过会话表关联获取用户的消息
	subQuery := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Select("session_id")

	// 总消息数
	err := DB.Model(&ChatMessage{}).Where("session_id IN (?) AND deleted_at IS NULL", subQuery).Count(&stats.TotalMessages).Error
	if err != nil {
		return nil, err
	}

	// 各角色消息数
	err = DB.Model(&ChatMessage{}).Where("session_id IN (?) AND role = 'user' AND deleted_at IS NULL", subQuery).Count(&stats.UserMessages).Error
	if err != nil {
		return nil, err
	}

	err = DB.Model(&ChatMessage{}).Where("session_id IN (?) AND role = 'assistant' AND deleted_at IS NULL", subQuery).Count(&stats.AssistantMessages).Error
	if err != nil {
		return nil, err
	}

	err = DB.Model(&ChatMessage{}).Where("session_id IN (?) AND role = 'system' AND deleted_at IS NULL", subQuery).Count(&stats.SystemMessages).Error
	if err != nil {
		return nil, err
	}

	// 总token数和总费用
	err = DB.Model(&ChatMessage{}).Where("session_id IN (?) AND deleted_at IS NULL", subQuery).
		Select("SUM(tokens) as total_tokens, SUM(cost) as total_cost").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// 计算平均token数
	if stats.TotalMessages > 0 {
		stats.AvgTokensPerMessage = float64(stats.TotalTokens) / float64(stats.TotalMessages)
	}

	return map[string]interface{}{
		"total_messages":         stats.TotalMessages,
		"user_messages":          stats.UserMessages,
		"assistant_messages":     stats.AssistantMessages,
		"system_messages":        stats.SystemMessages,
		"total_tokens":           stats.TotalTokens,
		"total_cost":             stats.TotalCost,
		"avg_tokens_per_message": stats.AvgTokensPerMessage,
	}, nil
}

// GetRecentMessages 获取最近的消息
func GetRecentMessages(userId int, limit int) ([]ChatMessage, error) {
	var messages []ChatMessage

	// 通过会话表关联获取用户的消息
	subQuery := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Select("session_id")

	err := DB.Where("session_id IN (?) AND deleted_at IS NULL", subQuery).
		Order("created_time DESC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

// GetErrorMessages 获取错误消息
func GetErrorMessages(userId int, page, pageSize int) ([]ChatMessage, int64, error) {
	var messages []ChatMessage
	var total int64

	// 通过会话表关联获取用户的消息
	subQuery := DB.Model(&ChatSession{}).Where("user_id = ? AND deleted_at IS NULL", userId).Select("session_id")

	// 获取总数
	err := DB.Model(&ChatMessage{}).Where("session_id IN (?) AND status = 2 AND deleted_at IS NULL", subQuery).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("session_id IN (?) AND status = 2 AND deleted_at IS NULL", subQuery).
		Order("created_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, err
}
