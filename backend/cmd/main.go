package main

import (
	"fmt"
	"os"

	"github.com/almarufh/mini-er/docs"
	"github.com/almarufh/mini-er/internal/config"
	"github.com/almarufh/mini-er/internal/router"
	"github.com/gin-gonic/gin"
)

// @title                       StockFlow Mini ERP API
// @version                     1.0
// @description                 Backend RESTful API Service untuk Manajemen Inventaris, Movement Stok, dan Sales Order (StockFlow Mini ERP System).
// @termsOfService              http://swagger.io/terms/

// @contact.name                Alma'ruf Hidayat
// @contact.email               almarufh@gmail.com

// @license.name                MIT License
// @license.url                 https://opensource.org/licenses/MIT

// @host                        localhost:5000
// @BasePath                    /
// @schemes                     http https

// @tag.name                    Products
// @tag.description             Management katalog produk, harga, dan filter stok

// @tag.description             Penyesuaian stok masuk (in) & keluar (out) serta riwayat movement
// @tag.name                    Stock Movements

// @tag.name                    Sales Orders
// @tag.description             Pencatatan order penjualan atomik dan pembatalan transaksi order
func main() {
	conf := config.LoadEnv()
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("HOST_SWAGG"), os.Getenv("PORT_SWAGG"))
	pool := conf.POSTGRES.Connect()
	defer pool.Close()
	server := gin.Default()
	router.NewRouter(server, pool, conf.ALLOWED_ORIGINS).Connect()
	server.Run(fmt.Sprintf("%s:%s", conf.HOST, conf.PORT))
}
