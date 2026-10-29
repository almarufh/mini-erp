import { FiBox, FiShoppingBag } from 'react-icons/fi';
import { Link, Outlet, useLocation } from 'react-router';

export default function Layout() {
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;

  return (
    <div className="min-h-screen bg-gray-light text-dark">
      {/* Header / Navbar */}
      <header className="sticky top-0 z-10 bg-secondary border-b border-border shadow-xs">
        <div className="app-container">
          <div className="flex justify-between h-16 items-center">
            {/* Brand Logo */}
            <div className="flex items-center space-x-3">
              <div className="bg-primary text-secondary px-2.5 py-1.5 rounded-lg font-bold text-sm tracking-wider">
                SF
              </div>
              <span className="text-xl font-bold text-dark tracking-tight">
                StockFlow
              </span>
            </div>

            {/* Nav Menu */}
            <nav className="flex space-x-2">
              <Link
                to="/"
                className={`flex items-center px-4 py-2 rounded-lg font-medium text-sm transition-all ${
                  isActive('/') || isActive('/products')
                    ? 'bg-indigo-50 text-primary font-semibold shadow-2xs'
                    : 'text-dark-gray hover:bg-slate-100 hover:text-dark'
                }`}
              >
                <FiBox className="mr-2 text-lg" />
                Products
              </Link>

              <Link
                to="/orders"
                className={`flex items-center px-4 py-2 rounded-lg font-medium text-sm transition-all ${
                  isActive('/orders')
                    ? 'bg-indigo-50 text-primary font-semibold shadow-2xs'
                    : 'text-dark-gray hover:bg-slate-100 hover:text-dark'
                }`}
              >
                <FiShoppingBag className="mr-2 text-lg" />
                Sales Orders
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content Area */}
      <main className="app-container py-8">
        <Outlet />
      </main>
    </div>
  );
}
