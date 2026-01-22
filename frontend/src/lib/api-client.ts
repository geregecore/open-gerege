import Cookies from 'js-cookie';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';

type RequestConfig = RequestInit & {
    version?: 'v1' | 'v2';
};

async function fetchWrapper<T>(endpoint: string, config: RequestConfig = {}): Promise<T> {
    const { version = 'v1', headers, ...customConfig } = config;

    const path = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    const url = `${API_BASE_URL}${path}`;

    const defaultHeaders: HeadersInit = {
        'Content-Type': 'application/json',
    };

    // Add Authorization header if token exists
    // We check localStorage or Cookie. 
    // Should ideally check if we are on client.
    let token: string | undefined;
    if (typeof window !== 'undefined') {
        token = Cookies.get('token') || localStorage.getItem('access_token') || undefined;
    }

    const authHeaders: HeadersInit = token ? { 'Authorization': `Bearer ${token}` } : {};

    const response = await fetch(url, {
        headers: {
            ...defaultHeaders,
            ...authHeaders,
            ...headers,
        },
        credentials: 'include',
        ...customConfig,
    });

    const data = await response.json();

    if (!response.ok) {
        // If 401, clear tokens and redirect to login
        if (response.status === 401 && typeof window !== 'undefined') {
            Cookies.remove('token');
            localStorage.removeItem('access_token');
            // Optionally redirect to login
            // window.location.href = '/login';
        }
        throw new Error(data.message || 'Something went wrong');
    }

    if (data.code && data.data) {
        return data.data;
    }

    return data;
}

export const apiClient = {
    get: <T>(endpoint: string, config?: RequestConfig) =>
        fetchWrapper<T>(endpoint, { ...config, method: 'GET' }),

    post: <T>(endpoint: string, body: any, config?: RequestConfig) =>
        fetchWrapper<T>(endpoint, { ...config, method: 'POST', body: JSON.stringify(body) }),

    put: <T>(endpoint: string, body: any, config?: RequestConfig) =>
        fetchWrapper<T>(endpoint, { ...config, method: 'PUT', body: JSON.stringify(body) }),

    delete: <T>(endpoint: string, config?: RequestConfig) =>
        fetchWrapper<T>(endpoint, { ...config, method: 'DELETE' }),
};
