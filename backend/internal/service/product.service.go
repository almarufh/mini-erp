package service

import (
	"context"
	"fmt"
	"math"

	"github.com/almarufh/mini-er/internal/dto"
	"github.com/almarufh/mini-er/internal/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	db *pgxpool.Pool
	pr *repo.ProductRepo
}

func NewProductService(db *pgxpool.Pool, pr *repo.ProductRepo) *ProductService {
	return &ProductService{
		db: db,
		pr: pr,
	}
}

func (ps *ProductService) GetProductsService(c context.Context, search string, page, limit int) (*dto.PaginatedProductsResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := ps.pr.GetProductsRepo(c, ps.db, search, limit, offset)
	if err != nil {
		return nil, err
	}

	productList := make([]dto.ProductResponse, 0)
	if products != nil {
		for _, p := range *products {
			productList = append(productList, dto.ProductResponse{
				ID:             p.ID,
				SKU:            p.SKU,
				Name:           p.Name,
				Price:          p.Price,
				Stock:          p.Stock,
				IsActive:       p.IsActive,
				CreatedAt:      p.CreatedAt,
				LastMovementAt: p.LastMovementAt,
			})
		}
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	nextPage := ""
	if page < totalPage {
		nextPage = fmt.Sprintf("?page=%d&limit=%d&search=%s", page+1, limit, search)
	}

	prevPage := ""
	if page > 1 && page <= totalPage+1 {
		prevPage = fmt.Sprintf("?page=%d&limit=%d&search=%s", page-1, limit, search)
	}

	return &dto.PaginatedProductsResponse{
		Data:      productList,
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}
