import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Routes that require authentication
const protectedRoutes = ["/profile", "/dashboard", "/settings"];

// Routes that are only for non-authenticated users
const authRoutes = ["/login", "/register"];

export function middleware(request: NextRequest) {
    const { pathname } = request.nextUrl;

    // Get session token from cookies
    const sessionToken = request.cookies.get("session")?.value || request.cookies.get("token")?.value;
    const isAuthenticated = !!sessionToken;

    // Check if current path is a protected route
    const isProtectedRoute = protectedRoutes.some(
        (route) => pathname === route || pathname.startsWith(`${route}/`)
    );

    // Check if current path is an auth route (login/register)
    const isAuthRoute = authRoutes.some(
        (route) => pathname === route || pathname.startsWith(`${route}/`)
    );

    // If trying to access protected route without authentication
    if (isProtectedRoute && !isAuthenticated) {
        const loginUrl = new URL("/login", request.url);
        // Add redirect parameter so user can be redirected back after login
        loginUrl.searchParams.set("redirect", pathname);
        return NextResponse.redirect(loginUrl);
    }

    // If authenticated user tries to access auth routes (login/register)
    if (isAuthRoute && isAuthenticated) {
        // Redirect to profile or the redirect parameter
        const redirectUrl = request.nextUrl.searchParams.get("redirect");
        const destination = redirectUrl || "/profile";
        return NextResponse.redirect(new URL(destination, request.url));
    }

    return NextResponse.next();
}

// Configure which routes the middleware runs on
export const config = {
    matcher: [
        /*
         * Match all request paths except for the ones starting with:
         * - api (API routes)
         * - _next/static (static files)
         * - _next/image (image optimization files)
         * - favicon.ico (favicon file)
         * - public assets
         */
        "/((?!api|_next/static|_next/image|favicon.ico|.*\\..*|_next).*)",
    ],
};
