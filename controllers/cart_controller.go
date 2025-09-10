package controllers

import (
	"purches-backend/models"
	"purches-backend/services"
	"purches-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	cartService *services.CartService
}

func NewCartController(cartService *services.CartService) *CartController {
	return &CartController{
		cartService: cartService,
	}
}

// GetCart 获取购物车
func (cc *CartController) GetCart(c *gin.Context) {
	response, err := cc.cartService.GetCart()
	if err != nil {
		utils.ResponseError(c, 500, "获取购物车失败", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", response)
}

// AddToCart 添加商品到购物车
func (cc *CartController) AddToCart(c *gin.Context) {
	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	item, err := cc.cartService.AddToCart(req)
	if err != nil {
		utils.ResponseError(c, 500, "添加失败", err.Error())
		return
	}

	utils.ResponseOK(c, "添加成功", item)
}

// UpdateCartItem 更新购物车商品数量
func (cc *CartController) UpdateCartItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		utils.ResponseError(c, 400, "商品ID格式错误", err.Error())
		return
	}

	var req models.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	err = cc.cartService.UpdateCartItem(itemID, req.Count)
	if err != nil {
		utils.ResponseError(c, 500, "更新失败", err.Error())
		return
	}

	utils.ResponseOK(c, "更新成功", nil)
}

// DeleteCartItem 删除购物车商品
func (cc *CartController) DeleteCartItem(c *gin.Context) {
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		utils.ResponseError(c, 400, "商品ID格式错误", err.Error())
		return
	}

	err = cc.cartService.DeleteCartItem(itemID)
	if err != nil {
		utils.ResponseError(c, 404, "购物车中没有该商品", err.Error())
		return
	}

	utils.ResponseOK(c, "删除成功", nil)
}

// ClearCart 清空购物车
func (cc *CartController) ClearCart(c *gin.Context) {
	err := cc.cartService.ClearCart()
	if err != nil {
		utils.ResponseError(c, 500, "清空失败", err.Error())
		return
	}

	utils.ResponseOK(c, "清空成功", nil)
}
