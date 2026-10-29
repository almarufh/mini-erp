import type { ApiResponse, CreateOrderReq, SalesOrder } from '../types';
import axiosInstance from './axiosInstance';

export const orderApi = {
  createOrder: async (data: CreateOrderReq) => {
    const response = await axiosInstance.post<ApiResponse<SalesOrder>>(
      '/api/orders',
      data
    );
    return response.data;
  },

  getOrderById: async (id: number) => {
    const response = await axiosInstance.get<ApiResponse<SalesOrder>>(
      `/api/orders/${id}`
    );
    return response.data;
  },

  cancelOrder: async (id: number) => {
    const response = await axiosInstance.post<ApiResponse<SalesOrder>>(
      `/api/orders/${id}/cancel`
    );
    return response.data;
  },

  getOrders: async () => {
    const response =
      await axiosInstance.get<ApiResponse<SalesOrder[]>>('/api/orders');
    return response.data;
  },
};
