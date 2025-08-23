package model

import (
	"time"
)

// Subscription 用户订阅表
type Subscription struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	UserID           int       `json:"user_id" gorm:"not null"`
	TopicName        string    `json:"topic_name" gorm:"not null;size:100"`
	TopicDescription string    `json:"topic_description" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status           int       `json:"status" gorm:"default:1"` // 1: 活跃, 0: 取消

	// 关联字段
	User         User                  `json:"user" gorm:"foreignKey:UserID"`
	Articles     []SubscriptionArticle `json:"articles" gorm:"foreignKey:SubscriptionID"`
	ArticleCount int                   `json:"article_count" gorm:"-"`
}

// SubscriptionArticle 订阅文章表
type SubscriptionArticle struct {
	ID             int        `json:"id" gorm:"primaryKey"`
	SubscriptionID int        `json:"subscription_id" gorm:"not null"`
	Title          string     `json:"title" gorm:"not null;size:255"`
	Content        string     `json:"content" gorm:"type:text"`
	Author         string     `json:"author" gorm:"size:100"`
	PublishedAt    *time.Time `json:"published_at"`
	ArticleURL     string     `json:"article_url" gorm:"size:500"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Status         int        `json:"status" gorm:"default:1"` // 1: 正常, 0: 删除

	// 关联字段
	Subscription Subscription `json:"subscription" gorm:"foreignKey:SubscriptionID"`
}

// TableName 指定表名
func (Subscription) TableName() string {
	return "subscriptions"
}

func (SubscriptionArticle) TableName() string {
	return "subscription_articles"
}

// GetUserSubscriptions 获取用户的所有订阅
func GetUserSubscriptions(userID int, page, pageSize int) ([]Subscription, int64, error) {
	var subscriptions []Subscription
	var total int64

	// 获取总数（包括已取消的订阅）
	err := DB.Model(&Subscription{}).Where("user_id = ? AND status = 1", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据（包括已取消的订阅）
	offset := (page - 1) * pageSize
	err = DB.Preload("Articles", "status = 1").
		Where("user_id = ? AND status = 1", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&subscriptions).Error

	if err != nil {
		return nil, 0, err
	}

	// 计算每个订阅的文章数量
	for i := range subscriptions {
		subscriptions[i].ArticleCount = len(subscriptions[i].Articles)
	}

	return subscriptions, total, nil
}

// GetSubscriptionByID 根据ID获取订阅
func GetSubscriptionByID(id int) (*Subscription, error) {
	var subscription Subscription
	err := DB.Preload("Articles", "status = 1").
		Where("id = ?", id).
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	subscription.ArticleCount = len(subscription.Articles)
	return &subscription, nil
}

// CreateSubscription 创建订阅
func CreateSubscription(subscription *Subscription) error {
	return DB.Create(subscription).Error
}

// UpdateSubscription 更新订阅
func UpdateSubscription(subscription *Subscription) error {
	return DB.Save(subscription).Error
}

// DeleteSubscription 删除订阅（软删除）
func DeleteSubscription(id int) error {
	return DB.Model(&Subscription{}).Where("id = ?", id).Update("status", 0).Error
}

// CancelSubscription 取消订阅
func CancelSubscription(id, userID int) error {
	return DB.Model(&Subscription{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", 0).Error
}

// CheckSubscriptionExists 检查用户是否已订阅某个主题
func CheckSubscriptionExists(userID int, topicName string) (bool, error) {
	var count int64
	err := DB.Model(&Subscription{}).
		Where("user_id = ? AND topic_name = ? AND status = 1", userID, topicName).
		Count(&count).Error
	return count > 0, err
}

// ReactivateSubscription 重新激活已取消的订阅
func ReactivateSubscription(userID int, topicName string) error {
	return DB.Model(&Subscription{}).
		Where("user_id = ? AND topic_name = ? AND status = 0", userID, topicName).
		Update("status", 1).Error
}

// GetSubscriptionArticles 获取订阅下的文章
func GetSubscriptionArticles(subscriptionID int, page, pageSize int) ([]SubscriptionArticle, int64, error) {
	var articles []SubscriptionArticle
	var total int64

	// 获取总数
	err := DB.Model(&SubscriptionArticle{}).
		Where("subscription_id = ? AND status = 1", subscriptionID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("subscription_id = ? AND status = 1", subscriptionID).
		Order("published_at DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error

	return articles, total, nil
}

// CreateSubscriptionArticle 创建订阅文章
func CreateSubscriptionArticle(article *SubscriptionArticle) error {
	return DB.Create(article).Error
}

// DeleteSubscriptionArticle 删除订阅文章（软删除）
func DeleteSubscriptionArticle(id int) error {
	return DB.Model(&SubscriptionArticle{}).Where("id = ?", id).Update("status", 0).Error
}
