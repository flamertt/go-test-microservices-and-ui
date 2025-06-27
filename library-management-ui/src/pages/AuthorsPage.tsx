import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { usePagination } from '../hooks/usePagination';
import type { Author } from '../types';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import '../styles/AuthorsPage.css';

const AuthorsPage: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [searchInput, setSearchInput] = useState('');

  // Pagination hook
  const {
    data: authors,
    loading,
    error,
    page,
    totalPages,
    total,
    pageSize,
    setPage,
    refetch
  } = usePagination<Author>({
    endpoint: '/api/authors',
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

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .substring(0, 2);
  };

  const renderPagination = () => {
    const pageNumbers = [];
    const startPage = Math.max(1, page - 2);
    const endPage = Math.min(totalPages, page + 2);

    for (let i = startPage; i <= endPage; i++) {
      pageNumbers.push(
        <button
          key={i}
          className={`page-num ${i === page ? 'active' : ''}`}
          onClick={() => setPage(i)}
        >
          {i}
        </button>
      );
    }

    return (
      <div className="authors-pagination">
        <button
          className="pagination-btn"
          onClick={() => setPage(page - 1)}
          disabled={page <= 1}
        >
          ← Önceki
        </button>
        
        <div className="page-numbers-authors">
          {pageNumbers}
        </div>
        
        <button
          className="pagination-btn"
          onClick={() => setPage(page + 1)}
          disabled={page >= totalPages}
        >
          Sonraki →
        </button>
        
        <div className="pagination-info-authors">
          Sayfa {page} / {totalPages} (Toplam {total} yazar)
        </div>
      </div>
    );
  };

  if (error) {
    return <ErrorMessage error={error} />;
  }

  return (
    <div className="authors-page">
      <div className="authors-header">
        <h1 className="authors-title">✍️ Yazarlar</h1>
        <p className="authors-subtitle">
          Kütüphanemizdeki yazarları keşfedin. Biografilerini okuyun ve eserlerini görüntüleyin.
        </p>
      </div>

      {/* Arama */}
      <div className="authors-search">
        <div className="search-grid">
          <div className="search-input-group">
            <label className="search-label">Yazar Ara</label>
            <input
              type="text"
              className="author-search-input"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Yazar adına göre ara..."
            />
          </div>
          
          <div className="search-buttons">
            <button className="btn-search" onClick={handleSearch}>
              🔍 Ara
            </button>
            <button className="btn-reset" onClick={handleReset}>
              🔄 Temizle
            </button>
          </div>
        </div>
      </div>

      {/* Loading state */}
      {loading && (
        <div className="authors-loading">
          <div className="authors-spinner"></div>
        </div>
      )}

      {/* Sonuçlar */}
      {!loading && (
        <>
          {renderPagination()}

          {!authors || authors.length === 0 ? (
            <div className="authors-empty">
              <div className="empty-icon">👤</div>
              <h3 className="authors-empty-title">Aradığınız kriterlere uygun yazar bulunamadı</h3>
              <p className="authors-empty-subtitle">Arama terimlerinizi değiştirerek tekrar deneyin</p>
            </div>
          ) : (
            <div className="authors-grid">
              {authors.map((author) => (
                <div key={author.id} className="author-card">
                  <div className="author-content">
                    <div className="author-header">
                      <div className="author-avatar">
                        {getInitials(author.name)}
                      </div>
                      <div className="author-info">
                        <h3 className="author-name">{author.name}</h3>
                      </div>
                    </div>

                    <p className="author-description">
                      Kütüphanemizde bulunan değerli yazarlarımızdan biri.
                    </p>

                    <div className="author-actions">
                      <Link
                        to={`/authors/${encodeURIComponent(author.name)}`}
                        className="author-view-btn"
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

export default AuthorsPage; 