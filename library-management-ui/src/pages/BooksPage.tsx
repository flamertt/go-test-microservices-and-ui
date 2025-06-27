import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { booksApi, genresApi } from '../services/api';
import type { Book, Genre, PaginatedResponse } from '../types';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import BookCard from '../components/BookCard';
import '../styles/BooksPage.css';

const BooksPage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedGenre, setSelectedGenre] = useState<string>('');
  const [selectedAuthor, setSelectedAuthor] = useState<string>('');
  const [searchInput, setSearchInput] = useState('');
  
  // Books data state
  const [books, setBooks] = useState<Book[]>([]);
  const [booksLoading, setBooksLoading] = useState(false);
  const [booksError, setBooksError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [total, setTotal] = useState(0);
  
  // Genres data state
  const [genres, setGenres] = useState<Genre[]>([]);
  const [genresLoading, setGenresLoading] = useState(false);

  // Fetch books function
  const fetchBooks = async () => {
    try {
      setBooksLoading(true);
      setBooksError(null);
      
      // API search endpoint'ini kullan
      let endpoint = '/books';
      const params = new URLSearchParams();
      
      if (searchTerm) {
        endpoint = `/books/search`;
        params.append('q', searchTerm);
      } else if (selectedGenre) {
        endpoint = `/books/category/${encodeURIComponent(selectedGenre)}`;
      } else if (selectedAuthor) {
        endpoint = `/books/author/${encodeURIComponent(selectedAuthor)}`;
      }
      
      // Sayfalama parametreleri
      params.append('page', page.toString());
      params.append('page_size', '20');
      
      const response = await fetch(`/api${endpoint}?${params.toString()}`);
      
      if (!response.ok) {
        throw new Error(`API Hatası: ${response.status}`);
      }
      
      const data = await response.json();
      setBooks(data.data.books || data.data || []);
      setTotalPages(data.data.total_pages || 1);
      setTotal(data.data.total || 0);
    } catch (error) {
      setBooksError(error instanceof Error ? error.message : 'Kitaplar yüklenirken hata oluştu');
    } finally {
      setBooksLoading(false);
    }
  };

  // Fetch genres function
  const fetchGenres = async () => {
    try {
      setGenresLoading(true);
      const response = await genresApi.getAll();
      
      // Response data'nın array olduğundan emin ol
      if (response && response.data) {
        const data: any = response.data;
        if (Array.isArray(data)) {
          setGenres(data);
        } else if (data.genres && Array.isArray(data.genres)) {
          setGenres(data.genres);
        } else if (data.data && Array.isArray(data.data)) {
          setGenres(data.data);
        } else {
          console.warn('Genres response format unexpected, setting empty array');
          setGenres([]);
        }
      } else {
        setGenres([]);
      }
    } catch (error) {
      console.error('Genres loading error:', error);
      setGenres([]); // Hata durumunda boş array
    } finally {
      setGenresLoading(false);
    }
  };

  // URL parametrelerini oku
  useEffect(() => {
    const genreParam = searchParams.get('genre');
    const authorParam = searchParams.get('author');
    const searchParam = searchParams.get('search');

    if (genreParam) {
      setSelectedGenre(genreParam);
      setSearchInput(''); // Tür filtresi varsa arama temizle
    }
    if (authorParam) {
      setSelectedAuthor(authorParam);
      setSearchInput(''); // Yazar filtresi varsa arama temizle
    }
    if (searchParam) {
      setSearchInput(searchParam);
      setSearchTerm(searchParam);
    }
  }, [searchParams]);

  // Fetch data when filters change
  useEffect(() => {
    fetchBooks();
  }, [page, searchTerm, selectedGenre, selectedAuthor]);

  // Fetch genres on mount
  useEffect(() => {
    fetchGenres();
  }, []);

  // Reset page when filters change
  useEffect(() => {
    setPage(1);
  }, [searchTerm, selectedGenre, selectedAuthor]);

  const handleGenreChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const genre = event.target.value;
    setSelectedGenre(genre);
    setSelectedAuthor(''); // Tür seçilince yazarı temizle
    
    // URL'yi güncelle
    const newParams = new URLSearchParams();
    if (genre) newParams.set('genre', genre);
    if (searchTerm) newParams.set('search', searchTerm);
    setSearchParams(newParams);
  };

  const handleSearch = () => {
    setSearchTerm(searchInput);
    setSelectedGenre(''); // Arama yapılınca filtreleri temizle
    setSelectedAuthor('');
    
    // URL'yi güncelle
    const newParams = new URLSearchParams();
    if (searchInput.trim()) newParams.set('search', searchInput.trim());
    setSearchParams(newParams);
  };

  const handleKeyPress = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };

  const handleReset = () => {
    setSearchInput('');
    setSearchTerm('');
    setSelectedGenre('');
    setSelectedAuthor('');
    setSearchParams(new URLSearchParams()); // URL'yi temizle
  };

  const renderFiltersInfo = () => {
    const activeFilters = [];
    if (selectedGenre) activeFilters.push(`Kategori: ${selectedGenre}`);
    if (selectedAuthor) activeFilters.push(`Yazar: ${selectedAuthor}`);
    if (searchTerm) activeFilters.push(`Arama: "${searchTerm}"`);

    if (activeFilters.length === 0) return null;

    return (
      <div className="active-filters">
        <span className="active-filters-label">🎯 Aktif Filtreler:</span>
        {activeFilters.map((filter, index) => (
          <span key={index} className="active-filter-tag">
            {filter}
          </span>
        ))}
        <button className="clear-filters-btn" onClick={handleReset}>
          ✖ Tümünü Temizle
        </button>
      </div>
    );
  };

  const renderPagination = () => {
    if (totalPages <= 1) return null;

    const pageNumbers = [];
    const startPage = Math.max(1, page - 2);
    const endPage = Math.min(totalPages, page + 2);

    for (let i = startPage; i <= endPage; i++) {
      pageNumbers.push(
        <button
          key={i}
          className={`books-page-num ${i === page ? 'active' : ''}`}
          onClick={() => setPage(i)}
        >
          {i}
        </button>
      );
    }

    return (
      <div className="books-pagination">
        <button
          className="books-pagination-btn"
          onClick={() => setPage(page - 1)}
          disabled={page <= 1}
        >
          ← Önceki
        </button>
        
        <div className="books-page-numbers">
          {pageNumbers}
        </div>
        
        <button
          className="books-pagination-btn"
          onClick={() => setPage(page + 1)}
          disabled={page >= totalPages}
        >
          Sonraki →
        </button>
        
        <div className="books-pagination-info">
          Sayfa {page} / {totalPages} (Toplam {total} kitap)
        </div>
      </div>
    );
  };

  if (genresLoading) {
    return <LoadingSpinner message="Filtreler yükleniyor..." />;
  }

  if (booksError) {
    return <ErrorMessage error={booksError} />;
  }

  return (
    <div className="books-page">
      <div className="books-header">
        <h1 className="books-title">📚 Kitaplar</h1>
        <p className="books-subtitle">
          Kütüphanemizdeki tüm kitapları keşfedin. Arama yapabilir, yazara veya türe göre filtreleyebilirsiniz.
        </p>
      </div>

      {/* Aktif Filtreler */}
      {renderFiltersInfo()}

      {/* Filtreler */}
      <div className="books-filters">
        <div className="books-filters-grid">
          <div className="books-input-group">
            <label className="books-input-label">Kitap Ara</label>
            <input
              type="text"
              className="books-search-input"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Kitap adı, yazar, yayınevi..."
            />
          </div>
          
          <div className="books-input-group">
            <label className="books-input-label">Kategori</label>
            <select
              className="books-select-input"
              value={selectedGenre}
              onChange={handleGenreChange}
            >
              <option value="">Tüm Kategoriler</option>
              {Array.isArray(genres) && genres.map((genre) => (
                <option key={genre.id} value={genre.name}>
                  {genre.name}
                </option>
              ))}
            </select>
          </div>
          
          <div className="books-button-group">
            <button className="books-btn-primary" onClick={handleSearch}>
              🔍 Ara
            </button>
            <button className="books-btn-secondary" onClick={handleReset}>
              🔄 Temizle
            </button>
          </div>
        </div>
      </div>

      {/* Loading state */}
      {booksLoading && (
        <div className="books-loading">
          <LoadingSpinner message="Kitaplar yükleniyor..." />
        </div>
      )}

      {/* Sonuçlar */}
      {!booksLoading && (
        <>


          {renderPagination()}

          {!books || books.length === 0 ? (
            <div className="books-empty">
              <div className="books-empty-icon">📚</div>
              <h3 className="books-empty-title">Aradığınız kriterlere uygun kitap bulunamadı</h3>
              <p className="books-empty-subtitle">Arama terimlerinizi değiştirerek tekrar deneyin</p>
            </div>
          ) : (
            <div className="books-grid">
              {books.map((book) => (
                <BookCard key={book.id} book={book} showAuthor={true} showCategory={true} />
              ))}
            </div>
          )}

          {renderPagination()}
        </>
      )}
    </div>
  );
};

export default BooksPage; 