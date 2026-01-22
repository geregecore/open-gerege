import { RegisterForm } from "@/features/auth/components/RegisterForm";
import Link from "next/link";

export default function RegisterPage() {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center p-24">
            <div className="w-full max-w-sm space-y-8">
                <div className="flex flex-col space-y-2 text-center">
                    <h1 className="text-2xl font-semibold tracking-tight">
                        Бүртгүүлэх
                    </h1>
                    <p className="text-sm text-muted-foreground">
                        Шинэ бүртгэл үүсгэхийн тулд доорх маягтыг бөглөнө үү
                    </p>
                </div>
                <RegisterForm />
                <div className="text-center text-sm">
                    <Link
                        href="/"
                        className="underline text-muted-foreground hover:text-primary"
                    >
                        Нүүр хуудас руу буцах
                    </Link>
                </div>
            </div>
        </div>
    );
}
