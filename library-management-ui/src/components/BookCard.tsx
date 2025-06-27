import React from 'react';
import { Link } from 'react-router-dom';
import type { Book } from '../types';
import '../styles/Components.css';

interface BookCardProps {
  book: Book;
  showAuthor?: boolean;
  showCategory?: boolean;
}

const BookCard: React.FC<BookCardProps> = ({ 
  book, 
  showAuthor = true, 
  showCategory = true 
}) => {
  return (
    <div className="book-card">
      <div className="book-card-content">
        <h3 className="book-card-title">{book.title}</h3>
        
        <div className="book-card-meta">
          {showAuthor && (
            <div className="book-card-author">
              👤 <span>{book.author}</span>
            </div>
          )}
          
          {showCategory && (
            <div className="book-card-category">
              📖 <span>{book.category_name}</span>
            </div>
          )}
          
          <div className="book-card-year">
            📅 <span>{book.released_year}</span>
          </div>
        </div>

        <div className="book-card-actions">
          <Link
            to={`/books/${book.id}`}
            className="book-card-btn primary"
          >
            Detayları Gör
          </Link>
          
          {showAuthor && (
            <Link
              to={`/authors/${encodeURIComponent(book.author)}`}
              className="book-card-btn outlined"
            >
              Yazar
            </Link>
          )}
        </div>
      </div>
    </div>
  );
};

export default BookCard; 