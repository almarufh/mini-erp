package router

import (
	"github.com/almarufh/mini-er/internal/control"
	"github.com/almarufh/mini-er/internal/middleware"
	"github.com/almarufh/mini-er/internal/repo"
	"github.com/almarufh/mini-er/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	api  *gin.RouterGroup
	pool *pgxpool.Pool
	m    *middleware.Middleware
	pc   *control.ProductControl
	oc   *control.OrderControl
}

func NewRouter(router *gin.Engine, pool *pgxpool.Pool, allowedOrigins []string) *Router {
	m := middleware.InitMiddleware(allowedOrigins)
	router.Use(m.Cors)
	api := router.Group("api")

	pr := repo.NewProductRepo()
	ps := service.NewProductService(pool, pr)
	pc := control.NewProductControl(ps)

	or := repo.NewOrderRepo()
	os := service.NewOrderService(pool, or, pr)
	oc := control.NewOrderControl(os)
	return &Router{
		pool: pool,
		m:    m,
		api:  api,
		pc:   pc,
		oc:   oc,
	}
}

func (r *Router) Connect() {
	r.Swagger()
	r.ProductsRouter()
	r.OrdersRouter()
}
