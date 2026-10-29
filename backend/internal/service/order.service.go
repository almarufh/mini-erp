package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/almarufh/mini-er/internal/dto"
	"github.com/almarufh/mini-er/internal/models"
	"github.com/almarufh/mini-er/internal/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	db *pgxpool.Pool
	or *repo.OrderRepo
	pr *repo.ProductRepo
}

func NewOrderService(db *pgxpool.Pool, or *repo.OrderRepo, pr *repo.ProductRepo) *OrderService {
	return &OrderService{
		db: db,
		or: or,
		pr: pr,
	}
}

func (os *OrderService) CreateOrderService(c context.Context, req *dto.CreateOrderReq) (*dto.SalesOrderResponse, error) {
	tx, err := os.db.Begin(c)
	if err != nil {
		return nil, errors.New("failed to start database transaction")
	}

	defer func() {
		_ = tx.Rollback(c)
	}()

	calculatedTotal := 0
	orderItemsToInsert := make([]models.SalesOrderItem, 0)

	for _, itemReq := range req.Items {
		product, err := os.or.GetProductBySKUForUpdateRepo(c, tx, itemReq.SKU)
		if err != nil {
			if err.Error() == "PRODUCT_NOT_FOUND" {
				return nil, fmt.Errorf("SKU_NOT_FOUND:%s", itemReq.SKU)
			}
			return nil, err
		}

		if product.Stock < itemReq.Qty {
			return nil, fmt.Errorf("INSUFFICIENT_STOCK:%s:%d", product.SKU, product.Stock)
		}

		itemTotal := product.Price * itemReq.Qty
		calculatedTotal += itemTotal

		newStock := product.Stock - itemReq.Qty
		if err := os.pr.UpdateStockRepo(c, tx, product.ID, newStock); err != nil {
			return nil, err
		}

		note := fmt.Sprintf("Sales Order for Customer: %s", req.Customer)
		movement := &models.StockMovement{
			ProductID: product.ID,
			Type:      "out",
			Qty:       itemReq.Qty,
			Note:      &note,
		}
		if err := os.pr.CreateStockMovementRepo(c, tx, movement); err != nil {
			return nil, err
		}

		orderItemsToInsert = append(orderItemsToInsert, models.SalesOrderItem{
			SKU:       itemReq.SKU,
			Qty:       itemReq.Qty,
			UnitPrice: product.Price,
		})
	}

	masterOrder := &models.SalesOrder{
		Customer: req.Customer,
		Status:   "CREATED",
		Total:    calculatedTotal,
	}

	createdOrder, err := os.or.CreateOrderRepo(c, tx, masterOrder)
	if err != nil {
		return nil, err
	}

	responseItems := make([]dto.SalesOrderItemResponse, 0)
	for _, item := range orderItemsToInsert {
		item.OrderID = createdOrder.ID
		createdItem, err := os.or.CreateOrderItemRepo(c, tx, &item)
		if err != nil {
			return nil, err
		}

		responseItems = append(responseItems, dto.SalesOrderItemResponse{
			ID:        createdItem.ID,
			OrderID:   createdItem.OrderID,
			SKU:       createdItem.SKU,
			Qty:       createdItem.Qty,
			UnitPrice: createdItem.UnitPrice,
		})
	}

	if err := tx.Commit(c); err != nil {
		return nil, errors.New("failed to commit order transaction")
	}

	return &dto.SalesOrderResponse{
		ID:          createdOrder.ID,
		Customer:    createdOrder.Customer,
		Status:      createdOrder.Status,
		Total:       createdOrder.Total,
		CreatedAt:   createdOrder.CreatedAt,
		CancelledAt: nil,
		Items:       responseItems,
	}, nil
}

func (os *OrderService) CancelOrderService(c context.Context, id int) (*dto.SalesOrderResponse, error) {
	tx, err := os.db.Begin(c)
	if err != nil {
		return nil, errors.New("failed to start database transaction")
	}

	defer func() {
		_ = tx.Rollback(c)
	}()

	order, err := os.or.GetOrderForUpdateRepo(c, tx, id)
	if err != nil {
		return nil, err
	}

	if order.Status == "CANCELLED" {
		return nil, errors.New("ORDER_ALREADY_CANCELLED")
	}

	items, err := os.or.GetOrderItemsByOrderIDRepo(c, tx, id)
	if err != nil {
		return nil, err
	}

	responseItems := make([]dto.SalesOrderItemResponse, 0)
	for _, item := range items {
		product, err := os.or.GetProductBySKUForUpdateRepo(c, tx, item.SKU)
		if err != nil {
			return nil, err
		}

		restoredStock := product.Stock + item.Qty
		if err := os.pr.UpdateStockRepo(c, tx, product.ID, restoredStock); err != nil {
			return nil, err
		}

		note := fmt.Sprintf("Restored stock from cancelled Order #%d", id)
		movement := &models.StockMovement{
			ProductID: product.ID,
			Type:      "in",
			Qty:       item.Qty,
			Note:      &note,
		}
		if err := os.pr.CreateStockMovementRepo(c, tx, movement); err != nil {
			return nil, err
		}

		responseItems = append(responseItems, dto.SalesOrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			SKU:       item.SKU,
			Qty:       item.Qty,
			UnitPrice: item.UnitPrice,
		})
	}

	cancelledOrder, err := os.or.CancelOrderRepo(c, tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, errors.New("failed to commit order cancellation transaction")
	}

	return &dto.SalesOrderResponse{
		ID:          cancelledOrder.ID,
		Customer:    cancelledOrder.Customer,
		Status:      cancelledOrder.Status,
		Total:       cancelledOrder.Total,
		CreatedAt:   cancelledOrder.CreatedAt,
		CancelledAt: cancelledOrder.CancelledAt,
		Items:       responseItems,
	}, nil
}

func (os *OrderService) GetOrderByIDService(c context.Context, id int) (*dto.SalesOrderResponse, error) {
	order, err := os.or.GetOrderByIDRepo(c, os.db, id)
	if err != nil {
		return nil, err
	}

	items, err := os.or.GetOrderItemsByOrderIDRepo(c, os.db, id)
	if err != nil {
		return nil, err
	}

	responseItems := make([]dto.SalesOrderItemResponse, 0)
	for _, item := range items {
		responseItems = append(responseItems, dto.SalesOrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			SKU:       item.SKU,
			Qty:       item.Qty,
			UnitPrice: item.UnitPrice,
		})
	}

	return &dto.SalesOrderResponse{
		ID:          order.ID,
		Customer:    order.Customer,
		Status:      order.Status,
		Total:       order.Total,
		CreatedAt:   order.CreatedAt,
		CancelledAt: order.CancelledAt,
		Items:       responseItems,
	}, nil
}

func (os *OrderService) GetOrdersService(c context.Context) ([]dto.SalesOrderResponse, error) {
	orders, err := os.or.GetOrdersRepo(c, os.db)
	if err != nil {
		return nil, err
	}

	res := make([]dto.SalesOrderResponse, 0)
	for _, o := range orders {
		res = append(res, dto.SalesOrderResponse{
			ID:          o.ID,
			Customer:    o.Customer,
			Status:      o.Status,
			Total:       o.Total,
			CreatedAt:   o.CreatedAt,
			CancelledAt: o.CancelledAt,
		})
	}

	return res, nil
}
