"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { authApi } from "../api";
import { loginSchema, type LoginFormData } from "../schemas";
import type { LoginResponse } from "../types";
import Cookies from "js-cookie";

export const LoginForm = () => {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Get redirect URL from query params
    const redirectUrl = searchParams.get("redirect") || "/profile";

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginFormData>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            email: "",
            password: "",
            rememberMe: false,
        },
    });

    const onSubmit = async (data: LoginFormData) => {
        setLoading(true);
        setError(null);

        try {
            const response = await authApi.loginLocal({
                email: data.email,
                password: data.password,
            });

            // The backend returns { code, message, data: { access_token: ... } }
            // The api-client usually returns 'data' if 'code' is present.
            // Extract token from response
            const responseData = response as LoginResponse & { data?: LoginResponse };
            const token = responseData.access_token || responseData.data?.access_token;

            if (token) {
                // Set cookie expiry based on rememberMe
                const expires = data.rememberMe ? 7 : 1; // 7 days or 1 day

                // Set cookie accessible to client
                Cookies.set("session", token, { expires });
                Cookies.set("token", token, { expires });

                // Also store in localStorage for backup
                localStorage.setItem("access_token", token);

                console.log("Login successful, token saved");
                router.push(redirectUrl);
            } else {
                throw new Error("No access token received");
            }
        } catch (err) {
            console.error(err);
            if (err instanceof Error) {
                // Handle specific error messages
                if (err.message.includes("invalid")) {
                    setError("Email эсвэл нууц үг буруу байна");
                } else if (err.message.includes("locked")) {
                    setError("Таны бүртгэл түр түгжигдсэн байна. Дараа дахин оролдоно уу.");
                } else if (err.message.includes("not active")) {
                    setError("Таны бүртгэл идэвхгүй байна. Админтай холбогдоно уу.");
                } else {
                    setError(err.message);
                }
            } else {
                setError("Нэвтрэлт амжилтгүй боллоо");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <form
            onSubmit={handleSubmit(onSubmit)}
            className="space-y-4 w-full max-w-sm"
            noValidate
        >
            {/* Email */}
            <div className="space-y-2">
                <label
                    htmlFor="email"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                    Email
                </label>
                <input
                    id="email"
                    type="email"
                    autoComplete="email"
                    aria-invalid={!!errors.email}
                    aria-describedby={errors.email ? "email-error" : undefined}
                    {...register("email")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    placeholder="name@example.com"
                />
                {errors.email && (
                    <p id="email-error" className="text-sm text-red-500" role="alert">
                        {errors.email.message}
                    </p>
                )}
            </div>

            {/* Password */}
            <div className="space-y-2">
                <div className="flex items-center justify-between">
                    <label
                        htmlFor="password"
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    >
                        Нууц үг
                    </label>
                    <Link
                        href="/forgot-password"
                        className="text-sm text-primary hover:underline"
                    >
                        Нууц үг мартсан?
                    </Link>
                </div>
                <input
                    id="password"
                    type="password"
                    autoComplete="current-password"
                    aria-invalid={!!errors.password}
                    aria-describedby={errors.password ? "password-error" : undefined}
                    {...register("password")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
                {errors.password && (
                    <p id="password-error" className="text-sm text-red-500" role="alert">
                        {errors.password.message}
                    </p>
                )}
            </div>

            {/* Remember Me */}
            <div className="flex items-center space-x-2">
                <input
                    id="rememberMe"
                    type="checkbox"
                    {...register("rememberMe")}
                    className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                />
                <label
                    htmlFor="rememberMe"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                    Намайг сана
                </label>
            </div>

            {/* Error message */}
            {error && (
                <div
                    className="p-3 rounded-md bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800"
                    role="alert"
                >
                    <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
                </div>
            )}

            {/* Submit button */}
            <Button type="submit" disabled={loading} className="w-full">
                {loading ? "Нэвтэрч байна..." : "Нэвтрэх"}
            </Button>

            {/* Register link */}
            <p className="text-center text-sm text-muted-foreground">
                Бүртгэлгүй юу?{" "}
                <Link href="/register" className="underline hover:text-primary">
                    Бүртгүүлэх
                </Link>
            </p>
        </form>
    );
};
