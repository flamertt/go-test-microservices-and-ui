import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authApi } from '../services/api';
import type { User, LoginRequest, RegisterRequest, AuthContextType } from '../types';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('auth_token'));
  const [loading, setLoading] = useState(true);

  // İlk yüklemede token kontrolü
  useEffect(() => {
    const initializeAuth = async () => {
      const savedToken = localStorage.getItem('auth_token');
      
      if (savedToken) {
        try {
          // Token'ı doğrula ve kullanıcı bilgilerini al
          const profileData = await authApi.getProfile();
          setUser(profileData);
          setToken(savedToken);
        } catch (error) {
          console.error('Token doğrulama hatası:', error);
          // Token geçersizse localStorage'dan kaldır
          localStorage.removeItem('auth_token');
          setToken(null);
          setUser(null);
        }
      }
      
      setLoading(false);
    };

    initializeAuth();
  }, []);

  const login = async (credentials: LoginRequest): Promise<void> => {
    try {
      setLoading(true);
      const response = await authApi.login(credentials);
      
      // Token ve kullanıcı bilgilerini kaydet
      localStorage.setItem('auth_token', response.token);
      setToken(response.token);
      setUser(response.user);
    } catch (error) {
      console.error('Giriş hatası:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const register = async (userData: RegisterRequest): Promise<void> => {
    try {
      setLoading(true);
      await authApi.register(userData);
      
      // Kayıt başarılıysa otomatik giriş yap
      await login({
        username: userData.username,
        password: userData.password,
      });
    } catch (error) {
      console.error('Kayıt hatası:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const logout = (): void => {
    localStorage.removeItem('auth_token');
    setToken(null);
    setUser(null);
  };

  const value: AuthContextType = {
    user,
    token,
    login,
    register,
    logout,
    isAuthenticated: !!user && !!token,
    loading,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}; 