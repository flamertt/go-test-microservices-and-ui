import { useState, useEffect } from 'react';

interface UseApiState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

export const useApi = <T>(
  endpoint: string,
  dependencies: any[] = []
): UseApiState<T> => {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    loading: true,
    error: null,
  });

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        setState(prev => ({ ...prev, loading: true, error: null }));
        const response = await fetch(endpoint);
        
        if (!response.ok) {
          throw new Error(`API Hatası: ${response.status} ${response.statusText}`);
        }
        
        const data = await response.json();
        
        if (isMounted) {
          setState({
            data: data.data, // API response'unda data property'si var
            loading: false,
            error: null,
          });
        }
      } catch (error: any) {
        if (isMounted) {
          setState({
            data: null,
            loading: false,
            error: error.message || 'Bilinmeyen bir hata oluştu',
          });
        }
      }
    };

    fetchData();

    return () => {
      isMounted = false;
    };
  }, [endpoint, ...dependencies]);

  return state;
}; 