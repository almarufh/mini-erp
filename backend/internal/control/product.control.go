package control

import (
	"net/http"
	"strconv"

	"github.com/almarufh/mini-er/internal/dto"
	"github.com/almarufh/mini-er/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductControl struct {
	ps *service.ProductService
}

func NewProductControl(ps *service.ProductService) *ProductControl {
	return &ProductControl{
		ps: ps,
	}
}

// Get Products
// @Summary      Ambil daftar produk
// @Description  Mengambil list produk dengan filter pencarian (SKU/Nama) dan pagination (page & limit)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        search   query     string  false  "Search by SKU or Name (Case-Insensitive)"
// @Param        page     query     int     false  "Page Number" default(1)
// @Param        limit    query     int     false  "Limit per Page" default(10)
// @Success      200      {object}  dto.PaginatedProductsResponse
// @Failure      500      {object}  dto.Response
// @Router       /api/products [get]
func (pc *ProductControl) GetProductsControl(c *gin.Context) {
	search := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	products, err := pc.ps.GetProductsService(c.Request.Context(), search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal mengambil data produk: ",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "get data success",
		Results: *products,
	})
}

// Create Product
// @Summary      Tambah Produk Baru
// @Description  Membuat data produk baru dengan validasi SKU unik, Nama, dan Harga >= 0
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateProductReq  true  "Payload Create Product"
// @Success      201   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Failure      409   {object}  dto.Response
// @Failure      500   {object}  dto.Response
// @Router       /api/products [post]
func (pc *ProductControl) CreateProductControl(c *gin.Context) {
	var req dto.CreateProductReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	if req.Price < 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Harga produk tidak boleh kurang dari 0",
		})
		return
	}

	res, err := pc.ps.CreateProductService(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "SKU_ALREADY_EXISTS" {
			c.JSON(http.StatusConflict, dto.Response{
				Success: false,
				Message: "SKU produk sudah terdaftar, gunakan SKU lain",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal membuat produk: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "Produk berhasil dibuat",
		Results: res,
	})
}

// Update Product
// @Summary      Update Data Produk
// @Description  Memperbarui nama, harga, dan status aktif produk berdasarkan ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id    path      int                   true  "Product ID"
// @Param        body  body      dto.UpdateProductReq  true  "Payload Update Product"
// @Success      200   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Failure      404   {object}  dto.Response
// @Failure      500   {object}  dto.Response
// @Router       /api/products/{id} [put]
func (pc *ProductControl) UpdateProductControl(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID produk tidak valid",
		})
		return
	}

	var req dto.UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validasi payload gagal: " + err.Error(),
		})
		return
	}

	if req.Price != nil && *req.Price < 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Harga produk tidak boleh kurang dari 0",
		})
		return
	}

	res, err := pc.ps.UpdateProductService(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "PRODUCT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "Produk tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal memperbarui produk: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Produk berhasil diperbarui",
		Results: res,
	})
}

// Delete Product
// @Summary      Hapus Produk
// @Description  Menghapus produk jika belum memiliki riwayat movement stok. Mengembalikan 409 jika memiliki riwayat.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content (Berhasil dihapus)"
// @Failure      400  {object}  dto.Response
// @Failure      404  {object}  dto.Response
// @Failure      409  {object}  dto.ErrorEnvelope  "PRODUCT_HAS_HISTORY"
// @Failure      500  {object}  dto.Response
// @Router       /api/products/{id} [delete]
func (pc *ProductControl) DeleteProductControl(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID produk tidak valid",
		})
		return
	}

	err = pc.ps.DeleteProductService(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "PRODUCT_HAS_HISTORY" {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code": "PRODUCT_HAS_HISTORY",
				},
			})
			return
		}

		if err.Error() == "PRODUCT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "Produk tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal menghapus produk: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Adjust Stock
// @Summary      Penyesuaian Stok Produk (In / Out)
// @Description  Menambah atau mengurangi stok produk secara atomik dalam transaksi database
// @Tags         Stock Movements
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "Product ID"
// @Param        body  body      dto.AdjustStockReq  true  "Payload Adjustment Stock"
// @Success      200   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Failure      404   {object}  dto.Response
// @Failure      409   {object}  dto.ErrorEnvelope  "INSUFFICIENT_STOCK"
// @Failure      500   {object}  dto.Response
// @Router       /api/products/{id}/stock [post]
func (pc *ProductControl) AdjustStockControl(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID produk tidak valid",
		})
		return
	}

	var req dto.AdjustStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validasi payload gagal: " + err.Error(),
		})
		return
	}

	availableStock, err := pc.ps.AdjustStockService(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "INSUFFICIENT_STOCK" {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":      "INSUFFICIENT_STOCK",
					"available": availableStock,
				},
			})
			return
		}

		if err.Error() == "PRODUCT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "Produk tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal menyesuaikan stok: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Penyesuaian stok berhasil disimpan",
	})
}
