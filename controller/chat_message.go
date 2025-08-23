package controller

import (
	"net/http"
	"strconv"

	"one-api/common"
	"one-api/dto"
	"one-api/model"

	"github.com/gin-gonic/gin"
)

// CreateChatMessage 创建聊天消息
func CreateChatMessage(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	var req dto.CreateChatMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证会话是否存在且属于当前用户
	session := &model.ChatSession{}
	if err := session.GetBySessionId(req.SessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限在此会话中发送消息")
		return
	}

	// 创建消息
	message := &model.ChatMessage{
		SessionId: req.SessionId,
		MessageId: req.MessageId,
		Role:      req.Role,
		Content:   req.Content,
		Tokens:    req.Tokens,
		Cost:      req.Cost,
		Status:    req.Status,
		ErrorMsg:  req.ErrorMsg,
	}

	if err := message.Insert(); err != nil {
		common.ApiErrorMsg(c, "创建消息失败: "+err.Error())
		return
	}

	// 更新会话统计信息
	// 这里需要重新计算会话的消息统计
	messages, err := model.GetSessionAllMessages(req.SessionId)
	if err == nil {
		totalMessages := len(messages)
		totalTokens := 0
		totalCost := 0.0
		for _, msg := range messages {
			totalTokens += msg.Tokens
			totalCost += msg.Cost
		}
		model.UpdateSessionStats(req.SessionId, totalMessages, totalTokens, totalCost)
	}

	// 转换为响应格式
	response := dto.ChatMessageResponse{
		Id:          message.Id,
		SessionId:   message.SessionId,
		MessageId:   message.MessageId,
		Role:        message.Role,
		Content:     message.Content,
		Tokens:      message.Tokens,
		Cost:        message.Cost,
		Status:      message.Status,
		ErrorMsg:    message.ErrorMsg,
		CreatedTime: message.CreatedTime,
		UpdatedTime: message.UpdatedTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "消息创建成功",
		"data":    response,
	})
}

// GetChatMessage 获取消息详情
func GetChatMessage(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	messageId := c.Param("message_id")

	// 获取消息
	message := &model.ChatMessage{}
	if err := message.GetByMessageId(messageId); err != nil {
		common.ApiErrorMsg(c, "消息不存在")
		return
	}

	// 验证权限（通过会话验证）
	session := &model.ChatSession{}
	if err := session.GetBySessionId(message.SessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限访问此消息")
		return
	}

	// 转换为响应格式
	response := dto.ChatMessageResponse{
		Id:          message.Id,
		SessionId:   message.SessionId,
		MessageId:   message.MessageId,
		Role:        message.Role,
		Content:     message.Content,
		Tokens:      message.Tokens,
		Cost:        message.Cost,
		Status:      message.Status,
		ErrorMsg:    message.ErrorMsg,
		CreatedTime: message.CreatedTime,
		UpdatedTime: message.UpdatedTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateChatMessage 更新消息
func UpdateChatMessage(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	messageId := c.Param("message_id")
	var req dto.UpdateChatMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取消息
	message := &model.ChatMessage{}
	if err := message.GetByMessageId(messageId); err != nil {
		common.ApiErrorMsg(c, "消息不存在")
		return
	}

	// 验证权限
	session := &model.ChatSession{}
	if err := session.GetBySessionId(message.SessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限修改此消息")
		return
	}

	// 更新消息
	message.Content = req.Content
	message.Tokens = req.Tokens
	message.Cost = req.Cost
	message.Status = req.Status
	message.ErrorMsg = req.ErrorMsg

	if err := message.Update(); err != nil {
		common.ApiErrorMsg(c, "更新消息失败: "+err.Error())
		return
	}

	// 更新会话统计信息
	messages, err := model.GetSessionAllMessages(message.SessionId)
	if err == nil {
		totalMessages := len(messages)
		totalTokens := 0
		totalCost := 0.0
		for _, msg := range messages {
			totalTokens += msg.Tokens
			totalCost += msg.Cost
		}
		model.UpdateSessionStats(message.SessionId, totalMessages, totalTokens, totalCost)
	}

	// 转换为响应格式
	response := dto.ChatMessageResponse{
		Id:          message.Id,
		SessionId:   message.SessionId,
		MessageId:   message.MessageId,
		Role:        message.Role,
		Content:     message.Content,
		Tokens:      message.Tokens,
		Cost:        message.Cost,
		Status:      message.Status,
		ErrorMsg:    message.ErrorMsg,
		CreatedTime: message.CreatedTime,
		UpdatedTime: message.UpdatedTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "消息更新成功",
		"data":    response,
	})
}

// DeleteChatMessage 删除消息
func DeleteChatMessage(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	messageId := c.Param("message_id")

	// 获取消息
	message := &model.ChatMessage{}
	if err := message.GetByMessageId(messageId); err != nil {
		common.ApiErrorMsg(c, "消息不存在")
		return
	}

	// 验证权限
	session := &model.ChatSession{}
	if err := session.GetBySessionId(message.SessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限删除此消息")
		return
	}

	// 删除消息
	if err := message.Delete(); err != nil {
		common.ApiErrorMsg(c, "删除消息失败: "+err.Error())
		return
	}

	// 更新会话统计信息
	messages, err := model.GetSessionAllMessages(message.SessionId)
	if err == nil {
		totalMessages := len(messages)
		totalTokens := 0
		totalCost := 0.0
		for _, msg := range messages {
			totalTokens += msg.Tokens
			totalCost += msg.Cost
		}
		model.UpdateSessionStats(message.SessionId, totalMessages, totalTokens, totalCost)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "消息删除成功",
	})
}

// GetSessionMessages 获取会话消息列表
func GetSessionMessages(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	sessionId := c.Param("session_id")

	// 验证会话权限
	session := &model.ChatSession{}
	if err := session.GetBySessionId(sessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限访问此会话")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 获取消息列表
	messages, total, err := model.GetSessionMessages(sessionId, page, size)
	if err != nil {
		common.ApiErrorMsg(c, "获取消息列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responseData []dto.ChatMessageResponse
	for _, message := range messages {
		responseData = append(responseData, dto.ChatMessageResponse{
			Id:          message.Id,
			SessionId:   message.SessionId,
			MessageId:   message.MessageId,
			Role:        message.Role,
			Content:     message.Content,
			Tokens:      message.Tokens,
			Cost:        message.Cost,
			Status:      message.Status,
			ErrorMsg:    message.ErrorMsg,
			CreatedTime: message.CreatedTime,
			UpdatedTime: message.UpdatedTime,
		})
	}

	response := dto.ChatMessageListResponse{
		Data:  responseData,
		Total: total,
		Page:  page,
		Size:  size,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetUserMessages 获取用户所有消息
func GetUserMessages(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 获取用户消息列表
	messages, total, err := model.GetUserMessages(userId, page, size)
	if err != nil {
		common.ApiErrorMsg(c, "获取消息列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responseData []dto.ChatMessageResponse
	for _, message := range messages {
		responseData = append(responseData, dto.ChatMessageResponse{
			Id:          message.Id,
			SessionId:   message.SessionId,
			MessageId:   message.MessageId,
			Role:        message.Role,
			Content:     message.Content,
			Tokens:      message.Tokens,
			Cost:        message.Cost,
			Status:      message.Status,
			ErrorMsg:    message.ErrorMsg,
			CreatedTime: message.CreatedTime,
			UpdatedTime: message.UpdatedTime,
		})
	}

	response := dto.ChatMessageListResponse{
		Data:  responseData,
		Total: total,
		Page:  page,
		Size:  size,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// SearchMessages 搜索消息
func SearchMessages(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	keyword := c.Query("keyword")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 搜索消息
	messages, total, err := model.SearchMessages(userId, keyword, page, size)
	if err != nil {
		common.ApiErrorMsg(c, "搜索消息失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responseData []dto.ChatMessageResponse
	for _, message := range messages {
		responseData = append(responseData, dto.ChatMessageResponse{
			Id:          message.Id,
			SessionId:   message.SessionId,
			MessageId:   message.MessageId,
			Role:        message.Role,
			Content:     message.Content,
			Tokens:      message.Tokens,
			Cost:        message.Cost,
			Status:      message.Status,
			ErrorMsg:    message.ErrorMsg,
			CreatedTime: message.CreatedTime,
			UpdatedTime: message.UpdatedTime,
		})
	}

	response := dto.ChatMessageListResponse{
		Data:  responseData,
		Total: total,
		Page:  page,
		Size:  size,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetMessageStats 获取消息统计
func GetMessageStats(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	sessionId := c.Param("session_id")

	// 验证会话权限
	session := &model.ChatSession{}
	if err := session.GetBySessionId(sessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限访问此会话")
		return
	}

	// 获取消息统计
	stats, err := model.GetMessageStats(sessionId)
	if err != nil {
		common.ApiErrorMsg(c, "获取消息统计失败: "+err.Error())
		return
	}

	response := dto.ChatMessageStatsResponse{
		TotalMessages:       int64(stats["total_messages"].(int64)),
		UserMessages:        int64(stats["user_messages"].(int64)),
		AssistantMessages:   int64(stats["assistant_messages"].(int64)),
		SystemMessages:      int64(stats["system_messages"].(int64)),
		TotalTokens:         int64(stats["total_tokens"].(int64)),
		TotalCost:           stats["total_cost"].(float64),
		AvgTokensPerMessage: stats["avg_tokens_per_message"].(float64),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetUserMessageStats 获取用户消息统计
func GetUserMessageStats(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 获取用户消息统计
	stats, err := model.GetUserMessageStats(userId)
	if err != nil {
		common.ApiErrorMsg(c, "获取用户消息统计失败: "+err.Error())
		return
	}

	response := dto.ChatMessageStatsResponse{
		TotalMessages:       int64(stats["total_messages"].(int64)),
		UserMessages:        int64(stats["user_messages"].(int64)),
		AssistantMessages:   int64(stats["assistant_messages"].(int64)),
		SystemMessages:      int64(stats["system_messages"].(int64)),
		TotalTokens:         int64(stats["total_tokens"].(int64)),
		TotalCost:           stats["total_cost"].(float64),
		AvgTokensPerMessage: stats["avg_tokens_per_message"].(float64),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetErrorMessages 获取错误消息
func GetErrorMessages(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 获取错误消息
	messages, total, err := model.GetErrorMessages(userId, page, size)
	if err != nil {
		common.ApiErrorMsg(c, "获取错误消息失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responseData []dto.ChatMessageResponse
	for _, message := range messages {
		responseData = append(responseData, dto.ChatMessageResponse{
			Id:          message.Id,
			SessionId:   message.SessionId,
			MessageId:   message.MessageId,
			Role:        message.Role,
			Content:     message.Content,
			Tokens:      message.Tokens,
			Cost:        message.Cost,
			Status:      message.Status,
			ErrorMsg:    message.ErrorMsg,
			CreatedTime: message.CreatedTime,
			UpdatedTime: message.UpdatedTime,
		})
	}

	response := dto.ChatMessageListResponse{
		Data:  responseData,
		Total: total,
		Page:  page,
		Size:  size,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}




