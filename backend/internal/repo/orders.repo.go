package repo

import (
	"context"
	"errors"
	"log"

	"github.com/almarufh/mini-er/internal/models"
	"github.com/jackc/pgx/v5"
)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

func (pr *OrderRepo) GetProductBySKUForUpdateRepo(c context.Context, db Database, sku string) (*models.Product, error) {
	query := `
		SELECT id, sku, name, price, stock, is_active, created_at 
		FROM products 
		WHERE sku = $1 
		FOR UPDATE;
	`

	var product models.Product
	err := db.QueryRow(c, query, sku).Scan(
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
		log.Printf("[GetProductBySKUForUpdateRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to lock product by sku")
	}

	return &product, nil
}

func (or *OrderRepo) CreateOrderRepo(c context.Context, db Database, order *models.SalesOrder) (*models.SalesOrder, error) {
	query := `
		INSERT INTO sales_orders (customer, status, total)
		VALUES ($1, $2, $3)
		RETURNING id, customer, status, total, created_at;
	`

	var created models.SalesOrder
	err := db.QueryRow(c, query, order.Customer, order.Status, order.Total).Scan(
		&created.ID,
		&created.Customer,
		&created.Status,
		&created.Total,
		&created.CreatedAt,
	)

	if err != nil {
		log.Printf("[CreateOrderRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to insert sales order")
	}

	return &created, nil
}

func (or *OrderRepo) CreateOrderItemRepo(c context.Context, db Database, item *models.SalesOrderItem) (*models.SalesOrderItem, error) {
	query := `
		INSERT INTO sales_order_items (order_id, sku, qty, unit_price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, order_id, sku, qty, unit_price;
	`

	var created models.SalesOrderItem
	err := db.QueryRow(c, query, item.OrderID, item.SKU, item.Qty, item.UnitPrice).Scan(
		&created.ID,
		&created.OrderID,
		&created.SKU,
		&created.Qty,
		&created.UnitPrice,
	)

	if err != nil {
		log.Printf("[CreateOrderItemRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to insert sales order item")
	}

	return &created, nil
}

func (or *OrderRepo) GetOrderForUpdateRepo(c context.Context, db Database, id int) (*models.SalesOrder, error) {
	query := `
		SELECT id, customer, status, total, created_at, cancelled_at
		FROM sales_orders
		WHERE id = $1
		FOR UPDATE;
	`

	var order models.SalesOrder
	err := db.QueryRow(c, query, id).Scan(
		&order.ID,
		&order.Customer,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ORDER_NOT_FOUND")
		}
		log.Printf("[GetOrderForUpdateRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to lock sales order")
	}

	return &order, nil
}

func (or *OrderRepo) GetOrderItemsByOrderIDRepo(c context.Context, db Database, orderID int) ([]models.SalesOrderItem, error) {
	query := `
		SELECT id, order_id, sku, qty, unit_price
		FROM sales_order_items
		WHERE order_id = $1;
	`

	rows, err := db.Query(c, query, orderID)
	if err != nil {
		log.Printf("[GetOrderItemsByOrderIDRepo] Query : %s", err.Error())
		return nil, errors.New("failed to query order items")
	}
	defer rows.Close()

	items := make([]models.SalesOrderItem, 0)
	for rows.Next() {
		var item models.SalesOrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.SKU,
			&item.Qty,
			&item.UnitPrice,
		)
		if err != nil {
			log.Printf("[GetOrderItemsByOrderIDRepo] Scan : %s", err.Error())
			return nil, errors.New("failed to scan order item row")
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetOrderItemsByOrderIDRepo] rows.Err : %s", err.Error())
		return nil, errors.New("error reading order item rows")
	}

	return items, nil
}

func (or *OrderRepo) CancelOrderRepo(c context.Context, db Database, id int) (*models.SalesOrder, error) {
	query := `
		UPDATE sales_orders
		SET status = 'CANCELLED', cancelled_at = NOW()
		WHERE id = $1
		RETURNING id, customer, status, total, created_at, cancelled_at;
	`

	var order models.SalesOrder
	err := db.QueryRow(c, query, id).Scan(
		&order.ID,
		&order.Customer,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.CancelledAt,
	)

	if err != nil {
		log.Printf("[CancelOrderRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to update order status to cancelled")
	}

	return &order, nil
}

func (or *OrderRepo) GetOrderByIDRepo(c context.Context, db Database, id int) (*models.SalesOrder, error) {
	query := `
		SELECT id, customer, status, total, created_at, cancelled_at
		FROM sales_orders
		WHERE id = $1;
	`

	var order models.SalesOrder
	err := db.QueryRow(c, query, id).Scan(
		&order.ID,
		&order.Customer,
		&order.Status,
		&order.Total,
		&order.CreatedAt,
		&order.CancelledAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("ORDER_NOT_FOUND")
		}
		log.Printf("[GetOrderByIDRepo] Scan : %s", err.Error())
		return nil, errors.New("failed to query sales order")
	}

	return &order, nil
}
