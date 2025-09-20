import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';
import '../pages/Dashboard.css';
import '../components/Modal.css';
import FloatingActionButton from '../components/FloatingActionButton';
import TransactionModal from '../components/TransactionModal';

function DashboardPage() {
  const [transactions, setTransactions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState(null);
  const navigate = useNavigate();

  const fetchData = async () => {
    const token = sessionStorage.getItem('token');
    if (!token) { handleLogout(); return; }
    try {
      const headers = { Authorization: `Bearer ${token}` };
      const [transRes, catRes] = await Promise.all([
        api.get('/transactions', { headers }),
        api.get('/categories', { headers }),
      ]);
      setTransactions(transRes.data || []);
      setCategories(catRes.data || []);
    } catch (error) {
      console.error('Gagal mengambil data:', error);
      if (error.response && error.response.status === 401) handleLogout();
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleOpenAddModal = () => {
    setEditingTransaction(null);
    setIsModalOpen(true);
  };
  
  const handleOpenEditModal = (transaction) => {
    setEditingTransaction(transaction);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };
  
  const handleSaveTransaction = async (data, transactionId) => {
    const token = sessionStorage.getItem('token');
    const headers = { Authorization: `Bearer ${token}` };
    try {
      if (transactionId) {
        await api.put(`/transactions/${transactionId}`, data, { headers });
      } else {
        await api.post('/transactions', data, { headers });
      }
      handleCloseModal();
      fetchData();
    } catch (error) {
      console.error('Gagal menyimpan transaksi:', error);
    }
  };

  const handleLogout = () => {
    sessionStorage.removeItem('token');
    navigate('/login');
  };

  const pendapatan = transactions.filter(t => t.Type === 'income').reduce((acc, t) => acc + t.Amount, 0);
  const pengeluaran = transactions.filter(t => t.Type === 'expense').reduce((acc, t) => acc + t.Amount, 0);
  const total = pendapatan - pengeluaran;
  const formatRupiah = (number) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(number);

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Dashboard Keuangan</h1>
        <button onClick={handleLogout}>Logout</button>
      </div>
      
      <div className="summary-cards">
        <div className="card"><h3>Pendapatan</h3><p className="amount income">{formatRupiah(pendapatan)}</p></div>
        <div className="card"><h3>Pengeluaran</h3><p className="amount expense">{formatRupiah(pengeluaran)}</p></div>
        <div className="card"><h3>Total</h3><p className="amount total">{formatRupiah(total)}</p></div>
      </div>
      
      <div className="transaction-section">
        <h2>TRANSAKSI</h2>
        <table className="transaction-table">
          <thead>
            <tr>
              <th>No</th><th>Deskripsi</th><th>Jumlah</th><th>Tipe</th><th>Tanggal</th><th>Kategori</th>
            </tr>
          </thead>
          <tbody>
            {transactions.map((t, index) => (
              <tr key={t.ID} onClick={() => handleOpenEditModal(t)}>
                <td>{index + 1}</td>
                <td>{t.Description}</td>
                <td>{formatRupiah(t.Amount)}</td>
                <td>{t.Type}</td>
                <td>{new Date(t.CreatedAt).toLocaleDateString('id-ID')}</td>
                <td>{t.category_name}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <FloatingActionButton onClick={handleOpenAddModal} />
      <TransactionModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSave={handleSaveTransaction}
        transactionToEdit={editingTransaction}
        categories={categories}
      />
    </div>
  );
}

export default DashboardPage;