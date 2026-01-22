import { NextRequest, NextResponse } from 'next/server'
import { cookies } from 'next/headers'

const BACKEND_BASE = process.env.API_PROXY_TARGET || 'http://localhost:8080'

/**
 * Auth Verify Handler
 * Proxies to backend /auth/verify and forwards cookies to frontend domain
 */
export async function GET(request: NextRequest) {
  const sid = request.nextUrl.searchParams.get('sid')

  if (!sid) {
    return NextResponse.redirect(new URL('/mn/home?error=no_sid', request.url))
  }

  try {
    // Call backend auth/verify
    const response = await fetch(`${BACKEND_BASE}/auth/verify?sid=${sid}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      redirect: 'manual', // Don't follow redirects automatically
    })

    // Get cookies from backend response
    const setCookieHeaders = response.headers.getSetCookie()

    // Log full response for debugging
    console.log('Backend auth/verify status:', response.status)
    console.log('Backend auth/verify headers:', Object.fromEntries(response.headers.entries()))
    console.log('Backend auth/verify cookies:', setCookieHeaders)

    // Try to get response body for debugging
    const responseText = await response.text()
    console.log('Backend auth/verify body:', responseText.substring(0, 500))

    // Create redirect response to /callback
    const redirectUrl = new URL('/mn/callback', request.url)
    const redirectResponse = NextResponse.redirect(redirectUrl)

    // Forward all Set-Cookie headers from backend
    for (const cookieHeader of setCookieHeaders) {
      redirectResponse.headers.append('Set-Cookie', cookieHeader)
    }

    // If backend returned cookies, also try to set them via cookies() API
    const cookieStore = await cookies()
    for (const cookieHeader of setCookieHeaders) {
      // Parse cookie name and value
      const [nameValue, ...rest] = cookieHeader.split(';')
      const [name, value] = nameValue.split('=')
      if (name && value) {
        // Parse cookie attributes from backend
        const attrs = rest.join(';').toLowerCase()
        const isHttpOnly = attrs.includes('httponly')
        const isSecure = attrs.includes('secure')

        // CSRF token should NOT be httpOnly (needs to be readable by JavaScript)
        const isCsrfCookie = name.toLowerCase().includes('csrf') || name.toLowerCase().includes('xsrf')

        cookieStore.set(name.trim(), value.trim(), {
          httpOnly: isCsrfCookie ? false : isHttpOnly,
          secure: isSecure || process.env.NODE_ENV === 'production',
          sameSite: 'lax',
          path: '/',
        })
      }
    }

    return redirectResponse
  } catch (error) {
    console.error('Auth verify error:', error)
    return NextResponse.redirect(new URL('/mn/home?error=verify_failed', request.url))
  }
}
