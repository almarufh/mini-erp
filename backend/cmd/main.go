package main

import (
	"fmt"

	"github.com/almarufh/mini-er/internal/config"
	"github.com/almarufh/mini-er/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.LoadEnv()
	pool := conf.POSTGRES.Connect()
	defer pool.Close()
	server := gin.Default()
	router.NewRouter(server, pool, conf.ALLOWED_ORIGINS).Connect()
	server.Run(fmt.Sprintf("%s:%s", conf.HOST, conf.PORT))
}
