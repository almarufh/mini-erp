export interface Product {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  created_at: string;
  last_movement_at?: string | null;
}

export interface PaginatedProductsResponse {
  data: Product[];
  page: number;
  limit: number;
  total_data: number;
  total_page: number;
  next_page: string;
  prev_page: string;
}

export interface CreateProductReq {
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active?: boolean;
}

export interface UpdateProductReq {
  name?: string;
  price?: number;
  is_active?: boolean;
}

export interface AdjustStockReq {
  type: 'in' | 'out';
  qty: number;
  note?: string;
}

export interface StockMovement {
  id: number;
  product_id: number;
  type: 'in' | 'out';
  qty: number;
  note?: string | null;
  created_at: string;
}

export interface ProductMovementsResponse {
  data: StockMovement[];
}

export interface CreateOrderItemReq {
  sku: string;
  qty: number;
}

export interface CreateOrderReq {
  customer: string;
  items: CreateOrderItemReq[];
}

export interface SalesOrderItem {
  id: number;
  order_id: number;
  sku: string;
  qty: number;
  unit_price: number;
}

export interface SalesOrder {
  id: number;
  customer: string;
  status: 'CREATED' | 'CANCELLED';
  total: number;
  created_at: string;
  cancelled_at?: string | null;
  items: SalesOrderItem[];
}

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  results?: T;
}

export interface ErrorEnvelope {
  error: {
    code: string;
    sku?: string;
    available?: number;
  };
}
