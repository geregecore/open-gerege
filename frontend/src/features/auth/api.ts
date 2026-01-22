import { apiClient } from '@/lib/api-client';
import { LocalLoginRequest, LoginResponse } from './types/index';

const SSO_CONFIG = {
    origin: process.env.NEXT_PUBLIC_SSO_ORIGIN || 'https://sso.gerege.mn',
    clientId: process.env.NEXT_PUBLIC_SSO_CLIENT_ID || 'GRG-CLI-01KCGT4564YJ6WM15VNP3Y1BFG',
    redirectUri: process.env.NEXT_PUBLIC_REDIRECT_URI || 'http://localhost:3000',
};

export const getSSOLoginUrl = () => {
    const callbackUrl = `${SSO_CONFIG.redirectUri}/callback`;
    const encodedUri = encodeURIComponent(callbackUrl);
    return `${SSO_CONFIG.origin}/auth?client_id=${SSO_CONFIG.clientId}&redirect_uri=${encodedUri}`;
};

export const authApi = {
    loginLocal: (data: LocalLoginRequest) => apiClient.post<LoginResponse>('/auth/local/login', data),
};
