import { BrowserRouter, Route, Routes } from 'react-router';
import Layout from './components/Layout';

// Placeholder sementara untuk halaman
const ProductsPage = () => (
  <div>
    <h1 className="text-2xl font-bold mb-4">Product Catalog</h1>
    <p className="text-slate-600">Daftar produk akan ditampilkan di sini.</p>
  </div>
);

const OrdersPage = () => (
  <div>
    <h1 className="text-2xl font-bold mb-4">Sales Orders</h1>
    <p className="text-slate-600">Daftar order akan ditampilkan di sini.</p>
  </div>
);

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={<Layout />}
        >
          <Route
            index
            element={<ProductsPage />}
          />
          <Route
            path="products"
            element={<ProductsPage />}
          />
          <Route
            path="orders"
            element={<OrdersPage />}
          />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
