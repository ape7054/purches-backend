package controllers

import (
	"math"
	"purches-backend/models"
	"purches-backend/services"
	"purches-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// GetProducts 获取商品列表
func (pc *ProductController) GetProducts(c *gin.Context) {
	var req models.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}

	products, total, err := pc.productService.GetProducts(req)
	if err != nil {
		utils.ResponseError(c, 500, "获取商品列表失败", err.Error())
		return
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	response := models.ProductListResponse{
		Products: products,
		Pagination: models.PaginationResponse{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	utils.ResponseOK(c, "获取成功", response)
}

// GetProduct 获取商品详情
func (pc *ProductController) GetProduct(c *gin.Context) {
	productID := c.Param("productId")
	id, err := strconv.Atoi(productID)
	if err != nil {
		utils.ResponseError(c, 400, "商品ID格式错误", err.Error())
		return
	}

	product, err := pc.productService.GetProductByID(id)
	if err != nil {
		utils.ResponseError(c, 404, "商品不存在", err.Error())
		return
	}

	utils.ResponseOK(c, "获取成功", product)
}

// ImportProducts 批量导入商品
func (pc *ProductController) ImportProducts(c *gin.Context) {
	var req models.ImportProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, 400, "请求参数错误", err.Error())
		return
	}

	createdProducts, err := pc.productService.ImportProducts(req.Products)
	if err != nil {
		utils.ResponseError(c, 500, "导入失败", err.Error())
		return
	}

	utils.ResponseOK(c, "导入成功", gin.H{
		"count":    len(createdProducts),
		"products": createdProducts,
	})
}
