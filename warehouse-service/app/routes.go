package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App, c *Container) {
	api := app.Group("/api/v1")

	warehouse := api.Group("/warehouses")
	warehouse.Post("/", c.WarehouseController.CreateWarehouse)
	warehouse.Get("/", c.WarehouseController.GetAllWarehouse)
	warehouse.Get("/:id", c.WarehouseController.GetWarehouseByID)
	warehouse.Put("/:id", c.WarehouseController.UpdateWarehouse)
	warehouse.Delete("/:id", c.WarehouseController.DeleteWarehouse)

	warehouseProduct := api.Group("/warehouse-products")
	warehouseProduct.Post("/:warehouse_id", c.WarehouseProductController.CreateWarehouseProduct)
	warehouseProduct.Get("/:warehouse_id", c.WarehouseProductController.GetDetailWarehouse)
	warehouseProduct.Get("/:warehouse_id/detail/:product_id", c.WarehouseProductController.GetWarehouseProductByWarehouseIDAndProductID)
	warehouseProduct.Put("/detail/:warehouse_product_id", c.WarehouseProductController.UpdateWarehouseProduct)
	warehouseProduct.Delete("/detail/:warehouse_product_id", c.WarehouseProductController.DeleteWarehouseProduct)
	warehouseProduct.Delete("/detail/products/:product_id", c.WarehouseProductController.DeleteAllWarehouseProductByProductID)
	warehouseProduct.Get("/detail/products/:product_id/total-stock", c.WarehouseProductController.GetProductTotalStock)
	warehouseProduct.Get("/detail/products/:product_id", c.WarehouseProductController.GetWarehouseProductByProductID)
	warehouseProduct.Get("/detail/products/:product_id/warehouse", c.WarehouseProductController.GetDetailWarehouseProductByID)

	api.Post("upload-warehouse", c.UploadController.UploadPhoto)
}
