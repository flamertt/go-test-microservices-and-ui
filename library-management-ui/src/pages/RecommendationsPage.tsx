import React, { useState, useEffect } from 'react';
import { recommendationsApi, type Recommendation } from '../services/api';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';
import BookCard from '../components/BookCard';
import '../styles/RecommendationsPage.css';

const RecommendationsPage: React.FC = () => {
  const [recommendations, setRecommendations] = useState<Recommendation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filterType, setFilterType] = useState<'general' | 'category' | 'author'>('general');
  const [filterValue, setFilterValue] = useState('');
  const [generatedAt, setGeneratedAt] = useState<string>('');
  const [selectedInfo, setSelectedInfo] = useState<string>(''); // Seçilen kategori/yazar bilgisi

  const fetchRecommendations = async () => {
    try {
      setLoading(true);
      setError(null);
      
      let response;
      
      switch (filterType) {
        case 'category':
          if (!filterValue.trim()) {
            // Random kategori önerisi
            response = await recommendationsApi.getRandomCategory(15);
            setRecommendations(response.data.recommendations);
            setSelectedInfo(`Rastgele seçilen kategori: ${response.data.category}`);
          } else {
            // Belirli kategori önerisi
            response = await recommendationsApi.getByCategory(filterValue.trim());
            setRecommendations(response.data.recommendations);
            setSelectedInfo(`Belirli kategori: ${response.data.category}`);
          }
          break;
          
        case 'author':
          if (!filterValue.trim()) {
            // Random yazar önerisi
            response = await recommendationsApi.getRandomAuthor(15);
            setRecommendations(response.data.recommendations);
            setSelectedInfo(`Rastgele seçilen yazar: ${response.data.author}`);
          } else {
            // Belirli yazar önerisi
            response = await recommendationsApi.getByAuthor(filterValue.trim());
            setRecommendations(response.data.recommendations);
            setSelectedInfo(`Belirli yazar: ${response.data.author}`);
          }
          break;
          
        default:
          // Random genel öneriler
          response = await recommendationsApi.getGeneral(15);
          setRecommendations(response.data.recommendations);
          setGeneratedAt(response.data.generated_at);
          setSelectedInfo('');
      }
    } catch (err) {
      console.error('API Error:', err);
      setError(
        err instanceof Error 
          ? `API Hatası: ${err.message}` 
          : 'Öneriler yüklenirken hata oluştu. Lütfen tekrar deneyin.'
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecommendations();
  }, [filterType]);

  const handleFilterSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    fetchRecommendations();
  };

  const handleQuickCategorySelect = (category: string) => {
    setFilterValue(category);
    setFilterType('category');
  };

  const handleQuickAuthorSelect = (author: string) => {
    setFilterValue(author);
    setFilterType('author');
  };

  const getRandomRecommendations = () => {
    setFilterValue(''); // Temizle ki random seçim yapsın
    fetchRecommendations();
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;

  return (
    <div className="recommendations-page">
      <div className="recommendations-header">
        <h1 className="recommendations-title">📚 Kitap Önerileri</h1>
        <p className="recommendations-subtitle">
          Size özel kitap önerileri keşfedin. Kategori, yazar veya genel önerileri görüntüleyebilirsiniz.
        </p>
        
        {generatedAt && (
          <p className="recommendations-updated">
            📅 Son güncelleme: {new Date(generatedAt).toLocaleString('tr-TR')}
          </p>
        )}
        {selectedInfo && (
          <p className="recommendations-filter-info">
            🎯 Seçilen: {selectedInfo}
          </p>
        )}
      </div>

      <div className="recommendations-filters">
        <div className="recommendations-tabs">
          <button 
            className={`recommendations-tab ${filterType === 'general' ? 'active' : ''}`}
            onClick={() => setFilterType('general')}
          >
            🎲 Genel Öneriler
          </button>
          <button 
            className={`recommendations-tab ${filterType === 'category' ? 'active' : ''}`}
            onClick={() => setFilterType('category')}
          >
            📖 Kategori Bazlı
          </button>
          <button 
            className={`recommendations-tab ${filterType === 'author' ? 'active' : ''}`}
            onClick={() => setFilterType('author')}
          >
            ✍️ Yazar Bazlı
          </button>
          <button 
            className="recommendations-tab recommendations-refresh-tab"
            onClick={getRandomRecommendations}
          >
            🔄 Yenile
          </button>
        </div>

        {(filterType === 'category' || filterType === 'author') && (
          <div className="recommendations-filter-container">
            <div className="recommendations-info-box">
              <p>💡 <strong>Varsayılan:</strong> Rastgele {filterType === 'category' ? 'kategori' : 'yazar'} seçilir</p>
              <p>📝 <strong>İsteğe bağlı:</strong> Belirli bir {filterType === 'category' ? 'kategori' : 'yazar'} seçebilirsiniz</p>
            </div>
            
            <form onSubmit={handleFilterSubmit} className="recommendations-filter-form">
              <div className="recommendations-input-group">
                <input
                  type="text"
                  value={filterValue}
                  onChange={(e) => setFilterValue(e.target.value)}
                  placeholder={
                    filterType === 'category' 
                      ? 'Belirli kategori (ör: literature, academic)' 
                      : 'Belirli yazar (ör: N.G. Kabal, Stefan Zweig)'
                  }
                  className="recommendations-filter-input"
                />
                <button type="submit" className="recommendations-filter-submit">
                  🔍 Filtrele
                </button>
              </div>
            </form>

            {/* Hızlı seçim butonları */}
            {filterType === 'category' && (
              <div className="recommendations-quick-select">
                <p className="quick-select-label">Popüler kategoriler:</p>
                <div className="quick-buttons">
                  <button onClick={() => handleQuickCategorySelect('literature')} className="quick-btn">📚 Literature</button>
                  <button onClick={() => handleQuickCategorySelect('academic')} className="quick-btn">🎓 Academic</button>
                  <button onClick={() => handleQuickCategorySelect('psychology')} className="quick-btn">🧠 Psychology</button>
                  <button onClick={() => handleQuickCategorySelect('science_and_engineering')} className="quick-btn">🔬 Science & Engineering</button>
                  <button onClick={() => handleQuickCategorySelect('business_and_economy')} className="quick-btn">💼 Business & Economy</button>
                </div>
              </div>
            )}

            {filterType === 'author' && (
              <div className="recommendations-quick-select">
                <p className="quick-select-label">Popüler yazarlar:</p>
                <div className="quick-buttons">
                  <button onClick={() => handleQuickAuthorSelect('N.G. Kabal')} className="quick-btn">✍️ N.G. Kabal</button>
                  <button onClick={() => handleQuickAuthorSelect('Ömer Seyfettin')} className="quick-btn">✍️ Ömer Seyfettin</button>
                  <button onClick={() => handleQuickAuthorSelect('Stefan Zweig')} className="quick-btn">✍️ Stefan Zweig</button>
                  <button onClick={() => handleQuickAuthorSelect('Franz Kafka')} className="quick-btn">✍️ Franz Kafka</button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      <div className="recommendations-content">
        {recommendations.length === 0 ? (
          <div className="recommendations-empty">
            <div className="recommendations-empty-icon">📭</div>
            <h3 className="recommendations-empty-title">Öneri bulunamadı</h3>
            <p className="recommendations-empty-subtitle">
              {filterType === 'general' 
                ? 'Henüz genel öneriler mevcut değil.'
                : `${selectedInfo} için öneri bulunamadı. Yeni seçim yapmayı deneyin.`
              }
            </p>
            <button onClick={getRandomRecommendations} className="recommendations-retry-btn">
              🎲 Tekrar Dene
            </button>
          </div>
        ) : (
          <>
           

            <div className="recommendations-grid">
              {recommendations.map((recommendation, index) => (
                <div key={`${recommendation.book.id}-${index}`} className="recommendation-item">
                  <div className="recommendation-badges">
                    <span className="recommendation-score-badge">
                      ⭐ {recommendation.score}/100
                    </span>
                    <span className="recommendation-reason-badge">
                      {recommendation.reason}
                    </span>
                  </div>
                  
                  <BookCard 
                    book={recommendation.book}
                    showAuthor={true}
                    showCategory={true}
                  />
                </div>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default RecommendationsPage; 