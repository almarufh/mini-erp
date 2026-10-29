package router

func (r *Router) ProductsRouter() {
	router := r.api.Group("products")

	router.GET("", r.pc.GetProductsControl)
	router.POST("", r.pc.CreateProductControl)
	router.PUT("/:id")
	router.DELETE("/:id")
	router.POST("/:id/stock")
	router.POST("/:id/movements")
}
