import {
  createAsyncThunk,
  createSlice,
  type PayloadAction,
} from '@reduxjs/toolkit';
import { orderApi } from '../../api/orderApi';
import type { CreateOrderReq, SalesOrder } from '../../types';

interface OrderState {
  orders: SalesOrder[];
  currentOrder: SalesOrder | null;
  isLoading: boolean;
  error: string | null;
  errorEnvelope: any | null;
}

const initialState: OrderState = {
  orders: [],
  currentOrder: null,
  isLoading: false,
  error: null,
  errorEnvelope: null,
};

export const createSalesOrder = createAsyncThunk(
  'orders/createOrder',
  async (payload: CreateOrderReq, { rejectWithValue }) => {
    try {
      const response = await orderApi.createOrder(payload);
      return response.results!;
    } catch (err: any) {
      if (err.response?.status === 409) {
        return rejectWithValue(err.response.data);
      }
      return rejectWithValue({
        message: err.response?.data?.message || 'Gagal membuat sales order',
      });
    }
  }
);

export const cancelSalesOrder = createAsyncThunk(
  'orders/cancelOrder',
  async (id: number, { rejectWithValue }) => {
    try {
      const response = await orderApi.cancelOrder(id);
      return response.results!;
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data?.message || 'Gagal membatalkan sales order'
      );
    }
  }
);

export const fetchOrders = createAsyncThunk(
  'orders/fetchOrders',
  async (_, { rejectWithValue }) => {
    try {
      const response = await orderApi.getOrders();
      return response.results || [];
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data?.message || 'Gagal mengambil data order'
      );
    }
  }
);

const orderSlice = createSlice({
  name: 'orders',
  initialState,
  reducers: {
    clearOrderError: (state) => {
      state.error = null;
      state.errorEnvelope = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchOrders.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(
        fetchOrders.fulfilled,
        (state, action: PayloadAction<SalesOrder[]>) => {
          state.isLoading = false;
          state.orders = action.payload;
        }
      )
      .addCase(fetchOrders.rejected, (state, action: any) => {
        state.isLoading = false;
        state.error = action.payload || 'Gagal mengambil data order';
      })
      .addCase(createSalesOrder.pending, (state) => {
        state.isLoading = true;
        state.error = null;
        state.errorEnvelope = null;
      })
      .addCase(
        createSalesOrder.fulfilled,
        (state, action: PayloadAction<SalesOrder>) => {
          state.isLoading = false;
          state.currentOrder = action.payload;
          state.orders.unshift(action.payload);
        }
      )
      .addCase(createSalesOrder.rejected, (state, action: any) => {
        state.isLoading = false;
        if (action.payload?.error) {
          state.errorEnvelope = action.payload.error;
        } else {
          state.error = action.payload?.message || 'Terjadi kesalahan';
        }
      })
      .addCase(
        cancelSalesOrder.fulfilled,
        (state, action: PayloadAction<SalesOrder>) => {
          const index = state.orders.findIndex(
            (o) => o.id === action.payload.id
          );
          if (index !== -1) {
            state.orders[index] = action.payload;
          }
        }
      );
  },
});

export const { clearOrderError } = orderSlice.actions;
export default orderSlice.reducer;
