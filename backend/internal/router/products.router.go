package router

func (r *Router) ProductsRouter() {
	router := r.api.Group("products")

	router.GET("", r.pc.GetProductsControl)
	router.POST("", r.pc.CreateProductControl)
	router.PUT("/:id", r.pc.UpdateProductControl)
	router.DELETE("/:id", r.pc.DeleteProductControl)
	router.POST("/:id/stock", r.pc.AdjustStockControl)
	router.POST("/:id/movements")
}
