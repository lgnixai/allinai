package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"one-api/model"
)

// GetTopics 获取话题列表
func GetTopics(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	topics, total, err := model.GetUserTopics(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取话题列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"topics": topics,
			"total":  total,
		},
	})
}

// CreateTopic 创建话题
func CreateTopic(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req struct {
		TopicName string `json:"topic_name" binding:"required"`
		Model     string `json:"model"`
		ChannelID int    `json:"channel_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Model == "" {
		req.Model = "gpt-3.5-turbo"
	}
	if req.ChannelID == 0 {
		req.ChannelID = 1
	}

	topic := &model.Topic{
		UserID:    userID,
		TopicName: req.TopicName,
		Model:     req.Model,
		ChannelID: req.ChannelID,
		Status:    1,
	}

	err := model.CreateTopic(topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建话题失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "话题创建成功",
		"data":    topic,
	})
}

// DeleteTopic 删除话题
func DeleteTopic(c *gin.Context) {
	userID := c.GetInt("user_id")
	topicID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "话题ID格式错误",
		})
		return
	}

	// 检查权限
	topic, err := model.GetTopicByID(topicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "话题不存在",
		})
		return
	}

	if topic.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限操作此话题",
		})
		return
	}

	err = model.DeleteTopic(topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除话题失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "话题已删除",
	})
}

// GetTopicMessages 获取话题下的消息
func GetTopicMessages(c *gin.Context) {
	userID := c.GetInt("user_id")
	topicID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "话题ID格式错误",
		})
		return
	}

	// 检查权限
	topic, err := model.GetTopicByID(topicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "话题不存在",
		})
		return
	}

	if topic.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限查看此话题",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	messages, total, err := model.GetTopicMessages(topicID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取消息列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"messages": messages,
			"total":    total,
		},
	})
}

// CreateMessage 创建消息
func CreateMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	topicID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "话题ID格式错误",
		})
		return
	}

	// 检查权限
	topic, err := model.GetTopicByID(topicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "话题不存在",
		})
		return
	}

	if topic.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限操作此话题",
		})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
		Role    string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认角色
	if req.Role == "" {
		req.Role = "user"
	}

	// 创建用户消息
	userMessage := &model.Message{
		TopicID: topicID,
		Role:    req.Role,
		Content: req.Content,
		Status:  1,
	}

	err = model.CreateMessage(userMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "发送消息失败: " + err.Error(),
		})
		return
	}

	// 生成AI回复
	aiResponse := model.GenerateAIResponse(req.Content, topic.TopicName)
	
	// 创建AI消息
	aiMessage := &model.Message{
		TopicID: topicID,
		Role:    "assistant",
		Content: aiResponse,
		Status:  1,
	}

	err = model.CreateMessage(aiMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成AI回复失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "消息发送成功",
		"data": gin.H{
			"user_message": userMessage,
			"ai_message":   aiMessage,
		},
	})
}
