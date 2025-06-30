// PostgreSQL veritabanı yapısına uygun Book interface'i
export interface Book {
  id: number;
  title: string;
  publisher: string;
  author: string;
  category_name: string;
  product_code: string;
  page_count: number;
  released_year: number;
}

// Author interface - sadece ID ve isim
export interface Author {
  id: number;
  name: string;
}

// Genre interface - basitleştirildi 
export interface Genre {
  id: number;
  name: string;
}

// API Response wrapper
export interface ApiResponse<T> {
  data: T;
}

// Paginated Response
export interface PaginatedResponse<T> {
  books?: T[];
  authors?: T[];
  categories?: T[];
  genres?: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// API Error response
export interface ApiError {
  error: string;
}

// Auth Types
export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterResponse {
  message: string;
  user: User;
}

export interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  register: (userData: RegisterRequest) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  loading: boolean;
}