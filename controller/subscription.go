package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"one-api/common"
	"one-api/dto"
	"one-api/model"
)

// GetUserSubscriptions 获取用户的订阅列表
func GetUserSubscriptions(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	// 获取订阅列表
	subscriptions, total, err := model.GetUserSubscriptions(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取订阅列表失败: " + err.Error(),
		})
		return
	}
	
	// 转换为响应格式
	var response dto.SubscriptionListResponse
	response.Total = total
	response.Page = page
	response.PageSize = pageSize
	
	for _, sub := range subscriptions {
		response.Subscriptions = append(response.Subscriptions, dto.SubscriptionResponse{
			ID:               sub.ID,
			UserID:           sub.UserID,
			TopicName:        sub.TopicName,
			TopicDescription: sub.TopicDescription,
			CreatedAt:        sub.CreatedAt,
			UpdatedAt:        sub.UpdatedAt,
			Status:           sub.Status,
			ArticleCount:     sub.ArticleCount,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// CreateSubscription 创建订阅
func CreateSubscription(c *gin.Context) {
	userID := c.GetInt("user_id")
	
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 检查是否已订阅该主题
	exists, err := model.CheckSubscriptionExists(userID, req.TopicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "检查订阅状态失败: " + err.Error(),
		})
		return
	}
	
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "您已经订阅了该主题",
		})
		return
	}

	// 检查是否有已取消的订阅，如果有则重新激活
	err = model.ReactivateSubscription(userID, req.TopicName)
	if err == nil {
		// 重新激活成功，获取订阅ID
		var existingSub model.Subscription
		err = model.DB.Where("user_id = ? AND topic_name = ?", userID, req.TopicName).First(&existingSub).Error
		if err == nil {
			// 生成模拟文章数据
			go func() {
				err := generateMockArticles(existingSub.ID, req.TopicName)
				if err != nil {
					common.SysError("生成模拟文章失败: " + err.Error())
				}
			}()
			
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "订阅重新激活成功",
				"data": gin.H{
					"id": existingSub.ID,
				},
			})
			return
		}
	}
	
	// 创建订阅
	subscription := &model.Subscription{
		UserID:           userID,
		TopicName:        req.TopicName,
		TopicDescription: req.TopicDescription,
		Status:           1,
	}
	
	err = model.CreateSubscription(subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建订阅失败: " + err.Error(),
		})
		return
	}

	// 生成模拟文章数据
	go func() {
		err := generateMockArticles(subscription.ID, req.TopicName)
		if err != nil {
			common.SysError("生成模拟文章失败: " + err.Error())
		}
	}()
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订阅创建成功",
		"data": gin.H{
			"id": subscription.ID,
		},
	})
}

// UpdateSubscription 更新订阅
func UpdateSubscription(c *gin.Context) {
	userID := c.GetInt("user_id")
	subscriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "订阅ID格式错误",
		})
		return
	}
	
	var req dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 获取订阅
	subscription, err := model.GetSubscriptionByID(subscriptionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	
	// 检查权限
	if subscription.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限修改此订阅",
		})
		return
	}
	
	// 更新订阅
	subscription.TopicDescription = req.TopicDescription
	err = model.UpdateSubscription(subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新订阅失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订阅更新成功",
	})
}

// CancelSubscription 取消订阅
func CancelSubscription(c *gin.Context) {
	userID := c.GetInt("user_id")
	subscriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "订阅ID格式错误",
		})
		return
	}
	
	err = model.CancelSubscription(subscriptionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "取消订阅失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订阅已取消",
	})
}

// ReactivateSubscription 重新激活订阅
func ReactivateSubscription(c *gin.Context) {
	userID := c.GetInt("user_id")
	subscriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "订阅ID格式错误",
		})
		return
	}
	
	// 检查权限
	subscription, err := model.GetSubscriptionByID(subscriptionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	
	if subscription.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限操作此订阅",
		})
		return
	}
	
	// 重新激活订阅
	err = model.ReactivateSubscription(userID, subscription.TopicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "重新激活订阅失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订阅重新激活成功",
	})
}

// DeleteSubscription 删除订阅
func DeleteSubscription(c *gin.Context) {
	userID := c.GetInt("user_id")
	subscriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "订阅ID格式错误",
		})
		return
	}
	
	// 检查权限
	subscription, err := model.GetSubscriptionByID(subscriptionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	
	if subscription.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限删除此订阅",
		})
		return
	}
	
	err = model.DeleteSubscription(subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除订阅失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "订阅已删除",
	})
}

// GetSubscriptionArticles 获取订阅下的文章
func GetSubscriptionArticles(c *gin.Context) {
	userID := c.GetInt("user_id")
	subscriptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "订阅ID格式错误",
		})
		return
	}
	
	// 检查权限
	subscription, err := model.GetSubscriptionByID(subscriptionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	
	if subscription.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "无权限查看此订阅的文章",
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
	
	// 获取文章列表
	articles, total, err := model.GetSubscriptionArticles(subscriptionID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取文章列表失败: " + err.Error(),
		})
		return
	}
	
	// 转换为响应格式
	var response dto.SubscriptionArticleListResponse
	response.Total = total
	response.Page = page
	response.PageSize = pageSize
	
	for _, article := range articles {
		response.Articles = append(response.Articles, dto.SubscriptionArticleResponse{
			ID:             article.ID,
			SubscriptionID: article.SubscriptionID,
			Title:          article.Title,
			Content:        article.Content,
			Author:         article.Author,
			PublishedAt:    article.PublishedAt,
			ArticleURL:     article.ArticleURL,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Status:         article.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetAllSubscriptionArticles 获取所有订阅文章
func GetAllSubscriptionArticles(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	
	// 获取文章列表
	articles, total, err := model.GetAllSubscriptionArticles(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取文章列表失败: " + err.Error(),
		})
		return
	}
	
	// 转换为响应格式
	var response dto.SubscriptionArticleListResponse
	response.Total = total
	response.Page = page
	response.PageSize = pageSize
	
	for _, article := range articles {
		response.Articles = append(response.Articles, dto.SubscriptionArticleResponse{
			ID:             article.ID,
			SubscriptionID: article.SubscriptionID,
			Title:          article.Title,
			Content:        article.Content,
			Author:         article.Author,
			PublishedAt:    article.PublishedAt,
			ArticleURL:     article.ArticleURL,
			CreatedAt:      article.CreatedAt,
			UpdatedAt:      article.UpdatedAt,
			Status:         article.Status,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// CreateSubscriptionArticle 创建订阅文章（管理员功能）
func CreateSubscriptionArticle(c *gin.Context) {
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
			"message": "无权限创建文章",
		})
		return
	}
	
	var req dto.CreateSubscriptionArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}
	
	// 检查订阅是否存在
	_, err = model.GetSubscriptionByID(req.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "订阅不存在",
		})
		return
	}
	
	// 创建文章
	article := &model.SubscriptionArticle{
		SubscriptionID: req.SubscriptionID,
		Title:          req.Title,
		Content:        req.Content,
		Author:         req.Author,
		PublishedAt:    &req.PublishedAt,
		ArticleURL:     req.ArticleURL,
		Status:         1,
	}
	
	err = model.CreateSubscriptionArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建文章失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文章创建成功",
		"data": gin.H{
			"id": article.ID,
		},
	})
}

// generateMockArticles 生成模拟文章数据
func generateMockArticles(subscriptionID int, topicName string) error {
	// 模拟文章标题模板
	titles := []string{
		"深入理解%s的核心概念",
		"%s技术的最新发展趋势",
		"从零开始学习%s",
		"%s实战案例分析",
		"%s最佳实践指南",
		"探索%s的无限可能",
		"%s技术深度解析",
		"掌握%s的关键要点",
	}

	// 模拟作者
	authors := []string{
		"技术专家",
		"资深工程师",
		"行业分析师",
		"技术博主",
		"研究学者",
	}

	// 模拟内容模板
	contents := []string{
		"本文深入探讨了%s的相关技术，从基础概念到高级应用，为读者提供全面的学习指南。",
		"随着技术的不断发展，%s领域出现了许多新的趋势和变化，本文将为您详细分析这些发展动态。",
		"对于初学者来说，%s可能看起来复杂，但通过本文的指导，您将能够快速掌握其核心要点。",
		"通过实际案例，本文展示了%s在实际项目中的应用，帮助读者更好地理解其价值。",
		"本文总结了%s的最佳实践，为开发者和技术团队提供实用的指导建议。",
	}

	// 生成3-5篇模拟文章
	articleCount := rand.Intn(3) + 3 // 3-5篇
	now := time.Now()

	for i := 0; i < articleCount; i++ {
		// 随机选择标题和内容模板
		titleTemplate := titles[rand.Intn(len(titles))]
		contentTemplate := contents[rand.Intn(len(contents))]
		author := authors[rand.Intn(len(authors))]

		// 生成标题和内容
		title := fmt.Sprintf(titleTemplate, topicName)
		content := fmt.Sprintf(contentTemplate, topicName)

		// 随机发布时间（过去30天内）
		publishedAt := now.AddDate(0, 0, -rand.Intn(30))

		article := &model.SubscriptionArticle{
			SubscriptionID: subscriptionID,
			Title:          title,
			Content:        content,
			Author:         author,
			PublishedAt:    &publishedAt,
			ArticleURL:     fmt.Sprintf("https://example.com/articles/%d", rand.Intn(10000)),
			Status:         1,
		}

		err := model.CreateSubscriptionArticle(article)
		if err != nil {
			return err
		}
	}

	return nil
}
