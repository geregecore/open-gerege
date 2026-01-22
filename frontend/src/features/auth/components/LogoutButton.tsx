"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { authApi } from "../api";
import Cookies from "js-cookie";

interface LogoutButtonProps {
    variant?: "default" | "destructive" | "outline" | "secondary" | "ghost" | "link";
    size?: "default" | "sm" | "lg" | "icon";
    className?: string;
    onLogoutSuccess?: () => void;
}

export const LogoutButton = ({
    variant = "outline",
    size = "default",
    className,
    onLogoutSuccess,
}: LogoutButtonProps) => {
    const router = useRouter();
    const [loading, setLoading] = useState(false);

    const handleLogout = async () => {
        setLoading(true);

        try {
            // Call logout API
            await authApi.logout();
        } catch (err) {
            // Even if API fails, clear local state
            console.error("Logout API error:", err);
        } finally {
            // Clear all auth-related cookies and storage
            Cookies.remove("session");
            Cookies.remove("token");
            localStorage.removeItem("access_token");
            localStorage.removeItem("user");

            setLoading(false);

            // Call success callback if provided
            if (onLogoutSuccess) {
                onLogoutSuccess();
            }

            // Redirect to login page
            router.push("/login");
        }
    };

    return (
        <Button
            type="button"
            variant={variant}
            size={size}
            className={className}
            onClick={handleLogout}
            disabled={loading}
            aria-label="Системээс гарах"
        >
            {loading ? (
                <>
                    <span className="animate-spin mr-2" aria-hidden="true">
                        <svg
                            className="w-4 h-4"
                            fill="none"
                            viewBox="0 0 24 24"
                        >
                            <circle
                                className="opacity-25"
                                cx="12"
                                cy="12"
                                r="10"
                                stroke="currentColor"
                                strokeWidth="4"
                            />
                            <path
                                className="opacity-75"
                                fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            />
                        </svg>
                    </span>
                    Гарж байна...
                </>
            ) : (
                <>
                    <svg
                        className="w-4 h-4 mr-2"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        aria-hidden="true"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                        />
                    </svg>
                    Гарах
                </>
            )}
        </Button>
    );
};
