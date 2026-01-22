"use client";

import { useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { apiClient } from "@/lib/api-client";

function CallbackContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const sid = searchParams.get('sid');
    const code = searchParams.get('code');
    const error = searchParams.get('error');

    useEffect(() => {
        if (error) {
            console.error("Auth error:", error);
            router.push('/?error=' + error);
            return;
        }

        const verifySession = async () => {
            const token = sid || code;
            if (!token) {
                // If we are here without token, but backend redirected us here, 
                // maybe cookie is already set?
                // Let's try to fetch profile directly.
                try {
                    // We'll just redirect to profile and let it handle the check
                    router.push('/profile');
                } catch (e) {
                    router.push('/');
                }
                return;
            }

            try {
                // New Flow: 
                // 1. We have SID. 
                // 2. We call /auth/verify?sid=... to explicitly set the cookie on our domain
                //    proxied via Next.js to backend.
                //    The backend /auth/verify returns 302 to callback URL.
                //    Wait, calling it via fetch might be wrong if it returns 302.
                //    But we are already ON the callback URL.

                // If the backend /auth/login -> SSO -> /auth/callback flow happened,
                // The backend /auth/callback ALREADY set the cookie (if on same domain).
                // But we are on localhost:3000 vs localhost:8080.
                // Cookies on localhost:8080 (backend) are NOT visible to localhost:3000 (frontend).

                // So we MUST use the sid to set cookie on localhost:3000.
                // We do this by hitting our Next.js API route /api/auth/verify?sid=...
                // which Proxies to Backend /auth/verify?sid=...
                // Backend Verify returns Set-Cookie.
                // Next.js Route forwards Set-Cookie.

                window.location.href = `/api/auth/verify?sid=${token}`;
            } catch (error) {
                console.error("Verification failed", error);
                router.push('/?error=verification_failed');
            }
        };

        verifySession();
    }, [sid, code, error, router]);

    return (
        <div className="flex min-h-screen items-center justify-center flex-col gap-4">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
            <p>Finalizing authentication...</p>
        </div>
    );
}

export default function CallbackPage() {
    return (
        <Suspense>
            <CallbackContent />
        </Suspense>
    );
}
