package router

import (
	"one-api/controller"
	"one-api/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.CORS())
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	apiRouter.Use(middleware.GlobalAPIRateLimit())
	{
		apiRouter.GET("/setup", controller.GetSetup)
		apiRouter.POST("/setup", controller.PostSetup)
		apiRouter.GET("/status", controller.GetStatus)
		apiRouter.GET("/uptime/status", controller.GetUptimeKumaStatus)
		apiRouter.GET("/models", middleware.UserAuth(), controller.DashboardListModels)
		apiRouter.GET("/status/test", middleware.AdminAuth(), controller.TestStatus)
		apiRouter.GET("/notice", controller.GetNotice)
		apiRouter.GET("/about", controller.GetAbout)
		//apiRouter.GET("/midjourney", controller.GetMidjourney)
		apiRouter.GET("/home_page_content", controller.GetHomePageContent)
		apiRouter.GET("/pricing", middleware.TryUserAuth(), controller.GetPricing)
		apiRouter.GET("/verification", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendEmailVerification)
		apiRouter.GET("/phone_verification", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendPhoneVerification)
		//apiRouter.GET("/reset_password", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendPasswordResetPhone)
		//apiRouter.POST("/user/reset", middleware.CriticalRateLimit(), controller.ResetPassword)
		//apiRouter.POST("/user/verify_reset_code", middleware.CriticalRateLimit(), controller.VerifyResetCode)
		//apiRouter.POST("/user/reset_password", middleware.CriticalRateLimit(), controller.ResetPasswordWithNewPassword)
		//apiRouter.GET("/oauth/github", middleware.CriticalRateLimit(), controller.GitHubOAuth)
		//apiRouter.GET("/oauth/oidc", middleware.CriticalRateLimit(), controller.OidcAuth)
		//apiRouter.GET("/oauth/linuxdo", middleware.CriticalRateLimit(), controller.LinuxdoOAuth)
		//apiRouter.GET("/oauth/state", middleware.CriticalRateLimit(), controller.GenerateOAuthCode)
		//apiRouter.GET("/oauth/wechat", middleware.CriticalRateLimit(), controller.WeChatAuth)
		//apiRouter.GET("/oauth/wechat/bind", middleware.CriticalRateLimit(), controller.WeChatBind)
		//apiRouter.GET("/oauth/email/bind", middleware.CriticalRateLimit(), controller.EmailBind)
		//apiRouter.GET("/oauth/telegram/login", middleware.CriticalRateLimit(), controller.TelegramLogin)
		//apiRouter.GET("/oauth/telegram/bind", middleware.CriticalRateLimit(), controller.TelegramBind)
		apiRouter.GET("/ratio_config", middleware.CriticalRateLimit(), controller.GetRatioConfig)

		apiRouter.POST("/stripe/webhook", controller.StripeWebhook)

		userRoute := apiRouter.Group("/user")
		{
			userRoute.POST("/register", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.Register)
			userRoute.POST("/login", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.Login)
			//userRoute.POST("/tokenlog", middleware.CriticalRateLimit(), controller.TokenLog)
			userRoute.GET("/logout", controller.Logout)
			userRoute.GET("/epay/notify", controller.EpayNotify)
			userRoute.GET("/groups", controller.GetUserGroups)

			selfRoute := userRoute.Group("/")
			selfRoute.Use(middleware.UserAuth())
			{
				selfRoute.GET("/self/groups", controller.GetUserGroups)
				selfRoute.GET("/self", controller.GetSelf)
				selfRoute.GET("/models", controller.GetUserModels)
				selfRoute.PUT("/self", controller.UpdateSelf)
				selfRoute.DELETE("/self", controller.DeleteSelf)
				selfRoute.GET("/token", controller.GenerateAccessToken)
				selfRoute.GET("/aff", controller.GetAffCode)
				selfRoute.POST("/topup", middleware.CriticalRateLimit(), controller.TopUp)
				selfRoute.POST("/pay", middleware.CriticalRateLimit(), controller.RequestEpay)
				selfRoute.POST("/amount", controller.RequestAmount)
				selfRoute.POST("/stripe/pay", middleware.CriticalRateLimit(), controller.RequestStripePay)
				selfRoute.POST("/stripe/amount", controller.RequestStripeAmount)
				selfRoute.POST("/aff_transfer", controller.TransferAffQuota)
				selfRoute.PUT("/setting", controller.UpdateUserSetting)
			}

			// 聊天会话相关路由
			chatSessionRoute := apiRouter.Group("/chat_sessions")
			chatSessionRoute.Use(middleware.UserAuth())
			{
				chatSessionRoute.POST("/", controller.CreateChatSession)               // 创建会话
				chatSessionRoute.GET("/:session_id", controller.GetChatSession)        // 获取会话详情
				chatSessionRoute.PUT("/:session_id", controller.UpdateChatSession)     // 更新会话
				chatSessionRoute.DELETE("/:session_id", controller.DeleteChatSession)  // 删除会话
				chatSessionRoute.GET("/user/sessions", controller.GetUserSessions)     // 获取用户会话列表
				chatSessionRoute.GET("/user/stats", controller.GetUserSessionStats)    // 获取用户会话统计
				chatSessionRoute.GET("/user/model_stats", controller.GetModelStats)    // 获取模型使用统计
				chatSessionRoute.GET("/search", controller.SearchUserSessions)         // 搜索会话
				chatSessionRoute.DELETE("/user/all", controller.DeleteUserAllSessions) // 删除用户所有会话
			}

			// 聊天消息相关路由
			chatMessageRoute := apiRouter.Group("/chat_messages")
			chatMessageRoute.Use(middleware.UserAuth())
			{
				chatMessageRoute.POST("/", controller.CreateChatMessage)                       // 创建消息
				chatMessageRoute.GET("/:message_id", controller.GetChatMessage)                // 获取消息详情
				chatMessageRoute.PUT("/:message_id", controller.UpdateChatMessage)             // 更新消息
				chatMessageRoute.DELETE("/:message_id", controller.DeleteChatMessage)          // 删除消息
				chatMessageRoute.GET("/session/:session_id", controller.GetSessionMessages)    // 获取会话消息列表
				chatMessageRoute.GET("/user/messages", controller.GetUserMessages)             // 获取用户所有消息
				chatMessageRoute.GET("/search", controller.SearchMessages)                     // 搜索消息
				chatMessageRoute.GET("/session/:session_id/stats", controller.GetMessageStats) // 获取会话消息统计
				chatMessageRoute.GET("/user/stats", controller.GetUserMessageStats)            // 获取用户消息统计
				chatMessageRoute.GET("/user/errors", controller.GetErrorMessages)              // 获取错误消息
			}

			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.AdminAuth())
			{
				adminRoute.GET("/", controller.GetAllUsers)
				adminRoute.GET("/search", controller.SearchUsers)
				adminRoute.GET("/:id", controller.GetUser)
				adminRoute.POST("/", controller.CreateUser)
				adminRoute.POST("/manage", controller.ManageUser)
				adminRoute.PUT("/", controller.UpdateUser)
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}
		optionRoute := apiRouter.Group("/option")
		optionRoute.Use(middleware.RootAuth())
		{
			optionRoute.GET("/", controller.GetOptions)
			optionRoute.PUT("/", controller.UpdateOption)
			optionRoute.POST("/rest_model_ratio", controller.ResetModelRatio)
			optionRoute.POST("/migrate_console_setting", controller.MigrateConsoleSetting) // 用于迁移检测的旧键，下个版本会删除
		}
		ratioSyncRoute := apiRouter.Group("/ratio_sync")
		ratioSyncRoute.Use(middleware.RootAuth())
		{
			ratioSyncRoute.GET("/channels", controller.GetSyncableChannels)
			ratioSyncRoute.POST("/fetch", controller.FetchUpstreamRatios)
		}
		channelRoute := apiRouter.Group("/channel")
		channelRoute.Use(middleware.AdminAuth())
		{
			channelRoute.GET("/", controller.GetAllChannels)
			channelRoute.GET("/search", controller.SearchChannels)
			channelRoute.GET("/models", controller.ChannelListModels)
			channelRoute.GET("/models_enabled", controller.EnabledListModels)
			channelRoute.GET("/:id", controller.GetChannel)
			channelRoute.GET("/test", controller.TestAllChannels)
			channelRoute.GET("/test/:id", controller.TestChannel)
			channelRoute.GET("/update_balance", controller.UpdateAllChannelsBalance)
			channelRoute.GET("/update_balance/:id", controller.UpdateChannelBalance)
			channelRoute.POST("/", controller.AddChannel)
			channelRoute.PUT("/", controller.UpdateChannel)
			channelRoute.DELETE("/disabled", controller.DeleteDisabledChannel)
			channelRoute.POST("/tag/disabled", controller.DisableTagChannels)
			channelRoute.POST("/tag/enabled", controller.EnableTagChannels)
			channelRoute.PUT("/tag", controller.EditTagChannels)
			channelRoute.DELETE("/:id", controller.DeleteChannel)
			channelRoute.POST("/batch", controller.DeleteChannelBatch)
			channelRoute.POST("/fix", controller.FixChannelsAbilities)
			channelRoute.GET("/fetch_models/:id", controller.FetchUpstreamModels)
			channelRoute.POST("/fetch_models", controller.FetchModels)
			channelRoute.POST("/batch/tag", controller.BatchSetChannelTag)
			channelRoute.GET("/tag/models", controller.GetTagModels)
			channelRoute.POST("/copy/:id", controller.CopyChannel)
		}
		tokenRoute := apiRouter.Group("/token")
		tokenRoute.Use(middleware.UserAuth())
		{
			tokenRoute.GET("/", controller.GetAllTokens)
			tokenRoute.GET("/search", controller.SearchTokens)
			tokenRoute.GET("/:id", controller.GetToken)
			tokenRoute.POST("/", controller.AddToken)
			tokenRoute.PUT("/", controller.UpdateToken)
			tokenRoute.DELETE("/:id", controller.DeleteToken)
			tokenRoute.POST("/batch", controller.DeleteTokenBatch)
		}
		redemptionRoute := apiRouter.Group("/redemption")
		redemptionRoute.Use(middleware.AdminAuth())
		{
			redemptionRoute.GET("/", controller.GetAllRedemptions)
			redemptionRoute.GET("/search", controller.SearchRedemptions)
			redemptionRoute.GET("/:id", controller.GetRedemption)
			redemptionRoute.POST("/", controller.AddRedemption)
			redemptionRoute.PUT("/", controller.UpdateRedemption)
			redemptionRoute.DELETE("/invalid", controller.DeleteInvalidRedemption)
			redemptionRoute.DELETE("/:id", controller.DeleteRedemption)
		}
		logRoute := apiRouter.Group("/log")
		logRoute.GET("/", middleware.AdminAuth(), controller.GetAllLogs)
		logRoute.DELETE("/", middleware.AdminAuth(), controller.DeleteHistoryLogs)
		logRoute.GET("/stat", middleware.AdminAuth(), controller.GetLogsStat)
		logRoute.GET("/self/stat", middleware.UserAuth(), controller.GetLogsSelfStat)
		logRoute.GET("/search", middleware.AdminAuth(), controller.SearchAllLogs)
		logRoute.GET("/self", middleware.UserAuth(), controller.GetUserLogs)
		logRoute.GET("/self/search", middleware.UserAuth(), controller.SearchUserLogs)

		dataRoute := apiRouter.Group("/data")
		dataRoute.GET("/", middleware.AdminAuth(), controller.GetAllQuotaDates)
		dataRoute.GET("/self", middleware.UserAuth(), controller.GetUserQuotaDates)

		logRoute.Use(middleware.CORS())
		{
			logRoute.GET("/token", controller.GetLogByKey)

		}
		groupRoute := apiRouter.Group("/group")
		groupRoute.Use(middleware.AdminAuth())
		{
			groupRoute.GET("/", controller.GetGroups)
		}
		mjRoute := apiRouter.Group("/mj")
		mjRoute.GET("/self", middleware.UserAuth(), controller.GetUserMidjourney)
		mjRoute.GET("/", middleware.AdminAuth(), controller.GetAllMidjourney)

		taskRoute := apiRouter.Group("/task")
		{
			taskRoute.GET("/self", middleware.UserAuth(), controller.GetUserTask)
			taskRoute.GET("/", middleware.AdminAuth(), controller.GetAllTask)
		}

		// 订阅相关路由
		subscriptionRoute := apiRouter.Group("/subscriptions")
		subscriptionRoute.Use(middleware.UserAuth())
		{
			subscriptionRoute.GET("/", controller.GetUserSubscriptions)                 // 获取用户订阅列表
			subscriptionRoute.POST("/", controller.CreateSubscription)                  // 创建订阅
			subscriptionRoute.PUT("/:id", controller.UpdateSubscription)                // 更新订阅
			subscriptionRoute.DELETE("/:id", controller.DeleteSubscription)             // 删除订阅
			subscriptionRoute.GET("/:id/articles", controller.GetSubscriptionArticles)  // 获取订阅下的文章
			subscriptionRoute.PUT("/:id/cancel", controller.CancelSubscription)         // 取消订阅
			subscriptionRoute.PUT("/:id/reactivate", controller.ReactivateSubscription) // 重新激活订阅
			subscriptionRoute.GET("/articles", controller.GetAllSubscriptionArticles)   // 获取当前用户订阅的所有文章
		}

		// 订阅文章管理路由（管理员功能）
		subscriptionArticleRoute := apiRouter.Group("/subscription_articles")
		subscriptionArticleRoute.Use(middleware.AdminAuth())
		{
			subscriptionArticleRoute.POST("/", controller.CreateSubscriptionArticle) // 创建订阅文章
		}

		// 系统推荐路由（用户相关，需要认证）
		apiRouter.GET("/user/recommendations", middleware.UserAuth(), controller.GetSystemRecommendations)            // 获取系统推荐列表
		apiRouter.POST("/user/recommendations/search", middleware.UserAuth(), controller.SearchSystemRecommendations) // 搜索系统推荐
		apiRouter.GET("/user/recommendations/:id", middleware.UserAuth(), controller.GetSystemRecommendationByID)     // 获取单个系统推荐
		apiRouter.GET("/user/welcome", middleware.UserAuth(), controller.GetWelcomePage)                              // 获取欢迎页面（首次访问）
		apiRouter.GET("/user/recommendations/change", middleware.UserAuth(), controller.GetRecommendationPage)        // 获取推荐页面（后续访问）

		// 系统推荐管理路由（管理员功能）
		recommendationRoute := apiRouter.Group("/recommendations")
		recommendationRoute.Use(middleware.AdminAuth())
		{
			recommendationRoute.POST("/", controller.CreateSystemRecommendation)      // 创建系统推荐
			recommendationRoute.PUT("/:id", controller.UpdateSystemRecommendation)    // 更新系统推荐
			recommendationRoute.DELETE("/:id", controller.DeleteSystemRecommendation) // 删除系统推荐
		}

		// 话题相关路由
		topicRoute := apiRouter.Group("/topics")
		topicRoute.Use(middleware.UserAuth())
		{
			topicRoute.GET("/", controller.GetTopics)                    // 获取话题列表
			topicRoute.POST("/", controller.CreateTopic)                 // 创建话题
			topicRoute.PUT("/:id", controller.UpdateTopicName)           // 更新话题名称
			topicRoute.DELETE("/:id", controller.DeleteTopic)            // 删除话题
			topicRoute.GET("/:id/messages", controller.GetTopicMessages) // 获取话题下的消息
			topicRoute.POST("/:id/messages", controller.CreateMessage)   // 发送消息
		}

		// API文档路由
		docsRoute := apiRouter.Group("/docs")
		{
			docsRoute.GET("/", controller.GetDocsIndex)    // 文档首页
			docsRoute.GET("/list", controller.GetDocsList) // 文档列表
			docsRoute.GET("/:type", controller.GetDocs)    // 具体文档页面
		}

		// 新的文档页面路由
		apiRouter.GET("/docs/docs.html", controller.GetDocsPage) // 新的文档页面
	}
}
