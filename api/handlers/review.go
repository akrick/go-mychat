package handlers

import (
	"context"
	"akrick.com/mychat/cache"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	OrderID         uint   `json:"order_id" binding:"required"`
	Rating          int    `json:"rating" binding:"required,min=1,max=5"`
	ServiceRating   int    `json:"service_rating" binding:"required,min=1,max=5"`
	Professionalism int   `json:"professionalism" binding:"required,min=1,max=5"`
	Effectiveness  int    `json:"effectiveness" binding:"required,min=1,max=5"`
	Content         string `json:"content" binding:"required,max=500"`
	IsAnonymous     bool   `json:"is_anonymous"`
}

type ReplyReviewRequest struct {
	ReplyContent string `json:"reply_content" binding:"required,max=500"`
}

// CreateReview godoc
// @Summary 创建评价
// @Description 用户对已完成的订单进行评价
// @Tags 评价
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateReviewRequest true "评价信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:评价成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 400 {object} map[string]interface{} "订单状态不允许评价"
// @Router /api/review/create [post]
func CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询订单信息
	var order models.Order
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查订单是否属于当前用户
	if order.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权评价此订单",
		})
		return
	}

	// 检查订单状态，只有已完成的订单才能评价
	if order.Status != models.OrderStatusCompleted {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "订单状态不允许评价",
		})
		return
	}

	// 检查是否已评价过
	var existingReview models.Review
	if err := database.DB.Where("order_id = ?", req.OrderID).First(&existingReview).Error; err == nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "该订单已评价",
		})
		return
	}

	// 创建评价
	review := models.Review{
		OrderID:         req.OrderID,
		OrderNo:         order.OrderNo,
		UserID:          userID.(uint),
		CounselorID:     order.CounselorID,
		Rating:          req.Rating,
		ServiceRating:   req.ServiceRating,
		Professionalism: req.Professionalism,
		Effectiveness:  req.Effectiveness,
		Content:         req.Content,
		IsAnonymous:     req.IsAnonymous,
		Status:          1,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建评价失败: " + err.Error(),
		})
		return
	}

	// 更新咨询师统计
	updateCounselorStatistics(c, order.CounselorID, req.Rating, 1)

	// 清除缓存
	ctx := context.Background()
	cache.DeleteOrderCache(ctx, order.ID)
	cache.DeleteCounselorCache(ctx, order.CounselorID)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "评价成功",
		"data": gin.H{
			"review_id": review.ID,
		},
	})
}

// GetReviewList godoc
// @Summary 获取咨询师评价列表
// @Description 获取指定咨询师的评价列表
// @Tags 评价
// @Accept json
// @Produce json
// @Param counselor_id path int true "咨询师ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param rating query int false "评分筛选"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{reviews,total}"
// @Router /api/review/counselor/{counselor_id} [get]
func GetReviewList(c *gin.Context) {
	counselorID := c.Param("counselor_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	rating := c.Query("rating")

	query := database.DB.Model(&models.Review{}).Where("counselor_id = ? AND status = 1", counselorID)

	// 评分筛选
	if rating != "" {
		query = query.Where("rating = ?", rating)
	}

	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("User").Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&reviews).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"reviews": reviews,
			"total":   total,
		},
	})
}

// GetReviewDetail godoc
// @Summary 获取评价详情
// @Description 获取评价详细信息
// @Tags 评价
// @Accept json
// @Produce json
// @Param id path int true "评价ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{review}"
// @Failure 404 {object} map[string]interface{} "评价不存在"
// @Router /api/review/{id} [get]
func GetReviewDetail(c *gin.Context) {
	reviewID := c.Param("id")

	var review models.Review
	if err := database.DB.Preload("Order").Preload("User").Preload("Counselor").First(&review, reviewID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "评价不存在",
		})
		return
	}

	// 如果是匿名评价且不是评价者本人或咨询师，隐藏用户信息
	if review.IsAnonymous {
		userID, exists := c.Get("user_id")
		if !exists || (userID.(uint) != review.UserID && userID.(uint) != review.CounselorID) {
			review.User = models.User{
				Username: "匿名用户",
			}
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": review,
	})
}

// ReplyReview godoc
// @Summary 咨询师回复评价
// @Description 咨询师回复用户评价
// @Tags 评价
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评价ID"
// @Param request body ReplyReviewRequest true "回复内容"
// @Success 200 {object} map[string]interface{} "code:200,msg:回复成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Router /api/review/{id}/reply [post]
func ReplyReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	reviewID := c.Param("id")
	ctx := context.Background()

	var req ReplyReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询评价信息
	var review models.Review
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "评价不存在",
		})
		return
	}

	// 检查是否为咨询师本人
	if review.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权回复此评价",
		})
		return
	}

	// 更新回复
	now := database.DB.NowFunc()
	if err := database.DB.Model(&review).Updates(map[string]interface{}{
		"reply_content": req.ReplyContent,
		"reply_time":    &now,
	}).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "回复失败: " + err.Error(),
		})
		return
	}

	// 清除缓存
	cache.DeleteCounselorCache(ctx, review.CounselorID)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "回复成功",
	})
}

// GetUserReviews godoc
// @Summary 获取用户评价列表
// @Description 获取当前用户的评价列表
// @Tags 评价
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{reviews,total}"
// @Router /api/review/my [get]
func GetUserReviews(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	query := database.DB.Model(&models.Review{}).Where("user_id = ?", userID)

	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Counselor").Preload("Order").Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&reviews).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "查询失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"reviews": reviews,
			"total":   total,
		},
	})
}

// GetCounselorStatistics godoc
// @Summary 获取咨询师统计数据
// @Description 获取指定咨询师的统计数据
// @Tags 评价
// @Accept json
// @Produce json
// @Param counselor_id path int true "咨询师ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{statistics}"
// @Router /api/review/counselor/{counselor_id}/statistics [get]
func GetCounselorStatistics(c *gin.Context) {
	counselorID := c.Param("counselor_id")

	var statistics models.CounselorStatistics
	if err := database.DB.Preload("Counselor").Where("counselor_id = ?", counselorID).First(&statistics).Error; err != nil {
		// 如果统计不存在，返回默认值
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": gin.H{
				"total_orders":      0,
				"completed_orders": 0,
				"cancelled_orders": 0,
				"total_duration":   0,
				"total_amount":     0,
				"review_count":     0,
				"avg_rating":       0,
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": statistics,
	})
}

// updateCounselorStatistics 更新咨询师统计数据
func updateCounselorStatistics(c *gin.Context, counselorID uint, rating int, isCompleted int) {
	_ = c // 参数暂时未使用

	var statistics models.CounselorStatistics

	// 查找或创建统计记录
	if err := database.DB.Where("counselor_id = ?", counselorID).First(&statistics).Error; err != nil {
		// 创建新统计
		statistics = models.CounselorStatistics{
			CounselorID:     counselorID,
			TotalOrders:     0,
			CompletedOrders: 0,
			CancelledOrders: 0,
			TotalDuration:   0,
			TotalAmount:     0,
			ReviewCount:     0,
			AvgRating:       0,
			SumRating:       0,
		}
		database.DB.Create(&statistics)
	}

	// 更新统计
	updates := map[string]interface{}{
		"total_orders":  statistics.TotalOrders + 1,
		"sum_rating":   statistics.SumRating + rating,
		"review_count": statistics.ReviewCount + 1,
	}

	if isCompleted == 1 {
		updates["completed_orders"] = statistics.CompletedOrders + 1
	}

	// 计算平均评分
	newReviewCount := statistics.ReviewCount + 1
	newAvgRating := float64(statistics.SumRating+rating) / float64(newReviewCount)
	updates["avg_rating"] = newAvgRating

	database.DB.Model(&statistics).Updates(updates)
}
