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

func (pr *ProductRepo) DeleteProductRepo(c context.Context, db Database, id int) error {
	checkQuery := `SELECT EXISTS (SELECT 1 FROM stock_movements WHERE product_id = $1);`
	var hasHistory bool
	err := db.QueryRow(c, checkQuery, id).Scan(&hasHistory)
	if err != nil {
		log.Printf("[DeleteProductRepo] QueryRow check history : %s", err.Error())
		return errors.New("failed to check product history")
	}

	if hasHistory {
		return errors.New("PRODUCT_HAS_HISTORY")
	}

	deleteQuery := `DELETE FROM products WHERE id = $1;`
	tag, err := db.Exec(c, deleteQuery, id)
	if err != nil {
		log.Printf("[DeleteProductRepo] Exec delete : %s", err.Error())
		return errors.New("failed to delete product")
	}

	if tag.RowsAffected() == 0 {
		return errors.New("PRODUCT_NOT_FOUND")
	}

	return nil
}

func (pr *ProductRepo) GetProductForUpdateRepo(c context.Context, db Database, id int) (*models.Product, error) {
	query := `
		SELECT id, sku, name, price, stock, is_active, created_at 
		FROM products 
		WHERE id = $1 
		FOR UPDATE;
	`

	var product models.Product
	err := db.QueryRow(c, query, id).Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("PRODUCT_NOT_FOUND")
		}
		log.Printf("[GetProductForUpdateRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to lock product row")
	}

	return &product, nil
}

func (pr *ProductRepo) UpdateStockRepo(c context.Context, db Database, id int, newStock int) error {
	query := `UPDATE products SET stock = $1 WHERE id = $2;`
	_, err := db.Exec(c, query, newStock, id)
	if err != nil {
		log.Printf("[UpdateStockRepo] Exec : %s", err.Error())
		return errors.New("failed to update product stock")
	}
	return nil
}

func (pr *ProductRepo) CreateStockMovementRepo(c context.Context, db Database, movement *models.StockMovement) error {
	query := `
		INSERT INTO stock_movements (product_id, type, qty, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at;
	`
	err := db.QueryRow(c, query, movement.ProductID, movement.Type, movement.Qty, movement.Note).Scan(&movement.ID, &movement.CreatedAt)
	if err != nil {
		log.Printf("[CreateStockMovementRepo] Scan : %s", err.Error())
		return errors.New("failed to record stock movement")
	}
	return nil
}
