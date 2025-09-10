package controllers

import (
	"purches-backend/services"
	"purches-backend/utils"

	"github.com/gin-gonic/gin"
)

type SupplierController struct {
	supplierService *services.SupplierService
}

func NewSupplierController(supplierService *services.SupplierService) *SupplierController {
	return &SupplierController{
		supplierService: supplierService,
	}
}

// GetSuppliers 获取供应商列表
func (sc *SupplierController) GetSuppliers(c *gin.Context) {
	response, err := sc.supplierService.GetSuppliers()
	if err != nil {
		utils.ResponseError(c, 500, "获取供应商列表失败", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", response)
}

// GetSupplier 获取供应商详情
func (sc *SupplierController) GetSupplier(c *gin.Context) {
	supplierName := c.Param("supplierName")

	response, err := sc.supplierService.GetSupplierDetail(supplierName)
	if err != nil {
		utils.ResponseError(c, 404, "供应商不存在", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", response)
}

// GetSupplierProducts 获取供应商的商品列表
func (sc *SupplierController) GetSupplierProducts(c *gin.Context) {
	supplierName := c.Param("supplierName")

	products, err := sc.supplierService.GetSupplierProducts(supplierName)
	if err != nil {
		utils.ResponseError(c, 500, "获取商品列表失败", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", products)
}

// GetSupplierOrders 获取供应商的订单列表
func (sc *SupplierController) GetSupplierOrders(c *gin.Context) {
	supplierName := c.Param("supplierName")

	orders, err := sc.supplierService.GetSupplierOrders(supplierName)
	if err != nil {
		utils.ResponseError(c, 500, "获取订单列表失败", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", orders)
}
