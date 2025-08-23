package dto

// CreateChatSessionRequest 创建会话请求
type CreateChatSessionRequest struct {
	Topic     string `json:"topic" binding:"required"` // 会话主题
	Model     string `json:"model" binding:"required"` // AI模型
	ChannelId int    `json:"channel_id"`               // 渠道ID
}

// UpdateChatSessionRequest 更新会话请求
type UpdateChatSessionRequest struct {
	Topic string `json:"topic" binding:"required"` // 会话主题
}

// ChatSessionResponse 会话响应
type ChatSessionResponse struct {
	Id              int     `json:"id"`                // 会话ID
	UserId          int     `json:"user_id"`           // 用户ID
	SessionId       string  `json:"session_id"`        // 会话唯一标识
	Topic           string  `json:"topic"`             // 会话主题
	Model           string  `json:"model"`             // AI模型
	TotalMessages   int     `json:"total_messages"`    // 消息总数
	TotalTokens     int     `json:"total_tokens"`      // 总token数
	TotalCost       float64 `json:"total_cost"`        // 总费用
	ChannelId       int     `json:"channel_id"`        // 渠道ID
	Status          int     `json:"status"`            // 状态：1-活跃, 2-已结束, 3-已删除
	CreatedTime     int64   `json:"created_time"`      // 创建时间
	UpdatedTime     int64   `json:"updated_time"`      // 更新时间
	LastMessageTime int64   `json:"last_message_time"` // 最后消息时间
}

// ChatSessionListResponse 会话列表响应
type ChatSessionListResponse struct {
	Data  []ChatSessionResponse `json:"data"`  // 会话列表
	Total int64                 `json:"total"` // 总数
	Page  int                   `json:"page"`  // 当前页
	Size  int                   `json:"size"`  // 每页大小
}

// ChatSessionStatsResponse 会话统计响应
type ChatSessionStatsResponse struct {
	TotalSessions         int64   `json:"total_sessions"`           // 总会话数
	ActiveSessions        int64   `json:"active_sessions"`          // 活跃会话数
	TotalMessages         int64   `json:"total_messages"`           // 总消息数
	TotalTokens           int64   `json:"total_tokens"`             // 总token数
	TotalCost             float64 `json:"total_cost"`               // 总费用
	AvgMessagesPerSession float64 `json:"avg_messages_per_session"` // 平均每会话消息数
}

// ModelStatsResponse 模型统计响应
type ModelStatsResponse struct {
	Model        string  `json:"model"`         // 模型名称
	SessionCount int64   `json:"session_count"` // 会话数
	MessageCount int64   `json:"message_count"` // 消息数
	TokenCount   int64   `json:"token_count"`   // token数
	TotalCost    float64 `json:"total_cost"`    // 总费用
}

// ModelStatsListResponse 模型统计列表响应
type ModelStatsListResponse struct {
	Data []ModelStatsResponse `json:"data"` // 模型统计列表
}


