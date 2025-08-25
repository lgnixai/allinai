package model

import (
	"fmt"
	"math/rand"
	"time"
)

// Topic 话题表
type Topic struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"not null"`
	TopicName string    `json:"topic_name" gorm:"not null;size:100"`
	Model     string    `json:"model" gorm:"size:50;default:'gpt-3.5-turbo'"`
	ChannelID int       `json:"channel_id" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status    int       `json:"status" gorm:"default:1"` // 1: 正常, 0: 删除

	// 关联字段
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Messages     []Message `json:"messages" gorm:"foreignKey:TopicID"`
	MessageCount int       `json:"message_count" gorm:"-"`
}

// Message 消息表
type Message struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	TopicID   int       `json:"topic_id" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;size:20"` // user, assistant
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status    int       `json:"status" gorm:"default:1"` // 1: 正常, 0: 删除
}

// TableName 指定表名
func (Topic) TableName() string {
	return "topics"
}

func (Message) TableName() string {
	return "messages"
}

// GetUserTopics 获取用户的所有话题
func GetUserTopics(userID int, page, pageSize int) ([]Topic, int64, error) {
	var topics []Topic
	var total int64

	// 获取总数
	err := DB.Model(&Topic{}).Where("user_id = ? AND status = 1", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Preload("Messages", "status = 1").
		Where("user_id = ? AND status = 1", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&topics).Error

	if err != nil {
		return nil, 0, err
	}

	// 计算每个话题的消息数量
	for i := range topics {
		topics[i].MessageCount = len(topics[i].Messages)
	}

	return topics, total, nil
}

// GetTopicByID 根据ID获取话题
func GetTopicByID(id int) (*Topic, error) {
	var topic Topic
	err := DB.Preload("Messages", "status = 1").
		Where("id = ? AND status = 1", id).
		First(&topic).Error
	if err != nil {
		return nil, err
	}
	topic.MessageCount = len(topic.Messages)
	return &topic, nil
}

// CreateTopic 创建话题
func CreateTopic(topic *Topic) error {
	return DB.Create(topic).Error
}

// DeleteTopic 删除话题（软删除）
func DeleteTopic(id int) error {
	return DB.Model(&Topic{}).Where("id = ?", id).Update("status", 0).Error
}

// GetTopicMessages 获取话题下的消息
func GetTopicMessages(topicID int, page, pageSize int) ([]Message, int64, error) {
	var messages []Message
	var total int64

	// 获取总数
	err := DB.Model(&Message{}).
		Where("topic_id = ? AND status = 1", topicID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("topic_id = ? AND status = 1", topicID).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	return messages, total, nil
}

// CreateMessage 创建消息
func CreateMessage(message *Message) error {
	return DB.Create(message).Error
}

// GenerateAIResponse 生成AI回复
func GenerateAIResponse(userMessage string, topicName string) string {
	// 这里可以集成真实的AI API，现在返回模拟回复
	responses := []string{
		"这是一个很好的问题！让我来为您详细解答。",
		"根据您的问题，我认为可以从以下几个方面来考虑。",
		"感谢您的提问！这个问题确实很有趣，我的看法是...",
		"您提到的这个观点很有见地，我想补充几点。",
		"这个问题涉及到多个层面，让我逐一为您分析。",
		"非常有趣的问题！让我从专业角度为您解答。",
		"您的问题很有深度，我来为您详细分析一下。",
		"这是一个值得深入探讨的话题，我的观点是...",
		"感谢您的提问，让我为您提供一些见解。",
		"您的问题很有启发性，我认为可以这样理解...",
	}

	// 随机选择回复
	rand.Seed(time.Now().UnixNano())
	response := responses[rand.Intn(len(responses))]

	// 在回复前加上话题名称作为前缀
	if topicName != "" {
		return fmt.Sprintf(`"%s": %s`, topicName, response)
	}

	return response
}
