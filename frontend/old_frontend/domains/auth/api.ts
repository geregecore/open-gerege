/**
 * Auth Domain API
 */

import api from '@/lib/api'
import type {
  LoginRequest,
  LoginResponse,
  LogoutRequest,
  RefreshTokenResponse,
} from './types'

const AUTH_BASE = '/auth'

export async function login(data: LoginRequest): Promise<LoginResponse> {
  return api.post<LoginResponse>(`${AUTH_BASE}/login`, { ...data })
}

export async function logout(data?: LogoutRequest): Promise<void> {
  return api.post<void>(`${AUTH_BASE}/logout`, data ? { ...data } : undefined)
}

export async function refreshToken(): Promise<RefreshTokenResponse> {
  return api.post<RefreshTokenResponse>(`${AUTH_BASE}/refresh`)
}

export async function checkSession(): Promise<boolean> {
  try {
    await api.get(`${AUTH_BASE}/session`, { hasToast: false })
    return true
  } catch {
    return false
  }
}

// ==========================================
// SSO Authentication
// ==========================================

/**
 * SSO Configuration
 *
 * Environment variables:
 * - NEXT_PUBLIC_SSO_ORIGIN: SSO server URL (default: https://sso.gerege.mn)
 * - NEXT_PUBLIC_SSO_CLIENT_ID: SSO client ID
 * - NEXT_PUBLIC_SSO_REDIRECT_URI: Backend URL where /auth/verify is handled
 */
export const SSO_CONFIG = {
  origin: process.env.NEXT_PUBLIC_SSO_ORIGIN || 'https://sso.gerege.mn',
  clientId: process.env.NEXT_PUBLIC_SSO_CLIENT_ID || 'GRG-CLI-01KCGT4564YJ6WM15VNP3Y1BFG',
  redirectUri: process.env.NEXT_PUBLIC_SSO_REDIRECT_URI || 'http://localhost:3000',
} as const

/**
 * SSO Embed URL for iframe popup
 */
export function getSSOEmbedUrl(): string {
  return `${SSO_CONFIG.origin}/embed/auth?client_id=${SSO_CONFIG.clientId}&mode=embed`
}

/**
 * SSO Redirect URL for direct redirect
 */
export function getSSORedirectUrl(): string {
  const callbackUrl = `${SSO_CONFIG.redirectUri}/callback`
  const encodedUri = encodeURIComponent(callbackUrl)
  return `${SSO_CONFIG.origin}/auth?client_id=${SSO_CONFIG.clientId}&redirect_uri=${encodedUri}`
}

/**
 * SSO Message Types
 */
export type SSOMessageType = 'SSO_AUTH_SUCCESS' | 'SSO_AUTH_CANCEL' | 'SSO_AUTH_ERROR'

export interface SSOMessage {
  type: SSOMessageType
  token?: string
  sid?: string
  error?: string
}

/**
 * Validate SSO Message origin
 */
export function isValidSSOMessage(event: MessageEvent): boolean {
  return event.origin === SSO_CONFIG.origin
}

/**
 * Parse SSO Message
 */
export function parseSSOMessage(event: MessageEvent): SSOMessage | null {
  if (!isValidSSOMessage(event)) return null
  const { type, token, sid, error } = event.data || {}
  if (!type) return null
  return { type, token, sid, error }
}
