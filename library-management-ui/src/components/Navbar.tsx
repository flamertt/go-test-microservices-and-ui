import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import '../styles/Navbar.css';

const Navbar: React.FC = () => {
  const location = useLocation();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  const isActive = (path: string) => location.pathname === path;

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  useEffect(() => {
    const handleScroll = () => {
      const isScrolled = window.scrollY > 20;
      setScrolled(isScrolled);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <nav className={`navbar ${scrolled ? 'scrolled' : ''}`}>
      <div className="navbar-container">
        <Link to="/" className="navbar-brand">
          <div className="brand-icon">🏛️</div>
          <span className="brand-text">Kütüphane Yönetim Sistemi</span>
        </Link>

        <div className="navbar-menu">
          <Link
            to="/books"
            className={`nav-link ${isActive('/books') ? 'active' : ''}`}
          >
            <span className="nav-icon">📚</span>
            Kitaplar
          </Link>
          
          <Link
            to="/authors"
            className={`nav-link ${isActive('/authors') ? 'active' : ''}`}
          >
            <span className="nav-icon">✍️</span>
            Yazarlar
          </Link>
          
          <Link
            to="/genres"
            className={`nav-link ${isActive('/genres') ? 'active' : ''}`}
          >
            <span className="nav-icon">🎭</span>
            Türler
          </Link>
          
          <Link
            to="/recommendations"
            className={`nav-link ${isActive('/recommendations') ? 'active' : ''}`}
          >
            <span className="nav-icon">🤖</span>
            Öneriler
          </Link>
        </div>

        <button 
          className="mobile-menu-toggle"
          onClick={toggleMobileMenu}
          aria-label="Mobile menüyü aç/kapat"
        >
          {mobileMenuOpen ? '✕' : '☰'}
        </button>

        <div className={`mobile-menu ${mobileMenuOpen ? 'open' : ''}`}>
          <Link
            to="/books"
            className={`mobile-nav-link ${isActive('/books') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <span className="nav-icon">📚</span>
            Kitaplar
          </Link>
          
          <Link
            to="/authors"
            className={`mobile-nav-link ${isActive('/authors') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <span className="nav-icon">✍️</span>
            Yazarlar
          </Link>
          
          <Link
            to="/genres"
            className={`mobile-nav-link ${isActive('/genres') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <span className="nav-icon">🎭</span>
            Türler
          </Link>
          
          <Link
            to="/recommendations"
            className={`mobile-nav-link ${isActive('/recommendations') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <span className="nav-icon">🤖</span>
            Öneriler
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default Navbar; 