package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"one-api/common"
	"one-api/constant"
	"one-api/middleware"
	"one-api/model"
	"one-api/setting"
	"one-api/setting/console_setting"
	"one-api/setting/operation_setting"
	"one-api/setting/system_setting"
	"strings"

	"github.com/gin-gonic/gin"
)

func TestStatus(c *gin.Context) {
	err := model.PingDB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"message": "数据库连接失败",
		})
		return
	}
	// 获取HTTP统计信息
	httpStats := middleware.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Server is running",
		"http_stats": httpStats,
	})
	return
}

func GetStatus(c *gin.Context) {

	cs := console_setting.GetConsoleSetting()

	data := gin.H{
		"version":                  common.Version,
		"start_time":               common.StartTime,
		"email_verification":       common.EmailVerificationEnabled,
		"github_oauth":             common.GitHubOAuthEnabled,
		"github_client_id":         common.GitHubClientId,
		"linuxdo_oauth":            common.LinuxDOOAuthEnabled,
		"linuxdo_client_id":        common.LinuxDOClientId,
		"telegram_oauth":           common.TelegramOAuthEnabled,
		"telegram_bot_name":        common.TelegramBotName,
		"system_name":              common.SystemName,
		"logo":                     common.Logo,
		"footer_html":              common.Footer,
		"wechat_qrcode":            common.WeChatAccountQRCodeImageURL,
		"wechat_login":             common.WeChatAuthEnabled,
		"server_address":           setting.ServerAddress,
		"price":                    setting.Price,
		"stripe_unit_price":        setting.StripeUnitPrice,
		"min_topup":                setting.MinTopUp,
		"stripe_min_topup":         setting.StripeMinTopUp,
		"turnstile_check":          common.TurnstileCheckEnabled,
		"turnstile_site_key":       common.TurnstileSiteKey,
		"top_up_link":              common.TopUpLink,
		"docs_link":                operation_setting.GetGeneralSetting().DocsLink,
		"quota_per_unit":           common.QuotaPerUnit,
		"display_in_currency":      common.DisplayInCurrencyEnabled,
		"enable_batch_update":      common.BatchUpdateEnabled,
		"enable_drawing":           common.DrawingEnabled,
		"enable_task":              common.TaskEnabled,
		"enable_data_export":       common.DataExportEnabled,
		"data_export_default_time": common.DataExportDefaultTime,
		"default_collapse_sidebar": common.DefaultCollapseSidebar,
		"enable_online_topup":      setting.PayAddress != "" && setting.EpayId != "" && setting.EpayKey != "",
		"enable_stripe_topup":      setting.StripeApiSecret != "" && setting.StripeWebhookSecret != "" && setting.StripePriceId != "",
		"mj_notify_enabled":        setting.MjNotifyEnabled,
		"chats":                    setting.Chats,
		"demo_site_enabled":        operation_setting.DemoSiteEnabled,
		"self_use_mode_enabled":    operation_setting.SelfUseModeEnabled,
		"default_use_auto_group":   setting.DefaultUseAutoGroup,
		"pay_methods":              setting.PayMethods,
		"usd_exchange_rate":        setting.USDExchangeRate,

		// 面板启用开关
		"api_info_enabled":      cs.ApiInfoEnabled,
		"uptime_kuma_enabled":   cs.UptimeKumaEnabled,
		"announcements_enabled": cs.AnnouncementsEnabled,
		"faq_enabled":           cs.FAQEnabled,

		"oidc_enabled":                system_setting.GetOIDCSettings().Enabled,
		"oidc_client_id":              system_setting.GetOIDCSettings().ClientId,
		"oidc_authorization_endpoint": system_setting.GetOIDCSettings().AuthorizationEndpoint,
		"setup":                       constant.Setup,
	}

	// 根据启用状态注入可选内容
	if cs.ApiInfoEnabled {
		data["api_info"] = console_setting.GetApiInfo()
	}
	if cs.AnnouncementsEnabled {
		data["announcements"] = console_setting.GetAnnouncements()
	}
	if cs.FAQEnabled {
		data["faq"] = console_setting.GetFAQ()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
	return
}

func GetNotice(c *gin.Context) {
	common.OptionMapRWMutex.RLock()
	defer common.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    common.OptionMap["Notice"],
	})
	return
}

func GetAbout(c *gin.Context) {
	common.OptionMapRWMutex.RLock()
	defer common.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    common.OptionMap["About"],
	})
	return
}

func GetMidjourney(c *gin.Context) {
	common.OptionMapRWMutex.RLock()
	defer common.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    common.OptionMap["Midjourney"],
	})
	return
}

func GetHomePageContent(c *gin.Context) {
	common.OptionMapRWMutex.RLock()
	defer common.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    common.OptionMap["HomePageContent"],
	})
	return
}

func SendEmailVerification(c *gin.Context) {
	email := c.Query("email")
	if err := common.Validate.Var(email, "required,email"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的邮箱地址",
		})
		return
	}
	localPart := parts[0]
	domainPart := parts[1]
	if common.EmailDomainRestrictionEnabled {
		allowed := false
		for _, domain := range common.EmailDomainWhitelist {
			if domainPart == domain {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "The administrator has enabled the email domain name whitelist, and your email address is not allowed due to special symbols or it's not in the whitelist.",
			})
			return
		}
	}
	if common.EmailAliasRestrictionEnabled {
		containsSpecialSymbols := strings.Contains(localPart, "+") || strings.Contains(localPart, ".")
		if containsSpecialSymbols {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "管理员已启用邮箱地址别名限制，您的邮箱地址由于包含特殊符号而被拒绝。",
			})
			return
		}
	}

	if model.IsEmailAlreadyTaken(email) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "邮箱地址已被占用",
		})
		return
	}
	code := common.GenerateVerificationCode(6)
	common.RegisterVerificationCodeWithKey(email, code, common.EmailVerificationPurpose)
	subject := fmt.Sprintf("%s邮箱验证邮件", common.SystemName)
	content := fmt.Sprintf("<p>您好，你正在进行%s邮箱验证。</p>"+
		"<p>您的验证码为: <strong>%s</strong></p>"+
		"<p>验证码 %d 分钟内有效，如果不是本人操作，请忽略。</p>", common.SystemName, code, common.VerificationValidMinutes)
	err := common.SendEmail(subject, email, content)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func SendPasswordResetPhone(c *gin.Context) {
	phone := c.Query("phone")
	if err := common.Validate.Var(phone, "required,len=11"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的手机号格式",
		})
		return
	}
	if !model.IsPhoneAlreadyTaken(phone) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "该手机号未注册",
		})
		return
	}
	code := "1111" //common.GenerateVerificationCode(6)
	common.RegisterVerificationCodeWithKey(phone, code, common.PasswordResetPurpose)

	// 这里可以集成短信服务发送验证码
	// 目前先返回验证码用于测试
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    code, // 实际生产环境中应该移除这个字段
	})
	return
}

type PasswordResetRequest struct {
	Phone string `json:"phone"`
	Token string `json:"token"`
}

func ResetPassword(c *gin.Context) {
	var req PasswordResetRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if req.Phone == "" || req.Token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if !common.VerifyCodeWithKey(req.Phone, req.Token, common.PasswordResetPurpose) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "重置验证码错误或已过期",
		})
		return
	}
	password := common.GenerateVerificationCode(12)
	err = model.ResetUserPasswordByPhone(req.Phone, password)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.DeleteKey(req.Phone, common.PasswordResetPurpose)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    password,
	})
	return
}

func SendPhoneVerification(c *gin.Context) {
	phone := c.Query("phone")
	purpose := c.Query("purpose") // 添加用途参数：register 或 login
	
	if err := common.Validate.Var(phone, "required,len=11"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的手机号格式",
		})
		return
	}

	isPhoneTaken := model.IsPhoneAlreadyTaken(phone)
	
	// 根据用途检查手机号状态
	if purpose == "login" {
		// 登录时，手机号必须已存在
		if !isPhoneTaken {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "手机号未注册",
			})
			return
		}
	} else {
		// 注册时（默认情况），手机号不能已存在
		if isPhoneTaken {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "手机号已被占用",
			})
			return
		}
	}

	code := "1111" //common.GenerateVerificationCode(6)
	common.RegisterVerificationCodeWithKey(phone, code, common.PhoneVerificationPurpose)

	// 这里可以集成短信服务发送验证码
	// 目前先返回验证码用于测试
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    code, // 实际生产环境中应该移除这个字段
	})
	return
}
