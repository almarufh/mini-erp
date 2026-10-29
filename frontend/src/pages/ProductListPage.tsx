import { useEffect, useState } from 'react';
import {
  FiActivity,
  FiCheckCircle,
  FiChevronLeft,
  FiChevronRight,
  FiPlus,
  FiRotateCw,
  FiSearch,
  FiTrash2,
  FiXCircle,
} from 'react-icons/fi';
import { productApi } from '../api/productApi';
import Modal from '../components/ui/Modal';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import {
  adjustProductStock,
  createProduct,
  fetchProductMovements,
  fetchProducts,
} from '../store/slices/productSlice';
import type { AdjustStockReq, CreateProductReq, Product } from '../types';

export default function ProductListPage() {
  const dispatch = useAppDispatch();
  const { products, pagination, movements, isLoading, error } = useAppSelector(
    (state) => state.products
  );

  // Search & Pagination State
  const [searchTerm, setSearchTerm] = useState('');
  const [page, setPage] = useState(1);

  // Modal States
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isStockModalOpen, setIsStockModalOpen] = useState(false);
  const [isMovementModalOpen, setIsMovementModalOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

  // Form Input States
  const [addForm, setAddForm] = useState<CreateProductReq>({
    sku: '',
    name: '',
    price: 0,
    stock: 0,
    is_active: true,
  });

  const [stockForm, setStockForm] = useState<AdjustStockReq>({
    type: 'in',
    qty: 1,
    note: '',
  });

  const [actionError, setActionError] = useState<string | null>(null);

  useEffect(() => {
    dispatch(fetchProducts({ search: searchTerm, page, limit: 10 }));
  }, [dispatch, searchTerm, page]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
    setPage(1);
  };

  // Submit Add Product
  const handleAddSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setActionError(null);
    try {
      await dispatch(createProduct(addForm)).unwrap();
      setIsAddModalOpen(false);
      setAddForm({ sku: '', name: '', price: 0, stock: 0, is_active: true });
      dispatch(fetchProducts({ search: searchTerm, page, limit: 10 }));
    } catch (err: any) {
      setActionError(err || 'Gagal menambahkan produk');
    }
  };

  // Submit Adjust Stock
  const handleStockSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProduct) return;
    setActionError(null);

    const res = await dispatch(
      adjustProductStock({ id: selectedProduct.id, data: stockForm })
    );

    if (adjustProductStock.fulfilled.match(res)) {
      setIsStockModalOpen(false);
      setStockForm({ type: 'in', qty: 1, note: '' });
      dispatch(fetchProducts({ search: searchTerm, page, limit: 10 }));
    } else if (adjustProductStock.rejected.match(res)) {
      const errData: any = res.payload;
      if (errData?.error?.code === 'INSUFFICIENT_STOCK') {
        setActionError(
          `Stok tidak mencukupi! Stok yang tersedia hanya ${errData.error.available}`
        );
      } else {
        setActionError(errData?.message || 'Gagal penyesuaian stok');
      }
    }
  };

  // Open Movements Modal
  const handleOpenMovements = (product: Product) => {
    setSelectedProduct(product);
    dispatch(fetchProductMovements(product.id));
    setIsMovementModalOpen(true);
  };

  // Handle Delete Product
  const handleDeleteProduct = async (id: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus produk ini?')) return;
    try {
      await productApi.deleteProduct(id);
      dispatch(fetchProducts({ search: searchTerm, page, limit: 10 }));
    } catch (err: any) {
      if (err.response?.status === 409) {
        alert(
          'Gagal Menghapus: Produk ini sudah memiliki riwayat movement stok!'
        );
      } else {
        alert(err.response?.data?.message || 'Gagal menghapus produk');
      }
    }
  };

  return (
    <div className="space-y-6">
      {/* Header Actions */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold text-dark">Katalog Produk</h1>
          <p className="text-sm text-dark-gray">
            Kelola data produk, persediaan stok, dan audit riwayat pergerakan.
          </p>
        </div>

        <button
          onClick={() => {
            setActionError(null);
            setIsAddModalOpen(true);
          }}
          className="flex items-center px-4 py-2.5 bg-primary hover:bg-primary-hover text-secondary font-semibold rounded-lg shadow-xs transition-colors cursor-pointer"
        >
          <FiPlus className="mr-2 text-lg" />
          Tambah Produk
        </button>
      </div>

      {/* Filter & Search Bar */}
      <div className="flex items-center justify-between bg-secondary p-4 rounded-xl border border-border shadow-2xs">
        <div className="relative w-full max-w-md">
          <FiSearch className="absolute left-3.5 top-1/2 -translate-y-1/2 text-dark-gray text-lg" />
          <input
            type="text"
            value={searchTerm}
            onChange={handleSearchChange}
            placeholder="Cari SKU atau nama produk..."
            className="w-full pl-10 pr-4 py-2 border border-border rounded-lg text-sm bg-gray-light focus:outline-none focus:border-primary focus:bg-secondary transition-all"
          />
        </div>
      </div>

      {/* Product Table */}
      <div className="bg-secondary rounded-xl border border-border shadow-2xs overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm border-collapse">
            <thead className="bg-gray-light border-b border-border text-dark-gray font-semibold">
              <tr>
                <th className="py-3.5 px-4">SKU</th>
                <th className="py-3.5 px-4">Nama Produk</th>
                <th className="py-3.5 px-4">Harga</th>
                <th className="py-3.5 px-4">Stok</th>
                <th className="py-3.5 px-4">Status</th>
                <th className="py-3.5 px-4 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {isLoading ? (
                <tr>
                  <td
                    colSpan={6}
                    className="py-8 text-center text-dark-gray"
                  >
                    Memuat data produk...
                  </td>
                </tr>
              ) : products.length === 0 ? (
                <tr>
                  <td
                    colSpan={6}
                    className="py-8 text-center text-dark-gray"
                  >
                    Tidak ada produk ditemukan.
                  </td>
                </tr>
              ) : (
                products.map((p) => (
                  <tr
                    key={p.id}
                    className="hover:bg-slate-50 transition-colors"
                  >
                    <td className="py-3.5 px-4 font-mono font-medium text-dark">
                      {p.sku}
                    </td>
                    <td className="py-3.5 px-4 font-semibold text-dark">
                      {p.name}
                    </td>
                    <td className="py-3.5 px-4">
                      Rp {p.price.toLocaleString('id-ID')}
                    </td>
                    <td className="py-3.5 px-4">
                      <span
                        className={`font-bold px-2.5 py-1 rounded-full text-xs ${
                          p.stock <= 5
                            ? 'bg-red-50 text-red font-extrabold'
                            : 'bg-emerald-50 text-emerald-700'
                        }`}
                      >
                        {p.stock} pcs
                      </span>
                    </td>
                    <td className="py-3.5 px-4">
                      {p.is_active ? (
                        <span className="inline-flex items-center text-xs font-semibold text-emerald-600 bg-emerald-50 px-2.5 py-1 rounded-full">
                          <FiCheckCircle className="mr-1" /> Aktif
                        </span>
                      ) : (
                        <span className="inline-flex items-center text-xs font-semibold text-slate-500 bg-slate-100 px-2.5 py-1 rounded-full">
                          <FiXCircle className="mr-1" /> Nonaktif
                        </span>
                      )}
                    </td>
                    <td className="py-3.5 px-4">
                      <div className="flex items-center justify-center space-x-2">
                        <button
                          onClick={() => {
                            setSelectedProduct(p);
                            setActionError(null);
                            setIsStockModalOpen(true);
                          }}
                          title="Adjust Stock"
                          className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors cursor-pointer"
                        >
                          <FiRotateCw className="text-base" />
                        </button>

                        <button
                          onClick={() => handleOpenMovements(p)}
                          title="Movement History"
                          className="p-2 text-amber-600 hover:bg-amber-50 rounded-lg transition-colors cursor-pointer"
                        >
                          <FiActivity className="text-base" />
                        </button>

                        <button
                          onClick={() => handleDeleteProduct(p.id)}
                          title="Hapus Produk"
                          className="p-2 text-red hover:bg-red-50 rounded-lg transition-colors cursor-pointer"
                        >
                          <FiTrash2 className="text-base" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination Footer */}
        {pagination && pagination.total_page > 1 && (
          <div className="flex items-center justify-between px-6 py-4 border-t border-border bg-gray-light">
            <span className="text-xs text-dark-gray">
              Halaman <span className="font-semibold">{pagination.page}</span>{' '}
              dari{' '}
              <span className="font-semibold">{pagination.total_page}</span> (
              {pagination.total_data} Produk)
            </span>

            <div className="flex space-x-2">
              <button
                disabled={!pagination.prev_page}
                onClick={() => setPage((prev) => prev - 1)}
                className="p-2 rounded-lg border border-border bg-secondary text-dark disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-100 transition-colors"
              >
                <FiChevronLeft />
              </button>
              <button
                disabled={!pagination.next_page}
                onClick={() => setPage((prev) => prev + 1)}
                className="p-2 rounded-lg border border-border bg-secondary text-dark disabled:opacity-40 disabled:cursor-not-allowed hover:bg-slate-100 transition-colors"
              >
                <FiChevronRight />
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Modal 1: Tambah Produk */}
      <Modal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        title="Tambah Produk Baru"
      >
        <form
          onSubmit={handleAddSubmit}
          className="space-y-4"
        >
          {actionError && (
            <div className="p-3 bg-red-50 border border-red-200 text-red text-xs rounded-lg">
              {actionError}
            </div>
          )}

          <div>
            <label className="block text-xs font-bold text-dark mb-1">
              SKU Produk
            </label>
            <input
              type="text"
              required
              value={addForm.sku}
              onChange={(e) => setAddForm({ ...addForm, sku: e.target.value })}
              placeholder="Contoh: SKU-KB-001"
              className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-bold text-dark mb-1">
              Nama Produk
            </label>
            <input
              type="text"
              required
              value={addForm.name}
              onChange={(e) => setAddForm({ ...addForm, name: e.target.value })}
              placeholder="Nama Keyboard / Barang"
              className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-bold text-dark mb-1">
                Harga (Rp)
              </label>
              <input
                type="number"
                min="0"
                required
                value={addForm.price}
                onChange={(e) =>
                  setAddForm({ ...addForm, price: Number(e.target.value) })
                }
                className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
              />
            </div>

            <div>
              <label className="block text-xs font-bold text-dark mb-1">
                Stok Awal
              </label>
              <input
                type="number"
                min="0"
                required
                value={addForm.stock}
                onChange={(e) =>
                  setAddForm({ ...addForm, stock: Number(e.target.value) })
                }
                className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
              />
            </div>
          </div>

          <div className="pt-4 flex justify-end space-x-2">
            <button
              type="button"
              onClick={() => setIsAddModalOpen(false)}
              className="px-4 py-2 border border-border rounded-lg text-xs font-semibold text-dark-gray hover:bg-slate-100 transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-primary hover:bg-primary-hover text-secondary text-xs font-bold rounded-lg transition-colors"
            >
              Simpan Produk
            </button>
          </div>
        </form>
      </Modal>

      {/* Modal 2: Adjust Stock */}
      <Modal
        isOpen={isStockModalOpen}
        onClose={() => setIsStockModalOpen(false)}
        title={`Penyesuaian Stok - ${selectedProduct?.name}`}
      >
        <form
          onSubmit={handleStockSubmit}
          className="space-y-4"
        >
          {actionError && (
            <div className="p-3 bg-red-50 border border-red-200 text-red text-xs font-semibold rounded-lg">
              {actionError}
            </div>
          )}

          <div>
            <label className="block text-xs font-bold text-dark mb-2">
              Jenis Adjustment
            </label>
            <div className="grid grid-cols-2 gap-3">
              <label
                className={`flex items-center justify-center p-2.5 rounded-lg border text-xs font-bold cursor-pointer transition-all ${
                  stockForm.type === 'in'
                    ? 'border-emerald-500 bg-emerald-50 text-emerald-700'
                    : 'border-border text-dark-gray hover:bg-slate-50'
                }`}
              >
                <input
                  type="radio"
                  name="type"
                  value="in"
                  checked={stockForm.type === 'in'}
                  onChange={() => setStockForm({ ...stockForm, type: 'in' })}
                  className="hidden"
                />
                Stock In (Tambah)
              </label>

              <label
                className={`flex items-center justify-center p-2.5 rounded-lg border text-xs font-bold cursor-pointer transition-all ${
                  stockForm.type === 'out'
                    ? 'border-red-500 bg-red-50 text-red'
                    : 'border-border text-dark-gray hover:bg-slate-50'
                }`}
              >
                <input
                  type="radio"
                  name="type"
                  value="out"
                  checked={stockForm.type === 'out'}
                  onChange={() => setStockForm({ ...stockForm, type: 'out' })}
                  className="hidden"
                />
                Stock Out (Kurang)
              </label>
            </div>
          </div>

          <div>
            <label className="block text-xs font-bold text-dark mb-1">
              Jumlah (Qty)
            </label>
            <input
              type="number"
              min="1"
              required
              value={stockForm.qty}
              onChange={(e) =>
                setStockForm({ ...stockForm, qty: Number(e.target.value) })
              }
              className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
            />
          </div>

          <div>
            <label className="block text-xs font-bold text-dark mb-1">
              Catatan / Note (Opsional)
            </label>
            <input
              type="text"
              value={stockForm.note || ''}
              onChange={(e) =>
                setStockForm({ ...stockForm, note: e.target.value })
              }
              placeholder="Restock supplier / barang rusak"
              className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
            />
          </div>

          <div className="pt-4 flex justify-end space-x-2">
            <button
              type="button"
              onClick={() => setIsStockModalOpen(false)}
              className="px-4 py-2 border border-border rounded-lg text-xs font-semibold text-dark-gray hover:bg-slate-100 transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-primary hover:bg-primary-hover text-secondary text-xs font-bold rounded-lg transition-colors"
            >
              Simpan Adjustment
            </button>
          </div>
        </form>
      </Modal>

      {/* Modal 3: Movement History Drawer */}
      <Modal
        isOpen={isMovementModalOpen}
        onClose={() => setIsMovementModalOpen(false)}
        title={`Riwayat Movement Stok - ${selectedProduct?.sku}`}
      >
        <div className="space-y-3 max-h-96 overflow-y-auto pr-1">
          {movements.length === 0 ? (
            <p className="text-xs text-dark-gray text-center py-6">
              Belum ada riwayat pergerakan stok.
            </p>
          ) : (
            movements.map((m) => (
              <div
                key={m.id}
                className="flex items-center justify-between p-3 rounded-lg border border-border bg-gray-light"
              >
                <div>
                  <span
                    className={`inline-block text-[10px] font-extrabold px-2 py-0.5 rounded-full uppercase mr-2 ${
                      m.type === 'in'
                        ? 'bg-emerald-100 text-emerald-800'
                        : 'bg-red-100 text-red'
                    }`}
                  >
                    {m.type}
                  </span>
                  <span className="text-xs font-bold text-dark">
                    {m.qty} pcs
                  </span>
                  <p className="text-[11px] text-dark-gray mt-0.5">
                    {m.note || 'Tidak ada catatan'}
                  </p>
                </div>

                <span className="text-[10px] text-dark-gray">
                  {new Date(m.created_at).toLocaleString('id-ID')}
                </span>
              </div>
            ))
          )}
        </div>
      </Modal>
    </div>
  );
}
