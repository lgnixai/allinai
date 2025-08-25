package dto

import "time"

// SystemRecommendationResponse 系统推荐响应
type SystemRecommendationResponse struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Category          string    `json:"category"`
	SubscriptionCount int       `json:"subscription_count"`
	ArticleCount      int       `json:"article_count"`
	Status            int       `json:"status"`
	SortOrder         int       `json:"sort_order"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// SystemRecommendationListResponse 系统推荐列表响应
type SystemRecommendationListResponse struct {
	Recommendations []SystemRecommendationResponse `json:"recommendations"`
	Total           int64                          `json:"total"`
	Page            int                            `json:"page"`
	PageSize        int                            `json:"page_size"`
}

// CreateSystemRecommendationRequest 创建系统推荐请求
type CreateSystemRecommendationRequest struct {
	Title     string `json:"title" binding:"required,max=255"`
	Description string `json:"description" binding:"max=1000"`
	Category  string `json:"category" binding:"max=100"`
	SortOrder int    `json:"sort_order"`
}

// UpdateSystemRecommendationRequest 更新系统推荐请求
type UpdateSystemRecommendationRequest struct {
	Title       string `json:"title" binding:"max=255"`
	Description string `json:"description" binding:"max=1000"`
	Category    string `json:"category" binding:"max=100"`
	Status      int    `json:"status"`
	SortOrder   int    `json:"sort_order"`
}

// WelcomePageResponse 欢迎页面响应
type WelcomePageResponse struct {
	WelcomeMessage string                        `json:"welcome_message"`
	Recommendations []SystemRecommendationResponse `json:"recommendations"`
}

// RecommendationPageResponse 推荐页面响应
type RecommendationPageResponse struct {
	Recommendations []SystemRecommendationResponse `json:"recommendations"`
}
