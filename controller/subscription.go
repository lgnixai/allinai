package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"one-api/common"
	"one-api/dto"
	"one-api/model"

	"github.com/gin-gonic/gin"
)

// GetUserSubscriptions 获取用户的订阅列表
func GetUserSubscriptions(c *gin.Context) {
	userID := c.GetInt("id")

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
	userID := c.GetInt("id")

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
	userID := c.GetInt("id")
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
	if req.TopicName != "" {
		subscription.TopicName = req.TopicName
	}
	if req.TopicDescription != "" {
		subscription.TopicDescription = req.TopicDescription
	}
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
	userID := c.GetInt("id")
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
	userID := c.GetInt("id")
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
	userID := c.GetInt("id")
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
	userID := c.GetInt("id")
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
			Summary:        article.Summary,
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

// GetAllSubscriptionArticles 获取当前用户订阅的所有文章
func GetAllSubscriptionArticles(c *gin.Context) {
	// 获取当前用户ID
	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未登录",
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

	// 获取当前用户的订阅文章列表
	articles, total, err := model.GetUserSubscriptionArticles(userID, page, pageSize)
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
			Summary:        article.Summary,
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
	userID := c.GetInt("id")
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
		Summary:        req.Summary,
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

	// 模拟概要模板
	summaries := []string{
		"本文全面介绍了%s技术的基础知识和核心概念，适合初学者和有一定基础的开发者阅读。文章从理论到实践，深入浅出地讲解了%s的关键要点和实际应用场景。",
		"随着技术的快速发展，%s领域涌现出许多新的趋势和突破性进展。本文详细分析了当前%s技术的最新发展方向，为读者提供前瞻性的技术洞察。",
		"对于想要入门%s技术的开发者来说，本文提供了一个系统性的学习路径。从基础概念到高级应用，逐步深入，帮助读者建立完整的知识体系。",
		"通过丰富的实际案例，本文展示了%s技术在不同场景下的应用效果。这些案例涵盖了从简单应用到复杂系统的各个方面，具有很强的实用价值。",
		"基于多年的实践经验，本文总结了%s技术的最佳实践和常见陷阱。这些经验教训可以帮助开发团队避免重复错误，提高开发效率。",
	}

	// 模拟内容模板（更长的文章内容）
	contents := []string{
		`在当今快速发展的技术世界中，%s已经成为了一个不可忽视的重要领域。本文将从多个维度深入探讨%s的核心概念、技术原理以及实际应用。

首先，让我们从基础概念开始。%s技术的核心在于其独特的设计理念和架构模式。这种技术不仅仅是一个工具，更是一种思维方式，它改变了我们解决问题的方法。通过深入理解%s的基本原理，我们可以更好地把握其发展方向。

在实际应用中，%s技术展现出了强大的适应性和扩展性。无论是小型项目还是大型企业级应用，%s都能够提供稳定可靠的解决方案。本文将通过具体的案例分析，展示%s在不同场景下的应用效果。

对于开发者来说，掌握%s技术不仅能够提升个人技能，还能够为团队带来更大的价值。通过系统性的学习和实践，我们可以将%s的优势发挥到极致，创造出更好的产品和服务。

最后，我们还需要关注%s技术的未来发展趋势。随着新技术的不断涌现，%s也在持续演进和完善。了解这些发展趋势，有助于我们做出更好的技术决策和规划。

总的来说，%s技术的学习是一个持续的过程，需要我们保持开放的心态和不断探索的精神。希望通过本文的介绍，能够为读者提供有价值的参考和启发。`,

		`随着技术的不断进步和创新，%s领域正在经历前所未有的变革。本文将深入分析当前%s技术的最新发展趋势，为读者提供前瞻性的技术洞察和行业分析。

在过去的几年里，%s技术取得了显著的突破和进展。这些进展不仅体现在技术本身的改进上，更重要的是在应用场景和商业模式上的创新。新的应用模式不断涌现，为%s技术的发展注入了新的活力。

从技术架构的角度来看，%s正在向着更加模块化、可扩展的方向发展。这种发展趋势使得%s技术能够更好地适应不同规模和复杂度的项目需求。同时，新的架构模式也为开发者提供了更多的选择和灵活性。

在性能优化方面，%s技术也在不断突破传统的限制。通过引入新的算法和优化策略，%s在处理大规模数据和复杂计算任务时表现出了更好的性能。这些优化不仅提升了用户体验，也为企业带来了更大的商业价值。

安全性是%s技术发展中的另一个重要方面。随着网络安全威胁的不断增加，%s技术也在不断加强其安全防护能力。新的安全机制和防护策略不断被引入，使得%s技术能够更好地保护用户数据和系统安全。

未来，%s技术将继续向着更加智能化、自动化的方向发展。人工智能和机器学习的融入将为%s技术带来新的可能性，使其能够更好地理解和适应用户需求。

总的来说，%s技术的发展前景广阔，但也面临着新的挑战和机遇。只有持续关注技术发展趋势，不断学习和创新，我们才能在激烈的竞争中保持优势。`,

		`对于许多初学者来说，%s技术可能看起来复杂且难以理解。然而，通过系统性的学习和正确的学习方法，任何人都能够掌握%s的核心概念和基本技能。本文将为初学者提供一个完整的学习路径和指导。

学习%s技术的第一步是建立正确的基础知识体系。这包括理解%s的基本概念、核心原理和主要特性。虽然这些基础知识可能看起来枯燥，但它们是后续深入学习的重要基础。通过深入理解这些基础概念，我们可以更好地把握%s技术的本质。

实践是学习%s技术的关键环节。仅仅掌握理论知识是不够的，我们需要通过实际的项目练习来巩固和深化所学知识。建议初学者从简单的项目开始，逐步增加项目的复杂度和规模。这样不仅可以积累实践经验，还能够建立学习的信心。

在学习过程中，我们还需要关注%s技术的最佳实践和常见陷阱。这些经验教训可以帮助我们避免重复错误，提高学习效率。同时，我们也应该关注%s技术社区的最新动态和讨论，这有助于我们了解技术的最新发展和应用趋势。

对于有一定基础的开发者来说，深入学习%s技术的高级特性和优化技巧是提升技能的重要途径。这些高级特性往往能够帮助我们解决更复杂的问题，提升系统的性能和稳定性。

最后，我们还需要培养持续学习的能力。技术发展日新月异，%s技术也在不断演进。只有保持开放的心态和不断探索的精神，我们才能够跟上技术发展的步伐，在职业生涯中保持竞争力。

总的来说，学习%s技术是一个循序渐进的过程，需要我们付出时间和努力。但只要我们坚持正确的学习方法，就一定能够掌握这项有价值的技术。`,

		`在实际的项目开发中，%s技术的应用效果往往需要通过具体的案例来验证和展示。本文将通过多个实际案例，详细分析%s技术在不同场景下的应用效果和实现方案。

第一个案例是一个企业级的数据处理系统。在这个项目中，我们使用%s技术来处理和分析大量的业务数据。通过合理的架构设计和优化策略，我们成功地构建了一个高性能、可扩展的数据处理平台。这个案例展示了%s技术在处理大规模数据时的强大能力。

第二个案例是一个移动应用的后端服务。在这个项目中，%s技术被用来构建API服务和数据处理逻辑。通过%s技术的特性，我们实现了高效的请求处理和响应机制，为用户提供了流畅的使用体验。这个案例说明了%s技术在构建现代Web服务中的重要作用。

第三个案例是一个机器学习平台的集成项目。在这个项目中，%s技术被用来处理和分析机器学习模型的输入输出数据。通过%s技术的灵活性和扩展性，我们成功地构建了一个支持多种机器学习算法的统一平台。这个案例展示了%s技术在人工智能领域的应用潜力。

通过这些实际案例的分析，我们可以看到%s技术在不同领域和场景下的广泛应用。每个案例都有其独特的技术挑战和解决方案，这些经验对于理解和应用%s技术具有重要的参考价值。

同时，这些案例也揭示了%s技术在实际应用中的一些注意事项和最佳实践。例如，性能优化、安全性考虑、可维护性设计等方面都需要我们在实际项目中认真考虑和规划。

总的来说，通过实际案例的学习，我们可以更好地理解%s技术的实际应用价值和潜力。这些案例不仅提供了技术参考，也为我们的项目开发提供了宝贵的经验教训。`,

		`基于多年的项目实践和团队协作经验，本文总结了%s技术的最佳实践和常见陷阱。这些经验教训可以帮助开发团队避免重复错误，提高开发效率和项目质量。

在项目规划阶段，我们需要注意%s技术的选型和架构设计。不同的项目需求和技术栈可能需要不同的%s技术方案。因此，在项目开始之前，我们需要充分了解项目需求，评估各种技术方案的优缺点，选择最适合的技术栈。

在开发过程中，代码质量和可维护性是重要的考虑因素。%s技术虽然强大，但如果使用不当，也可能导致代码复杂化和维护困难。因此，我们需要遵循良好的编码规范，注重代码的可读性和可维护性。

性能优化是%s技术应用中的另一个重要方面。虽然%s技术本身具有良好的性能特性，但在实际应用中，我们仍然需要注意性能优化。这包括合理的数据结构设计、算法优化、缓存策略等方面。

安全性是%s技术应用中不可忽视的重要问题。我们需要确保%s技术的使用符合安全规范，防止潜在的安全漏洞和风险。这包括数据加密、访问控制、输入验证等方面的安全措施。

在团队协作中，知识共享和技能培训是提高团队整体能力的重要途径。%s技术的学习和应用需要团队成员具备相应的技能和知识。因此，我们需要建立有效的知识共享机制，定期进行技术培训和经验交流。

测试和质量保证是确保%s技术应用质量的重要环节。我们需要建立完善的测试体系，包括单元测试、集成测试、性能测试等，确保%s技术应用的稳定性和可靠性。

总的来说，%s技术的最佳实践涵盖了项目开发的各个方面，从规划到实施，从开发到维护，都需要我们认真考虑和规划。只有遵循这些最佳实践，我们才能够充分发挥%s技术的优势，创造出高质量的产品和服务。`,
	}

	// 生成3-5篇模拟文章
	articleCount := rand.Intn(3) + 3 // 3-5篇
	now := time.Now()

	for i := 0; i < articleCount; i++ {
		// 随机选择标题、概要和内容模板
		titleTemplate := titles[rand.Intn(len(titles))]
		summaryTemplate := summaries[rand.Intn(len(summaries))]
		contentTemplate := contents[rand.Intn(len(contents))]
		author := authors[rand.Intn(len(authors))]

		// 生成标题、概要和内容
		title := fmt.Sprintf(titleTemplate, topicName)
		summary := fmt.Sprintf(summaryTemplate, topicName, topicName)
		content := fmt.Sprintf(contentTemplate, topicName)

		// 随机发布时间（过去30天内）
		publishedAt := now.AddDate(0, 0, -rand.Intn(30))

		article := &model.SubscriptionArticle{
			SubscriptionID: subscriptionID,
			Title:          title,
			Summary:        summary,
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
