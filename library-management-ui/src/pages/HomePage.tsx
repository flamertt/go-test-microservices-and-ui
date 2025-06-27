import React from 'react';
import { Link } from 'react-router-dom';
import '../styles/HomePage.css';

const HomePage: React.FC = () => {
  const features = [
    {
      title: 'Kitaplar',
      description: 'Kütüphanedeki tüm kitapları görüntüleyin ve arayın',
      icon: '📚',
      link: '/books',
    },
    {
      title: 'Yazarlar',
      description: 'Yazarları keşfedin ve eserlerini görüntüleyin',
      icon: '✍️',
      link: '/authors',
    },
    {
      title: 'Türler',
      description: 'Kitap türlerini keşfedin ve kategorilere göz atın',
      icon: '🎭',
      link: '/genres',
    },
    {
      title: 'Öneriler',
      description: 'Size özel kitap önerileri alın',
      icon: '🤖',
      link: '/recommendations',
    },
  ];

  return (
    <div className="home-page">
      {/* Hero Section */}
      <div className="hero-section">
        <div className="hero-content">
          <div className="hero-icon">📖</div>
          <h1 className="hero-title">Kütüphane</h1>
          <p className="hero-subtitle">Kitapları keşfedin, yazarları tanıyın</p>
          <Link to="/books" className="hero-cta">
            Başlayın
        </Link>
        </div>
      </div>

      {/* Features Section */}
      <div className="features-section">
        <div className="features-grid">
          {features.map((feature, index) => (
            <Link key={index} to={feature.link} className="feature-card">
              <div className="feature-icon">{feature.icon}</div>
              <h3 className="feature-title">{feature.title}</h3>
              <p className="feature-description">{feature.description}</p>
              </Link>
          ))}
        </div>
      </div>
      </div>
  );
};

export default HomePage; 