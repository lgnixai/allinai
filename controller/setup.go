package controller

import (
	"one-api/common"
	"one-api/constant"
	"one-api/model"
	"one-api/setting/operation_setting"
	"time"

	"github.com/gin-gonic/gin"
)

type Setup struct {
	Status       bool   `json:"status"`
	RootInit     bool   `json:"root_init"`
	DatabaseType string `json:"database_type"`
}

type SetupRequest struct {
	Username              string `json:"username"`
	Phone                 string `json:"phone"`
	PhoneVerificationCode string `json:"phone_verification_code"`
	SelfUseModeEnabled    bool   `json:"SelfUseModeEnabled"`
	DemoSiteEnabled       bool   `json:"DemoSiteEnabled"`
}

func GetSetup(c *gin.Context) {
	setup := Setup{
		Status: constant.Setup,
	}
	if constant.Setup {
		c.JSON(200, gin.H{
			"success": true,
			"data":    setup,
		})
		return
	}
	setup.RootInit = model.RootUserExists()
	if common.UsingMySQL {
		setup.DatabaseType = "mysql"
	}
	if common.UsingPostgreSQL {
		setup.DatabaseType = "postgres"
	}
	if common.UsingSQLite {
		setup.DatabaseType = "sqlite"
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    setup,
	})
}

func PostSetup(c *gin.Context) {
	// Check if setup is already completed
	if constant.Setup {
		c.JSON(400, gin.H{
			"success": false,
			"message": "系统已经初始化完成",
		})
		return
	}

	// Check if root user already exists
	rootExists := model.RootUserExists()

	var req SetupRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "请求参数有误",
		})
		return
	}

	// If root doesn't exist, validate and create admin account
	if !rootExists {
		// Validate username length: max 12 characters to align with model.User validation
		if len(req.Username) > 12 {
			c.JSON(400, gin.H{
				"success": false,
				"message": "用户名长度不能超过12个字符",
			})
			return
		}

		// Validate phone number
		if req.Phone == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "请输入手机号",
			})
			return
		}

		// Validate phone verification code
		if req.PhoneVerificationCode == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "请输入手机验证码",
			})
			return
		}

		// Check if phone number already exists
		var existingUser model.User
		err = model.DB.Where("phone = ?", req.Phone).First(&existingUser).Error
		if err == nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "手机号已被使用",
			})
			return
		}

		// Verify phone verification code
		expectedCode := "1111" // 这里应该从缓存或数据库中获取验证码
		if req.PhoneVerificationCode != expectedCode {
			c.JSON(400, gin.H{
				"success": false,
				"message": "手机验证码错误或已过期",
			})
			return
		}

		// Create root user
		rootUser := model.User{
			Username:    req.Username,
			Phone:       req.Phone,
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Root User",
			AccessToken: nil,
			Quota:       100000000,
		}
		err = model.DB.Create(&rootUser).Error
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "创建管理员账号失败: " + err.Error(),
			})
			return
		}
	}

	// Set operation modes
	operation_setting.SelfUseModeEnabled = req.SelfUseModeEnabled
	operation_setting.DemoSiteEnabled = req.DemoSiteEnabled

	// Save operation modes to database for persistence
	err = model.UpdateOption("SelfUseModeEnabled", boolToString(req.SelfUseModeEnabled))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "保存自用模式设置失败: " + err.Error(),
		})
		return
	}

	err = model.UpdateOption("DemoSiteEnabled", boolToString(req.DemoSiteEnabled))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "保存演示站点模式设置失败: " + err.Error(),
		})
		return
	}

	// Update setup status
	constant.Setup = true

	setup := model.Setup{
		Version:       common.Version,
		InitializedAt: time.Now().Unix(),
	}
	err = model.DB.Create(&setup).Error
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "系统初始化失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "系统初始化成功",
	})
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
