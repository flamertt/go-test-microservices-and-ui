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