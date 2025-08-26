package model

import (
	"time"
)

// Subscription 用户订阅表
type Subscription struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	CreateUserID     int       `json:"create_user_id" gorm:"column:create_user_id;not null"`
	TopicName        string    `json:"topic_name" gorm:"not null;size:100"`
	TopicDescription string    `json:"topic_description" gorm:"type:text"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status           int       `json:"status" gorm:"default:1"` // 1: 活跃, 0: 取消

	// 关联字段
	Articles     []SubscriptionArticle `json:"articles" gorm:"foreignKey:SubscriptionID"`
	ArticleCount int                   `json:"article_count" gorm:"-"`
}

// SubscriptionArticle 订阅文章表
type SubscriptionArticle struct {
	ID             int        `json:"id" gorm:"primaryKey"`
	SubscriptionID int        `json:"subscription_id" gorm:"not null"`
	Title          string     `json:"title" gorm:"not null;size:255"`
	Summary        string     `json:"summary" gorm:"type:text"` // 文章概要
	Content        string     `json:"content" gorm:"type:text"`
	Author         string     `json:"author" gorm:"size:100"`
	PublishedAt    *time.Time `json:"published_at"`
	ArticleURL     string     `json:"article_url" gorm:"size:500"`
	// 新增字段
	KeyPoints     string    `json:"key_points" gorm:"type:text"`                 // 重点提炼
	JournalName   string    `json:"journal_name" gorm:"size:200"`                // 期刊名称
	ReadCount     int       `json:"read_count" gorm:"default:0"`                 // 阅读次数
	CitationCount int       `json:"citation_count" gorm:"default:0"`             // 引用次数
	Rating        float64   `json:"rating" gorm:"default:0.0;type:decimal(3,1)"` // 评分
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status        int       `json:"status" gorm:"default:1"` // 1: 正常, 0: 删除

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
	err := DB.Model(&Subscription{}).Where("create_user_id = ? AND status = 1", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据（包括已取消的订阅）
	offset := (page - 1) * pageSize
	err = DB.Preload("Articles", "status = 1").
		Where("create_user_id = ? AND status = 1", userID).
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
		Where("id = ? AND create_user_id = ?", id, userID).
		Update("status", 0).Error
}

// CheckSubscriptionExists 检查用户是否已订阅某个主题
func CheckSubscriptionExists(userID int, topicName string) (bool, error) {
	var count int64
	err := DB.Model(&Subscription{}).
		Where("create_user_id = ? AND topic_name = ? AND status = 1", userID, topicName).
		Count(&count).Error
	return count > 0, err
}

// ReactivateSubscription 重新激活已取消的订阅
func ReactivateSubscription(userID int, topicName string) error {
	return DB.Model(&Subscription{}).
		Where("create_user_id = ? AND topic_name = ? AND status = 0", userID, topicName).
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

// GetAllSubscriptionArticles 获取所有订阅文章（分页）
func GetAllSubscriptionArticles(page, pageSize int) ([]SubscriptionArticle, int64, error) {
	var articles []SubscriptionArticle
	var total int64

	// 获取总数
	err := DB.Model(&SubscriptionArticle{}).
		Where("status = 1").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("status = 1").
		Order("published_at DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error

	return articles, total, nil
}

// GetAllSubscriptionArticlesWithSubscription 获取所有订阅文章（包含订阅信息）
func GetAllSubscriptionArticlesWithSubscription(page, pageSize int) ([]SubscriptionArticle, int64, error) {
	var articles []SubscriptionArticle
	var total int64

	// 获取总数
	err := DB.Model(&SubscriptionArticle{}).
		Joins("JOIN subscriptions ON subscription_articles.subscription_id = subscriptions.id").
		Where("subscription_articles.status = 1 AND subscriptions.status = 1").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Joins("JOIN subscriptions ON subscription_articles.subscription_id = subscriptions.id").
		Where("subscription_articles.status = 1 AND subscriptions.status = 1").
		Order("subscription_articles.published_at DESC, subscription_articles.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error

	return articles, total, nil
}

// UserSubscription 用户订阅关系表
type UserSubscription struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	UserID         int       `json:"user_id" gorm:"column:user_id;not null"`
	SubscriptionID int       `json:"subscription_id" gorm:"column:subscription_id;not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Status         int       `json:"status" gorm:"default:1"` // 1: 活跃, 0: 取消

	// 关联字段
	User         User         `json:"user" gorm:"foreignKey:UserID"`
	Subscription Subscription `json:"subscription" gorm:"foreignKey:SubscriptionID"`
}

// TableName 指定表名
func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

// CreateUserSubscription 创建用户订阅关系
func CreateUserSubscription(userSubscription *UserSubscription) error {
	return DB.Create(userSubscription).Error
}

// GetUserSubscriptionByUserAndSubscription 根据用户ID和订阅ID获取关系
func GetUserSubscriptionByUserAndSubscription(userID, subscriptionID int) (*UserSubscription, error) {
	var userSubscription UserSubscription
	err := DB.Where("user_id = ? AND subscription_id = ? AND status = 1", userID, subscriptionID).
		First(&userSubscription).Error
	if err != nil {
		return nil, err
	}
	return &userSubscription, nil
}

// GetUserSubscriptionByUserAndSubscriptionAnyStatus 根据用户ID和订阅ID获取关系（任何状态）
func GetUserSubscriptionByUserAndSubscriptionAnyStatus(userID, subscriptionID int) (*UserSubscription, error) {
	var userSubscription UserSubscription
	err := DB.Where("user_id = ? AND subscription_id = ?", userID, subscriptionID).
		First(&userSubscription).Error
	if err != nil {
		return nil, err
	}
	return &userSubscription, nil
}

// GetUserSubscriptionsByUserID 获取用户的所有订阅关系
func GetUserSubscriptionsByUserID(userID int, page, pageSize int) ([]UserSubscription, int64, error) {
	var userSubscriptions []UserSubscription
	var total int64

	// 获取总数
	err := DB.Model(&UserSubscription{}).
		Joins("JOIN subscriptions ON user_subscriptions.subscription_id = subscriptions.id").
		Where("user_subscriptions.user_id = ? AND user_subscriptions.status = 1 AND subscriptions.status = 1", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Preload("Subscription").
		Joins("JOIN subscriptions ON user_subscriptions.subscription_id = subscriptions.id").
		Where("user_subscriptions.user_id = ? AND user_subscriptions.status = 1 AND subscriptions.status = 1", userID).
		Order("user_subscriptions.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&userSubscriptions).Error

	return userSubscriptions, total, nil
}

// CancelUserSubscription 取消用户订阅关系
func CancelUserSubscription(userID, subscriptionID int) error {
	return DB.Model(&UserSubscription{}).
		Where("user_id = ? AND subscription_id = ?", userID, subscriptionID).
		Update("status", 0).Error
}

// ReactivateUserSubscription 重新激活用户订阅关系
func ReactivateUserSubscription(userID, subscriptionID int) error {
	return DB.Model(&UserSubscription{}).
		Where("user_id = ? AND subscription_id = ? AND status = 0", userID, subscriptionID).
		Update("status", 1).Error
}

// GetSubscriptionByTopicName 根据主题名称获取订阅
func GetSubscriptionByTopicName(topicName string) (*Subscription, error) {
	var subscription Subscription
	err := DB.Where("topic_name = ? AND status = 1", topicName).
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CreateSubscriptionWithUserRelation 创建订阅并建立用户关系
func CreateSubscriptionWithUserRelation(userID int, topicName, topicDescription string) (*Subscription, error) {
	// 开启事务
	tx := DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已存在相同主题的订阅
	existingSubscription, err := GetSubscriptionByTopicName(topicName)
	if err == nil {
		// 订阅已存在，检查用户是否已有关系
		existingUserSubscription, err := GetUserSubscriptionByUserAndSubscriptionAnyStatus(userID, existingSubscription.ID)
		if err == nil {
			// 关系已存在，检查状态
			if existingUserSubscription.Status == 1 {
				// 已经是活跃状态，回滚事务
				tx.Rollback()
				return existingSubscription, nil
			} else {
				// 状态为0，重新激活
				err = tx.Model(&UserSubscription{}).
					Where("user_id = ? AND subscription_id = ?", userID, existingSubscription.ID).
					Update("status", 1).Error
				if err != nil {
					tx.Rollback()
					return nil, err
				}
				// 提交事务
				tx.Commit()
				return existingSubscription, nil
			}
		}

		// 创建用户订阅关系
		userSubscription := &UserSubscription{
			UserID:         userID,
			SubscriptionID: existingSubscription.ID,
			Status:         1,
		}
		err = tx.Create(userSubscription).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 提交事务
		tx.Commit()
		return existingSubscription, nil
	}

	// 创建新的订阅
	subscription := &Subscription{
		CreateUserID:     0, // 订阅本身不关联特定用户，只是记录创建者ID
		TopicName:        topicName,
		TopicDescription: topicDescription,
		Status:           1,
	}
	err = tx.Create(subscription).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建用户订阅关系
	userSubscription := &UserSubscription{
		UserID:         userID,
		SubscriptionID: subscription.ID,
		Status:         1,
	}
	err = tx.Create(userSubscription).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	tx.Commit()
	return subscription, nil
}

// GetUserSubscriptionArticles 获取当前用户订阅的所有文章（使用新的关系表）
func GetUserSubscriptionArticles(userID int, page, pageSize int) ([]SubscriptionArticle, int64, error) {
	var articles []SubscriptionArticle
	var total int64

	// 获取总数
	err := DB.Model(&SubscriptionArticle{}).
		Joins("JOIN user_subscriptions ON subscription_articles.subscription_id = user_subscriptions.subscription_id").
		Joins("JOIN subscriptions ON subscription_articles.subscription_id = subscriptions.id").
		Where("subscription_articles.status = 1 AND user_subscriptions.status = 1 AND subscriptions.status = 1 AND user_subscriptions.user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Joins("JOIN user_subscriptions ON subscription_articles.subscription_id = user_subscriptions.subscription_id").
		Joins("JOIN subscriptions ON subscription_articles.subscription_id = subscriptions.id").
		Where("subscription_articles.status = 1 AND user_subscriptions.status = 1 AND subscriptions.status = 1 AND user_subscriptions.user_id = ?", userID).
		Order("subscription_articles.published_at DESC, subscription_articles.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error

	return articles, total, nil
}
