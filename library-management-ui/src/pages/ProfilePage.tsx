import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { authApi } from '../services/api';
import '../styles/ProfilePage.css';

const ProfilePage: React.FC = () => {
  const { user, logout } = useAuth();
  const [isChangingPassword, setIsChangingPassword] = useState(false);
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
        <div className="profile-avatar">👤</div>
        <div className="profile-info">
          <h1>Profil Bilgilerim</h1>
          <p>Hesap bilgilerinizi görüntüleyin ve yönetin</p>
        </div>
      </div>

      <div className="profile-content">
        <div className="profile-card">
          <div className="card-header">
            <h2>🔍 Hesap Bilgileri</h2>
          </div>
          <div className="card-content">
            <div className="info-group">
              <label>Kullanıcı Adı</label>
              <div className="info-value">{user.username}</div>
            </div>
            <div className="info-group">
              <label>E-posta Adresi</label>
              <div className="info-value">{user.email}</div>
            </div>
            <div className="info-group">
              <label>Hesap Oluşturma Tarihi</label>
              <div className="info-value">{formatDate(user.created_at)}</div>
            </div>
            <div className="info-group">
              <label>Son Güncelleme</label>
              <div className="info-value">{formatDate(user.updated_at)}</div>
            </div>
          </div>
        </div>

        <div className="profile-card">
          <div className="card-header">
            <h2>🔐 Güvenlik</h2>
          </div>
          <div className="card-content">
            {!isChangingPassword ? (
              <div className="security-actions">
                <button
                  className="secondary-button"
                  onClick={() => setIsChangingPassword(true)}
                >
                  Şifre Değiştir
                </button>
                <p className="security-hint">
                  Hesabınızın güvenliği için düzenli olarak şifrenizi değiştirin.
                </p>
              </div>
            ) : (
              <form onSubmit={handlePasswordSubmit} className="password-form">
                <div className="form-group">
                  <label htmlFor="oldPassword">Mevcut Şifre</label>
                  <input
                    type="password"
                    id="oldPassword"
                    name="oldPassword"
                    value={passwordData.oldPassword}
                    onChange={handlePasswordChange}
                    placeholder="Mevcut şifrenizi girin"
                    disabled={isLoading}
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="newPassword">Yeni Şifre</label>
                  <input
                    type="password"
                    id="newPassword"
                    name="newPassword"
                    value={passwordData.newPassword}
                    onChange={handlePasswordChange}
                    placeholder="Yeni şifrenizi girin (en az 6 karakter)"
                    disabled={isLoading}
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="confirmPassword">Yeni Şifre Tekrar</label>
                  <input
                    type="password"
                    id="confirmPassword"
                    name="confirmPassword"
                    value={passwordData.confirmPassword}
                    onChange={handlePasswordChange}
                    placeholder="Yeni şifrenizi tekrar girin"
                    disabled={isLoading}
                  />
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
                    }}
                    disabled={isLoading}
                  >
                    İptal
                  </button>
                  <button
                    type="submit"
                    className={`primary-button ${isLoading ? 'loading' : ''}`}
                    disabled={isLoading}
                  >
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
              {error ? '⚠️' : '✅'}
            </span>
            <span className="message-text">
              {error || message}
            </span>
          </div>
        )}

        <div className="profile-card">
          <div className="card-header">
            <h2>🚪 Oturum Yönetimi</h2>
          </div>
          <div className="card-content">
            <div className="logout-section">
              <button className="danger-button" onClick={handleLogout}>
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