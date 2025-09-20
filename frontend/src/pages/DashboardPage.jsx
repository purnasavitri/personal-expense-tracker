import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../api/axios';

function DashboardPage() {
  const [transactions, setTransactions] = useState([]);
  const [description, setDescription] = useState('');
  const [amount, setAmount] = useState('');
  const navigate = useNavigate();

  // Fungsi untuk mengambil data transaksi
  const fetchTransactions = async () => {
    const token = localStorage.getItem('token');
    try {
      const response = await api.get('/transactions', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setTransactions(response.data || []);
    } catch (error) {
      console.error('Gagal mengambil data transaksi:', error);
      // Jika token tidak valid, logout
      if (error.response && error.response.status === 401) {
        handleLogout();
      }
    }
  };

  useEffect(() => {
    fetchTransactions();
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };
  
  const handleAddTransaction = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    try {
      await api.post('/transactions', 
        { description, amount: parseFloat(amount), type: 'expense', category_id: 1 }, 
        { headers: { Authorization: `Bearer ${token}` } }
      );
      // Reset form dan muat ulang data
      setDescription('');
      setAmount('');
      fetchTransactions(); 
    } catch (error) {
      console.error('Gagal menambah transaksi:', error);
    }
  };

  return (
    <div>
      <h1>Dashboard</h1>
      <button onClick={handleLogout}>Logout</button>

      <h2>Tambah Transaksi Baru</h2>
      <form onSubmit={handleAddTransaction}>
        <input
          type="text"
          placeholder="Deskripsi"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
        />
        <input
          type="number"
          placeholder="Jumlah"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
          required
        />
        <button type="submit">Tambah</button>
      </form>
      
      <h2>Riwayat Transaksi</h2>
      <ul>
        {transactions.length > 0 ? (
          transactions.map((t) => (
            <li key={t.ID}>
              {t.Description}: Rp {t.Amount}
            </li>
          ))
        ) : (
          <p>Belum ada transaksi.</p>
        )}
      </ul>
    </div>
  );
}

export default DashboardPage;