package router

func (p *Router) OrdersRouter() {
	router := p.api.Group("orders")
	router.POST("")
	router.POST("/:id/cancel")
	router.GET("/:id")
}
