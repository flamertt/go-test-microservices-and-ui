import { useState, useEffect } from 'react';
import type { PaginatedResponse } from '../types';

interface UsePaginationState<T> {
  data: T[] | null;
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  loading: boolean;
  error: string | null;
}

interface UsePaginationParams {
  endpoint: string;
  initialPage?: number;
  initialPageSize?: number;
  searchTerm?: string;
  category?: string;
  author?: string;
}

export const usePagination = <T>({
  endpoint,
  initialPage = 1,
  initialPageSize = 50,
  searchTerm = '',
  category = '',
  author = ''
}: UsePaginationParams): UsePaginationState<T> & {
  setPage: (page: number) => void;
  nextPage: () => void;
  prevPage: () => void;
  refetch: () => void;
} => {
  const [state, setState] = useState<UsePaginationState<T>>({
    data: null,
    total: 0,
    page: initialPage,
    pageSize: initialPageSize,
    totalPages: 0,
    loading: true,
    error: null,
  });

  const fetchData = async (page: number) => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const params = new URLSearchParams({
        page: page.toString(),
        page_size: initialPageSize.toString(),
      });
      
      if (searchTerm.trim()) {
        params.append('search', searchTerm.trim());
      }
      
      if (category.trim()) {
        params.append('category', category.trim());
      }

      if (author.trim()) {
        params.append('author', author.trim());
      }

      const url = `${endpoint}?${params.toString()}`;
      const response = await fetch(url);
      
      if (!response.ok) {
        throw new Error(`API Hatası: ${response.status} ${response.statusText}`);
      }
      
      const result = await response.json();
      const paginatedData: PaginatedResponse<T> = result.data;
      
      // Data field'ını belirle (books, authors, categories, genres)
      let dataArray: T[] = [];
      if (paginatedData.books) {
        dataArray = paginatedData.books;
      } else if (paginatedData.authors) {
        dataArray = paginatedData.authors;
      } else if (paginatedData.categories) {
        dataArray = paginatedData.categories;
      } else if ((paginatedData as any).genres) {
        dataArray = (paginatedData as any).genres;
      }
      
      setState({
        data: dataArray,
        total: paginatedData.total,
        page: paginatedData.page,
        pageSize: paginatedData.page_size,
        totalPages: paginatedData.total_pages,
        loading: false,
        error: null,
      });
    } catch (error: any) {
      setState(prev => ({
        ...prev,
        data: null,
        loading: false,
        error: error.message || 'Bilinmeyen bir hata oluştu',
      }));
    }
  };

  useEffect(() => {
    fetchData(initialPage);
  }, [endpoint, searchTerm, category, author]);

  const setPage = (page: number) => {
    if (page >= 1 && page <= state.totalPages) {
      fetchData(page);
    }
  };

  const nextPage = () => {
    if (state.page < state.totalPages) {
      fetchData(state.page + 1);
    }
  };

  const prevPage = () => {
    if (state.page > 1) {
      fetchData(state.page - 1);
    }
  };

  const refetch = () => {
    fetchData(state.page);
  };

  return {
    ...state,
    setPage,
    nextPage,
    prevPage,
    refetch,
  };
}; 