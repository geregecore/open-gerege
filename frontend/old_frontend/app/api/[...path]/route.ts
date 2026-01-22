import { NextRequest, NextResponse } from 'next/server'

/**
 * 🔒 API Proxy Route
 * 
 * Бүх API хүсэлтийг backend руу proxy хийж, base URL-г нууна.
 * Network tab дээр зөвхөн `/api/*` харагдана.
 */

// API base URL - all endpoints at root
const getApiBase = () => {
  const target = process.env.API_PROXY_TARGET || process.env.NEXT_PUBLIC_API_BASE
  if (target) {
    if (target.startsWith('http://') || target.startsWith('https://')) {
      // Ensure /v1 suffix if missing
      if (target.endsWith('/api')) {
        return `${target}/v1`
      }
      return target
    }
  }

  // Local development fallback
  if (process.env.NODE_ENV === 'development') {
    return 'http://localhost:8080/api/v1'
  }

  return 'https://template.gerege.mn/api/v1'
}

const API_BASE = getApiBase()

export async function GET(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const params = await context.params
  return proxyRequest(request, params.path, 'GET')
}

export async function POST(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const params = await context.params
  return proxyRequest(request, params.path, 'POST')
}

export async function PUT(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const params = await context.params
  return proxyRequest(request, params.path, 'PUT')
}

export async function PATCH(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const params = await context.params
  return proxyRequest(request, params.path, 'PATCH')
}

export async function DELETE(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const params = await context.params
  return proxyRequest(request, params.path, 'DELETE')
}

async function proxyRequest(
  request: NextRequest,
  pathSegments: string[],
  method: string
) {
  let targetURL = ''
  try {
    const path = pathSegments.join('/')
    const queryString = request.nextUrl.search

    // Build target URL
    const base = API_BASE.replace(/\/+$/, '')
    const targetPath = path.startsWith('/') ? path : `/${path}`
    targetURL = `${base}${targetPath}${queryString}`

    // Ensure targetURL is absolute
    if (!targetURL.startsWith('http://') && !targetURL.startsWith('https://')) {
      // This shouldn't happen, but if it does, construct absolute URL
      const origin = request.headers.get('host') || 'localhost:3000'
      const protocol = request.headers.get('x-forwarded-proto') || 'http'
      targetURL = `${protocol}://${origin}${targetURL}`
    }

    // Get request body
    let body: BodyInit | undefined
    const contentType = request.headers.get('content-type')
    if (method !== 'GET' && method !== 'DELETE') {
      if (contentType?.includes('application/json')) {
        body = await request.text()
      } else if (contentType?.includes('multipart/form-data')) {
        body = await request.formData()
      } else {
        body = await request.text()
      }
    }

    // Forward headers (exclude host, connection, and origin-related headers)
    const headers = new Headers()
    const skipHeaders = ['host', 'connection', 'content-length', 'origin', 'referer']
    request.headers.forEach((value, key) => {
      if (!skipHeaders.includes(key.toLowerCase())) {
        headers.set(key, value)
      }
    })

    // Set Origin and Referer to backend URL (important for CSRF validation)
    headers.set('Origin', API_BASE)
    headers.set('Referer', API_BASE)

    // Extract CSRF token from cookies and add as header
    const cookieHeader = request.headers.get('cookie') || ''
    const csrfMatch = cookieHeader.match(/(?:csrf_token|_csrf|csrftoken|XSRF-TOKEN)=([^;]+)/)
    if (csrfMatch) {
      headers.set('X-CSRF-Token', csrfMatch[1])
      headers.set('X-XSRF-Token', csrfMatch[1])
    }

    // Log request details for debugging
    if (path.includes('logout')) {
      console.log('=== LOGOUT REQUEST DEBUG ===')
      console.log('Target URL:', targetURL)
      console.log('Method:', method)
      console.log('Request headers:', Object.fromEntries(headers.entries()))
      console.log('Request body:', body)
      console.log('Cookie header:', cookieHeader)
    }

    // Make request to backend
    const response = await fetch(targetURL, {
      method,
      headers,
      body,
      credentials: 'include',
    })

    // Log response for logout
    if (path.includes('logout')) {
      console.log('Response status:', response.status)
      console.log('Response headers:', Object.fromEntries(response.headers.entries()))
    }

    // Get response body
    const responseBody = await response.text()

    // Create response with same status and headers
    const proxyResponse = new NextResponse(responseBody, {
      status: response.status,
      statusText: response.statusText,
    })

    // Copy response headers
    response.headers.forEach((value, key) => {
      if (!['content-encoding', 'transfer-encoding'].includes(key.toLowerCase())) {
        proxyResponse.headers.set(key, value)
      }
    })

    return proxyResponse
  } catch (error) {
    console.error('Proxy error:', error)
    console.error('API_BASE:', API_BASE)
    console.error('Path segments:', pathSegments)
    console.error('Target URL:', targetURL)
    return NextResponse.json(
      {
        error: 'Proxy request failed',
        message: error instanceof Error ? error.message : 'Unknown error',
        details: {
          apiBase: API_BASE,
          path: pathSegments.join('/'),
          targetURL: targetURL
        }
      },
      { status: 500 }
    )
  }
}

