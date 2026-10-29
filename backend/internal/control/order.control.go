package control

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/almarufh/mini-er/internal/dto"
	"github.com/almarufh/mini-er/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderControl struct {
	os *service.OrderService
}

func NewOrderControl(os *service.OrderService) *OrderControl {
	return &OrderControl{
		os: os,
	}
}

// Create Sales Order
// @Summary      Pencatatan Sales Order Baru
// @Description  Membuat order penjualan baru, memotong stok produk secara atomik & menghitung total harga
// @Tags         Sales Orders
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateOrderReq  true  "Payload Create Order"
// @Success      201   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Failure      409   {object}  dto.ErrorEnvelope  "INSUFFICIENT_STOCK"
// @Failure      500   {object}  dto.Response
// @Router       /api/orders [post]
func (oc *OrderControl) CreateOrderControl(c *gin.Context) {
	var req dto.CreateOrderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validasi payload order gagal: " + err.Error(),
		})
		return
	}

	res, err := oc.os.CreateOrderService(c.Request.Context(), &req)
	if err != nil {
		errStr := err.Error()

		if strings.HasPrefix(errStr, "INSUFFICIENT_STOCK:") {
			parts := strings.Split(errStr, ":")
			if len(parts) == 3 {
				sku := parts[1]
				available, _ := strconv.Atoi(parts[2])

				c.JSON(http.StatusConflict, gin.H{
					"error": gin.H{
						"code":      "INSUFFICIENT_STOCK",
						"sku":       sku,
						"available": available,
					},
				})
				return
			}
		}

		if strings.HasPrefix(errStr, "SKU_NOT_FOUND:") {
			parts := strings.Split(errStr, ":")
			c.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "SKU produk tidak ditemukan: " + parts[1],
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal membuat sales order: " + errStr,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "Sales order berhasil dibuat",
		Results: res,
	})
}
