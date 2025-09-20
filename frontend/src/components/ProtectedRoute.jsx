import React from 'react';
import { Navigate } from 'react-router-dom';

function ProtectedRoute({ children }) {
  const token = localStorage.getItem('token');

  if (!token) {
    // Jika tidak ada token, arahkan ke halaman login
    return <Navigate to="/login" />;
  }

  return children; // Jika ada token, tampilkan halaman yang diminta
}

export default ProtectedRoute;