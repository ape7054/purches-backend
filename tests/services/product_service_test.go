package services

import (
	"purches-backend/models"
	"purches-backend/services"
	"purches-backend/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductService_GetProducts(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	productService := services.NewProductService(db)

	t.Run("获取所有商品", func(t *testing.T) {
		req := models.ProductListRequest{
			Page:  1,
			Limit: 10,
		}

		products, total, err := productService.GetProducts(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, products, 3)
		assert.Equal(t, "测试商品1", products[0].Name)
	})

	t.Run("按供应商筛选商品", func(t *testing.T) {
		req := models.ProductListRequest{
			Page:     1,
			Limit:    10,
			Supplier: "测试供应商A",
		}

		products, total, err := productService.GetProducts(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, products, 2)
		for _, product := range products {
			assert.Equal(t, "测试供应商A", product.Supplier)
		}
	})

	t.Run("搜索商品", func(t *testing.T) {
		req := models.ProductListRequest{
			Page:   1,
			Limit:  10,
			Search: "商品1",
		}

		products, total, err := productService.GetProducts(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, products, 1)
		assert.Equal(t, "测试商品1", products[0].Name)
	})

	t.Run("按状态筛选商品", func(t *testing.T) {
		req := models.ProductListRequest{
			Page:   1,
			Limit:  10,
			Status: "available",
		}

		products, total, err := productService.GetProducts(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, products, 3)
	})

	t.Run("分页测试", func(t *testing.T) {
		req := models.ProductListRequest{
			Page:  1,
			Limit: 2,
		}

		products, total, err := productService.GetProducts(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, products, 2)

		// 测试第二页
		req.Page = 2
		products2, total2, err := productService.GetProducts(req)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), total2)
		assert.Len(t, products2, 1)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}

func TestProductService_GetProductByID(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	productService := services.NewProductService(db)

	t.Run("获取存在的商品", func(t *testing.T) {
		product, err := productService.GetProductByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "测试商品1", product.Name)
		assert.Equal(t, 10.50, product.Price)
		assert.Equal(t, "个", product.Unit)
	})

	t.Run("获取不存在的商品", func(t *testing.T) {
		product, err := productService.GetProductByID(999)

		assert.Error(t, err)
		assert.Nil(t, product)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}

func TestProductService_ImportProducts(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据（需要先有供应商）
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	productService := services.NewProductService(db)

	t.Run("导入新商品", func(t *testing.T) {
		importProducts := []models.ImportProduct{
			{
				Name:        "导入测试商品1",
				Price:       20.00,
				Unit:        "盒",
				Description: "导入的测试商品描述1",
				Supplier:    "测试供应商A",
			},
			{
				Name:        "导入测试商品2",
				Price:       30.50,
				Unit:        "袋",
				Description: "导入的测试商品描述2",
				Supplier:    "测试供应商B",
			},
		}

		createdProducts, err := productService.ImportProducts(importProducts)

		assert.NoError(t, err)
		assert.Len(t, createdProducts, 2)
		assert.Equal(t, "导入测试商品1", createdProducts[0].Name)
		assert.Equal(t, 20.00, createdProducts[0].Price)
		assert.Equal(t, "available", createdProducts[0].Status)
	})

	t.Run("导入空商品列表", func(t *testing.T) {
		importProducts := []models.ImportProduct{}

		createdProducts, err := productService.ImportProducts(importProducts)

		assert.NoError(t, err)
		assert.Len(t, createdProducts, 0)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}
