CREATE TABLE IF NOT EXISTS sales_order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES sales_orders(id),
    sku TEXT NOT NULL REFERENCES products(sku),
    qty INTEGER NOT NULL CHECK (qty > 0),
    unit_price INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sales_order_items_order_id ON sales_order_items (order_id);