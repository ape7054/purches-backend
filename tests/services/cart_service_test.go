package services

import (
	"purches-backend/models"
	"purches-backend/services"
	"purches-backend/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartService_GetCart(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	cartService := services.NewCartService(db)

	t.Run("获取空购物车", func(t *testing.T) {
		response, err := cartService.GetCart()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.Items)
		assert.Equal(t, 0, response.Summary.TotalItems)
		assert.Equal(t, 0.0, response.Summary.TotalPrice)
		assert.Equal(t, 0, response.Summary.SupplierCount)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}

func TestCartService_AddToCart(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	cartService := services.NewCartService(db)

	t.Run("添加新商品到购物车", func(t *testing.T) {
		req := models.AddToCartRequest{
			ProductID: 1,
			Count:     2,
		}

		item, err := cartService.AddToCart(req)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, 1, item.ProductID)
		assert.Equal(t, "测试商品1", item.Name)
		assert.Equal(t, 2, item.Count)
		assert.Equal(t, 10.50, item.Price)
		assert.Equal(t, 21.0, item.TotalPrice)
		assert.Equal(t, "测试供应商A", item.ShopName)
	})

	t.Run("添加已存在商品到购物车", func(t *testing.T) {
		// 先添加一次
		req := models.AddToCartRequest{
			ProductID: 2,
			Count:     1,
		}
		_, err := cartService.AddToCart(req)
		require.NoError(t, err)

		// 再添加一次相同商品
		req2 := models.AddToCartRequest{
			ProductID: 2,
			Count:     3,
		}

		item, err := cartService.AddToCart(req2)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, 2, item.ProductID)
		assert.Equal(t, 4, item.Count) // 1 + 3 = 4
		assert.Equal(t, 25.00, item.Price)
		assert.Equal(t, 100.0, item.TotalPrice) // 25 * 4 = 100
	})

	t.Run("添加不存在的商品", func(t *testing.T) {
		req := models.AddToCartRequest{
			ProductID: 999,
			Count:     1,
		}

		item, err := cartService.AddToCart(req)

		assert.Error(t, err)
		assert.Nil(t, item)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}

func TestCartService_UpdateCartItem(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	cartService := services.NewCartService(db)

	// 先添加商品到购物车
	req := models.AddToCartRequest{
		ProductID: 1,
		Count:     2,
	}
	item, err := cartService.AddToCart(req)
	require.NoError(t, err)

	t.Run("更新商品数量", func(t *testing.T) {
		err := cartService.UpdateCartItem(item.ID, 5)

		assert.NoError(t, err)

		// 验证更新结果
		response, err := cartService.GetCart()
		require.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, 5, response.Items[0].Count)
		assert.Equal(t, 52.5, response.Items[0].TotalPrice) // 10.5 * 5 = 52.5
	})

	t.Run("设置数量为0删除商品", func(t *testing.T) {
		err := cartService.UpdateCartItem(item.ID, 0)

		assert.NoError(t, err)

		// 验证商品被删除
		response, err := cartService.GetCart()
		require.NoError(t, err)
		assert.Empty(t, response.Items)
	})

	t.Run("更新不存在的购物车商品", func(t *testing.T) {
		err := cartService.UpdateCartItem(999, 1)

		assert.Error(t, err)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}

func TestCartService_DeleteCartItem(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	cartService := services.NewCartService(db)

	// 先添加商品到购物车
	req := models.AddToCartRequest{
		ProductID: 1,
		Count:     2,
	}
	item, err := cartService.AddToCart(req)
	require.NoError(t, err)

	t.Run("删除存在的购物车商品", func(t *testing.T) {
		err := cartService.DeleteCartItem(item.ID)

		assert.NoError(t, err)

		// 验证商品被删除
		response, err := cartService.GetCart()
		require.NoError(t, err)
		assert.Empty(t, response.Items)
	})

	t.Run("删除不存在的购物车商品", func(t *testing.T) {
		err := cartService.DeleteCartItem(999)

		assert.Error(t, err)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}

func TestCartService_ClearCart(t *testing.T) {
	// 设置测试数据库
	db, err := testdata.SetupTestDB()
	require.NoError(t, err)

	// 创建测试数据
	err = testdata.SeedTestData(db)
	require.NoError(t, err)

	// 创建服务实例
	cartService := services.NewCartService(db)

	// 先添加几个商品到购物车
	items := []models.AddToCartRequest{
		{ProductID: 1, Count: 2},
		{ProductID: 2, Count: 3},
		{ProductID: 3, Count: 1},
	}

	for _, item := range items {
		_, err := cartService.AddToCart(item)
		require.NoError(t, err)
	}

	// 验证购物车有商品
	response, err := cartService.GetCart()
	require.NoError(t, err)
	assert.Len(t, response.Items, 3)

	t.Run("清空购物车", func(t *testing.T) {
		err := cartService.ClearCart()

		assert.NoError(t, err)

		// 验证购物车被清空
		response, err := cartService.GetCart()
		require.NoError(t, err)
		assert.Empty(t, response.Items)
		assert.Equal(t, 0, response.Summary.TotalItems)
		assert.Equal(t, 0.0, response.Summary.TotalPrice)
	})

	// 清理测试数据
	err = testdata.CleanupTestDB(db)
	require.NoError(t, err)
}
