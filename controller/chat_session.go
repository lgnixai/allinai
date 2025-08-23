package controller

import (
	"net/http"
	"strconv"

	"one-api/common"
	"one-api/dto"
	"one-api/model"

	"github.com/gin-gonic/gin"
)

// CreateChatSession 创建聊天会话
func CreateChatSession(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	var req dto.CreateChatSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "请求参数错误: "+err.Error())
		return
	}

	// 生成会话ID
	sessionId := common.GenerateUUID()

	// 创建会话
	session := &model.ChatSession{
		UserId:    userId,
		SessionId: sessionId,
		Topic:     req.Topic,
		Model:     req.Model,
		ChannelId: req.ChannelId,
		Status:    1, // 活跃状态
	}

	if err := session.Insert(); err != nil {
		common.ApiErrorMsg(c, "创建会话失败: "+err.Error())
		return
	}

	// 转换为响应格式
	response := dto.ChatSessionResponse{
		Id:              session.Id,
		UserId:          session.UserId,
		SessionId:       session.SessionId,
		Topic:           session.Topic,
		Model:           session.Model,
		TotalMessages:   session.TotalMessages,
		TotalTokens:     session.TotalTokens,
		TotalCost:       session.TotalCost,
		ChannelId:       session.ChannelId,
		Status:          session.Status,
		CreatedTime:     session.CreatedTime,
		UpdatedTime:     session.UpdatedTime,
		LastMessageTime: session.LastMessageTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "会话创建成功",
		"data":    response,
	})
}

// GetChatSession 获取会话详情
func GetChatSession(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	sessionId := c.Param("session_id")

	// 获取会话
	session := &model.ChatSession{}
	if err := session.GetBySessionId(sessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}

	// 检查权限
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限访问此会话")
		return
	}

	// 转换为响应格式
	response := dto.ChatSessionResponse{
		Id:              session.Id,
		UserId:          session.UserId,
		SessionId:       session.SessionId,
		Topic:           session.Topic,
		Model:           session.Model,
		TotalMessages:   session.TotalMessages,
		TotalTokens:     session.TotalTokens,
		TotalCost:       session.TotalCost,
		ChannelId:       session.ChannelId,
		Status:          session.Status,
		CreatedTime:     session.CreatedTime,
		UpdatedTime:     session.UpdatedTime,
		LastMessageTime: session.LastMessageTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// UpdateChatSession 更新会话
func UpdateChatSession(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	sessionId := c.Param("session_id")
	var req dto.UpdateChatSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取会话
	session := &model.ChatSession{}
	if err := session.GetBySessionId(sessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}

	// 检查权限
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限修改此会话")
		return
	}

	// 更新会话
	session.Topic = req.Topic
	if err := session.Update(); err != nil {
		common.ApiErrorMsg(c, "更新会话失败: "+err.Error())
		return
	}

	// 转换为响应格式
	response := dto.ChatSessionResponse{
		Id:              session.Id,
		UserId:          session.UserId,
		SessionId:       session.SessionId,
		Topic:           session.Topic,
		Model:           session.Model,
		TotalMessages:   session.TotalMessages,
		TotalTokens:     session.TotalTokens,
		TotalCost:       session.TotalCost,
		ChannelId:       session.ChannelId,
		Status:          session.Status,
		CreatedTime:     session.CreatedTime,
		UpdatedTime:     session.UpdatedTime,
		LastMessageTime: session.LastMessageTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "会话更新成功",
		"data":    response,
	})
}

// DeleteChatSession 删除会话
func DeleteChatSession(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	sessionId := c.Param("session_id")

	// 获取会话
	session := &model.ChatSession{}
	if err := session.GetBySessionId(sessionId); err != nil {
		common.ApiErrorMsg(c, "会话不存在")
		return
	}

	// 检查权限
	if session.UserId != userId {
		common.ApiErrorMsg(c, "无权限删除此会话")
		return
	}

	// 删除会话及其所有消息
	if err := model.DeleteUserSession(userId, sessionId); err != nil {
		common.ApiErrorMsg(c, "删除会话失败: "+err.Error())
		return
	}

	// 删除会话的所有消息
	if err := model.DeleteSessionMessages(sessionId); err != nil {
		common.ApiErrorMsg(c, "删除会话消息失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "会话删除成功",
	})
}

// GetUserSessions 获取用户会话列表
func GetUserSessions(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取会话列表
	sessions, total, err := model.GetUserSessions(userId, page, size)
	if err != nil {
		common.ApiErrorMsg(c, "获取会话列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responseData []dto.ChatSessionResponse
	for _, session := range sessions {
		responseData = append(responseData, dto.ChatSessionResponse{
			Id:              session.Id,
			UserId:          session.UserId,
			SessionId:       session.SessionId,
			Topic:           session.Topic,
			Model:           session.Model,
			TotalMessages:   session.TotalMessages,
			TotalTokens:     session.TotalTokens,
			TotalCost:       session.TotalCost,
			ChannelId:       session.ChannelId,
			Status:          session.Status,
			CreatedTime:     session.CreatedTime,
			UpdatedTime:     session.UpdatedTime,
			LastMessageTime: session.LastMessageTime,
		})
	}

	response := dto.ChatSessionListResponse{
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

// GetUserSessionStats 获取用户会话统计
func GetUserSessionStats(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 获取统计信息
	stats, err := model.GetUserSessionStats(userId)
	if err != nil {
		common.ApiErrorMsg(c, "获取统计信息失败: "+err.Error())
		return
	}

	response := dto.ChatSessionStatsResponse{
		TotalSessions:         int64(stats["total_sessions"].(int64)),
		ActiveSessions:        int64(stats["active_sessions"].(int64)),
		TotalMessages:         int64(stats["total_messages"].(int64)),
		TotalTokens:           int64(stats["total_tokens"].(int64)),
		TotalCost:             stats["total_cost"].(float64),
		AvgMessagesPerSession: stats["avg_messages_per_session"].(float64),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetModelStats 获取模型使用统计
func GetModelStats(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 获取模型统计
	stats, err := model.GetModelStats(userId)
	if err != nil {
		common.ApiErrorMsg(c, "获取模型统计失败: "+err.Error())
		return
	}

	var responseData []dto.ModelStatsResponse
	for _, stat := range stats {
		responseData = append(responseData, dto.ModelStatsResponse{
			Model:        stat["model"].(string),
			SessionCount: int64(stat["session_count"].(int64)),
			MessageCount: int64(stat["message_count"].(int64)),
			TokenCount:   int64(stat["token_count"].(int64)),
			TotalCost:    stat["total_cost"].(float64),
		})
	}

	response := dto.ModelStatsListResponse{
		Data: responseData,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// SearchUserSessions 搜索用户会话
func SearchUserSessions(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")
	keyword := c.Query("keyword")

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 搜索会话
	sessions, total, err := model.SearchUserSessions(userId, keyword, page, size)
	if err != nil {
		common.ApiErrorMsg(c, "搜索会话失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responseData []dto.ChatSessionResponse
	for _, session := range sessions {
		responseData = append(responseData, dto.ChatSessionResponse{
			Id:              session.Id,
			UserId:          session.UserId,
			SessionId:       session.SessionId,
			Topic:           session.Topic,
			Model:           session.Model,
			TotalMessages:   session.TotalMessages,
			TotalTokens:     session.TotalTokens,
			TotalCost:       session.TotalCost,
			ChannelId:       session.ChannelId,
			Status:          session.Status,
			CreatedTime:     session.CreatedTime,
			UpdatedTime:     session.UpdatedTime,
			LastMessageTime: session.LastMessageTime,
		})
	}

	response := dto.ChatSessionListResponse{
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

// DeleteUserAllSessions 删除用户所有会话
func DeleteUserAllSessions(c *gin.Context) {
	userId := common.GetContextKeyInt(c, "id")

	// 删除用户所有会话
	if err := model.DeleteUserAllSessions(userId); err != nil {
		common.ApiErrorMsg(c, "删除所有会话失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "所有会话删除成功",
	})
}


