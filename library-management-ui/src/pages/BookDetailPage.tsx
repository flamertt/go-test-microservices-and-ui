import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import type { Book } from '../types';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import '../styles/BookDetailPage.css';

const BookDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const bookId = parseInt(id || '0');

  const { data: book, loading, error } = useApi<Book>(
    `/api/books/${bookId}`,
    [bookId]
  );

  if (loading) {
    return (
      <div className="book-loading-container">
        <LoadingSpinner message="Kitap detayları yükleniyor..." />
      </div>
    );
  }

  if (error || !book) {
    return (
      <div className="book-error-container">
        <button
          className="book-back-button"
          onClick={() => navigate('/books')}
        >
          ← Kitaplara Dön
        </button>
        <ErrorMessage 
          error={error || 'Kitap bulunamadı'} 
          title="Kitap Detayları Yüklenemiyor"
        />
      </div>
    );
  }

  return (
    <div className="book-detail-page">
      {/* Geri Dön Butonu */}
      <button
        className="book-back-button"
        onClick={() => navigate('/books')}
      >
        ← Kitaplara Dön
      </button>

      <div className="book-detail-card">
        {/* Kitap Başlığı ve Temel Bilgiler */}
        <div className="book-header">
          <div className="book-icon">📚</div>
          <div className="book-title-section">
            <h1>{book.title}</h1>
            <p className="book-subtitle">
              {book.author} tarafından yazılan bu eser {book.category_name} kategorisinde yer almaktadır.
            </p>
            
            <div className="book-meta-chips">
              <div className="book-chip author">
                👤 {book.author}
              </div>
              <div className="book-chip category">
                📖 {book.category_name}
              </div>
              <div className="book-chip year">
                📅 {book.released_year}
              </div>
            </div>
          </div>
        </div>

        {/* Kitap Açıklaması */}
        <div className="book-description">
          <h6>📖 Kitap Hakkında</h6>
          <p className="book-description-text">
            {`"${book.title}" ${book.author} tarafından ${book.released_year} yılında yazılmış ${book.page_count} sayfalık bir eserdir. ${book.publisher} yayınevi tarafından yayınlanmış olan bu kitap ${book.category_name} kategorisinde yer almaktadır ve ürün kodu ${book.product_code}'dir.`}
          </p>
        </div>

        {/* Detaylı Bilgiler */}
        <div className="book-details-panel">
          <h6>📋 Detaylı Bilgiler</h6>
          <div className="book-details-grid">
            <div className="book-detail-item">
              <span className="book-detail-label">📅 Yayın Yılı:</span>
              <span className="book-detail-value">{book.released_year}</span>
            </div>
            <div className="book-detail-item">
              <span className="book-detail-label">📄 Sayfa Sayısı:</span>
              <span className="book-detail-value">{book.page_count} sayfa</span>
            </div>
            <div className="book-detail-item">
              <span className="book-detail-label">🏢 Yayınevi:</span>
              <span className="book-detail-value">{book.publisher}</span>
            </div>
            <div className="book-detail-item">
              <span className="book-detail-label">🔖 Ürün Kodu:</span>
              <span className="book-detail-value">{book.product_code}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BookDetailPage; 