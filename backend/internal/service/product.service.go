package service

import (
	"context"
	"fmt"
	"math"

	"github.com/almarufh/mini-er/internal/dto"
	"github.com/almarufh/mini-er/internal/models"
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

func (ps *ProductService) CreateProductService(c context.Context, req *dto.CreateProductReq) (*dto.ProductResponse, error) {

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	productModel := &models.Product{
		SKU:      req.SKU,
		Name:     req.Name,
		Price:    req.Price,
		Stock:    req.Stock,
		IsActive: isActive,
	}

	createdProduct, err := ps.pr.CreateProductRepo(c, ps.db, productModel)
	if err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:             createdProduct.ID,
		SKU:            createdProduct.SKU,
		Name:           createdProduct.Name,
		Price:          createdProduct.Price,
		Stock:          createdProduct.Stock,
		IsActive:       createdProduct.IsActive,
		CreatedAt:      createdProduct.CreatedAt,
		LastMovementAt: nil,
	}, nil
}

func (ps *ProductService) UpdateProductService(c context.Context, id int, req *dto.UpdateProductReq) (*dto.ProductResponse, error) {
	updatedProduct, err := ps.pr.UpdateProductRepo(c, ps.db, id, req)
	if err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:             updatedProduct.ID,
		SKU:            updatedProduct.SKU,
		Name:           updatedProduct.Name,
		Price:          updatedProduct.Price,
		Stock:          updatedProduct.Stock,
		IsActive:       updatedProduct.IsActive,
		CreatedAt:      updatedProduct.CreatedAt,
		LastMovementAt: nil,
	}, nil
}

func (ps *ProductService) DeleteProductService(c context.Context, id int) error {
	return ps.pr.DeleteProductRepo(c, ps.db, id)
}
