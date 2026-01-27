package handlers

import (
	"context"
	"fmt"
	"time"
	"akrick.com/mychat/cache"
	"akrick.com/mychat/database"
	"akrick.com/mychat/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	CounselorID  uint      `json:"counselor_id" binding:"required"`
	Duration     int       `json:"duration" binding:"required,min=15,max=180"`
	ScheduleTime time.Time `json:"schedule_time" binding:"required"`
	Notes        string    `json:"notes"`
}

type UpdateOrderRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1 2 3 4"`
}

// CreateOrder godoc
// @Summary 创建订单
// @Description 用户创建心理咨询订单
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "订单信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:创建成功,data:{order}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "咨询师不存在"
// @Router /api/order/create [post]
func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询咨询师信息
	var counselor models.Counselor
	if err := database.DB.First(&counselor, req.CounselorID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "咨询师不存在",
		})
		return
	}

	// 检查咨询师状态
	if counselor.Status != 1 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "咨询师暂不可用",
		})
		return
	}

	// 计算订单金额
	amount := float64(req.Duration) * counselor.Price

	// 生成订单号
	orderNo := fmt.Sprintf("ORD%d%d", time.Now().Unix(), userID.(uint))

	// 创建订单
	order := models.Order{
		OrderNo:      orderNo,
		UserID:       userID.(uint),
		CounselorID:  req.CounselorID,
		Duration:     req.Duration,
		Amount:       amount,
		Status:       models.OrderStatusPending,
		ScheduleTime: req.ScheduleTime,
		Notes:        req.Notes,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "创建订单失败: " + err.Error(),
		})
		return
	}

	// 清除用户和咨询师订单列表缓存
	if cache.Rdb != nil {
		cache.InvalidateUserOrdersCache(context.Background(), order.UserID)
		cache.InvalidateCounselorOrdersCache(context.Background(), order.CounselorID)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建订单成功",
		"data": gin.H{
			"order_id":  order.ID,
			"order_no":  order.OrderNo,
			"amount":    order.Amount,
			"status":    order.Status,
			"create_at": order.CreatedAt,
		},
	})
}

// GetOrderDetail godoc
// @Summary 获取订单详情
// @Description 获取订单详细信息（使用Redis缓存和SingleFlight防穿透）
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{order}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "无权访问"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Router /api/order/{id} [get]
func GetOrderDetail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("id")

	ctx := context.Background()
	var order models.Order
	var fromCache bool

	// 尝试从缓存获取订单信息
	if cache.Rdb != nil {
		id := parseUint(orderID)
		orderData, err := cache.GetOrderWithCache(ctx, id)
		if err == nil {
			order = models.Order{
				ID:           orderData.ID,
				OrderNo:      orderData.OrderNo,
				UserID:       orderData.UserID,
				CounselorID:  orderData.CounselorID,
				Duration:     orderData.Duration,
				Amount:       orderData.Amount,
				Status:       orderData.Status,
				ScheduleTime: orderData.ScheduleTime,
				Notes:        orderData.Notes,
				PayTime:      orderData.PayTime,
				CreatedAt:    orderData.CreatedAt,
				UpdatedAt:    orderData.UpdatedAt,
			}
			fromCache = true
		}
	}

	// 缓存未命中，从数据库查询并加载关联
	if !fromCache {
		if err := database.DB.Preload("User").Preload("Counselor").First(&order, orderID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{
					"code": 404,
					"msg":  "订单不存在",
				})
			} else {
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "查询失败: " + err.Error(),
				})
			}
			return
		}
	} else {
		// 从缓存获取后，仍需加载关联的用户和咨询师信息
		database.DB.Model(&order).Association("User").Find(&order.User)
		database.DB.Model(&order).Association("Counselor").Find(&order.Counselor)
	}

	// 检查权限：只有订单用户和咨询师可以查看
	if order.UserID != userID.(uint) && order.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权访问此订单",
		})
		return
	}

	msg := "获取成功"
	if fromCache {
		msg = "获取成功（来自缓存）"
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  msg,
		"data": gin.H{
			"id":            order.ID,
			"order_no":      order.OrderNo,
			"user":          order.User,
			"counselor":     order.Counselor,
			"duration":      order.Duration,
			"amount":        order.Amount,
			"status":        order.Status,
			"schedule_time": order.ScheduleTime,
			"notes":         order.Notes,
			"pay_time":      order.PayTime,
			"created_at":    order.CreatedAt,
			"updated_at":    order.UpdatedAt,
		},
	})
}

// GetUserOrders godoc
// @Summary 获取用户订单列表
// @Description 获取当前用户的订单列表
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query int false "订单状态:0-待支付,1-已支付,2-已完成,3-已取消,4-已退款"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{orders,total}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Router /api/order/list [get]
func GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	status := c.Query("status")

	query := database.DB.Model(&models.Order{}).Where("user_id = ?", userID)

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("Counselor").Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&orders).Error; err != nil {
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
			"orders": orders,
			"total":  total,
		},
	})
}

// GetCounselorOrders godoc
// @Summary 获取咨询师的订单列表
// @Description 咨询师查看自己的订单列表
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query int false "订单状态:0-待支付,1-已支付,2-已完成,3-已取消,4-已退款"
// @Success 200 {object} map[string]interface{} "code:200,msg:获取成功,data:{orders,total}"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "不是咨询师"
// @Router /api/counselor/orders [get]
func GetCounselorOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 检查是否为咨询师
	var counselor models.Counselor
	if err := database.DB.First(&counselor, userID).Error; err != nil {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "不是咨询师",
		})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	status := c.Query("status")

	query := database.DB.Model(&models.Order{}).Where("counselor_id = ?", counselor.ID)

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := 0
	if page == "1" {
		offset = 0
	} else {
		offset = (parseInt(page) - 1) * parseInt(pageSize)
	}

	if err := query.Preload("User").Offset(offset).Limit(parseInt(pageSize)).Order("created_at DESC").Find(&orders).Error; err != nil {
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
			"orders": orders,
			"total":  total,
		},
	})
}

// UpdateOrderStatus godoc
// @Summary 更新订单状态
// @Description 更新订单状态（仅咨询师或管理员可操作）
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param request body UpdateOrderRequest true "状态信息"
// @Success 200 {object} map[string]interface{} "code:200,msg:更新成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Router /api/order/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("id")
	ctx := context.Background()

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 查询订单
	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查权限：只有咨询师可以更新订单状态
	if order.CounselorID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权操作此订单",
		})
		return
	}

	// 更新订单状态
	updates := map[string]interface{}{
		"status": req.Status,
	}

	// 如果是已支付状态，记录支付时间
	if req.Status == models.OrderStatusPaid {
		now := time.Now()
		updates["pay_time"] = &now
	}

	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "更新失败: " + err.Error(),
		})
		return
	}

	// 删除缓存
	if cache.Rdb != nil {
		cache.DeleteOrderCache(ctx, order.ID)
		cache.InvalidateUserOrdersCache(ctx, order.UserID)
		cache.InvalidateCounselorOrdersCache(ctx, order.CounselorID)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// CancelOrder godoc
// @Summary 取消订单
// @Description 用户取消订单（仅待支付或已支付状态可取消）
// @Tags 订单
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} map[string]interface{} "code:200,msg:取消成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 403 {object} map[string]interface{} "无权操作"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Failure 400 {object} map[string]interface{} "订单状态不允许取消"
// @Router /api/order/{id}/cancel [post]
func CancelOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderID := c.Param("id")
	ctx := context.Background()

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(404, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 检查权限
	if order.UserID != userID.(uint) {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "无权操作此订单",
		})
		return
	}

	// 检查订单状态
	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusPaid {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "当前状态不允许取消订单",
		})
		return
	}

	if err := database.DB.Model(&order).Update("status", models.OrderStatusCancelled).Error; err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "取消失败: " + err.Error(),
		})
		return
	}

	// 删除缓存
	if cache.Rdb != nil {
		cache.DeleteOrderCache(ctx, order.ID)
		cache.InvalidateUserOrdersCache(ctx, order.UserID)
		cache.InvalidateCounselorOrdersCache(ctx, order.CounselorID)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "取消成功",
	})
}
