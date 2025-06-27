import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useApi } from '../hooks/useApi';
import type { Book, Genre } from '../types';
import BookCard from '../components/BookCard';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import '../styles/GenreDetailPage.css';

// Genre Detail Response type
interface GenreDetailResponse {
  genre: Genre;
  books: Book[];
  book_count: number;
  page: number;
  page_size: number;
  total_pages: number;
  total: number;
}

const GenreDetailPage: React.FC = () => {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const genreName = decodeURIComponent(name || '');
  
  // Manuel sayfalama state
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 20;

  // Genre detail API (sayfalama parametreleri ile)
  const { data: genreDetailData, loading, error } = useApi<GenreDetailResponse>(
    `/api/genres/detail/${encodeURIComponent(genreName)}?page=${currentPage}&page_size=${pageSize}`,
    [genreName, currentPage]
  );

  const getGenreIcon = (name: string) => {
    // Kategori adına göre farklı renkler
    const colorMap: { [key: string]: string } = {
      'literature': 'literature',
      'science': 'science', 
      'history': 'history',
      'philosophy': 'philosophy',
      'psychology': 'psychology',
      'academic': 'academic',
      'fiction': 'fiction',
      'non-fiction': 'non-fiction'
    };
    
    const lowerName = name.toLowerCase();
    const colorClass = Object.keys(colorMap).find(key => 
      lowerName.includes(key)
    );
    
    return colorClass || '';
  };

  const renderPagination = () => {
    if (!genreDetailData || genreDetailData.total_pages <= 1) return null;

    const { page, total_pages, total } = genreDetailData;
    const pageNumbers = [];
    const startPage = Math.max(1, page - 2);
    const endPage = Math.min(total_pages, page + 2);

    for (let i = startPage; i <= endPage; i++) {
      pageNumbers.push(
        <button
          key={i}
          className={`genre-page-number ${i === page ? 'active' : ''}`}
          onClick={() => setCurrentPage(i)}
        >
          {i}
        </button>
      );
    }

    return (
      <div className="genre-pagination">
        <button
          className="genre-pagination-btn"
          onClick={() => setCurrentPage(page - 1)}
          disabled={page <= 1}
        >
          ← Önceki
        </button>
        
        {pageNumbers}
        
        <button
          className="genre-pagination-btn"
          onClick={() => setCurrentPage(page + 1)}
          disabled={page >= total_pages}
        >
          Sonraki →
        </button>
        
        <span className="genre-pagination-info">
          Sayfa {page} / {total_pages} (Toplam {total} kitap)
        </span>
      </div>
    );
  };

  if (loading) {
    return <LoadingSpinner message="Tür detayları yükleniyor..." />;
  }

  if (error || !genreDetailData) {
    return (
      <div className="genre-error-container">
        <button
          className="genre-back-button"
          onClick={() => navigate('/genres')}
        >
          ← Türlere Dön
        </button>
        <ErrorMessage 
          error={error || 'Tür bulunamadı'} 
          title="Tür Detayları Yüklenemiyor"
        />
      </div>
    );
  }

  // API response'dan verileri al
  const { genre, books = [], book_count = 0, page = 1, total_pages = 1 } = genreDetailData;
  
  // Genre objesinin varlığını kontrol et
  if (!genre || !genre.name) {
    return (
      <div className="genre-error-container">
        <button
          className="genre-back-button"
          onClick={() => navigate('/genres')}
        >
          ← Türlere Dön
        </button>
        <ErrorMessage 
          error="Tür verisi bulunamadı" 
          title="Tür Detayları Yüklenemiyor"
        />
      </div>
    );
  }

  return (
    <div className="genre-detail-page">
      {/* Geri Dön Butonu */}
      <button
        className="genre-back-button"
        onClick={() => navigate('/genres')}
      >
        ← Türlere Dön
      </button>

      <div className="genre-detail-container">
        {/* Ana İçerik */}
        <div className="genre-main-content">
          <div className="genre-info-card">
            {/* Tür Bilgileri */}
            <div className="genre-header">
              <div className={`genre-avatar ${getGenreIcon(genre.name)}`}>
                📚
              </div>
              <div className="genre-title-section">
                <h1>{genre.name}</h1>
                <div className="genre-chips">
                  <div className="genre-chip primary">
                    📖 Kategori
                  </div>
                  <div className="genre-chip secondary">
                    📚 {book_count} Kitap
                  </div>
                </div>
              </div>
            </div>

            <hr className="genre-divider" />

            {/* Açıklama */}
            <div className="genre-description-section">
              <h6>Kategori Hakkında</h6>
              <p className="genre-description-text">
                {`${genre.name} kategorisinde toplam ${book_count} kitap bulunmaktadır. Bu kategorideki kitaplar çeşitli yazarlar tarafından kaleme alınmış olup, farklı konuları ele almaktadır. Aşağıda bu kategoriye ait kitapları sayfalı olarak görüntüleyebilirsiniz.`}
              </p>
            </div>
          </div>

          {/* Kategorinin Kitapları */}
          <div className="genre-books-section">
            <h5>Bu Kategorideki Kitaplar ({book_count} kitap - Sayfa {page})</h5>
            
            {/* Üst Sayfalama */}
            {renderPagination()}
            
            {books && books.length > 0 ? (
              <>
                <div className="genre-books-grid">
                  {books.map((book) => (
                    <BookCard 
                      key={book.id}
                      book={book} 
                      showAuthor={true} 
                      showCategory={false}
                    />
                  ))}
                </div>
                
                {/* Alt Sayfalama */}
                {renderPagination()}
              </>
            ) : (
              <div className="no-books-card">
                <h6>Bu kategoride kitap bulunmuyor</h6>
              </div>
            )}
          </div>
        </div>

        {/* Yan Panel */}
        <div className="genre-sidebar">
          <div className="genre-info-panel">
            <h6>Kategori Bilgileri</h6>

            <div className="genre-info-items">
              <div className="genre-info-item">
                <span className="genre-info-label">Toplam Kitap</span>
                <span className="genre-info-value large">{book_count}</span>
              </div>

              <div className="genre-info-item">
                <span className="genre-info-label">Şu Anki Sayfa</span>
                <span className="genre-info-value">{page} / {total_pages}</span>
              </div>

              <div className="genre-info-item">
                <span className="genre-info-label">Kategori Adı</span>
                <span className="genre-info-value">{genre.name}</span>
              </div>

              <div className="genre-info-item">
                <span className="genre-info-label">Durum</span>
                <span className="genre-info-value">Aktif Kategori</span>
              </div>
            </div>

            <div className="genre-action-buttons">
              <Link
                to={`/books?genre=${encodeURIComponent(genre.name)}`}
                className="genre-btn primary"
              >
                📚 Tüm Kitapları Filtrele
              </Link>

              <button
                className="genre-btn outlined"
                onClick={() => navigate('/genres')}
              >
                📖 Diğer Kategoriler
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default GenreDetailPage; 