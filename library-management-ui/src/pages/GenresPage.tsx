import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { usePagination } from '../hooks/usePagination';
import type { Genre } from '../types';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import '../styles/GenresPage.css';

const GenresPage: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [searchInput, setSearchInput] = useState('');

  // Pagination hook for genres
  const {
    data: genres,
    loading: genresLoading,
    error: genresError,
    page,
    totalPages,
    total,
    pageSize,
    setPage,
    refetch
  } = usePagination<Genre>({
    endpoint: '/api/genres',
    searchTerm,
  });

  const handleSearch = () => {
    setSearchTerm(searchInput);
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  const handleReset = () => {
    setSearchInput('');
    setSearchTerm('');
  };

  const renderPagination = () => {
    const pageNumbers = [];
    const startPage = Math.max(1, page - 2);
    const endPage = Math.min(totalPages, page + 2);

    for (let i = startPage; i <= endPage; i++) {
      pageNumbers.push(
        <button
          key={i}
          className={`page-num-genres ${i === page ? 'active' : ''}`}
          onClick={() => setPage(i)}
        >
          {i}
        </button>
      );
    }

    return (
      <div className="genres-pagination">
        <button
          className="pagination-btn-genres"
          onClick={() => setPage(page - 1)}
          disabled={page <= 1}
        >
          ← Önceki
        </button>
        
        <div className="page-numbers-genres">
          {pageNumbers}
        </div>
        
        <button
          className="pagination-btn-genres"
          onClick={() => setPage(page + 1)}
          disabled={page >= totalPages}
        >
          Sonraki →
        </button>
        
        <div className="pagination-info-genres">
          Sayfa {page} / {totalPages} (Toplam {total} kategori)
        </div>
      </div>
    );
  };

  if (genresError) {
    return <ErrorMessage error={genresError} />;
  }

  return (
    <div className="genres-page">
      <div className="genres-header">
        <h1 className="genres-title">🎭 Kitap Türleri</h1>
        <p className="genres-subtitle">
          Kütüphanemizdeki kitap türlerini keşfedin ve kategorilere göre kitapları bulun.
        </p>
      </div>

      {/* Arama */}
      <div className="genres-search">
        <div className="genre-search-grid">
          <div className="genre-search-group">
            <label className="genre-search-label">Kategori Ara</label>
            <input
              type="text"
              className="genre-search-input"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Kategori adına göre ara..."
            />
          </div>
          
          <div className="genre-search-buttons">
            <button className="btn-genre-search" onClick={handleSearch}>
              🔍 Ara
            </button>
            <button className="btn-genre-reset" onClick={handleReset}>
              🔄 Temizle
            </button>
          </div>
        </div>
      </div>

      {/* Loading state */}
      {genresLoading && (
        <div className="genres-loading">
          <div className="genres-spinner"></div>
        </div>
      )}

      {/* Sonuçlar */}
      {!genresLoading && (
        <>
          {renderPagination()}

          {!genres || genres.length === 0 ? (
            <div className="genres-empty">
              <div className="genres-empty-icon">🎭</div>
              <h3 className="genres-empty-title">Aradığınız kriterlere uygun tür bulunamadı</h3>
              <p className="genres-empty-subtitle">Arama terimlerinizi değiştirerek tekrar deneyin</p>
            </div>
          ) : (
            <div className="genres-grid">
              {genres.map((genre, index) => (
                <div key={genre.id} className="genre-card">
                  <div className="genre-content">
                    <div className="genre-header">
                      <div className="genre-icon">
                        🏷️
                      </div>
                      <h3 className="genre-name">{genre.name}</h3>
                    </div>

                    <p className="genre-description">
                      Bu kategorideki kitapları keşfetmek için aşağıdaki butona tıklayın.
                    </p>

                    <div className="genre-badge">
                      📚 Kitapları gör
                    </div>

                    <div className="genre-actions">
                      <Link
                        to={`/genres/${encodeURIComponent(genre.name)}`}
                        className="genre-view-btn"
                      >
                        📖 Detayları Görüntüle
                      </Link>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {renderPagination()}
        </>
      )}
    </div>
  );
};

export default GenresPage; 