package dto

import "time"

// CreateSubscriptionRequest 创建订阅请求
type CreateSubscriptionRequest struct {
	TopicName        string `json:"topic_name" binding:"required,max=100"`
	TopicDescription string `json:"topic_description" binding:"max=1000"`
}

// UpdateSubscriptionRequest 更新订阅请求
type UpdateSubscriptionRequest struct {
	TopicName        string `json:"topic_name" binding:"max=100"`
	TopicDescription string `json:"topic_description" binding:"max=1000"`
}

// SubscriptionResponse 订阅响应
type SubscriptionResponse struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	TopicName        string    `json:"topic_name"`
	TopicDescription string    `json:"topic_description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Status           int       `json:"status"`
	ArticleCount     int       `json:"article_count"`
}

// SubscriptionListResponse 订阅列表响应
type SubscriptionListResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
}

// CreateSubscriptionArticleRequest 创建订阅文章请求
type CreateSubscriptionArticleRequest struct {
	SubscriptionID int       `json:"subscription_id" binding:"required"`
	Title          string    `json:"title" binding:"required,max=255"`
	Summary        string    `json:"summary" binding:"max=1000"` // 文章概要
	Content        string    `json:"content" binding:"max=10000"`
	Author         string    `json:"author" binding:"max=100"`
	PublishedAt    time.Time `json:"published_at"`
	ArticleURL     string    `json:"article_url" binding:"max=500"`
}

// SubscriptionArticleResponse 订阅文章响应
type SubscriptionArticleResponse struct {
	ID             int        `json:"id"`
	SubscriptionID int        `json:"subscription_id"`
	Title          string     `json:"title"`
	Summary        string     `json:"summary"` // 文章概要
	Content        string     `json:"content"`
	Author         string     `json:"author"`
	PublishedAt    *time.Time `json:"published_at"`
	ArticleURL     string     `json:"article_url"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Status         int        `json:"status"`
}

// SubscriptionArticleListResponse 订阅文章列表响应
type SubscriptionArticleListResponse struct {
	Articles []SubscriptionArticleResponse `json:"articles"`
	Total    int64                         `json:"total"`
	Page     int                           `json:"page"`
	PageSize int                           `json:"page_size"`
}
