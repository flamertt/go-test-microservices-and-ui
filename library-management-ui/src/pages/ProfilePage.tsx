import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { authApi } from '../services/api';
import { 
  FaUser,
  FaEnvelope,
  FaCalendarAlt,
  FaClock,
  FaLock,
  FaKey,
  FaSignOutAlt,
  FaEdit,
  FaTimes,
  FaEye,
  FaEyeSlash,
  FaCheck,
  FaExclamationTriangle,
  FaUserShield,
  FaInfoCircle
} from 'react-icons/fa';
import '../styles/ProfilePage.css';

const ProfilePage: React.FC = () => {
  const { user, logout } = useAuth();
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [showOldPassword, setShowOldPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [passwordData, setPasswordData] = useState({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData(prev => ({
      ...prev,
      [name]: value,
    }));
    // Mesajları temizle
    if (error) setError('');
    if (message) setMessage('');
  };

  const validatePasswordForm = (): boolean => {
    const { oldPassword, newPassword, confirmPassword } = passwordData;

    if (!oldPassword.trim() || !newPassword.trim() || !confirmPassword.trim()) {
      setError('Lütfen tüm alanları doldurun');
      return false;
    }

    if (newPassword.length < 6) {
      setError('Yeni şifre en az 6 karakter olmalıdır');
      return false;
    }

    if (newPassword !== confirmPassword) {
      setError('Yeni şifreler eşleşmiyor');
      return false;
    }

    if (oldPassword === newPassword) {
      setError('Yeni şifre eski şifreden farklı olmalıdır');
      return false;
    }

    return true;
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validatePasswordForm()) {
      return;
    }

    setIsLoading(true);
    setError('');
    setMessage('');

    try {
      await authApi.changePassword(passwordData.oldPassword, passwordData.newPassword);
      setMessage('Şifreniz başarıyla değiştirildi');
      setPasswordData({
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
      });
      setIsChangingPassword(false);
    } catch (err: any) {
      setError(err.message || 'Şifre değiştirirken bir hata oluştu');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('tr-TR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (!user) {
    return (
      <div className="profile-container">
        <div className="profile-error">
          <h2>Profil bilgisi bulunamadı</h2>
          <p>Lütfen tekrar giriş yapın.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <div className="profile-avatar">
          <FaUser />
        </div>
        <div className="profile-info">
          <h1>Profil Bilgilerim</h1>
          <p>Hesap bilgilerinizi görüntüleyin ve yönetin</p>
        </div>
      </div>

      <div className="profile-content">
        <div className="profile-card">
          <div className="card-header">
            <h2>
              <FaInfoCircle />
              Hesap Bilgileri
            </h2>
          </div>
          <div className="card-content">
            <div className="info-group">
              <label>
                <FaUser className="info-icon" />
                Kullanıcı Adı
              </label>
              <div className="info-value">{user.username}</div>
            </div>
            <div className="info-group">
              <label>
                <FaEnvelope className="info-icon" />
                E-posta Adresi
              </label>
              <div className="info-value">{user.email}</div>
            </div>
            <div className="info-group">
              <label>
                <FaCalendarAlt className="info-icon" />
                Hesap Oluşturma Tarihi
              </label>
              <div className="info-value">{formatDate(user.created_at)}</div>
            </div>
            <div className="info-group">
              <label>
                <FaClock className="info-icon" />
                Son Güncelleme
              </label>
              <div className="info-value">{formatDate(user.updated_at)}</div>
            </div>
          </div>
        </div>

        <div className="profile-card">
          <div className="card-header">
            <h2>
              <FaUserShield />
              Güvenlik
            </h2>
          </div>
          <div className="card-content">
            {!isChangingPassword ? (
              <div className="security-actions">
                <button
                  className="secondary-button"
                  onClick={() => setIsChangingPassword(true)}
                >
                  <FaEdit />
                  Şifre Değiştir
                </button>
                <p className="security-hint">
                  Hesabınızın güvenliği için düzenli olarak şifrenizi değiştirin.
                </p>
              </div>
            ) : (
              <form onSubmit={handlePasswordSubmit} className="password-form">
                <div className="form-group">
                  <label htmlFor="oldPassword">
                    <FaLock />
                    Mevcut Şifre
                  </label>
                  <div className="password-input-wrapper">
                  <input
                      type={showOldPassword ? "text" : "password"}
                    id="oldPassword"
                    name="oldPassword"
                    value={passwordData.oldPassword}
                    onChange={handlePasswordChange}
                    placeholder="Mevcut şifrenizi girin"
                    disabled={isLoading}
                  />
                    <button
                      type="button"
                      className="password-toggle"
                      onClick={() => setShowOldPassword(!showOldPassword)}
                    >
                      {showOldPassword ? <FaEyeSlash /> : <FaEye />}
                    </button>
                  </div>
                </div>

                <div className="form-group">
                  <label htmlFor="newPassword">
                    <FaKey />
                    Yeni Şifre
                  </label>
                  <div className="password-input-wrapper">
                  <input
                      type={showNewPassword ? "text" : "password"}
                    id="newPassword"
                    name="newPassword"
                    value={passwordData.newPassword}
                    onChange={handlePasswordChange}
                    placeholder="Yeni şifrenizi girin (en az 6 karakter)"
                    disabled={isLoading}
                  />
                    <button
                      type="button"
                      className="password-toggle"
                      onClick={() => setShowNewPassword(!showNewPassword)}
                    >
                      {showNewPassword ? <FaEyeSlash /> : <FaEye />}
                    </button>
                  </div>
                </div>

                <div className="form-group">
                  <label htmlFor="confirmPassword">
                    <FaCheck />
                    Yeni Şifre Tekrar
                  </label>
                  <div className="password-input-wrapper">
                  <input
                      type={showConfirmPassword ? "text" : "password"}
                    id="confirmPassword"
                    name="confirmPassword"
                    value={passwordData.confirmPassword}
                    onChange={handlePasswordChange}
                    placeholder="Yeni şifrenizi tekrar girin"
                    disabled={isLoading}
                  />
                    <button
                      type="button"
                      className="password-toggle"
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    >
                      {showConfirmPassword ? <FaEyeSlash /> : <FaEye />}
                    </button>
                  </div>
                </div>

                <div className="form-actions">
                  <button
                    type="button"
                    className="secondary-button"
                    onClick={() => {
                      setIsChangingPassword(false);
                      setPasswordData({
                        oldPassword: '',
                        newPassword: '',
                        confirmPassword: '',
                      });
                      setError('');
                      setMessage('');
                      setShowOldPassword(false);
                      setShowNewPassword(false);
                      setShowConfirmPassword(false);
                    }}
                    disabled={isLoading}
                  >
                    <FaTimes />
                    İptal
                  </button>
                  <button
                    type="submit"
                    className={`primary-button ${isLoading ? 'loading' : ''}`}
                    disabled={isLoading}
                  >
                    <FaCheck />
                    {isLoading ? 'Değiştiriliyor...' : 'Şifre Değiştir'}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>

        {(message || error) && (
          <div className={`message-card ${error ? 'error' : 'success'}`}>
            <span className="message-icon">
              {error ? <FaExclamationTriangle /> : <FaCheck />}
            </span>
            <span className="message-text">
              {error || message}
            </span>
          </div>
        )}

        <div className="profile-card">
          <div className="card-header">
            <h2>
              <FaSignOutAlt />
              Oturum Yönetimi
            </h2>
          </div>
          <div className="card-content">
            <div className="logout-section">
              <button className="danger-button" onClick={handleLogout}>
                <FaSignOutAlt />
                Çıkış Yap
              </button>
              <p className="logout-hint">
                Hesabınızdan güvenli bir şekilde çıkış yapın.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage; 