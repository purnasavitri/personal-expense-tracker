import React, { useState, useEffect } from 'react';
import api from '../api/axios';

function TransactionModal({ isOpen, onClose, onSave, onDelete, transactionToEdit, categories }) {
  const [formData, setFormData] = useState({
    description: '',
    amount: '',
    type: 'expense',
    category_id: '',
  });
  const [newCategoryName, setNewCategoryName] = useState('');
  const [isAddingNewCategory, setIsAddingNewCategory] = useState(false);

  useEffect(() => {
    if (isOpen) {
      if (transactionToEdit) {
        setFormData({
          description: transactionToEdit.Description,
          amount: transactionToEdit.Amount,
          type: transactionToEdit.Type,
          category_id: transactionToEdit.CategoryID,
        });
      } else {
        setFormData({
          description: '',
          amount: '',
          type: 'expense',
          category_id: categories.length > 0 ? categories[0].ID : '',
        });
      }
      setIsAddingNewCategory(false);
      setNewCategoryName('');
    }
  }, [isOpen, transactionToEdit, categories]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name === 'category_id' && value === '--add-new--') {
      setIsAddingNewCategory(true);
    } else {
      setIsAddingNewCategory(false);
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    let finalCategoryId = formData.category_id;

    if (isAddingNewCategory && newCategoryName) {
      const token = sessionStorage.getItem('token');
      try {
        const response = await api.post('/categories', 
          { name: newCategoryName },
          { headers: { Authorization: `Bearer ${token}` } }
        );
        finalCategoryId = response.data.ID;
      } catch (error) {
        console.error('Gagal membuat kategori baru:', error);
        return;
      }
    }
    
    const finalData = { ...formData, amount: parseFloat(formData.amount), category_id: parseInt(finalCategoryId) };
    onSave(finalData, transactionToEdit ? transactionToEdit.ID : null);
  };

  const handleDeleteClick = () => {
    onDelete(transactionToEdit.ID);
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close-button" onClick={onClose}>&times;</button>
        <form onSubmit={handleSubmit} className="modal-form">
          <h2>{transactionToEdit ? 'Edit Transaksi' : 'Tambah Transaksi Baru'}</h2>
          <input name="description" type="text" placeholder="Deskripsi" value={formData.description} onChange={handleChange} required />
          <input name="amount" type="number" placeholder="Jumlah" value={formData.amount} onChange={handleChange} required />
          <select name="type" value={formData.type} onChange={handleChange}>
            <option value="expense">Pengeluaran</option>
            <option value="income">Pendapatan</option>
          </select>
          <select name="category_id" value={isAddingNewCategory ? '--add-new--' : formData.category_id} onChange={handleChange} required>
            {categories.map((cat) => (<option key={cat.ID} value={cat.ID}>{cat.Name}</option>))}
            <option value="--add-new--">-- Tambah Kategori Baru --</option>
          </select>
          {isAddingNewCategory && (
            <input
              type="text"
              placeholder="Nama Kategori Baru"
              className="new-category-input"
              value={newCategoryName}
              onChange={(e) => setNewCategoryName(e.target.value)}
              required
            />
          )}
          
          <div className="modal-actions">
            {transactionToEdit && (
              <button type="button" className="modal-delete-button" onClick={handleDeleteClick} title="Hapus Transaksi">
                Hapus
              </button>
            )}
            <button type="submit" className="save-button" style={{ marginLeft: 'auto' }}>Simpan</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default TransactionModal;