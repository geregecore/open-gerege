"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { authApi } from "@/features/auth/api";

type VerificationStatus = "loading" | "success" | "error";

export default function VerifyEmailPage() {
    const params = useParams();
    const router = useRouter();
    const [status, setStatus] = useState<VerificationStatus>("loading");
    const [errorMessage, setErrorMessage] = useState<string>("");

    const token = params.token as string;

    useEffect(() => {
        // Flag to handle cleanup and prevent state updates after unmount
        let isMounted = true;

        const verifyEmail = async () => {
            if (!token) {
                if (isMounted) {
                    setStatus("error");
                    setErrorMessage("Баталгаажуулах токен олдсонгүй");
                }
                return;
            }

            try {
                await authApi.verifyEmail(token);
                if (isMounted) {
                    setStatus("success");
                    // Redirect to login after 3 seconds
                    setTimeout(() => {
                        if (isMounted) {
                            router.push("/login");
                        }
                    }, 3000);
                }
            } catch (err) {
                console.error("Email verification failed:", err);
                if (isMounted) {
                    setStatus("error");
                    if (err instanceof Error) {
                        if (err.message.includes("invalid") || err.message.includes("expired")) {
                            setErrorMessage("Баталгаажуулах холбоос хүчингүй эсвэл хугацаа дууссан байна");
                        } else {
                            setErrorMessage(err.message);
                        }
                    } else {
                        setErrorMessage("Баталгаажуулалт амжилтгүй боллоо. Дахин оролдоно уу.");
                    }
                }
            }
        };

        verifyEmail();

        return () => {
            isMounted = false;
        };
    }, [token, router]);

    return (
        <div className="flex min-h-screen flex-col items-center justify-center p-24">
            <div className="w-full max-w-sm space-y-8 text-center">
                {status === "loading" && (
                    <div role="status" aria-live="polite">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
                        <h1 className="mt-4 text-xl font-semibold">
                            Email баталгаажуулж байна...
                        </h1>
                        <p className="mt-2 text-sm text-muted-foreground">
                            Түр хүлээнэ үү
                        </p>
                    </div>
                )}

                {status === "success" && (
                    <div role="alert" aria-live="polite">
                        <div className="mx-auto w-12 h-12 bg-green-100 dark:bg-green-900/20 rounded-full flex items-center justify-center">
                            <svg
                                className="w-6 h-6 text-green-600 dark:text-green-400"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                aria-hidden="true"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M5 13l4 4L19 7"
                                />
                            </svg>
                        </div>
                        <h1 className="mt-4 text-xl font-semibold text-green-800 dark:text-green-200">
                            Email баталгаажлаа!
                        </h1>
                        <p className="mt-2 text-sm text-muted-foreground">
                            Таны email амжилттай баталгаажлаа.
                            Та одоо нэвтрэх боломжтой.
                        </p>
                        <p className="mt-4 text-sm text-muted-foreground">
                            Нэвтрэх хуудас руу автоматаар шилжих болно...
                        </p>
                        <div className="mt-6">
                            <Link
                                href="/login"
                                className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                            >
                                Нэвтрэх
                            </Link>
                        </div>
                    </div>
                )}

                {status === "error" && (
                    <div role="alert" aria-live="polite">
                        <div className="mx-auto w-12 h-12 bg-red-100 dark:bg-red-900/20 rounded-full flex items-center justify-center">
                            <svg
                                className="w-6 h-6 text-red-600 dark:text-red-400"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                                aria-hidden="true"
                            >
                                <path
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                    strokeWidth={2}
                                    d="M6 18L18 6M6 6l12 12"
                                />
                            </svg>
                        </div>
                        <h1 className="mt-4 text-xl font-semibold text-red-800 dark:text-red-200">
                            Баталгаажуулалт амжилтгүй
                        </h1>
                        <p className="mt-2 text-sm text-red-600 dark:text-red-400">
                            {errorMessage}
                        </p>
                        <div className="mt-6 space-y-2">
                            <Link
                                href="/login"
                                className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                            >
                                Нэвтрэх хуудас руу очих
                            </Link>
                            <p className="text-sm text-muted-foreground">
                                Асуудал байвал{" "}
                                <Link href="/register" className="underline hover:text-primary">
                                    дахин бүртгүүлэх
                                </Link>
                                {" "}эсвэл{" "}
                                <Link href="/contact" className="underline hover:text-primary">
                                    холбоо барих
                                </Link>
                            </p>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
