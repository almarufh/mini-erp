import type {
  AdjustStockReq,
  ApiResponse,
  CreateProductReq,
  PaginatedProductsResponse,
  Product,
  ProductMovementsResponse,
  UpdateProductReq,
} from '../types';

import axiosInstance from './axiosInstance';

export const productApi = {
  getProducts: async (search = '', page = 1, limit = 10) => {
    const response = await axiosInstance.get<
      ApiResponse<PaginatedProductsResponse>
    >('/api/products', {
      params: { search, page, limit },
    });
    return response.data;
  },

  createProduct: async (data: CreateProductReq) => {
    const response = await axiosInstance.post<ApiResponse<Product>>(
      '/api/products',
      data
    );
    return response.data;
  },

  updateProduct: async (id: number, data: UpdateProductReq) => {
    const response = await axiosInstance.put<ApiResponse<Product>>(
      `/api/products/${id}`,
      data
    );
    return response.data;
  },

  deleteProduct: async (id: number) => {
    const response = await axiosInstance.delete(`/api/products/${id}`);
    return response.data;
  },

  adjustStock: async (id: number, data: AdjustStockReq) => {
    const response = await axiosInstance.post<ApiResponse>(
      `/api/products/${id}/stock`,
      data
    );
    return response.data;
  },

  getProductMovements: async (id: number) => {
    const response = await axiosInstance.get<ProductMovementsResponse>(
      `/api/products/${id}/movements`
    );
    return response.data;
  },
};
