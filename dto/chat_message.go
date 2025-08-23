package dto

// CreateChatMessageRequest 创建消息请求
type CreateChatMessageRequest struct {
	SessionId string  `json:"session_id" binding:"required"` // 会话ID
	MessageId string  `json:"message_id" binding:"required"` // 消息ID
	Role      string  `json:"role" binding:"required"`       // 角色：user, assistant, system
	Content   string  `json:"content" binding:"required"`    // 消息内容
	Tokens    int     `json:"tokens"`                        // token数
	Cost      float64 `json:"cost"`                          // 费用
	Status    int     `json:"status"`                        // 状态：1-正常, 2-错误
	ErrorMsg  string  `json:"error_msg"`                     // 错误信息
}

// UpdateChatMessageRequest 更新消息请求
type UpdateChatMessageRequest struct {
	Content  string  `json:"content" binding:"required"` // 消息内容
	Tokens   int     `json:"tokens"`                     // token数
	Cost     float64 `json:"cost"`                       // 费用
	Status   int     `json:"status"`                     // 状态：1-正常, 2-错误
	ErrorMsg string  `json:"error_msg"`                  // 错误信息
}

// ChatMessageResponse 消息响应
type ChatMessageResponse struct {
	Id          int     `json:"id"`           // 消息ID
	SessionId   string  `json:"session_id"`   // 会话ID
	MessageId   string  `json:"message_id"`   // 消息唯一标识
	Role        string  `json:"role"`         // 角色：user, assistant, system
	Content     string  `json:"content"`      // 消息内容
	Tokens      int     `json:"tokens"`       // token数
	Cost        float64 `json:"cost"`         // 费用
	Status      int     `json:"status"`       // 状态：1-正常, 2-错误
	ErrorMsg    string  `json:"error_msg"`    // 错误信息
	CreatedTime int64   `json:"created_time"` // 创建时间
	UpdatedTime int64   `json:"updated_time"` // 更新时间
}

// ChatMessageListResponse 消息列表响应
type ChatMessageListResponse struct {
	Data  []ChatMessageResponse `json:"data"`  // 消息列表
	Total int64                 `json:"total"` // 总数
	Page  int                   `json:"page"`  // 当前页
	Size  int                   `json:"size"`  // 每页大小
}

// ChatMessageStatsResponse 消息统计响应
type ChatMessageStatsResponse struct {
	TotalMessages       int64   `json:"total_messages"`         // 总消息数
	UserMessages        int64   `json:"user_messages"`          // 用户消息数
	AssistantMessages   int64   `json:"assistant_messages"`     // 助手消息数
	SystemMessages      int64   `json:"system_messages"`        // 系统消息数
	TotalTokens         int64   `json:"total_tokens"`           // 总token数
	TotalCost           float64 `json:"total_cost"`             // 总费用
	AvgTokensPerMessage float64 `json:"avg_tokens_per_message"` // 平均每消息token数
}

// SendMessageRequest 发送消息请求（用于聊天）
type SendMessageRequest struct {
	SessionId string `json:"session_id" binding:"required"` // 会话ID
	Content   string `json:"content" binding:"required"`    // 消息内容
	Model     string `json:"model"`                         // AI模型（可选，如果不提供则使用会话默认模型）
}

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
	UserMessage      ChatMessageResponse `json:"user_message"`      // 用户消息
	AssistantMessage ChatMessageResponse `json:"assistant_message"` // 助手回复
	Session          ChatSessionResponse `json:"session"`           // 更新后的会话信息
}


