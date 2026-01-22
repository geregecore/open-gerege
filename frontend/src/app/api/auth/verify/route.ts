import { NextRequest, NextResponse } from 'next/server'
import { cookies } from 'next/headers'

const BACKEND_BASE = process.env.API_PROXY_TARGET || 'http://localhost:8080'

export async function GET(request: NextRequest) {
    const sid = request.nextUrl.searchParams.get('sid')
    const code = request.nextUrl.searchParams.get('code') // Support code if needed

    const token = sid || code;

    if (!token) {
        return NextResponse.redirect(new URL('/?error=no_sid', request.url))
    }

    try {
        // Call backend auth/verify
        // The old code used /auth/verify?sid=...
        const response = await fetch(`${BACKEND_BASE}/auth/verify?sid=${token}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'User-Agent': request.headers.get('user-agent') || '',
            },
            redirect: 'manual',
        })

        console.log('Backend verify status:', response.status);

        // If backend redirects (302), it implies success + cookies
        // Or if it returns 200 with cookies.
        // We should forward cookies.

        const setCookieHeaders = response.headers.getSetCookie()

        // Redirect to Profile
        const redirectUrl = new URL('/profile', request.url)
        const redirectResponse = NextResponse.redirect(redirectUrl)

        // Forward cookies
        for (const cookieHeader of setCookieHeaders) {
            redirectResponse.headers.append('Set-Cookie', cookieHeader)
        }

        // Also set via cookieStore for immediate effect in Middleware or Server Components
        const cookieStore = await cookies()
        for (const cookieHeader of setCookieHeaders) {
            const [nameValue, ...rest] = cookieHeader.split(';')
            const [name, value] = nameValue.split('=')
            if (name && value) {
                const attrs = rest.join(';').toLowerCase()
                const isHttpOnly = attrs.includes('httponly')
                const isSecure = attrs.includes('secure')
                const isCsrf = name.toLowerCase().includes('csrf') || name.toLowerCase().includes('xsrf')

                cookieStore.set(name.trim(), value.trim(), {
                    httpOnly: isCsrf ? false : isHttpOnly,
                    secure: isSecure || process.env.NODE_ENV === 'production',
                    sameSite: 'lax',
                    path: '/',
                })
            }
        }

        return redirectResponse
    } catch (error) {
        console.error('Auth verify error:', error)
        return NextResponse.redirect(new URL('/?error=verify_failed_exception', request.url))
    }
}
