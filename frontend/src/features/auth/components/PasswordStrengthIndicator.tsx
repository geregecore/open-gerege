"use client";

import { useMemo } from "react";
import zxcvbn from "zxcvbn";

interface PasswordStrengthIndicatorProps {
    password: string;
}

const strengthLabels = [
    "Маш сул",
    "Сул",
    "Дунд",
    "Хүчтэй",
    "Маш хүчтэй",
];

const strengthColors = [
    "bg-red-500",
    "bg-orange-500",
    "bg-yellow-500",
    "bg-lime-500",
    "bg-green-500",
];

export const PasswordStrengthIndicator = ({ password }: PasswordStrengthIndicatorProps) => {
    const result = useMemo(() => {
        if (!password) {
            return null;
        }
        return zxcvbn(password);
    }, [password]);

    if (!result || !password) {
        return null;
    }

    const { score } = result;

    return (
        <div className="space-y-2" role="status" aria-live="polite">
            {/* Strength bars */}
            <div className="flex gap-1" aria-hidden="true">
                {[0, 1, 2, 3, 4].map((index) => (
                    <div
                        key={index}
                        className={`h-1 flex-1 rounded-full transition-colors ${
                            index <= score
                                ? strengthColors[score]
                                : "bg-gray-200 dark:bg-gray-700"
                        }`}
                    />
                ))}
            </div>

            {/* Strength label */}
            <div className="flex items-center justify-between text-xs">
                <span
                    className={`font-medium ${
                        score <= 1
                            ? "text-red-600 dark:text-red-400"
                            : score === 2
                            ? "text-yellow-600 dark:text-yellow-400"
                            : "text-green-600 dark:text-green-400"
                    }`}
                >
                    Нууц үгийн хүч: {strengthLabels[score]}
                </span>
            </div>

            {/* Feedback suggestions */}
            {result.feedback.suggestions.length > 0 && (
                <ul className="text-xs text-muted-foreground space-y-1">
                    {result.feedback.suggestions.slice(0, 2).map((suggestion, index) => (
                        <li key={index} className="flex items-start gap-1">
                            <span aria-hidden="true">-</span>
                            <span>{translateSuggestion(suggestion)}</span>
                        </li>
                    ))}
                </ul>
            )}

            {/* Screen reader text */}
            <span className="sr-only">
                Нууц үгийн хүч: {strengthLabels[score]}.
                {score < 3 && " Илүү хүчтэй нууц үг сонгоно уу."}
            </span>
        </div>
    );
};

// Translate common zxcvbn suggestions to Mongolian
function translateSuggestion(suggestion: string): string {
    const translations: Record<string, string> = {
        "Use a few words, avoid common phrases": "Хэд хэдэн үг ашиглаж, түгээмэл хэллэгээс зайлсхий",
        "No need for symbols, digits, or uppercase letters": "Тэмдэг, тоо, том үсэг шаардлагагүй",
        "Add another word or two. Uncommon words are better.": "Нэг хоёр үг нэмнэ үү. Ховор үгс илүү сайн",
        "Straight rows of keys are easy to guess": "Дараалсан товчлуурууд таахад хялбар",
        "Short keyboard patterns are easy to guess": "Богино товчлуурын хээ таахад хялбар",
        "Use a longer keyboard pattern with more turns": "Илүү урт, эргэлттэй товчлуурын хээ ашигла",
        "Repeats like \"aaa\" are easy to guess": "\"aaa\" гэх мэт давталт таахад хялбар",
        "Repeats like \"abcabcabc\" are only slightly harder to guess than \"abc\"": "\"abcabcabc\" гэх мэт давталт \"abc\"-ээс бага зэрэг хүнд",
        "Sequences like \"abc\" or \"6543\" are easy to guess": "\"abc\" эсвэл \"6543\" гэх мэт дараалал таахад хялбар",
        "Recent years are easy to guess": "Сүүлийн жилүүд таахад хялбар",
        "Dates are often easy to guess": "Огноо ихэвчлэн таахад хялбар",
        "This is a top-10 common password": "Энэ нь хамгийн түгээмэл 10 нууц үгийн нэг",
        "This is a top-100 common password": "Энэ нь хамгийн түгээмэл 100 нууц үгийн нэг",
        "This is a very common password": "Энэ нь маш түгээмэл нууц үг",
        "This is similar to a commonly used password": "Энэ нь түгээмэл нууц үгтэй төстэй",
        "A word by itself is easy to guess": "Ганц үг таахад хялбар",
        "Names and surnames by themselves are easy to guess": "Нэр, овог дангаараа таахад хялбар",
        "Common names and surnames are easy to guess": "Түгээмэл нэр, овог таахад хялбар",
        "Capitalization doesn't help very much": "Том үсэг тийм ч их тус болохгүй",
        "All-uppercase is almost as easy to guess as all-lowercase": "Бүгд том үсэг бүгд жижиг үсэгтэй адил таахад хялбар",
        "Reversed words aren't much harder to guess": "Эргүүлсэн үг таахад тийм ч хэцүү биш",
        "Predictable substitutions like '@' instead of 'a' don't help very much": "'@' гэх мэт урьдчилан таамаглах орлуулалт тийм ч их тус болохгүй",
    };

    return translations[suggestion] || suggestion;
}
