import {
  createAsyncThunk,
  createSlice,
  type PayloadAction,
} from '@reduxjs/toolkit';
import { productApi } from '../../api/productApi';
import type {
  AdjustStockReq,
  CreateProductReq,
  PaginatedProductsResponse,
  Product,
  StockMovement,
} from '../../types';

interface ProductState {
  products: Product[];
  pagination: Omit<PaginatedProductsResponse, 'data'> | null;
  movements: StockMovement[];
  isLoading: boolean;
  error: string | null;
  errorDetails: any | null;
}

const initialState: ProductState = {
  products: [],
  pagination: null,
  movements: [],
  isLoading: false,
  error: null,
  errorDetails: null,
};

export const fetchProducts = createAsyncThunk(
  'products/fetchProducts',
  async (
    { search, page, limit }: { search?: string; page?: number; limit?: number },
    { rejectWithValue }
  ) => {
    try {
      const response = await productApi.getProducts(search, page, limit);
      return response.results!;
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data?.message || 'Gagal mengambil data produk'
      );
    }
  }
);

export const createProduct = createAsyncThunk(
  'products/createProduct',
  async (payload: CreateProductReq, { rejectWithValue }) => {
    try {
      const response = await productApi.createProduct(payload);
      return response.results!;
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data?.message || 'Gagal membuat produk'
      );
    }
  }
);

export const adjustProductStock = createAsyncThunk(
  'products/adjustStock',
  async (
    { id, data }: { id: number; data: AdjustStockReq },
    { rejectWithValue }
  ) => {
    try {
      const response = await productApi.adjustStock(id, data);
      return response;
    } catch (err: any) {
      if (err.response?.status === 409) {
        return rejectWithValue(err.response.data);
      }
      return rejectWithValue({
        message: err.response?.data?.message || 'Gagal penyesuaian stok',
      });
    }
  }
);

export const fetchProductMovements = createAsyncThunk(
  'products/fetchMovements',
  async (id: number, { rejectWithValue }) => {
    try {
      const response = await productApi.getProductMovements(id);
      return response.data;
    } catch (err: any) {
      return rejectWithValue(
        err.response?.data?.message || 'Gagal mengambil riwayat movement'
      );
    }
  }
);

const productSlice = createSlice({
  name: 'products',
  initialState,
  reducers: {
    clearProductError: (state) => {
      state.error = null;
      state.errorDetails = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchProducts.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(
        fetchProducts.fulfilled,
        (state, action: PayloadAction<PaginatedProductsResponse>) => {
          state.isLoading = false;
          state.products = action.payload.data;
          state.pagination = {
            page: action.payload.page,
            limit: action.payload.limit,
            total_data: action.payload.total_data,
            total_page: action.payload.total_page,
            next_page: action.payload.next_page,
            prev_page: action.payload.prev_page,
          };
        }
      )
      .addCase(fetchProducts.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      .addCase(
        fetchProductMovements.fulfilled,
        (state, action: PayloadAction<StockMovement[]>) => {
          state.movements = action.payload;
        }
      );
  },
});

export const { clearProductError } = productSlice.actions;
export default productSlice.reducer;
