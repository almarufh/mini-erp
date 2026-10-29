package router

import (
	"github.com/almarufh/mini-er/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	api  *gin.RouterGroup
	pool *pgxpool.Pool
	m    *middleware.Middleware
}

func NewRouter(router *gin.Engine, pool *pgxpool.Pool, allowedOrigins []string) *Router {
	api := router.Group("api")
	return &Router{
		pool: pool,
		m:    middleware.InitMiddleware(allowedOrigins),
		api:  api,
	}
}

func (r *Router) Connect() {
	r.ProductsRouter()
	r.OrdersRouter()
}
