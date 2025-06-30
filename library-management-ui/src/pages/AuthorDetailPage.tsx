import React from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import type { AuthorDetail } from '../services/api';
import BookCard from '../components/BookCard';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import '../styles/AuthorDetailPage.css';

const AuthorDetailPage: React.FC = () => {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const authorName = decodeURIComponent(name || '');

  // Author detail API (yazar + kitapları birlikte getirir)
  const { data: authorDetail, loading, error } = useApi<AuthorDetail>(
    `/api/authors/detail/${encodeURIComponent(authorName)}`,
    [authorName]
  );

  if (loading) {
    return <LoadingSpinner message="Yazar detayları yükleniyor..." />;
  }

  if (error || !authorDetail) {
    return (
      <div className="author-error-container">
        <button
          className="author-back-button"
          onClick={() => navigate('/authors')}
        >
          ← Yazarlara Dön
        </button>
        <ErrorMessage 
          error={error || 'Yazar bulunamadı'} 
          title="Yazar Detayları Yüklenemiyor"
        />
      </div>
    );
  }

  const { author, books, book_count } = authorDetail;

  return (
    <div className="author-detail-page">
      {/* Geri Dön Butonu */}
      <button
        className="author-back-button"
        onClick={() => navigate('/authors')}
      >
        ← Yazarlara Dön
      </button>

      <div className="author-detail-container">
        {/* Ana İçerik */}
        <div className="author-main-content">
          <div className="author-info-card">
            {/* Yazar Bilgileri */}
            <div className="author-header">
              <div className="author-avatar">
                👤
              </div>
              <div className="author-title-section">
                <h1>{author.name}</h1>
                <div className="author-chips">
                  <div className="author-chip primary">
                    📝 Yazar
                  </div>
                  <div className="author-chip secondary">
                    📚 {book_count} Kitap
                  </div>
                </div>
              </div>
            </div>

            <hr className="author-divider" />

            {/* Biyografi */}
            <div className="author-bio-section">
              <h6>Yazar Hakkında</h6>
              <p className="author-bio-text">
                {`${author.name} hakkında detaylı bilgi yakında eklenecektir. Bu yazar toplam ${book_count} kitap yazmıştır ve kütüphanemizde eserleri bulunmaktadır.`}
              </p>
            </div>
          </div>

          {/* Yazarın Kitapları */}
          <div className="author-books-section">
            <h5>Bu Yazarın Kitapları ({book_count} kitap)</h5>
            {books && books.length > 0 ? (
              <div className="author-books-grid">
                {books.map((book) => (
                  <BookCard 
                    key={book.id}
                    book={book} 
                    showAuthor={false} 
                    showCategory={true}
                  />
                ))}
              </div>
            ) : (
              <div className="no-books-card">
                <h6>Bu yazara ait kitap bulunmuyor</h6>
              </div>
            )}
          </div>
        </div>

        {/* Yan Panel */}
        <div className="author-sidebar">
          <div className="author-info-panel">
            <h6>Yazar Bilgileri</h6>

            <div className="author-info-items">
              <div className="author-info-item">
                <span className="author-info-label">Toplam Kitap</span>
                <span className="author-info-value large">{book_count}</span>
              </div>

              <div className="author-info-item">
                <span className="author-info-label">Yazar Adı</span>
                <span className="author-info-value">{author.name}</span>
              </div>

              <div className="author-info-item">
                <span className="author-info-label">Durum</span>
                <span className="author-info-value">Aktif Yazar</span>
              </div>
            </div>

            <div className="author-action-buttons">
              <Link
                to={`/books?author=${encodeURIComponent(author.name)}`}
                className="author-btn primary"
              >
                📚 Tüm Kitaplarını Gör
              </Link>

              <button
                className="author-btn outlined"
                onClick={() => navigate('/authors')}
              >
                ✍️ Diğer Yazarlar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuthorDetailPage; 