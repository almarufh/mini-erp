import { BrowserRouter, Route, Routes } from 'react-router';
import Layout from './components/Layout';
import OrderListPage from './pages/OrderListPage';
import ProductListPage from './pages/ProductListPage';

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
            element={<ProductListPage />}
          />
          <Route
            path="products"
            element={<ProductListPage />}
          />
          <Route
            path="orders"
            element={<OrderListPage />}
          />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
