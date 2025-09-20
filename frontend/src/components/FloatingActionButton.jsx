import React from 'react';

function FloatingActionButton({ onClick }) {
  return (
    <button className="fab" onClick={onClick}>
      +
    </button>
  );
}

export default FloatingActionButton;