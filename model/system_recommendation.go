package model

import "time"

// SystemRecommendation 系统推荐订阅表
type SystemRecommendation struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title" gorm:"not null;size:255"`           // 推荐主题标题
	Description     string    `json:"description" gorm:"type:text"`            // 主题描述
	Category        string    `json:"category" gorm:"size:100"`                // 分类（如：技术、商业、设计等）
	SubscriptionCount int     `json:"subscription_count" gorm:"default:0"`     // 订阅数
	ArticleCount    int       `json:"article_count" gorm:"default:0"`          // 相关文章数
	Status          int       `json:"status" gorm:"default:1"`                 // 状态：1-启用，0-禁用
	SortOrder       int       `json:"sort_order" gorm:"default:0"`             // 排序权重
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (SystemRecommendation) TableName() string {
	return "system_recommendations"
}

// GetSystemRecommendations 获取系统推荐列表
func GetSystemRecommendations(page, pageSize int) ([]SystemRecommendation, int64, error) {
	var recommendations []SystemRecommendation
	var total int64

	// 获取总数
	err := DB.Model(&SystemRecommendation{}).
		Where("status = 1").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where("status = 1").
		Order("sort_order DESC, subscription_count DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&recommendations).Error

	return recommendations, total, nil
}

// SearchSystemRecommendations 搜索系统推荐
func SearchSystemRecommendations(keyword string, page, pageSize int) ([]SystemRecommendation, int64, error) {
	var recommendations []SystemRecommendation
	var total int64

	// 构建搜索条件
	searchCondition := "status = 1 AND (title LIKE ? OR description LIKE ? OR category LIKE ?)"
	searchPattern := "%" + keyword + "%"

	// 获取总数
	err := DB.Model(&SystemRecommendation{}).
		Where(searchCondition, searchPattern, searchPattern, searchPattern).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = DB.Where(searchCondition, searchPattern, searchPattern, searchPattern).
		Order("sort_order DESC, subscription_count DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&recommendations).Error

	return recommendations, total, nil
}

// GetSystemRecommendationByID 根据ID获取系统推荐
func GetSystemRecommendationByID(id int) (*SystemRecommendation, error) {
	var recommendation SystemRecommendation
	err := DB.Where("id = ? AND status = 1", id).First(&recommendation).Error
	if err != nil {
		return nil, err
	}
	return &recommendation, nil
}

// CreateSystemRecommendation 创建系统推荐
func CreateSystemRecommendation(recommendation *SystemRecommendation) error {
	return DB.Create(recommendation).Error
}

// UpdateSystemRecommendation 更新系统推荐
func UpdateSystemRecommendation(recommendation *SystemRecommendation) error {
	return DB.Save(recommendation).Error
}

// DeleteSystemRecommendation 删除系统推荐（软删除）
func DeleteSystemRecommendation(id int) error {
	return DB.Model(&SystemRecommendation{}).Where("id = ?", id).Update("status", 0).Error
}

// UpdateSubscriptionCount 更新订阅数
func UpdateSubscriptionCount(id int, count int) error {
	return DB.Model(&SystemRecommendation{}).Where("id = ?", id).Update("subscription_count", count).Error
}

// UpdateArticleCount 更新文章数
func UpdateArticleCount(id int, count int) error {
	return DB.Model(&SystemRecommendation{}).Where("id = ?", id).Update("article_count", count).Error
}

// InitializeSystemRecommendations 初始化系统推荐数据
func InitializeSystemRecommendations() error {
	// 检查是否已有数据
	var count int64
	DB.Model(&SystemRecommendation{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，不重复初始化
	}

	// 系统推荐的默认主题
	recommendations := []SystemRecommendation{
		{
			Title:           "人工智能与机器学习",
			Description:     "探索AI和ML的最新发展，包括深度学习、自然语言处理、计算机视觉等前沿技术",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       100,
		},
		{
			Title:           "Web开发技术",
			Description:     "涵盖前端、后端、全栈开发技术，包括React、Vue、Node.js、Python等热门技术栈",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       95,
		},
		{
			Title:           "移动应用开发",
			Description:     "iOS、Android开发技术，包括Swift、Kotlin、React Native、Flutter等跨平台解决方案",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       90,
		},
		{
			Title:           "云计算与DevOps",
			Description:     "云服务、容器化、CI/CD、微服务架构等现代化开发和部署技术",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       85,
		},
		{
			Title:           "数据科学与分析",
			Description:     "大数据处理、数据分析、数据可视化、商业智能等相关技术和应用",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       80,
		},
		{
			Title:           "网络安全",
			Description:     "网络安全、信息安全、隐私保护、漏洞分析等安全相关技术和最佳实践",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       75,
		},
		{
			Title:           "创业与商业",
			Description:     "创业指导、商业模式、市场营销、产品管理等商业相关内容",
			Category:        "商业",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       70,
		},
		{
			Title:           "产品设计与用户体验",
			Description:     "产品设计、UI/UX设计、用户研究、交互设计等设计相关内容",
			Category:        "设计",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       65,
		},
		{
			Title:           "编程语言与框架",
			Description:     "各种编程语言的学习指南、框架介绍、最佳实践等",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       60,
		},
		{
			Title:           "数据库与存储",
			Description:     "关系型数据库、NoSQL、分布式存储、数据建模等技术",
			Category:        "技术",
			SubscriptionCount: 0,
			ArticleCount:    0,
			Status:          1,
			SortOrder:       55,
		},
	}

	// 批量创建
	for _, rec := range recommendations {
		err := CreateSystemRecommendation(&rec)
		if err != nil {
			return err
		}
	}

	return nil
}
