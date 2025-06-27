import React from 'react';
import '../styles/Components.css';

interface LoadingSpinnerProps {
  message?: string;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ 
  message = "Yükleniyor..." 
}) => {
  return (
    <div className="loading-spinner">
      <div className="spinner-icon">⏳</div>
      <div className="spinner-text">{message}</div>
    </div>
  );
};

export default LoadingSpinner; 