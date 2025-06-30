import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { 
  FaBook, 
  FaUsers, 
  FaTheaterMasks, 
  FaRobot,
  FaUser,
  FaSignOutAlt,
  FaSignInAlt,
  FaUserPlus,
  FaBars,
  FaTimes,
  FaBuilding
} from 'react-icons/fa';
import { HiSparkles } from 'react-icons/hi';
import '../styles/Navbar.css';

const Navbar: React.FC = () => {
  const location = useLocation();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const { user, isAuthenticated, logout } = useAuth();

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
          <div className="brand-icon">
            <FaBuilding />
          </div>
          <span className="brand-text">Kütüphane Yönetim Sistemi</span>
        </Link>

        <div className="navbar-menu">
          <Link
            to="/books"
            className={`nav-link books-link ${isActive('/books') ? 'active' : ''}`}
          >
            <FaBook className="nav-icon" />
            Kitaplar
          </Link>
          
          <Link
            to="/authors"
            className={`nav-link authors-link ${isActive('/authors') ? 'active' : ''}`}
          >
            <FaUsers className="nav-icon" />
            Yazarlar
          </Link>
          
          <Link
            to="/genres"
            className={`nav-link genres-link ${isActive('/genres') ? 'active' : ''}`}
          >
            <FaTheaterMasks className="nav-icon" />
            Türler
          </Link>
          
          <Link
            to="/recommendations"
            className={`nav-link recommendations-link ${isActive('/recommendations') ? 'active' : ''}`}
          >
            <div className="recommendation-icon">
              <HiSparkles className="nav-icon sparkle-icon" />
              <FaRobot className="nav-icon robot-icon" />
            </div>
            Öneriler
          </Link>
        </div>

        {/* Auth Buttons */}
        <div className="auth-section">
          {isAuthenticated ? (
            <div className="user-menu">
              <Link to="/profile" className="user-profile">
                <div className="user-avatar">
                  <FaUser />
                </div>
                <span className="user-name">{user?.username}</span>
              </Link>
              <button onClick={logout} className="logout-btn">
                <FaSignOutAlt />
                Çıkış
              </button>
            </div>
          ) : (
            <div className="auth-buttons">
              <Link to="/login" className="login-btn">
                <FaSignInAlt className="btn-icon" />
                Giriş Yap
              </Link>
              <Link to="/register" className="register-btn">
                <FaUserPlus className="btn-icon" />
                Kayıt Ol
              </Link>
            </div>
          )}
        </div>

        <button 
          className="mobile-menu-toggle"
          onClick={toggleMobileMenu}
          aria-label="Mobile menüyü aç/kapat"
        >
          {mobileMenuOpen ? <FaTimes /> : <FaBars />}
        </button>

        <div className={`mobile-menu ${mobileMenuOpen ? 'open' : ''}`}>
          <Link
            to="/books"
            className={`mobile-nav-link books-link ${isActive('/books') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <FaBook className="nav-icon" />
            Kitaplar
          </Link>
          
          <Link
            to="/authors"
            className={`mobile-nav-link authors-link ${isActive('/authors') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <FaUsers className="nav-icon" />
            Yazarlar
          </Link>
          
          <Link
            to="/genres"
            className={`mobile-nav-link genres-link ${isActive('/genres') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <FaTheaterMasks className="nav-icon" />
            Türler
          </Link>
          
          <Link
            to="/recommendations"
            className={`mobile-nav-link recommendations-link ${isActive('/recommendations') ? 'active' : ''}`}
            onClick={() => setMobileMenuOpen(false)}
          >
            <div className="recommendation-icon">
              <HiSparkles className="nav-icon sparkle-icon" />
              <FaRobot className="nav-icon robot-icon" />
            </div>
            Öneriler
          </Link>
          
          {/* Mobile Auth Section */}
          <div className="mobile-auth-section">
            {isAuthenticated ? (
              <>
                <Link
                  to="/profile"
                  className="mobile-nav-link"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  <FaUser className="nav-icon" />
                  Profil ({user?.username})
                </Link>
                <button
                  onClick={() => {
                    logout();
                    setMobileMenuOpen(false);
                  }}
                  className="mobile-nav-link logout-mobile"
                >
                  <FaSignOutAlt className="nav-icon" />
                  Çıkış Yap
                </button>
              </>
            ) : (
              <>
                <Link
                  to="/login"
                  className="mobile-nav-link"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  <FaSignInAlt className="nav-icon" />
                  Giriş Yap
                </Link>
                <Link
                  to="/register"
                  className="mobile-nav-link register-mobile"
                  onClick={() => setMobileMenuOpen(false)}
                >
                  <FaUserPlus className="nav-icon" />
                  Kayıt Ol
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar; 