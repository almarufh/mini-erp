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
