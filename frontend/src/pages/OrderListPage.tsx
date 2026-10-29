import { useEffect, useState } from 'react';
import {
  FiCheckCircle,
  FiEye,
  FiPlus,
  FiTrash2,
  FiXCircle,
} from 'react-icons/fi';
import { orderApi } from '../api/orderApi';
import Modal from '../components/ui/Modal';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import {
  cancelSalesOrder,
  clearOrderError,
  createSalesOrder,
  fetchOrders,
} from '../store/slices/orderSlice';
import { fetchProducts } from '../store/slices/productSlice';
import type { CreateOrderItemReq, CreateOrderReq, SalesOrder } from '../types';

export default function OrderListPage() {
  const dispatch = useAppDispatch();
  const { products } = useAppSelector((state) => state.products);
  const { orders, isLoading } = useAppSelector((state) => state.orders);

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<SalesOrder | null>(null);

  const [customer, setCustomer] = useState('');
  const [items, setItems] = useState<CreateOrderItemReq[]>([
    { sku: '', qty: 1 },
  ]);
  const [formError, setFormError] = useState<string | null>(null);

  useEffect(() => {
    dispatch(fetchOrders());
    dispatch(fetchProducts({ limit: 100 }));
  }, [dispatch]);

  const handleAddItemRow = () => {
    setItems([...items, { sku: '', qty: 1 }]);
  };

  const handleRemoveItemRow = (index: number) => {
    if (items.length === 1) return;
    setItems(items.filter((_, i) => i !== index));
  };

  const handleItemChange = (
    index: number,
    field: keyof CreateOrderItemReq,
    value: any
  ) => {
    const updated = [...items];
    updated[index] = { ...updated[index], [field]: value };
    setItems(updated);
  };

  const calculateEstimatedTotal = () => {
    return items.reduce((acc, item) => {
      const prod = products.find((p) => p.sku === item.sku);
      if (prod) {
        return acc + prod.price * (item.qty || 0);
      }
      return acc;
    }, 0);
  };

  const handleCreateOrderSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    if (!customer.trim()) {
      setFormError('Nama customer wajib diisi');
      return;
    }

    const hasEmptySKU = items.some((i) => !i.sku);
    if (hasEmptySKU) {
      setFormError('Harap pilih SKU untuk semua item belanjaan');
      return;
    }

    const payload: CreateOrderReq = { customer, items };

    const res = await dispatch(createSalesOrder(payload));

    if (createSalesOrder.fulfilled.match(res)) {
      setIsCreateModalOpen(false);
      setCustomer('');
      setItems([{ sku: '', qty: 1 }]);
      dispatch(fetchProducts({ limit: 100 }));
    } else if (createSalesOrder.rejected.match(res)) {
      const errData: any = res.payload;
      if (errData?.error?.code === 'INSUFFICIENT_STOCK') {
        setFormError(
          `Stok untuk SKU "${errData.error.sku}" tidak mencukupi! Tersisa: ${errData.error.available}`
        );
      } else {
        setFormError(errData?.message || 'Gagal membuat sales order');
      }
    }
  };

  const handleViewDetail = async (id: number) => {
    try {
      const res = await orderApi.getOrderById(id);
      if (res.results) {
        setSelectedOrder(res.results);
        setIsDetailModalOpen(true);
      }
    } catch {
      alert('Gagal mengambil detail order');
    }
  };

  const handleCancelOrder = async (id: number) => {
    if (
      !confirm(
        'Apakah Anda yakin ingin membatalkan order ini? Stok akan dikembalikan ke gudang.'
      )
    )
      return;

    const res = await dispatch(cancelSalesOrder(id));
    if (cancelSalesOrder.fulfilled.match(res)) {
      dispatch(fetchProducts({ limit: 100 }));
      alert('Order berhasil dibatalkan dan stok telah dikembalikan.');
    } else {
      alert('Gagal membatalkan order');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold text-dark">Sales Orders</h1>
          <p className="text-sm text-dark-gray">
            Pencatatan transaksi penjualan, pemotongan stok otomatis, dan
            pembatalan order.
          </p>
        </div>

        <button
          onClick={() => {
            setFormError(null);
            dispatch(clearOrderError());
            setIsCreateModalOpen(true);
          }}
          className="flex items-center px-4 py-2.5 bg-primary hover:bg-primary-hover text-secondary font-semibold rounded-lg shadow-xs transition-colors cursor-pointer"
        >
          <FiPlus className="mr-2 text-lg" />
          Buat Order Baru
        </button>
      </div>

      <div className="bg-secondary rounded-xl border border-border shadow-2xs overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm border-collapse">
            <thead className="bg-gray-light border-b border-border text-dark-gray font-semibold">
              <tr>
                <th className="py-3.5 px-4">Order ID</th>
                <th className="py-3.5 px-4">Customer</th>
                <th className="py-3.5 px-4">Total Price</th>
                <th className="py-3.5 px-4">Status</th>
                <th className="py-3.5 px-4">Tanggal Order</th>
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
                    Memuat data sales order...
                  </td>
                </tr>
              ) : orders.length === 0 ? (
                <tr>
                  <td
                    colSpan={6}
                    className="py-8 text-center text-dark-gray"
                  >
                    Belum ada sales order yang dibuat. Klik tombol di atas untuk
                    membuat order baru.
                  </td>
                </tr>
              ) : (
                orders.map((order) => (
                  <tr
                    key={order.id}
                    className="hover:bg-slate-50 transition-colors"
                  >
                    <td className="py-3.5 px-4 font-mono font-bold text-dark">
                      #{order.id}
                    </td>
                    <td className="py-3.5 px-4 font-semibold text-dark">
                      {order.customer}
                    </td>
                    <td className="py-3.5 px-4 font-bold text-indigo-600">
                      Rp {order.total.toLocaleString('id-ID')}
                    </td>
                    <td className="py-3.5 px-4">
                      {order.status === 'CREATED' ? (
                        <span className="inline-flex items-center text-xs font-semibold text-emerald-600 bg-emerald-50 px-2.5 py-1 rounded-full">
                          <FiCheckCircle className="mr-1" /> CREATED
                        </span>
                      ) : (
                        <span className="inline-flex items-center text-xs font-semibold text-red bg-red-50 px-2.5 py-1 rounded-full">
                          <FiXCircle className="mr-1" /> CANCELLED
                        </span>
                      )}
                    </td>
                    <td className="py-3.5 px-4 text-dark-gray text-xs">
                      {new Date(order.created_at).toLocaleString('id-ID')}
                    </td>
                    <td className="py-3.5 px-4">
                      <div className="flex items-center justify-center space-x-2">
                        <button
                          onClick={() => handleViewDetail(order.id)}
                          title="Lihat Detail Item"
                          className="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors cursor-pointer"
                        >
                          <FiEye className="text-base" />
                        </button>

                        {order.status === 'CREATED' && (
                          <button
                            onClick={() => handleCancelOrder(order.id)}
                            title="Batalkan Order"
                            className="p-2 text-red hover:bg-red-50 rounded-lg transition-colors cursor-pointer"
                          >
                            <FiXCircle className="text-base" />
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <Modal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        title="Buat Sales Order Baru"
      >
        <form
          onSubmit={handleCreateOrderSubmit}
          className="space-y-4"
        >
          {formError && (
            <div className="p-3 bg-red-50 border border-red-200 text-red text-xs font-semibold rounded-lg">
              {formError}
            </div>
          )}

          <div>
            <label className="block text-xs font-bold text-dark mb-1">
              Nama Customer
            </label>
            <input
              type="text"
              required
              value={customer}
              onChange={(e) => setCustomer(e.target.value)}
              placeholder="Contoh: Toko Maju Jaya"
              className="w-full px-3 py-2 border border-border rounded-lg text-sm bg-gray-light focus:bg-secondary focus:border-primary outline-none"
            />
          </div>

          <div>
            <div className="flex justify-between items-center mb-2">
              <label className="text-xs font-bold text-dark">
                Daftar Item Belanjaan
              </label>
              <button
                type="button"
                onClick={handleAddItemRow}
                className="text-xs text-primary font-bold hover:underline flex items-center cursor-pointer"
              >
                <FiPlus className="mr-1" /> Tambah Item
              </button>
            </div>

            <div className="space-y-2 max-h-56 overflow-y-auto pr-1">
              {items.map((item, idx) => (
                <div
                  key={idx}
                  className="flex items-center space-x-2 bg-gray-light p-2.5 rounded-lg border border-border"
                >
                  <select
                    required
                    value={item.sku}
                    onChange={(e) =>
                      handleItemChange(idx, 'sku', e.target.value)
                    }
                    className="flex-1 px-2.5 py-1.5 border border-border rounded-lg text-xs bg-secondary focus:border-primary outline-none"
                  >
                    <option value="">-- Pilih Produk (SKU) --</option>
                    {products
                      .filter((p) => p.is_active)
                      .map((p) => (
                        <option
                          key={p.id}
                          value={p.sku}
                        >
                          {p.sku} - {p.name} (Stok: {p.stock})
                        </option>
                      ))}
                  </select>

                  <input
                    type="number"
                    min="1"
                    required
                    value={item.qty}
                    onChange={(e) =>
                      handleItemChange(idx, 'qty', Number(e.target.value))
                    }
                    className="w-20 px-2 py-1.5 border border-border rounded-lg text-xs text-center bg-secondary outline-none"
                  />

                  {items.length > 1 && (
                    <button
                      type="button"
                      onClick={() => handleRemoveItemRow(idx)}
                      className="text-red hover:bg-red-50 p-1.5 rounded-lg transition-colors cursor-pointer"
                    >
                      <FiTrash2 />
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>

          <div className="p-3 bg-indigo-50 border border-indigo-100 rounded-lg flex justify-between items-center">
            <span className="text-xs font-bold text-indigo-900">
              Estimasi Total Bayar:
            </span>
            <span className="text-base font-extrabold text-primary">
              Rp {calculateEstimatedTotal().toLocaleString('id-ID')}
            </span>
          </div>

          <div className="pt-4 flex justify-end space-x-2">
            <button
              type="button"
              onClick={() => setIsCreateModalOpen(false)}
              className="px-4 py-2 border border-border rounded-lg text-xs font-semibold text-dark-gray hover:bg-slate-100 transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-4 py-2 bg-primary hover:bg-primary-hover text-secondary text-xs font-bold rounded-lg transition-colors cursor-pointer disabled:opacity-50"
            >
              {isLoading ? 'Memproses...' : 'Proses Order'}
            </button>
          </div>
        </form>
      </Modal>

      <Modal
        isOpen={isDetailModalOpen}
        onClose={() => setIsDetailModalOpen(false)}
        title={`Detail Order #${selectedOrder?.id}`}
      >
        {selectedOrder && (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-2 text-xs border-b border-border pb-3">
              <div>
                <span className="text-dark-gray block">Customer:</span>
                <span className="font-bold text-dark">
                  {selectedOrder.customer}
                </span>
              </div>
              <div>
                <span className="text-dark-gray block">Status:</span>
                <span
                  className={`font-bold ${
                    selectedOrder.status === 'CREATED'
                      ? 'text-emerald-600'
                      : 'text-red'
                  }`}
                >
                  {selectedOrder.status}
                </span>
              </div>
            </div>

            <div>
              <h4 className="text-xs font-bold text-dark mb-2">
                Item Belanjaan:
              </h4>
              <div className="space-y-2 max-h-48 overflow-y-auto">
                {selectedOrder.items?.map((item) => (
                  <div
                    key={item.id}
                    className="flex justify-between items-center p-2.5 bg-gray-light border border-border rounded-lg text-xs"
                  >
                    <div>
                      <span className="font-mono font-bold block">
                        {item.sku}
                      </span>
                      <span className="text-dark-gray">
                        {item.qty} pcs x Rp{' '}
                        {item.unit_price.toLocaleString('id-ID')}
                      </span>
                    </div>
                    <span className="font-bold text-dark">
                      Rp {(item.qty * item.unit_price).toLocaleString('id-ID')}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            <div className="pt-2 border-t border-border flex justify-between items-center">
              <span className="text-xs font-bold text-dark">
                Total Transaksi:
              </span>
              <span className="text-lg font-extrabold text-primary">
                Rp {selectedOrder.total.toLocaleString('id-ID')}
              </span>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}
