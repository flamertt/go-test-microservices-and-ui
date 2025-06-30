import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { 
  FaUserPlus, 
  FaUser, 
  FaEnvelope, 
  FaLock, 
  FaEye, 
  FaEyeSlash,
  FaExclamationTriangle,
  FaSpinner,
  FaCheck
} from 'react-icons/fa';
import '../styles/AuthPages.css';

const RegisterPage: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const { register, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  // Eğer kullanıcı zaten giriş yapmışsa ana sayfaya yönlendir
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/', { replace: true });
    }
  }, [isAuthenticated, navigate]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
    // Hata mesajını temizle
    if (error) setError('');
  };

  const validateForm = (): boolean => {
    const { username, email, password, confirmPassword } = formData;

    if (!username.trim() || !email.trim() || !password.trim() || !confirmPassword.trim()) {
      setError('Lütfen tüm alanları doldurun');
      return false;
    }

    if (username.length < 3) {
      setError('Kullanıcı adı en az 3 karakter olmalıdır');
      return false;
    }

    if (username.length > 50) {
      setError('Kullanıcı adı en fazla 50 karakter olmalıdır');
      return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setError('Geçerli bir e-posta adresi girin');
      return false;
    }

    if (password.length < 6) {
      setError('Şifre en az 6 karakter olmalıdır');
      return false;
    }

    if (password !== confirmPassword) {
      setError('Şifreler eşleşmiyor');
      return false;
    }

    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      await register({
        username: formData.username.trim(),
        email: formData.email.trim(),
        password: formData.password,
      });
      // Başarılı kayıt, useEffect tarafından yönlendirilecek
    } catch (err: any) {
      setError(err.message || 'Kayıt olunurken bir hata oluştu');
    } finally {
      setIsLoading(false);
    }
  };

  const getPasswordStrength = (password: string) => {
    if (password.length === 0) return { strength: 0, label: '' };
    if (password.length < 4) return { strength: 1, label: 'Zayıf' };
    if (password.length < 8) return { strength: 2, label: 'Orta' };
    return { strength: 3, label: 'Güçlü' };
  };

  const passwordStrength = getPasswordStrength(formData.password);

  return (
    <div className="auth-container">
      <div className="auth-background-decoration">
        <div className="floating-shape shape-1"></div>
        <div className="floating-shape shape-2"></div>
        <div className="floating-shape shape-3"></div>
      </div>
      
      <div className="auth-card">
        <div className="auth-header">
          <div className="auth-icon">
            <FaUserPlus />
          </div>
          <h1>Hesap Oluşturun</h1>
          <p>Kütüphane hesabı oluşturun ve keşfetmeye başlayın</p>
        </div>

        <form onSubmit={handleSubmit} className="auth-form">
          {error && (
            <div className="error-message">
              <FaExclamationTriangle className="error-icon" />
              {error}
            </div>
          )}

          <div className="form-group">
            <label htmlFor="username">
              <FaUser className="label-icon" />
              Kullanıcı Adı
            </label>
            <div className="input-wrapper">
              <FaUser className="input-icon" />
              <input
                type="text"
                id="username"
                name="username"
                value={formData.username}
                onChange={handleChange}
                placeholder="Kullanıcı adınızı girin (3-50 karakter)"
                disabled={isLoading}
                autoComplete="username"
              />
              {formData.username.length >= 3 && (
                <FaCheck className="input-success-icon" />
              )}
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="email">
              <FaEnvelope className="label-icon" />
              E-posta Adresi
            </label>
            <div className="input-wrapper">
              <FaEnvelope className="input-icon" />
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                placeholder="E-posta adresinizi girin"
                disabled={isLoading}
                autoComplete="email"
              />
              {formData.email.includes('@') && formData.email.includes('.') && (
                <FaCheck className="input-success-icon" />
              )}
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="password">
              <FaLock className="label-icon" />
              Şifre
            </label>
            <div className="input-wrapper">
              <FaLock className="input-icon" />
              <input
                type={showPassword ? 'text' : 'password'}
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="Şifrenizi girin (en az 6 karakter)"
                disabled={isLoading}
                autoComplete="new-password"
              />
              <button
                type="button"
                className="password-toggle"
                onClick={() => setShowPassword(!showPassword)}
                disabled={isLoading}
              >
                {showPassword ? <FaEyeSlash /> : <FaEye />}
              </button>
            </div>
            {formData.password && (
              <div className="password-strength">
                <div className={`strength-bar strength-${passwordStrength.strength}`}>
                  <div className="strength-fill"></div>
                </div>
                <span className="strength-label">{passwordStrength.label}</span>
              </div>
            )}
          </div>

          <div className="form-group">
            <label htmlFor="confirmPassword">
              <FaLock className="label-icon" />
              Şifre Tekrar
            </label>
            <div className="input-wrapper">
              <FaLock className="input-icon" />
              <input
                type={showConfirmPassword ? 'text' : 'password'}
                id="confirmPassword"
                name="confirmPassword"
                value={formData.confirmPassword}
                onChange={handleChange}
                placeholder="Şifrenizi tekrar girin"
                disabled={isLoading}
                autoComplete="new-password"
              />
              <button
                type="button"
                className="password-toggle"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                disabled={isLoading}
              >
                {showConfirmPassword ? <FaEyeSlash /> : <FaEye />}
              </button>
              {formData.confirmPassword && formData.password === formData.confirmPassword && (
                <FaCheck className="input-success-icon" />
              )}
            </div>
          </div>

          <button
            type="submit"
            className={`auth-button ${isLoading ? 'loading' : ''}`}
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <FaSpinner className="spinner-icon" />
                Kayıt oluşturuluyor...
              </>
            ) : (
              <>
                <FaUserPlus className="button-icon" />
                Kayıt Ol
              </>
            )}
          </button>
        </form>

        <div className="auth-footer">
          <p>
            Zaten hesabınız var mı?{' '}
            <Link to="/login" className="auth-link">
              Giriş yap
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage; 