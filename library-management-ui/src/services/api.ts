import type { Book, Author, Genre, ApiResponse, LoginRequest, RegisterRequest, LoginResponse, RegisterResponse, User } from '../types';

const API_BASE_URL = '/api';

// Generic fetch fonksiyonu
async function fetchApi<T>(endpoint: string): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`);
    
    if (!response.ok) {
      throw new Error(`API Hatası: ${response.status} ${response.statusText}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Fetch hatası:', error);
    throw error;
  }
}

// Auth için özel fetch fonksiyonu (artık Gateway üzerinden)
async function fetchAuthApi<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  try {
    const token = localStorage.getItem('auth_token');
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };
    
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers,
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Network error' }));
      throw new Error(errorData.message || errorData.error || `HTTP ${response.status}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Auth API hatası:', error);
    throw error;
  }
}

// Yeni tipler
export interface EnrichedBook extends Book {
  author_info?: {
    name: string;
    biography: string;
  };
}

export interface AuthorDetail {
  author: Author;
  books: Book[];
  book_count: number;
}

export interface GenreDetail {
  genre: Genre;
  books: Book[];
  book_count: number;
  page?: number;
  page_size?: number;
  total_pages?: number;
  total?: number;
}

export interface Recommendation {
  book: Book;
  reason: string;
  score: number;
}

export interface RecommendationResponse {
  recommendations: Recommendation[];
  total: number;
  generated_at: string;
}

// Books API
export const booksApi = {
  // Tüm kitapları getir
  getAll: (): Promise<ApiResponse<Book[]>> => 
    fetchApi<ApiResponse<Book[]>>('/books'),
  
  // 🆕 Zenginleştirilmiş kitap listesi (yazar bilgileri ile)
  getEnriched: (page = 1, pageSize = 10): Promise<ApiResponse<EnrichedBook[]>> => 
    fetchApi<ApiResponse<EnrichedBook[]>>(`/books/enriched?page=${page}&page_size=${pageSize}`),
  
  // ID'ye göre kitap getir (artık yazar bilgisi ile zenginleştirilmiş)
  getById: (id: number): Promise<ApiResponse<EnrichedBook>> => 
    fetchApi<ApiResponse<EnrichedBook>>(`/books/${id}`),
  
  // Yazara göre kitapları getir
  getByAuthor: (authorName: string): Promise<ApiResponse<Book[]>> => 
    fetchApi<ApiResponse<Book[]>>(`/books/author/${encodeURIComponent(authorName)}`),
  
  // Kategoriye göre kitapları getir
  getByCategory: (categoryName: string): Promise<ApiResponse<Book[]>> => 
    fetchApi<ApiResponse<Book[]>>(`/books/category/${encodeURIComponent(categoryName)}`),
};

// Authors API
export const authorsApi = {
  // Tüm yazarları getir
  getAll: (): Promise<ApiResponse<Author[]>> => 
    fetchApi<ApiResponse<Author[]>>('/authors'),
  
  // ID'ye göre yazar getir (artık kitapları ile birlikte)
  getById: (id: number): Promise<ApiResponse<AuthorDetail>> => 
    fetchApi<ApiResponse<AuthorDetail>>(`/authors/${id}`),
  
  // 🆕 Yazar detayı + kitapları
  getDetail: (name: string): Promise<ApiResponse<AuthorDetail>> => 
    fetchApi<ApiResponse<AuthorDetail>>(`/authors/detail/${encodeURIComponent(name)}`),
  
  // Yazar ara
  search: (name: string): Promise<ApiResponse<Author[]>> => 
    fetchApi<ApiResponse<Author[]>>(`/authors/search?name=${encodeURIComponent(name)}`),
};

// Genres API
export const genresApi = {
  // Tüm türleri getir
  getAll: (): Promise<ApiResponse<Genre[]>> => 
    fetchApi<ApiResponse<Genre[]>>('/genres'),
  
  // ID'ye göre tür getir (artık kitapları ile birlikte)
  getById: (id: number): Promise<ApiResponse<GenreDetail>> => 
    fetchApi<ApiResponse<GenreDetail>>(`/genres/${id}`),
  
  // 🆕 Tür detayı + kitapları
  getDetail: (name: string): Promise<ApiResponse<GenreDetail>> => 
    fetchApi<ApiResponse<GenreDetail>>(`/genres/detail/${encodeURIComponent(name)}`),
  
  // Tür ara
  search: (name: string): Promise<ApiResponse<Genre[]>> => 
    fetchApi<ApiResponse<Genre[]>>(`/genres/search?name=${encodeURIComponent(name)}`),
};

// 🆕 Recommendations API
export const recommendationsApi = {
  // Genel akıllı öneriler (random 15 kitap)
  getGeneral: (limit = 15): Promise<ApiResponse<RecommendationResponse>> => 
    fetchApi<ApiResponse<RecommendationResponse>>(`/recommendations?limit=${limit}`),
  
  // 🆕 Random kategori önerileri (rastgele kategori seçilir)
  getRandomCategory: (limit = 15): Promise<ApiResponse<{ category: string; recommendations: Recommendation[]; total: number; type: string }>> => 
    fetchApi(`/recommendations/by-category?limit=${limit}`),
  
  // 🆕 Random yazar önerileri (rastgele yazar seçilir)
  getRandomAuthor: (limit = 15): Promise<ApiResponse<{ author: string; recommendations: Recommendation[]; total: number; type: string }>> => 
    fetchApi(`/recommendations/by-author?limit=${limit}`),
  
  // Belirli kategoriye göre öneriler
  getByCategory: (category: string, limit = 5): Promise<ApiResponse<{ category: string; recommendations: Recommendation[]; total: number }>> => 
    fetchApi(`/recommendations/by-category?category=${encodeURIComponent(category)}&limit=${limit}`),
  
  // Belirli yazara göre öneriler
  getByAuthor: (author: string, limit = 5): Promise<ApiResponse<{ author: string; recommendations: Recommendation[]; total: number }>> => 
    fetchApi(`/recommendations/by-author?author=${encodeURIComponent(author)}&limit=${limit}`),
  
  // Trend öneriler
  getTrending: (limit = 15): Promise<ApiResponse<RecommendationResponse>> => 
    fetchApi<ApiResponse<RecommendationResponse>>(`/recommendations/trending?limit=${limit}`),
  
  // Öneri servisinin durumu
  getStatus: (): Promise<ApiResponse<{ recommendation_service: string; dependent_services: Record<string, string>; all_services_healthy: boolean }>> => 
    fetchApi<ApiResponse<any>>('/recommendations/status'),
};

// 🆕 Health Check API
export const healthApi = {
  // Tüm servislerin durumu
  getAll: (): Promise<ApiResponse<{ gateway: string; services: Record<string, string> }>> => 
    fetchApi<ApiResponse<any>>('/health'),
  
  // Öneri servisinin durumu
  getRecommendations: (): Promise<ApiResponse<any>> => 
    fetchApi<ApiResponse<any>>('/recommendations/status'),
};

// 🔐 Auth API
export const authApi = {
  // Kullanıcı kaydı
  register: (userData: RegisterRequest): Promise<RegisterResponse> =>
    fetchAuthApi<RegisterResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    }),
  
  // Kullanıcı girişi
  login: (credentials: LoginRequest): Promise<LoginResponse> =>
    fetchAuthApi<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    }),
  
  // Kullanıcı profili getir
  getProfile: (): Promise<User> =>
    fetchAuthApi<User>('/auth/profile'),
  
  // Token doğrula
  validateToken: (): Promise<{ valid: boolean; user_id: number; username: string; email: string }> =>
    fetchAuthApi('/auth/validate'),
  
  // Şifre değiştir
  changePassword: (oldPassword: string, newPassword: string): Promise<{ message: string }> =>
    fetchAuthApi('/auth/change-password', {
      method: 'POST',
      body: JSON.stringify({
        old_password: oldPassword,
        new_password: newPassword,
      }),
    }),
  
  // Token yenile
  refreshToken: (token: string): Promise<{ token: string }> =>
    fetchAuthApi('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ token }),
    }),
  
  // Kullanıcı bilgisi getir (ID ile)
  getUser: (id: number): Promise<User> =>
    fetchAuthApi<User>(`/auth/users/${id}`),
  
  // Auth servis health check
  healthCheck: (): Promise<{ status: string; service: string; message: string }> =>
    fetchAuthApi('/health'),
}; 