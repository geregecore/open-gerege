"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { authApi } from "../api";
import { registerSchema, type RegisterFormData } from "../schemas";
import { PasswordStrengthIndicator } from "./PasswordStrengthIndicator";

export const RegisterForm = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState(false);

    const {
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm<RegisterFormData>({
        resolver: zodResolver(registerSchema),
        defaultValues: {
            email: "",
            password: "",
            confirmPassword: "",
            firstName: "",
            lastName: "",
            acceptTerms: false as unknown as true,
        },
    });

    const password = watch("password");

    const onSubmit = async (data: RegisterFormData) => {
        setLoading(true);
        setError(null);

        try {
            await authApi.register({
                email: data.email,
                password: data.password,
                confirmPassword: data.confirmPassword,
                firstName: data.firstName,
                lastName: data.lastName,
                acceptTerms: data.acceptTerms,
            });

            setSuccess(true);
        } catch (err) {
            console.error(err);
            if (err instanceof Error) {
                // Handle specific error messages
                if (err.message.includes("email already registered")) {
                    setError("Энэ email хаяг бүртгэлтэй байна");
                } else {
                    setError(err.message);
                }
            } else {
                setError("Бүртгэл амжилтгүй боллоо. Дахин оролдоно уу.");
            }
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div
                className="space-y-4 w-full max-w-sm text-center"
                role="alert"
                aria-live="polite"
            >
                <div className="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg">
                    <h2 className="text-lg font-semibold text-green-800 dark:text-green-200">
                        Бүртгэл амжилттай!
                    </h2>
                    <p className="mt-2 text-sm text-green-700 dark:text-green-300">
                        Таны email хаяг руу баталгаажуулах холбоос илгээгдлээ.
                        Email-ээ шалгаж, холбоос дээр дарж бүртгэлээ баталгаажуулна уу.
                    </p>
                </div>
                <Link
                    href="/login"
                    className="inline-block text-sm text-primary hover:underline"
                >
                    Нэвтрэх хуудас руу буцах
                </Link>
            </div>
        );
    }

    return (
        <form
            onSubmit={handleSubmit(onSubmit)}
            className="space-y-4 w-full max-w-sm"
            noValidate
        >
            {/* First Name */}
            <div className="space-y-2">
                <label
                    htmlFor="firstName"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                    Нэр
                </label>
                <input
                    id="firstName"
                    type="text"
                    autoComplete="given-name"
                    aria-invalid={!!errors.firstName}
                    aria-describedby={errors.firstName ? "firstName-error" : undefined}
                    {...register("firstName")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    placeholder="Баяр"
                />
                {errors.firstName && (
                    <p id="firstName-error" className="text-sm text-red-500" role="alert">
                        {errors.firstName.message}
                    </p>
                )}
            </div>

            {/* Last Name */}
            <div className="space-y-2">
                <label
                    htmlFor="lastName"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                    Овог
                </label>
                <input
                    id="lastName"
                    type="text"
                    autoComplete="family-name"
                    aria-invalid={!!errors.lastName}
                    aria-describedby={errors.lastName ? "lastName-error" : undefined}
                    {...register("lastName")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    placeholder="Дорж"
                />
                {errors.lastName && (
                    <p id="lastName-error" className="text-sm text-red-500" role="alert">
                        {errors.lastName.message}
                    </p>
                )}
            </div>

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
                <label
                    htmlFor="password"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                    Нууц үг
                </label>
                <input
                    id="password"
                    type="password"
                    autoComplete="new-password"
                    aria-invalid={!!errors.password}
                    aria-describedby={errors.password ? "password-error" : "password-strength"}
                    {...register("password")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
                {errors.password && (
                    <p id="password-error" className="text-sm text-red-500" role="alert">
                        {errors.password.message}
                    </p>
                )}
                <div id="password-strength">
                    <PasswordStrengthIndicator password={password} />
                </div>
            </div>

            {/* Confirm Password */}
            <div className="space-y-2">
                <label
                    htmlFor="confirmPassword"
                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                    Нууц үг баталгаажуулах
                </label>
                <input
                    id="confirmPassword"
                    type="password"
                    autoComplete="new-password"
                    aria-invalid={!!errors.confirmPassword}
                    aria-describedby={errors.confirmPassword ? "confirmPassword-error" : undefined}
                    {...register("confirmPassword")}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                />
                {errors.confirmPassword && (
                    <p id="confirmPassword-error" className="text-sm text-red-500" role="alert">
                        {errors.confirmPassword.message}
                    </p>
                )}
            </div>

            {/* Terms and Conditions */}
            <div className="flex items-start space-x-2">
                <input
                    id="acceptTerms"
                    type="checkbox"
                    aria-invalid={!!errors.acceptTerms}
                    aria-describedby={errors.acceptTerms ? "acceptTerms-error" : undefined}
                    {...register("acceptTerms")}
                    className="mt-1 h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                />
                <div className="grid gap-1.5 leading-none">
                    <label
                        htmlFor="acceptTerms"
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    >
                        Үйлчилгээний нөхцөлийг зөвшөөрч байна
                    </label>
                    <p className="text-sm text-muted-foreground">
                        <Link href="/terms" className="underline hover:text-primary">
                            Үйлчилгээний нөхцөл
                        </Link>
                        {" "}болон{" "}
                        <Link href="/privacy" className="underline hover:text-primary">
                            Нууцлалын бодлого
                        </Link>
                        -ыг уншсан.
                    </p>
                </div>
            </div>
            {errors.acceptTerms && (
                <p id="acceptTerms-error" className="text-sm text-red-500" role="alert">
                    {errors.acceptTerms.message}
                </p>
            )}

            {/* Error message */}
            {error && (
                <p className="text-sm text-red-500" role="alert">
                    {error}
                </p>
            )}

            {/* Submit button */}
            <Button type="submit" disabled={loading} className="w-full">
                {loading ? "Бүртгэж байна..." : "Бүртгүүлэх"}
            </Button>

            {/* Login link */}
            <p className="text-center text-sm text-muted-foreground">
                Бүртгэлтэй юу?{" "}
                <Link href="/login" className="underline hover:text-primary">
                    Нэвтрэх
                </Link>
            </p>
        </form>
    );
};
