import React from 'react';
import '../styles/Components.css';

interface ErrorMessageProps {
  error: string;
  title?: string;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ 
  error, 
  title = "Bir hata oluştu" 
}) => {
  return (
    <div className="error-message">
      <div className="error-icon">❌</div>
      <div className="error-text">{title}: {error}</div>
    </div>
  );
};

export default ErrorMessage; 