package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// GetDocsIndex 获取文档首页
func GetDocsIndex(c *gin.Context) {
	// 提供HTML首页
	filePath := filepath.Join("docs", "index.html")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果HTML文件不存在，重定向到API文档
		c.Redirect(http.StatusMovedPermanently, "/api/docs/api")
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/html; charset=utf-8")

	// 返回HTML文件
	c.File(filePath)
}

// GetDocsPage 获取新的文档页面
func GetDocsPage(c *gin.Context) {
	// 提供新的文档页面
	filePath := filepath.Join("docs", "docs.html")

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 如果文件不存在，返回404
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "文档页面不存在",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/html; charset=utf-8")

	// 返回HTML文件
	c.File(filePath)
}

// GetDocs 获取具体文档
func GetDocs(c *gin.Context) {
	docType := c.Param("type")

	// 文档文件映射
	docFiles := map[string]string{
		"api":             "API_Documentation.md",
		"postman":         "Postman_Usage_Guide.md",
		"deployment":      "Deployment_Guide.md",
		"auto-deployment": "Auto_Deployment_Guide.md",
		"auth":            "api_auth.md",
	}

	fileName, exists := docFiles[docType]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "文档不存在",
		})
		return
	}

	// 构建文件路径
	filePath := filepath.Join("docs", fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "文档文件不存在",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.Header("Content-Disposition", "inline; filename="+fileName)

	// 返回文件内容
	c.File(filePath)
}

// GetDocsList 获取文档列表
func GetDocsList(c *gin.Context) {
	docs := []gin.H{
		{
			"key":      "api",
			"title":    "API 文档",
			"filename": "API_Documentation.md",
			"url":      "/api/docs/api",
		},
		{
			"key":      "postman",
			"title":    "Postman 使用指南",
			"filename": "Postman_Usage_Guide.md",
			"url":      "/api/docs/postman",
		},
		{
			"key":      "deployment",
			"title":    "部署指南",
			"filename": "Deployment_Guide.md",
			"url":      "/api/docs/deployment",
		},
		{
			"key":      "auto-deployment",
			"title":    "自动部署指南",
			"filename": "Auto_Deployment_Guide.md",
			"url":      "/api/docs/auto-deployment",
		},
		{
			"key":      "auth",
			"title":    "认证说明",
			"filename": "api_auth.md",
			"url":      "/api/docs/auth",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    docs,
	})
}
