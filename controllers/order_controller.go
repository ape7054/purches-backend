package controllers

import (
	"purches-backend/models"
	"purches-backend/services"
	"purches-backend/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// CreateOrder 提交订单（按供应商分组）
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	createdOrders, err := oc.orderService.CreateOrder(req)
	if err != nil {
		utils.ResponseError(c, 500, "创建订单失败", err.Error())
		return
	}

	response := models.CreateOrderResponse{
		Orders: createdOrders,
	}

	utils.ResponseOK(c, "订单提交成功", response)
}

// GetOrders 获取订单列表
func (oc *OrderController) GetOrders(c *gin.Context) {
	var req models.OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	response, err := oc.orderService.GetOrders(req)
	if err != nil {
		utils.ResponseError(c, 500, "获取订单列表失败", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", response)
}

// GetOrder 获取订单详情
func (oc *OrderController) GetOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	order, err := oc.orderService.GetOrderByID(orderID)
	if err != nil {
		utils.ResponseError(c, 404, "订单不存在", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", order)
}

// UpdateOrderStatus 更新订单状态
func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderId")

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	err := oc.orderService.UpdateOrderStatus(orderID, req)
	if err != nil {
		utils.ResponseError(c, 500, "更新失败", err.Error())
		return
	}

	utils.ResponseOK(c, "更新成功", nil)
}

// UpdateOrderFinalPrice 更新订单最终价格
func (oc *OrderController) UpdateOrderFinalPrice(c *gin.Context) {
	orderID := c.Param("orderId")

	var req models.UpdateOrderPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	err := oc.orderService.UpdateOrderFinalPrice(orderID, req.FinalPrice)
	if err != nil {
		utils.ResponseError(c, 500, "更新失败", err.Error())
		return
	}

	utils.ResponseOK(c, "价格更新成功", nil)
}

// ExportOrders 导出订单数据
func (oc *OrderController) ExportOrders(c *gin.Context) {
	var req models.ExportOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	orders, err := oc.orderService.ExportOrders(req)
	if err != nil {
		utils.ResponseError(c, 500, "导出失败", err.Error())
		return
	}

	utils.ResponseOK(c, "导出成功", orders)
}
