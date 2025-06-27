import React from 'react';
import '../styles/Components.css';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  totalItems?: number;
  itemsPerPage?: number;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
  totalItems,
  itemsPerPage,
}) => {
  if (totalPages <= 1) return null;

  const renderPageNumbers = () => {
    const pages = [];
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);

    // İlk sayfa
    if (startPage > 1) {
      pages.push(
        <button
          key={1}
          className="pagination-btn"
          onClick={() => onPageChange(1)}
        >
          1
        </button>
      );
      
      if (startPage > 2) {
        pages.push(<span key="start-ellipsis">...</span>);
      }
    }

    // Orta sayfalar
    for (let i = startPage; i <= endPage; i++) {
      pages.push(
        <button
          key={i}
          className={`pagination-btn ${i === currentPage ? 'active' : ''}`}
          onClick={() => onPageChange(i)}
        >
          {i}
        </button>
      );
    }

    // Son sayfa
    if (endPage < totalPages) {
      if (endPage < totalPages - 1) {
        pages.push(<span key="end-ellipsis">...</span>);
      }
      
      pages.push(
        <button
          key={totalPages}
          className="pagination-btn"
          onClick={() => onPageChange(totalPages)}
        >
          {totalPages}
        </button>
      );
    }

    return pages;
  };

  return (
    <div className="pagination-container">
      <div className="pagination-controls">
        <button
          className="pagination-btn"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage <= 1}
        >
          ← Önceki
        </button>
        
        {renderPageNumbers()}
        
        <button
          className="pagination-btn"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage >= totalPages}
        >
          Sonraki →
        </button>
      </div>
      
      {(totalItems !== undefined && itemsPerPage !== undefined) && (
        <div className="pagination-info">
          <span>
            Sayfa {currentPage} / {totalPages} 
            {totalItems > 0 && ` (Toplam ${totalItems} öğe)`}
          </span>
        </div>
      )}
    </div>
  );
};

export default Pagination; 