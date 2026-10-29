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

// Cancel Sales Order
// @Summary      Pembatalan Sales Order
// @Description  Membatalkan sales order, mengembalikan (restore) stok produk, dan mencatat stock movement tipe 'in'
// @Tags         Sales Orders
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      404  {object}  dto.Response
// @Failure      409  {object}  dto.Response  "ORDER_ALREADY_CANCELLED"
// @Failure      500  {object}  dto.Response
// @Router       /api/orders/{id}/cancel [post]
func (oc *OrderControl) CancelOrderControl(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID order tidak valid",
		})
		return
	}

	res, err := oc.os.CancelOrderService(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "ORDER_ALREADY_CANCELLED" {
			c.JSON(http.StatusConflict, dto.Response{
				Success: false,
				Message: "Order sudah pernah dibatalkan sebelumnya",
			})
			return
		}

		if err.Error() == "ORDER_NOT_FOUND" {
			c.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "Sales order tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal membatalkan sales order: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Sales order berhasil dibatalkan dan stok telah dikembalikan",
		Results: res,
	})
}

// Get Sales Order By ID
// @Summary      Detail Sales Order
// @Description  Mengambil detail sales order berdasarkan ID beserta daftar item, total harga, dan statusnya
// @Tags         Sales Orders
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.Response
// @Failure      404  {object}  dto.Response
// @Failure      500  {object}  dto.Response
// @Router       /api/orders/{id} [get]
func (oc *OrderControl) GetOrderByIDControl(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID order tidak valid",
		})
		return
	}

	res, err := oc.os.GetOrderByIDService(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "ORDER_NOT_FOUND" {
			c.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "Sales order tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal mengambil detail sales order: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil mengambil detail sales order",
		Results: res,
	})
}
