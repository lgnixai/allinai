package controller

import (
	"net/http"
	"strconv"

	"one-api/common"
	"one-api/dto"
	"one-api/model"

	"github.com/gin-gonic/gin"
)

// GetSystemRecommendations 获取系统推荐列表
func GetSystemRecommendations(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取推荐列表
	recommendations, total, err := model.GetSystemRecommendations(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取推荐列表失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var response dto.SystemRecommendationListResponse
	response.Total = total
	response.Page = page
	response.PageSize = pageSize

	for _, rec := range recommendations {
		response.Recommendations = append(response.Recommendations, dto.SystemRecommendationResponse{
			ID:                rec.ID,
			Title:             rec.Title,
			Description:       rec.Description,
			Category:          rec.Category,
			SubscriptionCount: rec.SubscriptionCount,
			ArticleCount:      rec.ArticleCount,
			Status:            rec.Status,
			SortOrder:         rec.SortOrder,
			CreatedAt:         rec.CreatedAt,
			UpdatedAt:         rec.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetSystemRecommendationByID 根据ID获取系统推荐
func GetSystemRecommendationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID参数",
		})
		return
	}

	recommendation, err := model.GetSystemRecommendationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "推荐主题不存在",
		})
		return
	}

	response := dto.SystemRecommendationResponse{
		ID:                recommendation.ID,
		Title:             recommendation.Title,
		Description:       recommendation.Description,
		Category:          recommendation.Category,
		SubscriptionCount: recommendation.SubscriptionCount,
		ArticleCount:      recommendation.ArticleCount,
		Status:            recommendation.Status,
		SortOrder:         recommendation.SortOrder,
		CreatedAt:         recommendation.CreatedAt,
		UpdatedAt:         recommendation.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// CreateSystemRecommendation 创建系统推荐（管理员功能）
func CreateSystemRecommendation(c *gin.Context) {
	// 检查管理员权限
	userID := c.GetInt("user_id")
	user, err := model.GetUserById(userID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户信息失败",
		})
		return
	}

	if user.Role < common.RoleAdminUser {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限创建系统推荐",
		})
		return
	}

	var req dto.CreateSystemRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 创建推荐
	recommendation := &model.SystemRecommendation{
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		SubscriptionCount: 0,
		ArticleCount:    0,
		Status:          1,
		SortOrder:       req.SortOrder,
	}

	err = model.CreateSystemRecommendation(recommendation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建推荐失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "推荐创建成功",
		"data": gin.H{
			"id": recommendation.ID,
		},
	})
}

// UpdateSystemRecommendation 更新系统推荐（管理员功能）
func UpdateSystemRecommendation(c *gin.Context) {
	// 检查管理员权限
	userID := c.GetInt("user_id")
	user, err := model.GetUserById(userID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户信息失败",
		})
		return
	}

	if user.Role < common.RoleAdminUser {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限更新系统推荐",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID参数",
		})
		return
	}

	var req dto.UpdateSystemRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取现有推荐
	recommendation, err := model.GetSystemRecommendationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "推荐主题不存在",
		})
		return
	}

	// 更新字段
	if req.Title != "" {
		recommendation.Title = req.Title
	}
	if req.Description != "" {
		recommendation.Description = req.Description
	}
	if req.Category != "" {
		recommendation.Category = req.Category
	}
	if req.Status >= 0 {
		recommendation.Status = req.Status
	}
	if req.SortOrder > 0 {
		recommendation.SortOrder = req.SortOrder
	}

	err = model.UpdateSystemRecommendation(recommendation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新推荐失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "推荐更新成功",
	})
}

// DeleteSystemRecommendation 删除系统推荐（管理员功能）
func DeleteSystemRecommendation(c *gin.Context) {
	// 检查管理员权限
	userID := c.GetInt("user_id")
	user, err := model.GetUserById(userID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户信息失败",
		})
		return
	}

	if user.Role < common.RoleAdminUser {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限删除系统推荐",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID参数",
		})
		return
	}

	err = model.DeleteSystemRecommendation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除推荐失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "推荐删除成功",
	})
}

// SearchSystemRecommendations 搜索系统推荐
func SearchSystemRecommendations(c *gin.Context) {
	// 获取搜索关键字
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "搜索关键字不能为空",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 搜索推荐列表
	recommendations, total, err := model.SearchSystemRecommendations(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "搜索推荐失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var response dto.SystemRecommendationListResponse
	response.Total = total
	response.Page = page
	response.PageSize = pageSize

	for _, rec := range recommendations {
		response.Recommendations = append(response.Recommendations, dto.SystemRecommendationResponse{
			ID:                rec.ID,
			Title:             rec.Title,
			Description:       rec.Description,
			Category:          rec.Category,
			SubscriptionCount: rec.SubscriptionCount,
			ArticleCount:      rec.ArticleCount,
			Status:            rec.Status,
			SortOrder:         rec.SortOrder,
			CreatedAt:         rec.CreatedAt,
			UpdatedAt:         rec.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"keyword": keyword,
	})
}

// GetWelcomePage 获取欢迎页面（首次访问）
func GetWelcomePage(c *gin.Context) {
	// 获取4个随机推荐
	recommendations, err := model.GetRandomSystemRecommendations(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取推荐失败: " + err.Error(),
		})
		return
	}

	// 欢迎消息
	welcomeMessage := "Hi, 杨博士,我是 Moyo 安排给你的科研合伙人, 我叫IU。今天是咱们俩第一次见面,为了可以更好的开展后面的工作,给你初步介绍下我现在可以做的事情。因为还不知道你想让我做什么,我根据你的专业帮你选择了几个可能感兴趣的话题。"

	// 转换为响应格式
	var response dto.WelcomePageResponse
	response.WelcomeMessage = welcomeMessage

	for _, rec := range recommendations {
		response.Recommendations = append(response.Recommendations, dto.SystemRecommendationResponse{
			ID:                rec.ID,
			Title:             rec.Title,
			Description:       rec.Description,
			Category:          rec.Category,
			SubscriptionCount: rec.SubscriptionCount,
			ArticleCount:      rec.ArticleCount,
			Status:            rec.Status,
			SortOrder:         rec.SortOrder,
			CreatedAt:         rec.CreatedAt,
			UpdatedAt:         rec.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetRecommendationPage 获取推荐页面（后续访问）
func GetRecommendationPage(c *gin.Context) {
	// 获取4个随机推荐
	recommendations, err := model.GetRandomSystemRecommendations(4)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取推荐失败: " + err.Error(),
		})
		return
	}

	// 转换为响应格式
	var response dto.RecommendationPageResponse

	for _, rec := range recommendations {
		response.Recommendations = append(response.Recommendations, dto.SystemRecommendationResponse{
			ID:                rec.ID,
			Title:             rec.Title,
			Description:       rec.Description,
			Category:          rec.Category,
			SubscriptionCount: rec.SubscriptionCount,
			ArticleCount:      rec.ArticleCount,
			Status:            rec.Status,
			SortOrder:         rec.SortOrder,
			CreatedAt:         rec.CreatedAt,
			UpdatedAt:         rec.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
