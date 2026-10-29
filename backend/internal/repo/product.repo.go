package repo

import (
	"context"
	"errors"
	"log"

	"github.com/almarufh/mini-er/internal/dto"
	"github.com/almarufh/mini-er/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type ProductRepo struct{}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

func (pr *ProductRepo) GetProductsRepo(c context.Context, db Database, search string, limit, offset int) (*[]models.Product, int, error) {
	query := `
        SELECT 
            p.id, 
            p.sku, 
            p.name, 
            p.price, 
            p.stock, 
            p.is_active, 
            p.created_at,
            MAX(sm.created_at) AS last_movement_at,
            COUNT(*) OVER() AS total_count
        FROM products p
        LEFT JOIN stock_movements sm ON p.id = sm.product_id
        WHERE 
            ($1::text IS NULL OR $1 = '' OR LOWER(p.sku) LIKE '%' || LOWER($1) || '%' OR LOWER(p.name) LIKE '%' || LOWER($1) || '%')
        GROUP BY p.id
        ORDER BY p.created_at DESC
        LIMIT $2 OFFSET $3;
    `

	rows, err := db.Query(c, query, search, limit, offset)
	if err != nil {
		log.Printf("[GetProductsRepo] Query : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, 0, err
		}
		return nil, 0, errors.New("internal server error")
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	totalCount := 0

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.SKU,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.IsActive,
			&product.CreatedAt,
			&product.LastMovementAt,
			&totalCount,
		)
		if err != nil {
			log.Printf("[GetProductsRepo] rows.Next : %s", err.Error())
			return nil, 0, errors.New("failed to scan query result")
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetProductsRepo] rows.Err : %s", err.Error())
		return nil, 0, errors.New("error reading rows")
	}

	return &products, totalCount, nil
}

func (pr *ProductRepo) CreateProductRepo(c context.Context, db Database, product *models.Product) (*models.Product, error) {
	query := `
		INSERT INTO products (sku, name, price, stock, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, sku, name, price, stock, is_active, created_at;
	`

	var createdProduct models.Product
	err := db.QueryRow(
		c,
		query,
		product.SKU,
		product.Name,
		product.Price,
		product.Stock,
		product.IsActive,
	).Scan(
		&createdProduct.ID,
		&createdProduct.SKU,
		&createdProduct.Name,
		&createdProduct.Price,
		&createdProduct.Stock,
		&createdProduct.IsActive,
		&createdProduct.CreatedAt,
	)

	if err != nil {
		log.Printf("[CreateProductRepo] QueryRow Scan : %s", err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("SKU_ALREADY_EXISTS")
		}

		return nil, errors.New("failed to create product")
	}

	return &createdProduct, nil
}

func (pr *ProductRepo) UpdateProductRepo(c context.Context, db Database, id int, req *dto.UpdateProductReq) (*models.Product, error) {
	query := `
		UPDATE products 
		SET 
			name = COALESCE($1, name),
			price = COALESCE($2, price),
			is_active = COALESCE($3, is_active)
		WHERE id = $4
		RETURNING id, sku, name, price, stock, is_active, created_at;
	`

	var updatedProduct models.Product
	err := db.QueryRow(
		c,
		query,
		req.Name,
		req.Price,
		req.IsActive,
		id,
	).Scan(
		&updatedProduct.ID,
		&updatedProduct.SKU,
		&updatedProduct.Name,
		&updatedProduct.Price,
		&updatedProduct.Stock,
		&updatedProduct.IsActive,
		&updatedProduct.CreatedAt,
	)

	if err != nil {
		log.Printf("[UpdateProductRepo] QueryRow Scan : %s", err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("PRODUCT_NOT_FOUND")
		}
		return nil, errors.New("failed to update product")
	}

	return &updatedProduct, nil
}
