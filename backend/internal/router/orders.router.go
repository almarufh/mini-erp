package router

func (r *Router) OrdersRouter() {
	router := r.api.Group("orders")
	router.POST("", r.oc.CreateOrderControl)
	router.POST("/:id/cancel", r.oc.CancelOrderControl)
	router.GET("/:id")
}
